package cosem

import (
	"fmt"

	"github.com/ViewWay/dlms-cosem-go/core"
)

// ScheduleEntry represents a single schedule entry
type ScheduleEntry struct {
	Index          uint16
	Enable         bool
	ScriptSelector uint8
	StartTime      core.CosemTime
	ValidForDays   uint8
}

// Schedule is the COSEM Schedule interface class (IC 9).
type Schedule struct {
	LogicalName core.ObisCode
	Entries     []ScheduleEntry
	Version     uint8
}

func (s *Schedule) ClassID() uint16               { return core.ClassIDSchedule }
func (s *Schedule) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *Schedule) GetVersion() uint8             { return s.Version }
func (s *Schedule) Access(attr int) core.CosemAttributeAccess {
	switch attr {
	case 1, 2:
		return core.CosemAttributeAccess{Access: core.AccessReadWrite}
	default:
		return core.CosemAttributeAccess{Access: core.AccessNone}
	}
}

func (s *Schedule) MarshalBinary() ([]byte, error) {
	arr := core.ArrayData{}
	for _, e := range s.Entries {
		arr = append(arr, core.StructureData{
			core.UnsignedIntegerData(uint8(e.Index)),
			core.BooleanData(e.Enable),
			core.UnsignedIntegerData(e.ScriptSelector),
			core.TimeData{Value: e.StartTime},
			core.UnsignedIntegerData(e.ValidForDays),
		})
	}
	return arr.ToBytes(), nil
}

func (s *Schedule) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	arr, ok := elem.(core.ArrayData)
	if !ok {
		return fmt.Errorf("schedule: expected array")
	}
	s.Entries = make([]ScheduleEntry, len(arr))
	for i, item := range arr {
		st, ok := item.(core.StructureData)
		if !ok || len(st) < 5 {
			return fmt.Errorf("schedule: invalid entry at index %d", i)
		}
		s.Entries[i].Index = uint16(st[0].(core.UnsignedIntegerData))
		s.Entries[i].Enable = bool(st[1].(core.BooleanData))
		s.Entries[i].ScriptSelector = uint8(st[2].(core.UnsignedIntegerData))
		s.Entries[i].StartTime = st[3].(core.TimeData).Value
		s.Entries[i].ValidForDays = uint8(st[4].(core.UnsignedIntegerData))
	}
	return nil
}

func (s *Schedule) AddEntry(entry ScheduleEntry) {
	s.Entries = append(s.Entries, entry)
}

func (s *Schedule) EntryCount() int {
	return len(s.Entries)
}

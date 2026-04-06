package cosem

import (
	"fmt"

	"github.com/ViewWay/dlms-cosem-go/core"
)

// ActionScheduleEntry represents an action schedule entry
type ActionScheduleEntry struct {
	ExecutedScriptLN core.ObisCode
	ExecutedAt       core.CosemDateTime
}

// SingleActionSchedule is the COSEM Single Action Schedule interface class (IC 71).
type SingleActionSchedule struct {
	LogicalName    core.ObisCode
	Entries        []ActionScheduleEntry
	ExecutedScript core.ObisCode
	Type           uint8
	Version        uint8
}

func (sas *SingleActionSchedule) ClassID() uint16               { return core.ClassIDSingleActionSchedule }
func (sas *SingleActionSchedule) GetLogicalName() core.ObisCode { return sas.LogicalName }
func (sas *SingleActionSchedule) GetVersion() uint8             { return sas.Version }
func (sas *SingleActionSchedule) Access(attr int) core.CosemAttributeAccess {
	switch attr {
	case 1:
		return core.CosemAttributeAccess{Access: core.AccessRead}
	case 2, 3:
		return core.CosemAttributeAccess{Access: core.AccessReadWrite}
	default:
		return core.CosemAttributeAccess{Access: core.AccessNone}
	}
}

func (sas *SingleActionSchedule) MarshalBinary() ([]byte, error) {
	arr := core.ArrayData{}
	for _, e := range sas.Entries {
		arr = append(arr, core.StructureData{
			core.OctetStringData(e.ExecutedScriptLN[:]),
			core.DateTimeData{Value: e.ExecutedAt},
		})
	}
	return arr.ToBytes(), nil
}

func (sas *SingleActionSchedule) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	arr, ok := elem.(core.ArrayData)
	if !ok {
		return fmt.Errorf("single_action_schedule: expected array")
	}
	sas.Entries = make([]ActionScheduleEntry, len(arr))
	for i, item := range arr {
		st, ok := item.(core.StructureData)
		if !ok || len(st) < 2 {
			return fmt.Errorf("single_action_schedule: invalid entry at index %d", i)
		}
		sas.Entries[i].ExecutedScriptLN = core.ObisCode(st[0].(core.OctetStringData))
		sas.Entries[i].ExecutedAt = st[1].(core.DateTimeData).Value
	}
	return nil
}

func (sas *SingleActionSchedule) AddEntry(entry ActionScheduleEntry) {
	sas.Entries = append(sas.Entries, entry)
}

func (sas *SingleActionSchedule) EntryCount() int {
	return len(sas.Entries)
}

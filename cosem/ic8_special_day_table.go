package cosem

import (
	"fmt"

	"github.com/ViewWay/dlms-cosem-go/core"
)

// SpecialDayEntry represents a special day entry
type SpecialDayEntry struct {
	Index         uint16
	SpecialDay    core.CosemDate
	DayProfileRef string
}

// SpecialDaysTable is the COSEM Special Days Table interface class (IC 8).
type SpecialDaysTable struct {
	LogicalName core.ObisCode
	Entries     []SpecialDayEntry
	Version     uint8
}

func (sdt *SpecialDaysTable) ClassID() uint16               { return core.ClassIDSpecialDaysTable }
func (sdt *SpecialDaysTable) GetLogicalName() core.ObisCode { return sdt.LogicalName }
func (sdt *SpecialDaysTable) GetVersion() uint8             { return sdt.Version }
func (sdt *SpecialDaysTable) Access(attr int) core.CosemAttributeAccess {
	switch attr {
	case 2:
		return core.CosemAttributeAccess{Access: core.AccessReadWrite}
	default:
		return core.CosemAttributeAccess{Access: core.AccessRead}
	}
}

func (sdt *SpecialDaysTable) MarshalBinary() ([]byte, error) {
	arr := core.ArrayData{}
	for _, e := range sdt.Entries {
		arr = append(arr, core.StructureData{
			core.UnsignedLongData(e.Index),
			core.DateData{Value: e.SpecialDay},
			core.VisibleStringData(e.DayProfileRef),
		})
	}
	return arr.ToBytes(), nil
}

func (sdt *SpecialDaysTable) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	arr, ok := elem.(core.ArrayData)
	if !ok {
		return fmt.Errorf("special_days_table: expected array")
	}
	sdt.Entries = make([]SpecialDayEntry, len(arr))
	for i, item := range arr {
		st, ok := item.(core.StructureData)
		if !ok || len(st) < 3 {
			return fmt.Errorf("special_days_table: invalid entry at index %d", i)
		}
		sdt.Entries[i].Index = uint16(st[0].(core.UnsignedLongData))
		sdt.Entries[i].SpecialDay = st[1].(core.DateData).Value
		sdt.Entries[i].DayProfileRef = string(st[2].(core.VisibleStringData))
	}
	return nil
}

func (sdt *SpecialDaysTable) AddEntry(entry SpecialDayEntry) {
	sdt.Entries = append(sdt.Entries, entry)
}

func (sdt *SpecialDaysTable) EntryCount() int {
	return len(sdt.Entries)
}

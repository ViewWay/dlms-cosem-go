package cosem

import (
	"fmt"

	"github.com/ViewWay/dlms-cosem-go/core"
)

// ValueEntry represents a value entry in the table
type ValueEntry struct {
	Index     uint16
	Value     core.DlmsData
	Timestamp *core.CosemDateTime
}

// ValueDescriptor describes a value in the table
type ValueDescriptor struct {
	Index       uint16
	Description string
	Unit        uint8
	Scaler      int8
}

// ValueTable is the COSEM Value Table interface class.
type ValueTable struct {
	LogicalName core.ObisCode
	Values      []ValueEntry
	Descriptors []ValueDescriptor
	Version     uint8
}

func (vt *ValueTable) ClassID() uint16               { return core.ClassIDValueTable }
func (vt *ValueTable) GetLogicalName() core.ObisCode { return vt.LogicalName }
func (vt *ValueTable) GetVersion() uint8             { return vt.Version }
func (vt *ValueTable) Access(attr int) core.CosemAttributeAccess {
	switch attr {
	case 1, 2, 3:
		return core.CosemAttributeAccess{Access: core.AccessRead}
	default:
		return core.CosemAttributeAccess{Access: core.AccessNone}
	}
}

func (vt *ValueTable) MarshalBinary() ([]byte, error) {
	arr := core.ArrayData{}
	for _, v := range vt.Values {
		var ts core.DlmsData = core.NullData{}
		if v.Timestamp != nil {
			ts = core.DateTimeData{Value: *v.Timestamp}
		}
		arr = append(arr, core.StructureData{
			core.LongData(int16(v.Index)),
			v.Value,
			ts,
		})
	}
	return arr.ToBytes(), nil
}

func (vt *ValueTable) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	arr, ok := elem.(core.ArrayData)
	if !ok {
		return fmt.Errorf("value_table: expected array")
	}
	vt.Values = make([]ValueEntry, len(arr))
	for i, item := range arr {
		st, ok := item.(core.StructureData)
		if !ok || len(st) < 3 {
			return fmt.Errorf("value_table: invalid entry at index %d", i)
		}
		vt.Values[i].Index = uint16(st[0].(core.LongData))
		vt.Values[i].Value = st[1]
		if _, isNull := st[2].(core.NullData); !isNull {
			if dt, ok := st[2].(core.DateTimeData); ok {
				vt.Values[i].Timestamp = &dt.Value
			}
		}
	}
	return nil
}

func (vt *ValueTable) AddValue(entry ValueEntry) {
	vt.Values = append(vt.Values, entry)
}

func (vt *ValueTable) AddDescriptor(desc ValueDescriptor) {
	vt.Descriptors = append(vt.Descriptors, desc)
}

func (vt *ValueTable) ValueCount() int { return len(vt.Values) }
func (vt *ValueTable) DescriptorCount() int { return len(vt.Descriptors) }

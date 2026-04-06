package cosem

import (
	"fmt"

	"github.com/ViewWay/dlms-cosem-go/core"
)

// StatusMappingEntry represents a status bit mapping entry.
type StatusMappingEntry struct {
	BitPosition uint8
	Description string
}

// StatusMapping is the COSEM Status Mapping interface class (IC 55).
type StatusMapping struct {
	LogicalName   core.ObisCode
	StatusWord    uint32
	MappingTable  []StatusMappingEntry
	Version       uint8
}

func (sm *StatusMapping) ClassID() uint16               { return core.ClassIDStatusMapping }
func (sm *StatusMapping) GetLogicalName() core.ObisCode { return sm.LogicalName }
func (sm *StatusMapping) GetVersion() uint8             { return sm.Version }
func (sm *StatusMapping) Access(attr int) core.CosemAttributeAccess {
	switch attr {
	case 1, 2:
		return core.CosemAttributeAccess{Access: core.AccessRead}
	case 3:
		return core.CosemAttributeAccess{Access: core.AccessReadWrite}
	default:
		return core.CosemAttributeAccess{Access: core.AccessNone}
	}
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (sm *StatusMapping) MarshalBinary() ([]byte, error) {
	return core.DoubleLongUnsignedData(sm.StatusWord).ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (sm *StatusMapping) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	if v, ok := elem.(core.DoubleLongUnsignedData); ok {
		sm.StatusWord = uint32(v)
	}
	return nil
}

// GetBitStatus checks if a specific bit is set.
func (sm *StatusMapping) GetBitStatus(bit uint8) bool {
	if bit < 32 {
		return (sm.StatusWord & (1 << bit)) != 0
	}
	return false
}

// SetBit sets or clears a specific bit.
func (sm *StatusMapping) SetBit(bit uint8, value bool) {
	if bit < 32 {
		if value {
			sm.StatusWord |= 1 << bit
		} else {
			sm.StatusWord &= ^(1 << bit)
		}
	}
}

// AddMapping adds a status mapping entry.
func (sm *StatusMapping) AddMapping(entry StatusMappingEntry) {
	sm.MappingTable = append(sm.MappingTable, entry)
}

// MarshalMappingTable encodes the mapping table as DLMS data.
func (sm *StatusMapping) MarshalMappingTable() ([]byte, error) {
	arr := core.ArrayData{}
	for _, m := range sm.MappingTable {
		arr = append(arr, core.StructureData{
			core.UnsignedIntegerData(m.BitPosition),
			core.VisibleStringData(m.Description),
		})
	}
	return arr.ToBytes(), nil
}

// UnmarshalMappingTable decodes the mapping table from DLMS data.
func (sm *StatusMapping) UnmarshalMappingTable(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	arr, ok := elem.(core.ArrayData)
	if !ok {
		return fmt.Errorf("status_mapping: expected array for mapping table")
	}
	sm.MappingTable = make([]StatusMappingEntry, len(arr))
	for i, item := range arr {
		st, ok := item.(core.StructureData)
		if !ok || len(st) < 2 {
			continue
		}
		sm.MappingTable[i].BitPosition = uint8(st[0].(core.UnsignedIntegerData))
		if v, ok := st[1].(core.VisibleStringData); ok {
			sm.MappingTable[i].Description = string(v)
		}
	}
	return nil
}

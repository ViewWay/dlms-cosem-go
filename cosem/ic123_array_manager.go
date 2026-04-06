package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// ArrayManager is the COSEM ArrayManager interface class (IC 123).
type ArrayManager struct {
	LogicalName core.ObisCode
	Version     uint8
	ManagedObjectList                   []byte
	AllocationSize                      uint32
	AllocatedData                       []byte
}

func (s *ArrayManager) ClassID() uint16               { return core.ClassIDArrayManager }
func (s *ArrayManager) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *ArrayManager) GetVersion() uint8             { return s.Version }

func (s *ArrayManager) Access(attr int) core.CosemAttributeAccess {
	switch attr {
	case 1:
		return core.CosemAttributeAccess{Access: core.AccessRead}
	case 2:
		return core.CosemAttributeAccess{Access: core.AccessReadWrite}
	default:
		return core.CosemAttributeAccess{Access: core.AccessNone}
	}
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s *ArrayManager) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.OctetStringData(s.ManagedObjectList),
		core.DoubleLongUnsignedData(s.AllocationSize),
		core.OctetStringData(s.AllocatedData),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ArrayManager) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return nil
	}
	if 0 < len(st) {
		if v, ok := st[0].(core.OctetStringData); ok { s.ManagedObjectList = []byte(v) }
	}
	if 1 < len(st) {
		if v, ok := st[1].(core.DoubleLongUnsignedData); ok { s.AllocationSize = uint32(v) }
	}
	if 2 < len(st) {
		if v, ok := st[2].(core.OctetStringData); ok { s.AllocatedData = []byte(v) }
	}
	return nil
}


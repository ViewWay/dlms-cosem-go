package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// RegisterTable is the COSEM RegisterTable interface class (IC 61).
type RegisterTable struct {
	LogicalName core.ObisCode
	Version     uint8
	TableCellValues                     []byte
}

func (s *RegisterTable) ClassID() uint16               { return core.ClassIDRegisterTable }
func (s *RegisterTable) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *RegisterTable) GetVersion() uint8             { return s.Version }

func (s *RegisterTable) Access(attr int) core.CosemAttributeAccess {
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
func (s *RegisterTable) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.OctetStringData(s.TableCellValues),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *RegisterTable) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return nil
	}
	if 0 < len(st) {
		if v, ok := st[0].(core.OctetStringData); ok { s.TableCellValues = []byte(v) }
	}
	return nil
}


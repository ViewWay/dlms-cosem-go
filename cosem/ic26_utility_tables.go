package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// UtilityTables is the COSEM UtilityTables interface class (IC 26).
type UtilityTables struct {
	LogicalName core.ObisCode
	Version     uint8
	TableCellValues                     []byte
}

func (s *UtilityTables) ClassID() uint16               { return core.ClassIDUtilityTables }
func (s *UtilityTables) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *UtilityTables) GetVersion() uint8             { return s.Version }

func (s *UtilityTables) Access(attr int) core.CosemAttributeAccess {
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
func (s *UtilityTables) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.OctetStringData(s.TableCellValues),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *UtilityTables) UnmarshalBinary(data []byte) error {
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


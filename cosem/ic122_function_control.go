package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// FunctionControl is the COSEM FunctionControl interface class (IC 122).
type FunctionControl struct {
	LogicalName core.ObisCode
	Version     uint8
	FunctionList                        []byte
}

func (s *FunctionControl) ClassID() uint16               { return core.ClassIDFunctionControl }
func (s *FunctionControl) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *FunctionControl) GetVersion() uint8             { return s.Version }

func (s *FunctionControl) Access(attr int) core.CosemAttributeAccess {
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
func (s *FunctionControl) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.OctetStringData(s.FunctionList),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *FunctionControl) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return nil
	}
	if 0 < len(st) {
		if v, ok := st[0].(core.OctetStringData); ok { s.FunctionList = []byte(v) }
	}
	return nil
}


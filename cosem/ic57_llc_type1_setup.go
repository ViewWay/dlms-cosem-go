package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// LLCType1Setup is the COSEM LLCType1Setup interface class (IC 57).
type LLCType1Setup struct {
	LogicalName core.ObisCode
	Version     uint8
	LLCType1Enable                      bool
	LLCType1Parameters                  []byte
}

func (s *LLCType1Setup) ClassID() uint16               { return core.ClassIDLLCType1Setup }
func (s *LLCType1Setup) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *LLCType1Setup) GetVersion() uint8             { return s.Version }

func (s *LLCType1Setup) Access(attr int) core.CosemAttributeAccess {
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
func (s *LLCType1Setup) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.BooleanData(s.LLCType1Enable),
		core.OctetStringData(s.LLCType1Parameters),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *LLCType1Setup) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return nil
	}
	if 0 < len(st) {
		if v, ok := st[0].(core.BooleanData); ok { s.LLCType1Enable = bool(v) }
	}
	if 1 < len(st) {
		if v, ok := st[1].(core.OctetStringData); ok { s.LLCType1Parameters = []byte(v) }
	}
	return nil
}


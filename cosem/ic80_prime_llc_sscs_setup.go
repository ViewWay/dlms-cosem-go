package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// PrimeLLCSSCSSetup is the COSEM PrimeLLCSSCSSetup interface class (IC 80).
type PrimeLLCSSCSSetup struct {
	LogicalName core.ObisCode
	Version     uint8
	SSCSType                            uint8
	SSCSEnable                          bool
	SSCSResponseTime                    uint16
}

func (s *PrimeLLCSSCSSetup) ClassID() uint16               { return core.ClassIDPrimeLLCSSCSSetup }
func (s *PrimeLLCSSCSSetup) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *PrimeLLCSSCSSetup) GetVersion() uint8             { return s.Version }

func (s *PrimeLLCSSCSSetup) Access(attr int) core.CosemAttributeAccess {
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
func (s *PrimeLLCSSCSSetup) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.UnsignedIntegerData(s.SSCSType),
		core.BooleanData(s.SSCSEnable),
		core.UnsignedLongData(s.SSCSResponseTime),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *PrimeLLCSSCSSetup) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return nil
	}
	if 0 < len(st) {
		if v, ok := st[0].(core.UnsignedIntegerData); ok { s.SSCSType = uint8(v) }
	}
	if 1 < len(st) {
		if v, ok := st[1].(core.BooleanData); ok { s.SSCSEnable = bool(v) }
	}
	if 2 < len(st) {
		if v, ok := st[2].(core.UnsignedLongData); ok { s.SSCSResponseTime = uint16(v) }
	}
	return nil
}


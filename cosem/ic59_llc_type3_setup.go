package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// LLCType3Setup is the COSEM LLCType3Setup interface class (IC 59).
type LLCType3Setup struct {
	LogicalName core.ObisCode
	Version     uint8
	LLCType3Enable                      bool
	LLCType3Parameters                  []byte
	LLCType3WindowSize                  uint16
	LLCType3RetryCount                  uint8
}

func (s *LLCType3Setup) ClassID() uint16               { return core.ClassIDLLCType3Setup }
func (s *LLCType3Setup) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *LLCType3Setup) GetVersion() uint8             { return s.Version }

func (s *LLCType3Setup) Access(attr int) core.CosemAttributeAccess {
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
func (s *LLCType3Setup) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.BooleanData(s.LLCType3Enable),
		core.OctetStringData(s.LLCType3Parameters),
		core.UnsignedLongData(s.LLCType3WindowSize),
		core.UnsignedIntegerData(s.LLCType3RetryCount),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *LLCType3Setup) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return nil
	}
	if 0 < len(st) {
		if v, ok := st[0].(core.BooleanData); ok { s.LLCType3Enable = bool(v) }
	}
	if 1 < len(st) {
		if v, ok := st[1].(core.OctetStringData); ok { s.LLCType3Parameters = []byte(v) }
	}
	if 2 < len(st) {
		if v, ok := st[2].(core.UnsignedLongData); ok { s.LLCType3WindowSize = uint16(v) }
	}
	if 3 < len(st) {
		if v, ok := st[3].(core.UnsignedIntegerData); ok { s.LLCType3RetryCount = uint8(v) }
	}
	return nil
}


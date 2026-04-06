package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// LLCType2Setup is the COSEM LLCType2Setup interface class (IC 58).
type LLCType2Setup struct {
	LogicalName core.ObisCode
	Version     uint8
	LLCType2Enable                      bool
	LLCType2Parameters                  []byte
	LLCType2WindowSize                  uint16
	LLCType2RetryCount                  uint8
}

func (s *LLCType2Setup) ClassID() uint16               { return core.ClassIDLLCType2Setup }
func (s *LLCType2Setup) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *LLCType2Setup) GetVersion() uint8             { return s.Version }

func (s *LLCType2Setup) Access(attr int) core.CosemAttributeAccess {
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
func (s *LLCType2Setup) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.BooleanData(s.LLCType2Enable),
		core.OctetStringData(s.LLCType2Parameters),
		core.UnsignedLongData(s.LLCType2WindowSize),
		core.UnsignedIntegerData(s.LLCType2RetryCount),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *LLCType2Setup) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return nil
	}
	if 0 < len(st) {
		if v, ok := st[0].(core.BooleanData); ok { s.LLCType2Enable = bool(v) }
	}
	if 1 < len(st) {
		if v, ok := st[1].(core.OctetStringData); ok { s.LLCType2Parameters = []byte(v) }
	}
	if 2 < len(st) {
		if v, ok := st[2].(core.UnsignedLongData); ok { s.LLCType2WindowSize = uint16(v) }
	}
	if 3 < len(st) {
		if v, ok := st[3].(core.UnsignedIntegerData); ok { s.LLCType2RetryCount = uint8(v) }
	}
	return nil
}


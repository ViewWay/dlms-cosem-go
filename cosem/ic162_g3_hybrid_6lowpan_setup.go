package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// G3Hybrid6LoWPANSetup is the COSEM G3Hybrid6LoWPANSetup interface class (IC 162).
type G3Hybrid6LoWPANSetup struct {
	LogicalName core.ObisCode
	Version     uint8
	SixLoWPANEnable                     bool
	SixLoWPANMTU                        uint16
	SixLoWPANFragmentationTimeout       uint16
	SixLoWPANFragmentationRetries       uint8
}

func (s *G3Hybrid6LoWPANSetup) ClassID() uint16               { return core.ClassIDG3Hybrid6LoWPANSetup }
func (s *G3Hybrid6LoWPANSetup) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *G3Hybrid6LoWPANSetup) GetVersion() uint8             { return s.Version }

func (s *G3Hybrid6LoWPANSetup) Access(attr int) core.CosemAttributeAccess {
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
func (s *G3Hybrid6LoWPANSetup) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.BooleanData(s.SixLoWPANEnable),
		core.UnsignedLongData(s.SixLoWPANMTU),
		core.UnsignedLongData(s.SixLoWPANFragmentationTimeout),
		core.UnsignedIntegerData(s.SixLoWPANFragmentationRetries),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *G3Hybrid6LoWPANSetup) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return nil
	}
	if 0 < len(st) {
		if v, ok := st[0].(core.BooleanData); ok { s.SixLoWPANEnable = bool(v) }
	}
	if 1 < len(st) {
		if v, ok := st[1].(core.UnsignedLongData); ok { s.SixLoWPANMTU = uint16(v) }
	}
	if 2 < len(st) {
		if v, ok := st[2].(core.UnsignedLongData); ok { s.SixLoWPANFragmentationTimeout = uint16(v) }
	}
	if 3 < len(st) {
		if v, ok := st[3].(core.UnsignedIntegerData); ok { s.SixLoWPANFragmentationRetries = uint8(v) }
	}
	return nil
}


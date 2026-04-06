package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// HSPLCHDLCSSASSetup is the COSEM HSPLCHDLCSSASSetup interface class (IC 143).
type HSPLCHDLCSSASSetup struct {
	LogicalName core.ObisCode
	Version     uint8
	HDLCSSASEnable                      bool
	HDLCSSASMTU                         uint16
	HDLCSSASFragmentationTimeout        uint16
}

func (s *HSPLCHDLCSSASSetup) ClassID() uint16               { return core.ClassIDHSPLCHDLCSSASSetup }
func (s *HSPLCHDLCSSASSetup) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *HSPLCHDLCSSASSetup) GetVersion() uint8             { return s.Version }

func (s *HSPLCHDLCSSASSetup) Access(attr int) core.CosemAttributeAccess {
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
func (s *HSPLCHDLCSSASSetup) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.BooleanData(s.HDLCSSASEnable),
		core.UnsignedLongData(s.HDLCSSASMTU),
		core.UnsignedLongData(s.HDLCSSASFragmentationTimeout),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *HSPLCHDLCSSASSetup) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return nil
	}
	if 0 < len(st) {
		if v, ok := st[0].(core.BooleanData); ok { s.HDLCSSASEnable = bool(v) }
	}
	if 1 < len(st) {
		if v, ok := st[1].(core.UnsignedLongData); ok { s.HDLCSSASMTU = uint16(v) }
	}
	if 2 < len(st) {
		if v, ok := st[2].(core.UnsignedLongData); ok { s.HDLCSSASFragmentationTimeout = uint16(v) }
	}
	return nil
}


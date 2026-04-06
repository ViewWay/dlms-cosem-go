package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// HSPLCIPSSASSetup is the COSEM HSPLCIPSSASSetup interface class (IC 142).
type HSPLCIPSSASSetup struct {
	LogicalName core.ObisCode
	Version     uint8
	IPSSASEnable                        bool
	IPSSASMTU                           uint16
	IPSSASFragmentationTimeout          uint16
}

func (s *HSPLCIPSSASSetup) ClassID() uint16               { return core.ClassIDHSPLCIPSSASSetup }
func (s *HSPLCIPSSASSetup) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *HSPLCIPSSASSetup) GetVersion() uint8             { return s.Version }

func (s *HSPLCIPSSASSetup) Access(attr int) core.CosemAttributeAccess {
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
func (s *HSPLCIPSSASSetup) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.BooleanData(s.IPSSASEnable),
		core.UnsignedLongData(s.IPSSASMTU),
		core.UnsignedLongData(s.IPSSASFragmentationTimeout),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *HSPLCIPSSASSetup) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return nil
	}
	if 0 < len(st) {
		if v, ok := st[0].(core.BooleanData); ok { s.IPSSASEnable = bool(v) }
	}
	if 1 < len(st) {
		if v, ok := st[1].(core.UnsignedLongData); ok { s.IPSSASMTU = uint16(v) }
	}
	if 2 < len(st) {
		if v, ok := st[2].(core.UnsignedLongData); ok { s.IPSSASFragmentationTimeout = uint16(v) }
	}
	return nil
}


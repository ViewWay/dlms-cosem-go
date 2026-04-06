package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// HSPLCCPASSetup is the COSEM HSPLCCPASSetup interface class (IC 141).
type HSPLCCPASSetup struct {
	LogicalName core.ObisCode
	Version     uint8
	CPASEnable                          bool
	CPASMTU                             uint16
	CPASFragmentationTimeout            uint16
}

func (s *HSPLCCPASSetup) ClassID() uint16               { return core.ClassIDHSPLCCPASSetup }
func (s *HSPLCCPASSetup) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *HSPLCCPASSetup) GetVersion() uint8             { return s.Version }

func (s *HSPLCCPASSetup) Access(attr int) core.CosemAttributeAccess {
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
func (s *HSPLCCPASSetup) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.BooleanData(s.CPASEnable),
		core.UnsignedLongData(s.CPASMTU),
		core.UnsignedLongData(s.CPASFragmentationTimeout),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *HSPLCCPASSetup) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return nil
	}
	if 0 < len(st) {
		if v, ok := st[0].(core.BooleanData); ok { s.CPASEnable = bool(v) }
	}
	if 1 < len(st) {
		if v, ok := st[1].(core.UnsignedLongData); ok { s.CPASMTU = uint16(v) }
	}
	if 2 < len(st) {
		if v, ok := st[2].(core.UnsignedLongData); ok { s.CPASFragmentationTimeout = uint16(v) }
	}
	return nil
}


package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// PrimeMACSetup is the COSEM PrimeMACSetup interface class (IC 82).
type PrimeMACSetup struct {
	LogicalName core.ObisCode
	Version     uint8
	MACAddress                          [6]byte
	MACFrameCounter                     uint32
	MACKey                              []byte
	MACSwitch                           uint8
	MACSecurityEnabled                  bool
	MACSecurityLevel                    uint8
}

func (s *PrimeMACSetup) ClassID() uint16               { return core.ClassIDPrimeMACSetup }
func (s *PrimeMACSetup) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *PrimeMACSetup) GetVersion() uint8             { return s.Version }

func (s *PrimeMACSetup) Access(attr int) core.CosemAttributeAccess {
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
func (s *PrimeMACSetup) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.OctetStringData(s.MACAddress[:]),
		core.DoubleLongUnsignedData(s.MACFrameCounter),
		core.OctetStringData(s.MACKey),
		core.UnsignedIntegerData(s.MACSwitch),
		core.BooleanData(s.MACSecurityEnabled),
		core.UnsignedIntegerData(s.MACSecurityLevel),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *PrimeMACSetup) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return nil
	}
	if 0 < len(st) {
		if v, ok := st[0].(core.OctetStringData); ok && len(v) == 6 { copy(s.MACAddress[:], v) }
	}
	if 1 < len(st) {
		if v, ok := st[1].(core.DoubleLongUnsignedData); ok { s.MACFrameCounter = uint32(v) }
	}
	if 2 < len(st) {
		if v, ok := st[2].(core.OctetStringData); ok { s.MACKey = []byte(v) }
	}
	if 3 < len(st) {
		if v, ok := st[3].(core.UnsignedIntegerData); ok { s.MACSwitch = uint8(v) }
	}
	if 4 < len(st) {
		if v, ok := st[4].(core.BooleanData); ok { s.MACSecurityEnabled = bool(v) }
	}
	if 5 < len(st) {
		if v, ok := st[5].(core.UnsignedIntegerData); ok { s.MACSecurityLevel = uint8(v) }
	}
	return nil
}


package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// ZigBeeTunnelSetup is the COSEM ZigBeeTunnelSetup interface class (IC 105).
type ZigBeeTunnelSetup struct {
	LogicalName core.ObisCode
	Version     uint8
	TunnelAddress                       []byte
	TunnelPort                          uint16
}

func (s *ZigBeeTunnelSetup) ClassID() uint16               { return core.ClassIDZigBeeTunnelSetup }
func (s *ZigBeeTunnelSetup) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *ZigBeeTunnelSetup) GetVersion() uint8             { return s.Version }

func (s *ZigBeeTunnelSetup) Access(attr int) core.CosemAttributeAccess {
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
func (s *ZigBeeTunnelSetup) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.OctetStringData(s.TunnelAddress),
		core.UnsignedLongData(s.TunnelPort),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ZigBeeTunnelSetup) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return nil
	}
	if 0 < len(st) {
		if v, ok := st[0].(core.OctetStringData); ok { s.TunnelAddress = []byte(v) }
	}
	if 1 < len(st) {
		if v, ok := st[1].(core.UnsignedLongData); ok { s.TunnelPort = uint16(v) }
	}
	return nil
}


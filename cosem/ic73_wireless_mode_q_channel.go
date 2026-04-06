package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// WirelessModeQChannel is the COSEM WirelessModeQChannel interface class (IC 73).
type WirelessModeQChannel struct {
	LogicalName core.ObisCode
	Version     uint8
	ChannelInfo                         []byte
}

func (s *WirelessModeQChannel) ClassID() uint16               { return core.ClassIDWirelessModeQChannel }
func (s *WirelessModeQChannel) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *WirelessModeQChannel) GetVersion() uint8             { return s.Version }

func (s *WirelessModeQChannel) Access(attr int) core.CosemAttributeAccess {
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
func (s *WirelessModeQChannel) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.OctetStringData(s.ChannelInfo),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *WirelessModeQChannel) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return nil
	}
	if 0 < len(st) {
		if v, ok := st[0].(core.OctetStringData); ok { s.ChannelInfo = []byte(v) }
	}
	return nil
}


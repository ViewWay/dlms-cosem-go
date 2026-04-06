package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// ZigBeeNetworkControl is the COSEM ZigBeeNetworkControl interface class (IC 104).
type ZigBeeNetworkControl struct {
	LogicalName core.ObisCode
	Version     uint8
	NetworkMode                         uint8
	PANID                               uint16
	ExtendedPANID                       [8]byte
	Channel                             uint8
	PermitDuration                      uint16
	DeviceTimeout                       uint16
	RouterCapacity                      uint8
	EndDeviceCapacity                   uint8
	TrustCenterAddress                  [8]byte
	TrustCenterMasterKey                []byte
	ActiveNetworkKeySeqNumber           uint8
	NetworkKey                          []byte
	LinkKey                             []byte
}

func (s *ZigBeeNetworkControl) ClassID() uint16               { return core.ClassIDZigBeeNetworkControl }
func (s *ZigBeeNetworkControl) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *ZigBeeNetworkControl) GetVersion() uint8             { return s.Version }

func (s *ZigBeeNetworkControl) Access(attr int) core.CosemAttributeAccess {
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
func (s *ZigBeeNetworkControl) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.UnsignedIntegerData(s.NetworkMode),
		core.UnsignedLongData(s.PANID),
		core.OctetStringData(s.ExtendedPANID[:]),
		core.UnsignedIntegerData(s.Channel),
		core.UnsignedLongData(s.PermitDuration),
		core.UnsignedLongData(s.DeviceTimeout),
		core.UnsignedIntegerData(s.RouterCapacity),
		core.UnsignedIntegerData(s.EndDeviceCapacity),
		core.OctetStringData(s.TrustCenterAddress[:]),
		core.OctetStringData(s.TrustCenterMasterKey),
		core.UnsignedIntegerData(s.ActiveNetworkKeySeqNumber),
		core.OctetStringData(s.NetworkKey),
		core.OctetStringData(s.LinkKey),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ZigBeeNetworkControl) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return nil
	}
	if 0 < len(st) {
		if v, ok := st[0].(core.UnsignedIntegerData); ok { s.NetworkMode = uint8(v) }
	}
	if 1 < len(st) {
		if v, ok := st[1].(core.UnsignedLongData); ok { s.PANID = uint16(v) }
	}
	if 2 < len(st) {
		if v, ok := st[2].(core.OctetStringData); ok && len(v) == 8 { copy(s.ExtendedPANID[:], v) }
	}
	if 3 < len(st) {
		if v, ok := st[3].(core.UnsignedIntegerData); ok { s.Channel = uint8(v) }
	}
	if 4 < len(st) {
		if v, ok := st[4].(core.UnsignedLongData); ok { s.PermitDuration = uint16(v) }
	}
	if 5 < len(st) {
		if v, ok := st[5].(core.UnsignedLongData); ok { s.DeviceTimeout = uint16(v) }
	}
	if 6 < len(st) {
		if v, ok := st[6].(core.UnsignedIntegerData); ok { s.RouterCapacity = uint8(v) }
	}
	if 7 < len(st) {
		if v, ok := st[7].(core.UnsignedIntegerData); ok { s.EndDeviceCapacity = uint8(v) }
	}
	if 8 < len(st) {
		if v, ok := st[8].(core.OctetStringData); ok && len(v) == 8 { copy(s.TrustCenterAddress[:], v) }
	}
	if 9 < len(st) {
		if v, ok := st[9].(core.OctetStringData); ok { s.TrustCenterMasterKey = []byte(v) }
	}
	if 10 < len(st) {
		if v, ok := st[10].(core.UnsignedIntegerData); ok { s.ActiveNetworkKeySeqNumber = uint8(v) }
	}
	if 11 < len(st) {
		if v, ok := st[11].(core.OctetStringData); ok { s.NetworkKey = []byte(v) }
	}
	if 12 < len(st) {
		if v, ok := st[12].(core.OctetStringData); ok { s.LinkKey = []byte(v) }
	}
	return nil
}


package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// WiSUNSetup is the COSEM WiSUNSetup interface class (IC 95).
type WiSUNSetup struct {
	LogicalName core.ObisCode
	Version     uint8
	PHYOperatingMode                    uint8
	NetworkMode                         uint8
	PANID                               uint16
	RoutingMethod                       uint8
	RoutingMethodParameters             []byte
	PHYOperatingModeList                []byte
	ChannelFunction                     uint8
	ChannelHoppingMode                  uint8
	UnicastDwellTime                    uint16
	BroadcastDwellTime                  uint16
	BroadcastInterval                   uint32
	BroadcastSequenceNumber             uint8
	MeshHeaderSequenceNumber            uint8
	RoutingTable                        []byte
	RoutingTableUpdateTime              []byte
}

func (s *WiSUNSetup) ClassID() uint16               { return core.ClassIDWiSUNSetup }
func (s *WiSUNSetup) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *WiSUNSetup) GetVersion() uint8             { return s.Version }

func (s *WiSUNSetup) Access(attr int) core.CosemAttributeAccess {
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
func (s *WiSUNSetup) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.UnsignedIntegerData(s.PHYOperatingMode),
		core.UnsignedIntegerData(s.NetworkMode),
		core.UnsignedLongData(s.PANID),
		core.UnsignedIntegerData(s.RoutingMethod),
		core.OctetStringData(s.RoutingMethodParameters),
		core.OctetStringData(s.PHYOperatingModeList),
		core.UnsignedIntegerData(s.ChannelFunction),
		core.UnsignedIntegerData(s.ChannelHoppingMode),
		core.UnsignedLongData(s.UnicastDwellTime),
		core.UnsignedLongData(s.BroadcastDwellTime),
		core.DoubleLongUnsignedData(s.BroadcastInterval),
		core.UnsignedIntegerData(s.BroadcastSequenceNumber),
		core.UnsignedIntegerData(s.MeshHeaderSequenceNumber),
		core.OctetStringData(s.RoutingTable),
		core.OctetStringData(s.RoutingTableUpdateTime),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *WiSUNSetup) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return nil
	}
	if 0 < len(st) {
		if v, ok := st[0].(core.UnsignedIntegerData); ok { s.PHYOperatingMode = uint8(v) }
	}
	if 1 < len(st) {
		if v, ok := st[1].(core.UnsignedIntegerData); ok { s.NetworkMode = uint8(v) }
	}
	if 2 < len(st) {
		if v, ok := st[2].(core.UnsignedLongData); ok { s.PANID = uint16(v) }
	}
	if 3 < len(st) {
		if v, ok := st[3].(core.UnsignedIntegerData); ok { s.RoutingMethod = uint8(v) }
	}
	if 4 < len(st) {
		if v, ok := st[4].(core.OctetStringData); ok { s.RoutingMethodParameters = []byte(v) }
	}
	if 5 < len(st) {
		if v, ok := st[5].(core.OctetStringData); ok { s.PHYOperatingModeList = []byte(v) }
	}
	if 6 < len(st) {
		if v, ok := st[6].(core.UnsignedIntegerData); ok { s.ChannelFunction = uint8(v) }
	}
	if 7 < len(st) {
		if v, ok := st[7].(core.UnsignedIntegerData); ok { s.ChannelHoppingMode = uint8(v) }
	}
	if 8 < len(st) {
		if v, ok := st[8].(core.UnsignedLongData); ok { s.UnicastDwellTime = uint16(v) }
	}
	if 9 < len(st) {
		if v, ok := st[9].(core.UnsignedLongData); ok { s.BroadcastDwellTime = uint16(v) }
	}
	if 10 < len(st) {
		if v, ok := st[10].(core.DoubleLongUnsignedData); ok { s.BroadcastInterval = uint32(v) }
	}
	if 11 < len(st) {
		if v, ok := st[11].(core.UnsignedIntegerData); ok { s.BroadcastSequenceNumber = uint8(v) }
	}
	if 12 < len(st) {
		if v, ok := st[12].(core.UnsignedIntegerData); ok { s.MeshHeaderSequenceNumber = uint8(v) }
	}
	if 13 < len(st) {
		if v, ok := st[13].(core.OctetStringData); ok { s.RoutingTable = []byte(v) }
	}
	if 14 < len(st) {
		if v, ok := st[14].(core.OctetStringData); ok { s.RoutingTableUpdateTime = []byte(v) }
	}
	return nil
}


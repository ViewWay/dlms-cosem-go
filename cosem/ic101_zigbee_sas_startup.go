package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// ZigBeeSASStartup is the COSEM ZigBeeSASStartup interface class (IC 101).
type ZigBeeSASStartup struct {
	LogicalName core.ObisCode
	Version     uint8
	StartupControl                      uint8
	ChannelMask                         uint32
	ScanDuration                        uint8
	ScanAttempts                        uint8
	ScanAttemptsTimeout                 uint16
	Channel                             uint8
	SecurityLevel                       uint8
	PreconfiguredLinkKey                []byte
	NetworkKey                          []byte
	NetworkKeyEnable                    bool
	UseInsecureJoin                     bool
	PermitDuration                      uint16
	DeviceTimeout                       uint16
}

func (s *ZigBeeSASStartup) ClassID() uint16               { return core.ClassIDZigBeeSASStartup }
func (s *ZigBeeSASStartup) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *ZigBeeSASStartup) GetVersion() uint8             { return s.Version }

func (s *ZigBeeSASStartup) Access(attr int) core.CosemAttributeAccess {
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
func (s *ZigBeeSASStartup) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.UnsignedIntegerData(s.StartupControl),
		core.DoubleLongUnsignedData(s.ChannelMask),
		core.UnsignedIntegerData(s.ScanDuration),
		core.UnsignedIntegerData(s.ScanAttempts),
		core.UnsignedLongData(s.ScanAttemptsTimeout),
		core.UnsignedIntegerData(s.Channel),
		core.UnsignedIntegerData(s.SecurityLevel),
		core.OctetStringData(s.PreconfiguredLinkKey),
		core.OctetStringData(s.NetworkKey),
		core.BooleanData(s.NetworkKeyEnable),
		core.BooleanData(s.UseInsecureJoin),
		core.UnsignedLongData(s.PermitDuration),
		core.UnsignedLongData(s.DeviceTimeout),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ZigBeeSASStartup) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return nil
	}
	if 0 < len(st) {
		if v, ok := st[0].(core.UnsignedIntegerData); ok { s.StartupControl = uint8(v) }
	}
	if 1 < len(st) {
		if v, ok := st[1].(core.DoubleLongUnsignedData); ok { s.ChannelMask = uint32(v) }
	}
	if 2 < len(st) {
		if v, ok := st[2].(core.UnsignedIntegerData); ok { s.ScanDuration = uint8(v) }
	}
	if 3 < len(st) {
		if v, ok := st[3].(core.UnsignedIntegerData); ok { s.ScanAttempts = uint8(v) }
	}
	if 4 < len(st) {
		if v, ok := st[4].(core.UnsignedLongData); ok { s.ScanAttemptsTimeout = uint16(v) }
	}
	if 5 < len(st) {
		if v, ok := st[5].(core.UnsignedIntegerData); ok { s.Channel = uint8(v) }
	}
	if 6 < len(st) {
		if v, ok := st[6].(core.UnsignedIntegerData); ok { s.SecurityLevel = uint8(v) }
	}
	if 7 < len(st) {
		if v, ok := st[7].(core.OctetStringData); ok { s.PreconfiguredLinkKey = []byte(v) }
	}
	if 8 < len(st) {
		if v, ok := st[8].(core.OctetStringData); ok { s.NetworkKey = []byte(v) }
	}
	if 9 < len(st) {
		if v, ok := st[9].(core.BooleanData); ok { s.NetworkKeyEnable = bool(v) }
	}
	if 10 < len(st) {
		if v, ok := st[10].(core.BooleanData); ok { s.UseInsecureJoin = bool(v) }
	}
	if 11 < len(st) {
		if v, ok := st[11].(core.UnsignedLongData); ok { s.PermitDuration = uint16(v) }
	}
	if 12 < len(st) {
		if v, ok := st[12].(core.UnsignedLongData); ok { s.DeviceTimeout = uint16(v) }
	}
	return nil
}


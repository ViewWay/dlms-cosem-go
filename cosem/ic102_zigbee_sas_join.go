package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// ZigBeeSASJoin is the COSEM ZigBeeSASJoin interface class (IC 102).
type ZigBeeSASJoin struct {
	LogicalName core.ObisCode
	Version     uint8
	JoinControl                         uint8
	RejoinInterval                      uint16
	MaxRejoinInterval                   uint16
	SecurityLevel                       uint8
	NetworkKeyEnable                    bool
	PreconfiguredLinkKey                []byte
	TrustCenterAddress                  [8]byte
	TrustCenterMasterKey                []byte
	ActiveNetworkKeySeqNumber           uint8
	LinkKey                             []byte
}

func (s *ZigBeeSASJoin) ClassID() uint16               { return core.ClassIDZigBeeSASJoin }
func (s *ZigBeeSASJoin) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *ZigBeeSASJoin) GetVersion() uint8             { return s.Version }

func (s *ZigBeeSASJoin) Access(attr int) core.CosemAttributeAccess {
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
func (s *ZigBeeSASJoin) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.UnsignedIntegerData(s.JoinControl),
		core.UnsignedLongData(s.RejoinInterval),
		core.UnsignedLongData(s.MaxRejoinInterval),
		core.UnsignedIntegerData(s.SecurityLevel),
		core.BooleanData(s.NetworkKeyEnable),
		core.OctetStringData(s.PreconfiguredLinkKey),
		core.OctetStringData(s.TrustCenterAddress[:]),
		core.OctetStringData(s.TrustCenterMasterKey),
		core.UnsignedIntegerData(s.ActiveNetworkKeySeqNumber),
		core.OctetStringData(s.LinkKey),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ZigBeeSASJoin) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return nil
	}
	if 0 < len(st) {
		if v, ok := st[0].(core.UnsignedIntegerData); ok { s.JoinControl = uint8(v) }
	}
	if 1 < len(st) {
		if v, ok := st[1].(core.UnsignedLongData); ok { s.RejoinInterval = uint16(v) }
	}
	if 2 < len(st) {
		if v, ok := st[2].(core.UnsignedLongData); ok { s.MaxRejoinInterval = uint16(v) }
	}
	if 3 < len(st) {
		if v, ok := st[3].(core.UnsignedIntegerData); ok { s.SecurityLevel = uint8(v) }
	}
	if 4 < len(st) {
		if v, ok := st[4].(core.BooleanData); ok { s.NetworkKeyEnable = bool(v) }
	}
	if 5 < len(st) {
		if v, ok := st[5].(core.OctetStringData); ok { s.PreconfiguredLinkKey = []byte(v) }
	}
	if 6 < len(st) {
		if v, ok := st[6].(core.OctetStringData); ok && len(v) == 8 { copy(s.TrustCenterAddress[:], v) }
	}
	if 7 < len(st) {
		if v, ok := st[7].(core.OctetStringData); ok { s.TrustCenterMasterKey = []byte(v) }
	}
	if 8 < len(st) {
		if v, ok := st[8].(core.UnsignedIntegerData); ok { s.ActiveNetworkKeySeqNumber = uint8(v) }
	}
	if 9 < len(st) {
		if v, ok := st[9].(core.OctetStringData); ok { s.LinkKey = []byte(v) }
	}
	return nil
}


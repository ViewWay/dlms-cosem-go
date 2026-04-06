package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// G3MACSetup is the COSEM G3MACSetup interface class (IC 91).
type G3MACSetup struct {
	LogicalName core.ObisCode
	Version     uint8
	MACAddress                          [6]byte
	MACFrameCounter                     uint32
	MACKey                              []byte
	MACSwitch                           uint8
	MACSecurityEnabled                  bool
	MACSecurityLevel                    uint8
	MACRoutingMode                      uint8
	MACTxPower                          uint8
	MACTxRetries                        uint8
	MACTxAckTimeout                     uint32
	MACTxDataRate                       uint8
	MACTxPowerControl                   bool
}

func (s *G3MACSetup) ClassID() uint16               { return core.ClassIDG3MACSetup }
func (s *G3MACSetup) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *G3MACSetup) GetVersion() uint8             { return s.Version }

func (s *G3MACSetup) Access(attr int) core.CosemAttributeAccess {
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
func (s *G3MACSetup) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.OctetStringData(s.MACAddress[:]),
		core.DoubleLongUnsignedData(s.MACFrameCounter),
		core.OctetStringData(s.MACKey),
		core.UnsignedIntegerData(s.MACSwitch),
		core.BooleanData(s.MACSecurityEnabled),
		core.UnsignedIntegerData(s.MACSecurityLevel),
		core.UnsignedIntegerData(s.MACRoutingMode),
		core.UnsignedIntegerData(s.MACTxPower),
		core.UnsignedIntegerData(s.MACTxRetries),
		core.DoubleLongUnsignedData(s.MACTxAckTimeout),
		core.UnsignedIntegerData(s.MACTxDataRate),
		core.BooleanData(s.MACTxPowerControl),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *G3MACSetup) UnmarshalBinary(data []byte) error {
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
	if 6 < len(st) {
		if v, ok := st[6].(core.UnsignedIntegerData); ok { s.MACRoutingMode = uint8(v) }
	}
	if 7 < len(st) {
		if v, ok := st[7].(core.UnsignedIntegerData); ok { s.MACTxPower = uint8(v) }
	}
	if 8 < len(st) {
		if v, ok := st[8].(core.UnsignedIntegerData); ok { s.MACTxRetries = uint8(v) }
	}
	if 9 < len(st) {
		if v, ok := st[9].(core.DoubleLongUnsignedData); ok { s.MACTxAckTimeout = uint32(v) }
	}
	if 10 < len(st) {
		if v, ok := st[10].(core.UnsignedIntegerData); ok { s.MACTxDataRate = uint8(v) }
	}
	if 11 < len(st) {
		if v, ok := st[11].(core.BooleanData); ok { s.MACTxPowerControl = bool(v) }
	}
	return nil
}


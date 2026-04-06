package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// WiSUNDiagnostic is the COSEM WiSUNDiagnostic interface class (IC 96).
type WiSUNDiagnostic struct {
	LogicalName core.ObisCode
	Version     uint8
	MessagesSent                        uint32
	MessagesReceived                    uint32
	MessagesFailed                      uint32
	MessagesRetransmitted               uint32
	PHYTxTotal                          uint32
	PHYRxTotal                          uint32
	PHYTxError                          uint32
	PHYRxError                          uint32
	MACTxTotal                          uint32
	MACRxTotal                          uint32
	MACTxError                          uint32
	MACRxError                          uint32
}

func (s *WiSUNDiagnostic) ClassID() uint16               { return core.ClassIDWiSUNDiagnostic }
func (s *WiSUNDiagnostic) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *WiSUNDiagnostic) GetVersion() uint8             { return s.Version }

func (s *WiSUNDiagnostic) Access(attr int) core.CosemAttributeAccess {
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
func (s *WiSUNDiagnostic) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.DoubleLongUnsignedData(s.MessagesSent),
		core.DoubleLongUnsignedData(s.MessagesReceived),
		core.DoubleLongUnsignedData(s.MessagesFailed),
		core.DoubleLongUnsignedData(s.MessagesRetransmitted),
		core.DoubleLongUnsignedData(s.PHYTxTotal),
		core.DoubleLongUnsignedData(s.PHYRxTotal),
		core.DoubleLongUnsignedData(s.PHYTxError),
		core.DoubleLongUnsignedData(s.PHYRxError),
		core.DoubleLongUnsignedData(s.MACTxTotal),
		core.DoubleLongUnsignedData(s.MACRxTotal),
		core.DoubleLongUnsignedData(s.MACTxError),
		core.DoubleLongUnsignedData(s.MACRxError),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *WiSUNDiagnostic) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return nil
	}
	if 0 < len(st) {
		if v, ok := st[0].(core.DoubleLongUnsignedData); ok { s.MessagesSent = uint32(v) }
	}
	if 1 < len(st) {
		if v, ok := st[1].(core.DoubleLongUnsignedData); ok { s.MessagesReceived = uint32(v) }
	}
	if 2 < len(st) {
		if v, ok := st[2].(core.DoubleLongUnsignedData); ok { s.MessagesFailed = uint32(v) }
	}
	if 3 < len(st) {
		if v, ok := st[3].(core.DoubleLongUnsignedData); ok { s.MessagesRetransmitted = uint32(v) }
	}
	if 4 < len(st) {
		if v, ok := st[4].(core.DoubleLongUnsignedData); ok { s.PHYTxTotal = uint32(v) }
	}
	if 5 < len(st) {
		if v, ok := st[5].(core.DoubleLongUnsignedData); ok { s.PHYRxTotal = uint32(v) }
	}
	if 6 < len(st) {
		if v, ok := st[6].(core.DoubleLongUnsignedData); ok { s.PHYTxError = uint32(v) }
	}
	if 7 < len(st) {
		if v, ok := st[7].(core.DoubleLongUnsignedData); ok { s.PHYRxError = uint32(v) }
	}
	if 8 < len(st) {
		if v, ok := st[8].(core.DoubleLongUnsignedData); ok { s.MACTxTotal = uint32(v) }
	}
	if 9 < len(st) {
		if v, ok := st[9].(core.DoubleLongUnsignedData); ok { s.MACRxTotal = uint32(v) }
	}
	if 10 < len(st) {
		if v, ok := st[10].(core.DoubleLongUnsignedData); ok { s.MACTxError = uint32(v) }
	}
	if 11 < len(st) {
		if v, ok := st[11].(core.DoubleLongUnsignedData); ok { s.MACRxError = uint32(v) }
	}
	return nil
}


package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// LoRaWANDiagnostic is the COSEM LoRaWANDiagnostic interface class (IC 129).
type LoRaWANDiagnostic struct {
	LogicalName core.ObisCode
	Version     uint8
	MessagesSent                        uint32
	MessagesReceived                    uint32
	MessagesFailed                      uint32
	MessagesRetransmitted               uint32
}

func (s *LoRaWANDiagnostic) ClassID() uint16               { return core.ClassIDLoRaWANDiagnostic }
func (s *LoRaWANDiagnostic) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *LoRaWANDiagnostic) GetVersion() uint8             { return s.Version }

func (s *LoRaWANDiagnostic) Access(attr int) core.CosemAttributeAccess {
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
func (s *LoRaWANDiagnostic) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.DoubleLongUnsignedData(s.MessagesSent),
		core.DoubleLongUnsignedData(s.MessagesReceived),
		core.DoubleLongUnsignedData(s.MessagesFailed),
		core.DoubleLongUnsignedData(s.MessagesRetransmitted),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *LoRaWANDiagnostic) UnmarshalBinary(data []byte) error {
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
	return nil
}


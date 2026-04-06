package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// IEC14908ProtocolStatus is the COSEM IEC14908ProtocolStatus interface class (IC 132).
type IEC14908ProtocolStatus struct {
	LogicalName core.ObisCode
	Version     uint8
	ProtocolStatus                      uint8
	ConnectionStatus                    uint8
	CommunicationStatistics             []byte
}

func (s *IEC14908ProtocolStatus) ClassID() uint16               { return core.ClassIDIEC14908ProtocolStatus }
func (s *IEC14908ProtocolStatus) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *IEC14908ProtocolStatus) GetVersion() uint8             { return s.Version }

func (s *IEC14908ProtocolStatus) Access(attr int) core.CosemAttributeAccess {
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
func (s *IEC14908ProtocolStatus) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.UnsignedIntegerData(s.ProtocolStatus),
		core.UnsignedIntegerData(s.ConnectionStatus),
		core.OctetStringData(s.CommunicationStatistics),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *IEC14908ProtocolStatus) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return nil
	}
	if 0 < len(st) {
		if v, ok := st[0].(core.UnsignedIntegerData); ok { s.ProtocolStatus = uint8(v) }
	}
	if 1 < len(st) {
		if v, ok := st[1].(core.UnsignedIntegerData); ok { s.ConnectionStatus = uint8(v) }
	}
	if 2 < len(st) {
		if v, ok := st[2].(core.OctetStringData); ok { s.CommunicationStatistics = []byte(v) }
	}
	return nil
}


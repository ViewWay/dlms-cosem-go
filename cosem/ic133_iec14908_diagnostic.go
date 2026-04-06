package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// IEC14908Diagnostic is the COSEM IEC14908Diagnostic interface class (IC 133).
type IEC14908Diagnostic struct {
	LogicalName core.ObisCode
	Version     uint8
	MessagesSent                        uint32
	MessagesReceived                    uint32
	MessagesFailed                      uint32
	CRCErrors                           uint32
	Timeouts                            uint32
}

func (s *IEC14908Diagnostic) ClassID() uint16               { return core.ClassIDIEC14908Diagnostic }
func (s *IEC14908Diagnostic) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *IEC14908Diagnostic) GetVersion() uint8             { return s.Version }

func (s *IEC14908Diagnostic) Access(attr int) core.CosemAttributeAccess {
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
func (s *IEC14908Diagnostic) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.DoubleLongUnsignedData(s.MessagesSent),
		core.DoubleLongUnsignedData(s.MessagesReceived),
		core.DoubleLongUnsignedData(s.MessagesFailed),
		core.DoubleLongUnsignedData(s.CRCErrors),
		core.DoubleLongUnsignedData(s.Timeouts),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *IEC14908Diagnostic) UnmarshalBinary(data []byte) error {
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
		if v, ok := st[3].(core.DoubleLongUnsignedData); ok { s.CRCErrors = uint32(v) }
	}
	if 4 < len(st) {
		if v, ok := st[4].(core.DoubleLongUnsignedData); ok { s.Timeouts = uint32(v) }
	}
	return nil
}


package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// IEC14908ProtocolSetup is the COSEM IEC14908ProtocolSetup interface class (IC 131).
type IEC14908ProtocolSetup struct {
	LogicalName core.ObisCode
	Version     uint8
	ProtocolMode                        uint8
	ProtocolVersion                     uint8
	ProtocolParameters                  []byte
}

func (s *IEC14908ProtocolSetup) ClassID() uint16               { return core.ClassIDIEC14908ProtocolSetup }
func (s *IEC14908ProtocolSetup) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *IEC14908ProtocolSetup) GetVersion() uint8             { return s.Version }

func (s *IEC14908ProtocolSetup) Access(attr int) core.CosemAttributeAccess {
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
func (s *IEC14908ProtocolSetup) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.UnsignedIntegerData(s.ProtocolMode),
		core.UnsignedIntegerData(s.ProtocolVersion),
		core.OctetStringData(s.ProtocolParameters),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *IEC14908ProtocolSetup) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return nil
	}
	if 0 < len(st) {
		if v, ok := st[0].(core.UnsignedIntegerData); ok { s.ProtocolMode = uint8(v) }
	}
	if 1 < len(st) {
		if v, ok := st[1].(core.UnsignedIntegerData); ok { s.ProtocolVersion = uint8(v) }
	}
	if 2 < len(st) {
		if v, ok := st[2].(core.OctetStringData); ok { s.ProtocolParameters = []byte(v) }
	}
	return nil
}


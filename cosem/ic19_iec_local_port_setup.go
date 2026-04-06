package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// IECLocalPortSetup is the COSEM IECLocalPortSetup interface class (IC 19).
type IECLocalPortSetup struct {
	LogicalName core.ObisCode
	Version     uint8
	DefaultMode                         uint8
	DefaultBaud                         uint16
	BaudRate                            uint16
	LocalPortState                      uint8
}

func (s *IECLocalPortSetup) ClassID() uint16               { return core.ClassIDIECLocalPortSetup }
func (s *IECLocalPortSetup) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *IECLocalPortSetup) GetVersion() uint8             { return s.Version }

func (s *IECLocalPortSetup) Access(attr int) core.CosemAttributeAccess {
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
func (s *IECLocalPortSetup) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.UnsignedIntegerData(s.DefaultMode),
		core.UnsignedLongData(s.DefaultBaud),
		core.UnsignedLongData(s.BaudRate),
		core.UnsignedIntegerData(s.LocalPortState),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *IECLocalPortSetup) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return nil
	}
	if 0 < len(st) {
		if v, ok := st[0].(core.UnsignedIntegerData); ok { s.DefaultMode = uint8(v) }
	}
	if 1 < len(st) {
		if v, ok := st[1].(core.UnsignedLongData); ok { s.DefaultBaud = uint16(v) }
	}
	if 2 < len(st) {
		if v, ok := st[2].(core.UnsignedLongData); ok { s.BaudRate = uint16(v) }
	}
	if 3 < len(st) {
		if v, ok := st[3].(core.UnsignedIntegerData); ok { s.LocalPortState = uint8(v) }
	}
	return nil
}


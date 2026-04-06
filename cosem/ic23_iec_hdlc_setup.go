package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// IECHDLCSetup is the COSEM IECHDLCSetup interface class (IC 23).
type IECHDLCSetup struct {
	LogicalName core.ObisCode
	Version     uint8
	CommSpeed                           uint32
	WindowSizeTransmit                  uint8
	WindowSizeReceive                   uint8
	MaxInfoFieldLengthTransmit          uint16
	MaxInfoFieldLengthReceive           uint16
	InterOctetTimeOut                   uint16
	InterCharacterTimeOut               uint16
	InactivityTimeOut                   uint16
}

func (s *IECHDLCSetup) ClassID() uint16               { return core.ClassIDIECHDLCSetup }
func (s *IECHDLCSetup) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *IECHDLCSetup) GetVersion() uint8             { return s.Version }

func (s *IECHDLCSetup) Access(attr int) core.CosemAttributeAccess {
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
func (s *IECHDLCSetup) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.DoubleLongUnsignedData(s.CommSpeed),
		core.UnsignedIntegerData(s.WindowSizeTransmit),
		core.UnsignedIntegerData(s.WindowSizeReceive),
		core.UnsignedLongData(s.MaxInfoFieldLengthTransmit),
		core.UnsignedLongData(s.MaxInfoFieldLengthReceive),
		core.UnsignedLongData(s.InterOctetTimeOut),
		core.UnsignedLongData(s.InterCharacterTimeOut),
		core.UnsignedLongData(s.InactivityTimeOut),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *IECHDLCSetup) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return nil
	}
	if 0 < len(st) {
		if v, ok := st[0].(core.DoubleLongUnsignedData); ok { s.CommSpeed = uint32(v) }
	}
	if 1 < len(st) {
		if v, ok := st[1].(core.UnsignedIntegerData); ok { s.WindowSizeTransmit = uint8(v) }
	}
	if 2 < len(st) {
		if v, ok := st[2].(core.UnsignedIntegerData); ok { s.WindowSizeReceive = uint8(v) }
	}
	if 3 < len(st) {
		if v, ok := st[3].(core.UnsignedLongData); ok { s.MaxInfoFieldLengthTransmit = uint16(v) }
	}
	if 4 < len(st) {
		if v, ok := st[4].(core.UnsignedLongData); ok { s.MaxInfoFieldLengthReceive = uint16(v) }
	}
	if 5 < len(st) {
		if v, ok := st[5].(core.UnsignedLongData); ok { s.InterOctetTimeOut = uint16(v) }
	}
	if 6 < len(st) {
		if v, ok := st[6].(core.UnsignedLongData); ok { s.InterCharacterTimeOut = uint16(v) }
	}
	if 7 < len(st) {
		if v, ok := st[7].(core.UnsignedLongData); ok { s.InactivityTimeOut = uint16(v) }
	}
	return nil
}


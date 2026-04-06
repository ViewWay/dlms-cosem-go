package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// IECTwistedPairSetup is the COSEM IECTwistedPairSetup interface class (IC 24).
type IECTwistedPairSetup struct {
	LogicalName core.ObisCode
	Version     uint8
	Mode                                uint8
	Speed                               uint8
	ConvertTime                         uint16
	Repetitions                         uint8
	Reference                           uint8
}

func (s *IECTwistedPairSetup) ClassID() uint16               { return core.ClassIDIECTwistedPairSetup }
func (s *IECTwistedPairSetup) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *IECTwistedPairSetup) GetVersion() uint8             { return s.Version }

func (s *IECTwistedPairSetup) Access(attr int) core.CosemAttributeAccess {
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
func (s *IECTwistedPairSetup) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.UnsignedIntegerData(s.Mode),
		core.UnsignedIntegerData(s.Speed),
		core.UnsignedLongData(s.ConvertTime),
		core.UnsignedIntegerData(s.Repetitions),
		core.UnsignedIntegerData(s.Reference),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *IECTwistedPairSetup) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return nil
	}
	if 0 < len(st) {
		if v, ok := st[0].(core.UnsignedIntegerData); ok { s.Mode = uint8(v) }
	}
	if 1 < len(st) {
		if v, ok := st[1].(core.UnsignedIntegerData); ok { s.Speed = uint8(v) }
	}
	if 2 < len(st) {
		if v, ok := st[2].(core.UnsignedLongData); ok { s.ConvertTime = uint16(v) }
	}
	if 3 < len(st) {
		if v, ok := st[3].(core.UnsignedIntegerData); ok { s.Repetitions = uint8(v) }
	}
	if 4 < len(st) {
		if v, ok := st[4].(core.UnsignedIntegerData); ok { s.Reference = uint8(v) }
	}
	return nil
}


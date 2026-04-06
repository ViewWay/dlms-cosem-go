package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// AutoAnswer is the COSEM AutoAnswer interface class (IC 28).
type AutoAnswer struct {
	LogicalName core.ObisCode
	Version     uint8
	Mode                                uint8
	ListeningWindow                     []byte
	NumberOfCalls                       uint16
	NumberOfRings                       uint8
	Answered                            bool
}

func (s *AutoAnswer) ClassID() uint16               { return core.ClassIDAutoAnswer }
func (s *AutoAnswer) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *AutoAnswer) GetVersion() uint8             { return s.Version }

func (s *AutoAnswer) Access(attr int) core.CosemAttributeAccess {
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
func (s *AutoAnswer) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.UnsignedIntegerData(s.Mode),
		core.OctetStringData(s.ListeningWindow),
		core.UnsignedLongData(s.NumberOfCalls),
		core.UnsignedIntegerData(s.NumberOfRings),
		core.BooleanData(s.Answered),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *AutoAnswer) UnmarshalBinary(data []byte) error {
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
		if v, ok := st[1].(core.OctetStringData); ok { s.ListeningWindow = []byte(v) }
	}
	if 2 < len(st) {
		if v, ok := st[2].(core.UnsignedLongData); ok { s.NumberOfCalls = uint16(v) }
	}
	if 3 < len(st) {
		if v, ok := st[3].(core.UnsignedIntegerData); ok { s.NumberOfRings = uint8(v) }
	}
	if 4 < len(st) {
		if v, ok := st[4].(core.BooleanData); ok { s.Answered = bool(v) }
	}
	return nil
}


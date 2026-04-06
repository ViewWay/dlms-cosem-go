package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// ModemConfiguration is the COSEM ModemConfiguration interface class (IC 27).
type ModemConfiguration struct {
	LogicalName core.ObisCode
	Version     uint8
	ModemType                           uint8
	Initialized                         bool
	ModemInitializationStrings          []byte
	ModemInitializationResponseTimeout  uint32
}

func (s *ModemConfiguration) ClassID() uint16               { return core.ClassIDModemConfiguration }
func (s *ModemConfiguration) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *ModemConfiguration) GetVersion() uint8             { return s.Version }

func (s *ModemConfiguration) Access(attr int) core.CosemAttributeAccess {
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
func (s *ModemConfiguration) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.UnsignedIntegerData(s.ModemType),
		core.BooleanData(s.Initialized),
		core.OctetStringData(s.ModemInitializationStrings),
		core.DoubleLongUnsignedData(s.ModemInitializationResponseTimeout),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ModemConfiguration) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return nil
	}
	if 0 < len(st) {
		if v, ok := st[0].(core.UnsignedIntegerData); ok { s.ModemType = uint8(v) }
	}
	if 1 < len(st) {
		if v, ok := st[1].(core.BooleanData); ok { s.Initialized = bool(v) }
	}
	if 2 < len(st) {
		if v, ok := st[2].(core.OctetStringData); ok { s.ModemInitializationStrings = []byte(v) }
	}
	if 3 < len(st) {
		if v, ok := st[3].(core.DoubleLongUnsignedData); ok { s.ModemInitializationResponseTimeout = uint32(v) }
	}
	return nil
}


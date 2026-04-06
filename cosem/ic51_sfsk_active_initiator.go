package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// SFSKActiveInitiator is the COSEM SFSKActiveInitiator interface class (IC 51).
type SFSKActiveInitiator struct {
	LogicalName core.ObisCode
	Version     uint8
	ActiveInitiator                     []byte
	ActiveInitiatorCount                uint16
}

func (s *SFSKActiveInitiator) ClassID() uint16               { return core.ClassIDSFSKActiveInitiator }
func (s *SFSKActiveInitiator) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *SFSKActiveInitiator) GetVersion() uint8             { return s.Version }

func (s *SFSKActiveInitiator) Access(attr int) core.CosemAttributeAccess {
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
func (s *SFSKActiveInitiator) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.OctetStringData(s.ActiveInitiator),
		core.UnsignedLongData(s.ActiveInitiatorCount),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *SFSKActiveInitiator) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return nil
	}
	if 0 < len(st) {
		if v, ok := st[0].(core.OctetStringData); ok { s.ActiveInitiator = []byte(v) }
	}
	if 1 < len(st) {
		if v, ok := st[1].(core.UnsignedLongData); ok { s.ActiveInitiatorCount = uint16(v) }
	}
	return nil
}


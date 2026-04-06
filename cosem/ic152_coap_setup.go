package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// CoAPSetup is the COSEM CoAPSetup interface class (IC 152).
type CoAPSetup struct {
	LogicalName core.ObisCode
	Version     uint8
	CoAPServerAddress                   string
	CoAPServerPort                      uint16
	ResponseTimeout                     uint32
	MaxRetransmit                       uint8
	AckTimeout                          uint32
	AckRandomFactor                     uint16
}

func (s *CoAPSetup) ClassID() uint16               { return core.ClassIDCoAPSetup }
func (s *CoAPSetup) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *CoAPSetup) GetVersion() uint8             { return s.Version }

func (s *CoAPSetup) Access(attr int) core.CosemAttributeAccess {
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
func (s *CoAPSetup) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.VisibleStringData(s.CoAPServerAddress),
		core.UnsignedLongData(s.CoAPServerPort),
		core.DoubleLongUnsignedData(s.ResponseTimeout),
		core.UnsignedIntegerData(s.MaxRetransmit),
		core.DoubleLongUnsignedData(s.AckTimeout),
		core.UnsignedLongData(s.AckRandomFactor),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *CoAPSetup) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return nil
	}
	if 0 < len(st) {
		if v, ok := st[0].(core.VisibleStringData); ok { s.CoAPServerAddress = string(v) }
	}
	if 1 < len(st) {
		if v, ok := st[1].(core.UnsignedLongData); ok { s.CoAPServerPort = uint16(v) }
	}
	if 2 < len(st) {
		if v, ok := st[2].(core.DoubleLongUnsignedData); ok { s.ResponseTimeout = uint32(v) }
	}
	if 3 < len(st) {
		if v, ok := st[3].(core.UnsignedIntegerData); ok { s.MaxRetransmit = uint8(v) }
	}
	if 4 < len(st) {
		if v, ok := st[4].(core.DoubleLongUnsignedData); ok { s.AckTimeout = uint32(v) }
	}
	if 5 < len(st) {
		if v, ok := st[5].(core.UnsignedLongData); ok { s.AckRandomFactor = uint16(v) }
	}
	return nil
}


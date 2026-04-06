package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// ZigBeeSASAPSFragmentation is the COSEM ZigBeeSASAPSFragmentation interface class (IC 103).
type ZigBeeSASAPSFragmentation struct {
	LogicalName core.ObisCode
	Version     uint8
	FragmentationEnabled                bool
	WindowSize                          uint8
	InterFrameDelay                     uint16
}

func (s *ZigBeeSASAPSFragmentation) ClassID() uint16               { return core.ClassIDZigBeeSASAPSFragmentation }
func (s *ZigBeeSASAPSFragmentation) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *ZigBeeSASAPSFragmentation) GetVersion() uint8             { return s.Version }

func (s *ZigBeeSASAPSFragmentation) Access(attr int) core.CosemAttributeAccess {
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
func (s *ZigBeeSASAPSFragmentation) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.BooleanData(s.FragmentationEnabled),
		core.UnsignedIntegerData(s.WindowSize),
		core.UnsignedLongData(s.InterFrameDelay),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ZigBeeSASAPSFragmentation) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return nil
	}
	if 0 < len(st) {
		if v, ok := st[0].(core.BooleanData); ok { s.FragmentationEnabled = bool(v) }
	}
	if 1 < len(st) {
		if v, ok := st[1].(core.UnsignedIntegerData); ok { s.WindowSize = uint8(v) }
	}
	if 2 < len(st) {
		if v, ok := st[2].(core.UnsignedLongData); ok { s.InterFrameDelay = uint16(v) }
	}
	return nil
}


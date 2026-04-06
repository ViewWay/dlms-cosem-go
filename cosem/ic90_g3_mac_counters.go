package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// G3MACCounters is the COSEM G3MACCounters interface class (IC 90).
type G3MACCounters struct {
	LogicalName core.ObisCode
	Version     uint8
	MACTxPacketCount                    uint32
	MACRxPacketCount                    uint32
	MACCRCErrorCount                    uint32
	MACTxTime                           uint32
	MACRxTime                           uint32
}

func (s *G3MACCounters) ClassID() uint16               { return core.ClassIDG3MACCounters }
func (s *G3MACCounters) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *G3MACCounters) GetVersion() uint8             { return s.Version }

func (s *G3MACCounters) Access(attr int) core.CosemAttributeAccess {
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
func (s *G3MACCounters) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.DoubleLongUnsignedData(s.MACTxPacketCount),
		core.DoubleLongUnsignedData(s.MACRxPacketCount),
		core.DoubleLongUnsignedData(s.MACCRCErrorCount),
		core.DoubleLongUnsignedData(s.MACTxTime),
		core.DoubleLongUnsignedData(s.MACRxTime),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *G3MACCounters) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return nil
	}
	if 0 < len(st) {
		if v, ok := st[0].(core.DoubleLongUnsignedData); ok { s.MACTxPacketCount = uint32(v) }
	}
	if 1 < len(st) {
		if v, ok := st[1].(core.DoubleLongUnsignedData); ok { s.MACRxPacketCount = uint32(v) }
	}
	if 2 < len(st) {
		if v, ok := st[2].(core.DoubleLongUnsignedData); ok { s.MACCRCErrorCount = uint32(v) }
	}
	if 3 < len(st) {
		if v, ok := st[3].(core.DoubleLongUnsignedData); ok { s.MACTxTime = uint32(v) }
	}
	if 4 < len(st) {
		if v, ok := st[4].(core.DoubleLongUnsignedData); ok { s.MACRxTime = uint32(v) }
	}
	return nil
}


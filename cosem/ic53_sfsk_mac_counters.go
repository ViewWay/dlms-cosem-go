package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// SFSKMACCounters is the COSEM SFSKMACCounters interface class (IC 53).
type SFSKMACCounters struct {
	LogicalName core.ObisCode
	Version     uint8
	TxPacketCount                       uint32
	RxPacketCount                       uint32
	CRCErrorCount                       uint32
	TxTime                              uint32
	RxTime                              uint32
}

func (s *SFSKMACCounters) ClassID() uint16               { return core.ClassIDSFSKMACCounters }
func (s *SFSKMACCounters) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *SFSKMACCounters) GetVersion() uint8             { return s.Version }

func (s *SFSKMACCounters) Access(attr int) core.CosemAttributeAccess {
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
func (s *SFSKMACCounters) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.DoubleLongUnsignedData(s.TxPacketCount),
		core.DoubleLongUnsignedData(s.RxPacketCount),
		core.DoubleLongUnsignedData(s.CRCErrorCount),
		core.DoubleLongUnsignedData(s.TxTime),
		core.DoubleLongUnsignedData(s.RxTime),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *SFSKMACCounters) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return nil
	}
	if 0 < len(st) {
		if v, ok := st[0].(core.DoubleLongUnsignedData); ok { s.TxPacketCount = uint32(v) }
	}
	if 1 < len(st) {
		if v, ok := st[1].(core.DoubleLongUnsignedData); ok { s.RxPacketCount = uint32(v) }
	}
	if 2 < len(st) {
		if v, ok := st[2].(core.DoubleLongUnsignedData); ok { s.CRCErrorCount = uint32(v) }
	}
	if 3 < len(st) {
		if v, ok := st[3].(core.DoubleLongUnsignedData); ok { s.TxTime = uint32(v) }
	}
	if 4 < len(st) {
		if v, ok := st[4].(core.DoubleLongUnsignedData); ok { s.RxTime = uint32(v) }
	}
	return nil
}


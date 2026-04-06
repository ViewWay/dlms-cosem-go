package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// PrimeMACCounters is the COSEM PrimeMACCounters interface class (IC 84).
type PrimeMACCounters struct {
	LogicalName core.ObisCode
	Version     uint8
	MACTxTotal                          uint32
	MACRxTotal                          uint32
	MACTxError                          uint32
	MACRxError                          uint32
	MACTxDropped                        uint32
	MACRxDropped                        uint32
}

func (s *PrimeMACCounters) ClassID() uint16               { return core.ClassIDPrimeMACCounters }
func (s *PrimeMACCounters) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *PrimeMACCounters) GetVersion() uint8             { return s.Version }

func (s *PrimeMACCounters) Access(attr int) core.CosemAttributeAccess {
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
func (s *PrimeMACCounters) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.DoubleLongUnsignedData(s.MACTxTotal),
		core.DoubleLongUnsignedData(s.MACRxTotal),
		core.DoubleLongUnsignedData(s.MACTxError),
		core.DoubleLongUnsignedData(s.MACRxError),
		core.DoubleLongUnsignedData(s.MACTxDropped),
		core.DoubleLongUnsignedData(s.MACRxDropped),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *PrimeMACCounters) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return nil
	}
	if 0 < len(st) {
		if v, ok := st[0].(core.DoubleLongUnsignedData); ok { s.MACTxTotal = uint32(v) }
	}
	if 1 < len(st) {
		if v, ok := st[1].(core.DoubleLongUnsignedData); ok { s.MACRxTotal = uint32(v) }
	}
	if 2 < len(st) {
		if v, ok := st[2].(core.DoubleLongUnsignedData); ok { s.MACTxError = uint32(v) }
	}
	if 3 < len(st) {
		if v, ok := st[3].(core.DoubleLongUnsignedData); ok { s.MACRxError = uint32(v) }
	}
	if 4 < len(st) {
		if v, ok := st[4].(core.DoubleLongUnsignedData); ok { s.MACTxDropped = uint32(v) }
	}
	if 5 < len(st) {
		if v, ok := st[5].(core.DoubleLongUnsignedData); ok { s.MACRxDropped = uint32(v) }
	}
	return nil
}


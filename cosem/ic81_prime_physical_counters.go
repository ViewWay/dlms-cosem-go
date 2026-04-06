package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// PrimePhysicalCounters is the COSEM PrimePhysicalCounters interface class (IC 81).
type PrimePhysicalCounters struct {
	LogicalName core.ObisCode
	Version     uint8
	PhyTxDrop                           uint32
	PhyRxTotal                          uint32
	PhyRxCRCError                       uint32
	PhyTxTotal                          uint32
}

func (s *PrimePhysicalCounters) ClassID() uint16               { return core.ClassIDPrimePhysicalCounters }
func (s *PrimePhysicalCounters) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *PrimePhysicalCounters) GetVersion() uint8             { return s.Version }

func (s *PrimePhysicalCounters) Access(attr int) core.CosemAttributeAccess {
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
func (s *PrimePhysicalCounters) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.DoubleLongUnsignedData(s.PhyTxDrop),
		core.DoubleLongUnsignedData(s.PhyRxTotal),
		core.DoubleLongUnsignedData(s.PhyRxCRCError),
		core.DoubleLongUnsignedData(s.PhyTxTotal),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *PrimePhysicalCounters) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return nil
	}
	if 0 < len(st) {
		if v, ok := st[0].(core.DoubleLongUnsignedData); ok { s.PhyTxDrop = uint32(v) }
	}
	if 1 < len(st) {
		if v, ok := st[1].(core.DoubleLongUnsignedData); ok { s.PhyRxTotal = uint32(v) }
	}
	if 2 < len(st) {
		if v, ok := st[2].(core.DoubleLongUnsignedData); ok { s.PhyRxCRCError = uint32(v) }
	}
	if 3 < len(st) {
		if v, ok := st[3].(core.DoubleLongUnsignedData); ok { s.PhyTxTotal = uint32(v) }
	}
	return nil
}


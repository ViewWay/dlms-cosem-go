package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// Charge is the COSEM Charge interface class (IC 113).
type Charge struct {
	LogicalName core.ObisCode
	Version     uint8
	TotalAmountPaid                     uint32
	ChargeType                          uint8
	Priority                            uint8
	UnitCharge                          uint32
	ChargeConfiguration                 []byte
	LastChargeAt                        []byte
}

func (s *Charge) ClassID() uint16               { return core.ClassIDCharge }
func (s *Charge) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *Charge) GetVersion() uint8             { return s.Version }

func (s *Charge) Access(attr int) core.CosemAttributeAccess {
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
func (s *Charge) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.DoubleLongUnsignedData(s.TotalAmountPaid),
		core.UnsignedIntegerData(s.ChargeType),
		core.UnsignedIntegerData(s.Priority),
		core.DoubleLongUnsignedData(s.UnitCharge),
		core.OctetStringData(s.ChargeConfiguration),
		core.OctetStringData(s.LastChargeAt),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *Charge) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return nil
	}
	if 0 < len(st) {
		if v, ok := st[0].(core.DoubleLongUnsignedData); ok { s.TotalAmountPaid = uint32(v) }
	}
	if 1 < len(st) {
		if v, ok := st[1].(core.UnsignedIntegerData); ok { s.ChargeType = uint8(v) }
	}
	if 2 < len(st) {
		if v, ok := st[2].(core.UnsignedIntegerData); ok { s.Priority = uint8(v) }
	}
	if 3 < len(st) {
		if v, ok := st[3].(core.DoubleLongUnsignedData); ok { s.UnitCharge = uint32(v) }
	}
	if 4 < len(st) {
		if v, ok := st[4].(core.OctetStringData); ok { s.ChargeConfiguration = []byte(v) }
	}
	if 5 < len(st) {
		if v, ok := st[5].(core.OctetStringData); ok { s.LastChargeAt = []byte(v) }
	}
	return nil
}


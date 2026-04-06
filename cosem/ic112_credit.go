package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// Credit is the COSEM Credit interface class (IC 112).
type Credit struct {
	LogicalName core.ObisCode
	Version     uint8
	CurrentCreditAmount                 uint32
	Type                                uint8
	Priority                            uint8
	WarningThreshold                    uint32
	Limit                               uint32
	CreditConfiguration                 []byte
	Status                              uint8
	PresetCreditAmount                  uint32
	AvailableCredit                     uint32
	Amount                              uint32
	CreditThreshold                     uint32
}

func (s *Credit) ClassID() uint16               { return core.ClassIDCredit }
func (s *Credit) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *Credit) GetVersion() uint8             { return s.Version }

func (s *Credit) Access(attr int) core.CosemAttributeAccess {
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
func (s *Credit) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.DoubleLongUnsignedData(s.CurrentCreditAmount),
		core.UnsignedIntegerData(s.Type),
		core.UnsignedIntegerData(s.Priority),
		core.DoubleLongUnsignedData(s.WarningThreshold),
		core.DoubleLongUnsignedData(s.Limit),
		core.OctetStringData(s.CreditConfiguration),
		core.UnsignedIntegerData(s.Status),
		core.DoubleLongUnsignedData(s.PresetCreditAmount),
		core.DoubleLongUnsignedData(s.AvailableCredit),
		core.DoubleLongUnsignedData(s.Amount),
		core.DoubleLongUnsignedData(s.CreditThreshold),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *Credit) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return nil
	}
	if 0 < len(st) {
		if v, ok := st[0].(core.DoubleLongUnsignedData); ok { s.CurrentCreditAmount = uint32(v) }
	}
	if 1 < len(st) {
		if v, ok := st[1].(core.UnsignedIntegerData); ok { s.Type = uint8(v) }
	}
	if 2 < len(st) {
		if v, ok := st[2].(core.UnsignedIntegerData); ok { s.Priority = uint8(v) }
	}
	if 3 < len(st) {
		if v, ok := st[3].(core.DoubleLongUnsignedData); ok { s.WarningThreshold = uint32(v) }
	}
	if 4 < len(st) {
		if v, ok := st[4].(core.DoubleLongUnsignedData); ok { s.Limit = uint32(v) }
	}
	if 5 < len(st) {
		if v, ok := st[5].(core.OctetStringData); ok { s.CreditConfiguration = []byte(v) }
	}
	if 6 < len(st) {
		if v, ok := st[6].(core.UnsignedIntegerData); ok { s.Status = uint8(v) }
	}
	if 7 < len(st) {
		if v, ok := st[7].(core.DoubleLongUnsignedData); ok { s.PresetCreditAmount = uint32(v) }
	}
	if 8 < len(st) {
		if v, ok := st[8].(core.DoubleLongUnsignedData); ok { s.AvailableCredit = uint32(v) }
	}
	if 9 < len(st) {
		if v, ok := st[9].(core.DoubleLongUnsignedData); ok { s.Amount = uint32(v) }
	}
	if 10 < len(st) {
		if v, ok := st[10].(core.DoubleLongUnsignedData); ok { s.CreditThreshold = uint32(v) }
	}
	return nil
}


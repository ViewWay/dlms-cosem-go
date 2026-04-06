package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// SFSKPhyMACSetup is the COSEM SFSKPhyMACSetup interface class (IC 50).
type SFSKPhyMACSetup struct {
	LogicalName core.ObisCode
	Version     uint8
	MACAddress                          [6]byte
	PhyList                             []byte
	TxLevel                             uint8
	RxLevel                             uint8
	FrequencyBand                       uint8
	DataRate                            uint8
	TxFrequency                         uint32
	RxFrequency                         uint32
}

func (s *SFSKPhyMACSetup) ClassID() uint16               { return core.ClassIDSFSKPhyMACSetup }
func (s *SFSKPhyMACSetup) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *SFSKPhyMACSetup) GetVersion() uint8             { return s.Version }

func (s *SFSKPhyMACSetup) Access(attr int) core.CosemAttributeAccess {
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
func (s *SFSKPhyMACSetup) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.OctetStringData(s.MACAddress[:]),
		core.OctetStringData(s.PhyList),
		core.UnsignedIntegerData(s.TxLevel),
		core.UnsignedIntegerData(s.RxLevel),
		core.UnsignedIntegerData(s.FrequencyBand),
		core.UnsignedIntegerData(s.DataRate),
		core.DoubleLongUnsignedData(s.TxFrequency),
		core.DoubleLongUnsignedData(s.RxFrequency),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *SFSKPhyMACSetup) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return nil
	}
	if 0 < len(st) {
		if v, ok := st[0].(core.OctetStringData); ok && len(v) == 6 { copy(s.MACAddress[:], v) }
	}
	if 1 < len(st) {
		if v, ok := st[1].(core.OctetStringData); ok { s.PhyList = []byte(v) }
	}
	if 2 < len(st) {
		if v, ok := st[2].(core.UnsignedIntegerData); ok { s.TxLevel = uint8(v) }
	}
	if 3 < len(st) {
		if v, ok := st[3].(core.UnsignedIntegerData); ok { s.RxLevel = uint8(v) }
	}
	if 4 < len(st) {
		if v, ok := st[4].(core.UnsignedIntegerData); ok { s.FrequencyBand = uint8(v) }
	}
	if 5 < len(st) {
		if v, ok := st[5].(core.UnsignedIntegerData); ok { s.DataRate = uint8(v) }
	}
	if 6 < len(st) {
		if v, ok := st[6].(core.DoubleLongUnsignedData); ok { s.TxFrequency = uint32(v) }
	}
	if 7 < len(st) {
		if v, ok := st[7].(core.DoubleLongUnsignedData); ok { s.RxFrequency = uint32(v) }
	}
	return nil
}


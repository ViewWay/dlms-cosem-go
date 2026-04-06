package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// MBusSlavePortSetup is the COSEM MBusSlavePortSetup interface class (IC 25).
type MBusSlavePortSetup struct {
	LogicalName core.ObisCode
	Version     uint8
	PrimaryAddress                      uint8
	IdentificationNumber                []byte
	ManufacturerID                      []byte
	DeviceVersion                       uint8
	DeviceType                          uint8
	AccessNumber                        []byte
	Status                              uint8
	Aligned                             bool
}

func (s *MBusSlavePortSetup) ClassID() uint16               { return core.ClassIDMBusSlavePortSetup }
func (s *MBusSlavePortSetup) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *MBusSlavePortSetup) GetVersion() uint8             { return s.Version }

func (s *MBusSlavePortSetup) Access(attr int) core.CosemAttributeAccess {
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
func (s *MBusSlavePortSetup) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.UnsignedIntegerData(s.PrimaryAddress),
		core.OctetStringData(s.IdentificationNumber),
		core.OctetStringData(s.ManufacturerID),
		core.UnsignedIntegerData(s.DeviceVersion),
		core.UnsignedIntegerData(s.DeviceType),
		core.OctetStringData(s.AccessNumber),
		core.UnsignedIntegerData(s.Status),
		core.BooleanData(s.Aligned),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *MBusSlavePortSetup) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return nil
	}
	if 0 < len(st) {
		if v, ok := st[0].(core.UnsignedIntegerData); ok { s.PrimaryAddress = uint8(v) }
	}
	if 1 < len(st) {
		if v, ok := st[1].(core.OctetStringData); ok { s.IdentificationNumber = []byte(v) }
	}
	if 2 < len(st) {
		if v, ok := st[2].(core.OctetStringData); ok { s.ManufacturerID = []byte(v) }
	}
	if 3 < len(st) {
		if v, ok := st[3].(core.UnsignedIntegerData); ok { s.DeviceVersion = uint8(v) }
	}
	if 4 < len(st) {
		if v, ok := st[4].(core.UnsignedIntegerData); ok { s.DeviceType = uint8(v) }
	}
	if 5 < len(st) {
		if v, ok := st[5].(core.OctetStringData); ok { s.AccessNumber = []byte(v) }
	}
	if 6 < len(st) {
		if v, ok := st[6].(core.UnsignedIntegerData); ok { s.Status = uint8(v) }
	}
	if 7 < len(st) {
		if v, ok := st[7].(core.BooleanData); ok { s.Aligned = bool(v) }
	}
	return nil
}


package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// MBusClient is the COSEM M-Bus Client interface class (IC 26).
type MBusClient struct {
	LogicalName          core.ObisCode
	MBusPort             uint8
	PrimaryAddress       uint8
	IdentificationNumber uint32
	ManufacturerID       string
	Version              uint8
	DeviceType           uint8
	AccessNumber         uint8
	Status               uint8
	VersionField         uint8
}

func (m *MBusClient) ClassID() uint16               { return core.ClassIDMBusClient }
func (m *MBusClient) GetLogicalName() core.ObisCode { return m.LogicalName }
func (m *MBusClient) GetVersion() uint8             { return m.Version }
func (m *MBusClient) Access(attr int) core.CosemAttributeAccess {
	switch attr {
	case 1:
		return core.CosemAttributeAccess{Access: core.AccessRead}
	case 2, 3, 4, 5, 6, 7, 8, 9, 10:
		return core.CosemAttributeAccess{Access: core.AccessReadWrite}
	default:
		return core.CosemAttributeAccess{Access: core.AccessNone}
	}
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (m *MBusClient) MarshalBinary() ([]byte, error) {
	return core.UnsignedIntegerData(m.PrimaryAddress).ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (m *MBusClient) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	if v, ok := elem.(core.UnsignedIntegerData); ok {
		m.PrimaryAddress = uint8(v)
	}
	return nil
}

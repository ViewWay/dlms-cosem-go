package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// IPv6Setup is the COSEM IPv6Setup interface class (IC 48).
type IPv6Setup struct {
	LogicalName core.ObisCode
	Version     uint8
	AddressConfigurationMode            uint8
	LinkLocalAddress                    [16]byte
	AddressList                         []byte
}

func (s *IPv6Setup) ClassID() uint16               { return core.ClassIDIPv6Setup }
func (s *IPv6Setup) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *IPv6Setup) GetVersion() uint8             { return s.Version }

func (s *IPv6Setup) Access(attr int) core.CosemAttributeAccess {
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
func (s *IPv6Setup) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.UnsignedIntegerData(s.AddressConfigurationMode),
		core.OctetStringData(s.LinkLocalAddress[:]),
		core.OctetStringData(s.AddressList),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *IPv6Setup) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return nil
	}
	if 0 < len(st) {
		if v, ok := st[0].(core.UnsignedIntegerData); ok { s.AddressConfigurationMode = uint8(v) }
	}
	if 1 < len(st) {
		if v, ok := st[1].(core.OctetStringData); ok && len(v) == 16 { copy(s.LinkLocalAddress[:], v) }
	}
	if 2 < len(st) {
		if v, ok := st[2].(core.OctetStringData); ok { s.AddressList = []byte(v) }
	}
	return nil
}


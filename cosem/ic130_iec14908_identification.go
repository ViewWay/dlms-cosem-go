package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// IEC14908Identification is the COSEM IEC14908Identification interface class (IC 130).
type IEC14908Identification struct {
	LogicalName core.ObisCode
	Version     uint8
	DomainAddress                       []byte
	SubnetAddress                       uint8
	NodeAddress                         uint8
}

func (s *IEC14908Identification) ClassID() uint16               { return core.ClassIDIEC14908Identification }
func (s *IEC14908Identification) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *IEC14908Identification) GetVersion() uint8             { return s.Version }

func (s *IEC14908Identification) Access(attr int) core.CosemAttributeAccess {
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
func (s *IEC14908Identification) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.OctetStringData(s.DomainAddress),
		core.UnsignedIntegerData(s.SubnetAddress),
		core.UnsignedIntegerData(s.NodeAddress),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *IEC14908Identification) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return nil
	}
	if 0 < len(st) {
		if v, ok := st[0].(core.OctetStringData); ok { s.DomainAddress = []byte(v) }
	}
	if 1 < len(st) {
		if v, ok := st[1].(core.UnsignedIntegerData); ok { s.SubnetAddress = uint8(v) }
	}
	if 2 < len(st) {
		if v, ok := st[2].(core.UnsignedIntegerData); ok { s.NodeAddress = uint8(v) }
	}
	return nil
}


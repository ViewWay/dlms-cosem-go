package cosem

import (
	"fmt"

	"github.com/ViewWay/dlms-cosem-go/core"
)

// IPv4Setup is the COSEM IPv4 Setup interface class (IC 40).
type IPv4Setup struct {
	LogicalName  core.ObisCode
	IPAddress    [4]byte
	SubnetMask   [4]byte
	Gateway      [4]byte
	PrimaryDNS   [4]byte
	SecondaryDNS [4]byte
	Version      uint8
}

func (s *IPv4Setup) ClassID() uint16               { return core.ClassIDIPv4Setup }
func (s *IPv4Setup) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *IPv4Setup) GetVersion() uint8             { return s.Version }
func (s *IPv4Setup) Access(attr int) core.CosemAttributeAccess {
	switch attr {
	case 1:
		return core.CosemAttributeAccess{Access: core.AccessRead}
	case 2, 3, 4, 5, 6:
		return core.CosemAttributeAccess{Access: core.AccessReadWrite}
	default:
		return core.CosemAttributeAccess{Access: core.AccessNone}
	}
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s *IPv4Setup) MarshalBinary() ([]byte, error) {
	return core.OctetStringData(s.IPAddress[:]).ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *IPv4Setup) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	if v, ok := elem.(core.OctetStringData); ok && len(v) == 4 {
		s.IPAddress = [4]byte{v[0], v[1], v[2], v[3]}
	}
	return nil
}

// SetIPAddress sets the IP address.
func (s *IPv4Setup) SetIPAddress(ip [4]byte) { s.IPAddress = ip }

// SetSubnetMask sets the subnet mask.
func (s *IPv4Setup) SetSubnetMask(mask [4]byte) { s.SubnetMask = mask }

// SetGateway sets the gateway.
func (s *IPv4Setup) SetGateway(gw [4]byte) { s.Gateway = gw }

// MarshalObjectList returns all IP fields as a structure.
func (s *IPv4Setup) MarshalObjectList() ([]byte, error) {
	return core.StructureData{
		core.OctetStringData(s.IPAddress[:]),
		core.OctetStringData(s.SubnetMask[:]),
		core.OctetStringData(s.Gateway[:]),
		core.OctetStringData(s.PrimaryDNS[:]),
		core.OctetStringData(s.SecondaryDNS[:]),
	}.ToBytes(), nil
}

// UnmarshalObjectList parses all IP fields from a structure.
func (s *IPv4Setup) UnmarshalObjectList(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return fmt.Errorf("ipv4_setup: expected structure")
	}
	for i, f := range []*[4]byte{&s.IPAddress, &s.SubnetMask, &s.Gateway, &s.PrimaryDNS, &s.SecondaryDNS} {
		if i < len(st) {
			if v, ok := st[i].(core.OctetStringData); ok && len(v) == 4 {
				*f = [4]byte{v[0], v[1], v[2], v[3]}
			}
		}
	}
	return nil
}

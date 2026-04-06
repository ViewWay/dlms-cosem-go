package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// PrimeMACNetworkAdmin is the COSEM PrimeMACNetworkAdmin interface class (IC 85).
type PrimeMACNetworkAdmin struct {
	LogicalName core.ObisCode
	Version     uint8
	NetworkID                           [8]byte
	NetworkRole                         uint8
	NetworkState                        uint8
	NetworkDeviceList                   []byte
}

func (s *PrimeMACNetworkAdmin) ClassID() uint16               { return core.ClassIDPrimeMACNetworkAdmin }
func (s *PrimeMACNetworkAdmin) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *PrimeMACNetworkAdmin) GetVersion() uint8             { return s.Version }

func (s *PrimeMACNetworkAdmin) Access(attr int) core.CosemAttributeAccess {
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
func (s *PrimeMACNetworkAdmin) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.OctetStringData(s.NetworkID[:]),
		core.UnsignedIntegerData(s.NetworkRole),
		core.UnsignedIntegerData(s.NetworkState),
		core.OctetStringData(s.NetworkDeviceList),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *PrimeMACNetworkAdmin) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return nil
	}
	if 0 < len(st) {
		if v, ok := st[0].(core.OctetStringData); ok && len(v) == 8 { copy(s.NetworkID[:], v) }
	}
	if 1 < len(st) {
		if v, ok := st[1].(core.UnsignedIntegerData); ok { s.NetworkRole = uint8(v) }
	}
	if 2 < len(st) {
		if v, ok := st[2].(core.UnsignedIntegerData); ok { s.NetworkState = uint8(v) }
	}
	if 3 < len(st) {
		if v, ok := st[3].(core.OctetStringData); ok { s.NetworkDeviceList = []byte(v) }
	}
	return nil
}


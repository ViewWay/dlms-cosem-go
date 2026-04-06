package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// MACAddressSetup is the COSEM MACAddressSetup interface class (IC 43).
type MACAddressSetup struct {
	LogicalName core.ObisCode
	Version     uint8
	MACAddress                          [6]byte
}

func (s *MACAddressSetup) ClassID() uint16               { return core.ClassIDMACAddressSetup }
func (s *MACAddressSetup) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *MACAddressSetup) GetVersion() uint8             { return s.Version }

func (s *MACAddressSetup) Access(attr int) core.CosemAttributeAccess {
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
func (s *MACAddressSetup) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.OctetStringData(s.MACAddress[:]),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *MACAddressSetup) UnmarshalBinary(data []byte) error {
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
	return nil
}


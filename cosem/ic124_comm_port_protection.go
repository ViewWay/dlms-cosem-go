package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// CommPortProtection is the COSEM CommPortProtection interface class (IC 124).
type CommPortProtection struct {
	LogicalName core.ObisCode
	Version     uint8
	PortProtectionParameters            []byte
}

func (s *CommPortProtection) ClassID() uint16               { return core.ClassIDCommPortProtection }
func (s *CommPortProtection) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *CommPortProtection) GetVersion() uint8             { return s.Version }

func (s *CommPortProtection) Access(attr int) core.CosemAttributeAccess {
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
func (s *CommPortProtection) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.OctetStringData(s.PortProtectionParameters),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *CommPortProtection) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return nil
	}
	if 0 < len(st) {
		if v, ok := st[0].(core.OctetStringData); ok { s.PortProtectionParameters = []byte(v) }
	}
	return nil
}


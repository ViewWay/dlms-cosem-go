package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// COSEMDataProtection is the COSEM COSEMDataProtection interface class (IC 30).
type COSEMDataProtection struct {
	LogicalName core.ObisCode
	Version     uint8
	ProtectionParametersGet             []byte
	ProtectionParametersSet             []byte
	ProtectionBuffer                    []byte
}

func (s *COSEMDataProtection) ClassID() uint16               { return core.ClassIDCOSEMDataProtection }
func (s *COSEMDataProtection) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *COSEMDataProtection) GetVersion() uint8             { return s.Version }

func (s *COSEMDataProtection) Access(attr int) core.CosemAttributeAccess {
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
func (s *COSEMDataProtection) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.OctetStringData(s.ProtectionParametersGet),
		core.OctetStringData(s.ProtectionParametersSet),
		core.OctetStringData(s.ProtectionBuffer),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *COSEMDataProtection) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return nil
	}
	if 0 < len(st) {
		if v, ok := st[0].(core.OctetStringData); ok { s.ProtectionParametersGet = []byte(v) }
	}
	if 1 < len(st) {
		if v, ok := st[1].(core.OctetStringData); ok { s.ProtectionParametersSet = []byte(v) }
	}
	if 2 < len(st) {
		if v, ok := st[2].(core.OctetStringData); ok { s.ProtectionBuffer = []byte(v) }
	}
	return nil
}


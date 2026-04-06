package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// IEC62055Attributes is the COSEM IEC62055Attributes interface class (IC 116).
type IEC62055Attributes struct {
	LogicalName core.ObisCode
	Version     uint8
	STSKeyIdentificationNo              []byte
	STSKeyRevisionNo                    []byte
	STSKeyExpiryDate                    []byte
	STSTokenCarrierIdentification       []byte
	STSTokenDecoderKeyStatus            uint8
}

func (s *IEC62055Attributes) ClassID() uint16               { return core.ClassIDIEC62055Attributes }
func (s *IEC62055Attributes) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *IEC62055Attributes) GetVersion() uint8             { return s.Version }

func (s *IEC62055Attributes) Access(attr int) core.CosemAttributeAccess {
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
func (s *IEC62055Attributes) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.OctetStringData(s.STSKeyIdentificationNo),
		core.OctetStringData(s.STSKeyRevisionNo),
		core.OctetStringData(s.STSKeyExpiryDate),
		core.OctetStringData(s.STSTokenCarrierIdentification),
		core.UnsignedIntegerData(s.STSTokenDecoderKeyStatus),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *IEC62055Attributes) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return nil
	}
	if 0 < len(st) {
		if v, ok := st[0].(core.OctetStringData); ok { s.STSKeyIdentificationNo = []byte(v) }
	}
	if 1 < len(st) {
		if v, ok := st[1].(core.OctetStringData); ok { s.STSKeyRevisionNo = []byte(v) }
	}
	if 2 < len(st) {
		if v, ok := st[2].(core.OctetStringData); ok { s.STSKeyExpiryDate = []byte(v) }
	}
	if 3 < len(st) {
		if v, ok := st[3].(core.OctetStringData); ok { s.STSTokenCarrierIdentification = []byte(v) }
	}
	if 4 < len(st) {
		if v, ok := st[4].(core.UnsignedIntegerData); ok { s.STSTokenDecoderKeyStatus = uint8(v) }
	}
	return nil
}


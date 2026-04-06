package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// IECKeyType represents the key type
type IECKeyType uint8

const (
	IECKeyTypeNone       IECKeyType = 0
	IECKeyTypeECDSA_P256 IECKeyType = 1
	IECKeyTypeECDSA_P384 IECKeyType = 2
	IECKeyTypeECDSA_P521 IECKeyType = 3
	IECKeyTypeRSA_2048   IECKeyType = 4
	IECKeyTypeRSA_3072   IECKeyType = 5
	IECKeyTypeRSA_4096   IECKeyType = 6
)

// IECPublicKey is the COSEM IEC Public Key interface class (IC 32).
type IECPublicKey struct {
	LogicalName    core.ObisCode
	KeyID          uint16
	KeyType        IECKeyType
	PublicKeyValue []byte
	KeyUsage       uint8
	ValidityStart  core.CosemDate
	ValidityEnd    core.CosemDate
	Version        uint8
}

func (pk *IECPublicKey) ClassID() uint16               { return core.ClassIDIECPublicKey }
func (pk *IECPublicKey) GetLogicalName() core.ObisCode { return pk.LogicalName }
func (pk *IECPublicKey) GetVersion() uint8             { return pk.Version }
func (pk *IECPublicKey) Access(attr int) core.CosemAttributeAccess {
	switch attr {
	case 2, 3, 4, 5:
		return core.CosemAttributeAccess{Access: core.AccessRead}
	case 6, 7:
		return core.CosemAttributeAccess{Access: core.AccessReadWrite}
	default:
		return core.CosemAttributeAccess{Access: core.AccessNone}
	}
}

func (pk *IECPublicKey) MarshalBinary() ([]byte, error) {
	return core.OctetStringData(pk.PublicKeyValue).ToBytes(), nil
}

func (pk *IECPublicKey) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	if v, ok := elem.(core.OctetStringData); ok {
		pk.PublicKeyValue = v
	}
	return nil
}

func (pk *IECPublicKey) SetPublicKey(keyType IECKeyType, value []byte) {
	pk.KeyType = keyType
	pk.PublicKeyValue = value
}

func (pk *IECPublicKey) IsValid() bool {
	return len(pk.PublicKeyValue) > 0
}

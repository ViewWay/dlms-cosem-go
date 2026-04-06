package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// Account is the COSEM Account interface class (IC 22).
type Account struct {
	LogicalName       core.ObisCode
	VendorInfo        []byte
	PaymentMode       core.DlmsData
	AccountStatus     core.DlmsData
	CreditReference   core.DlmsData
	ChargeReference   core.DlmsData
	CurrentCredit     core.DlmsData
	AvailableCredit   core.DlmsData
	Currency          core.DlmsData
	LowCreditWarning  core.DlmsData
	LowCreditLimit    core.DlmsData
	OverdraftLimit    core.DlmsData
	ActivationDate    core.CosemDate
	DeactivationDate  core.CosemDate
	Version           uint8
}

func (a *Account) ClassID() uint16               { return core.ClassIDAccount }
func (a *Account) GetLogicalName() core.ObisCode { return a.LogicalName }
func (a *Account) GetVersion() uint8             { return a.Version }
func (a *Account) Access(attr int) core.CosemAttributeAccess {
	switch attr {
	case 1, 2:
		return core.CosemAttributeAccess{Access: core.AccessRead}
	case 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13:
		return core.CosemAttributeAccess{Access: core.AccessReadWrite}
	default:
		return core.CosemAttributeAccess{Access: core.AccessNone}
	}
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (a *Account) MarshalBinary() ([]byte, error) {
	return core.OctetStringData(a.VendorInfo).ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (a *Account) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	if v, ok := elem.(core.OctetStringData); ok {
		a.VendorInfo = v
	}
	return nil
}

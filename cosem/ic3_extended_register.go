package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// ExtendedRegister is the COSEM Extended Register interface class (IC 3).
type ExtendedRegister struct {
	LogicalName core.ObisCode
	Value       core.DlmsData
	Scaler      int8
	Unit        uint8
	Status      core.DlmsData
	CaptureTime core.CosemDateTime
	Version     uint8
}

func (e *ExtendedRegister) ClassID() uint16               { return core.ClassIDExtendedRegister }
func (e *ExtendedRegister) GetLogicalName() core.ObisCode { return e.LogicalName }
func (e *ExtendedRegister) GetVersion() uint8             { return e.Version }
func (e *ExtendedRegister) Access(attr int) core.CosemAttributeAccess {
	switch attr {
	case 2:
		return core.CosemAttributeAccess{Access: core.AccessRead}
	case 3:
		return core.CosemAttributeAccess{Access: core.AccessRead}
	default:
		return core.CosemAttributeAccess{Access: core.AccessNone}
	}
}

func (e *ExtendedRegister) MarshalBinary() ([]byte, error) {
	return e.Value.ToBytes(), nil
}

func (e *ExtendedRegister) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	e.Value = elem
	return nil
}

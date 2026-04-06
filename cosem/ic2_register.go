package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// Register is the COSEM Register interface class (IC 2).
type Register struct {
	LogicalName core.ObisCode
	Value       core.DlmsData
	Scaler      int8
	Unit        uint8
	Status      core.DlmsData
	Version     uint8
}

func (r *Register) ClassID() uint16               { return core.ClassIDRegister }
func (r *Register) GetLogicalName() core.ObisCode { return r.LogicalName }
func (r *Register) GetVersion() uint8             { return r.Version }
func (r *Register) Access(attr int) core.CosemAttributeAccess {
	switch attr {
	case 2:
		return core.CosemAttributeAccess{Access: core.AccessRead}
	case 3:
		return core.CosemAttributeAccess{Access: core.AccessRead}
	default:
		return core.CosemAttributeAccess{Access: core.AccessNone}
	}
}

// MarshalBinary encodes the register value.
func (r *Register) MarshalBinary() ([]byte, error) {
	if r.Value == nil {
		return core.DoubleLongUnsignedData(0).ToBytes(), nil
	}
	return r.Value.ToBytes(), nil
}

// UnmarshalBinary decodes the register value.
func (r *Register) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	r.Value = elem
	return nil
}

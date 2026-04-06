package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// ValueDisplay is the COSEM Value Display interface class (IC 30).
type ValueDisplay struct {
	LogicalName    core.ObisCode
	ValueToDisplay core.DlmsData
	Status         uint8
	Version        uint8
}

func (vd *ValueDisplay) ClassID() uint16               { return core.ClassIDValueDisplay }
func (vd *ValueDisplay) GetLogicalName() core.ObisCode { return vd.LogicalName }
func (vd *ValueDisplay) GetVersion() uint8             { return vd.Version }
func (vd *ValueDisplay) Access(attr int) core.CosemAttributeAccess {
	switch attr {
	case 1:
		return core.CosemAttributeAccess{Access: core.AccessRead}
	case 2:
		return core.CosemAttributeAccess{Access: core.AccessReadWrite}
	case 3:
		return core.CosemAttributeAccess{Access: core.AccessRead}
	default:
		return core.CosemAttributeAccess{Access: core.AccessNone}
	}
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (vd *ValueDisplay) MarshalBinary() ([]byte, error) {
	if vd.ValueToDisplay == nil {
		return core.NullData{}.ToBytes(), nil
	}
	return vd.ValueToDisplay.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (vd *ValueDisplay) UnmarshalBinary(data []byte) error {
	if len(data) == 0 {
		vd.ValueToDisplay = core.NullData{}
		return nil
	}
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	vd.ValueToDisplay = elem
	return nil
}

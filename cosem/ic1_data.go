package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// Data is the COSEM Data interface class (IC 1).
type Data struct {
	LogicalName core.ObisCode
	Value       core.DlmsData
	Version     uint8
}

func (d *Data) ClassID() uint16               { return core.ClassIDData }
func (d *Data) GetLogicalName() core.ObisCode { return d.LogicalName }
func (d *Data) GetVersion() uint8             { return d.Version }
func (d *Data) Access(attr int) core.CosemAttributeAccess {
	switch attr {
	case 1:
		return core.CosemAttributeAccess{Access: core.AccessReadWrite}
	case 2:
		return core.CosemAttributeAccess{Access: core.AccessRead}
	default:
		return core.CosemAttributeAccess{Access: core.AccessNone}
	}
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (d *Data) MarshalBinary() ([]byte, error) {
	if d.Value == nil {
		return core.NullData{}.ToBytes(), nil
	}
	return d.Value.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (d *Data) UnmarshalBinary(data []byte) error {
	if len(data) == 0 {
		d.Value = core.NullData{}
		return nil
	}
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	d.Value = elem
	return nil
}

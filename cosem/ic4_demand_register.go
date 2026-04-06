package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// DemandRegister is the COSEM Demand Register interface class (IC 4).
type DemandRegister struct {
	LogicalName     core.ObisCode
	CurrentValue    core.DlmsData
	LastValue       core.DlmsData
	Scaler          int8
	Unit            uint8
	Status          core.DlmsData
	CaptureTime     core.CosemDateTime
	StartTime       core.CosemDateTime
	Period          uint16
	NumberOfPeriods uint8
	Version         uint8
}

func (d *DemandRegister) ClassID() uint16               { return core.ClassIDDemandRegister }
func (d *DemandRegister) GetLogicalName() core.ObisCode { return d.LogicalName }
func (d *DemandRegister) GetVersion() uint8             { return d.Version }
func (d *DemandRegister) Access(attr int) core.CosemAttributeAccess {
	switch attr {
	case 2, 3, 4:
		return core.CosemAttributeAccess{Access: core.AccessRead}
	default:
		return core.CosemAttributeAccess{Access: core.AccessNone}
	}
}

func (d *DemandRegister) MarshalBinary() ([]byte, error) {
	return d.CurrentValue.ToBytes(), nil
}

func (d *DemandRegister) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	d.CurrentValue = elem
	return nil
}

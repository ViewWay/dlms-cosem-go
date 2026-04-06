package cosem

import (
	"fmt"

	"github.com/ViewWay/dlms-cosem-go/core"
)

// RegisterMonitor is the COSEM Register Monitor interface class (IC 13).
type RegisterMonitor struct {
	LogicalName      core.ObisCode
	Thresholds       []core.DlmsData
	ThresholdActions []core.DlmsData
	MonitoredValue   core.DlmsData
	Version          uint8
}

func (rm *RegisterMonitor) ClassID() uint16               { return core.ClassIDRegisterMonitor }
func (rm *RegisterMonitor) GetLogicalName() core.ObisCode { return rm.LogicalName }
func (rm *RegisterMonitor) GetVersion() uint8             { return rm.Version }
func (rm *RegisterMonitor) Access(attr int) core.CosemAttributeAccess {
	switch attr {
	case 2, 3, 4:
		return core.CosemAttributeAccess{Access: core.AccessRead}
	default:
		return core.CosemAttributeAccess{Access: core.AccessNone}
	}
}

func (rm *RegisterMonitor) MarshalBinary() ([]byte, error) {
	arr := core.ArrayData{}
	for _, t := range rm.Thresholds {
		arr = append(arr, t)
	}
	return arr.ToBytes(), nil
}

func (rm *RegisterMonitor) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	arr, ok := elem.(core.ArrayData)
	if !ok {
		return fmt.Errorf("register_monitor: expected array")
	}
	rm.Thresholds = []core.DlmsData(arr)
	return nil
}

func (rm *RegisterMonitor) AddThreshold(threshold core.DlmsData) {
	rm.Thresholds = append(rm.Thresholds, threshold)
}

func (rm *RegisterMonitor) ThresholdCount() int {
	return len(rm.Thresholds)
}

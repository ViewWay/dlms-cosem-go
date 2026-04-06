package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// LimiterAction represents limiter action settings
type LimiterAction struct {
	ScriptLN       core.ObisCode
	ScriptSelector uint8
}

// Limiter is the COSEM Limiter interface class (IC 17).
type Limiter struct {
	LogicalName      core.ObisCode
	MonitoredValue   core.DlmsData
	ThresholdActive  core.DlmsData
	ThresholdNormal  core.DlmsData
	MinOverThreshold uint16
	MinUnderThreshold uint16
	Actions          []LimiterAction
	EmergencyProfile core.ObisCode
	Version          uint8
}

func (l *Limiter) ClassID() uint16               { return core.ClassIDLimiter }
func (l *Limiter) GetLogicalName() core.ObisCode { return l.LogicalName }
func (l *Limiter) GetVersion() uint8             { return l.Version }
func (l *Limiter) Access(attr int) core.CosemAttributeAccess {
	switch attr {
	case 2, 3, 4, 5, 6:
		return core.CosemAttributeAccess{Access: core.AccessRead}
	default:
		return core.CosemAttributeAccess{Access: core.AccessNone}
	}
}

func (l *Limiter) MarshalBinary() ([]byte, error) {
	return l.MonitoredValue.ToBytes(), nil
}

func (l *Limiter) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	l.MonitoredValue = elem
	return nil
}

func (l *Limiter) SetThreshold(active, normal core.DlmsData) {
	l.ThresholdActive = active
	l.ThresholdNormal = normal
}

func (l *Limiter) AddAction(action LimiterAction) {
	l.Actions = append(l.Actions, action)
}

func (l *Limiter) ActionCount() int {
	return len(l.Actions)
}

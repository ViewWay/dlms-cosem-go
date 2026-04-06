package cosem

import (
	"fmt"
	"time"

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

// numericValue extracts a float64 from common numeric DlmsData types.
func numericValue(d core.DlmsData) (float64, error) {
	switch v := d.(type) {
	case core.DoubleLongUnsignedData:
		return float64(v), nil
	case core.DoubleLongData:
		return float64(v), nil
	case core.LongData:
		return float64(v), nil
	case core.UnsignedLongData:
		return float64(v), nil
	case core.IntegerData:
		return float64(v), nil
	case core.UnsignedIntegerData:
		return float64(v), nil
	case core.Long64Data:
		return float64(v), nil
	case core.UnsignedLong64Data:
		return float64(v), nil
	case core.Float32Data:
		return float64(v), nil
	case core.Float64Data:
		return float64(v), nil
	default:
		return 0, fmt.Errorf("limiter: unsupported data type %T", d)
	}
}

// CheckThreshold compares MonitoredValue against ThresholdActive and
// ThresholdNormal. Returns (over, under):
//   - over:  MonitoredValue >= ThresholdActive
//   - under: MonitoredValue <= ThresholdNormal
func (l *Limiter) CheckThreshold() (over bool, under bool) {
	if l.MonitoredValue == nil {
		return false, false
	}
	val, err := numericValue(l.MonitoredValue)
	if err != nil {
		return false, false
	}
	if l.ThresholdActive != nil {
		act, err := numericValue(l.ThresholdActive)
		if err == nil && val >= act {
			over = true
		}
	}
	if l.ThresholdNormal != nil {
		nrm, err := numericValue(l.ThresholdNormal)
		if err == nil && val <= nrm {
			under = true
		}
	}
	return
}

// ActivateEmergency sets the EmergencyProfile to the given script logical name,
// effectively activating emergency limiting mode.
func (l *Limiter) ActivateEmergency(profile core.ObisCode) {
	l.EmergencyProfile = profile
}

// DeactivateEmergency clears the EmergencyProfile.
func (l *Limiter) DeactivateEmergency() {
	l.EmergencyProfile = core.ObisCode{}
}

// EvaluateThresholds performs a comprehensive threshold evaluation.
// It checks MonitoredValue against thresholds and the duration against
// MinOverThreshold / MinUnderThreshold to determine the action to take.
func (l *Limiter) EvaluateThresholds(duration time.Duration) LimiterAction {
	over, under := l.CheckThreshold()

	// Emergency profile takes priority
	if l.EmergencyProfile != (core.ObisCode{}) {
		return LimiterAction{
			ScriptLN:       l.EmergencyProfile,
			ScriptSelector: 0,
		}
	}

	// Check over-threshold duration
	if over && l.MinOverThreshold > 0 {
		minDur := time.Duration(l.MinOverThreshold) * time.Second
		if duration >= minDur {
			return l.pickAction(0) // first action = over-threshold action
		}
	}

	// Check under-threshold duration
	if under && l.MinUnderThreshold > 0 {
		minDur := time.Duration(l.MinUnderThreshold) * time.Second
		if duration >= minDur {
			return l.pickAction(1) // second action = under-threshold action
		}
	}

	return LimiterAction{} // no action
}

// pickAction returns the action at index i, or empty if not available.
func (l *Limiter) pickAction(i int) LimiterAction {
	if i < len(l.Actions) {
		return l.Actions[i]
	}
	return LimiterAction{}
}

// ThresholdStatus returns a string representation of the current status.
func (l *Limiter) ThresholdStatus() string {
	over, under := l.CheckThreshold()
	val, _ := numericValue(l.MonitoredValue)
	switch {
	case over && under:
		return fmt.Sprintf("BOTH (value=%.2f)", val)
	case over:
		return fmt.Sprintf("OVER (value=%.2f)", val)
	case under:
		return fmt.Sprintf("UNDER (value=%.2f)", val)
	default:
		return fmt.Sprintf("NORMAL (value=%.2f)", val)
	}
}

// IsEmergencyActive returns true if an emergency profile is set.
func (l *Limiter) IsEmergencyActive() bool {
	return l.EmergencyProfile != (core.ObisCode{})
}

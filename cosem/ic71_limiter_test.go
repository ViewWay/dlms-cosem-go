package cosem

import (
	"testing"
	"time"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestLimiter_ClassID(t *testing.T) {
	l := &Limiter{}
	if l.ClassID() != 71 {
		t.Errorf("ClassID() = %d, want 71", l.ClassID())
	}
	if l.ClassID() != core.ClassIDLimiter {
		t.Error("ClassID mismatch with const")
	}
}

func TestLimiter_New(t *testing.T) {
	l := &Limiter{
		LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255},
	}
	if l.MonitoredValue != nil {
		t.Error("MonitoredValue should default nil")
	}
}

func TestLimiter_MarshalBinary(t *testing.T) {
	tests := []struct {
		name  string
		value core.DlmsData
	}{
		{"zero", core.DoubleLongUnsignedData(0)},
		{"value", core.DoubleLongUnsignedData(5000)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Limiter{MonitoredValue: tt.value}
			b, err := l.MarshalBinary()
			if err != nil {
				t.Fatal(err)
			}
			l2 := &Limiter{}
			if err := l2.UnmarshalBinary(b); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestLimiter_CheckThreshold(t *testing.T) {
	l := &Limiter{
		MonitoredValue:  core.DoubleLongUnsignedData(5000),
		ThresholdActive:  core.DoubleLongUnsignedData(5000),
		ThresholdNormal:  core.DoubleLongUnsignedData(1000),
	}
	over, under := l.CheckThreshold()
	if !over {
		t.Error("expected over=true")
	}
	if under {
		t.Error("expected under=false")
	}

	l.MonitoredValue = core.DoubleLongUnsignedData(500)
	over, under = l.CheckThreshold()
	if over {
		t.Error("expected over=false")
	}
	if !under {
		t.Error("expected under=true")
	}

	l.MonitoredValue = core.DoubleLongUnsignedData(3000)
	over, under = l.CheckThreshold()
	if over || under {
		t.Error("expected both false (normal range)")
	}

	// nil thresholds
	l2 := &Limiter{MonitoredValue: core.DoubleLongUnsignedData(100)}
	over, under = l2.CheckThreshold()
	if over || under {
		t.Error("expected both false with nil thresholds")
	}
}

func TestLimiter_Emergency(t *testing.T) {
	l := &Limiter{}
	if l.IsEmergencyActive() {
		t.Error("emergency should not be active initially")
	}

	profile := core.ObisCode{0, 0, 0, 0, 1, 255}
	l.ActivateEmergency(profile)
	if !l.IsEmergencyActive() {
		t.Error("emergency should be active")
	}
	if l.EmergencyProfile != profile {
		t.Error("profile mismatch")
	}

	l.DeactivateEmergency()
	if l.IsEmergencyActive() {
		t.Error("emergency should be deactivated")
	}
}

func TestLimiter_EvaluateThresholds(t *testing.T) {
	l := &Limiter{
		MonitoredValue:    core.DoubleLongUnsignedData(5000),
		ThresholdActive:   core.DoubleLongUnsignedData(5000),
		MinOverThreshold:  10,
		Actions: []LimiterAction{
			{ScriptLN: core.ObisCode{0, 0, 0, 0, 1, 255}, ScriptSelector: 1},
		},
	}

	// Duration not met — no action
	action := l.EvaluateThresholds(5 * time.Second)
	if action.ScriptLN != (core.ObisCode{}) {
		t.Error("expected no action when duration not met")
	}

	// Duration met — should trigger over-threshold action
	action = l.EvaluateThresholds(10 * time.Second)
	if action.ScriptLN != (core.ObisCode{0, 0, 0, 0, 1, 255}) {
		t.Error("expected over-threshold action")
	}

	// Emergency takes priority
	emergency := core.ObisCode{0, 0, 0, 0, 9, 255}
	l.ActivateEmergency(emergency)
	action = l.EvaluateThresholds(10 * time.Second)
	if action.ScriptLN != emergency {
		t.Error("emergency profile should take priority")
	}
}

func TestLimiter_ThresholdStatus(t *testing.T) {
	l := &Limiter{
		MonitoredValue:  core.DoubleLongUnsignedData(5000),
		ThresholdActive:  core.DoubleLongUnsignedData(5000),
	}
	s := l.ThresholdStatus()
	if s == "" {
		t.Error("expected non-empty status")
	}
}

func TestLimiter_Fields(t *testing.T) {
	l := &Limiter{
		LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255},
		Version:     0,
	}
	if l.GetVersion() != 0 {
		t.Error("GetVersion mismatch")
	}
}

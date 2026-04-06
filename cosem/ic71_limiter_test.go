package cosem

import (
	"testing"

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

func TestLimiter_Fields(t *testing.T) {
	l := &Limiter{
		LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255},
		Version:     0,
	}
	if l.GetVersion() != 0 {
		t.Error("GetVersion mismatch")
	}
}

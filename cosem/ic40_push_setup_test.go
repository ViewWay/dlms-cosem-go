package cosem

import (
	"testing"
	"time"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestPushSetup_ClassID(t *testing.T) {
	ps := &PushSetup{}
	if ps.ClassID() != 40 {
		t.Errorf("ClassID() = %d, want 40", ps.ClassID())
	}
	if ps.ClassID() != core.ClassIDPushSetup {
		t.Error("ClassID mismatch with const")
	}
}

func TestPushSetup_New(t *testing.T) {
	ps := &PushSetup{
		LogicalName:    core.ObisCode{0, 0, 41, 0, 0, 255},
		Service:        1,
		NumberOfRetries: 3,
		Version:        1,
	}
	if ps.Service != 1 {
		t.Error("Service mismatch")
	}
	if ps.NumberOfRetries != 3 {
		t.Error("NumberOfRetries mismatch")
	}
}

func TestPushSetup_MarshalBinary(t *testing.T) {
	ps := &PushSetup{
		PushObjectList: []PushObject{
			{ClassID: 1, LogicalName: core.ObisCode{0, 0, 1, 0, 0, 255}, Attribute: 2},
		},
	}
	b, err := ps.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	ps2 := &PushSetup{}
	if err := ps2.UnmarshalBinary(b); err != nil {
		t.Fatal(err)
	}
}

func TestPushSetup_ExecutePush(t *testing.T) {
	ps := &PushSetup{
		PushObjectList: []PushObject{
			{ClassID: 1, LogicalName: core.ObisCode{0, 0, 1, 0, 0, 255}, Attribute: 2},
		},
		Destination: []byte{0x01, 0x02},
		Service:     1,
	}
	payload, dest, err := ps.ExecutePush()
	if err != nil {
		t.Fatal(err)
	}
	if len(payload) == 0 {
		t.Error("expected non-empty payload")
	}
	if dest.Service != 1 {
		t.Errorf("dest.Service = %d, want 1", dest.Service)
	}
}

func TestPushSetup_ExecutePush_Empty(t *testing.T) {
	ps := &PushSetup{}
	_, _, err := ps.ExecutePush()
	if err == nil {
		t.Error("expected error for empty object list")
	}
}

func TestPushSetup_IsWithinWindow(t *testing.T) {
	ps := &PushSetup{RandomisationStartInterval: 60}
	windowEnd := time.Date(2026, 4, 6, 12, 0, 0, 0, time.Local)

	// Inside window (30s before end)
	if !ps.IsWithinWindow(windowEnd.Add(-30*time.Second), windowEnd) {
		t.Error("should be within window")
	}
	// At boundary
	if !ps.IsWithinWindow(windowEnd.Add(-60*time.Second), windowEnd) {
		t.Error("should be at window start")
	}
	// Outside window (too early)
	if ps.IsWithinWindow(windowEnd.Add(-61*time.Second), windowEnd) {
		t.Error("should be outside window")
	}
	// Outside window (too late)
	if ps.IsWithinWindow(windowEnd.Add(time.Second), windowEnd) {
		t.Error("should be outside window")
	}
}

func TestPushSetup_ShouldRetry(t *testing.T) {
	ps := &PushSetup{NumberOfRetries: 3}
	if !ps.ShouldRetry(1) {
		t.Error("attempt 1 should retry")
	}
	if !ps.ShouldRetry(3) {
		t.Error("attempt 3 should retry")
	}
	if ps.ShouldRetry(4) {
		t.Error("attempt 4 should not retry")
	}
	if ps.ShouldRetry(0) {
		t.Error("attempt 0 should not retry")
	}
}

func TestPushSetup_GetRetryDelay(t *testing.T) {
	ps := &PushSetup{RepetitionDelay: 10}
	for i := 0; i < 5; i++ {
		d := ps.GetRetryDelay(1)
		if d < 10*time.Second {
			t.Errorf("delay %v less than base 10s", d)
		}
	}
}

func TestPushSetup_Fields(t *testing.T) {
	ps := &PushSetup{
		LogicalName: core.ObisCode{0, 0, 41, 0, 0, 255},
		Destination: []byte{0x01, 0x02},
		Version:     2,
	}
	if ps.GetVersion() != 2 {
		t.Error("GetVersion mismatch")
	}
	if len(ps.Destination) != 2 {
		t.Error("Destination length mismatch")
	}
}

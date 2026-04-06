package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestRegisterMonitor_ClassID(t *testing.T) {
	rm := &RegisterMonitor{}
	if rm.ClassID() != 21 {
		t.Errorf("ClassID() = %d, want 21", rm.ClassID())
	}
	if rm.ClassID() != core.ClassIDRegisterMonitor {
		t.Error("ClassID mismatch with const")
	}
}

func TestRegisterMonitor_New(t *testing.T) {
	rm := &RegisterMonitor{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if rm.Thresholds == nil {
		t.Log("Thresholds nil by default - ok")
	}
}

func TestRegisterMonitor_MarshalBinary(t *testing.T) {
	rm := &RegisterMonitor{}
	b, err := rm.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	rm2 := &RegisterMonitor{}
	if err := rm2.UnmarshalBinary(b); err != nil {
		t.Fatal(err)
	}
}

func TestRegisterMonitor_Fields(t *testing.T) {
	rm := &RegisterMonitor{
		LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255},
		Version:     0,
	}
	if rm.GetVersion() != 0 {
		t.Error("GetVersion mismatch")
	}
}

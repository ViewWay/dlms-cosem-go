package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestSingleActionSchedule_ClassID(t *testing.T) {
	sas := &SingleActionSchedule{}
	if sas.ClassID() != 22 {
		t.Errorf("ClassID() = %d, want 22", sas.ClassID())
	}
	if sas.ClassID() != core.ClassIDSingleActionSchedule {
		t.Error("ClassID mismatch with const")
	}
}

func TestSingleActionSchedule_New(t *testing.T) {
	sas := &SingleActionSchedule{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if sas.Entries == nil {
		t.Log("Entries nil by default - ok")
	}
}

func TestSingleActionSchedule_MarshalBinary(t *testing.T) {
	sas := &SingleActionSchedule{}
	b, err := sas.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	sas2 := &SingleActionSchedule{}
	if err := sas2.UnmarshalBinary(b); err != nil {
		t.Fatal(err)
	}
}

func TestSingleActionSchedule_Fields(t *testing.T) {
	sas := &SingleActionSchedule{
		LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255},
		Version:     0,
	}
	if sas.GetVersion() != 0 {
		t.Error("GetVersion mismatch")
	}
}

package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestSchedule_ClassID(t *testing.T) {
	s := &Schedule{}
	if s.ClassID() != 10 {
		t.Errorf("ClassID() = %d, want 10", s.ClassID())
	}
	if s.ClassID() != core.ClassIDSchedule {
		t.Error("ClassID mismatch with const")
	}
}

func TestSchedule_New(t *testing.T) {
	s := &Schedule{LogicalName: core.ObisCode{0, 0, 10, 0, 0, 255}}
	if s.Entries == nil {
		t.Log("Entries nil by default - ok")
	}
}

func TestSchedule_MarshalBinary(t *testing.T) {
	s := &Schedule{}
	b, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &Schedule{}
	if err := s2.UnmarshalBinary(b); err != nil {
		t.Fatal(err)
	}
}

func TestSchedule_Fields(t *testing.T) {
	s := &Schedule{
		LogicalName: core.ObisCode{0, 0, 10, 0, 0, 255},
		Version:     0,
	}
	if s.GetVersion() != 0 {
		t.Error("GetVersion mismatch")
	}
}

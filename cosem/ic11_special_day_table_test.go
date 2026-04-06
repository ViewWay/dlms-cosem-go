package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestSpecialDaysTable_ClassID(t *testing.T) {
	s := &SpecialDaysTable{}
	if s.ClassID() != 11 {
		t.Errorf("ClassID() = %d, want 11", s.ClassID())
	}
	if s.ClassID() != core.ClassIDSpecialDaysTable {
		t.Error("ClassID mismatch with const")
	}
}

func TestSpecialDaysTable_New(t *testing.T) {
	s := &SpecialDaysTable{LogicalName: core.ObisCode{0, 0, 11, 0, 0, 255}}
	if s.Entries == nil {
		t.Log("Entries nil by default - ok")
	}
}

func TestSpecialDaysTable_MarshalBinary(t *testing.T) {
	s := &SpecialDaysTable{}
	b, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &SpecialDaysTable{}
	if err := s2.UnmarshalBinary(b); err != nil {
		t.Fatal(err)
	}
}

func TestSpecialDaysTable_Fields(t *testing.T) {
	s := &SpecialDaysTable{
		LogicalName: core.ObisCode{0, 0, 11, 0, 0, 255},
		Version:     0,
	}
	if s.GetVersion() != 0 {
		t.Error("GetVersion mismatch")
	}
}

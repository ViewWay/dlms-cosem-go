package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestWiSUNSetup_ClassID(t *testing.T) {
	s := &WiSUNSetup{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDWiSUNSetup {
		t.Errorf("expected %d, got %d", core.ClassIDWiSUNSetup, s.ClassID())
	}
}

func TestWiSUNSetup_MarshalUnmarshalBinary(t *testing.T) {
	s := &WiSUNSetup{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &WiSUNSetup{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

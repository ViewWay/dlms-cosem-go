package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestCharge_ClassID(t *testing.T) {
	s := &Charge{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDCharge {
		t.Errorf("expected %d, got %d", core.ClassIDCharge, s.ClassID())
	}
}

func TestCharge_MarshalUnmarshalBinary(t *testing.T) {
	s := &Charge{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &Charge{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

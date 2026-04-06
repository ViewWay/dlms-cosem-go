package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestArbitrator_ClassID(t *testing.T) {
	s := &Arbitrator{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDArbitrator {
		t.Errorf("expected %d, got %d", core.ClassIDArbitrator, s.ClassID())
	}
}

func TestArbitrator_MarshalUnmarshalBinary(t *testing.T) {
	s := &Arbitrator{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &Arbitrator{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

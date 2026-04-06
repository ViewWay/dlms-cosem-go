package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestCoAPSetup_ClassID(t *testing.T) {
	s := &CoAPSetup{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDCoAPSetup {
		t.Errorf("expected %d, got %d", core.ClassIDCoAPSetup, s.ClassID())
	}
}

func TestCoAPSetup_MarshalUnmarshalBinary(t *testing.T) {
	s := &CoAPSetup{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &CoAPSetup{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

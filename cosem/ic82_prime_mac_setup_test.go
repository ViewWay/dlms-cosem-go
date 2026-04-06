package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestPrimeMACSetup_ClassID(t *testing.T) {
	s := &PrimeMACSetup{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDPrimeMACSetup {
		t.Errorf("expected %d, got %d", core.ClassIDPrimeMACSetup, s.ClassID())
	}
}

func TestPrimeMACSetup_MarshalUnmarshalBinary(t *testing.T) {
	s := &PrimeMACSetup{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &PrimeMACSetup{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

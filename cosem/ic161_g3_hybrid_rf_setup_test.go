package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestG3HybridRFSetup_ClassID(t *testing.T) {
	s := &G3HybridRFSetup{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDG3HybridRFSetup {
		t.Errorf("expected %d, got %d", core.ClassIDG3HybridRFSetup, s.ClassID())
	}
}

func TestG3HybridRFSetup_MarshalUnmarshalBinary(t *testing.T) {
	s := &G3HybridRFSetup{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &G3HybridRFSetup{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestG3Hybrid6LoWPANSetup_ClassID(t *testing.T) {
	s := &G3Hybrid6LoWPANSetup{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDG3Hybrid6LoWPANSetup {
		t.Errorf("expected %d, got %d", core.ClassIDG3Hybrid6LoWPANSetup, s.ClassID())
	}
}

func TestG3Hybrid6LoWPANSetup_MarshalUnmarshalBinary(t *testing.T) {
	s := &G3Hybrid6LoWPANSetup{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &G3Hybrid6LoWPANSetup{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

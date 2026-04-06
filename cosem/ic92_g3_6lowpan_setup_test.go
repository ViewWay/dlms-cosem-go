package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestG36LoWPANSetup_ClassID(t *testing.T) {
	s := &G36LoWPANSetup{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDG36LoWPANSetup {
		t.Errorf("expected %d, got %d", core.ClassIDG36LoWPANSetup, s.ClassID())
	}
}

func TestG36LoWPANSetup_MarshalUnmarshalBinary(t *testing.T) {
	s := &G36LoWPANSetup{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &G36LoWPANSetup{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

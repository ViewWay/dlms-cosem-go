package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestHSPLCIPSSASSetup_ClassID(t *testing.T) {
	s := &HSPLCIPSSASSetup{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDHSPLCIPSSASSetup {
		t.Errorf("expected %d, got %d", core.ClassIDHSPLCIPSSASSetup, s.ClassID())
	}
}

func TestHSPLCIPSSASSetup_MarshalUnmarshalBinary(t *testing.T) {
	s := &HSPLCIPSSASSetup{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &HSPLCIPSSASSetup{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

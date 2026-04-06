package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestHSPLCCPASSetup_ClassID(t *testing.T) {
	s := &HSPLCCPASSetup{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDHSPLCCPASSetup {
		t.Errorf("expected %d, got %d", core.ClassIDHSPLCCPASSetup, s.ClassID())
	}
}

func TestHSPLCCPASSetup_MarshalUnmarshalBinary(t *testing.T) {
	s := &HSPLCCPASSetup{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &HSPLCCPASSetup{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

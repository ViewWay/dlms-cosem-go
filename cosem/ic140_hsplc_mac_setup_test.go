package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestHSPLCMACSetup_ClassID(t *testing.T) {
	s := &HSPLCMACSetup{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDHSPLCMACSetup {
		t.Errorf("expected %d, got %d", core.ClassIDHSPLCMACSetup, s.ClassID())
	}
}

func TestHSPLCMACSetup_MarshalUnmarshalBinary(t *testing.T) {
	s := &HSPLCMACSetup{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &HSPLCMACSetup{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

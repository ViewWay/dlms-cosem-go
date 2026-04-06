package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestSCHCLPWANSetup_ClassID(t *testing.T) {
	s := &SCHCLPWANSetup{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDSCHCLPWANSetup {
		t.Errorf("expected %d, got %d", core.ClassIDSCHCLPWANSetup, s.ClassID())
	}
}

func TestSCHCLPWANSetup_MarshalUnmarshalBinary(t *testing.T) {
	s := &SCHCLPWANSetup{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &SCHCLPWANSetup{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestSFSKPhyMACSetup_ClassID(t *testing.T) {
	s := &SFSKPhyMACSetup{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDSFSKPhyMACSetup {
		t.Errorf("expected %d, got %d", core.ClassIDSFSKPhyMACSetup, s.ClassID())
	}
}

func TestSFSKPhyMACSetup_MarshalUnmarshalBinary(t *testing.T) {
	s := &SFSKPhyMACSetup{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &SFSKPhyMACSetup{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestMBusMasterPortSetup_ClassID(t *testing.T) {
	s := &MBusMasterPortSetup{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDMBusMasterPortSetup {
		t.Errorf("expected %d, got %d", core.ClassIDMBusMasterPortSetup, s.ClassID())
	}
}

func TestMBusMasterPortSetup_MarshalUnmarshalBinary(t *testing.T) {
	s := &MBusMasterPortSetup{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &MBusMasterPortSetup{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

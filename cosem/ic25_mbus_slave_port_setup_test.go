package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestMBusSlavePortSetup_ClassID(t *testing.T) {
	s := &MBusSlavePortSetup{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDMBusSlavePortSetup {
		t.Errorf("expected %d, got %d", core.ClassIDMBusSlavePortSetup, s.ClassID())
	}
}

func TestMBusSlavePortSetup_MarshalUnmarshalBinary(t *testing.T) {
	s := &MBusSlavePortSetup{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &MBusSlavePortSetup{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

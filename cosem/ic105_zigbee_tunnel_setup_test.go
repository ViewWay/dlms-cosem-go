package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestZigBeeTunnelSetup_ClassID(t *testing.T) {
	s := &ZigBeeTunnelSetup{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDZigBeeTunnelSetup {
		t.Errorf("expected %d, got %d", core.ClassIDZigBeeTunnelSetup, s.ClassID())
	}
}

func TestZigBeeTunnelSetup_MarshalUnmarshalBinary(t *testing.T) {
	s := &ZigBeeTunnelSetup{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &ZigBeeTunnelSetup{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

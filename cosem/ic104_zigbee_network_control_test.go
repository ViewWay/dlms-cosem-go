package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestZigBeeNetworkControl_ClassID(t *testing.T) {
	s := &ZigBeeNetworkControl{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDZigBeeNetworkControl {
		t.Errorf("expected %d, got %d", core.ClassIDZigBeeNetworkControl, s.ClassID())
	}
}

func TestZigBeeNetworkControl_MarshalUnmarshalBinary(t *testing.T) {
	s := &ZigBeeNetworkControl{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &ZigBeeNetworkControl{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

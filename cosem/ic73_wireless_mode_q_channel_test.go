package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestWirelessModeQChannel_ClassID(t *testing.T) {
	s := &WirelessModeQChannel{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDWirelessModeQChannel {
		t.Errorf("expected %d, got %d", core.ClassIDWirelessModeQChannel, s.ClassID())
	}
}

func TestWirelessModeQChannel_MarshalUnmarshalBinary(t *testing.T) {
	s := &WirelessModeQChannel{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &WirelessModeQChannel{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

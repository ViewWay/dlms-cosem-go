package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestMACAddressSetup_ClassID(t *testing.T) {
	s := &MACAddressSetup{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDMACAddressSetup {
		t.Errorf("expected %d, got %d", core.ClassIDMACAddressSetup, s.ClassID())
	}
}

func TestMACAddressSetup_MarshalUnmarshalBinary(t *testing.T) {
	s := &MACAddressSetup{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &MACAddressSetup{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

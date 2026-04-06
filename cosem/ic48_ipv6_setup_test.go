package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestIPv6Setup_ClassID(t *testing.T) {
	s := &IPv6Setup{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDIPv6Setup {
		t.Errorf("expected %d, got %d", core.ClassIDIPv6Setup, s.ClassID())
	}
}

func TestIPv6Setup_MarshalUnmarshalBinary(t *testing.T) {
	s := &IPv6Setup{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &IPv6Setup{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

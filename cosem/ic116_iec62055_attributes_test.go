package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestIEC62055Attributes_ClassID(t *testing.T) {
	s := &IEC62055Attributes{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDIEC62055Attributes {
		t.Errorf("expected %d, got %d", core.ClassIDIEC62055Attributes, s.ClassID())
	}
}

func TestIEC62055Attributes_MarshalUnmarshalBinary(t *testing.T) {
	s := &IEC62055Attributes{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &IEC62055Attributes{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

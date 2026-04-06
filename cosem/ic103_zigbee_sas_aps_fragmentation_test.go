package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestZigBeeSASAPSFragmentation_ClassID(t *testing.T) {
	s := &ZigBeeSASAPSFragmentation{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDZigBeeSASAPSFragmentation {
		t.Errorf("expected %d, got %d", core.ClassIDZigBeeSASAPSFragmentation, s.ClassID())
	}
}

func TestZigBeeSASAPSFragmentation_MarshalUnmarshalBinary(t *testing.T) {
	s := &ZigBeeSASAPSFragmentation{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &ZigBeeSASAPSFragmentation{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

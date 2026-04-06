package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestZigBeeSASJoin_ClassID(t *testing.T) {
	s := &ZigBeeSASJoin{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDZigBeeSASJoin {
		t.Errorf("expected %d, got %d", core.ClassIDZigBeeSASJoin, s.ClassID())
	}
}

func TestZigBeeSASJoin_MarshalUnmarshalBinary(t *testing.T) {
	s := &ZigBeeSASJoin{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &ZigBeeSASJoin{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

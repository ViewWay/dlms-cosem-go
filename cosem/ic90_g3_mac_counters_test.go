package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestG3MACCounters_ClassID(t *testing.T) {
	s := &G3MACCounters{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDG3MACCounters {
		t.Errorf("expected %d, got %d", core.ClassIDG3MACCounters, s.ClassID())
	}
}

func TestG3MACCounters_MarshalUnmarshalBinary(t *testing.T) {
	s := &G3MACCounters{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &G3MACCounters{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestG3HybridRFCounters_ClassID(t *testing.T) {
	s := &G3HybridRFCounters{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDG3HybridRFCounters {
		t.Errorf("expected %d, got %d", core.ClassIDG3HybridRFCounters, s.ClassID())
	}
}

func TestG3HybridRFCounters_MarshalUnmarshalBinary(t *testing.T) {
	s := &G3HybridRFCounters{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &G3HybridRFCounters{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

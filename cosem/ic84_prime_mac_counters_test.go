package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestPrimeMACCounters_ClassID(t *testing.T) {
	s := &PrimeMACCounters{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDPrimeMACCounters {
		t.Errorf("expected %d, got %d", core.ClassIDPrimeMACCounters, s.ClassID())
	}
}

func TestPrimeMACCounters_MarshalUnmarshalBinary(t *testing.T) {
	s := &PrimeMACCounters{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &PrimeMACCounters{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

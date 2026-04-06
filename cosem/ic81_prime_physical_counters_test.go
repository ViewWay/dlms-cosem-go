package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestPrimePhysicalCounters_ClassID(t *testing.T) {
	s := &PrimePhysicalCounters{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDPrimePhysicalCounters {
		t.Errorf("expected %d, got %d", core.ClassIDPrimePhysicalCounters, s.ClassID())
	}
}

func TestPrimePhysicalCounters_MarshalUnmarshalBinary(t *testing.T) {
	s := &PrimePhysicalCounters{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &PrimePhysicalCounters{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

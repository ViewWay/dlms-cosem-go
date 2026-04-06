package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestSFSKMACCounters_ClassID(t *testing.T) {
	s := &SFSKMACCounters{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDSFSKMACCounters {
		t.Errorf("expected %d, got %d", core.ClassIDSFSKMACCounters, s.ClassID())
	}
}

func TestSFSKMACCounters_MarshalUnmarshalBinary(t *testing.T) {
	s := &SFSKMACCounters{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &SFSKMACCounters{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

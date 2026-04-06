package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestSFSKActiveInitiator_ClassID(t *testing.T) {
	s := &SFSKActiveInitiator{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDSFSKActiveInitiator {
		t.Errorf("expected %d, got %d", core.ClassIDSFSKActiveInitiator, s.ClassID())
	}
}

func TestSFSKActiveInitiator_MarshalUnmarshalBinary(t *testing.T) {
	s := &SFSKActiveInitiator{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &SFSKActiveInitiator{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

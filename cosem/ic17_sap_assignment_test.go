package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestSAPAssignment_ClassID(t *testing.T) {
	s := &SAPAssignment{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDSAPAssignment {
		t.Errorf("expected %d, got %d", core.ClassIDSAPAssignment, s.ClassID())
	}
}

func TestSAPAssignment_MarshalUnmarshalBinary(t *testing.T) {
	s := &SAPAssignment{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &SAPAssignment{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

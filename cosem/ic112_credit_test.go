package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestCredit_ClassID(t *testing.T) {
	s := &Credit{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDCredit {
		t.Errorf("expected %d, got %d", core.ClassIDCredit, s.ClassID())
	}
}

func TestCredit_MarshalUnmarshalBinary(t *testing.T) {
	s := &Credit{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &Credit{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

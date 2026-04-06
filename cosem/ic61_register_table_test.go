package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestRegisterTable_ClassID(t *testing.T) {
	s := &RegisterTable{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDRegisterTable {
		t.Errorf("expected %d, got %d", core.ClassIDRegisterTable, s.ClassID())
	}
}

func TestRegisterTable_MarshalUnmarshalBinary(t *testing.T) {
	s := &RegisterTable{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &RegisterTable{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

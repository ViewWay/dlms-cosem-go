package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestFunctionControl_ClassID(t *testing.T) {
	s := &FunctionControl{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDFunctionControl {
		t.Errorf("expected %d, got %d", core.ClassIDFunctionControl, s.ClassID())
	}
}

func TestFunctionControl_MarshalUnmarshalBinary(t *testing.T) {
	s := &FunctionControl{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &FunctionControl{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

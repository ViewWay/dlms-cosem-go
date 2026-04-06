package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestPrimeMACFuncParams_ClassID(t *testing.T) {
	s := &PrimeMACFuncParams{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDPrimeMACFuncParams {
		t.Errorf("expected %d, got %d", core.ClassIDPrimeMACFuncParams, s.ClassID())
	}
}

func TestPrimeMACFuncParams_MarshalUnmarshalBinary(t *testing.T) {
	s := &PrimeMACFuncParams{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &PrimeMACFuncParams{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestLLCType1Setup_ClassID(t *testing.T) {
	s := &LLCType1Setup{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDLLCType1Setup {
		t.Errorf("expected %d, got %d", core.ClassIDLLCType1Setup, s.ClassID())
	}
}

func TestLLCType1Setup_MarshalUnmarshalBinary(t *testing.T) {
	s := &LLCType1Setup{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &LLCType1Setup{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

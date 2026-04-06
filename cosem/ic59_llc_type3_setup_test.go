package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestLLCType3Setup_ClassID(t *testing.T) {
	s := &LLCType3Setup{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDLLCType3Setup {
		t.Errorf("expected %d, got %d", core.ClassIDLLCType3Setup, s.ClassID())
	}
}

func TestLLCType3Setup_MarshalUnmarshalBinary(t *testing.T) {
	s := &LLCType3Setup{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &LLCType3Setup{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

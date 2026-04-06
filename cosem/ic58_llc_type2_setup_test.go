package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestLLCType2Setup_ClassID(t *testing.T) {
	s := &LLCType2Setup{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDLLCType2Setup {
		t.Errorf("expected %d, got %d", core.ClassIDLLCType2Setup, s.ClassID())
	}
}

func TestLLCType2Setup_MarshalUnmarshalBinary(t *testing.T) {
	s := &LLCType2Setup{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &LLCType2Setup{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

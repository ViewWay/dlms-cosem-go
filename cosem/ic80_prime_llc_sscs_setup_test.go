package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestPrimeLLCSSCSSetup_ClassID(t *testing.T) {
	s := &PrimeLLCSSCSSetup{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDPrimeLLCSSCSSetup {
		t.Errorf("expected %d, got %d", core.ClassIDPrimeLLCSSCSSetup, s.ClassID())
	}
}

func TestPrimeLLCSSCSSetup_MarshalUnmarshalBinary(t *testing.T) {
	s := &PrimeLLCSSCSSetup{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &PrimeLLCSSCSSetup{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

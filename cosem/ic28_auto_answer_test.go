package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestAutoAnswer_ClassID(t *testing.T) {
	s := &AutoAnswer{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDAutoAnswer {
		t.Errorf("expected %d, got %d", core.ClassIDAutoAnswer, s.ClassID())
	}
}

func TestAutoAnswer_MarshalUnmarshalBinary(t *testing.T) {
	s := &AutoAnswer{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &AutoAnswer{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

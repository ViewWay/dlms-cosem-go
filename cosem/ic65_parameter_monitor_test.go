package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestParameterMonitor_ClassID(t *testing.T) {
	s := &ParameterMonitor{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDParameterMonitor {
		t.Errorf("expected %d, got %d", core.ClassIDParameterMonitor, s.ClassID())
	}
}

func TestParameterMonitor_MarshalUnmarshalBinary(t *testing.T) {
	s := &ParameterMonitor{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &ParameterMonitor{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

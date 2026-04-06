package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestLTEMonitoring_ClassID(t *testing.T) {
	s := &LTEMonitoring{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDLTEMonitoring {
		t.Errorf("expected %d, got %d", core.ClassIDLTEMonitoring, s.ClassID())
	}
}

func TestLTEMonitoring_MarshalUnmarshalBinary(t *testing.T) {
	s := &LTEMonitoring{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &LTEMonitoring{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

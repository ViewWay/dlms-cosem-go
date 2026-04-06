package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestZigBeeSASStartup_ClassID(t *testing.T) {
	s := &ZigBeeSASStartup{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDZigBeeSASStartup {
		t.Errorf("expected %d, got %d", core.ClassIDZigBeeSASStartup, s.ClassID())
	}
}

func TestZigBeeSASStartup_MarshalUnmarshalBinary(t *testing.T) {
	s := &ZigBeeSASStartup{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &ZigBeeSASStartup{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

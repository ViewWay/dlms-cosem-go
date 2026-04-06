package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestG3MACSetup_ClassID(t *testing.T) {
	s := &G3MACSetup{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDG3MACSetup {
		t.Errorf("expected %d, got %d", core.ClassIDG3MACSetup, s.ClassID())
	}
}

func TestG3MACSetup_MarshalUnmarshalBinary(t *testing.T) {
	s := &G3MACSetup{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &G3MACSetup{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

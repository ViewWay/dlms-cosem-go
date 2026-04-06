package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestIECTwistedPairSetup_ClassID(t *testing.T) {
	s := &IECTwistedPairSetup{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDIECTwistedPairSetup {
		t.Errorf("expected %d, got %d", core.ClassIDIECTwistedPairSetup, s.ClassID())
	}
}

func TestIECTwistedPairSetup_MarshalUnmarshalBinary(t *testing.T) {
	s := &IECTwistedPairSetup{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &IECTwistedPairSetup{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

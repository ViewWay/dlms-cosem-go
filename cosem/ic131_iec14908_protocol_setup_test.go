package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestIEC14908ProtocolSetup_ClassID(t *testing.T) {
	s := &IEC14908ProtocolSetup{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDIEC14908ProtocolSetup {
		t.Errorf("expected %d, got %d", core.ClassIDIEC14908ProtocolSetup, s.ClassID())
	}
}

func TestIEC14908ProtocolSetup_MarshalUnmarshalBinary(t *testing.T) {
	s := &IEC14908ProtocolSetup{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &IEC14908ProtocolSetup{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

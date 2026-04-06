package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestIEC14908ProtocolStatus_ClassID(t *testing.T) {
	s := &IEC14908ProtocolStatus{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDIEC14908ProtocolStatus {
		t.Errorf("expected %d, got %d", core.ClassIDIEC14908ProtocolStatus, s.ClassID())
	}
}

func TestIEC14908ProtocolStatus_MarshalUnmarshalBinary(t *testing.T) {
	s := &IEC14908ProtocolStatus{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &IEC14908ProtocolStatus{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

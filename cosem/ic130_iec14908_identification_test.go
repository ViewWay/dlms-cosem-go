package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestIEC14908Identification_ClassID(t *testing.T) {
	s := &IEC14908Identification{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDIEC14908Identification {
		t.Errorf("expected %d, got %d", core.ClassIDIEC14908Identification, s.ClassID())
	}
}

func TestIEC14908Identification_MarshalUnmarshalBinary(t *testing.T) {
	s := &IEC14908Identification{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &IEC14908Identification{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

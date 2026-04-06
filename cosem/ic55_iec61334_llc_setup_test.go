package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestIEC61334LLCSetup_ClassID(t *testing.T) {
	s := &IEC61334LLCSetup{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDIEC61334LLCSetup {
		t.Errorf("expected %d, got %d", core.ClassIDIEC61334LLCSetup, s.ClassID())
	}
}

func TestIEC61334LLCSetup_MarshalUnmarshalBinary(t *testing.T) {
	s := &IEC61334LLCSetup{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &IEC61334LLCSetup{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

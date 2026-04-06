package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestWiSUNDiagnostic_ClassID(t *testing.T) {
	s := &WiSUNDiagnostic{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDWiSUNDiagnostic {
		t.Errorf("expected %d, got %d", core.ClassIDWiSUNDiagnostic, s.ClassID())
	}
}

func TestWiSUNDiagnostic_MarshalUnmarshalBinary(t *testing.T) {
	s := &WiSUNDiagnostic{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &WiSUNDiagnostic{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

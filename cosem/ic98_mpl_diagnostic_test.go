package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestMPLDiagnostic_ClassID(t *testing.T) {
	s := &MPLDiagnostic{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDMPLDiagnostic {
		t.Errorf("expected %d, got %d", core.ClassIDMPLDiagnostic, s.ClassID())
	}
}

func TestMPLDiagnostic_MarshalUnmarshalBinary(t *testing.T) {
	s := &MPLDiagnostic{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &MPLDiagnostic{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

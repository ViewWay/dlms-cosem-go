package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestRPLDiagnostic_ClassID(t *testing.T) {
	s := &RPLDiagnostic{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDRPLDiagnostic {
		t.Errorf("expected %d, got %d", core.ClassIDRPLDiagnostic, s.ClassID())
	}
}

func TestRPLDiagnostic_MarshalUnmarshalBinary(t *testing.T) {
	s := &RPLDiagnostic{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &RPLDiagnostic{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

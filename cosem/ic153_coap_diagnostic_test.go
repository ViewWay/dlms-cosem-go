package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestCoAPDiagnostic_ClassID(t *testing.T) {
	s := &CoAPDiagnostic{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDCoAPDiagnostic {
		t.Errorf("expected %d, got %d", core.ClassIDCoAPDiagnostic, s.ClassID())
	}
}

func TestCoAPDiagnostic_MarshalUnmarshalBinary(t *testing.T) {
	s := &CoAPDiagnostic{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &CoAPDiagnostic{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

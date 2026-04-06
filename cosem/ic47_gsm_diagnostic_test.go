package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestGSMDiagnostic_ClassID(t *testing.T) {
	s := &GSMDiagnostic{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDGSMDiagnostic {
		t.Errorf("expected %d, got %d", core.ClassIDGSMDiagnostic, s.ClassID())
	}
}

func TestGSMDiagnostic_MarshalUnmarshalBinary(t *testing.T) {
	s := &GSMDiagnostic{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &GSMDiagnostic{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

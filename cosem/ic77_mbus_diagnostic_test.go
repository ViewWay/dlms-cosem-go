package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestMBusDiagnostic_ClassID(t *testing.T) {
	s := &MBusDiagnostic{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDMBusDiagnostic {
		t.Errorf("expected %d, got %d", core.ClassIDMBusDiagnostic, s.ClassID())
	}
}

func TestMBusDiagnostic_MarshalUnmarshalBinary(t *testing.T) {
	s := &MBusDiagnostic{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &MBusDiagnostic{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestLoRaWANDiagnostic_ClassID(t *testing.T) {
	s := &LoRaWANDiagnostic{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDLoRaWANDiagnostic {
		t.Errorf("expected %d, got %d", core.ClassIDLoRaWANDiagnostic, s.ClassID())
	}
}

func TestLoRaWANDiagnostic_MarshalUnmarshalBinary(t *testing.T) {
	s := &LoRaWANDiagnostic{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &LoRaWANDiagnostic{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestDLMSMBusPortSetup_ClassID(t *testing.T) {
	s := &DLMSMBusPortSetup{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDDLMSMBusPortSetup {
		t.Errorf("expected %d, got %d", core.ClassIDDLMSMBusPortSetup, s.ClassID())
	}
}

func TestDLMSMBusPortSetup_MarshalUnmarshalBinary(t *testing.T) {
	s := &DLMSMBusPortSetup{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &DLMSMBusPortSetup{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

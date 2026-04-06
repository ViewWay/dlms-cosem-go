package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestSMTPSetup_ClassID(t *testing.T) {
	s := &SMTPSetup{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDSMTPSetup {
		t.Errorf("expected %d, got %d", core.ClassIDSMTPSetup, s.ClassID())
	}
}

func TestSMTPSetup_MarshalUnmarshalBinary(t *testing.T) {
	s := &SMTPSetup{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &SMTPSetup{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestDataProtection_ClassID(t *testing.T) {
	s := &COSEMDataProtection{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDCOSEMDataProtection {
		t.Errorf("expected %d, got %d", core.ClassIDCOSEMDataProtection, s.ClassID())
	}
}

func TestDataProtection_MarshalUnmarshalBinary(t *testing.T) {
	s := &COSEMDataProtection{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &COSEMDataProtection{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

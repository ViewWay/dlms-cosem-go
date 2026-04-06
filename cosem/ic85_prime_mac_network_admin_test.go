package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestPrimeMACNetworkAdmin_ClassID(t *testing.T) {
	s := &PrimeMACNetworkAdmin{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDPrimeMACNetworkAdmin {
		t.Errorf("expected %d, got %d", core.ClassIDPrimeMACNetworkAdmin, s.ClassID())
	}
}

func TestPrimeMACNetworkAdmin_MarshalUnmarshalBinary(t *testing.T) {
	s := &PrimeMACNetworkAdmin{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &PrimeMACNetworkAdmin{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestTokenGateway_ClassID(t *testing.T) {
	s := &TokenGateway{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDTokenGateway {
		t.Errorf("expected %d, got %d", core.ClassIDTokenGateway, s.ClassID())
	}
}

func TestTokenGateway_MarshalUnmarshalBinary(t *testing.T) {
	s := &TokenGateway{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &TokenGateway{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

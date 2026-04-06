package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestAssociationLN_ClassID(t *testing.T) {
	a := &AssociationLN{}
	if a.ClassID() != 15 {
		t.Errorf("ClassID() = %d, want 15", a.ClassID())
	}
	if a.ClassID() != core.ClassIDAssociationLN {
		t.Error("ClassID mismatch with const")
	}
}

func TestAssociationLN_New(t *testing.T) {
	a := &AssociationLN{
		LogicalName: core.ObisCode{0, 0, 40, 0, 0, 255},
	}
	if a.ClientSAP != 0 {
		t.Error("ClientSAP should default 0")
	}
	if a.ServerSAP != 0 {
		t.Error("ServerSAP should default 0")
	}
}

func TestAssociationLN_MarshalBinary(t *testing.T) {
	a := &AssociationLN{}
	b, err := a.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	a2 := &AssociationLN{}
	if err := a2.UnmarshalBinary(b); err != nil {
		t.Fatal(err)
	}
}

func TestAssociationLN_Fields(t *testing.T) {
	a := &AssociationLN{
		LogicalName: core.ObisCode{0, 0, 40, 0, 0, 255},
		Version:     2,
	}
	if a.GetVersion() != 2 {
		t.Error("GetVersion mismatch")
	}
}

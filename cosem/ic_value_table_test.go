package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestValueTable_ClassID(t *testing.T) {
	vt := &ValueTable{}
	if vt.ClassID() != 258 {
		t.Errorf("ClassID() = %d, want 258", vt.ClassID())
	}
	if vt.ClassID() != core.ClassIDValueTable {
		t.Error("ClassID mismatch with const")
	}
}

func TestValueTable_New(t *testing.T) {
	vt := &ValueTable{
		LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255},
	}
	if vt.Values != nil {
		t.Error("Values should default nil")
	}
	if vt.Descriptors != nil {
		t.Error("Descriptors should default nil")
	}
}

func TestValueTable_MarshalBinary(t *testing.T) {
	vt := &ValueTable{
		Values: []ValueEntry{
			{Index: 1, Value: core.DoubleLongUnsignedData(42)},
		},
	}
	b, err := vt.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	vt2 := &ValueTable{}
	if err := vt2.UnmarshalBinary(b); err != nil {
		t.Fatal(err)
	}
}

func TestValueTable_Fields(t *testing.T) {
	vt := &ValueTable{
		LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255},
		Descriptors: []ValueDescriptor{
			{Index: 0, Description: "Voltage", Unit: 27, Scaler: -1},
		},
		Version: 0,
	}
	if len(vt.Descriptors) != 1 {
		t.Error("Descriptors length mismatch")
	}
	if vt.Descriptors[0].Description != "Voltage" {
		t.Error("Descriptor description mismatch")
	}
}

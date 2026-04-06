package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestPushSetup_ClassID(t *testing.T) {
	ps := &PushSetup{}
	if ps.ClassID() != 40 {
		t.Errorf("ClassID() = %d, want 40", ps.ClassID())
	}
	if ps.ClassID() != core.ClassIDPushSetup {
		t.Error("ClassID mismatch with const")
	}
}

func TestPushSetup_New(t *testing.T) {
	ps := &PushSetup{
		LogicalName:    core.ObisCode{0, 0, 41, 0, 0, 255},
		Service:        1,
		NumberOfRetries: 3,
		Version:        1,
	}
	if ps.Service != 1 {
		t.Error("Service mismatch")
	}
	if ps.NumberOfRetries != 3 {
		t.Error("NumberOfRetries mismatch")
	}
}

func TestPushSetup_MarshalBinary(t *testing.T) {
	ps := &PushSetup{
		PushObjectList: []PushObject{
			{ClassID: 1, LogicalName: core.ObisCode{0, 0, 1, 0, 0, 255}, Attribute: 2},
		},
	}
	b, err := ps.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	ps2 := &PushSetup{}
	if err := ps2.UnmarshalBinary(b); err != nil {
		t.Fatal(err)
	}
}

func TestPushSetup_Fields(t *testing.T) {
	ps := &PushSetup{
		LogicalName: core.ObisCode{0, 0, 41, 0, 0, 255},
		Destination: []byte{0x01, 0x02},
		Version:     2,
	}
	if ps.GetVersion() != 2 {
		t.Error("GetVersion mismatch")
	}
	if len(ps.Destination) != 2 {
		t.Error("Destination length mismatch")
	}
}

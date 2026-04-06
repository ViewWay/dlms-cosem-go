package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestDisconnectControl_ClassID(t *testing.T) {
	dc := &DisconnectControl{}
	if dc.ClassID() != 70 {
		t.Errorf("ClassID() = %d, want 70", dc.ClassID())
	}
	if dc.ClassID() != core.ClassIDDisconnectControl {
		t.Error("ClassID mismatch with const")
	}
}

func TestDisconnectControl_New(t *testing.T) {
	dc := &DisconnectControl{
		LogicalName: core.ObisCode{0, 0, 96, 0, 0, 255},
	}
	if dc.ControlState != 0 {
		t.Error("ControlState should default 0")
	}
}

func TestDisconnectControl_MarshalBinary(t *testing.T) {
	tests := []struct {
		name  string
		state uint8
	}{
		{"disconnected", 0},
		{"connected", 1},
		{"armed", 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dc := &DisconnectControl{ControlState: DisconnectState(tt.state)}
			b, err := dc.MarshalBinary()
			if err != nil {
				t.Fatal(err)
			}
			dc2 := &DisconnectControl{}
			if err := dc2.UnmarshalBinary(b); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestDisconnectControl_Fields(t *testing.T) {
	dc := &DisconnectControl{
		LogicalName:   core.ObisCode{0, 0, 96, 0, 0, 255},
		ControlState:  2,
		OutputState:   1,
		Version:       1,
	}
	if dc.OutputState != 1 {
		t.Error("OutputState mismatch")
	}
}

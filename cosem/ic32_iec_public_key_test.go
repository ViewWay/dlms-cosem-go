package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestIECPublicKey_ClassID(t *testing.T) {
	pk := &IECPublicKey{}
	if pk.ClassID() != 256 {
		t.Errorf("ClassID() = %d, want 256", pk.ClassID())
	}
	if pk.ClassID() != core.ClassIDIECPublicKey {
		t.Error("ClassID mismatch with const")
	}
}

func TestIECPublicKey_New(t *testing.T) {
	pk := &IECPublicKey{
		LogicalName: core.ObisCode{0, 0, 43, 0, 0, 255},
	}
	if pk.PublicKeyValue != nil {
		t.Error("default PublicKeyValue should be nil")
	}
}

func TestIECPublicKey_MarshalBinary(t *testing.T) {
	tests := []struct {
		name string
		key  []byte
	}{
		{"empty", []byte{}},
		{"small", []byte{0x01, 0x02, 0x03}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pk := &IECPublicKey{PublicKeyValue: tt.key}
			b, err := pk.MarshalBinary()
			if err != nil {
				t.Fatal(err)
			}
			pk2 := &IECPublicKey{}
			if err := pk2.UnmarshalBinary(b); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestIECPublicKey_Fields(t *testing.T) {
	pk := &IECPublicKey{
		LogicalName:    core.ObisCode{0, 0, 43, 0, 0, 255},
		PublicKeyValue: []byte{0xDE, 0xAD, 0xBE, 0xEF},
		Version:        0,
	}
	if len(pk.PublicKeyValue) != 4 {
		t.Error("PublicKeyValue length mismatch")
	}
}

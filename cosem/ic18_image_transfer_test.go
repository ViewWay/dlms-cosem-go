package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestImageTransfer_ClassID(t *testing.T) {
	it := &ImageTransfer{}
	if it.ClassID() != 18 {
		t.Errorf("ClassID() = %d, want 18", it.ClassID())
	}
	if it.ClassID() != core.ClassIDImageTransfer {
		t.Error("ClassID mismatch with const")
	}
}

func TestImageTransfer_New(t *testing.T) {
	it := &ImageTransfer{
		LogicalName: core.ObisCode{0, 0, 44, 0, 0, 255},
	}
	if it.ImageBlockSize != 0 {
		t.Error("ImageBlockSize should default 0")
	}
}

func TestImageTransfer_MarshalBinary(t *testing.T) {
	tests := []struct {
		name   string
		status uint8
	}{
		{"idle", 0},
		{"initiated", 1},
		{"verifying", 6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			it := &ImageTransfer{TransferStatus: ImageTransferStatus(tt.status)}
			b, err := it.MarshalBinary()
			if err != nil {
				t.Fatal(err)
			}
			it2 := &ImageTransfer{}
			if err := it2.UnmarshalBinary(b); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestImageTransfer_Fields(t *testing.T) {
	it := &ImageTransfer{
		LogicalName:   core.ObisCode{0, 0, 44, 0, 0, 255},
		ImageBlockSize: 200,
		Version:       1,
	}
	if it.ImageBlockSize != 200 {
		t.Error("ImageBlockSize mismatch")
	}
}

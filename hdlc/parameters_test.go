package hdlc

import (
	"testing"
)

func TestHdlcParameterTypeConstants(t *testing.T) {
	// Test parameter IDs according to Green Book
	if ParamMaxInfoFieldLengthTX != 0x05 {
		t.Errorf("ParamMaxInfoFieldLengthTX should be 0x05, got %02x", ParamMaxInfoFieldLengthTX)
	}
	if ParamMaxInfoFieldLengthRX != 0x06 {
		t.Errorf("ParamMaxInfoFieldLengthRX should be 0x06, got %02x", ParamMaxInfoFieldLengthRX)
	}
	if ParamWindowSizeTX != 0x07 {
		t.Errorf("ParamWindowSizeTX should be 0x07, got %02x", ParamWindowSizeTX)
	}
	if ParamWindowSizeRX != 0x08 {
		t.Errorf("ParamWindowSizeRX should be 0x08, got %02x", ParamWindowSizeRX)
	}
}

func TestHdlcParametersGreenValidation(t *testing.T) {
	tests := []struct {
		name    string
		params  HdlcParametersGreen
		wantErr bool
	}{
		{
			name: "valid defaults",
			params: HdlcParametersGreen{
				MaxInfoLengthTX: 128,
				MaxInfoLengthRX: 128,
				WindowSizeTX:    1,
				WindowSizeRX:    1,
			},
			wantErr: false,
		},
		{
			name: "valid max values",
			params: HdlcParametersGreen{
				MaxInfoLengthTX: 2048,
				MaxInfoLengthRX: 2048,
				WindowSizeTX:    7,
				WindowSizeRX:    7,
			},
			wantErr: false,
		},
		{
			name: "invalid window size TX below minimum",
			params: HdlcParametersGreen{
				WindowSizeTX: 0,
			},
			wantErr: true,
		},
		{
			name: "invalid window size TX above maximum",
			params: HdlcParametersGreen{
				WindowSizeTX: 8,
			},
			wantErr: true,
		},
		{
			name: "invalid max info length below minimum",
			params: HdlcParametersGreen{
				MaxInfoLengthTX: 127,
			},
			wantErr: true,
		},
		{
			name: "invalid max info length above maximum",
			params: HdlcParametersGreen{
				MaxInfoLengthTX: 2049,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.params.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHdlcParametersGreenEncoding(t *testing.T) {
	params := &HdlcParametersGreen{
		MaxInfoLengthTX: 128,
		MaxInfoLengthRX: 128,
		WindowSizeTX:    1,
		WindowSizeRX:    1,
	}

	encoded, err := params.Encode()
	if err != nil {
		t.Fatalf("Encode() failed: %v", err)
	}

	// Check header
	if len(encoded) < 3 {
		t.Fatalf("encoded data too short: %d bytes", len(encoded))
	}
	if encoded[0] != FormatIdentifier {
		t.Errorf("format identifier should be %02x, got %02x", FormatIdentifier, encoded[0])
	}
	if encoded[1] != GroupIdentifier {
		t.Errorf("group identifier should be %02x, got %02x", GroupIdentifier, encoded[1])
	}

	// Check group length
	groupLength := int(encoded[2])
	if len(encoded) != 3+groupLength {
		t.Errorf("encoded length %d doesn't match 3 + group length %d", len(encoded), groupLength)
	}
}

func TestHdlcParametersGreenRoundtrip(t *testing.T) {
	original := &HdlcParametersGreen{
		MaxInfoLengthTX: 512,
		MaxInfoLengthRX: 1024,
		WindowSizeTX:    3,
		WindowSizeRX:    5,
	}

	encoded, err := original.Encode()
	if err != nil {
		t.Fatalf("Encode() failed: %v", err)
	}

	decoded, err := ParseHdlcParametersGreen(encoded)
	if err != nil {
		t.Fatalf("ParseHdlcParametersGreen() failed: %v", err)
	}

	if decoded.MaxInfoLengthTX != original.MaxInfoLengthTX {
		t.Errorf("MaxInfoLengthTX mismatch: got %d, want %d", decoded.MaxInfoLengthTX, original.MaxInfoLengthTX)
	}
	if decoded.MaxInfoLengthRX != original.MaxInfoLengthRX {
		t.Errorf("MaxInfoLengthRX mismatch: got %d, want %d", decoded.MaxInfoLengthRX, original.MaxInfoLengthRX)
	}
	if decoded.WindowSizeTX != original.WindowSizeTX {
		t.Errorf("WindowSizeTX mismatch: got %d, want %d", decoded.WindowSizeTX, original.WindowSizeTX)
	}
	if decoded.WindowSizeRX != original.WindowSizeRX {
		t.Errorf("WindowSizeRX mismatch: got %d, want %d", decoded.WindowSizeRX, original.WindowSizeRX)
	}
}

func TestHdlcParametersGreenLargeValues(t *testing.T) {
	// Test with values > 255 to ensure 2-byte encoding
	params := &HdlcParametersGreen{
		MaxInfoLengthTX: 2048,
		MaxInfoLengthRX: 1024,
		WindowSizeTX:    7,
		WindowSizeRX:    7,
	}

	encoded, err := params.Encode()
	if err != nil {
		t.Fatalf("Encode() failed: %v", err)
	}

	decoded, err := ParseHdlcParametersGreen(encoded)
	if err != nil {
		t.Fatalf("ParseHdlcParametersGreen() failed: %v", err)
	}

	if decoded.MaxInfoLengthTX != 2048 {
		t.Errorf("MaxInfoLengthTX: got %d, want 2048", decoded.MaxInfoLengthTX)
	}
	if decoded.MaxInfoLengthRX != 1024 {
		t.Errorf("MaxInfoLengthRX: got %d, want 1024", decoded.MaxInfoLengthRX)
	}
}

func TestNegotiateParametersGreen(t *testing.T) {
	client := &HdlcParametersGreen{
		MaxInfoLengthTX: 1024,
		MaxInfoLengthRX: 2048,
		WindowSizeTX:    5,
		WindowSizeRX:    7,
	}

	server := &HdlcParametersGreen{
		MaxInfoLengthTX: 512,
		MaxInfoLengthRX: 512,
		WindowSizeTX:    7,
		WindowSizeRX:    3,
	}

	negotiated := NegotiateParametersGreen(client, server)

	// Client TX (1024) with server RX (512) -> should be 512
	if negotiated.MaxInfoLengthTX != 512 {
		t.Errorf("MaxInfoLengthTX: got %d, want 512", negotiated.MaxInfoLengthTX)
	}

	// Client TX window (5) with server RX window (3) -> should be 3
	if negotiated.WindowSizeTX != 3 {
		t.Errorf("WindowSizeTX: got %d, want 3", negotiated.WindowSizeTX)
	}

	// Client RX (2048) with server TX (512) -> should be 512
	if negotiated.MaxInfoLengthRX != 512 {
		t.Errorf("MaxInfoLengthRX: got %d, want 512", negotiated.MaxInfoLengthRX)
	}

	// Client RX window (7) with server TX window (7) -> should be 7
	if negotiated.WindowSizeRX != 7 {
		t.Errorf("WindowSizeRX: got %d, want 7", negotiated.WindowSizeRX)
	}
}

func TestGreenBookExampleEncoding(t *testing.T) {
	// Test against Green Book example from section 8.4.5.3.2
	// Example with default values (1-byte for max info length)
	params := &HdlcParametersGreen{
		MaxInfoLengthTX: 128, // 0x80
		MaxInfoLengthRX: 128, // 0x80
		WindowSizeTX:    1,
		WindowSizeRX:    1,
	}

	encoded, err := params.Encode()
	if err != nil {
		t.Fatalf("Encode() failed: %v", err)
	}

	// Check that encoding starts with correct header
	if encoded[0] != 0x81 || encoded[1] != 0x80 {
		t.Errorf("Header mismatch: got %02x %02x, want 81 80", encoded[0], encoded[1])
	}

	// Parse and verify
	decoded, err := ParseHdlcParametersGreen(encoded)
	if err != nil {
		t.Fatalf("ParseHdlcParametersGreen() failed: %v", err)
	}

	if decoded.MaxInfoLengthTX != 128 {
		t.Errorf("MaxInfoLengthTX: got %d, want 128", decoded.MaxInfoLengthTX)
	}
}

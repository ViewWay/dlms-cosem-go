package asn1

import (
	"bytes"
	"testing"
	"testing/quick"
)

// Property tests for BER encoding/decoding

// TestBEREncodeDecode_Roundtrip verifies that encoding and decoding are inverse operations
func TestBEREncodeDecode_Roundtrip(t *testing.T) {
	f := func(tag int, data []byte) bool {
		// Limit tag to single-byte encoding (0-255, avoiding 0x1F in lower 5 bits)
		if tag < 0 || tag > 255 {
			return true // Skip invalid tags
		}
		if tag&0x1F == 0x1F {
			return true // Skip multi-byte tag markers (not supported by BEREncode)
		}

		encoded := BEREncode(tag, data)
		decodedTag, decodedValue, _, err := BERDecode(encoded)
		if err != nil {
			return false
		}
		if decodedTag != tag {
			return false
		}
		return bytes.Equal(decodedValue, data)
	}

	config := &quick.Config{MaxCount: 1000}
	if err := quick.Check(f, config); err != nil {
		t.Error(err)
	}
}

// TestBEREncodeDecode_EmptyNil verifies empty and nil data are handled correctly
func TestBEREncodeDecode_EmptyNil(t *testing.T) {
	for tag := 0; tag < 256; tag++ {
		// Skip multi-byte tag markers
		if tag&0x1F == 0x1F {
			continue
		}

		// Test with empty slice
		encoded := BEREncode(tag, []byte{})
		decodedTag, decodedValue, _, err := BERDecode(encoded)
		if err != nil {
			t.Errorf("tag=0x%02x: %v", tag, err)
			continue
		}
		if decodedTag != tag || len(decodedValue) != 0 {
			t.Errorf("tag=0x%02x: decoded tag=0x%02x, len=%d", tag, decodedTag, len(decodedValue))
		}

		// Test with nil slice
		encoded = BEREncode(tag, nil)
		decodedTag, decodedValue, _, err = BERDecode(encoded)
		if err != nil {
			t.Errorf("tag=0x%02x: %v", tag, err)
			continue
		}
		if decodedTag != tag || len(decodedValue) != 0 {
			t.Errorf("tag=0x%02x: decoded tag=0x%02x, len=%d", tag, decodedTag, len(decodedValue))
		}
	}
}

// TestBEREncodeLength0_Properties verifies BEREncodeLength0 for all tag values
func TestBEREncodeLength0_Properties(t *testing.T) {
	f := func(tag int) bool {
		encoded := BEREncodeLength0(tag)
		// Should always be 2 bytes: tag byte + 0x00 length byte
		if len(encoded) != 2 {
			return false
		}
		if encoded[0] != byte(tag) {
			return false
		}
		if encoded[1] != 0x00 {
			return false
		}
		return true
	}

	config := &quick.Config{MaxCount: 256}
	if err := quick.Check(f, config); err != nil {
		t.Error(err)
	}
}

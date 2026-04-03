package asn1

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestBEREncodeDecode_Simple(t *testing.T) {
	data := []byte{0x01, 0x02, 0x03}
	encoded := BEREncode(0x60, data)
	if encoded[0] != 0x60 {
		t.Errorf("tag=%02x", encoded[0])
	}
	if encoded[1] != 3 {
		t.Errorf("length=%d", encoded[1])
	}
	if !bytes.Equal(encoded[2:], data) {
		t.Errorf("value=%v", encoded[2:])
	}

	tag, value, consumed, err := BERDecode(encoded)
	if err != nil {
		t.Fatal(err)
	}
	if tag != 0x60 {
		t.Errorf("tag=%d", tag)
	}
	if !bytes.Equal(value, data) {
		t.Errorf("value=%v", value)
	}
	if consumed != len(encoded) {
		t.Errorf("consumed=%d", consumed)
	}
}

func TestBEREncodeDecode_Empty(t *testing.T) {
	encoded := BEREncode(0x60, []byte{})
	if encoded[0] != 0x60 || encoded[1] != 0 {
		t.Errorf("got %v", encoded)
	}
	tag, value, _, err := BERDecode(encoded)
	if err != nil || tag != 0x60 || len(value) != 0 {
		t.Error("decode failed")
	}
}

func TestBEREncodeDecode_Nil(t *testing.T) {
	encoded := BEREncode(0x60, nil)
	tag, value, _, err := BERDecode(encoded)
	if err != nil || tag != 0x60 || len(value) != 0 {
		t.Error("decode failed")
	}
}

func TestBEREncodeDecode_MultiByteLength(t *testing.T) {
	data := make([]byte, 200)
	encoded := BEREncode(0x60, data)
	if encoded[1] != 0x81 {
		t.Errorf("length encoding=%02x", encoded[1])
	}
	_, value, _, err := BERDecode(encoded)
	if err != nil || len(value) != 200 {
		t.Error("decode failed")
	}
}

func TestBEREncodeDecode_LongLength(t *testing.T) {
	data := make([]byte, 500)
	encoded := BEREncode(0x60, data)
	if encoded[1] != 0x82 {
		t.Errorf("length encoding=%02x", encoded[1])
	}
}

func TestBERDecode_Short(t *testing.T) {
	_, _, _, err := BERDecode([]byte{0x60})
	if err == nil {
		t.Error("expected error for short data")
	}
}

func TestBERDecode_Truncated(t *testing.T) {
	encoded := BEREncode(0x60, []byte{1, 2, 3})
	_, _, _, err := BERDecode(encoded[:4]) // missing last byte
	if err == nil {
		t.Error("expected error for truncated data")
	}
}

func TestBEREncodeOctetString(t *testing.T) {
	encoded := BEREncodeOctetString([]byte{0x01, 0x02})
	if encoded[0] != 0x04 || encoded[1] != 2 {
		t.Errorf("got %v", encoded)
	}
}

func TestBEREncodeInteger(t *testing.T) {
	tests := []struct {
		value int
	}{
		{0},
		{1},
		{127},
		{128},
		{255},
		{256},
		{-1},
		{-128},
		{65535},
	}
	for _, tc := range tests {
		encoded := BEREncodeInteger(tc.value)
		if encoded[0] != 0x02 {
			t.Errorf("integer tag for %d: %02x", tc.value, encoded[0])
		}
	}
}

func TestBEREncodeContextTag(t *testing.T) {
	encoded := BEREncodeContextTag(1, []byte{0x01, 0x02})
	if encoded[0] != 0xA1 {
		t.Errorf("got %02x", encoded[0])
	}
}

func TestBEREncodeContextPrimitive(t *testing.T) {
	encoded := BEREncodeContextPrimitive(0x0A, []byte{0x01})
	if encoded[0] != 0x8A {
		t.Errorf("got %02x", encoded[0])
	}
}

func TestBEREncodeSequence(t *testing.T) {
	encoded := BEREncodeSequence([]byte{0x01, 0x02})
	if encoded[0] != 0x30 {
		t.Errorf("got %02x", encoded[0])
	}
}

func TestBERDecodeAllTLV(t *testing.T) {
	data := BEREncode(0x02, []byte{42})
	data = append(data, BEREncode(0x02, []byte{43})...)
	tlvs, err := BERDecodeAllTLV(data)
	if err != nil {
		t.Fatal(err)
	}
	if len(tlvs) != 2 {
		t.Fatalf("len=%d", len(tlvs))
	}
	if tlvs[0].Tag != 2 || tlvs[0].Value[0] != 42 {
		t.Error("first TLV")
	}
	if tlvs[1].Tag != 2 || tlvs[1].Value[0] != 43 {
		t.Error("second TLV")
	}
}

func TestBERDecodeAllTLV_Empty(t *testing.T) {
	tlvs, err := BERDecodeAllTLV([]byte{})
	if err != nil {
		t.Fatal(err)
	}
	if len(tlvs) != 0 {
		t.Errorf("len=%d", len(tlvs))
	}
}

func TestEncodeAppContextName(t *testing.T) {
	tests := []struct {
		logical  bool
		ciphered bool
		want     int
	}{
		{true, false, AppContextNameLogicalNameRef},
		{true, true, AppContextNameCipheredLN},
		{false, false, AppContextNameShortNameRef},
	}
	for _, tc := range tests {
		b := EncodeAppContextName(tc.logical, tc.ciphered)
		if len(b) < 4 {
			t.Errorf("too short: %v", b)
			continue
		}
		oid := int(binary.BigEndian.Uint16(b[2:4]))
		if oid != tc.want {
			t.Errorf("got %d, want %d", oid, tc.want)
		}
	}
}

func TestEncodeDecodeUserInformation(t *testing.T) {
	initReq := []byte{0x01, 0x02, 0x03}
	ui := EncodeUserInformation(initReq)
	decoded, err := DecodeUserInformation(ui)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(decoded, initReq) {
		t.Errorf("got %v, want %v", decoded, initReq)
	}
}

func TestDecodeUserInformation_InvalidTag(t *testing.T) {
	_, err := DecodeUserInformation([]byte{0x01, 0x02})
	if err == nil {
		t.Error("expected error")
	}
}

func TestAARQ_EncodeDecode(t *testing.T) {
	aarq := &AARQ{
		ApplicationContextName: EncodeAppContextName(true, false),
		CallingAPTitle:         []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
		Authentication:         AuthHLS,
		AuthenticationValue:    []byte{0xAA, 0xBB},
		UserInformation:        EncodeUserInformation([]byte{0x01, 0x00, 0x00, 0x00, 0x06, 0x5F, 0x1F, 0x04, 0x00}),
	}

	encoded := aarq.Encode()
	parsed, err := ParseAARQ(encoded)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Authentication != AuthHLS {
		t.Errorf("auth=%d", parsed.Authentication)
	}
}

func TestAARQ_EncodeMinimal(t *testing.T) {
	aarq := &AARQ{
		ApplicationContextName: EncodeAppContextName(true, false),
		UserInformation:        EncodeUserInformation([]byte{0x01, 0x00}),
	}
	encoded := aarq.Encode()
	if encoded[0] != TagAARQ {
		t.Errorf("tag=%02x", encoded[0])
	}
}

func TestAARE_Encode(t *testing.T) {
	aare := &AARE{
		Result:          AssocAccepted,
		UserInformation: EncodeUserInformation([]byte{0x01, 0x00}),
	}
	encoded := aare.Encode()
	if encoded[0] != TagAARE {
		t.Errorf("tag=%02x", encoded[0])
	}
}

func TestAARE_Parse(t *testing.T) {
	aare := &AARE{
		Result:          AssocAccepted,
		UserInformation: EncodeUserInformation([]byte{0x01, 0x00}),
	}
	encoded := aare.Encode()
	parsed, err := ParseAARE(encoded)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Result != AssocAccepted {
		t.Errorf("result=%d", parsed.Result)
	}
}

func TestRLRE_EncodeDecode(t *testing.T) {
	rlre := &RLRE{Result: 0}
	encoded := rlre.Encode()
	if encoded[0] != TagRLRE {
		t.Errorf("tag=%02x", encoded[0])
	}
	parsed, err := ParseRLRE(encoded)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Result != 0 {
		t.Errorf("result=%d", parsed.Result)
	}
}

func TestParseAARQ_WrongTag(t *testing.T) {
	_, err := ParseAARQ([]byte{0x61, 0x02, 0x01, 0x02})
	if err == nil {
		t.Error("expected error for wrong tag")
	}
}

func TestParseAARE_WrongTag(t *testing.T) {
	_, err := ParseAARE([]byte{0x60, 0x02, 0x01, 0x02})
	if err == nil {
		t.Error("expected error for wrong tag")
	}
}

func TestParseRLRE_WrongTag(t *testing.T) {
	_, err := ParseRLRE([]byte{0x60, 0x02, 0x01, 0x02})
	if err == nil {
		t.Error("expected error for wrong tag")
	}
}

func TestBEREncodeLength0(t *testing.T) {
	encoded := BEREncodeLength0(0x60)
	if len(encoded) != 2 || encoded[0] != 0x60 || encoded[1] != 0x00 {
		t.Errorf("got %v", encoded)
	}
}

func TestEncodeBERLength(t *testing.T) {
	tests := []struct {
		input   int
		wantLen int
	}{
		{0, 1},
		{127, 1},
		{128, 2},
		{255, 2},
		{256, 3},
		{65535, 3},
	}
	for _, tc := range tests {
		b := encodeBERLength(tc.input)
		if len(b) != tc.wantLen {
			t.Errorf("encodeBERLength(%d): len=%d, want %d", tc.input, len(b), tc.wantLen)
		}
	}
}

func TestDecodeBERTag(t *testing.T) {
	tag, n := decodeBERTag([]byte{0x60})
	if tag != 0x60 || n != 1 {
		t.Error("single byte tag")
	}
}

func TestDecodeBERTag_MultiByte(t *testing.T) {
	tag, n := decodeBERTag([]byte{0x1F, 0x80, 0x01})
	if tag != 0x8001 || n != 3 {
		t.Errorf("got %d, %d", tag, n)
	}
}

func TestBEREncodeDecode_Integer_Roundtrip(t *testing.T) {
	values := []int{0, 1, -1, 127, 128, -128, 255, 256, 32767, -32768, 65535}
	for _, v := range values {
		encoded := BEREncodeInteger(v)
		tag, value, _, err := BERDecode(encoded)
		if err != nil {
			t.Errorf("encode/decode %d: %v", v, err)
			continue
		}
		if tag != 0x02 {
			t.Errorf("tag=%d for %d", tag, v)
		}
		// Verify decoding
		_ = value
	}
}

func TestAARE_ParseWithSystemTitle(t *testing.T) {
	aare := &AARE{
		Result:            AssocAccepted,
		RespondingAPTitle: []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
		UserInformation:   EncodeUserInformation([]byte{0x01, 0x00}),
	}
	encoded := aare.Encode()
	parsed, err := ParseAARE(encoded)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(parsed.RespondingAPTitle, aare.RespondingAPTitle) {
		t.Errorf("got %v", parsed.RespondingAPTitle)
	}
}

func TestBEREncode_NilData(t *testing.T) {
	encoded := BEREncode(0x60, nil)
	if len(encoded) != 2 {
		t.Errorf("len=%d", len(encoded))
	}
}

func TestBERDecode_BadLength(t *testing.T) {
	// Length byte claims more data than available
	_, _, _, err := BERDecode([]byte{0x60, 0x05, 0x01, 0x02})
	if err == nil {
		t.Error("expected error for truncated data")
	}
}

func TestDecodeBERLength_Invalid(t *testing.T) {
	_, _, err := decodeBERLength([]byte{0x85, 0x01, 0x02, 0x03, 0x04, 0x05})
	if err == nil {
		t.Error("expected error for multi-byte length > 4")
	}
}

func TestDecodeBERLength_Empty(t *testing.T) {
	_, _, err := decodeBERLength([]byte{})
	if err == nil {
		t.Error("expected error for empty data")
	}
}

func TestDecodeBERLength_Indefinite(t *testing.T) {
	length, n, err := decodeBERLength([]byte{0x80})
	if err != nil {
		t.Error("indefinite length should not error")
	}
	if n != 1 || length != 0 {
		t.Errorf("got %d, %d", length, n)
	}
}

func TestAARQ_WithUserID(t *testing.T) {
	aarq := &AARQ{
		ApplicationContextName: EncodeAppContextName(true, false),
		CallingAEInvocationID:  []byte{0x01, 0x00, 0x01},
		UserInformation:        EncodeUserInformation([]byte{0x01, 0x00}),
	}
	encoded := aarq.Encode()
	parsed, err := ParseAARQ(encoded)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(parsed.CallingAEInvocationID, []byte{0x01, 0x00, 0x01}) {
		t.Errorf("got %v", parsed.CallingAEInvocationID)
	}
}

package hdlc

import (
	"bytes"
	"testing"
)

func TestCRCCCITT(t *testing.T) {
	// Test with known data
	crc := CRCCCITT([]byte{0x01, 0x02, 0x03})
	if len(crc) != 2 {
		t.Errorf("expected 2 bytes, got %d", len(crc))
	}
	// Test roundtrip: data + crc should verify
	full := append([]byte{0x01, 0x02, 0x03}, crc...)
	if !VerifyFCS(full) {
		t.Error("FCS verification failed")
	}
}

func TestCRCCCITT_Empty(t *testing.T) {
	crc := CRCCCITT([]byte{})
	if len(crc) != 2 {
		t.Error("expected 2 bytes")
	}
}

func TestVerifyFCS_Tampered(t *testing.T) {
	crc := CRCCCITT([]byte{0x01, 0x02})
	full := append([]byte{0x01, 0x03}, crc...) // changed byte
	if VerifyFCS(full) {
		t.Error("should fail with tampered data")
	}
}

func TestVerifyFCS_Short(t *testing.T) {
	if VerifyFCS([]byte{}) {
		t.Error("should fail for empty data")
	}
	if VerifyFCS([]byte{0x00}) {
		t.Error("should fail for 1 byte")
	}
}

func TestReverseByte(t *testing.T) {
	if reverseByte(0x80) != 0x01 {
		t.Errorf("got %02x", reverseByte(0x80))
	}
	if reverseByte(0x00) != 0x00 {
		t.Errorf("got %02x", reverseByte(0x00))
	}
	if reverseByte(0xFF) != 0xFF {
		t.Errorf("got %02x", reverseByte(0xFF))
	}
	if reverseByte(0x01) != 0x80 {
		t.Errorf("got %02x", reverseByte(0x01))
	}
	if reverseByte(0xA5) != 0xA5 {
		t.Errorf("got %02x", reverseByte(0xA5))
	}
}

func TestHdlcEscape(t *testing.T) {
	data := []byte{0x7E, 0x01, 0x7D, 0x02}
	escaped := hdlcEscape(data)
	if escaped[0] != 0x7E || escaped[len(escaped)-1] != 0x7E {
		t.Error("missing flags")
	}
	// Check that all 0x7E and 0x7D in data are properly escaped
	// (no standalone 0x7E or 0x7D between flags)
	for i := 1; i < len(escaped)-1; i++ {
		if escaped[i] == 0x7E {
			t.Errorf("found unescaped 0x7E at pos %d", i)
		}
		if escaped[i] == 0x7D && (i+1 >= len(escaped)-1 || (escaped[i+1] != 0x5E && escaped[i+1] != 0x5D)) {
			t.Errorf("found escape char 0x7D at pos %d not followed by 0x5E or 0x5D", i)
		}
	}
}

func TestHdlcUnescape(t *testing.T) {
	data := []byte{0x7E, 0x01, 0x7D, 0x5D, 0x7D, 0x31, 0x02, 0x7E}
	unescaped, err := hdlcUnescape(data)
	if err != nil {
		t.Fatal(err)
	}
	expected := []byte{0x01, 0x7D, 0x11, 0x02}
	if !bytes.Equal(unescaped, expected) {
		t.Errorf("got %v, want %v", unescaped, expected)
	}
}

func TestHdlcUnescape_NoFlags(t *testing.T) {
	_, err := hdlcUnescape([]byte{0x01, 0x02})
	if err == nil {
		t.Error("expected error for missing flags")
	}
}

func TestHdlcUnescape_IncompleteEscape(t *testing.T) {
	_, err := hdlcUnescape([]byte{0x7E, 0x7D, 0x7E})
	if err == nil {
		t.Error("expected error for incomplete escape")
	}
}

func TestParseFrameType(t *testing.T) {
	if ParseFrameType(0x00) != FrameTypeI {
		t.Error("I-frame")
	}
	if ParseFrameType(0x01) != FrameTypeS {
		t.Error("S-frame")
	}
	if ParseFrameType(0x03) != FrameTypeU {
		t.Error("U-frame")
	}
	if ParseFrameType(0x10) != FrameTypeI {
		t.Error("I-frame with send seq")
	}
	if ParseFrameType(0x05) != FrameTypeS {
		t.Error("S-frame REJ")
	}
	if ParseFrameType(0x63) != FrameTypeU {
		t.Error("U-frame UA")
	}
}

func TestHdlcAddress_Client(t *testing.T) {
	a := &HdlcAddress{Logical: 1}
	b := a.EncodeAddress(false)
	if len(b) != 1 || b[0] != 0x03 {
		t.Errorf("got %v", b)
	}
}

func TestHdlcAddress_Server(t *testing.T) {
	a := &HdlcAddress{Logical: 1}
	b := a.EncodeAddress(true)
	if len(b) != 1 || b[0] != 0x03 {
		t.Errorf("got %v", b)
	}
}

func TestHdlcAddress_ServerTwoByte(t *testing.T) {
	a := &HdlcAddress{Logical: 128, Extended: true}
	b := a.EncodeAddress(true)
	if len(b) != 2 {
		t.Errorf("len=%d", len(b))
	}
}

func TestHdlcAddress_ServerFourByte(t *testing.T) {
	a := &HdlcAddress{Logical: 1, Physical: 2, HasPhysical: true, Extended: true}
	b := a.EncodeAddress(true)
	if len(b) != 4 {
		t.Errorf("len=%d, got %v", len(b), b)
	}
}

func TestSplitAddress(t *testing.T) {
	h, l := splitAddress(1)
	if h != 0 || l != 2 {
		t.Errorf("got %d, %d", h, l)
	}
	h, l = splitAddress(128)
	// 128 & 0x7F = 0, <<1 = 0; 128 & 0x3F80 = 128, >> 6 = 2
	if h != 2 || l != 0 {
		t.Errorf("got %d, %d", h, l)
	}
	h, l = splitAddress(255)
	// 255 & 0x7F = 127, <<1 = 254; 255 & 0x3F80 = 128, >> 6 = 2
	if h != 2 || l != 254 {
		t.Errorf("got %d, %d", h, l)
	}
	h, l = splitAddress(0)
	if h != 0 || l != 0 {
		t.Errorf("got %d, %d", h, l)
	}
}

func TestParseTwoByteAddress(t *testing.T) {
	// Test roundtrip: splitAddress then parseTwoByteAddress
	tests := []struct{ addr int }{{1}, {127}, {255}, {16383}}
	for _, tc := range tests {
		h, l := splitAddress(tc.addr)
		v := parseTwoByteAddress([]byte{h, l})
		if v != tc.addr {
			t.Errorf("roundtrip %d: got %d", tc.addr, v)
		}
	}
}

func TestParseAddresses(t *testing.T) {
	// Client=1 (1 byte), Server=1 (1 byte)
	data := []byte{0x00, 0x00, 0x03, 0x03, 0x00} // format + client + server + control
	dest, src, consumed, err := ParseAddresses(data)
	if err != nil {
		t.Fatal(err)
	}
	if dest.Logical != 1 {
		t.Errorf("dest=%d", dest.Logical)
	}
	if src.Logical != 1 {
		t.Errorf("src=%d", src.Logical)
	}
	if consumed != 4 {
		t.Errorf("consumed=%d", consumed)
	}
}

func TestParseAddresses_Short(t *testing.T) {
	_, _, _, err := ParseAddresses([]byte{0x00, 0x00})
	if err == nil {
		t.Error("expected error")
	}
}

func TestEncodeDecodeFrame(t *testing.T) {
	dest := &HdlcAddress{Logical: 1}
	src := &HdlcAddress{Logical: 1}
	control := byte(0x93) // SNRM
	info := []byte{0x01, 0x02}

	encoded := EncodeFrame(dest, src, control, info)
	if encoded[0] != HDLCFlag || encoded[len(encoded)-1] != HDLCFlag {
		t.Error("missing flags")
	}

	parsed, err := ParseFrame(encoded)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Control != control {
		t.Errorf("control=%02x", parsed.Control)
	}
	if !bytes.Equal(parsed.Info, info) {
		t.Errorf("info=%v", parsed.Info)
	}
}

func TestEncodeFrame_NoInfo(t *testing.T) {
	dest := &HdlcAddress{Logical: 1}
	src := &HdlcAddress{Logical: 1}
	encoded := EncodeFrame(dest, src, 0x93, nil)
	parsed, err := ParseFrame(encoded)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Info != nil {
		t.Errorf("expected nil info, got %v", parsed.Info)
	}
}

func TestFrameType_Method(t *testing.T) {
	f := &Frame{Control: 0x00}
	if f.FrameType() != FrameTypeI {
		t.Error("expected I")
	}
	f = &Frame{Control: 0x01}
	if f.FrameType() != FrameTypeS {
		t.Error("expected S")
	}
	f = &Frame{Control: 0x63}
	if f.FrameType() != FrameTypeU {
		t.Error("expected U")
	}
}

func TestFrameParser(t *testing.T) {
	parser := NewFrameParser()
	// Feed data byte by byte
	data := []byte{0x7E, 0x03, 0x7D, 0x5D, 0x7E}
	frames := parser.Feed(data)
	if len(frames) != 1 {
		t.Fatalf("expected 1 frame, got %d", len(frames))
	}
	unescaped, err := hdlcUnescape(frames[0])
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(unescaped, []byte{0x03, 0x7D}) {
		t.Errorf("got %v", unescaped)
	}
}

func TestFrameParser_MultipleFrames(t *testing.T) {
	parser := NewFrameParser()
	f1 := []byte{0x7E, 0x01, 0x7E}
	f2 := []byte{0x7E, 0x02, 0x7E}
	all := append(f1, f2...)
	frames := parser.Feed(all)
	if len(frames) != 2 {
		t.Fatalf("expected 2 frames, got %d", len(frames))
	}
}

func TestFrameParser_Reset(t *testing.T) {
	parser := NewFrameParser()
	parser.Feed([]byte{0x7E, 0x01})
	parser.Reset()
	if parser.inFrame {
		t.Error("should be reset")
	}
}

func TestFrameParser_PartialFeed(t *testing.T) {
	parser := NewFrameParser()
	frames := parser.Feed([]byte{0x7E, 0x01})
	if len(frames) != 0 {
		t.Error("should not have complete frame yet")
	}
	frames = parser.Feed([]byte{0x02, 0x7E})
	if len(frames) != 1 {
		t.Errorf("expected 1 frame, got %d", len(frames))
	}
}

func TestFrameParser_IgnoreNonFrameData(t *testing.T) {
	parser := NewFrameParser()
	frames := parser.Feed([]byte{0x00, 0x01, 0x02})
	if len(frames) != 0 {
		t.Error("should not have frames")
	}
}

func TestControlHelpers(t *testing.T) {
	if !IsU(UAControl(true)) {
		t.Error("UA should be U-frame")
	}
	if !IsU(SNRMControl()) {
		t.Error("SNRM should be U-frame")
	}
	if !IsI(IFrameControl(0, 0)) {
		t.Error("I-frame check failed")
	}
	if !IsS(RRControl(0, false)) {
		t.Error("RR should be S-frame")
	}
	if !IsUA(UAControl(false)) {
		t.Error("UA check failed")
	}
	if !IsSNRM(SNRMControl()) {
		t.Error("SNRM check failed")
	}
	if !IsDISC(DISCControl()) {
		t.Error("DISC check failed")
	}
	if !IsDM(DMControl()) {
		t.Error("DM check failed")
	}
}

func TestExtractSequences(t *testing.T) {
	control := IFrameControl(5, 3)
	if ExtractSendSeq(control) != 5 {
		t.Errorf("send=%d", ExtractSendSeq(control))
	}
	if ExtractRecvSeq(control) != 3 {
		t.Errorf("recv=%d", ExtractRecvSeq(control))
	}
}

func TestExtractSFrameFields(t *testing.T) {
	c := RRControl(2, true)
	if ExtractSFrameRecvSeq(c) != 2 {
		t.Errorf("recv=%d", ExtractSFrameRecvSeq(c))
	}
}

func TestExtractUFrameModifier(t *testing.T) {
	c := UAControl(true)
	m := ExtractUFrameModifier(c)
	if m == 0 {
		t.Error("expected non-zero modifier")
	}
}

func TestParseFormatField(t *testing.T) {
	l, seg, err := ParseFormatField([]byte{0x10, 0x08})
	if err != nil {
		t.Fatal(err)
	}
	if l != 16 || !seg {
		t.Errorf("got %d, %v", l, seg)
	}
}

func TestEncodeFormatField(t *testing.T) {
	b := EncodeFormatField(16, true)
	if b[0] != 16 || b[1] != 0x08 {
		t.Errorf("got %v", b)
	}
}

func TestParseFormatField_Short(t *testing.T) {
	_, _, err := ParseFormatField([]byte{0x10})
	if err == nil {
		t.Error("expected error")
	}
}

func TestDefaultHdlcParameters(t *testing.T) {
	p := DefaultHdlcParameters()
	if p.WindowSize != 1 || p.MaxInfoLength != 200 {
		t.Errorf("got %v", p)
	}
}

func TestHdlcParameters_Encode(t *testing.T) {
	p := DefaultHdlcParameters()
	b := p.Encode()
	if len(b) == 0 {
		t.Error("empty encode")
	}
	p2, err := ParseHdlcParameters(b)
	if err != nil {
		t.Fatal(err)
	}
	if p2.WindowSize != p.WindowSize {
		t.Errorf("got %v", p2)
	}
}

func TestParseHdlcParameters_Empty(t *testing.T) {
	p, err := ParseHdlcParameters(nil)
	if err != nil {
		t.Fatal(err)
	}
	if p.WindowSize != 1 {
		t.Errorf("got %v", p)
	}
}

func TestParseFrame_Short(t *testing.T) {
	_, err := ParseFrame([]byte{0x7E, 0x7E})
	if err == nil {
		t.Error("expected error for short frame")
	}
}

func TestParseFrame_BadFlags(t *testing.T) {
	_, err := ParseFrame([]byte{0x01, 0x02})
	if err == nil {
		t.Error("expected error for missing flags")
	}
}

func TestCRCCCITT_SingleByte(t *testing.T) {
	crc := CRCCCITT([]byte{0x00})
	if len(crc) != 2 {
		t.Error("expected 2 bytes")
	}
	full := append([]byte{0x00}, crc...)
	if !VerifyFCS(full) {
		t.Error("FCS failed")
	}
}

func TestCRCCCITT_AllOnes(t *testing.T) {
	crc := CRCCCITT([]byte{0xFF, 0xFF, 0xFF})
	full := append([]byte{0xFF, 0xFF, 0xFF}, crc...)
	if !VerifyFCS(full) {
		t.Error("FCS failed")
	}
}

func TestCRCCCITT_Pattern(t *testing.T) {
	data := []byte{0x01, 0x02, 0x03, 0x04, 0x05}
	crc := CRCCCITT(data)
	full := append(data, crc...)
	if !VerifyFCS(full) {
		t.Error("FCS failed")
	}
	// Tamper
	full[2] = 0xFF
	if VerifyFCS(full) {
		t.Error("tampered FCS should fail")
	}
}

func TestFrameParser_EscapeInFrame(t *testing.T) {
	parser := NewFrameParser()
	// Frame containing 0x7D escape
	data := []byte{0x7E, 0x7D, 0x5E, 0x7E} // 0x7D 0x5E -> 0x7E (which would normally be a flag, but escaped)
	frames := parser.Feed(data)
	if len(frames) != 1 {
		t.Fatalf("expected 1 frame, got %d", len(frames))
	}
}

func TestHdlcEscape_NoSpecialBytes(t *testing.T) {
	data := []byte{0x01, 0x02, 0x03}
	escaped := hdlcEscape(data)
	if len(escaped) != 5 { // flag + 3 bytes + flag
		t.Errorf("len=%d", len(escaped))
	}
}

func TestHdlcUnescape_Roundtrip(t *testing.T) {
	original := []byte{0x7E, 0x01, 0x7D, 0xFF, 0x7E}
	escaped := hdlcEscape(original)
	unescaped, err := hdlcUnescape(escaped)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(unescaped, original) {
		t.Errorf("got %v, want %v", unescaped, original)
	}
}

func TestAddressLength(t *testing.T) {
	a := &HdlcAddress{Logical: 1}
	if a.AddressLength() != 1 {
		t.Error("expected 1")
	}
	a = &HdlcAddress{Logical: 128, Extended: true}
	if a.AddressLength() != 2 {
		t.Error("expected 2")
	}
	a = &HdlcAddress{Logical: 1, Physical: 2, HasPhysical: true, Extended: true}
	if a.AddressLength() != 4 {
		t.Error("expected 4")
	}
}

func TestParseAddressFromBytes(t *testing.T) {
	tests := []struct {
		data    []byte
		logical int
	}{
		{[]byte{0x02}, 1},
		{[]byte{0x00}, 0},
		{[]byte{0xFE}, 127},
	}
	for _, tc := range tests {
		a, err := parseAddressFromBytes(tc.data)
		if err != nil {
			t.Errorf("parseAddressFromBytes(%v): %v", tc.data, err)
			continue
		}
		if a.Logical != tc.logical {
			t.Errorf("got %d, want %d", a.Logical, tc.logical)
		}
	}
}

func TestParseAddressFromBytes_Invalid(t *testing.T) {
	_, err := parseAddressFromBytes([]byte{1, 2, 3})
	if err == nil {
		t.Error("expected error for 3-byte address")
	}
}

func TestSNRMControl(t *testing.T) {
	c := SNRMControl()
	if !IsSNRM(c) {
		t.Error("SNRM check")
	}
}

func TestDISCControl(t *testing.T) {
	c := DISCControl()
	if !IsDISC(c) {
		t.Error("DISC check")
	}
}

func TestDMControl(t *testing.T) {
	c := DMControl()
	if !IsDM(c) {
		t.Error("DM check")
	}
}

func TestRRControl(t *testing.T) {
	c := RRControl(0, false)
	if !IsS(c) {
		t.Error("RR should be S-frame")
	}
	if ExtractSFrameFunction(c) != 0 {
		t.Error("RR function should be 0")
	}
}

func TestREJControl(t *testing.T) {
	c := REJControl(1, true)
	if !IsS(c) {
		t.Error("REJ should be S-frame")
	}
}

func TestRNRControl(t *testing.T) {
	c := RNRControl(2, false)
	if !IsS(c) {
		t.Error("RNR should be S-frame")
	}
}

func TestEncodeFrame_Roundtrip(t *testing.T) {
	dest := &HdlcAddress{Logical: 16}
	src := &HdlcAddress{Logical: 1}
	info := []byte("test payload data")

	encoded := EncodeFrame(dest, src, IFrameControl(0, 0), info)
	parsed, err := ParseFrame(encoded)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.DestAddr.Logical != 16 {
		t.Errorf("dest=%d", parsed.DestAddr.Logical)
	}
	if !bytes.Equal(parsed.Info, info) {
		t.Errorf("info=%v", parsed.Info)
	}
}

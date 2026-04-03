package hdlc

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

const (
	HDLCFlag = 0x7E
)

// Escape and unescape HDLC frames.
func hdlcEscape(data []byte) []byte {
	var buf bytes.Buffer
	buf.WriteByte(HDLCFlag)
	for _, b := range data {
		if b == HDLCFlag || b == 0x7D {
			buf.WriteByte(0x7D)
			buf.WriteByte(b ^ 0x20)
		} else {
			buf.WriteByte(b)
		}
	}
	buf.WriteByte(HDLCFlag)
	return buf.Bytes()
}

func hdlcUnescape(data []byte) ([]byte, error) {
	if len(data) < 2 || data[0] != HDLCFlag || data[len(data)-1] != HDLCFlag {
		return nil, fmt.Errorf("hdlc: frame not enclosed by flags")
	}
	var buf bytes.Buffer
	for i := 1; i < len(data)-1; i++ {
		if data[i] == 0x7D {
			if i+1 >= len(data)-1 {
				return nil, fmt.Errorf("hdlc: incomplete escape sequence")
			}
			buf.WriteByte(data[i+1] ^ 0x20)
			i++
		} else {
			buf.WriteByte(data[i])
		}
	}
	return buf.Bytes(), nil
}

// FrameType identifies the type of HDLC frame.
type FrameType int

const (
	FrameTypeI FrameType = iota // Information
	FrameTypeS                  // Supervisory
	FrameTypeU                  // Unnumbered
)

// ParseFrameType determines frame type from control byte.
func ParseFrameType(control byte) FrameType {
	if control&0x01 == 0 {
		return FrameTypeI
	}
	if control&0x02 == 0 {
		return FrameTypeS
	}
	return FrameTypeU
}

// HdlcAddress represents an HDLC address.
type HdlcAddress struct {
	Logical     int
	Physical    int
	HasPhysical bool
	Extended    bool
}

// AddressLength returns the number of bytes in the encoded address.
func (a *HdlcAddress) AddressLength() int {
	if a.HasPhysical {
		return 4
	}
	if a.Logical > 0x7F {
		return 2
	}
	if a.Logical == 0 && a.Extended {
		return 1
	}
	return 1
}

// EncodeAddress encodes an HDLC address (client or server type).
func (a *HdlcAddress) EncodeAddress(serverSide bool) []byte {
	if !serverSide {
		// Client: single byte, LSB=1
		return []byte{byte(a.Logical<<1) | 0x01}
	}
	// Server side
	if a.HasPhysical {
		lh, ll := splitAddress(a.Logical)
		ph, pl := splitAddress(a.Physical)
		pl |= 0x01
		return []byte{lh, ll, ph, pl}
	}
	if a.Logical > 0x7F {
		h, l := splitAddress(a.Logical)
		l |= 0x01
		return []byte{h, l}
	}
	return []byte{byte(a.Logical<<1) | 0x01}
}

func splitAddress(addr int) (byte, byte) {
	if addr > 0x7F {
		lower := byte((addr & 0x7F) << 1)
		upper := byte((addr & 0x3F80) >> 6)
		return upper, lower
	}
	return 0, byte(addr << 1)
}

func parseTwoByteAddress(data []byte) int {
	upper := int(data[0] >> 1)
	lower := int(data[1] >> 1)
	return lower + (upper << 7)
}

// ParseAddresses parses destination and source addresses from frame data.
func ParseAddresses(frameData []byte) (dest, src *HdlcAddress, consumed int, err error) {
	if len(frameData) < 4 {
		return nil, nil, 0, fmt.Errorf("hdlc: frame too short for addresses")
	}
	// Destination address
	destLen := findAddressEnd(frameData, 2)
	destAddr, err := parseAddressFromBytes(frameData[2 : 2+destLen])
	if err != nil {
		return nil, nil, 0, err
	}
	// Source address
	srcStart := 2 + destLen
	if srcStart >= len(frameData) {
		return nil, nil, 0, fmt.Errorf("hdlc: no room for source address")
	}
	srcLen := findAddressEnd(frameData, srcStart)
	srcAddr, err := parseAddressFromBytes(frameData[srcStart : srcStart+srcLen])
	if err != nil {
		return nil, nil, 0, err
	}
	return destAddr, srcAddr, srcStart + srcLen, nil
}

func findAddressEnd(data []byte, start int) int {
	for i := start; i < start+4 && i < len(data); i++ {
		if data[i]&0x01 != 0 {
			return i - start + 1
		}
	}
	return 1
}

func parseAddressFromBytes(data []byte) (*HdlcAddress, error) {
	addr := &HdlcAddress{}
	switch len(data) {
	case 1:
		addr.Logical = int(data[0] >> 1)
	case 2:
		addr.Logical = parseTwoByteAddress(data)
		addr.Extended = true
	case 4:
		addr.Logical = parseTwoByteAddress(data[0:2])
		addr.Physical = parseTwoByteAddress(data[2:4])
		addr.HasPhysical = true
		addr.Extended = true
	default:
		return nil, fmt.Errorf("hdlc: invalid address length %d", len(data))
	}
	return addr, nil
}

// Frame represents a parsed HDLC frame.
type Frame struct {
	DestAddr *HdlcAddress
	SrcAddr  *HdlcAddress
	Control  byte
	HCS      []byte
	Info     []byte
	FCS      []byte
}

// FrameType returns the type of this frame.
func (f *Frame) FrameType() FrameType {
	return ParseFrameType(f.Control)
}

// ParseFrame parses a raw HDLC frame (with flags).
func ParseFrame(data []byte) (*Frame, error) {
	unescaped, err := hdlcUnescape(data)
	if err != nil {
		return nil, err
	}

	if len(unescaped) < 7 {
		return nil, fmt.Errorf("hdlc: frame too short (%d bytes)", len(unescaped))
	}

	// Format field: 2 bytes
	frameLen := int(unescaped[0])
	segmented := unescaped[1]&0x08 != 0
	_ = segmented

	// Addresses
	dest, src, addrEnd, err := ParseAddresses(unescaped)
	if err != nil {
		return nil, err
	}

	if addrEnd >= len(unescaped) {
		return nil, fmt.Errorf("hdlc: no control byte")
	}
	control := unescaped[addrEnd]

	f := &Frame{
		DestAddr: dest,
		SrcAddr:  src,
		Control:  control,
	}

	// Check if there's information
	hasInfo := (frameLen - (addrEnd - 2) - 1) > 0
	if hasInfo {
		hcsStart := addrEnd + 1
		if hcsStart+2 > len(unescaped) {
			return nil, fmt.Errorf("hdlc: no HCS")
		}
		f.HCS = unescaped[hcsStart : hcsStart+2]
		infoEnd := len(unescaped) - 2
		if infoEnd > hcsStart+2 {
			f.Info = unescaped[hcsStart+2 : infoEnd]
		}
	}

	f.FCS = unescaped[len(unescaped)-2:]

	return f, nil
}

// EncodeFrame encodes an HDLC frame.
func EncodeFrame(dest, src *HdlcAddress, control byte, info []byte) []byte {
	var buf bytes.Buffer

	// Header content (without HCS and FCS)
	destBytes := dest.EncodeAddress(src != nil)
	srcBytes := src.EncodeAddress(dest != nil)
	headerLen := 2 + len(destBytes) + len(srcBytes) + 1 // format + dest + src + control
	frameLen := headerLen
	if len(info) > 0 {
		frameLen += 2 + len(info) // HCS + info
	}

	buf.WriteByte(byte(frameLen))
	buf.WriteByte(0x00) // no segmentation
	buf.Write(destBytes)
	buf.Write(srcBytes)
	buf.WriteByte(control)

	if len(info) > 0 {
		hcs := CRCCCITT(buf.Bytes())
		buf.Write(hcs)
		buf.Write(info)
	}

	fcs := CRCCCITT(buf.Bytes())
	buf.Write(fcs)

	return hdlcEscape(buf.Bytes())
}

// Control field helpers for S-frames
func SFrameControl(rBit bool, fBit bool, function byte) byte {
	var c byte
	if rBit {
		c |= 0x02
	}
	if fBit {
		c |= 0x10
	}
	c |= (function & 0x0F)
	return c
}

// Control field helpers for U-frames
func UFrameControl(fBit bool, modifierA, modifierB byte) byte {
	var c byte
	if fBit {
		c |= 0x10
	}
	c |= 0x03 // U-frame marker
	c |= (modifierB << 2)
	c |= (modifierA << 3)
	return c
}

// I-frame control
func IFrameControl(sendSeq, recvSeq int) byte {
	return byte((sendSeq << 1) | (recvSeq & 0x07))
}

// FrameParser reads HDLC frames from a stream.
type FrameParser struct {
	buf     []byte
	inFrame bool
	escaped bool
}

// NewFrameParser creates a new stream parser.
func NewFrameParser() *FrameParser {
	return &FrameParser{}
}

// Feed adds data to the parser and returns any complete frames found.
func (p *FrameParser) Feed(data []byte) [][]byte {
	var frames [][]byte
	for _, b := range data {
		if !p.inFrame {
			if b == HDLCFlag {
				p.inFrame = true
				p.buf = p.buf[:0]
				p.buf = append(p.buf, b)
			}
			continue
		}

		p.buf = append(p.buf, b)

		if b == HDLCFlag {
			if len(p.buf) > 2 {
				frame := make([]byte, len(p.buf))
				copy(frame, p.buf)
				frames = append(frames, frame)
			}
			p.inFrame = false
			p.escaped = false
		}
	}
	return frames
}

// Reset clears the parser state.
func (p *FrameParser) Reset() {
	p.buf = p.buf[:0]
	p.inFrame = false
}

// StreamFrames reads frames from an io.Reader until an error or close.
func StreamFrames(r io.Reader) ([][]byte, error) {
	parser := NewFrameParser()
	buf := make([]byte, 4096)
	var allFrames [][]byte
	for {
		n, err := r.Read(buf)
		if n > 0 {
			frames := parser.Feed(buf[:n])
			allFrames = append(allFrames, frames...)
		}
		if err != nil {
			if err == io.EOF {
				return allFrames, nil
			}
			return allFrames, err
		}
	}
}

// ParseFormatField parses the 2-byte format field.
func ParseFormatField(data []byte) (frameLength int, segmented bool, err error) {
	if len(data) < 2 {
		return 0, false, fmt.Errorf("hdlc: format field too short")
	}
	return int(data[0]), data[1]&0x08 != 0, nil
}

// EncodeFormatField encodes the 2-byte format field.
func EncodeFormatField(frameLength int, segmented bool) []byte {
	b := make([]byte, 2)
	b[0] = byte(frameLength)
	if segmented {
		b[1] = 0x08
	}
	return b
}

// SNRM control byte
func SNRMControl() byte { return 0x93 }

// UA control byte
func UAControl(fBit bool) byte {
	if fBit {
		return 0x73
	}
	return 0x63
}

// DISC control byte
func DISCControl() byte { return 0x53 }

// DM control byte
func DMControl() byte { return 0x1F }

// RR (Receiver Ready) S-frame
func RRControl(recvSeq int, fBit bool) byte {
	c := byte((recvSeq&0x07)<<5) | 0x01
	if fBit {
		c |= 0x10
	}
	return c
}

// REJ (Reject) S-frame
func REJControl(recvSeq int, fBit bool) byte {
	c := byte((recvSeq&0x07)<<5) | 0x05
	if fBit {
		c |= 0x10
	}
	return c
}

// RNR (Receiver Not Ready) S-frame
func RNRControl(recvSeq int, fBit bool) byte {
	c := byte((recvSeq&0x07)<<5) | 0x09
	if fBit {
		c |= 0x10
	}
	return c
}

// IsUA checks if control byte is UA.
func IsUA(control byte) bool {
	return (control & 0xEF) == 0x63
}

// IsSNRM checks if control byte is SNRM.
func IsSNRM(control byte) bool {
	return control == 0x93
}

// IsDISC checks if control byte is DISC.
func IsDISC(control byte) bool {
	return control == 0x53
}

// IsDM checks if control byte is DM.
func IsDM(control byte) bool {
	return control == 0x0F || control == 0x1F
}

// IsI checks if frame is I-frame.
func IsI(control byte) bool {
	return control&0x01 == 0
}

// IsS checks if frame is S-frame.
func IsS(control byte) bool {
	return (control&0x03 == 0x01) && (control&0x02 == 0)
}

// IsU checks if frame is U-frame.
func IsU(control byte) bool {
	return control&0x03 == 0x03
}

// ExtractSendSeq extracts send sequence from I-frame control byte.
func ExtractSendSeq(control byte) int {
	return int((control >> 1) & 0x7F)
}

// ExtractRecvSeq extracts receive sequence from I-frame control byte.
func ExtractRecvSeq(control byte) int {
	return int(control & 0x07)
}

// ExtractSFrameRecvSeq extracts receive sequence from S-frame control byte.
func ExtractSFrameRecvSeq(control byte) int {
	return int(control>>5) & 0x07
}

// ExtractSFrameFunction extracts S-frame function.
func ExtractSFrameFunction(control byte) byte {
	return (control >> 2) & 0x03
}

// ExtractUFrameModifier extracts U-frame modifier bits.
func ExtractUFrameModifier(control byte) byte {
	return ((control >> 2) & 0x03) | ((control >> 3) & 0x04)
}

// HdlcParameters represents HDLC parameter negotiation.
type HdlcParameters struct {
	WindowSize    int
	MaxInfoLength int
	MaxTransmit   int
	RetryCount    int
}

// DefaultHdlcParameters returns default HDLC parameters.
func DefaultHdlcParameters() *HdlcParameters {
	return &HdlcParameters{
		WindowSize:    1,
		MaxInfoLength: 200,
		MaxTransmit:   20,
		RetryCount:    3,
	}
}

// Encode encodes parameters as TLV.
func (p *HdlcParameters) Encode() []byte {
	var buf bytes.Buffer
	// Window size
	buf.WriteByte(0x01) // tag
	buf.WriteByte(0x00) // length 0 (window size is in next byte)
	buf.WriteByte(0x01) // value 1
	// Max info length
	buf.WriteByte(0x02)
	buf.WriteByte(0x00)
	binary.Write(&buf, binary.BigEndian, uint16(p.MaxInfoLength))
	// Max transmit
	buf.WriteByte(0x03)
	buf.WriteByte(0x00)
	buf.WriteByte(byte(p.MaxTransmit))
	// Retry count
	buf.WriteByte(0x04)
	buf.WriteByte(0x00)
	buf.WriteByte(byte(p.RetryCount))
	return buf.Bytes()
}

// ParseHdlcParameters parses HDLC parameters from info field.
func ParseHdlcParameters(info []byte) (*HdlcParameters, error) {
	p := DefaultHdlcParameters()
	i := 0
	for i < len(info) {
		if i+1 >= len(info) {
			break
		}
		tag := info[i]
		length := int(info[i+1])
		i += 2

		if length == 0 {
			// Single byte value follows
			if i >= len(info) {
				break
			}
			val := info[i]
			i++
			switch tag {
			case 0x01:
				p.WindowSize = int(val)
			case 0x03:
				p.MaxTransmit = int(val)
			case 0x04:
				p.RetryCount = int(val)
			}
		} else {
			if i+length > len(info) {
				break
			}
			data := info[i : i+length]
			i += length
			switch tag {
			case 0x02:
				if len(data) >= 2 {
					p.MaxInfoLength = int(binary.BigEndian.Uint16(data))
				}
			}
		}
	}
	return p, nil
}

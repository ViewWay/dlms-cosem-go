package hdlc

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// HDLC parameter type constants as defined in Green Book Edition 9,// Section 8.4.5.3.2
const (
	// Format and group identifiers for SNRM/UA information field
	FormatIdentifier = 0x81
	GroupIdentifier  = 0x80

	// Parameter IDs for HDLC parameter negotiation
	ParamMaxInfoFieldLengthTX = 0x05 // Maximum information field length - transmit
	ParamMaxInfoFieldLengthRX = 0x06 // Maximum information field length - receive
	ParamWindowSizeTX         = 0x07 // Window size - transmit
	ParamWindowSizeRX         = 0x08 // Window size - receive
)

// HdlcParametersGreen represents HDLC parameters according to Green Book specification.
// This structure uses separate TX/RX parameters for window size and max info length.
type HdlcParametersGreen struct {
	MaxInfoLengthTX int // Maximum information field length for transmission (128-2048)
	MaxInfoLengthRX int // Maximum information field length for reception (128-2048)
	WindowSizeTX    int // Window size for transmission (1-7)
	WindowSizeRX    int // Window size for reception (1-7)
}

// DefaultHdlcParametersGreen returns default HDLC parameters.
func DefaultHdlcParametersGreen() *HdlcParametersGreen {
	return &HdlcParametersGreen{
		MaxInfoLengthTX: 128,
		MaxInfoLengthRX: 128,
		WindowSizeTX:    1,
		WindowSizeRX:    1,
	}
}

// Validate validates HDLC parameters according to Green Book ranges.
func (p *HdlcParametersGreen) Validate() error {
	if p.WindowSizeTX < 1 || p.WindowSizeTX > 7 {
		return fmt.Errorf("WindowSizeTX %d out of range [1, 7]", p.WindowSizeTX)
	}
	if p.WindowSizeRX < 1 || p.WindowSizeRX > 7 {
		return fmt.Errorf("WindowSizeRX %d out of range [1, 7]", p.WindowSizeRX)
	}
	if p.MaxInfoLengthTX < 128 || p.MaxInfoLengthTX > 2048 {
		return fmt.Errorf("MaxInfoLengthTX %d out of range [128, 2048]", p.MaxInfoLengthTX)
	}
	if p.MaxInfoLengthRX < 128 || p.MaxInfoLengthRX > 2048 {
		return fmt.Errorf("MaxInfoLengthRX %d out of range [128, 2048]", p.MaxInfoLengthRX)
	}
	return nil
}

// Encode encodes parameters as TLV with Green Book format.
// Returns bytes in format: 0x81, 0x80, group_length, params...
func (p *HdlcParametersGreen) Encode() ([]byte, error) {
	if err := p.Validate(); err != nil {
		return nil, err
	}

	var paramsBuf bytes.Buffer

	// Encode MaxInfoLengthTX (0x05)
	paramsBuf.WriteByte(ParamMaxInfoFieldLengthTX)
	if p.MaxInfoLengthTX <= 255 {
		paramsBuf.WriteByte(0x01) // length = 1
		paramsBuf.WriteByte(byte(p.MaxInfoLengthTX))
	} else {
		paramsBuf.WriteByte(0x02) // length = 2
		binary.Write(&paramsBuf, binary.BigEndian, uint16(p.MaxInfoLengthTX))
	}

	// Encode MaxInfoLengthRX (0x06)
	paramsBuf.WriteByte(ParamMaxInfoFieldLengthRX)
	if p.MaxInfoLengthRX <= 255 {
		paramsBuf.WriteByte(0x01) // length = 1
		paramsBuf.WriteByte(byte(p.MaxInfoLengthRX))
	} else {
		paramsBuf.WriteByte(0x02) // length = 2
		binary.Write(&paramsBuf, binary.BigEndian, uint16(p.MaxInfoLengthRX))
	}

	// Encode WindowSizeTX (0x07)
	paramsBuf.WriteByte(ParamWindowSizeTX)
	paramsBuf.WriteByte(0x01) // length = 1
	paramsBuf.WriteByte(byte(p.WindowSizeTX))

	// Encode WindowSizeRX (0x08)
	paramsBuf.WriteByte(ParamWindowSizeRX)
	paramsBuf.WriteByte(0x01) // length = 1
	paramsBuf.WriteByte(byte(p.WindowSizeRX))

	// Build complete frame with header
	var buf bytes.Buffer
	buf.WriteByte(FormatIdentifier) // 0x81
	buf.WriteByte(GroupIdentifier)  // 0x80
	buf.WriteByte(byte(paramsBuf.Len())) // group length
	buf.Write(paramsBuf.Bytes())

	return buf.Bytes(), nil
}

// ParseHdlcParametersGreen parses HDLC parameters from Green Book format.
// Expected format: 0x81, 0x80, group_length, params...
func ParseHdlcParametersGreen(data []byte) (*HdlcParametersGreen, error) {
	p := DefaultHdlcParametersGreen()

	if len(data) < 3 {
		return nil, fmt.Errorf("data too short for Green Book format")
	}

	// Check format and group identifiers
	if data[0] != FormatIdentifier || data[1] != GroupIdentifier {
		return nil, fmt.Errorf("invalid format/group identifiers: got %02x %02x, want %02x %02x",
			data[0], data[1], FormatIdentifier, GroupIdentifier)
	}

	groupLength := int(data[2])
	if len(data) < 3+groupLength {
		return nil, fmt.Errorf("data too short for group length %d", groupLength)
	}

	// Parse parameters
	i := 3
	for i < len(data) {
		if i+1 >= len(data) {
			break
		}
		tag := data[i]
		length := int(data[i+1])
		i += 2

		if i+length > len(data) {
			return nil, fmt.Errorf("parameter value extends beyond data")
		}

		valueData := data[i : i+length]
		i += length

		switch tag {
		case ParamMaxInfoFieldLengthTX:
			if length == 1 {
				p.MaxInfoLengthTX = int(valueData[0])
			} else if length == 2 {
				p.MaxInfoLengthTX = int(binary.BigEndian.Uint16(valueData))
			}
		case ParamMaxInfoFieldLengthRX:
			if length == 1 {
				p.MaxInfoLengthRX = int(valueData[0])
			} else if length == 2 {
				p.MaxInfoLengthRX = int(binary.BigEndian.Uint16(valueData))
			}
		case ParamWindowSizeTX:
			if length >= 1 {
				p.WindowSizeTX = int(valueData[0])
			}
		case ParamWindowSizeRX:
			if length >= 1 {
				p.WindowSizeRX = int(valueData[0])
			}
		}
	}

	return p, p.Validate()
}

// NegotiateParametersGreen negotiates HDLC parameters between client and server.
// Note: Client TX parameters correspond to server RX and vice versa.
func NegotiateParametersGreen(client, server *HdlcParametersGreen) *HdlcParametersGreen {
	negotiated := &HdlcParametersGreen{
		// Client TX negotiates with server RX
		MaxInfoLengthTX: min(client.MaxInfoLengthTX, server.MaxInfoLengthRX),
		WindowSizeTX:    min(client.WindowSizeTX, server.WindowSizeRX),
		// Client RX negotiates with server TX
		MaxInfoLengthRX: min(client.MaxInfoLengthRX, server.MaxInfoLengthTX),
		WindowSizeRX:    min(client.WindowSizeRX, server.WindowSizeTX),
	}
	return negotiated
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

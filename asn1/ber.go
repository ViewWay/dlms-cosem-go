package asn1

import (
	"encoding/binary"
	"fmt"
)

// BER encoding/decoding for DLMS/COSEM ACSE (AARQ/AARE/RLRE).

// BEREncode encodes a TLV in BER format.
func BEREncode(tag int, data []byte) []byte {
	if data == nil || len(data) == 0 {
		return BEREncodeLength0(tag)
	}
	out := make([]byte, 0, 1+len(data)+4)
	// Tag
	if tag >= 0x1F {
		// Context class tag: use as-is for DLMS tags
		out = append(out, byte(tag))
	} else {
		out = append(out, byte(tag))
	}
	// Length
	out = append(out, encodeBERLength(len(data))...)
	// Value
	out = append(out, data...)
	return out
}

func BEREncodeLength0(tag int) []byte {
	return []byte{byte(tag), 0x00}
}

func encodeBERLength(length int) []byte {
	if length < 0x80 {
		return []byte{byte(length)}
	}
	if length < 0x100 {
		return []byte{0x81, byte(length)}
	}
	if length < 0x10000 {
		return []byte{0x82, byte(length >> 8), byte(length)}
	}
	return []byte{0x83, byte(length >> 16), byte(length >> 8), byte(length)}
}

// BERDecode decodes a TLV from BER data. Returns tag, value, and bytes consumed.
func BERDecode(data []byte) (int, []byte, int, error) {
	if len(data) < 2 {
		return 0, nil, 0, fmt.Errorf("asn1: data too short")
	}
	// Tag
	tag, tagLen := decodeBERTag(data)
	// Length
	length, lengthLen, err := decodeBERLength(data[tagLen:])
	if err != nil {
		return 0, nil, 0, err
	}
	totalLen := tagLen + lengthLen + length
	if totalLen > len(data) {
		return 0, nil, 0, fmt.Errorf("asn1: data truncated (need %d, have %d)", totalLen, len(data))
	}
	value := data[tagLen+lengthLen : tagLen+lengthLen+length]
	return tag, value, totalLen, nil
}

func decodeBERTag(data []byte) (int, int) {
	// Check if lower 5 bits are NOT 0x1F (single-byte tag number)
	if data[0]&0x1F != 0x1F {
		return int(data[0]), 1
	}
	// Multi-byte tag: lower 5 bits are 0x1F, tag number follows in base-128
	if len(data) < 2 {
		return int(data[0]), 1
	}
	var tag int
	i := 1
	for i < len(data) {
		b := data[i]
		tag = (tag << 7) | int(b&0x7F)
		i++
		if b&0x80 == 0 {
			break // Last byte (no continuation bit)
		}
		if i >= len(data) {
			// Incomplete multi-byte tag, return what we have
			break
		}
	}
	return tag, i
}

func decodeBERLength(data []byte) (int, int, error) {
	if len(data) == 0 {
		return 0, 0, fmt.Errorf("asn1: no length byte")
	}
	first := data[0]
	if first < 0x80 {
		return int(first), 1, nil
	}
	numBytes := int(first & 0x7F)
	if numBytes == 0 {
		return 0, 1, nil // indefinite length (not supported in practice)
	}
	if numBytes > 4 || len(data) < 1+numBytes {
		return 0, 0, fmt.Errorf("asn1: invalid length encoding")
	}
	var length int
	for i := 0; i < numBytes; i++ {
		length = length<<8 | int(data[1+i])
	}
	return length, 1 + numBytes, nil
}

// BEREncodeOctetString encodes an OCTET STRING.
func BEREncodeOctetString(data []byte) []byte {
	return BEREncode(0x04, data)
}

// BEREncodeInteger encodes an INTEGER.
func BEREncodeInteger(value int) []byte {
	var b []byte
	if value == 0 {
		b = []byte{0}
	} else {
		n := value
		if n < 0 {
			n = -n
		}
		for n > 0 {
			b = append([]byte{byte(n)}, b...)
			n >>= 8
		}
		if value < 0 {
			// Add sign bit
			if b[0]&0x80 != 0 {
				b = append([]byte{0xFF}, b...)
			} else {
				b[0] |= 0x80
			}
		}
	}
	return BEREncode(0x02, b)
}

// BEREncodeContextTag encodes with context-specific tag.
func BEREncodeContextTag(tag int, data []byte) []byte {
	constructedTag := 0xA0 | (tag & 0x1F)
	return BEREncode(constructedTag, data)
}

// BEREncodeContextPrimitive encodes a primitive context-specific tag.
func BEREncodeContextPrimitive(tag int, data []byte) []byte {
	primitiveTag := 0x80 | (tag & 0x1F)
	return BEREncode(primitiveTag, data)
}

// BEREncodeSequence encodes a SEQUENCE.
func BEREncodeSequence(data []byte) []byte {
	return BEREncode(0x30, data)
}

// BERDecodeAllTLV decodes all TLVs from data.
func BERDecodeAllTLV(data []byte) ([]struct {
	Tag   int
	Value []byte
}, error) {
	var result []struct {
		Tag   int
		Value []byte
	}
	pos := 0
	for pos < len(data) {
		tag, value, consumed, err := BERDecode(data[pos:])
		if err != nil {
			return nil, err
		}
		result = append(result, struct {
			Tag   int
			Value []byte
		}{tag, value})
		pos += consumed
	}
	return result, nil
}

// --- DLMS ACSE types ---

// AppContextName constants
const (
	AppContextNameLogicalNameRef = 0x0007
	AppContextNameShortNameRef   = 0x0006
	AppContextNameCipheredLN     = 0x0008
)

// AARQ (Application Association Request) tag = 0x60
const TagAARQ = 0x60
const TagAARE = 0x61
const TagRLRQ = 0x62
const TagRLRE = 0x64

// Context tags for AARQ/AARE fields
const (
	CTagApplicationContextName       = 0xA1
	CTagCalledAPTitle                = 0xA2
	CTagCalledAEQualifier            = 0xA3
	CTagCalledAPInvocationIdentifier = 0xA4
	CTagCalledAEInvocationIdentifier = 0xA5
	CTagCallingAPTitle               = 0xA6
	CTagCallingAEQualifier           = 0xA7
	CTagCallingAPInvocationID        = 0xA8
	CTagCallingAEInvocationID        = 0xA9
	CTagSenderACSERequirements       = 0x8A
	CTagMechanismName                = 0x8B
	CTagCallingAuthenticationValue   = 0xAC
	CTagImplementationInfo           = 0xBD
	CTagResult                       = 0x8A
	CTagResultSourceDiagnostic       = 0x8B
	CTagUserInformation              = 0xBE
)

// Authentication mechanism values
const (
	AuthNone      = 0
	AuthLLS       = 1
	AuthHLS       = 2
	AuthHLSMD5    = 3
	AuthHLSSHA1   = 4
	AuthHLSGMAC   = 5
	AuthHLSSHA256 = 6
	AuthHLSECDSA  = 7
)

// AssociationResult
const (
	AssocAccepted     = 0
	AssocRejectedPerm = 1
	AssocRejectedTemp = 2
)

// AARQ represents an Application Association Request.
type AARQ struct {
	ApplicationContextName []byte
	CallingAPTitle         []byte // System title
	CallingAEQualifier     []byte // Public cert
	Authentication         int
	AuthenticationValue    []byte
	CallingAEInvocationID  []byte // User ID
	UserInformation        []byte
}

// Encode encodes AARQ to BER bytes.
func (a *AARQ) Encode() []byte {
	var content []byte

	if a.ApplicationContextName != nil {
		content = append(content, BEREncodeContextTag(1, a.ApplicationContextName)...)
	}
	if a.CallingAPTitle != nil {
		content = append(content, BEREncodeContextTag(6, BEREncodeOctetString(a.CallingAPTitle))...)
	}
	if a.CallingAEQualifier != nil {
		content = append(content, BEREncodeContextTag(7, BEREncodeOctetString(a.CallingAEQualifier))...)
	}
	if a.CallingAEInvocationID != nil {
		content = append(content, BEREncodeContextTag(9, a.CallingAEInvocationID)...)
	}
	if a.Authentication != AuthNone {
		// sender ACSE requirements
		content = append(content, BEREncodeContextPrimitive(0x0A, []byte{0x01})...)
		// mechanism name
		mech := make([]byte, 4)
		binary.BigEndian.PutUint32(mech, uint32(a.Authentication))
		content = append(content, BEREncodeContextPrimitive(0x0B, mech)...)
	}
	if a.AuthenticationValue != nil {
		content = append(content, BEREncodeContextPrimitive(0x2C, a.AuthenticationValue)...)
	}
	if a.UserInformation != nil {
		content = append(content, BEREncodeContextTag(0x1E, a.UserInformation)...)
	}

	return BEREncode(TagAARQ, content)
}

// AARE represents an Application Association Response.
type AARE struct {
	Result                 int
	ResultSourceDiagnostic []byte
	ApplicationContextName []byte
	RespondingAPTitle      []byte
	RespondingAEQualifier  []byte
	Authentication         int
	AuthenticationValue    []byte
	UserInformation        []byte
}

// Encode encodes AARE to BER bytes.
func (a *AARE) Encode() []byte {
	var content []byte

	content = append(content, BEREncodeContextPrimitive(0x0A, []byte{byte(a.Result)})...)

	if a.ResultSourceDiagnostic != nil {
		content = append(content, BEREncodeContextPrimitive(0x0B, a.ResultSourceDiagnostic)...)
	}

	if a.ApplicationContextName != nil {
		content = append(content, BEREncodeContextTag(1, a.ApplicationContextName)...)
	}
	if a.RespondingAPTitle != nil {
		content = append(content, BEREncodeContextTag(6, BEREncodeOctetString(a.RespondingAPTitle))...)
	}
	if a.RespondingAEQualifier != nil {
		content = append(content, BEREncodeContextTag(7, BEREncodeOctetString(a.RespondingAEQualifier))...)
	}
	if a.Authentication != AuthNone && a.AuthenticationValue != nil {
		content = append(content, BEREncodeContextPrimitive(0x2C, a.AuthenticationValue)...)
	}
	if a.UserInformation != nil {
		content = append(content, BEREncodeContextTag(0x1E, a.UserInformation)...)
	}

	return BEREncode(TagAARE, content)
}

// RLRE represents a Release Response.
type RLRE struct {
	Result          int
	UserInformation []byte
}

// Encode encodes RLRE to BER bytes.
func (r *RLRE) Encode() []byte {
	var content []byte
	content = append(content, BEREncodeContextPrimitive(0x0A, []byte{byte(r.Result)})...)
	if r.UserInformation != nil {
		content = append(content, BEREncodeContextTag(0x1E, r.UserInformation)...)
	}
	return BEREncode(TagRLRE, content)
}

// ParseAARE parses an AARE from BER bytes.
func ParseAARE(data []byte) (*AARE, error) {
	tag, content, _, err := BERDecode(data)
	if err != nil {
		return nil, err
	}
	if tag != TagAARE {
		return nil, fmt.Errorf("asn1: expected AARE tag 0x61, got 0x%02x", tag)
	}

	aare := &AARE{}
	tlvs, err := BERDecodeAllTLV(content)
	if err != nil {
		return nil, err
	}

	for _, tlv := range tlvs {
		switch tlv.Tag {
		case 0x8A:
			if len(tlv.Value) > 0 {
				aare.Result = int(tlv.Value[0])
			}
		case 0x8B:
			aare.ResultSourceDiagnostic = tlv.Value
		case CTagApplicationContextName:
			aare.ApplicationContextName = tlv.Value
		case CTagCallingAPTitle: // responding AP title
			if len(tlv.Value) >= 2 && tlv.Value[0] == 0x04 {
				aare.RespondingAPTitle = tlv.Value[2:]
			}
		case CTagUserInformation:
			aare.UserInformation = tlv.Value
		}
	}

	return aare, nil
}

// ParseAARQ parses an AARQ from BER bytes.
func ParseAARQ(data []byte) (*AARQ, error) {
	tag, content, _, err := BERDecode(data)
	if err != nil {
		return nil, err
	}
	if tag != TagAARQ {
		return nil, fmt.Errorf("asn1: expected AARQ tag 0x60, got 0x%02x", tag)
	}

	aarq := &AARQ{}
	tlvs, err := BERDecodeAllTLV(content)
	if err != nil {
		return nil, err
	}

	for _, tlv := range tlvs {
		switch tlv.Tag {
		case CTagApplicationContextName:
			aarq.ApplicationContextName = tlv.Value
		case CTagCallingAPTitle:
			if len(tlv.Value) >= 2 && tlv.Value[0] == 0x04 {
				aarq.CallingAPTitle = tlv.Value[2:]
			}
		case CTagCallingAEQualifier:
			if len(tlv.Value) >= 2 && tlv.Value[0] == 0x04 {
				aarq.CallingAEQualifier = tlv.Value[2:]
			}
		case CTagCallingAEInvocationID:
			aarq.CallingAEInvocationID = tlv.Value
		case CTagSenderACSERequirements:
			// Authentication required
		case CTagMechanismName:
			if len(tlv.Value) >= 4 {
				aarq.Authentication = int(binary.BigEndian.Uint32(tlv.Value[:4]))
			}
		case CTagCallingAuthenticationValue:
			aarq.AuthenticationValue = tlv.Value
		case CTagUserInformation:
			aarq.UserInformation = tlv.Value
		}
	}

	return aarq, nil
}

// ParseRLRE parses an RLRE from BER bytes.
func ParseRLRE(data []byte) (*RLRE, error) {
	tag, content, _, err := BERDecode(data)
	if err != nil {
		return nil, err
	}
	if tag != TagRLRE {
		return nil, fmt.Errorf("asn1: expected RLRE tag 0x64, got 0x%02x", tag)
	}

	rlre := &RLRE{}
	tlvs, err := BERDecodeAllTLV(content)
	if err != nil {
		return nil, err
	}

	for _, tlv := range tlvs {
		switch tlv.Tag {
		case 0x8A:
			if len(tlv.Value) > 0 {
				rlre.Result = int(tlv.Value[0])
			}
		case CTagUserInformation:
			rlre.UserInformation = tlv.Value
		}
	}

	return rlre, nil
}

// EncodeAppContextName encodes an application context name OID.
func EncodeAppContextName(logicalNameRef bool, ciphered bool) []byte {
	var oid int
	if logicalNameRef {
		if ciphered {
			oid = AppContextNameCipheredLN
		} else {
			oid = AppContextNameLogicalNameRef
		}
	} else {
		oid = AppContextNameShortNameRef
	}
	oidBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(oidBytes, uint16(oid))
	// BER encode OID: 06 02 xx xx
	return []byte{0x06, 0x02, oidBytes[0], oidBytes[1]}
}

// EncodeUserInformation wraps initiate request in user information.
func EncodeUserInformation(initiateRequest []byte) []byte {
	// user-information = context-constructed 30
	// content = OCTET STRING wrapping initiate request
	octetString := BEREncodeOctetString(initiateRequest)
	return BEREncodeContextTag(0x1E, octetString)
}

// DecodeUserInformation extracts the initiate request from user information.
func DecodeUserInformation(data []byte) ([]byte, error) {
	tag, content, _, err := BERDecode(data)
	if err != nil {
		return nil, err
	}
	if tag != CTagUserInformation {
		return nil, fmt.Errorf("asn1: expected user information tag, got 0x%02x", tag)
	}
	// Decode inner OCTET STRING
	_, value, _, err := BERDecode(content)
	if err != nil {
		return nil, err
	}
	return value, nil
}

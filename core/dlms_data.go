package core

import (
	"encoding/binary"
	"fmt"
	"math"
	"time"
)

// VariableLength indicates a variable-length data element.
const VariableLength = -1

// DlmsTag constants for DLMS/COSEM data types.
const (
	TagNull               = 0
	TagArray              = 1
	TagStructure          = 2
	TagBoolean            = 3
	TagBitString          = 4
	TagDoubleLong         = 5
	TagDoubleLongUnsigned = 6
	TagOctetString        = 9
	TagVisibleString      = 10
	TagUTF8String         = 12
	TagBCD                = 13
	TagInteger            = 15
	TagLong               = 16
	TagUnsignedInteger    = 17
	TagUnsignedLong       = 18
	TagCompactArray       = 19
	TagLong64             = 20
	TagUnsignedLong64     = 21
	TagEnum               = 22
	TagFloat32            = 23
	TagFloat64            = 24
	TagDateTime           = 25
	TagDate               = 26
	TagTime               = 27
	TagDontCare           = 255
)

// DlmsData is the interface for all DLMS/COSEM data types.
type DlmsData interface {
	Tag() byte
	ToBytes() []byte
	ToPython() interface{}
}

// DlmsDataFromBytes parses a single DlmsData element from bytes.
// Returns the parsed element and the number of bytes consumed.
func DlmsDataFromBytes(data []byte) (DlmsData, int, error) {
	if len(data) < 1 {
		return nil, 0, fmt.Errorf("dlms: insufficient data")
	}
	tag := data[0]
	consumed := 1

	switch tag {
	case TagNull:
		return NullData{}, consumed, nil
	case TagBoolean:
		if len(data) < 2 {
			return nil, 0, fmt.Errorf("dlms: insufficient data for boolean")
		}
		return BooleanData(data[1] != 0), consumed + 1, nil
	case TagBitString:
		length, n := decodeVariableLength(data[consumed:])
		consumed += n
		if len(data) < consumed+int(length) {
			return nil, 0, fmt.Errorf("dlms: insufficient data for bitstring")
		}
		val := make([]byte, length)
		copy(val, data[consumed:consumed+int(length)])
		return BitStringData(val), consumed + int(length), nil
	case TagDoubleLong:
		if len(data) < 5 {
			return nil, 0, fmt.Errorf("dlms: insufficient data for double-long")
		}
		v := int32(binary.BigEndian.Uint32(data[1:5]))
		return DoubleLongData(v), consumed + 4, nil
	case TagDoubleLongUnsigned:
		if len(data) < 5 {
			return nil, 0, fmt.Errorf("dlms: insufficient data for double-long-unsigned")
		}
		v := binary.BigEndian.Uint32(data[1:5])
		return DoubleLongUnsignedData(v), consumed + 4, nil
	case TagOctetString:
		length, n := decodeVariableLength(data[consumed:])
		consumed += n
		if len(data) < consumed+int(length) {
			return nil, 0, fmt.Errorf("dlms: insufficient data for octet-string")
		}
		val := make([]byte, length)
		copy(val, data[consumed:consumed+int(length)])
		return OctetStringData(val), consumed + int(length), nil
	case TagVisibleString:
		length, n := decodeVariableLength(data[consumed:])
		consumed += n
		if len(data) < consumed+int(length) {
			return nil, 0, fmt.Errorf("dlms: insufficient data for visible-string")
		}
		val := string(data[consumed : consumed+int(length)])
		return VisibleStringData(val), consumed + int(length), nil
	case TagUTF8String:
		length, n := decodeVariableLength(data[consumed:])
		consumed += n
		if len(data) < consumed+int(length) {
			return nil, 0, fmt.Errorf("dlms: insufficient data for utf8-string")
		}
		val := string(data[consumed : consumed+int(length)])
		return UTF8StringData(val), consumed + int(length), nil
	case TagBCD:
		length, n := decodeVariableLength(data[consumed:])
		consumed += n
		if len(data) < consumed+int(length) {
			return nil, 0, fmt.Errorf("dlms: insufficient data for bcd")
		}
		val := make([]byte, length)
		copy(val, data[consumed:consumed+int(length)])
		return BCDData(val), consumed + int(length), nil
	case TagInteger:
		if len(data) < 2 {
			return nil, 0, fmt.Errorf("dlms: insufficient data for integer")
		}
		return IntegerData(int8(data[1])), consumed + 1, nil
	case TagLong:
		if len(data) < 3 {
			return nil, 0, fmt.Errorf("dlms: insufficient data for long")
		}
		v := int16(binary.BigEndian.Uint16(data[1:3]))
		return LongData(v), consumed + 2, nil
	case TagUnsignedInteger:
		if len(data) < 2 {
			return nil, 0, fmt.Errorf("dlms: insufficient data for unsigned-integer")
		}
		return UnsignedIntegerData(data[1]), consumed + 1, nil
	case TagUnsignedLong:
		if len(data) < 3 {
			return nil, 0, fmt.Errorf("dlms: insufficient data for unsigned-long")
		}
		v := binary.BigEndian.Uint16(data[1:3])
		return UnsignedLongData(v), consumed + 2, nil
	case TagLong64:
		if len(data) < 9 {
			return nil, 0, fmt.Errorf("dlms: insufficient data for long64")
		}
		v := int64(binary.BigEndian.Uint64(data[1:9]))
		return Long64Data(v), consumed + 8, nil
	case TagUnsignedLong64:
		if len(data) < 9 {
			return nil, 0, fmt.Errorf("dlms: insufficient data for unsigned-long64")
		}
		v := binary.BigEndian.Uint64(data[1:9])
		return UnsignedLong64Data(v), consumed + 8, nil
	case TagEnum:
		if len(data) < 2 {
			return nil, 0, fmt.Errorf("dlms: insufficient data for enum")
		}
		return EnumData(data[1]), consumed + 1, nil
	case TagFloat32:
		if len(data) < 5 {
			return nil, 0, fmt.Errorf("dlms: insufficient data for float32")
		}
		bits := binary.BigEndian.Uint32(data[1:5])
		v := math.Float32frombits(bits)
		return Float32Data(v), consumed + 4, nil
	case TagFloat64:
		if len(data) < 9 {
			return nil, 0, fmt.Errorf("dlms: insufficient data for float64")
		}
		bits := binary.BigEndian.Uint64(data[1:9])
		v := math.Float64frombits(bits)
		return Float64Data(v), consumed + 8, nil
	case TagDateTime:
		if len(data) < 13 {
			return nil, 0, fmt.Errorf("dlms: insufficient data for datetime")
		}
		dt, err := CosemDateTimeFromBytes(data[1:13])
		if err != nil {
			return nil, 0, err
		}
		return DateTimeData{Value: dt}, consumed + 12, nil
	case TagDate:
		if len(data) < 6 {
			return nil, 0, fmt.Errorf("dlms: insufficient data for date")
		}
		d := CosemDateFromBytes(data[1:6])
		return DateData{Value: d}, consumed + 5, nil
	case TagTime:
		if len(data) < 5 {
			return nil, 0, fmt.Errorf("dlms: insufficient data for time")
		}
		t := CosemTimeFromBytes(data[1:5])
		return TimeData{Value: t}, consumed + 4, nil
	case TagCompactArray:
		// CompactArray: tag + length + type_description + octet_string content
		length, n := decodeVariableLength(data[consumed:])
		consumed += n
		if len(data) < consumed+int(length) {
			return nil, 0, fmt.Errorf("dlms: insufficient data for compact-array")
		}
		val := make([]byte, length)
		copy(val, data[consumed:consumed+int(length)])
		return CompactArrayData(val), consumed + int(length), nil
	case TagDontCare:
		return DontCareData{}, consumed, nil
	case TagArray:
		count, n := decodeVariableLength(data[consumed:])
		consumed += n
		elements := make([]DlmsData, count)
		for i := 0; i < int(count); i++ {
			elem, nc, err := DlmsDataFromBytes(data[consumed:])
			if err != nil {
				return nil, 0, err
			}
			elements[i] = elem
			consumed += nc
		}
		return ArrayData(elements), consumed, nil
	case TagStructure:
		count, n := decodeVariableLength(data[consumed:])
		consumed += n
		elements := make([]DlmsData, count)
		for i := 0; i < int(count); i++ {
			elem, nc, err := DlmsDataFromBytes(data[consumed:])
			if err != nil {
				return nil, 0, err
			}
			elements[i] = elem
			consumed += nc
		}
		return StructureData(elements), consumed, nil
	default:
		return nil, 0, fmt.Errorf("dlms: unknown tag 0x%02x", tag)
	}
}

func decodeVariableLength(data []byte) (uint32, int) {
	if len(data) == 0 {
		return 0, 0
	}
	first := data[0]
	if first&0x80 == 0 {
		return uint32(first), 1
	}
	count := first & 0x7f
	if len(data) < int(count)+1 {
		return 0, 1
	}
	var val uint32
	for i := 0; i < int(count); i++ {
		val = val<<8 | uint32(data[1+i])
	}
	return val, int(count) + 1
}

func EncodeVariableLength(val uint32) []byte {
	if val < 0x80 {
		return []byte{byte(val)}
	}
	if val < 0x100 {
		return []byte{0x81, byte(val)}
	}
	if val < 0x10000 {
		return []byte{0x82, byte(val >> 8), byte(val)}
	}
	if val < 0x1000000 {
		return []byte{0x83, byte(val >> 16), byte(val >> 8), byte(val)}
	}
	return []byte{0x84, byte(val >> 24), byte(val >> 16), byte(val >> 8), byte(val)}
}

// --- Individual data types ---

// NullData represents a null value (tag 0).
type NullData struct{}

func (NullData) Tag() byte             { return TagNull }
func (NullData) ToBytes() []byte       { return []byte{TagNull} }
func (NullData) ToPython() interface{} { return nil }

// BooleanData represents a boolean (tag 3).
type BooleanData bool

func (BooleanData) Tag() byte { return TagBoolean }
func (b BooleanData) ToBytes() []byte {
	if b {
		return []byte{TagBoolean, 1}
	}
	return []byte{TagBoolean, 0}
}
func (b BooleanData) ToPython() interface{} { return bool(b) }

// BitStringData represents a bit string (tag 4).
type BitStringData []byte

func (BitStringData) Tag() byte { return TagBitString }
func (b BitStringData) ToBytes() []byte {
	out := []byte{TagBitString}
	out = append(out, EncodeVariableLength(uint32(len(b)))...)
	out = append(out, b...)
	return out
}
func (b BitStringData) ToPython() interface{} { return []byte(b) }

// DoubleLongData represents a signed 32-bit integer (tag 5).
type DoubleLongData int32

func (DoubleLongData) Tag() byte { return TagDoubleLong }
func (d DoubleLongData) ToBytes() []byte {
	buf := make([]byte, 5)
	buf[0] = TagDoubleLong
	binary.BigEndian.PutUint32(buf[1:], uint32(d))
	return buf
}
func (d DoubleLongData) ToPython() interface{} { return int32(d) }

// DoubleLongUnsignedData represents an unsigned 32-bit integer (tag 6).
type DoubleLongUnsignedData uint32

func (DoubleLongUnsignedData) Tag() byte { return TagDoubleLongUnsigned }
func (d DoubleLongUnsignedData) ToBytes() []byte {
	buf := make([]byte, 5)
	buf[0] = TagDoubleLongUnsigned
	binary.BigEndian.PutUint32(buf[1:], uint32(d))
	return buf
}
func (d DoubleLongUnsignedData) ToPython() interface{} { return uint32(d) }

// OctetStringData represents a byte string (tag 9).
type OctetStringData []byte

func (OctetStringData) Tag() byte { return TagOctetString }
func (o OctetStringData) ToBytes() []byte {
	out := []byte{TagOctetString}
	out = append(out, EncodeVariableLength(uint32(len(o)))...)
	out = append(out, o...)
	return out
}
func (o OctetStringData) ToPython() interface{} { return []byte(o) }

// VisibleStringData represents an ASCII string (tag 10).
type VisibleStringData string

func (VisibleStringData) Tag() byte { return TagVisibleString }
func (v VisibleStringData) ToBytes() []byte {
	out := []byte{TagVisibleString}
	b := []byte(string(v))
	out = append(out, EncodeVariableLength(uint32(len(b)))...)
	out = append(out, b...)
	return out
}
func (v VisibleStringData) ToPython() interface{} { return string(v) }

// UTF8StringData represents a UTF-8 string (tag 12).
type UTF8StringData string

func (UTF8StringData) Tag() byte { return TagUTF8String }
func (u UTF8StringData) ToBytes() []byte {
	out := []byte{TagUTF8String}
	b := []byte(string(u))
	out = append(out, EncodeVariableLength(uint32(len(b)))...)
	out = append(out, b...)
	return out
}
func (u UTF8StringData) ToPython() interface{} { return string(u) }

// BCDData represents BCD encoded data (tag 13).
type BCDData []byte

func (BCDData) Tag() byte { return TagBCD }
func (b BCDData) ToBytes() []byte {
	out := []byte{TagBCD}
	out = append(out, EncodeVariableLength(uint32(len(b)))...)
	out = append(out, b...)
	return out
}
func (b BCDData) ToPython() interface{} { return []byte(b) }

// IntegerData represents a signed 8-bit integer (tag 15).
type IntegerData int8

func (IntegerData) Tag() byte               { return TagInteger }
func (i IntegerData) ToBytes() []byte       { return []byte{TagInteger, byte(i)} }
func (i IntegerData) ToPython() interface{} { return int8(i) }

// LongData represents a signed 16-bit integer (tag 16).
type LongData int16

func (LongData) Tag() byte { return TagLong }
func (l LongData) ToBytes() []byte {
	buf := make([]byte, 3)
	buf[0] = TagLong
	binary.BigEndian.PutUint16(buf[1:], uint16(l))
	return buf
}
func (l LongData) ToPython() interface{} { return int16(l) }

// UnsignedIntegerData represents an unsigned 8-bit integer (tag 17).
type UnsignedIntegerData byte

func (UnsignedIntegerData) Tag() byte               { return TagUnsignedInteger }
func (u UnsignedIntegerData) ToBytes() []byte       { return []byte{TagUnsignedInteger, byte(u)} }
func (u UnsignedIntegerData) ToPython() interface{} { return byte(u) }

// UnsignedLongData represents an unsigned 16-bit integer (tag 18).
type UnsignedLongData uint16

func (UnsignedLongData) Tag() byte { return TagUnsignedLong }
func (u UnsignedLongData) ToBytes() []byte {
	buf := make([]byte, 3)
	buf[0] = TagUnsignedLong
	binary.BigEndian.PutUint16(buf[1:], uint16(u))
	return buf
}
func (u UnsignedLongData) ToPython() interface{} { return uint16(u) }

// CompactArrayData represents a compact array (tag 19).
type CompactArrayData []byte

func (CompactArrayData) Tag() byte { return TagCompactArray }
func (c CompactArrayData) ToBytes() []byte {
	out := []byte{TagCompactArray}
	out = append(out, EncodeVariableLength(uint32(len(c)))...)
	out = append(out, c...)
	return out
}
func (c CompactArrayData) ToPython() interface{} { return []byte(c) }

// Long64Data represents a signed 64-bit integer (tag 20).
type Long64Data int64

func (Long64Data) Tag() byte { return TagLong64 }
func (l Long64Data) ToBytes() []byte {
	buf := make([]byte, 9)
	buf[0] = TagLong64
	binary.BigEndian.PutUint64(buf[1:], uint64(l))
	return buf
}
func (l Long64Data) ToPython() interface{} { return int64(l) }

// UnsignedLong64Data represents an unsigned 64-bit integer (tag 21).
type UnsignedLong64Data uint64

func (UnsignedLong64Data) Tag() byte { return TagUnsignedLong64 }
func (u UnsignedLong64Data) ToBytes() []byte {
	buf := make([]byte, 9)
	buf[0] = TagUnsignedLong64
	binary.BigEndian.PutUint64(buf[1:], uint64(u))
	return buf
}
func (u UnsignedLong64Data) ToPython() interface{} { return uint64(u) }

// EnumData represents an enum value (tag 22).
type EnumData byte

func (EnumData) Tag() byte               { return TagEnum }
func (e EnumData) ToBytes() []byte       { return []byte{TagEnum, byte(e)} }
func (e EnumData) ToPython() interface{} { return byte(e) }

// Float32Data represents an IEEE 754 float32 (tag 23).
type Float32Data float32

func (Float32Data) Tag() byte { return TagFloat32 }
func (f Float32Data) ToBytes() []byte {
	buf := make([]byte, 5)
	buf[0] = TagFloat32
	binary.BigEndian.PutUint32(buf[1:], math.Float32bits(float32(f)))
	return buf
}
func (f Float32Data) ToPython() interface{} { return float32(f) }

// Float64Data represents an IEEE 754 float64 (tag 24).
type Float64Data float64

func (Float64Data) Tag() byte { return TagFloat64 }
func (f Float64Data) ToBytes() []byte {
	buf := make([]byte, 9)
	buf[0] = TagFloat64
	binary.BigEndian.PutUint64(buf[1:], math.Float64bits(float64(f)))
	return buf
}
func (f Float64Data) ToPython() interface{} { return float64(f) }

// DateTimeData represents a COSEM DateTime (tag 25).
type DateTimeData struct {
	Value CosemDateTime
}

func (DateTimeData) Tag() byte { return TagDateTime }
func (d DateTimeData) ToBytes() []byte {
	return append([]byte{TagDateTime}, d.Value.ToBytes()...)
}
func (d DateTimeData) ToPython() interface{} { return d.Value }

// DateData represents a COSEM Date (tag 26).
type DateData struct {
	Value CosemDate
}

func (DateData) Tag() byte { return TagDate }
func (d DateData) ToBytes() []byte {
	return append([]byte{TagDate}, d.Value.ToBytes()...)
}
func (d DateData) ToPython() interface{} { return d.Value }

// TimeData represents a COSEM Time (tag 27).
type TimeData struct {
	Value CosemTime
}

func (TimeData) Tag() byte { return TagTime }
func (t TimeData) ToBytes() []byte {
	return append([]byte{TagTime}, t.Value.ToBytes()...)
}
func (t TimeData) ToPython() interface{} { return t.Value }

// DontCareData represents a don't-care value (tag 255).
type DontCareData struct{}

func (DontCareData) Tag() byte             { return TagDontCare }
func (DontCareData) ToBytes() []byte       { return []byte{TagDontCare} }
func (DontCareData) ToPython() interface{} { return nil }

// ArrayData represents an array of DlmsData elements (tag 1).
type ArrayData []DlmsData

func (ArrayData) Tag() byte { return TagArray }
func (a ArrayData) ToBytes() []byte {
	out := []byte{TagArray}
	out = append(out, EncodeVariableLength(uint32(len(a)))...)
	for _, elem := range a {
		out = append(out, elem.ToBytes()...)
	}
	return out
}
func (a ArrayData) ToPython() interface{} {
	result := make([]interface{}, len(a))
	for i, elem := range a {
		result[i] = elem.ToPython()
	}
	return result
}

// StructureData represents a structure of DlmsData elements (tag 2).
type StructureData []DlmsData

func (StructureData) Tag() byte { return TagStructure }
func (s StructureData) ToBytes() []byte {
	out := []byte{TagStructure}
	out = append(out, EncodeVariableLength(uint32(len(s)))...)
	for _, elem := range s {
		out = append(out, elem.ToBytes()...)
	}
	return out
}
func (s StructureData) ToPython() interface{} {
	result := make([]interface{}, len(s))
	for i, elem := range s {
		result[i] = elem.ToPython()
	}
	return result
}

// NewCosemDateTime returns a CosemDateTime from time.Time.
func NewCosemDateTime(t time.Time) CosemDateTime {
	return CosemDateTime{
		Year:   uint16(t.Year()),
		Month:  uint8(t.Month()),
		Day:    uint8(t.Day()),
		Hour:   uint8(t.Hour()),
		Minute: uint8(t.Minute()),
		Second: uint8(t.Second()),
	}
}

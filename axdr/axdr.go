package axdr

import (
	"encoding/binary"
	"fmt"

	"github.com/ViewWay/dlms-cosem-go/core"
)

// Decoder decodes AXDR-encoded data.
type Decoder struct {
	data []byte
	pos  int
}

// NewDecoder creates a new AXDR decoder.
func NewDecoder(data []byte) *Decoder {
	return &Decoder{data: data, pos: 0}
}

// Remaining returns bytes remaining.
func (d *Decoder) Remaining() int {
	return len(d.data) - d.pos
}

// Empty returns true if no data remains.
func (d *Decoder) Empty() bool {
	return d.pos >= len(d.data)
}

// GetBytes reads n bytes.
func (d *Decoder) GetBytes(n int) ([]byte, error) {
	if d.pos+n > len(d.data) {
		return nil, fmt.Errorf("axdr: insufficient data (need %d, have %d)", n, d.Remaining())
	}
	b := d.data[d.pos : d.pos+n]
	d.pos += n
	return b, nil
}

// GetByte reads one byte.
func (d *Decoder) GetByte() (byte, error) {
	b, err := d.GetBytes(1)
	if err != nil {
		return 0, err
	}
	return b[0], nil
}

// GetAXDRLength reads a variable-length integer.
func (d *Decoder) GetAXDRLength() (uint32, error) {
	if d.Empty() {
		return 0, fmt.Errorf("axdr: no data for length")
	}
	first, err := d.GetByte()
	if err != nil {
		return 0, err
	}
	if first&0x80 == 0 {
		return uint32(first), nil
	}
	count := first & 0x7F
	if count == 0 || count > 4 {
		return 0, fmt.Errorf("axdr: invalid length encoding 0x%02x", first)
	}
	lengthData, err := d.GetBytes(int(count))
	if err != nil {
		return 0, err
	}
	var val uint32
	for _, b := range lengthData {
		val = val<<8 | uint32(b)
	}
	return val, nil
}

// GetOptional reads an optional value. Returns nil if indicator is 0x00.
func (d *Decoder) GetOptional() ([]byte, error) {
	indicator, err := d.GetByte()
	if err != nil {
		return nil, err
	}
	if indicator == 0x00 {
		return nil, nil
	}
	return d.GetBytes(int(indicator))
}

// GetDefault reads a defaultable value. Returns the default if indicator is 0x00.
func (d *Decoder) GetDefault(defaultValue []byte) ([]byte, error) {
	indicator, err := d.GetByte()
	if err != nil {
		return nil, err
	}
	if indicator == 0x00 {
		return defaultValue, nil
	}
	return d.GetBytes(int(indicator))
}

// GetUint8 reads an unsigned 8-bit integer.
func (d *Decoder) GetUint8() (uint8, error) {
	b, err := d.GetByte()
	return b, err
}

// GetInt8 reads a signed 8-bit integer.
func (d *Decoder) GetInt8() (int8, error) {
	b, err := d.GetByte()
	return int8(b), err
}

// GetUint16 reads an unsigned 16-bit integer (big-endian).
func (d *Decoder) GetUint16() (uint16, error) {
	b, err := d.GetBytes(2)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint16(b), nil
}

// GetInt16 reads a signed 16-bit integer.
func (d *Decoder) GetInt16() (int16, error) {
	v, err := d.GetUint16()
	return int16(v), err
}

// GetUint32 reads an unsigned 32-bit integer.
func (d *Decoder) GetUint32() (uint32, error) {
	b, err := d.GetBytes(4)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint32(b), nil
}

// GetInt32 reads a signed 32-bit integer.
func (d *Decoder) GetInt32() (int32, error) {
	v, err := d.GetUint32()
	return int32(v), err
}

// GetUint64 reads an unsigned 64-bit integer.
func (d *Decoder) GetUint64() (uint64, error) {
	b, err := d.GetBytes(8)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint64(b), nil
}

// GetInt64 reads a signed 64-bit integer.
func (d *Decoder) GetInt64() (int64, error) {
	v, err := d.GetUint64()
	return int64(v), err
}

// GetRaw reads remaining bytes.
func (d *Decoder) GetRaw() []byte {
	remaining := d.data[d.pos:]
	d.pos = len(d.data)
	return remaining
}

// Peek returns next byte without consuming.
func (d *Decoder) Peek() (byte, error) {
	if d.Empty() {
		return 0, fmt.Errorf("axdr: no data")
	}
	return d.data[d.pos], nil
}

// GetDlmsData reads a single DLMS data element.
func (d *Decoder) GetDlmsData() (core.DlmsData, error) {
	data := d.data[d.pos:]
	elem, consumed, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return nil, err
	}
	d.pos += consumed
	return elem, nil
}

// GetDlmsDataArray reads a DLMS array.
func (d *Decoder) GetDlmsDataArray() ([]core.DlmsData, error) {
	count, err := d.GetAXDRLength()
	if err != nil {
		return nil, err
	}
	result := make([]core.DlmsData, count)
	for i := 0; i < int(count); i++ {
		elem, err := d.GetDlmsData()
		if err != nil {
			return nil, err
		}
		result[i] = elem
	}
	return result, nil
}

// GetDlmsDataStructure reads a DLMS structure.
func (d *Decoder) GetDlmsDataStructure() ([]core.DlmsData, error) {
	return d.GetDlmsDataArray()
}

// Encoder encodes AXDR data.
type Encoder struct {
	buf []byte
}

// NewEncoder creates a new AXDR encoder.
func NewEncoder() *Encoder {
	return &Encoder{buf: make([]byte, 0, 64)}
}

// Bytes returns encoded data.
func (e *Encoder) Bytes() []byte {
	return e.buf
}

// Len returns encoded length.
func (e *Encoder) Len() int {
	return len(e.buf)
}

// WriteBytes writes raw bytes.
func (e *Encoder) WriteBytes(data []byte) {
	e.buf = append(e.buf, data...)
}

// WriteByte writes a single byte.
func (e *Encoder) WriteByte(b byte) error {
	e.buf = append(e.buf, b)
	return nil
}

// WriteUint8 writes an unsigned 8-bit integer.
func (e *Encoder) WriteUint8(v uint8) {
	e.WriteByte(v)
}

// WriteInt8 writes a signed 8-bit integer.
func (e *Encoder) WriteInt8(v int8) {
	e.WriteByte(uint8(v))
}

// WriteUint16 writes an unsigned 16-bit integer (big-endian).
func (e *Encoder) WriteUint16(v uint16) {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, v)
	e.WriteBytes(b)
}

// WriteInt16 writes a signed 16-bit integer.
func (e *Encoder) WriteInt16(v int16) {
	e.WriteUint16(uint16(v))
}

// WriteUint32 writes an unsigned 32-bit integer.
func (e *Encoder) WriteUint32(v uint32) {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, v)
	e.WriteBytes(b)
}

// WriteInt32 writes a signed 32-bit integer.
func (e *Encoder) WriteInt32(v int32) {
	e.WriteUint32(uint32(v))
}

// WriteUint64 writes an unsigned 64-bit integer.
func (e *Encoder) WriteUint64(v uint64) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	e.WriteBytes(b)
}

// WriteInt64 writes a signed 64-bit integer.
func (e *Encoder) WriteInt64(v int64) {
	e.WriteUint64(uint64(v))
}

// WriteAXDRLength writes a variable-length integer.
func (e *Encoder) WriteAXDRLength(v uint32) {
	e.WriteBytes(core.EncodeVariableLength(v))
}

// WriteOptional writes an optional value.
func (e *Encoder) WriteOptional(data []byte) {
	if data == nil {
		e.WriteByte(0x00)
	} else {
		e.WriteByte(byte(len(data)))
		e.WriteBytes(data)
	}
}

// WriteDefault writes a defaultable value.
func (e *Encoder) WriteDefault(data []byte, defaultValue []byte) {
	if data == nil || bytesEqual(data, defaultValue) {
		e.WriteByte(0x00)
	} else {
		e.WriteByte(byte(len(data)))
		e.WriteBytes(data)
	}
}

// WriteDlmsData writes a DLMS data element.
func (e *Encoder) WriteDlmsData(data core.DlmsData) {
	e.WriteBytes(data.ToBytes())
}

// WriteDlmsArray writes a DLMS array.
func (e *Encoder) WriteDlmsArray(elements []core.DlmsData) {
	e.WriteByte(core.TagArray)
	e.WriteAXDRLength(uint32(len(elements)))
	for _, elem := range elements {
		e.WriteDlmsData(elem)
	}
}

// WriteDlmsStructure writes a DLMS structure.
func (e *Encoder) WriteDlmsStructure(elements []core.DlmsData) {
	e.WriteByte(core.TagStructure)
	e.WriteAXDRLength(uint32(len(elements)))
	for _, elem := range elements {
		e.WriteDlmsData(elem)
	}
}

// WriteOctetString writes an octet string (with length).
func (e *Encoder) WriteOctetString(data []byte) {
	e.WriteByte(core.TagOctetString)
	e.WriteAXDRLength(uint32(len(data)))
	e.WriteBytes(data)
}

// WriteVisibleString writes a visible string.
func (e *Encoder) WriteVisibleString(s string) {
	e.WriteByte(core.TagVisibleString)
	b := []byte(s)
	e.WriteAXDRLength(uint32(len(b)))
	e.WriteBytes(b)
}

// WriteDateTime writes a COSEM DateTime.
func (e *Encoder) WriteDateTime(dt core.CosemDateTime) {
	e.WriteByte(core.TagDateTime)
	e.WriteBytes(dt.ToBytes())
}

// WriteDate writes a COSEM Date.
func (e *Encoder) WriteDate(d core.CosemDate) {
	e.WriteByte(core.TagDate)
	e.WriteBytes(d.ToBytes())
}

// WriteTime writes a COSEM Time.
func (e *Encoder) WriteTime(t core.CosemTime) {
	e.WriteByte(core.TagTime)
	e.WriteBytes(t.ToBytes())
}

func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

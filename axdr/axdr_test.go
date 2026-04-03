package axdr

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestDecoder_GetBytes(t *testing.T) {
	d := NewDecoder([]byte{0x01, 0x02, 0x03})
	b, err := d.GetBytes(2)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(b, []byte{0x01, 0x02}) {
		t.Errorf("got %v", b)
	}
	if d.Remaining() != 1 {
		t.Errorf("remaining=%d", d.Remaining())
	}
}

func TestDecoder_GetBytes_Insufficient(t *testing.T) {
	d := NewDecoder([]byte{0x01})
	_, err := d.GetBytes(5)
	if err == nil {
		t.Error("expected error")
	}
}

func TestDecoder_GetByte(t *testing.T) {
	d := NewDecoder([]byte{0x42})
	b, err := d.GetByte()
	if err != nil {
		t.Fatal(err)
	}
	if b != 0x42 {
		t.Errorf("got %02x", b)
	}
}

func TestDecoder_Empty(t *testing.T) {
	d := NewDecoder([]byte{})
	if !d.Empty() {
		t.Error("expected empty")
	}
	if d.Remaining() != 0 {
		t.Errorf("remaining=%d", d.Remaining())
	}
}

func TestDecoder_GetAXDRLength_SingleByte(t *testing.T) {
	tests := []struct {
		data []byte
		want uint32
	}{
		{[]byte{0x00}, 0},
		{[]byte{0x7F}, 127},
		{[]byte{0x05}, 5},
	}
	for _, tc := range tests {
		d := NewDecoder(tc.data)
		v, err := d.GetAXDRLength()
		if err != nil {
			t.Errorf("GetAXDRLength(%v): %v", tc.data, err)
			continue
		}
		if v != tc.want {
			t.Errorf("got %d, want %d", v, tc.want)
		}
	}
}

func TestDecoder_GetAXDRLength_MultiByte(t *testing.T) {
	tests := []struct {
		data []byte
		want uint32
	}{
		{[]byte{0x81, 0x80}, 128},
		{[]byte{0x81, 0xFF}, 255},
		{[]byte{0x82, 0x01, 0x00}, 256},
		{[]byte{0x82, 0xFF, 0xFF}, 65535},
		{[]byte{0x83, 0x01, 0x00, 0x00}, 65536},
		{[]byte{0x84, 0x01, 0x00, 0x00, 0x00}, 16777216},
	}
	for _, tc := range tests {
		d := NewDecoder(tc.data)
		v, err := d.GetAXDRLength()
		if err != nil {
			t.Errorf("GetAXDRLength(%v): %v", tc.data, err)
			continue
		}
		if v != tc.want {
			t.Errorf("got %d, want %d", v, tc.want)
		}
	}
}

func TestDecoder_GetAXDRLength_Invalid(t *testing.T) {
	d := NewDecoder([]byte{})
	_, err := d.GetAXDRLength()
	if err == nil {
		t.Error("expected error for empty data")
	}
	d = NewDecoder([]byte{0x80}) // count=0 invalid
	_, err = d.GetAXDRLength()
	if err == nil {
		t.Error("expected error for count=0")
	}
}

func TestDecoder_GetOptional_Present(t *testing.T) {
	d := NewDecoder([]byte{0x02, 0xAB, 0xCD})
	data, err := d.GetOptional()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(data, []byte{0xAB, 0xCD}) {
		t.Errorf("got %v", data)
	}
}

func TestDecoder_GetOptional_Absent(t *testing.T) {
	d := NewDecoder([]byte{0x00})
	data, err := d.GetOptional()
	if err != nil {
		t.Fatal(err)
	}
	if data != nil {
		t.Errorf("got %v", data)
	}
}

func TestDecoder_GetDefault(t *testing.T) {
	d := NewDecoder([]byte{0x00})
	data, err := d.GetDefault([]byte{0x01, 0x02})
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(data, []byte{0x01, 0x02}) {
		t.Errorf("got %v", data)
	}
}

func TestDecoder_GetDefault_NonDefault(t *testing.T) {
	d := NewDecoder([]byte{0x02, 0xAB, 0xCD})
	data, err := d.GetDefault([]byte{0x01, 0x02})
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(data, []byte{0xAB, 0xCD}) {
		t.Errorf("got %v", data)
	}
}

func TestDecoder_IntegerTypes(t *testing.T) {
	data := make([]byte, 0)
	data = append(data, 0xFF)                                           // uint8 = 255
	data = append(data, 0x80)                                           // int8 = -128
	data = append(data, 0x12, 0x34)                                     // uint16 = 0x1234
	data = append(data, 0xFF, 0xFF)                                     // int16 = -1
	data = append(data, 0x01, 0x02, 0x03, 0x04)                         // uint32
	data = append(data, 0xFF, 0xFF, 0xFF, 0xFF)                         // int32 = -1
	data = append(data, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08) // uint64
	data = append(data, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF) // int64 = -1

	d := NewDecoder(data)

	v, _ := d.GetUint8()
	if v != 255 {
		t.Errorf("uint8: %d", v)
	}
	v8, _ := d.GetInt8()
	if v8 != -128 {
		t.Errorf("int8: %d", v8)
	}
	v16, _ := d.GetUint16()
	if v16 != 0x1234 {
		t.Errorf("uint16: %d", v16)
	}
	v16s, _ := d.GetInt16()
	if v16s != -1 {
		t.Errorf("int16: %d", v16s)
	}
	v32, _ := d.GetUint32()
	if v32 != 0x01020304 {
		t.Errorf("uint32: %d", v32)
	}
	v32s, _ := d.GetInt32()
	if v32s != -1 {
		t.Errorf("int32: %d", v32s)
	}
	v64, _ := d.GetUint64()
	if v64 != 0x0102030405060708 {
		t.Errorf("uint64: %d", v64)
	}
	v64s, _ := d.GetInt64()
	if v64s != -1 {
		t.Errorf("int64: %d", v64s)
	}
}

func TestDecoder_IntegerTypes_Errors(t *testing.T) {
	d := NewDecoder([]byte{0x01})
	_, err := d.GetUint16()
	if err == nil {
		t.Error("expected error")
	}
	_, err = d.GetInt16()
	if err == nil {
		t.Error("expected error")
	}
	_, err = d.GetUint32()
	if err == nil {
		t.Error("expected error")
	}
	_, err = d.GetInt32()
	if err == nil {
		t.Error("expected error")
	}
	_, err = d.GetUint64()
	if err == nil {
		t.Error("expected error")
	}
	_, err = d.GetInt64()
	if err == nil {
		t.Error("expected error")
	}
}

func TestDecoder_Peek(t *testing.T) {
	d := NewDecoder([]byte{0x42, 0x43})
	b, err := d.Peek()
	if err != nil {
		t.Fatal(err)
	}
	if b != 0x42 {
		t.Errorf("got %02x", b)
	}
	// Should not consume
	b, err = d.Peek()
	if err != nil || b != 0x42 {
		t.Error("peek consumed data")
	}
}

func TestDecoder_Peek_Empty(t *testing.T) {
	d := NewDecoder([]byte{})
	_, err := d.Peek()
	if err == nil {
		t.Error("expected error")
	}
}

func TestDecoder_GetRaw(t *testing.T) {
	d := NewDecoder([]byte{0x01, 0x02, 0x03})
	_, _ = d.GetByte()
	raw := d.GetRaw()
	if !bytes.Equal(raw, []byte{0x02, 0x03}) {
		t.Errorf("got %v", raw)
	}
	if !d.Empty() {
		t.Error("should be empty")
	}
}

func TestEncoder_Basic(t *testing.T) {
	e := NewEncoder()
	e.WriteByte(0x01)
	e.WriteBytes([]byte{0x02, 0x03})
	b := e.Bytes()
	if !bytes.Equal(b, []byte{0x01, 0x02, 0x03}) {
		t.Errorf("got %v", b)
	}
}

func TestEncoder_Integers(t *testing.T) {
	e := NewEncoder()
	e.WriteUint8(255)
	e.WriteInt8(-1)
	e.WriteUint16(0x1234)
	e.WriteInt16(-1)
	e.WriteUint32(0xDEADBEEF)
	e.WriteInt32(-1)
	e.WriteUint64(0x0102030405060708)
	e.WriteInt64(-1)
	b := e.Bytes()
	if len(b) != 1+1+2+2+4+4+8+8 {
		t.Errorf("len=%d", len(b))
	}
}

func TestEncoder_WriteAXDRLength(t *testing.T) {
	tests := []uint32{0, 1, 127, 128, 255, 256, 65535, 65536, 16777215}
	for _, v := range tests {
		e := NewEncoder()
		e.WriteAXDRLength(v)
		d := NewDecoder(e.Bytes())
		got, err := d.GetAXDRLength()
		if err != nil {
			t.Errorf("roundtrip %d: %v", v, err)
			continue
		}
		if got != v {
			t.Errorf("roundtrip %d: got %d", v, got)
		}
	}
}

func TestEncoder_WriteOptional(t *testing.T) {
	e := NewEncoder()
	e.WriteOptional(nil)
	e.WriteOptional([]byte{0x01, 0x02})
	b := e.Bytes()
	if !bytes.Equal(b, []byte{0x00, 0x02, 0x01, 0x02}) {
		t.Errorf("got %v", b)
	}
}

func TestEncoder_WriteDefault(t *testing.T) {
	e := NewEncoder()
	e.WriteDefault(nil, []byte{0x01})
	b := e.Bytes()
	if b[0] != 0x00 {
		t.Errorf("got %v", b)
	}
}

func TestEncoder_WriteDlmsData(t *testing.T) {
	e := NewEncoder()
	e.WriteDlmsData(core.IntegerData(42))
	b := e.Bytes()
	if !bytes.Equal(b, core.IntegerData(42).ToBytes()) {
		t.Errorf("got %v", b)
	}
}

func TestEncoder_WriteDlmsArray(t *testing.T) {
	e := NewEncoder()
	e.WriteDlmsArray([]core.DlmsData{core.IntegerData(1), core.IntegerData(2)})
	b := e.Bytes()
	if b[0] != core.TagArray {
		t.Errorf("tag=%02x", b[0])
	}
}

func TestEncoder_WriteDlmsStructure(t *testing.T) {
	e := NewEncoder()
	e.WriteDlmsStructure([]core.DlmsData{core.IntegerData(1), core.VisibleStringData("test")})
	b := e.Bytes()
	if b[0] != core.TagStructure {
		t.Errorf("tag=%02x", b[0])
	}
}

func TestEncoder_WriteOctetString(t *testing.T) {
	e := NewEncoder()
	e.WriteOctetString([]byte{0x01, 0x02})
	b := e.Bytes()
	parsed, _, err := core.DlmsDataFromBytes(b)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(parsed.(core.OctetStringData), []byte{0x01, 0x02}) {
		t.Error("mismatch")
	}
}

func TestEncoder_WriteVisibleString(t *testing.T) {
	e := NewEncoder()
	e.WriteVisibleString("hello")
	b := e.Bytes()
	parsed, _, err := core.DlmsDataFromBytes(b)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.(core.VisibleStringData) != "hello" {
		t.Error("mismatch")
	}
}

func TestEncoder_WriteDateTime(t *testing.T) {
	e := NewEncoder()
	dt := core.CosemDateTime{Year: 2024, Month: 1, Day: 1}
	e.WriteDateTime(dt)
	b := e.Bytes()
	parsed, _, err := core.DlmsDataFromBytes(b)
	if err != nil {
		t.Fatal(err)
	}
	pdt := parsed.(core.DateTimeData)
	if pdt.Value.Year != 2024 {
		t.Errorf("got %v", pdt.Value)
	}
}

func TestEncoder_WriteDate(t *testing.T) {
	e := NewEncoder()
	d := core.CosemDate{Year: 2024, Month: 6, Day: 15}
	e.WriteDate(d)
	b := e.Bytes()
	parsed, _, err := core.DlmsDataFromBytes(b)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.(core.DateData).Value.Year != 2024 {
		t.Error("mismatch")
	}
}

func TestEncoder_WriteTime(t *testing.T) {
	e := NewEncoder()
	tm := core.CosemTime{Hour: 12, Minute: 30, Second: 45}
	e.WriteTime(tm)
	b := e.Bytes()
	parsed, _, err := core.DlmsDataFromBytes(b)
	if err != nil {
		t.Fatal(err)
	}
	pt := parsed.(core.TimeData)
	if pt.Value.Hour != 12 {
		t.Error("mismatch")
	}
}

func TestDecoder_GetDlmsData(t *testing.T) {
	data := core.IntegerData(42).ToBytes()
	d := NewDecoder(data)
	elem, err := d.GetDlmsData()
	if err != nil {
		t.Fatal(err)
	}
	if elem.(core.IntegerData) != 42 {
		t.Errorf("got %v", elem)
	}
}

func TestDecoder_GetDlmsDataArray(t *testing.T) {
	arr := core.ArrayData{core.IntegerData(1), core.IntegerData(2), core.IntegerData(3)}
	data := arr.ToBytes()
	d := NewDecoder(data)
	// Skip tag and length
	_, _ = d.GetByte() // tag
	elements, err := d.GetDlmsDataArray()
	if err != nil {
		t.Fatal(err)
	}
	if len(elements) != 3 {
		t.Fatalf("len=%d", len(elements))
	}
}

func TestDecoder_GetDlmsData_Error(t *testing.T) {
	d := NewDecoder([]byte{0x07}) // unknown tag
	_, err := d.GetDlmsData()
	if err == nil {
		t.Error("expected error")
	}
}

func TestEncoder_Len(t *testing.T) {
	e := NewEncoder()
	if e.Len() != 0 {
		t.Error("expected 0")
	}
	e.WriteByte(0x01)
	if e.Len() != 1 {
		t.Error("expected 1")
	}
}

func TestRoundtrip_Uint8(t *testing.T) {
	e := NewEncoder()
	e.WriteUint8(200)
	d := NewDecoder(e.Bytes())
	v, err := d.GetUint8()
	if err != nil || v != 200 {
		t.Error("roundtrip failed")
	}
}

func TestRoundtrip_Uint16(t *testing.T) {
	e := NewEncoder()
	e.WriteUint16(50000)
	d := NewDecoder(e.Bytes())
	v, err := d.GetUint16()
	if err != nil || v != 50000 {
		t.Error("roundtrip failed")
	}
}

func TestRoundtrip_Uint32(t *testing.T) {
	e := NewEncoder()
	e.WriteUint32(3000000000)
	d := NewDecoder(e.Bytes())
	v, err := d.GetUint32()
	if err != nil || v != 3000000000 {
		t.Error("roundtrip failed")
	}
}

func TestRoundtrip_Uint64(t *testing.T) {
	e := NewEncoder()
	e.WriteUint64(0xFFFFFFFFFFFFFFFF)
	d := NewDecoder(e.Bytes())
	v, err := d.GetUint64()
	if err != nil || v != 0xFFFFFFFFFFFFFFFF {
		t.Error("roundtrip failed")
	}
}

func TestRoundtrip_DlmsData_Multiple(t *testing.T) {
	e := NewEncoder()
	e.WriteDlmsData(core.DoubleLongData(-1000))
	e.WriteDlmsData(core.VisibleStringData("test"))
	e.WriteDlmsData(core.BooleanData(true))

	d := NewDecoder(e.Bytes())

	elem1, _ := d.GetDlmsData()
	if elem1.(core.DoubleLongData) != -1000 {
		t.Error("elem1")
	}
	elem2, _ := d.GetDlmsData()
	if elem2.(core.VisibleStringData) != "test" {
		t.Error("elem2")
	}
	elem3, _ := d.GetDlmsData()
	if elem3.(core.BooleanData) != true {
		t.Error("elem3")
	}
}

func TestDecoder_GetOptional_Error(t *testing.T) {
	d := NewDecoder([]byte{})
	_, err := d.GetOptional()
	if err == nil {
		t.Error("expected error")
	}
}

func TestDecoder_GetDefault_Error(t *testing.T) {
	d := NewDecoder([]byte{})
	_, err := d.GetDefault(nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestEncoder_WriteDefault_NonDefault(t *testing.T) {
	e := NewEncoder()
	e.WriteDefault([]byte{0x05, 0x06}, []byte{0x01, 0x02})
	b := e.Bytes()
	if b[0] != 0x02 {
		t.Errorf("got %v", b)
	}
}

func TestDecoder_GetByte_Error(t *testing.T) {
	d := NewDecoder([]byte{})
	_, err := d.GetByte()
	if err == nil {
		t.Error("expected error")
	}
}

func TestNewEncoder(t *testing.T) {
	e := NewEncoder()
	if e == nil {
		t.Error("nil encoder")
	}
	if e.Len() != 0 {
		t.Error("non-zero length")
	}
}

func TestNewDecoder(t *testing.T) {
	d := NewDecoder(nil)
	if d == nil {
		t.Error("nil decoder")
	}
	if !d.Empty() {
		t.Error("should be empty")
	}
}

func TestRoundtrip_FloatTypes(t *testing.T) {
	e := NewEncoder()
	data := make([]byte, 12)
	data[0] = core.TagFloat32
	data[1] = 0x42 // float32 tag
	data[2] = core.TagFloat64
	data[3] = 0x44 // float64 tag
	e.WriteBytes(data)

	d := NewDecoder(e.Bytes())
	// Read float32 tag
	b1, _ := d.GetByte()
	if b1 != core.TagFloat32 {
		t.Error("wrong tag")
	}
}

func TestEncoder_ResetBehavior(t *testing.T) {
	e1 := NewEncoder()
	e1.WriteByte(0x01)
	e1.WriteByte(0x02)

	e2 := NewEncoder()
	if e2.Len() != 0 {
		t.Error("new encoder should be empty")
	}
}

func TestEncoder_WriteDefault_NilData(t *testing.T) {
	e := NewEncoder()
	e.WriteDefault(nil, []byte{0x01})
	b := e.Bytes()
	if len(b) != 1 || b[0] != 0x00 {
		t.Errorf("got %v", b)
	}
}

func TestEncoder_WriteOptional_Nil(t *testing.T) {
	e := NewEncoder()
	e.WriteOptional(nil)
	b := e.Bytes()
	if len(b) != 1 || b[0] != 0x00 {
		t.Errorf("got %v", b)
	}
}

// Helper for complex tests
func TestDecoder_MixedReads(t *testing.T) {
	data := []byte{
		0x42,       // byte
		0x01, 0x02, // uint16
		0x03, 0x04, 0x05, 0x06, // uint32
		0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, // uint64
	}
	d := NewDecoder(data)

	b, _ := d.GetByte()
	if b != 0x42 {
		t.Error("byte")
	}
	v16, _ := d.GetUint16()
	if v16 != 0x0102 {
		t.Error("uint16")
	}
	v32, _ := d.GetUint32()
	if v32 != 0x03040506 {
		t.Error("uint32")
	}
	v64, _ := d.GetUint64()
	if v64 != 0x0708090A0B0C0D0E {
		t.Error("uint64")
	}
	if !d.Empty() {
		t.Error("should be empty")
	}
}

func TestDecoder_VariousLengths(t *testing.T) {
	// Test all AXDR length encodings
	tests := []uint32{0, 1, 127, 128, 255, 256, 65535, 65536, 16777215, 16777216, 0xFFFFFFFF}
	for _, v := range tests {
		e := NewEncoder()
		e.WriteAXDRLength(v)
		d := NewDecoder(e.Bytes())
		got, err := d.GetAXDRLength()
		if err != nil {
			t.Errorf("length %d: %v", v, err)
			continue
		}
		if got != v {
			t.Errorf("length %d: got %d", v, got)
		}
	}
}

func TestEncoder_WriteBytes_Multiple(t *testing.T) {
	e := NewEncoder()
	e.WriteBytes([]byte{1, 2})
	e.WriteBytes([]byte{3, 4})
	if !bytes.Equal(e.Bytes(), []byte{1, 2, 3, 4}) {
		t.Error("concat failed")
	}
}

func TestBytesEqual(t *testing.T) {
	if !bytesEqual([]byte{1, 2, 3}, []byte{1, 2, 3}) {
		t.Error("equal")
	}
	if bytesEqual([]byte{1, 2}, []byte{1, 2, 3}) {
		t.Error("different lengths")
	}
	if bytesEqual([]byte{1, 2, 3}, []byte{1, 2, 4}) {
		t.Error("different content")
	}
}

// Ensure error messages are useful
func TestDecoder_ErrorMessages(t *testing.T) {
	d := NewDecoder([]byte{0x01})
	_, err := d.GetBytes(10)
	if err == nil {
		t.Error("expected error")
	}
	msg := err.Error()
	if len(msg) == 0 {
		t.Error("empty error message")
	}
	_ = fmt.Sprintf("error: %v", err) // ensure it's formattable
}

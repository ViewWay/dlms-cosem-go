package core

import (
	"bytes"
	"encoding/binary"
	"math"
	"testing"
	"time"
)

func TestObis_Parse(t *testing.T) {
	tests := []struct {
		input            string
		a, b, c, d, e, f byte
		err              bool
	}{
		{"0.0.1.0.0.255", 0, 0, 1, 0, 0, 255, false},
		{"1.0.1.8.0.255", 1, 0, 1, 8, 0, 255, false},
		{"1.0.1.8.0", 1, 0, 1, 8, 0, 255, false}, // 5-part
		{"255.255.255.255.255.255", 255, 255, 255, 255, 255, 255, false},
		{"0.0.1", 0, 0, 0, 0, 0, 0, true},       // too few
		{"a.b.c.d.e.f", 0, 0, 0, 0, 0, 0, true}, // invalid
	}
	for _, tc := range tests {
		oc, err := ParseObis(tc.input)
		if tc.err {
			if err == nil {
				t.Errorf("ParseObis(%q): expected error", tc.input)
			}
			continue
		}
		if err != nil {
			t.Errorf("ParseObis(%q): %v", tc.input, err)
			continue
		}
		if oc.A() != tc.a || oc.B() != tc.b || oc.C() != tc.c || oc.D() != tc.d || oc.E() != tc.e || oc.F() != tc.f {
			t.Errorf("ParseObis(%q) = %v, want [%d,%d,%d,%d,%d,%d]", tc.input, oc, tc.a, tc.b, tc.c, tc.d, tc.e, tc.f)
		}
	}
}

func TestObis_String(t *testing.T) {
	oc := MustParseObis("0.0.1.0.0.255")
	if s := oc.String(); s != "0.0.1.0.0.255" {
		t.Errorf("String() = %q, want %q", s, "0.0.1.0.0.255")
	}
}

func TestObis_Compare(t *testing.T) {
	a := MustParseObis("1.0.0.0.0.255")
	b := MustParseObis("1.0.1.0.0.255")
	if a.Compare(b) >= 0 {
		t.Error("expected a < b")
	}
	if b.Compare(a) <= 0 {
		t.Error("expected b > a")
	}
	if a.Compare(a) != 0 {
		t.Error("expected a == a")
	}
}

func TestObis_FromBytes(t *testing.T) {
	oc, err := ObisFromBytes([]byte{0, 0, 1, 0, 0, 255})
	if err != nil {
		t.Fatal(err)
	}
	if oc.String() != "0.0.1.0.0.255" {
		t.Errorf("got %s", oc.String())
	}
	_, err = ObisFromBytes([]byte{1, 2, 3})
	if err == nil {
		t.Error("expected error for short bytes")
	}
}

func TestNullData(t *testing.T) {
	d := NullData{}
	if d.Tag() != TagNull {
		t.Error("wrong tag")
	}
	b := d.ToBytes()
	if !bytes.Equal(b, []byte{0}) {
		t.Errorf("got %v", b)
	}
	if d.ToPython() != nil {
		t.Error("expected nil")
	}
}

func TestBooleanData(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		d := BooleanData(true)
		if d.Tag() != TagBoolean {
			t.Error("wrong tag")
		}
		b := d.ToBytes()
		if !bytes.Equal(b, []byte{3, 1}) {
			t.Errorf("got %v", b)
		}
		if d.ToPython() != true {
			t.Error("expected true")
		}
	})
	t.Run("false", func(t *testing.T) {
		d := BooleanData(false)
		if !bytes.Equal(d.ToBytes(), []byte{3, 0}) {
			t.Errorf("got %v", d.ToBytes())
		}
	})
}

func TestBitStringData(t *testing.T) {
	d := BitStringData{0xFF, 0x00}
	if d.Tag() != TagBitString {
		t.Error("wrong tag")
	}
	b := d.ToBytes()
	if !bytes.Equal(b, []byte{4, 2, 0xFF, 0x00}) {
		t.Errorf("got %v", b)
	}
}

func TestDoubleLongData(t *testing.T) {
	d := DoubleLongData(-123456)
	b := d.ToBytes()
	if len(b) != 5 || b[0] != TagDoubleLong {
		t.Errorf("got %v", b)
	}
	// roundtrip
	parsed, n, err := DlmsDataFromBytes(b)
	if err != nil {
		t.Fatal(err)
	}
	if n != 5 {
		t.Errorf("consumed %d", n)
	}
	if parsed.(DoubleLongData) != d {
		t.Errorf("got %v", parsed)
	}
}

func TestDoubleLongUnsignedData(t *testing.T) {
	d := DoubleLongUnsignedData(4000000000)
	b := d.ToBytes()
	parsed, _, err := DlmsDataFromBytes(b)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.(DoubleLongUnsignedData) != d {
		t.Errorf("got %v", parsed)
	}
}

func TestOctetStringData(t *testing.T) {
	d := OctetStringData{0x01, 0x02, 0x03}
	b := d.ToBytes()
	if !bytes.Equal(b, []byte{9, 3, 0x01, 0x02, 0x03}) {
		t.Errorf("got %v", b)
	}
	parsed, _, err := DlmsDataFromBytes(b)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(parsed.(OctetStringData), d) {
		t.Errorf("got %v", parsed)
	}
}

func TestVisibleStringData(t *testing.T) {
	d := VisibleStringData("hello")
	b := d.ToBytes()
	if !bytes.Equal(b, []byte{10, 5, 'h', 'e', 'l', 'l', 'o'}) {
		t.Errorf("got %v", b)
	}
	parsed, _, err := DlmsDataFromBytes(b)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.(VisibleStringData) != d {
		t.Errorf("got %v", parsed)
	}
}

func TestUTF8StringData(t *testing.T) {
	d := UTF8StringData("世界")
	b := d.ToBytes()
	parsed, _, err := DlmsDataFromBytes(b)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.(UTF8StringData) != d {
		t.Errorf("got %v", parsed)
	}
}

func TestIntegerData(t *testing.T) {
	d := IntegerData(-1)
	b := d.ToBytes()
	if !bytes.Equal(b, []byte{15, 0xFF}) {
		t.Errorf("got %v", b)
	}
	parsed, _, err := DlmsDataFromBytes(b)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.(IntegerData) != d {
		t.Errorf("got %v", parsed)
	}
}

func TestLongData(t *testing.T) {
	d := LongData(-1000)
	b := d.ToBytes()
	parsed, _, err := DlmsDataFromBytes(b)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.(LongData) != d {
		t.Errorf("got %v", parsed)
	}
}

func TestUnsignedIntegerData(t *testing.T) {
	d := UnsignedIntegerData(200)
	b := d.ToBytes()
	parsed, _, err := DlmsDataFromBytes(b)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.(UnsignedIntegerData) != d {
		t.Errorf("got %v", parsed)
	}
}

func TestUnsignedLongData(t *testing.T) {
	d := UnsignedLongData(60000)
	b := d.ToBytes()
	parsed, _, err := DlmsDataFromBytes(b)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.(UnsignedLongData) != d {
		t.Errorf("got %v", parsed)
	}
}

func TestLong64Data(t *testing.T) {
	d := Long64Data(-123456789012)
	b := d.ToBytes()
	parsed, _, err := DlmsDataFromBytes(b)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.(Long64Data) != d {
		t.Errorf("got %v", parsed)
	}
}

func TestUnsignedLong64Data(t *testing.T) {
	d := UnsignedLong64Data(18446744073709551615)
	b := d.ToBytes()
	parsed, _, err := DlmsDataFromBytes(b)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.(UnsignedLong64Data) != d {
		t.Errorf("got %v", parsed)
	}
}

func TestEnumData(t *testing.T) {
	d := EnumData(5)
	b := d.ToBytes()
	if !bytes.Equal(b, []byte{22, 5}) {
		t.Errorf("got %v", b)
	}
}

func TestFloat32Data(t *testing.T) {
	d := Float32Data(3.14)
	b := d.ToBytes()
	if len(b) != 5 || b[0] != TagFloat32 {
		t.Errorf("got %v", b)
	}
	parsed, _, err := DlmsDataFromBytes(b)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(float64(parsed.(Float32Data)-d)) > 0.001 {
		t.Errorf("got %v, want %v", parsed, d)
	}
}

func TestFloat64Data(t *testing.T) {
	d := Float64Data(3.14159265358979)
	b := d.ToBytes()
	parsed, _, err := DlmsDataFromBytes(b)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(float64(parsed.(Float64Data))-float64(d)) > 1e-10 {
		t.Errorf("got %v, want %v", parsed, d)
	}
}

func TestArrayData(t *testing.T) {
	a := ArrayData{IntegerData(1), IntegerData(2), IntegerData(3)}
	b := a.ToBytes()
	parsed, _, err := DlmsDataFromBytes(b)
	if err != nil {
		t.Fatal(err)
	}
	arr := parsed.(ArrayData)
	if len(arr) != 3 {
		t.Fatalf("len=%d", len(arr))
	}
	if arr[0].(IntegerData) != 1 || arr[2].(IntegerData) != 3 {
		t.Errorf("got %v", arr)
	}
}

func TestStructureData(t *testing.T) {
	s := StructureData{UnsignedLongData(100), VisibleStringData("test")}
	b := s.ToBytes()
	parsed, _, err := DlmsDataFromBytes(b)
	if err != nil {
		t.Fatal(err)
	}
	st := parsed.(StructureData)
	if len(st) != 2 {
		t.Fatalf("len=%d", len(st))
	}
}

func TestBCDData(t *testing.T) {
	d := BCDData{0x12, 0x34}
	b := d.ToBytes()
	parsed, _, err := DlmsDataFromBytes(b)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(parsed.(BCDData), d) {
		t.Errorf("got %v", parsed)
	}
}

func TestCompactArrayData(t *testing.T) {
	d := CompactArrayData{0x01, 0x02, 0x03}
	b := d.ToBytes()
	parsed, _, err := DlmsDataFromBytes(b)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(parsed.(CompactArrayData), d) {
		t.Errorf("got %v", parsed)
	}
}

func TestDontCareData(t *testing.T) {
	d := DontCareData{}
	if d.Tag() != TagDontCare {
		t.Error("wrong tag")
	}
	if !bytes.Equal(d.ToBytes(), []byte{255}) {
		t.Errorf("got %v", d.ToBytes())
	}
}

func TestDlmsDataFromBytes_UnknownTag(t *testing.T) {
	_, _, err := DlmsDataFromBytes([]byte{0x07}) // tag 7 doesn't exist
	if err == nil {
		t.Error("expected error for unknown tag")
	}
}

func TestDlmsDataFromBytes_InsufficientData(t *testing.T) {
	_, _, err := DlmsDataFromBytes([]byte{})
	if err == nil {
		t.Error("expected error for empty data")
	}
}

func TestDateTimeData(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	dt := NewCosemDateTime(now)
	dd := DateTimeData{Value: dt}
	b := dd.ToBytes()
	parsed, _, err := DlmsDataFromBytes(b)
	if err != nil {
		t.Fatal(err)
	}
	pdt := parsed.(DateTimeData)
	if pdt.Value.Year != dt.Year || pdt.Value.Month != dt.Month || pdt.Value.Day != dt.Day {
		t.Errorf("got %v", pdt.Value)
	}
}

func TestDateData(t *testing.T) {
	d := DateData{Value: CosemDate{Year: 2024, Month: 1, Day: 15}}
	b := d.ToBytes()
	parsed, _, err := DlmsDataFromBytes(b)
	if err != nil {
		t.Fatal(err)
	}
	pd := parsed.(DateData)
	if pd.Value.Year != 2024 || pd.Value.Month != 1 || pd.Value.Day != 15 {
		t.Errorf("got %v", pd.Value)
	}
}

func TestTimeData(t *testing.T) {
	d := TimeData{Value: CosemTime{Hour: 12, Minute: 30, Second: 45}}
	b := d.ToBytes()
	parsed, _, err := DlmsDataFromBytes(b)
	if err != nil {
		t.Fatal(err)
	}
	pt := parsed.(TimeData)
	if pt.Value.Hour != 12 || pt.Value.Minute != 30 || pt.Value.Second != 45 {
		t.Errorf("got %v", pt.Value)
	}
}

func TestEncodeDecodeVariableLength(t *testing.T) {
	tests := []uint32{0, 1, 127, 128, 255, 256, 65535, 65536, 16777215, 16777216, 0xFFFFFFFF}
	for _, v := range tests {
		encoded := EncodeVariableLength(v)
		decoded, n := decodeVariableLength(encoded)
		if decoded != v {
			t.Errorf("encode/decode %d: got %d, consumed %d", v, decoded, n)
		}
	}
}

func TestCosemDateTime_Roundtrip(t *testing.T) {
	dt := CosemDateTime{
		Year: 2024, Month: 6, Day: 15,
		Hour: 10, Minute: 30, Second: 45,
		Hundredths: 50, Deviation: 120, ClockStatus: 0x00,
	}
	b := dt.ToBytes()
	if len(b) != 12 {
		t.Errorf("len=%d", len(b))
	}
	parsed, err := CosemDateTimeFromBytes(b)
	if err != nil {
		t.Fatal(err)
	}
	if parsed != dt {
		t.Errorf("got %v, want %v", parsed, dt)
	}
}

func TestCosemDateTime_Invalid(t *testing.T) {
	dt := CosemDateTime{Year: 0xFFFF, Month: 0xFF, Day: 0xFF}
	if !dt.IsInvalid() {
		t.Error("expected invalid")
	}
}

func TestCosemDateTime_ToTime(t *testing.T) {
	dt := CosemDateTime{Year: 2024, Month: 1, Day: 1, Hour: 0, Deviation: 0}
	tm := dt.ToTime()
	if tm.Year() != 2024 {
		t.Errorf("got %v", tm)
	}
}

func TestCosemDate_Roundtrip(t *testing.T) {
	d := CosemDate{Year: 2024, Month: 12, Day: 25, DayOfWeek: 3}
	b := d.ToBytes()
	if len(b) != 5 {
		t.Errorf("len=%d", len(b))
	}
	parsed := CosemDateFromBytes(b)
	if parsed != d {
		t.Errorf("got %v, want %v", parsed, d)
	}
}

func TestCosemDate_Invalid(t *testing.T) {
	d := CosemDate{Year: 0xFFFF, Month: 0xFF, Day: 0xFF}
	if !d.IsInvalid() {
		t.Error("expected invalid")
	}
}

func TestCosemTime_Roundtrip(t *testing.T) {
	tm := CosemTime{Hour: 23, Minute: 59, Second: 59, Hundredths: 99}
	b := tm.ToBytes()
	if len(b) != 4 {
		t.Errorf("len=%d", len(b))
	}
	parsed := CosemTimeFromBytes(b)
	if parsed != tm {
		t.Errorf("got %v, want %v", parsed, tm)
	}
}

func TestCosemTime_Invalid(t *testing.T) {
	tm := CosemTime{Hour: 0xFF, Minute: 0xFF, Second: 0xFF}
	if !tm.IsInvalid() {
		t.Error("expected invalid")
	}
}

func TestCosemTime_ToDuration(t *testing.T) {
	tm := CosemTime{Hour: 1, Minute: 30, Second: 0}
	d := tm.ToTimeDuration()
	if d != 90*time.Minute {
		t.Errorf("got %v", d)
	}
}

func TestCosemDateTimeFromBytes_Short(t *testing.T) {
	_, err := CosemDateTimeFromBytes([]byte{1, 2, 3})
	if err == nil {
		t.Error("expected error")
	}
}

func TestAccessSelector_String(t *testing.T) {
	tests := []struct {
		s    AccessSelector
		want string
	}{
		{AccessSelectorNoSelector, "no_selector"},
		{AccessSelectorRangeDescriptor, "range_descriptor"},
		{AccessSelectorEntryDescriptor, "entry_descriptor"},
		{AccessSelectorProfileGenericBuffer, "profile_generic_buffer"},
		{AccessSelectorSelectiveAccess, "selective_access"},
		{AccessSelector(99), "unknown(99)"},
	}
	for _, tc := range tests {
		if got := tc.s.String(); got != tc.want {
			t.Errorf("AccessSelector(%d).String() = %q, want %q", tc.s, got, tc.want)
		}
	}
}

func TestClassIDName(t *testing.T) {
	tests := []struct {
		id   uint16
		want string
	}{
		{ClassIDData, "Data"},
		{ClassIDRegister, "Register"},
		{ClassIDClock, "Clock"},
		{ClassIDSecuritySetup, "SecuritySetup"},
		{ClassIDProfileGeneric, "ProfileGeneric"},
		{9999, "Unknown(0x270f)"},
	}
	for _, tc := range tests {
		if got := ClassIDName(tc.id); got != tc.want {
			t.Errorf("ClassIDName(%d) = %q, want %q", tc.id, got, tc.want)
		}
	}
}

func TestNestedArrayStructure(t *testing.T) {
	// Array of structures
	inner := StructureData{IntegerData(1), VisibleStringData("a")}
	arr := ArrayData{inner, inner}
	b := arr.ToBytes()
	parsed, _, err := DlmsDataFromBytes(b)
	if err != nil {
		t.Fatal(err)
	}
	parr := parsed.(ArrayData)
	if len(parr) != 2 {
		t.Fatalf("len=%d", len(parr))
	}
}

func TestEmptyArray(t *testing.T) {
	a := ArrayData{}
	b := a.ToBytes()
	if !bytes.Equal(b, []byte{1, 0}) {
		t.Errorf("got %v", b)
	}
	parsed, _, err := DlmsDataFromBytes(b)
	if err != nil {
		t.Fatal(err)
	}
	if len(parsed.(ArrayData)) != 0 {
		t.Error("expected empty")
	}
}

func TestEmptyStructure(t *testing.T) {
	s := StructureData{}
	b := s.ToBytes()
	if !bytes.Equal(b, []byte{2, 0}) {
		t.Errorf("got %v", b)
	}
}

func TestDoubleLongZero(t *testing.T) {
	d := DoubleLongData(0)
	b := d.ToBytes()
	if b[1] != 0 || b[2] != 0 || b[3] != 0 || b[4] != 0 {
		t.Errorf("got %v", b)
	}
}

func TestDoubleLongMax(t *testing.T) {
	d := DoubleLongData(2147483647)
	b := d.ToBytes()
	v := int32(binary.BigEndian.Uint32(b[1:5]))
	if v != 2147483647 {
		t.Errorf("got %d", v)
	}
}

func TestDoubleLongMin(t *testing.T) {
	d := DoubleLongData(-2147483648)
	b := d.ToBytes()
	v := int32(binary.BigEndian.Uint32(b[1:5]))
	if v != -2147483648 {
		t.Errorf("got %d", v)
	}
}

func TestToPython(t *testing.T) {
	var nd NullData
	if nd.ToPython() != nil {
		t.Error("null")
	}
	if BooleanData(true).ToPython() != true {
		t.Error("bool")
	}
	if DoubleLongData(42).ToPython() != int32(42) {
		t.Error("doublelong")
	}
	if VisibleStringData("abc").ToPython() != "abc" {
		t.Error("visible string")
	}
	arr := ArrayData{IntegerData(1), IntegerData(2)}
	py := arr.ToPython().([]interface{})
	if len(py) != 2 {
		t.Error("array")
	}
}

func TestNewCosemDateTime(t *testing.T) {
	now := time.Date(2024, 6, 15, 10, 30, 0, 0, time.UTC)
	dt := NewCosemDateTime(now)
	if dt.Year != 2024 || dt.Month != 6 || dt.Day != 15 {
		t.Errorf("got %v", dt)
	}
}

func TestMustParseObis(t *testing.T) {
	oc := MustParseObis("1.0.1.8.0.255")
	if oc.A() != 1 {
		t.Error("wrong A")
	}
}

func TestMustParseObis_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic")
		}
	}()
	MustParseObis("invalid")
}

func TestCosemDateTimeToTime_Invalid(t *testing.T) {
	dt := CosemDateTime{Year: 0xFFFF, Month: 0xFF, Day: 0xFF}
	if !dt.ToTime().IsZero() {
		t.Error("expected zero time")
	}
}

func TestCosemDateTimeToTime_WithDeviation(t *testing.T) {
	dt := CosemDateTime{Year: 2024, Month: 1, Day: 1, Deviation: 60}
	tm := dt.ToTime()
	_, offset := tm.Zone()
	if offset != 3600 {
		t.Errorf("offset=%d", offset)
	}
}

func TestCosemDateTimeToTime_NegativeDeviation(t *testing.T) {
	dt := CosemDateTime{Year: 2024, Month: 1, Day: 1, Deviation: -120}
	tm := dt.ToTime()
	_, offset := tm.Zone()
	if offset != -7200 {
		t.Errorf("offset=%d", offset)
	}
}

func TestCosemDateTimeToTime_UndefinedDeviation(t *testing.T) {
	dt := CosemDateTime{Year: 2024, Month: 1, Day: 1, Deviation: 0x8000}
	tm := dt.ToTime()
	if tm.Year() != 2024 {
		t.Errorf("got %v", tm)
	}
}

func TestObisBytes(t *testing.T) {
	oc := MustParseObis("0.1.2.3.4.5")
	b := oc.Bytes()
	if len(b) != 6 {
		t.Errorf("len=%d", len(b))
	}
	if b[0] != 0 || b[5] != 5 {
		t.Errorf("got %v", b)
	}
}

package cosem

import (
	"bytes"
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestData_ClassID(t *testing.T) {
	d := &Data{LogicalName: core.ObisData}
	if d.ClassID() != core.ClassIDData {
		t.Error("wrong class ID")
	}
	if d.GetLogicalName().String() != core.ObisData.String() {
		t.Error("wrong logical name")
	}
}

func TestData_MarshalBinary(t *testing.T) {
	d := &Data{Value: core.IntegerData(42)}
	b, err := d.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(b, core.IntegerData(42).ToBytes()) {
		t.Errorf("got %v", b)
	}
}

func TestData_UnmarshalBinary(t *testing.T) {
	d := &Data{}
	err := d.UnmarshalBinary(core.DoubleLongData(12345).ToBytes())
	if err != nil {
		t.Fatal(err)
	}
	if d.Value.(core.DoubleLongData) != 12345 {
		t.Error("wrong value")
	}
}

func TestData_UnmarshalBinary_Empty(t *testing.T) {
	d := &Data{}
	err := d.UnmarshalBinary([]byte{})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := d.Value.(core.NullData); !ok {
		t.Error("expected null")
	}
}

func TestData_UnmarshalBinary_Invalid(t *testing.T) {
	d := &Data{}
	err := d.UnmarshalBinary([]byte{0x07}) // unknown tag
	if err == nil {
		t.Error("expected error")
	}
}

func TestData_NilValue(t *testing.T) {
	d := &Data{Value: nil}
	b, err := d.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	if b[0] != core.TagNull {
		t.Errorf("got %v", b)
	}
}

func TestRegister_ClassID(t *testing.T) {
	r := &Register{LogicalName: core.ObisActivePowerPlus}
	if r.ClassID() != core.ClassIDRegister {
		t.Error("wrong class ID")
	}
}

func TestRegister_MarshalBinary(t *testing.T) {
	r := &Register{Value: core.DoubleLongUnsignedData(12345)}
	b, err := r.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	parsed, _, _ := core.DlmsDataFromBytes(b)
	if parsed.(core.DoubleLongUnsignedData) != 12345 {
		t.Error("wrong value")
	}
}

func TestRegister_NilValue(t *testing.T) {
	r := &Register{Value: nil}
	b, err := r.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	if b[0] != core.TagDoubleLongUnsigned {
		t.Errorf("got %02x", b[0])
	}
}

func TestRegister_UnmarshalBinary(t *testing.T) {
	r := &Register{}
	err := r.UnmarshalBinary(core.DoubleLongUnsignedData(99999).ToBytes())
	if err != nil {
		t.Fatal(err)
	}
	if r.Value.(core.DoubleLongUnsignedData) != 99999 {
		t.Error("wrong value")
	}
}

func TestExtendedRegister_ClassID(t *testing.T) {
	e := &ExtendedRegister{}
	if e.ClassID() != core.ClassIDExtendedRegister {
		t.Error("wrong class ID")
	}
}

func TestExtendedRegister_MarshalUnmarshal(t *testing.T) {
	e := &ExtendedRegister{Value: core.Long64Data(-500)}
	b, err := e.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	var e2 ExtendedRegister
	err = e2.UnmarshalBinary(b)
	if err != nil {
		t.Fatal(err)
	}
	if e2.Value.(core.Long64Data) != -500 {
		t.Error("roundtrip failed")
	}
}

func TestDemandRegister_ClassID(t *testing.T) {
	d := &DemandRegister{}
	if d.ClassID() != core.ClassIDDemandRegister {
		t.Error("wrong class ID")
	}
}

func TestDemandRegister_MarshalUnmarshal(t *testing.T) {
	d := &DemandRegister{CurrentValue: core.DoubleLongData(100)}
	b, err := d.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	var d2 DemandRegister
	err = d2.UnmarshalBinary(b)
	if err != nil {
		t.Fatal(err)
	}
	if d2.CurrentValue.(core.DoubleLongData) != 100 {
		t.Error("roundtrip failed")
	}
}

func TestClock_ClassID(t *testing.T) {
	c := &Clock{LogicalName: core.ObisClock}
	if c.ClassID() != core.ClassIDClock {
		t.Error("wrong class ID")
	}
}

func TestClock_MarshalUnmarshal(t *testing.T) {
	dt := core.NewCosemDateTime(testTime())
	c := &Clock{Time: dt}
	b, err := c.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	var c2 Clock
	err = c2.UnmarshalBinary(b)
	if err != nil {
		t.Fatal(err)
	}
	if c2.Time.Year != dt.Year {
		t.Error("roundtrip failed")
	}
}

func TestClock_UnmarshalBinary_Invalid(t *testing.T) {
	c := &Clock{}
	err := c.UnmarshalBinary(core.IntegerData(42).ToBytes())
	if err == nil {
		t.Error("expected error for non-datetime")
	}
}

func TestProfileGeneric_ClassID(t *testing.T) {
	p := &ProfileGeneric{}
	if p.ClassID() != core.ClassIDProfileGeneric {
		t.Error("wrong class ID")
	}
}

func TestProfileGeneric_MarshalUnmarshal(t *testing.T) {
	p := &ProfileGeneric{
		Buffer: []core.StructureData{
			{core.IntegerData(1), core.DoubleLongUnsignedData(100)},
			{core.IntegerData(2), core.DoubleLongUnsignedData(200)},
		},
	}
	b, err := p.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	var p2 ProfileGeneric
	err = p2.UnmarshalBinary(b)
	if err != nil {
		t.Fatal(err)
	}
	if len(p2.Buffer) != 2 {
		t.Fatalf("len=%d", len(p2.Buffer))
	}
}

func TestProfileGeneric_UnmarshalBinary_Invalid(t *testing.T) {
	p := &ProfileGeneric{}
	err := p.UnmarshalBinary(core.IntegerData(42).ToBytes())
	if err == nil {
		t.Error("expected error")
	}
}

func TestProfileGeneric_UnmarshalBinary_NonStructure(t *testing.T) {
	p := &ProfileGeneric{}
	arr := core.ArrayData{core.IntegerData(1), core.IntegerData(2)} // not structures
	err := p.UnmarshalBinary(arr.ToBytes())
	if err == nil {
		t.Error("expected error for non-structure elements")
	}
}

func TestSecuritySetup_ClassID(t *testing.T) {
	s := &SecuritySetup{}
	if s.ClassID() != core.ClassIDSecuritySetup {
		t.Error("wrong class ID")
	}
}

func TestSecuritySetup_MarshalUnmarshal(t *testing.T) {
	s := &SecuritySetup{
		SecuritySuite:     1,
		EncryptionKey:     []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10},
		AuthenticationKey: []byte{0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99},
	}
	b, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	var s2 SecuritySetup
	err = s2.UnmarshalBinary(b)
	if err != nil {
		t.Fatal(err)
	}
	if s2.SecuritySuite != 1 {
		t.Errorf("suite=%d", s2.SecuritySuite)
	}
	if !bytes.Equal(s2.EncryptionKey, s.EncryptionKey) {
		t.Error("encryption key mismatch")
	}
}

func TestSecuritySetup_UnmarshalBinary_Short(t *testing.T) {
	s := &SecuritySetup{}
	err := s.UnmarshalBinary(core.IntegerData(42).ToBytes())
	if err == nil {
		t.Error("expected error")
	}
}

func TestTariffPlan_ClassID(t *testing.T) {
	tp := &TariffPlan{}
	if tp.ClassID() != core.ClassIDTariffSchedule {
		t.Error("wrong class ID")
	}
}

func TestTariffPlan_MarshalUnmarshal(t *testing.T) {
	tp := &TariffPlan{PlanName: "TOU"}
	b, err := tp.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	var tp2 TariffPlan
	err = tp2.UnmarshalBinary(b)
	if err != nil {
		t.Fatal(err)
	}
	if tp2.PlanName != "TOU" {
		t.Errorf("got %s", tp2.PlanName)
	}
}

func TestRS485Setup_MarshalUnmarshal(t *testing.T) {
	r := &RS485Setup{DefaultBaudRate: 9600, CommSpeed: 1}
	b, err := r.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	var r2 RS485Setup
	err = r2.UnmarshalBinary(b)
	if err != nil {
		t.Fatal(err)
	}
	if r2.DefaultBaudRate != 9600 || r2.CommSpeed != 1 {
		t.Errorf("got %d, %d", r2.DefaultBaudRate, r2.CommSpeed)
	}
}

func TestRS485Setup_UnmarshalBinary_Short(t *testing.T) {
	r := &RS485Setup{}
	err := r.UnmarshalBinary(core.IntegerData(42).ToBytes())
	if err == nil {
		t.Error("expected error")
	}
}

func TestInfraredSetup_MarshalUnmarshal(t *testing.T) {
	i := &InfraredSetup{DefaultBaudRate: 2400, CommSpeed: 0}
	b, err := i.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	var i2 InfraredSetup
	err = i2.UnmarshalBinary(b)
	if err != nil {
		t.Fatal(err)
	}
	if i2.DefaultBaudRate != 2400 {
		t.Errorf("got %d", i2.DefaultBaudRate)
	}
}

func TestInfraredSetup_UnmarshalBinary_Short(t *testing.T) {
	i := &InfraredSetup{}
	err := i.UnmarshalBinary(core.IntegerData(42).ToBytes())
	if err == nil {
		t.Error("expected error")
	}
}

func TestNBIoTSetup_MarshalUnmarshal(t *testing.T) {
	n := &NBIoTSetup{APN: "internet.mnc000.mcc460.gprs"}
	b, err := n.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	var n2 NBIoTSetup
	err = n2.UnmarshalBinary(b)
	if err != nil {
		t.Fatal(err)
	}
	if n2.APN != "internet.mnc000.mcc460.gprs" {
		t.Errorf("got %s", n2.APN)
	}
}

func TestLoRaWANSetup_MarshalUnmarshal(t *testing.T) {
	l := &LoRaWANSetup{
		DevEUI: []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
		AppEUI: []byte{0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x00, 0x11},
		AppKey: []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10},
	}
	b, err := l.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	var l2 LoRaWANSetup
	err = l2.UnmarshalBinary(b)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(l2.DevEUI, l.DevEUI) {
		t.Error("DevEUI mismatch")
	}
	if !bytes.Equal(l2.AppKey, l.AppKey) {
		t.Error("AppKey mismatch")
	}
}

func TestLoRaWANSetup_UnmarshalBinary_Short(t *testing.T) {
	l := &LoRaWANSetup{}
	err := l.UnmarshalBinary(core.IntegerData(42).ToBytes())
	if err == nil {
		t.Error("expected error")
	}
}

func TestAccess_Default(t *testing.T) {
	d := &Data{}
	a := d.Access(99)
	if a.Access != core.AccessNone {
		t.Error("unknown attr should be none")
	}
}

func TestLPSetup_ClassID(t *testing.T) {
	lp := &LPSetup{}
	if lp.ClassID() != core.ClassIDLpSetup {
		t.Error("wrong class ID")
	}
}

func TestSeasonProfile(t *testing.T) {
	sp := SeasonProfile{
		SeasonName:     "Summer",
		WeekProfileRef: "WP1",
	}
	if sp.SeasonName != "Summer" {
		t.Error("wrong name")
	}
}

func TestDaySchedule(t *testing.T) {
	ds := DaySchedule{
		StartTime: core.CosemTime{Hour: 8, Minute: 0},
		ScriptRef: "Script1",
	}
	if ds.ScriptRef != "Script1" {
		t.Error("wrong ref")
	}
}

func TestDayProfile(t *testing.T) {
	dp := DayProfile{
		DayProfileName: "Weekday",
		Schedule: []DaySchedule{
			{core.CosemTime{Hour: 8}, "S1"},
			{core.CosemTime{Hour: 17}, "S2"},
		},
	}
	if len(dp.Schedule) != 2 {
		t.Error("wrong schedule count")
	}
}

func TestWeekProfile(t *testing.T) {
	wp := WeekProfile{
		ProfileName: "Default",
		Monday:      "DP1",
		Tuesday:     "DP1",
	}
	if wp.Monday != "DP1" {
		t.Error("wrong day ref")
	}
}

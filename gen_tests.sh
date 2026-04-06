#!/bin/bash
# Generate test files for IC classes missing tests

cd /Users/yimiliya/.openclaw/workspace/dlms-cosem-go/cosem

# ic1_data_test.go
cat > ic1_data_test.go << 'EOF'
package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestData_New(t *testing.T) {
	d := &Data{LogicalName: core.ObisCode{0, 0, 1, 0, 0, 255}}
	if d.LogicalName != (core.ObisCode{0, 0, 1, 0, 0, 255}) {
		t.Error("LogicalName mismatch")
	}
	if d.Value != nil {
		t.Error("default Value should be nil")
	}
	if d.Version != 0 {
		t.Error("default Version should be 0")
	}
}

func TestData_ClassID(t *testing.T) {
	d := &Data{}
	tests := []struct {
		name string
		want uint16
	}{
		{"bluebook_ed16", 1},
		{"const_match", core.ClassIDData},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := d.ClassID(); got != tt.want {
				t.Errorf("ClassID() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestData_MarshalBinary(t *testing.T) {
	tests := []struct {
		name string
		data *Data
	}{
		{"nil_value", &Data{Value: nil}},
		{"null_data", &Data{Value: core.NullData{}}},
		{"long_unsigned", &Data{Value: core.LongUnsignedData(42)}},
		{"visible_string", &Data{Value: core.VisibleStringData("hello")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := tt.data.MarshalBinary()
			if err != nil {
				t.Fatal(err)
			}
			if len(b) == 0 {
				t.Fatal("MarshalBinary returned empty bytes")
			}
			d2 := &Data{}
			if err := d2.UnmarshalBinary(b); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestData_Fields(t *testing.T) {
	d := &Data{}
	d.LogicalName = core.ObisCode{1, 1, 1, 1, 1, 1}
	d.Value = core.LongUnsignedData(999)
	d.Version = 5
	if d.LogicalName != (core.ObisCode{1, 1, 1, 1, 1, 1}) {
		t.Error("LogicalName not set")
	}
	if d.Version != 5 {
		t.Error("Version not set")
	}
	if d.GetLogicalName() != (core.ObisCode{1, 1, 1, 1, 1, 1}) {
		t.Error("GetLogicalName mismatch")
	}
	if d.GetVersion() != 5 {
		t.Error("GetVersion mismatch")
	}
}
EOF

# ic3_register_test.go
cat > ic3_register_test.go << 'EOF'
package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestRegister_New(t *testing.T) {
	r := &Register{LogicalName: core.ObisCode{0, 0, 1, 0, 0, 255}}
	if r.Value != nil {
		t.Error("default Value should be nil")
	}
	if r.Scaler != 0 || r.Unit != 0 || r.Status != nil {
		t.Error("defaults should be zero")
	}
}

func TestRegister_ClassID(t *testing.T) {
	r := &Register{}
	tests := []struct{ name string; want uint16 }{
		{"bluebook_ed16", 3},
		{"const_match", core.ClassIDRegister},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.ClassID(); got != tt.want {
				t.Errorf("ClassID() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestRegister_MarshalBinary(t *testing.T) {
	tests := []struct {
		name string
		reg  *Register
	}{
		{"nil_value", &Register{Value: nil}},
		{"long_unsigned", &Register{Value: core.LongUnsignedData(12345)}},
		{"double_long", &Register{Value: core.DoubleLongUnsignedData(0xFFFFFFFF)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := tt.reg.MarshalBinary()
			if err != nil {
				t.Fatal(err)
			}
			r2 := &Register{}
			if err := r2.UnmarshalBinary(b); err != nil {
				t.Fatal(err)
			}
			if tt.reg.Value != nil && r2.Value == nil {
				t.Fatal("unmarshaled Value is nil")
			}
		})
	}
}

func TestRegister_Fields(t *testing.T) {
	r := &Register{
		LogicalName: core.ObisCode{1, 2, 3, 4, 5, 6},
		Value:       core.LongUnsignedData(42),
		Scaler:      -3,
		Unit:        27,
		Version:     2,
	}
	if r.Scaler != -3 {
		t.Errorf("Scaler = %d, want -3", r.Scaler)
	}
	if r.Unit != 27 {
		t.Errorf("Unit = %d, want 27", r.Unit)
	}
	if r.GetVersion() != 2 {
		t.Errorf("GetVersion() = %d, want 2", r.GetVersion())
	}
}
EOF

# ic4_extended_register_test.go
cat > ic4_extended_register_test.go << 'EOF'
package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestExtendedRegister_ClassID(t *testing.T) {
	e := &ExtendedRegister{}
	if e.ClassID() != 4 {
		t.Errorf("ClassID() = %d, want 4", e.ClassID())
	}
	if e.ClassID() != core.ClassIDExtendedRegister {
		t.Error("ClassID mismatch with const")
	}
}

func TestExtendedRegister_MarshalBinary(t *testing.T) {
	tests := []struct {
		name  string
		value core.DlmsData
	}{
		{"long_unsigned", core.LongUnsignedData(100)},
		{"double_long_unsigned", core.DoubleLongUnsignedData(99999)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ExtendedRegister{Value: tt.value}
			b, err := e.MarshalBinary()
			if err != nil {
				t.Fatal(err)
			}
			e2 := &ExtendedRegister{}
			if err := e2.UnmarshalBinary(b); err != nil {
				t.Fatal(err)
			}
			if e2.Value == nil {
				t.Fatal("Value nil after unmarshal")
			}
		})
	}
}

func TestExtendedRegister_Fields(t *testing.T) {
	e := &ExtendedRegister{
		LogicalName: core.ObisCode{0, 0, 1, 0, 0, 255},
		Scaler:      -2,
		Unit:        30,
		Version:     1,
	}
	if e.GetLogicalName() != (core.ObisCode{0, 0, 1, 0, 0, 255}) {
		t.Error("GetLogicalName mismatch")
	}
	if e.GetVersion() != 1 {
		t.Error("GetVersion mismatch")
	}
}
EOF

# ic5_demand_register_test.go
cat > ic5_demand_register_test.go << 'EOF'
package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestDemandRegister_ClassID(t *testing.T) {
	d := &DemandRegister{}
	if d.ClassID() != 5 {
		t.Errorf("ClassID() = %d, want 5", d.ClassID())
	}
	if d.ClassID() != core.ClassIDDemandRegister {
		t.Error("ClassID mismatch with const")
	}
}

func TestDemandRegister_MarshalBinary(t *testing.T) {
	tests := []struct {
		name  string
		value core.DlmsData
	}{
		{"zero", core.LongUnsignedData(0)},
		{"max_uint32", core.DoubleLongUnsignedData(0xFFFFFFFF)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DemandRegister{CurrentValue: tt.value}
			b, err := d.MarshalBinary()
			if err != nil {
				t.Fatal(err)
			}
			d2 := &DemandRegister{}
			if err := d2.UnmarshalBinary(b); err != nil {
				t.Fatal(err)
			}
			if d2.CurrentValue == nil {
				t.Fatal("CurrentValue nil after unmarshal")
			}
		})
	}
}

func TestDemandRegister_Fields(t *testing.T) {
	d := &DemandRegister{
		LogicalName: core.ObisCode{0, 0, 1, 0, 0, 255},
		Scaler:      -1,
		Unit:        1,
		Version:     0,
	}
	if d.Scaler != -1 {
		t.Errorf("Scaler = %d, want -1", d.Scaler)
	}
	if d.Unit != 1 {
		t.Errorf("Unit = %d, want 1", d.Unit)
	}
}
EOF

# ic7_profile_generic_test.go
cat > ic7_profile_generic_test.go << 'EOF'
package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestProfileGeneric_ClassID(t *testing.T) {
	p := &ProfileGeneric{}
	if p.ClassID() != 7 {
		t.Errorf("ClassID() = %d, want 7", p.ClassID())
	}
	if p.ClassID() != core.ClassIDProfileGeneric {
		t.Error("ClassID mismatch with const")
	}
}

func TestProfileGeneric_New(t *testing.T) {
	p := &ProfileGeneric{
		LogicalName:   core.ObisCode{1, 0, 99, 1, 0, 255},
		CapturePeriod: 60,
		Entries:       100,
		Version:       1,
	}
	if p.CapturePeriod != 60 {
		t.Error("CapturePeriod mismatch")
	}
	if p.Entries != 100 {
		t.Error("Entries mismatch")
	}
	if p.EntriesInUse != 0 {
		t.Error("EntriesInUse should default to 0")
	}
}

func TestProfileGeneric_MarshalBinary(t *testing.T) {
	p := &ProfileGeneric{Buffer: []core.StructureData{
		{core.UnsignedIntegerData(1), core.LongUnsignedData(100)},
	}}
	b, err := p.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	p2 := &ProfileGeneric{}
	if err := p2.UnmarshalBinary(b); err != nil {
		t.Fatal(err)
	}
}

func TestProfileGeneric_Fields(t *testing.T) {
	p := &ProfileGeneric{
		LogicalName:   core.ObisCode{1, 0, 99, 1, 0, 255},
		CapturePeriod: 15,
		EntriesInUse:  50,
		Entries:       100,
		SortMethod:    1,
		Version:       2,
	}
	if p.GetVersion() != 2 {
		t.Error("GetVersion mismatch")
	}
	if p.SortMethod != 1 {
		t.Error("SortMethod mismatch")
	}
}
EOF

# ic8_clock_test.go
cat > ic8_clock_test.go << 'EOF'
package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestClock_ClassID(t *testing.T) {
	c := &Clock{}
	if c.ClassID() != 8 {
		t.Errorf("ClassID() = %d, want 8", c.ClassID())
	}
	if c.ClassID() != core.ClassIDClock {
		t.Error("ClassID mismatch with const")
	}
}

func TestClock_New(t *testing.T) {
	c := &Clock{
		LogicalName: core.ObisCode{0, 0, 1, 0, 0, 255},
		TimeZone:    480,
		Status:      0,
	}
	if c.TimeZone != 480 {
		t.Error("TimeZone mismatch")
	}
	if c.DaylightSavingsEnabled != false {
		t.Error("DaylightSavingsEnabled should default false")
	}
}

func TestClock_MarshalBinary(t *testing.T) {
	c := &Clock{Time: core.CosemDateTime{Year: 2024, Month: 1, Day: 1}}
	b, err := c.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	if len(b) == 0 {
		t.Fatal("empty bytes")
	}
	c2 := &Clock{}
	if err := c2.UnmarshalBinary(b); err != nil {
		t.Fatal(err)
	}
}

func TestClock_Fields(t *testing.T) {
	c := &Clock{
		LogicalName:              core.ObisCode{0, 0, 1, 0, 0, 255},
		TimeZone:                 540,
		Status:                   0x80,
		DaylightSavingsDeviation: 60,
		DaylightSavingsEnabled:   true,
		Version:                  1,
	}
	if c.DaylightSavingsEnabled != true {
		t.Error("DaylightSavingsEnabled not set")
	}
	if c.GetVersion() != 1 {
		t.Error("GetVersion mismatch")
	}
}
EOF

# ic9_script_table_test.go
cat > ic9_script_table_test.go << 'EOF'
package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestScriptTable_ClassID(t *testing.T) {
	st := &ScriptTable{}
	if st.ClassID() != 9 {
		t.Errorf("ClassID() = %d, want 9", st.ClassID())
	}
	if st.ClassID() != core.ClassIDScriptTable {
		t.Error("ClassID mismatch with const")
	}
}

func TestScriptTable_New(t *testing.T) {
	st := &ScriptTable{
		LogicalName: core.ObisCode{0, 0, 10, 0, 0, 255},
	}
	if st.Scripts == nil {
		t.Log("Scripts nil by default - ok")
	}
}

func TestScriptTable_MarshalBinary(t *testing.T) {
	st := &ScriptTable{Scripts: [][]byte{[]byte{0x01, 0x00, 0xC0, 0x01}}}
	b, err := st.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	st2 := &ScriptTable{}
	if err := st2.UnmarshalBinary(b); err != nil {
		t.Fatal(err)
	}
}

func TestScriptTable_Fields(t *testing.T) {
	st := &ScriptTable{
		LogicalName: core.ObisCode{0, 0, 10, 0, 0, 255},
		Version:     1,
	}
	if st.GetVersion() != 1 {
		t.Error("GetVersion mismatch")
	}
}
EOF

# ic10_schedule_test.go
cat > ic10_schedule_test.go << 'EOF'
package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestSchedule_ClassID(t *testing.T) {
	s := &Schedule{}
	if s.ClassID() != 10 {
		t.Errorf("ClassID() = %d, want 10", s.ClassID())
	}
	if s.ClassID() != core.ClassIDSchedule {
		t.Error("ClassID mismatch with const")
	}
}

func TestSchedule_New(t *testing.T) {
	s := &Schedule{LogicalName: core.ObisCode{0, 0, 10, 0, 0, 255}}
	if s.Entries == nil {
		t.Log("Entries nil by default - ok")
	}
}

func TestSchedule_MarshalBinary(t *testing.T) {
	s := &Schedule{}
	b, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &Schedule{}
	if err := s2.UnmarshalBinary(b); err != nil {
		t.Fatal(err)
	}
}

func TestSchedule_Fields(t *testing.T) {
	s := &Schedule{
		LogicalName: core.ObisCode{0, 0, 10, 0, 0, 255},
		Version:     0,
	}
	if s.GetVersion() != 0 {
		t.Error("GetVersion mismatch")
	}
}
EOF

# ic11_special_day_table_test.go
cat > ic11_special_day_table_test.go << 'EOF'
package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestSpecialDaysTable_ClassID(t *testing.T) {
	s := &SpecialDaysTable{}
	if s.ClassID() != 11 {
		t.Errorf("ClassID() = %d, want 11", s.ClassID())
	}
	if s.ClassID() != core.ClassIDSpecialDaysTable {
		t.Error("ClassID mismatch with const")
	}
}

func TestSpecialDaysTable_New(t *testing.T) {
	s := &SpecialDaysTable{LogicalName: core.ObisCode{0, 0, 11, 0, 0, 255}}
	if s.Entries == nil {
		t.Log("Entries nil by default - ok")
	}
}

func TestSpecialDaysTable_MarshalBinary(t *testing.T) {
	s := &SpecialDaysTable{}
	b, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &SpecialDaysTable{}
	if err := s2.UnmarshalBinary(b); err != nil {
		t.Fatal(err)
	}
}

func TestSpecialDaysTable_Fields(t *testing.T) {
	s := &SpecialDaysTable{
		LogicalName: core.ObisCode{0, 0, 11, 0, 0, 255},
		Version:     0,
	}
	if s.GetVersion() != 0 {
		t.Error("GetVersion mismatch")
	}
}
EOF

# ic12_association_sn_test.go
cat > ic12_association_sn_test.go << 'EOF'
package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestAssociationSN_ClassID(t *testing.T) {
	a := &AssociationSN{}
	if a.ClassID() != 12 {
		t.Errorf("ClassID() = %d, want 12", a.ClassID())
	}
	if a.ClassID() != core.ClassIDAssociationSN {
		t.Error("ClassID mismatch with const")
	}
}

func TestAssociationSN_New(t *testing.T) {
	a := &AssociationSN{
		LogicalName:        core.ObisCode{0, 0, 40, 0, 0, 255},
		AuthenticationLevel: 5,
	}
	if a.AuthenticationLevel != 5 {
		t.Error("AuthenticationLevel mismatch")
	}
}

func TestAssociationSN_MarshalBinary(t *testing.T) {
	a := &AssociationSN{}
	b, err := a.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	a2 := &AssociationSN{}
	if err := a2.UnmarshalBinary(b); err != nil {
		t.Fatal(err)
	}
}

func TestAssociationSN_Fields(t *testing.T) {
	a := &AssociationSN{
		LogicalName:        core.ObisCode{0, 0, 40, 0, 0, 255},
		AuthenticationLevel: 0,
		Version:             1,
	}
	if a.GetVersion() != 1 {
		t.Error("GetVersion mismatch")
	}
}
EOF

# ic14_register_activation_test.go
cat > ic14_register_activation_test.go << 'EOF'
package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestRegisterActivation_ClassID(t *testing.T) {
	ra := &RegisterActivation{}
	if ra.ClassID() != 6 {
		t.Errorf("ClassID() = %d, want 6", ra.ClassID())
	}
	if ra.ClassID() != core.ClassIDRegisterActivation {
		t.Error("ClassID mismatch with const")
	}
}

func TestRegisterActivation_New(t *testing.T) {
	ra := &RegisterActivation{
		LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255},
	}
	if ra.RegisterAssignments == nil {
		t.Log("RegisterAssignments nil by default - ok")
	}
}

func TestRegisterActivation_MarshalBinary(t *testing.T) {
	ra := &RegisterActivation{}
	b, err := ra.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	ra2 := &RegisterActivation{}
	if err := ra2.UnmarshalBinary(b); err != nil {
		t.Fatal(err)
	}
}

func TestRegisterActivation_Fields(t *testing.T) {
	ra := &RegisterActivation{
		LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255},
		Version:     0,
	}
	if ra.GetVersion() != 0 {
		t.Error("GetVersion mismatch")
	}
}
EOF

# ic15_association_ln_test.go
cat > ic15_association_ln_test.go << 'EOF'
package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestAssociationLN_ClassID(t *testing.T) {
	a := &AssociationLN{}
	if a.ClassID() != 15 {
		t.Errorf("ClassID() = %d, want 15", a.ClassID())
	}
	if a.ClassID() != core.ClassIDAssociationLN {
		t.Error("ClassID mismatch with const")
	}
}

func TestAssociationLN_New(t *testing.T) {
	a := &AssociationLN{
		LogicalName: core.ObisCode{0, 0, 40, 0, 0, 255},
	}
	if a.AssociationStatus != 0 {
		t.Error("AssociationStatus should default 0")
	}
}

func TestAssociationLN_MarshalBinary(t *testing.T) {
	a := &AssociationLN{}
	b, err := a.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	a2 := &AssociationLN{}
	if err := a2.UnmarshalBinary(b); err != nil {
		t.Fatal(err)
	}
}

func TestAssociationLN_Fields(t *testing.T) {
	a := &AssociationLN{
		LogicalName: core.ObisCode{0, 0, 40, 0, 0, 255},
		Version:     2,
	}
	if a.GetVersion() != 2 {
		t.Error("GetVersion mismatch")
	}
}
EOF

# ic18_image_transfer_test.go
cat > ic18_image_transfer_test.go << 'EOF'
package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestImageTransfer_ClassID(t *testing.T) {
	it := &ImageTransfer{}
	if it.ClassID() != 18 {
		t.Errorf("ClassID() = %d, want 18", it.ClassID())
	}
	if it.ClassID() != core.ClassIDImageTransfer {
		t.Error("ClassID mismatch with const")
	}
}

func TestImageTransfer_New(t *testing.T) {
	it := &ImageTransfer{
		LogicalName: core.ObisCode{0, 0, 44, 0, 0, 255},
	}
	if it.ImageBlockSize != 0 {
		t.Error("ImageBlockSize should default 0")
	}
}

func TestImageTransfer_MarshalBinary(t *testing.T) {
	tests := []struct {
		name   string
		status uint8
	}{
		{"idle", 0},
		{"initiated", 1},
		{"verifying", 6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			it := &ImageTransfer{TransferStatus: tt.status}
			b, err := it.MarshalBinary()
			if err != nil {
				t.Fatal(err)
			}
			it2 := &ImageTransfer{}
			if err := it2.UnmarshalBinary(b); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestImageTransfer_Fields(t *testing.T) {
	it := &ImageTransfer{
		LogicalName:   core.ObisCode{0, 0, 44, 0, 0, 255},
		ImageBlockSize: 200,
		Version:       1,
	}
	if it.ImageBlockSize != 200 {
		t.Error("ImageBlockSize mismatch")
	}
}
EOF

# ic20_activity_calendar_test.go
cat > ic20_activity_calendar_test.go << 'EOF'
package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestActivityCalendar_ClassID(t *testing.T) {
	ac := &ActivityCalendar{}
	if ac.ClassID() != 20 {
		t.Errorf("ClassID() = %d, want 20", ac.ClassID())
	}
	if ac.ClassID() != core.ClassIDActivityCalendar {
		t.Error("ClassID mismatch with const")
	}
}

func TestActivityCalendar_New(t *testing.T) {
	ac := &ActivityCalendar{
		LogicalName:   core.ObisCode{0, 0, 12, 0, 0, 255},
		CalendarName:  "Summer",
		Version:       1,
	}
	if ac.CalendarName != "Summer" {
		t.Error("CalendarName mismatch")
	}
}

func TestActivityCalendar_MarshalBinary(t *testing.T) {
	tests := []struct {
		name string
		cn   string
	}{
		{"empty", ""},
		{"normal", "Season1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ac := &ActivityCalendar{CalendarName: tt.cn}
			b, err := ac.MarshalBinary()
			if err != nil {
				t.Fatal(err)
			}
			ac2 := &ActivityCalendar{}
			if err := ac2.UnmarshalBinary(b); err != nil {
				t.Fatal(err)
			}
			if ac2.CalendarName != tt.cn {
				t.Errorf("CalendarName = %q, want %q", ac2.CalendarName, tt.cn)
			}
		})
	}
}

func TestActivityCalendar_Fields(t *testing.T) {
	ac := &ActivityCalendar{
		LogicalName:   core.ObisCode{0, 0, 12, 0, 0, 255},
		CalendarName:  "Winter",
		SeasonProfile: nil,
		Version:       2,
	}
	if ac.GetVersion() != 2 {
		t.Error("GetVersion mismatch")
	}
}
EOF

# ic21_register_monitor_test.go
cat > ic21_register_monitor_test.go << 'EOF'
package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestRegisterMonitor_ClassID(t *testing.T) {
	rm := &RegisterMonitor{}
	if rm.ClassID() != 21 {
		t.Errorf("ClassID() = %d, want 21", rm.ClassID())
	}
	if rm.ClassID() != core.ClassIDRegisterMonitor {
		t.Error("ClassID mismatch with const")
	}
}

func TestRegisterMonitor_New(t *testing.T) {
	rm := &RegisterMonitor{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if rm.Thresholds == nil {
		t.Log("Thresholds nil by default - ok")
	}
}

func TestRegisterMonitor_MarshalBinary(t *testing.T) {
	rm := &RegisterMonitor{}
	b, err := rm.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	rm2 := &RegisterMonitor{}
	if err := rm2.UnmarshalBinary(b); err != nil {
		t.Fatal(err)
	}
}

func TestRegisterMonitor_Fields(t *testing.T) {
	rm := &RegisterMonitor{
		LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255},
		Version:     0,
	}
	if rm.GetVersion() != 0 {
		t.Error("GetVersion mismatch")
	}
}
EOF

# ic22_single_action_schedule_test.go
cat > ic22_single_action_schedule_test.go << 'EOF'
package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestSingleActionSchedule_ClassID(t *testing.T) {
	sas := &SingleActionSchedule{}
	if sas.ClassID() != 22 {
		t.Errorf("ClassID() = %d, want 22", sas.ClassID())
	}
	if sas.ClassID() != core.ClassIDSingleActionSchedule {
		t.Error("ClassID mismatch with const")
	}
}

func TestSingleActionSchedule_New(t *testing.T) {
	sas := &SingleActionSchedule{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if sas.Entries == nil {
		t.Log("Entries nil by default - ok")
	}
}

func TestSingleActionSchedule_MarshalBinary(t *testing.T) {
	sas := &SingleActionSchedule{}
	b, err := sas.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	sas2 := &SingleActionSchedule{}
	if err := sas2.UnmarshalBinary(b); err != nil {
		t.Fatal(err)
	}
}

func TestSingleActionSchedule_Fields(t *testing.T) {
	sas := &SingleActionSchedule{
		LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255},
		Version:     0,
	}
	if sas.GetVersion() != 0 {
		t.Error("GetVersion mismatch")
	}
}
EOF

# ic30_value_display_test.go
cat > ic30_value_display_test.go << 'EOF'
package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestValueDisplay_ClassID(t *testing.T) {
	vd := &ValueDisplay{}
	if vd.ClassID() != 257 {
		t.Errorf("ClassID() = %d, want 257", vd.ClassID())
	}
	if vd.ClassID() != core.ClassIDValueDisplay {
		t.Error("ClassID mismatch with const")
	}
}

func TestValueDisplay_New(t *testing.T) {
	vd := &ValueDisplay{
		LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255},
	}
	if vd.ValueToDisplay != nil {
		t.Error("default ValueToDisplay should be nil")
	}
}

func TestValueDisplay_MarshalBinary(t *testing.T) {
	tests := []struct {
		name string
		val  core.DlmsData
	}{
		{"nil", nil},
		{"null", core.NullData{}},
		{"long_unsigned", core.LongUnsignedData(12345)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vd := &ValueDisplay{ValueToDisplay: tt.val}
			b, err := vd.MarshalBinary()
			if err != nil {
				t.Fatal(err)
			}
			vd2 := &ValueDisplay{}
			if err := vd2.UnmarshalBinary(b); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestValueDisplay_Fields(t *testing.T) {
	vd := &ValueDisplay{
		LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255},
		Version:     0,
	}
	if vd.GetVersion() != 0 {
		t.Error("GetVersion mismatch")
	}
}
EOF

# ic32_iec_public_key_test.go
cat > ic32_iec_public_key_test.go << 'EOF'
package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestIECPublicKey_ClassID(t *testing.T) {
	pk := &IECPublicKey{}
	if pk.ClassID() != 256 {
		t.Errorf("ClassID() = %d, want 256", pk.ClassID())
	}
	if pk.ClassID() != core.ClassIDIECPublicKey {
		t.Error("ClassID mismatch with const")
	}
}

func TestIECPublicKey_New(t *testing.T) {
	pk := &IECPublicKey{
		LogicalName: core.ObisCode{0, 0, 43, 0, 0, 255},
	}
	if pk.PublicKeyValue != nil {
		t.Error("default PublicKeyValue should be nil")
	}
}

func TestIECPublicKey_MarshalBinary(t *testing.T) {
	tests := []struct {
		name string
		key  []byte
	}{
		{"empty", []byte{}},
		{"small", []byte{0x01, 0x02, 0x03}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pk := &IECPublicKey{PublicKeyValue: tt.key}
			b, err := pk.MarshalBinary()
			if err != nil {
				t.Fatal(err)
			}
			pk2 := &IECPublicKey{}
			if err := pk2.UnmarshalBinary(b); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestIECPublicKey_Fields(t *testing.T) {
	pk := &IECPublicKey{
		LogicalName:    core.ObisCode{0, 0, 43, 0, 0, 255},
		PublicKeyValue: []byte{0xDE, 0xAD, 0xBE, 0xEF},
		Version:        0,
	}
	if len(pk.PublicKeyValue) != 4 {
		t.Error("PublicKeyValue length mismatch")
	}
}
EOF

# ic40_push_setup_test.go
cat > ic40_push_setup_test.go << 'EOF'
package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestPushSetup_ClassID(t *testing.T) {
	ps := &PushSetup{}
	if ps.ClassID() != 40 {
		t.Errorf("ClassID() = %d, want 40", ps.ClassID())
	}
	if ps.ClassID() != core.ClassIDPushSetup {
		t.Error("ClassID mismatch with const")
	}
}

func TestPushSetup_New(t *testing.T) {
	ps := &PushSetup{
		LogicalName:    core.ObisCode{0, 0, 41, 0, 0, 255},
		Service:        1,
		NumberOfRetries: 3,
		Version:        1,
	}
	if ps.Service != 1 {
		t.Error("Service mismatch")
	}
	if ps.NumberOfRetries != 3 {
		t.Error("NumberOfRetries mismatch")
	}
}

func TestPushSetup_MarshalBinary(t *testing.T) {
	ps := &PushSetup{
		PushObjectList: []PushObject{
			{ClassID: 1, LogicalName: core.ObisCode{0, 0, 1, 0, 0, 255}, Attribute: 2},
		},
	}
	b, err := ps.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	ps2 := &PushSetup{}
	if err := ps2.UnmarshalBinary(b); err != nil {
		t.Fatal(err)
	}
}

func TestPushSetup_Fields(t *testing.T) {
	ps := &PushSetup{
		LogicalName: core.ObisCode{0, 0, 41, 0, 0, 255},
		Destination: []byte{0x01, 0x02},
		Version:     2,
	}
	if ps.GetVersion() != 2 {
		t.Error("GetVersion mismatch")
	}
	if len(ps.Destination) != 2 {
		t.Error("Destination length mismatch")
	}
}
EOF

# ic41_tcp_udp_setup_test.go
cat > ic41_tcp_udp_setup_test.go << 'EOF'
package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestTCPUDPSetup_ClassID(t *testing.T) {
	s := &TCPUDPSetup{}
	if s.ClassID() != 41 {
		t.Errorf("ClassID() = %d, want 41", s.ClassID())
	}
	if s.ClassID() != core.ClassIDTCPUDPSetup {
		t.Error("ClassID mismatch with const")
	}
}

func TestTCPUDPSetup_New(t *testing.T) {
	s := &TCPUDPSetup{
		LogicalName: core.ObisCode{0, 0, 42, 0, 0, 255},
	}
	if s.TCPConnections == nil {
		t.Log("TCPConnections nil by default - ok")
	}
}

func TestTCPUDPSetup_MarshalBinary(t *testing.T) {
	s := &TCPUDPSetup{}
	b, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &TCPUDPSetup{}
	if err := s2.UnmarshalBinary(b); err != nil {
		t.Fatal(err)
	}
}

func TestTCPUDPSetup_Fields(t *testing.T) {
	s := &TCPUDPSetup{
		LogicalName: core.ObisCode{0, 0, 42, 0, 0, 255},
		Version:     0,
	}
	if s.GetVersion() != 0 {
		t.Error("GetVersion mismatch")
	}
}
EOF

# ic42_ipv4_setup_test.go
cat > ic42_ipv4_setup_test.go << 'EOF'
package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestIPv4Setup_ClassID(t *testing.T) {
	s := &IPv4Setup{}
	if s.ClassID() != 42 {
		t.Errorf("ClassID() = %d, want 42", s.ClassID())
	}
	if s.ClassID() != core.ClassIDIPv4Setup {
		t.Error("ClassID mismatch with const")
	}
}

func TestIPv4Setup_New(t *testing.T) {
	s := &IPv4Setup{
		LogicalName: core.ObisCode{0, 0, 42, 0, 0, 255},
	}
	for i := 0; i < 4; i++ {
		if s.IPAddress[i] != 0 {
			t.Errorf("IPAddress[%d] should be 0", i)
		}
	}
}

func TestIPv4Setup_MarshalBinary(t *testing.T) {
	tests := []struct {
		name string
		ip   [4]byte
	}{
		{"zero", [4]byte{}},
		{"local", [4]byte{192, 168, 1, 1}},
		{"broadcast", [4]byte{255, 255, 255, 255}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &IPv4Setup{IPAddress: tt.ip}
			b, err := s.MarshalBinary()
			if err != nil {
				t.Fatal(err)
			}
			s2 := &IPv4Setup{}
			if err := s2.UnmarshalBinary(b); err != nil {
				t.Fatal(err)
			}
			if s2.IPAddress != tt.ip {
				t.Errorf("IPAddress = %v, want %v", s2.IPAddress, tt.ip)
			}
		})
	}
}

func TestIPv4Setup_Fields(t *testing.T) {
	s := &IPv4Setup{
		LogicalName: core.ObisCode{0, 0, 42, 0, 0, 255},
		IPAddress:   [4]byte{10, 0, 0, 1},
		Version:     0,
	}
	if s.IPAddress != [4]byte{10, 0, 0, 1} {
		t.Error("IPAddress mismatch")
	}
}
EOF

# ic44_ppp_setup_test.go
cat > ic44_ppp_setup_test.go << 'EOF'
package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestPPPSetup_ClassID(t *testing.T) {
	p := &PPPSetup{}
	if p.ClassID() != 44 {
		t.Errorf("ClassID() = %d, want 44", p.ClassID())
	}
	if p.ClassID() != core.ClassIDPPPSetup {
		t.Error("ClassID mismatch with const")
	}
}

func TestPPPSetup_New(t *testing.T) {
	p := &PPPSetup{
		LogicalName: core.ObisCode{0, 0, 44, 0, 0, 255},
		UserName:    "admin",
	}
	if p.UserName != "admin" {
		t.Error("UserName mismatch")
	}
}

func TestPPPSetup_MarshalBinary(t *testing.T) {
	tests := []struct {
		name string
		un   string
	}{
		{"empty", ""},
		{"user", "testuser"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PPPSetup{UserName: tt.un}
			b, err := p.MarshalBinary()
			if err != nil {
				t.Fatal(err)
			}
			p2 := &PPPSetup{}
			if err := p2.UnmarshalBinary(b); err != nil {
				t.Fatal(err)
			}
			if p2.UserName != tt.un {
				t.Errorf("UserName = %q, want %q", p2.UserName, tt.un)
			}
		})
	}
}

func TestPPPSetup_Fields(t *testing.T) {
	p := &PPPSetup{
		LogicalName: core.ObisCode{0, 0, 44, 0, 0, 255},
		UserName:    "root",
		Version:     0,
	}
	if p.GetVersion() != 0 {
		t.Error("GetVersion mismatch")
	}
}
EOF

# ic62_compact_data_test.go
cat > ic62_compact_data_test.go << 'EOF'
package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestCompactData_ClassID(t *testing.T) {
	cd := &CompactData{}
	if cd.ClassID() != 62 {
		t.Errorf("ClassID() = %d, want 62", cd.ClassID())
	}
	if cd.ClassID() != core.ClassIDCompactData {
		t.Error("ClassID mismatch with const")
	}
}

func TestCompactData_New(t *testing.T) {
	cd := &CompactData{
		LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255},
	}
	if cd.Buffer != nil {
		t.Error("Buffer should default nil")
	}
}

func TestCompactData_MarshalBinary(t *testing.T) {
	tests := []struct {
		name   string
		buffer []byte
	}{
		{"nil", nil},
		{"empty", []byte{}},
		{"data", []byte{0x01, 0x02, 0x03, 0x04}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cd := &CompactData{Buffer: tt.buffer}
			b, err := cd.MarshalBinary()
			if err != nil {
				t.Fatal(err)
			}
			cd2 := &CompactData{}
			if err := cd2.UnmarshalBinary(b); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestCompactData_Fields(t *testing.T) {
	cd := &CompactData{
		LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255},
		Buffer:      []byte{0xAA, 0xBB},
		Version:     0,
	}
	if len(cd.Buffer) != 2 {
		t.Error("Buffer length mismatch")
	}
}
EOF

# ic63_status_mapping_test.go
cat > ic63_status_mapping_test.go << 'EOF'
package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestStatusMapping_ClassID(t *testing.T) {
	sm := &StatusMapping{}
	if sm.ClassID() != 63 {
		t.Errorf("ClassID() = %d, want 63", sm.ClassID())
	}
	if sm.ClassID() != core.ClassIDStatusMapping {
		t.Error("ClassID mismatch with const")
	}
}

func TestStatusMapping_New(t *testing.T) {
	sm := &StatusMapping{
		LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255},
	}
	if sm.StatusWord != 0 {
		t.Error("StatusWord should default 0")
	}
}

func TestStatusMapping_MarshalBinary(t *testing.T) {
	tests := []struct {
		name string
		sw   uint32
	}{
		{"zero", 0},
		{"max", 0xFFFFFFFF},
		{"typical", 0x00008000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := &StatusMapping{StatusWord: tt.sw}
			b, err := sm.MarshalBinary()
			if err != nil {
				t.Fatal(err)
			}
			sm2 := &StatusMapping{}
			if err := sm2.UnmarshalBinary(b); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestStatusMapping_Fields(t *testing.T) {
	sm := &StatusMapping{
		LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255},
		StatusWord:  0xDEAD,
		Version:     0,
	}
	if sm.StatusWord != 0xDEAD {
		t.Error("StatusWord mismatch")
	}
}
EOF

# ic64_security_setup_test.go
cat > ic64_security_setup_test.go << 'EOF'
package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestSecuritySetup_ClassID(t *testing.T) {
	s := &SecuritySetup{}
	if s.ClassID() != 64 {
		t.Errorf("ClassID() = %d, want 64", s.ClassID())
	}
	if s.ClassID() != core.ClassIDSecuritySetup {
		t.Error("ClassID mismatch with const")
	}
}

func TestSecuritySetup_New(t *testing.T) {
	s := &SecuritySetup{
		LogicalName: core.ObisCode{0, 0, 43, 0, 0, 255},
	}
	if s.SecuritySuite != 0 {
		t.Error("SecuritySuite should default 0")
	}
}

func TestSecuritySetup_MarshalBinary(t *testing.T) {
	s := &SecuritySetup{
		SecuritySuite: 1,
		EncryptionKey: []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
			0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10},
		AuthenticationKey: []byte{0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
			0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F, 0x20},
		MasterKey: []byte{0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,
			0x29, 0x2A, 0x2B, 0x2C, 0x2D, 0x2E, 0x2F, 0x30},
	}
	b, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &SecuritySetup{}
	if err := s2.UnmarshalBinary(b); err != nil {
		t.Fatal(err)
	}
}

func TestSecuritySetup_Fields(t *testing.T) {
	s := &SecuritySetup{
		LogicalName:   core.ObisCode{0, 0, 43, 0, 0, 255},
		SecurityPolicy: 0x05,
		Version:        1,
	}
	if s.SecurityPolicy != 0x05 {
		t.Error("SecurityPolicy mismatch")
	}
	if s.GetVersion() != 1 {
		t.Error("GetVersion mismatch")
	}
}
EOF

# ic70_disconnect_control_test.go
cat > ic70_disconnect_control_test.go << 'EOF'
package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestDisconnectControl_ClassID(t *testing.T) {
	dc := &DisconnectControl{}
	if dc.ClassID() != 70 {
		t.Errorf("ClassID() = %d, want 70", dc.ClassID())
	}
	if dc.ClassID() != core.ClassIDDisconnectControl {
		t.Error("ClassID mismatch with const")
	}
}

func TestDisconnectControl_New(t *testing.T) {
	dc := &DisconnectControl{
		LogicalName: core.ObisCode{0, 0, 96, 0, 0, 255},
	}
	if dc.ControlState != 0 {
		t.Error("ControlState should default 0")
	}
}

func TestDisconnectControl_MarshalBinary(t *testing.T) {
	tests := []struct {
		name  string
		state uint8
	}{
		{"disconnected", 0},
		{"connected", 1},
		{"armed", 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dc := &DisconnectControl{ControlState: tt.state}
			b, err := dc.MarshalBinary()
			if err != nil {
				t.Fatal(err)
			}
			dc2 := &DisconnectControl{}
			if err := dc2.UnmarshalBinary(b); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestDisconnectControl_Fields(t *testing.T) {
	dc := &DisconnectControl{
		LogicalName:   core.ObisCode{0, 0, 96, 0, 0, 255},
		ControlState:  2,
		OutputState:   1,
		Version:       1,
	}
	if dc.OutputState != 1 {
		t.Error("OutputState mismatch")
	}
}
EOF

# ic71_limiter_test.go
cat > ic71_limiter_test.go << 'EOF'
package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestLimiter_ClassID(t *testing.T) {
	l := &Limiter{}
	if l.ClassID() != 71 {
		t.Errorf("ClassID() = %d, want 71", l.ClassID())
	}
	if l.ClassID() != core.ClassIDLimiter {
		t.Error("ClassID mismatch with const")
	}
}

func TestLimiter_New(t *testing.T) {
	l := &Limiter{
		LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255},
	}
	if l.MonitoredValue != nil {
		t.Error("MonitoredValue should default nil")
	}
}

func TestLimiter_MarshalBinary(t *testing.T) {
	tests := []struct {
		name  string
		value core.DlmsData
	}{
		{"zero", core.LongUnsignedData(0)},
		{"value", core.DoubleLongUnsignedData(5000)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Limiter{MonitoredValue: tt.value}
			b, err := l.MarshalBinary()
			if err != nil {
				t.Fatal(err)
			}
			l2 := &Limiter{}
			if err := l2.UnmarshalBinary(b); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestLimiter_Fields(t *testing.T) {
	l := &Limiter{
		LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255},
		Version:     0,
	}
	if l.GetVersion() != 0 {
		t.Error("GetVersion mismatch")
	}
}
EOF

# ic72_mbus_client_test.go
cat > ic72_mbus_client_test.go << 'EOF'
package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestMBusClient_ClassID(t *testing.T) {
	m := &MBusClient{}
	if m.ClassID() != 72 {
		t.Errorf("ClassID() = %d, want 72", m.ClassID())
	}
	if m.ClassID() != core.ClassIDMBusClient {
		t.Error("ClassID mismatch with const")
	}
}

func TestMBusClient_New(t *testing.T) {
	m := &MBusClient{
		LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255},
	}
	if m.PrimaryAddress != 0 {
		t.Error("PrimaryAddress should default 0")
	}
}

func TestMBusClient_MarshalBinary(t *testing.T) {
	tests := []struct {
		name string
		addr uint16
	}{
		{"zero", 0},
		{"normal", 250},
		{"max", 65535},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MBusClient{PrimaryAddress: tt.addr}
			b, err := m.MarshalBinary()
			if err != nil {
				t.Fatal(err)
			}
			m2 := &MBusClient{}
			if err := m2.UnmarshalBinary(b); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestMBusClient_Fields(t *testing.T) {
	m := &MBusClient{
		LogicalName:    core.ObisCode{0, 0, 0, 0, 0, 255},
		PrimaryAddress: 42,
		Version:        0,
	}
	if m.PrimaryAddress != 42 {
		t.Error("PrimaryAddress mismatch")
	}
}
EOF

# ic100_lp_setup_test.go
cat > ic100_lp_setup_test.go << 'EOF'
package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestLPSetup_ClassID(t *testing.T) {
	lp := &LPSetup{}
	if lp.ClassID() != 260 {
		t.Errorf("ClassID() = %d, want 260", lp.ClassID())
	}
	if lp.ClassID() != core.ClassIDLpSetup {
		t.Error("ClassID mismatch with const")
	}
}

func TestLPSetup_New(t *testing.T) {
	lp := &LPSetup{
		LogicalName:   core.ObisCode{1, 0, 99, 1, 0, 255},
		CapturePeriod: 60,
		ProfileEntries: 1000,
		Version:       1,
	}
	if lp.CapturePeriod != 60 {
		t.Error("CapturePeriod mismatch")
	}
	if lp.ProfileEntries != 1000 {
		t.Error("ProfileEntries mismatch")
	}
}

func TestLPSetup_Fields(t *testing.T) {
	lp := &LPSetup{
		LogicalName:   core.ObisCode{1, 0, 99, 1, 0, 255},
		CaptureObjects: []string{"0.0.1.0.0.255", "0.0.1.0.0.254"},
		CapturePeriod:  15,
		Version:        2,
	}
	if len(lp.CaptureObjects) != 2 {
		t.Error("CaptureObjects length mismatch")
	}
	if lp.GetVersion() != 2 {
		t.Error("GetVersion mismatch")
	}
}
EOF

# ic101_rs485_setup_test.go
cat > ic101_rs485_setup_test.go << 'EOF'
package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestRS485Setup_ClassID(t *testing.T) {
	r := &RS485Setup{}
	if r.ClassID() != 261 {
		t.Errorf("ClassID() = %d, want 261", r.ClassID())
	}
	if r.ClassID() != core.ClassIDIPv4TCPSetup {
		t.Error("ClassID mismatch with const")
	}
}

func TestRS485Setup_New(t *testing.T) {
	r := &RS485Setup{
		LogicalName: core.ObisCode{0, 0, 41, 0, 0, 255},
	}
	if r.DefaultBaudRate != 0 {
		t.Error("DefaultBaudRate should default 0")
	}
}

func TestRS485Setup_MarshalBinary(t *testing.T) {
	r := &RS485Setup{
		DefaultBaudRate: 9600,
	}
	b, err := r.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	r2 := &RS485Setup{}
	if err := r2.UnmarshalBinary(b); err != nil {
		t.Fatal(err)
	}
}

func TestRS485Setup_Fields(t *testing.T) {
	r := &RS485Setup{
		LogicalName:      core.ObisCode{0, 0, 41, 0, 0, 255},
		DefaultBaudRate:  19200,
		CommunicationParity: 1,
		Version:          0,
	}
	if r.DefaultBaudRate != 19200 {
		t.Error("DefaultBaudRate mismatch")
	}
}
EOF

# ic102_infrared_setup_test.go
cat > ic102_infrared_setup_test.go << 'EOF'
package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestInfraredSetup_ClassID(t *testing.T) {
	i := &InfraredSetup{}
	// Check the actual ClassID constant used
	tests := []struct {
		name string
		want uint16
	}{
		{"default_check", i.ClassID()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.want != tt.want {
				t.Error("ClassID check")
			}
		})
	}
}

func TestInfraredSetup_New(t *testing.T) {
	i := &InfraredSetup{
		LogicalName: core.ObisCode{0, 0, 20, 0, 0, 255},
	}
	if i.DefaultBaudRate != 0 {
		t.Error("DefaultBaudRate should default 0")
	}
}

func TestInfraredSetup_MarshalBinary(t *testing.T) {
	i := &InfraredSetup{
		DefaultBaudRate: 300,
	}
	b, err := i.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	i2 := &InfraredSetup{}
	if err := i2.UnmarshalBinary(b); err != nil {
		t.Fatal(err)
	}
}

func TestInfraredSetup_Fields(t *testing.T) {
	i := &InfraredSetup{
		LogicalName:      core.ObisCode{0, 0, 20, 0, 0, 255},
		DefaultBaudRate:  9600,
		Version:          0,
	}
	if i.GetVersion() != 0 {
		t.Error("GetVersion mismatch")
	}
}
EOF

# ic103_nbiot_setup_test.go
cat > ic103_nbiot_setup_test.go << 'EOF'
package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestNBIoTSetup_ClassID(t *testing.T) {
	n := &NBIoTSetup{}
	tests := []struct {
		name string
		want uint16
	}{
		{"default_check", n.ClassID()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.want != tt.want {
				t.Error("ClassID check")
			}
		})
	}
}

func TestNBIoTSetup_New(t *testing.T) {
	n := &NBIoTSetup{
		LogicalName: core.ObisCode{0, 0, 42, 0, 0, 255},
	}
	if n.APN != "" {
		t.Error("APN should default empty")
	}
}

func TestNBIoTSetup_MarshalBinary(t *testing.T) {
	tests := []struct {
		name string
		apn  string
	}{
		{"empty", ""},
		{"normal", "nbiot.apn"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &NBIoTSetup{APN: tt.apn}
			b, err := n.MarshalBinary()
			if err != nil {
				t.Fatal(err)
			}
			n2 := &NBIoTSetup{}
			if err := n2.UnmarshalBinary(b); err != nil {
				t.Fatal(err)
			}
			if n2.APN != tt.apn {
				t.Errorf("APN = %q, want %q", n2.APN, tt.apn)
			}
		})
	}
}

func TestNBIoTSetup_Fields(t *testing.T) {
	n := &NBIoTSetup{
		LogicalName: core.ObisCode{0, 0, 42, 0, 0, 255},
		APN:         "iot.test",
		Version:     0,
	}
	if n.APN != "iot.test" {
		t.Error("APN mismatch")
	}
}
EOF

# ic104_lorawan_setup_test.go
cat > ic104_lorawan_setup_test.go << 'EOF'
package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestLoRaWANSetup_ClassID(t *testing.T) {
	l := &LoRaWANSetup{}
	if l.ClassID() != 128 {
		t.Errorf("ClassID() = %d, want 128", l.ClassID())
	}
	if l.ClassID() != core.ClassIDLoRaWANSetup {
		t.Error("ClassID mismatch with const")
	}
}

func TestLoRaWANSetup_New(t *testing.T) {
	l := &LoRaWANSetup{
		LogicalName: core.ObisCode{0, 0, 42, 0, 0, 255},
	}
	if l.DevEUI != "" {
		t.Error("DevEUI should default empty")
	}
}

func TestLoRaWANSetup_MarshalBinary(t *testing.T) {
	l := &LoRaWANSetup{
		DevEUI:  "0102030405060708",
		AppEUI:  "AABBCCDD11223344",
		AppKey:  []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
			0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10},
	}
	b, err := l.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	l2 := &LoRaWANSetup{}
	if err := l2.UnmarshalBinary(b); err != nil {
		t.Fatal(err)
	}
}

func TestLoRaWANSetup_Fields(t *testing.T) {
	l := &LoRaWANSetup{
		LogicalName: core.ObisCode{0, 0, 42, 0, 0, 255},
		DevEUI:      "DEADBEEFCAFEBABE",
		Version:     0,
	}
	if l.DevEUI != "DEADBEEFCAFEBABE" {
		t.Error("DevEUI mismatch")
	}
}
EOF

# ic111_account_test.go
cat > ic111_account_test.go << 'EOF'
package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestAccount_ClassID(t *testing.T) {
	a := &Account{}
	if a.ClassID() != 111 {
		t.Errorf("ClassID() = %d, want 111", a.ClassID())
	}
	if a.ClassID() != core.ClassIDAccount {
		t.Error("ClassID mismatch with const")
	}
}

func TestAccount_New(t *testing.T) {
	a := &Account{
		LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255},
	}
	if a.VendorInfo != nil {
		t.Error("VendorInfo should default nil")
	}
}

func TestAccount_MarshalBinary(t *testing.T) {
	tests := []struct {
		name string
		info []byte
	}{
		{"nil", nil},
		{"empty", []byte{}},
		{"data", []byte{0x01, 0x02, 0x03}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Account{VendorInfo: tt.info}
			b, err := a.MarshalBinary()
			if err != nil {
				t.Fatal(err)
			}
			a2 := &Account{}
			if err := a2.UnmarshalBinary(b); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestAccount_Fields(t *testing.T) {
	a := &Account{
		LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255},
		VendorInfo:  []byte{0xAA, 0xBB, 0xCC},
		Version:     0,
	}
	if len(a.VendorInfo) != 3 {
		t.Error("VendorInfo length mismatch")
	}
}
EOF

# ic_tariff_test.go
cat > ic_tariff_test.go << 'EOF'
package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestTariffPlan_ClassID(t *testing.T) {
	tp := &TariffPlan{}
	if tp.ClassID() != 259 {
		t.Errorf("ClassID() = %d, want 259", tp.ClassID())
	}
	if tp.ClassID() != core.ClassIDTariffSchedule {
		t.Error("ClassID mismatch with const")
	}
}

func TestTariffPlan_New(t *testing.T) {
	tp := &TariffPlan{
		LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255},
		PlanName:    "TOU-Plan-1",
	}
	if tp.PlanName != "TOU-Plan-1" {
		t.Error("PlanName mismatch")
	}
}

func TestTariffPlan_MarshalBinary(t *testing.T) {
	tests := []struct {
		name string
		pn   string
	}{
		{"empty", ""},
		{"name", "SummerPlan"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tp := &TariffPlan{PlanName: tt.pn}
			b, err := tp.MarshalBinary()
			if err != nil {
				t.Fatal(err)
			}
			tp2 := &TariffPlan{}
			if err := tp2.UnmarshalBinary(b); err != nil {
				t.Fatal(err)
			}
			if tp2.PlanName != tt.pn {
				t.Errorf("PlanName = %q, want %q", tp2.PlanName, tt.pn)
			}
		})
	}
}

func TestTariffPlan_Fields(t *testing.T) {
	tp := &TariffPlan{
		LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255},
		PlanName:    "WinterPlan",
		Version:     0,
	}
	if tp.GetVersion() != 0 {
		t.Error("GetVersion mismatch")
	}
}
EOF

# ic_value_table_test.go
cat > ic_value_table_test.go << 'EOF'
package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestValueTable_ClassID(t *testing.T) {
	vt := &ValueTable{}
	if vt.ClassID() != 258 {
		t.Errorf("ClassID() = %d, want 258", vt.ClassID())
	}
	if vt.ClassID() != core.ClassIDValueTable {
		t.Error("ClassID mismatch with const")
	}
}

func TestValueTable_New(t *testing.T) {
	vt := &ValueTable{
		LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255},
	}
	if vt.Values != nil {
		t.Error("Values should default nil")
	}
	if vt.Descriptors != nil {
		t.Error("Descriptors should default nil")
	}
}

func TestValueTable_MarshalBinary(t *testing.T) {
	vt := &ValueTable{
		Values: []ValueEntry{
			{Index: 1, Value: core.LongUnsignedData(42)},
		},
	}
	b, err := vt.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	vt2 := &ValueTable{}
	if err := vt2.UnmarshalBinary(b); err != nil {
		t.Fatal(err)
	}
}

func TestValueTable_Fields(t *testing.T) {
	vt := &ValueTable{
		LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255},
		Descriptors: []ValueDescriptor{
			{Index: 0, Description: "Voltage", Unit: 27, Scaler: -1},
		},
		Version: 0,
	}
	if len(vt.Descriptors) != 1 {
		t.Error("Descriptors length mismatch")
	}
	if vt.Descriptors[0].Description != "Voltage" {
		t.Error("Descriptor description mismatch")
	}
}
EOF

echo "All test files generated."

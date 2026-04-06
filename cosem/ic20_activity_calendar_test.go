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
		SeasonProfileActive: nil,
		Version:       2,
	}
	if ac.GetVersion() != 2 {
		t.Error("GetVersion mismatch")
	}
}

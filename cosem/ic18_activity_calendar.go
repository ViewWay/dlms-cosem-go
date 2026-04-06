package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// SeasonProfileEntry represents a season profile
type SeasonProfileEntry struct {
	SeasonName     string
	SeasonStart    core.CosemDate
	WeekProfileRef string
}

// WeekProfileEntry represents a week profile
type WeekProfileEntry struct {
	ProfileName string
	Monday      string
	Tuesday     string
	Wednesday   string
	Thursday    string
	Friday      string
	Saturday    string
	Sunday      string
}

// ActivityCalendar is the COSEM Activity Calendar interface class (IC 18).
type ActivityCalendar struct {
	LogicalName             core.ObisCode
	CalendarName            string
	ActivatePassiveCalendar core.CosemDateTime
	SeasonProfileActive     []SeasonProfileEntry
	SeasonProfilePassive    []SeasonProfileEntry
	WeekProfileActive       []WeekProfileEntry
	WeekProfilePassive      []WeekProfileEntry
	DayProfileActiveName    string
	DayProfilePassiveName   string
	Version                 uint8
}

func (ac *ActivityCalendar) ClassID() uint16               { return core.ClassIDActivityCalendar }
func (ac *ActivityCalendar) GetLogicalName() core.ObisCode { return ac.LogicalName }
func (ac *ActivityCalendar) GetVersion() uint8             { return ac.Version }
func (ac *ActivityCalendar) Access(attr int) core.CosemAttributeAccess {
	switch attr {
	case 1:
		return core.CosemAttributeAccess{Access: core.AccessRead}
	case 2, 3:
		return core.CosemAttributeAccess{Access: core.AccessReadWrite}
	default:
		return core.CosemAttributeAccess{Access: core.AccessRead}
	}
}

func (ac *ActivityCalendar) MarshalBinary() ([]byte, error) {
	return core.VisibleStringData(ac.CalendarName).ToBytes(), nil
}

func (ac *ActivityCalendar) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	if v, ok := elem.(core.VisibleStringData); ok {
		ac.CalendarName = string(v)
	}
	return nil
}

func (ac *ActivityCalendar) SetCalendarName(name string) {
	ac.CalendarName = name
}

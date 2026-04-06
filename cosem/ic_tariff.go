package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// TariffPlan is a COSEM utility table.
type TariffPlan struct {
	LogicalName core.ObisCode
	PlanName    string
	Version     uint8
}

func (tp *TariffPlan) ClassID() uint16               { return core.ClassIDTariffSchedule }
func (tp *TariffPlan) GetLogicalName() core.ObisCode { return tp.LogicalName }
func (tp *TariffPlan) GetVersion() uint8             { return tp.Version }
func (tp *TariffPlan) Access(attr int) core.CosemAttributeAccess {
	return core.CosemAttributeAccess{Access: core.AccessRead}
}

func (tp *TariffPlan) MarshalBinary() ([]byte, error) {
	return core.VisibleStringData(tp.PlanName).ToBytes(), nil
}

func (tp *TariffPlan) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	if v, ok := elem.(core.VisibleStringData); ok {
		tp.PlanName = string(v)
	}
	return nil
}

// TariffTable contains tariff entries.
type TariffTable struct {
	LogicalName core.ObisCode
	Entries     []TariffEntry
	Version     uint8
}

func (tt *TariffTable) ClassID() uint16               { return core.ClassIDTariffSchedule }
func (tt *TariffTable) GetLogicalName() core.ObisCode { return tt.LogicalName }
func (tt *TariffTable) GetVersion() uint8             { return tt.Version }
func (tt *TariffTable) Access(attr int) core.CosemAttributeAccess {
	return core.CosemAttributeAccess{Access: core.AccessRead}
}

// TariffEntry represents a single tariff entry.
type TariffEntry struct {
	TariffIndex uint8
	Name        string
}

// SeasonProfile represents a season profile entry.
type SeasonProfile struct {
	SeasonName     string
	SeasonStart    core.CosemDate
	WeekProfileRef string
}

// WeekProfile represents a week profile entry.
type WeekProfile struct {
	ProfileName string
	Monday      string
	Tuesday     string
	Wednesday   string
	Thursday    string
	Friday      string
	Saturday    string
	Sunday      string
}

// DayProfile represents a day profile entry.
type DayProfile struct {
	DayProfileName string
	Schedule       []DaySchedule
}

// DaySchedule represents a day schedule entry.
type DaySchedule struct {
	StartTime core.CosemTime
	ScriptRef string
}

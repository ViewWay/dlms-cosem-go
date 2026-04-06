package cosem

import (
	"fmt"

	"github.com/ViewWay/dlms-cosem-go/core"
)

// Clock is the COSEM Clock interface class (IC 6).
type Clock struct {
	LogicalName              core.ObisCode
	Time                     core.CosemDateTime
	TimeZone                 int16
	Status                   uint8
	DaylightSavingsBegin     core.CosemDate
	DaylightSavingsEnd       core.CosemDate
	DaylightSavingsDeviation int16
	DaylightSavingsEnabled   bool
	ClockBase                core.CosemDateTime
	Version                  uint8
}

func (c *Clock) ClassID() uint16               { return core.ClassIDClock }
func (c *Clock) GetLogicalName() core.ObisCode { return c.LogicalName }
func (c *Clock) GetVersion() uint8             { return c.Version }
func (c *Clock) Access(attr int) core.CosemAttributeAccess {
	switch attr {
	case 2, 3, 4, 9, 10:
		return core.CosemAttributeAccess{Access: core.AccessReadWrite}
	default:
		return core.CosemAttributeAccess{Access: core.AccessRead}
	}
}

func (c *Clock) MarshalBinary() ([]byte, error) {
	return core.DateTimeData{Value: c.Time}.ToBytes(), nil
}

func (c *Clock) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	dt, ok := elem.(core.DateTimeData)
	if !ok {
		return fmt.Errorf("clock: expected DateTimeData")
	}
	c.Time = dt.Value
	return nil
}

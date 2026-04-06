package cosem

import (
	"fmt"

	"github.com/ViewWay/dlms-cosem-go/core"
)

// InfraredSetup is the infrared port setup.
type InfraredSetup struct {
	LogicalName     core.ObisCode
	DefaultBaudRate uint16
	CommSpeed       uint8
	Version         uint8
}

func (i *InfraredSetup) ClassID() uint16               { return core.ClassIDIECLocalPortSetup }
func (i *InfraredSetup) GetLogicalName() core.ObisCode { return i.LogicalName }
func (i *InfraredSetup) GetVersion() uint8             { return i.Version }
func (i *InfraredSetup) Access(attr int) core.CosemAttributeAccess {
	return core.CosemAttributeAccess{Access: core.AccessReadWrite}
}

func (i *InfraredSetup) MarshalBinary() ([]byte, error) {
	st := core.StructureData{
		core.UnsignedLongData(i.DefaultBaudRate),
		core.EnumData(i.CommSpeed),
	}
	return st.ToBytes(), nil
}

func (i *InfraredSetup) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok || len(st) < 2 {
		return fmt.Errorf("infrared_setup: expected structure")
	}
	if v, ok := st[0].(core.UnsignedLongData); ok {
		i.DefaultBaudRate = uint16(v)
	}
	if v, ok := st[1].(core.EnumData); ok {
		i.CommSpeed = uint8(v)
	}
	return nil
}

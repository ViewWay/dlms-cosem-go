package cosem

import (
	"fmt"

	"github.com/ViewWay/dlms-cosem-go/core"
)

// RS485Setup is the RS-485 port setup.
type RS485Setup struct {
	LogicalName     core.ObisCode
	DefaultBaudRate uint16
	CommSpeed       uint8
	Version         uint8
}

func (r *RS485Setup) ClassID() uint16               { return core.ClassIDIECHDLCSetup }
func (r *RS485Setup) GetLogicalName() core.ObisCode { return r.LogicalName }
func (r *RS485Setup) GetVersion() uint8             { return r.Version }
func (r *RS485Setup) Access(attr int) core.CosemAttributeAccess {
	return core.CosemAttributeAccess{Access: core.AccessReadWrite}
}

func (r *RS485Setup) MarshalBinary() ([]byte, error) {
	st := core.StructureData{
		core.UnsignedLongData(r.DefaultBaudRate),
		core.EnumData(r.CommSpeed),
	}
	return st.ToBytes(), nil
}

func (r *RS485Setup) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok || len(st) < 2 {
		return fmt.Errorf("rs485_setup: expected structure with 2 elements")
	}
	if v, ok := st[0].(core.UnsignedLongData); ok {
		r.DefaultBaudRate = uint16(v)
	}
	if v, ok := st[1].(core.EnumData); ok {
		r.CommSpeed = uint8(v)
	}
	return nil
}

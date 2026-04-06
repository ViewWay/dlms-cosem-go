package cosem

import (
	"fmt"

	"github.com/ViewWay/dlms-cosem-go/core"
)

// LoRaWANSetup is the LoRaWAN port setup.
type LoRaWANSetup struct {
	LogicalName core.ObisCode
	DevEUI      []byte
	AppEUI      []byte
	AppKey      []byte
	Version     uint8
}

func (l *LoRaWANSetup) ClassID() uint16               { return core.ClassIDLorawanSetup }
func (l *LoRaWANSetup) GetLogicalName() core.ObisCode { return l.LogicalName }
func (l *LoRaWANSetup) GetVersion() uint8             { return l.Version }
func (l *LoRaWANSetup) Access(attr int) core.CosemAttributeAccess {
	return core.CosemAttributeAccess{Access: core.AccessReadWrite}
}

func (l *LoRaWANSetup) MarshalBinary() ([]byte, error) {
	st := core.StructureData{
		core.OctetStringData(l.DevEUI),
		core.OctetStringData(l.AppEUI),
		core.OctetStringData(l.AppKey),
	}
	return st.ToBytes(), nil
}

func (l *LoRaWANSetup) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok || len(st) < 3 {
		return fmt.Errorf("lorawan_setup: expected structure with 3 elements")
	}
	if v, ok := st[0].(core.OctetStringData); ok {
		l.DevEUI = v
	}
	if v, ok := st[1].(core.OctetStringData); ok {
		l.AppEUI = v
	}
	if v, ok := st[2].(core.OctetStringData); ok {
		l.AppKey = v
	}
	return nil
}

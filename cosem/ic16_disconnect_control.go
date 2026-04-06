package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// DisconnectState represents the disconnect control state
type DisconnectState uint8

const (
	DisconnectStateDisconnected       DisconnectState = 0
	DisconnectStateConnected          DisconnectState = 1
	DisconnectStateReadyForDisconnect DisconnectState = 2
	DisconnectStateReadyForReconnect  DisconnectState = 3
	DisconnectStateArmed              DisconnectState = 4
)

// DisconnectControl is the COSEM Disconnect Control interface class (IC 16).
type DisconnectControl struct {
	LogicalName  core.ObisCode
	ControlState DisconnectState
	OutputState  DisconnectState
	ControlValue uint8
	Version      uint8
}

func (dc *DisconnectControl) ClassID() uint16               { return core.ClassIDDisconnectControl }
func (dc *DisconnectControl) GetLogicalName() core.ObisCode { return dc.LogicalName }
func (dc *DisconnectControl) GetVersion() uint8             { return dc.Version }
func (dc *DisconnectControl) Access(attr int) core.CosemAttributeAccess {
	switch attr {
	case 2, 3, 4:
		return core.CosemAttributeAccess{Access: core.AccessRead}
	default:
		return core.CosemAttributeAccess{Access: core.AccessNone}
	}
}

func (dc *DisconnectControl) MarshalBinary() ([]byte, error) {
	return core.EnumData(dc.ControlState).ToBytes(), nil
}

func (dc *DisconnectControl) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	if v, ok := elem.(core.EnumData); ok {
		dc.ControlState = DisconnectState(v)
	}
	return nil
}

func (dc *DisconnectControl) Disconnect() error {
	dc.ControlState = DisconnectStateDisconnected
	dc.OutputState = DisconnectStateDisconnected
	return nil
}

func (dc *DisconnectControl) Reconnect() error {
	dc.ControlState = DisconnectStateConnected
	dc.OutputState = DisconnectStateConnected
	return nil
}

func (dc *DisconnectControl) Arm() error {
	dc.ControlState = DisconnectStateArmed
	return nil
}

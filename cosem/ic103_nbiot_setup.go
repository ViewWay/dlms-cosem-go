package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// NBIoTSetup is the NB-IoT port setup.
type NBIoTSetup struct {
	LogicalName core.ObisCode
	APN         string
	Version     uint8
}

func (n *NBIoTSetup) ClassID() uint16               { return core.ClassIDGPRSModemSetup }
func (n *NBIoTSetup) GetLogicalName() core.ObisCode { return n.LogicalName }
func (n *NBIoTSetup) GetVersion() uint8             { return n.Version }
func (n *NBIoTSetup) Access(attr int) core.CosemAttributeAccess {
	return core.CosemAttributeAccess{Access: core.AccessReadWrite}
}

func (n *NBIoTSetup) MarshalBinary() ([]byte, error) {
	return core.VisibleStringData(n.APN).ToBytes(), nil
}

func (n *NBIoTSetup) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	if v, ok := elem.(core.VisibleStringData); ok {
		n.APN = string(v)
	}
	return nil
}

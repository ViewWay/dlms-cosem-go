package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// PPPSetup is the COSEM PPP Setup interface class (IC 45).
type PPPSetup struct {
	LogicalName             core.ObisCode
	UserName                string
	Password                string
	PhoneNumber             string
	AuthenticationProtocol  uint8
	Enabled                 bool
	Version                 uint8
}

func (p *PPPSetup) ClassID() uint16               { return core.ClassIDPPPSetup }
func (p *PPPSetup) GetLogicalName() core.ObisCode { return p.LogicalName }
func (p *PPPSetup) GetVersion() uint8             { return p.Version }
func (p *PPPSetup) Access(attr int) core.CosemAttributeAccess {
	switch attr {
	case 1:
		return core.CosemAttributeAccess{Access: core.AccessRead}
	case 2, 3, 4, 5, 6:
		return core.CosemAttributeAccess{Access: core.AccessReadWrite}
	default:
		return core.CosemAttributeAccess{Access: core.AccessNone}
	}
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (p *PPPSetup) MarshalBinary() ([]byte, error) {
	return core.VisibleStringData(p.UserName).ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (p *PPPSetup) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	if v, ok := elem.(core.VisibleStringData); ok {
		p.UserName = string(v)
	}
	return nil
}

package cosem

import (
	"fmt"

	"github.com/ViewWay/dlms-cosem-go/core"
)

// SecuritySetup is the COSEM Security Setup interface class (IC 29).
type SecuritySetup struct {
	LogicalName       core.ObisCode
	SecuritySuite     uint8
	EncryptionKey     []byte
	AuthenticationKey []byte
	MasterKey         []byte
	SecurityPolicy    uint8
	Version           uint8
}

func (s *SecuritySetup) ClassID() uint16               { return core.ClassIDSecuritySetup }
func (s *SecuritySetup) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *SecuritySetup) GetVersion() uint8             { return s.Version }
func (s *SecuritySetup) Access(attr int) core.CosemAttributeAccess {
	switch attr {
	case 2:
		return core.CosemAttributeAccess{Access: core.AccessRead}
	case 3, 4, 5:
		return core.CosemAttributeAccess{Access: core.AccessReadWrite, Set: core.AccessHigh}
	default:
		return core.CosemAttributeAccess{Access: core.AccessNone}
	}
}

func (s *SecuritySetup) MarshalBinary() ([]byte, error) {
	st := core.StructureData{
		core.UnsignedIntegerData(s.SecuritySuite),
		core.OctetStringData(s.EncryptionKey),
		core.OctetStringData(s.AuthenticationKey),
	}
	return st.ToBytes(), nil
}

func (s *SecuritySetup) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok || len(st) < 1 {
		return fmt.Errorf("security_setup: expected structure")
	}
	if v, ok := st[0].(core.UnsignedIntegerData); ok {
		s.SecuritySuite = uint8(v)
	}
	if len(st) > 1 {
		if v, ok := st[1].(core.OctetStringData); ok {
			s.EncryptionKey = v
		}
	}
	if len(st) > 2 {
		if v, ok := st[2].(core.OctetStringData); ok {
			s.AuthenticationKey = v
		}
	}
	return nil
}

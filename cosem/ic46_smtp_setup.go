package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// SMTPSetup is the COSEM SMTPSetup interface class (IC 46).
type SMTPSetup struct {
	LogicalName core.ObisCode
	Version     uint8
	ServerAddress                       string
	UserName                            string
	Password                            string
	SMTPPort                            uint16
	SenderAddress                       string
}

func (s *SMTPSetup) ClassID() uint16               { return core.ClassIDSMTPSetup }
func (s *SMTPSetup) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *SMTPSetup) GetVersion() uint8             { return s.Version }

func (s *SMTPSetup) Access(attr int) core.CosemAttributeAccess {
	switch attr {
	case 1:
		return core.CosemAttributeAccess{Access: core.AccessRead}
	case 2:
		return core.CosemAttributeAccess{Access: core.AccessReadWrite}
	default:
		return core.CosemAttributeAccess{Access: core.AccessNone}
	}
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s *SMTPSetup) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.VisibleStringData(s.ServerAddress),
		core.VisibleStringData(s.UserName),
		core.VisibleStringData(s.Password),
		core.UnsignedLongData(s.SMTPPort),
		core.VisibleStringData(s.SenderAddress),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *SMTPSetup) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return nil
	}
	if 0 < len(st) {
		if v, ok := st[0].(core.VisibleStringData); ok { s.ServerAddress = string(v) }
	}
	if 1 < len(st) {
		if v, ok := st[1].(core.VisibleStringData); ok { s.UserName = string(v) }
	}
	if 2 < len(st) {
		if v, ok := st[2].(core.VisibleStringData); ok { s.Password = string(v) }
	}
	if 3 < len(st) {
		if v, ok := st[3].(core.UnsignedLongData); ok { s.SMTPPort = uint16(v) }
	}
	if 4 < len(st) {
		if v, ok := st[4].(core.VisibleStringData); ok { s.SenderAddress = string(v) }
	}
	return nil
}


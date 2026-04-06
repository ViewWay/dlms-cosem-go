package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// SFSKMACSyncTimeouts is the COSEM SFSKMACSyncTimeouts interface class (IC 52).
type SFSKMACSyncTimeouts struct {
	LogicalName core.ObisCode
	Version     uint8
	Timeouts                            []byte
}

func (s *SFSKMACSyncTimeouts) ClassID() uint16               { return core.ClassIDSFSKMACSyncTimeouts }
func (s *SFSKMACSyncTimeouts) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *SFSKMACSyncTimeouts) GetVersion() uint8             { return s.Version }

func (s *SFSKMACSyncTimeouts) Access(attr int) core.CosemAttributeAccess {
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
func (s *SFSKMACSyncTimeouts) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.OctetStringData(s.Timeouts),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *SFSKMACSyncTimeouts) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return nil
	}
	if 0 < len(st) {
		if v, ok := st[0].(core.OctetStringData); ok { s.Timeouts = []byte(v) }
	}
	return nil
}


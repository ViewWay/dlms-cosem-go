package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// SAPAssignment is the COSEM SAPAssignment interface class (IC 17).
type SAPAssignment struct {
	LogicalName core.ObisCode
	Version     uint8
	SAPAssignmentList                   []byte
}

func (s *SAPAssignment) ClassID() uint16               { return core.ClassIDSAPAssignment }
func (s *SAPAssignment) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *SAPAssignment) GetVersion() uint8             { return s.Version }

func (s *SAPAssignment) Access(attr int) core.CosemAttributeAccess {
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
func (s *SAPAssignment) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.OctetStringData(s.SAPAssignmentList),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *SAPAssignment) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return nil
	}
	if 0 < len(st) {
		if v, ok := st[0].(core.OctetStringData); ok { s.SAPAssignmentList = []byte(v) }
	}
	return nil
}


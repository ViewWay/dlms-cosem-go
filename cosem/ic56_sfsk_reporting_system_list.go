package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// SFSKReportingSystemList is the COSEM SFSKReportingSystemList interface class (IC 56).
type SFSKReportingSystemList struct {
	LogicalName core.ObisCode
	Version     uint8
	ReportingSystemList                 []byte
}

func (s *SFSKReportingSystemList) ClassID() uint16               { return core.ClassIDSFSKReportingSystemList }
func (s *SFSKReportingSystemList) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *SFSKReportingSystemList) GetVersion() uint8             { return s.Version }

func (s *SFSKReportingSystemList) Access(attr int) core.CosemAttributeAccess {
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
func (s *SFSKReportingSystemList) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.OctetStringData(s.ReportingSystemList),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *SFSKReportingSystemList) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return nil
	}
	if 0 < len(st) {
		if v, ok := st[0].(core.OctetStringData); ok { s.ReportingSystemList = []byte(v) }
	}
	return nil
}


package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// RPLDiagnostic is the COSEM RPLDiagnostic interface class (IC 97).
type RPLDiagnostic struct {
	LogicalName core.ObisCode
	Version     uint8
	ParentAddress                       [6]byte
	ParentRank                          uint16
	ParentLinkMetric                    uint16
	ParentLinkMetricType                uint8
	ParentSwitches                      uint32
	ChildrenAddresses                   []byte
	ChildrenRanks                       []byte
	DAOMessagesSent                     uint32
	DAOMessagesReceived                 uint32
	DIOMessagesSent                     uint32
	DIOMessagesReceived                 uint32
}

func (s *RPLDiagnostic) ClassID() uint16               { return core.ClassIDRPLDiagnostic }
func (s *RPLDiagnostic) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *RPLDiagnostic) GetVersion() uint8             { return s.Version }

func (s *RPLDiagnostic) Access(attr int) core.CosemAttributeAccess {
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
func (s *RPLDiagnostic) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.OctetStringData(s.ParentAddress[:]),
		core.UnsignedLongData(s.ParentRank),
		core.UnsignedLongData(s.ParentLinkMetric),
		core.UnsignedIntegerData(s.ParentLinkMetricType),
		core.DoubleLongUnsignedData(s.ParentSwitches),
		core.OctetStringData(s.ChildrenAddresses),
		core.OctetStringData(s.ChildrenRanks),
		core.DoubleLongUnsignedData(s.DAOMessagesSent),
		core.DoubleLongUnsignedData(s.DAOMessagesReceived),
		core.DoubleLongUnsignedData(s.DIOMessagesSent),
		core.DoubleLongUnsignedData(s.DIOMessagesReceived),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *RPLDiagnostic) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return nil
	}
	if 0 < len(st) {
		if v, ok := st[0].(core.OctetStringData); ok && len(v) == 6 { copy(s.ParentAddress[:], v) }
	}
	if 1 < len(st) {
		if v, ok := st[1].(core.UnsignedLongData); ok { s.ParentRank = uint16(v) }
	}
	if 2 < len(st) {
		if v, ok := st[2].(core.UnsignedLongData); ok { s.ParentLinkMetric = uint16(v) }
	}
	if 3 < len(st) {
		if v, ok := st[3].(core.UnsignedIntegerData); ok { s.ParentLinkMetricType = uint8(v) }
	}
	if 4 < len(st) {
		if v, ok := st[4].(core.DoubleLongUnsignedData); ok { s.ParentSwitches = uint32(v) }
	}
	if 5 < len(st) {
		if v, ok := st[5].(core.OctetStringData); ok { s.ChildrenAddresses = []byte(v) }
	}
	if 6 < len(st) {
		if v, ok := st[6].(core.OctetStringData); ok { s.ChildrenRanks = []byte(v) }
	}
	if 7 < len(st) {
		if v, ok := st[7].(core.DoubleLongUnsignedData); ok { s.DAOMessagesSent = uint32(v) }
	}
	if 8 < len(st) {
		if v, ok := st[8].(core.DoubleLongUnsignedData); ok { s.DAOMessagesReceived = uint32(v) }
	}
	if 9 < len(st) {
		if v, ok := st[9].(core.DoubleLongUnsignedData); ok { s.DIOMessagesSent = uint32(v) }
	}
	if 10 < len(st) {
		if v, ok := st[10].(core.DoubleLongUnsignedData); ok { s.DIOMessagesReceived = uint32(v) }
	}
	return nil
}


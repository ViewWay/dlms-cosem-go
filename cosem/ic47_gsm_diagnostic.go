package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// GSMDiagnostic is the COSEM GSMDiagnostic interface class (IC 47).
type GSMDiagnostic struct {
	LogicalName core.ObisCode
	Version     uint8
	Operator                            string
	Status                              uint8
	CircuitSwitchStatus                 uint8
	PacketSwitchStatus                  uint8
	CellID                              []byte
	LocationArea                        []byte
	VCI                                 []byte
	MCC                                 []byte
	MNC                                 []byte
	BaseStationID                       []byte
}

func (s *GSMDiagnostic) ClassID() uint16               { return core.ClassIDGSMDiagnostic }
func (s *GSMDiagnostic) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *GSMDiagnostic) GetVersion() uint8             { return s.Version }

func (s *GSMDiagnostic) Access(attr int) core.CosemAttributeAccess {
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
func (s *GSMDiagnostic) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.VisibleStringData(s.Operator),
		core.UnsignedIntegerData(s.Status),
		core.UnsignedIntegerData(s.CircuitSwitchStatus),
		core.UnsignedIntegerData(s.PacketSwitchStatus),
		core.OctetStringData(s.CellID),
		core.OctetStringData(s.LocationArea),
		core.OctetStringData(s.VCI),
		core.OctetStringData(s.MCC),
		core.OctetStringData(s.MNC),
		core.OctetStringData(s.BaseStationID),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *GSMDiagnostic) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return nil
	}
	if 0 < len(st) {
		if v, ok := st[0].(core.VisibleStringData); ok { s.Operator = string(v) }
	}
	if 1 < len(st) {
		if v, ok := st[1].(core.UnsignedIntegerData); ok { s.Status = uint8(v) }
	}
	if 2 < len(st) {
		if v, ok := st[2].(core.UnsignedIntegerData); ok { s.CircuitSwitchStatus = uint8(v) }
	}
	if 3 < len(st) {
		if v, ok := st[3].(core.UnsignedIntegerData); ok { s.PacketSwitchStatus = uint8(v) }
	}
	if 4 < len(st) {
		if v, ok := st[4].(core.OctetStringData); ok { s.CellID = []byte(v) }
	}
	if 5 < len(st) {
		if v, ok := st[5].(core.OctetStringData); ok { s.LocationArea = []byte(v) }
	}
	if 6 < len(st) {
		if v, ok := st[6].(core.OctetStringData); ok { s.VCI = []byte(v) }
	}
	if 7 < len(st) {
		if v, ok := st[7].(core.OctetStringData); ok { s.MCC = []byte(v) }
	}
	if 8 < len(st) {
		if v, ok := st[8].(core.OctetStringData); ok { s.MNC = []byte(v) }
	}
	if 9 < len(st) {
		if v, ok := st[9].(core.OctetStringData); ok { s.BaseStationID = []byte(v) }
	}
	return nil
}


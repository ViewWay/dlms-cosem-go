package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// SCHCLPWANSetup is the COSEM SCHCLPWANSetup interface class (IC 126).
type SCHCLPWANSetup struct {
	LogicalName core.ObisCode
	Version     uint8
	DeviceID                            []byte
	RuleID                              uint8
	Direction                           uint8
	Mode                                uint8
	DTAGSize                            uint8
	DTAG                                uint8
	WSize                               uint8
	Bitmap                              []byte
	FCN                                 uint8
	MIC                                 []byte
	Payload                             []byte
}

func (s *SCHCLPWANSetup) ClassID() uint16               { return core.ClassIDSCHCLPWANSetup }
func (s *SCHCLPWANSetup) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *SCHCLPWANSetup) GetVersion() uint8             { return s.Version }

func (s *SCHCLPWANSetup) Access(attr int) core.CosemAttributeAccess {
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
func (s *SCHCLPWANSetup) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.OctetStringData(s.DeviceID),
		core.UnsignedIntegerData(s.RuleID),
		core.UnsignedIntegerData(s.Direction),
		core.UnsignedIntegerData(s.Mode),
		core.UnsignedIntegerData(s.DTAGSize),
		core.UnsignedIntegerData(s.DTAG),
		core.UnsignedIntegerData(s.WSize),
		core.OctetStringData(s.Bitmap),
		core.UnsignedIntegerData(s.FCN),
		core.OctetStringData(s.MIC),
		core.OctetStringData(s.Payload),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *SCHCLPWANSetup) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return nil
	}
	if 0 < len(st) {
		if v, ok := st[0].(core.OctetStringData); ok { s.DeviceID = []byte(v) }
	}
	if 1 < len(st) {
		if v, ok := st[1].(core.UnsignedIntegerData); ok { s.RuleID = uint8(v) }
	}
	if 2 < len(st) {
		if v, ok := st[2].(core.UnsignedIntegerData); ok { s.Direction = uint8(v) }
	}
	if 3 < len(st) {
		if v, ok := st[3].(core.UnsignedIntegerData); ok { s.Mode = uint8(v) }
	}
	if 4 < len(st) {
		if v, ok := st[4].(core.UnsignedIntegerData); ok { s.DTAGSize = uint8(v) }
	}
	if 5 < len(st) {
		if v, ok := st[5].(core.UnsignedIntegerData); ok { s.DTAG = uint8(v) }
	}
	if 6 < len(st) {
		if v, ok := st[6].(core.UnsignedIntegerData); ok { s.WSize = uint8(v) }
	}
	if 7 < len(st) {
		if v, ok := st[7].(core.OctetStringData); ok { s.Bitmap = []byte(v) }
	}
	if 8 < len(st) {
		if v, ok := st[8].(core.UnsignedIntegerData); ok { s.FCN = uint8(v) }
	}
	if 9 < len(st) {
		if v, ok := st[9].(core.OctetStringData); ok { s.MIC = []byte(v) }
	}
	if 10 < len(st) {
		if v, ok := st[10].(core.OctetStringData); ok { s.Payload = []byte(v) }
	}
	return nil
}


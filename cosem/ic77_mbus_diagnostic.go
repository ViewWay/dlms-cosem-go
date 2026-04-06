package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// MBusDiagnostic is the COSEM MBusDiagnostic interface class (IC 77).
type MBusDiagnostic struct {
	LogicalName core.ObisCode
	Version     uint8
	BusStuck                            bool
	StopByteFound                       bool
	FrameFormatError                    uint32
	FrameLengthError                    uint32
	CRCError                            uint32
	FrameTimeout                        uint32
	SlaveList                           []byte
}

func (s *MBusDiagnostic) ClassID() uint16               { return core.ClassIDMBusDiagnostic }
func (s *MBusDiagnostic) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *MBusDiagnostic) GetVersion() uint8             { return s.Version }

func (s *MBusDiagnostic) Access(attr int) core.CosemAttributeAccess {
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
func (s *MBusDiagnostic) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.BooleanData(s.BusStuck),
		core.BooleanData(s.StopByteFound),
		core.DoubleLongUnsignedData(s.FrameFormatError),
		core.DoubleLongUnsignedData(s.FrameLengthError),
		core.DoubleLongUnsignedData(s.CRCError),
		core.DoubleLongUnsignedData(s.FrameTimeout),
		core.OctetStringData(s.SlaveList),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *MBusDiagnostic) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return nil
	}
	if 0 < len(st) {
		if v, ok := st[0].(core.BooleanData); ok { s.BusStuck = bool(v) }
	}
	if 1 < len(st) {
		if v, ok := st[1].(core.BooleanData); ok { s.StopByteFound = bool(v) }
	}
	if 2 < len(st) {
		if v, ok := st[2].(core.DoubleLongUnsignedData); ok { s.FrameFormatError = uint32(v) }
	}
	if 3 < len(st) {
		if v, ok := st[3].(core.DoubleLongUnsignedData); ok { s.FrameLengthError = uint32(v) }
	}
	if 4 < len(st) {
		if v, ok := st[4].(core.DoubleLongUnsignedData); ok { s.CRCError = uint32(v) }
	}
	if 5 < len(st) {
		if v, ok := st[5].(core.DoubleLongUnsignedData); ok { s.FrameTimeout = uint32(v) }
	}
	if 6 < len(st) {
		if v, ok := st[6].(core.OctetStringData); ok { s.SlaveList = []byte(v) }
	}
	return nil
}


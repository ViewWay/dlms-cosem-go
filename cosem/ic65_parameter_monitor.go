package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// ParameterMonitor is the COSEM ParameterMonitor interface class (IC 65).
type ParameterMonitor struct {
	LogicalName core.ObisCode
	Version     uint8
	CapturedValue                       []byte
	CapturedTime                        []byte
	Status                              uint8
	CaptureMethod                       uint8
}

func (s *ParameterMonitor) ClassID() uint16               { return core.ClassIDParameterMonitor }
func (s *ParameterMonitor) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *ParameterMonitor) GetVersion() uint8             { return s.Version }

func (s *ParameterMonitor) Access(attr int) core.CosemAttributeAccess {
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
func (s *ParameterMonitor) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.OctetStringData(s.CapturedValue),
		core.OctetStringData(s.CapturedTime),
		core.UnsignedIntegerData(s.Status),
		core.UnsignedIntegerData(s.CaptureMethod),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ParameterMonitor) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return nil
	}
	if 0 < len(st) {
		if v, ok := st[0].(core.OctetStringData); ok { s.CapturedValue = []byte(v) }
	}
	if 1 < len(st) {
		if v, ok := st[1].(core.OctetStringData); ok { s.CapturedTime = []byte(v) }
	}
	if 2 < len(st) {
		if v, ok := st[2].(core.UnsignedIntegerData); ok { s.Status = uint8(v) }
	}
	if 3 < len(st) {
		if v, ok := st[3].(core.UnsignedIntegerData); ok { s.CaptureMethod = uint8(v) }
	}
	return nil
}


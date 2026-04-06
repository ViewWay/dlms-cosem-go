package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// DLMSMBusPortSetup is the COSEM DLMSMBusPortSetup interface class (IC 76).
type DLMSMBusPortSetup struct {
	LogicalName core.ObisCode
	Version     uint8
	MBusPortReference                   core.ObisCode
	ListenPort                          uint16
	SlaveDevices                        []byte
	ClientActive                        bool
}

func (s *DLMSMBusPortSetup) ClassID() uint16               { return core.ClassIDDLMSMBusPortSetup }
func (s *DLMSMBusPortSetup) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *DLMSMBusPortSetup) GetVersion() uint8             { return s.Version }

func (s *DLMSMBusPortSetup) Access(attr int) core.CosemAttributeAccess {
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
func (s *DLMSMBusPortSetup) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.OctetStringData(s.MBusPortReference[:]),
		core.UnsignedLongData(s.ListenPort),
		core.OctetStringData(s.SlaveDevices),
		core.BooleanData(s.ClientActive),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *DLMSMBusPortSetup) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return nil
	}
	if 0 < len(st) {
		if v, ok := st[0].(core.OctetStringData); ok && len(v) == 6 { copy(s.MBusPortReference[:], v) }
	}
	if 1 < len(st) {
		if v, ok := st[1].(core.UnsignedLongData); ok { s.ListenPort = uint16(v) }
	}
	if 2 < len(st) {
		if v, ok := st[2].(core.OctetStringData); ok { s.SlaveDevices = []byte(v) }
	}
	if 3 < len(st) {
		if v, ok := st[3].(core.BooleanData); ok { s.ClientActive = bool(v) }
	}
	return nil
}


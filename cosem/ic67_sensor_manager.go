package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// SensorManager is the COSEM SensorManager interface class (IC 67).
type SensorManager struct {
	LogicalName core.ObisCode
	Version     uint8
	SensorList                          []byte
}

func (s *SensorManager) ClassID() uint16               { return core.ClassIDSensorManager }
func (s *SensorManager) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *SensorManager) GetVersion() uint8             { return s.Version }

func (s *SensorManager) Access(attr int) core.CosemAttributeAccess {
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
func (s *SensorManager) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.OctetStringData(s.SensorList),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *SensorManager) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return nil
	}
	if 0 < len(st) {
		if v, ok := st[0].(core.OctetStringData); ok { s.SensorList = []byte(v) }
	}
	return nil
}


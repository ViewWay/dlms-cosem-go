package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// PrimeAppIdentification is the COSEM PrimeAppIdentification interface class (IC 86).
type PrimeAppIdentification struct {
	LogicalName core.ObisCode
	Version     uint8
	ApplicationName                     string
	ApplicationVersion                  string
}

func (s *PrimeAppIdentification) ClassID() uint16               { return core.ClassIDPrimeAppIdentification }
func (s *PrimeAppIdentification) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *PrimeAppIdentification) GetVersion() uint8             { return s.Version }

func (s *PrimeAppIdentification) Access(attr int) core.CosemAttributeAccess {
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
func (s *PrimeAppIdentification) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.VisibleStringData(s.ApplicationName),
		core.VisibleStringData(s.ApplicationVersion),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *PrimeAppIdentification) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return nil
	}
	if 0 < len(st) {
		if v, ok := st[0].(core.VisibleStringData); ok { s.ApplicationName = string(v) }
	}
	if 1 < len(st) {
		if v, ok := st[1].(core.VisibleStringData); ok { s.ApplicationVersion = string(v) }
	}
	return nil
}


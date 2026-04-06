package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// IEC61334LLCSetup is the COSEM IEC61334LLCSetup interface class (IC 55).
type IEC61334LLCSetup struct {
	LogicalName core.ObisCode
	Version     uint8
	LLCType1Enable                      bool
	LLCType2Enable                      bool
	LLCType3Enable                      bool
}

func (s *IEC61334LLCSetup) ClassID() uint16               { return core.ClassIDIEC61334LLCSetup }
func (s *IEC61334LLCSetup) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *IEC61334LLCSetup) GetVersion() uint8             { return s.Version }

func (s *IEC61334LLCSetup) Access(attr int) core.CosemAttributeAccess {
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
func (s *IEC61334LLCSetup) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.BooleanData(s.LLCType1Enable),
		core.BooleanData(s.LLCType2Enable),
		core.BooleanData(s.LLCType3Enable),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *IEC61334LLCSetup) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return nil
	}
	if 0 < len(st) {
		if v, ok := st[0].(core.BooleanData); ok { s.LLCType1Enable = bool(v) }
	}
	if 1 < len(st) {
		if v, ok := st[1].(core.BooleanData); ok { s.LLCType2Enable = bool(v) }
	}
	if 2 < len(st) {
		if v, ok := st[2].(core.BooleanData); ok { s.LLCType3Enable = bool(v) }
	}
	return nil
}


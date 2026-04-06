package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// MBusMasterPortSetup is the COSEM MBusMasterPortSetup interface class (IC 74).
type MBusMasterPortSetup struct {
	LogicalName core.ObisCode
	Version     uint8
	CommSpeed                           uint32
	RecTimeout                          uint16
	SendTimeout                         uint16
}

func (s *MBusMasterPortSetup) ClassID() uint16               { return core.ClassIDMBusMasterPortSetup }
func (s *MBusMasterPortSetup) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *MBusMasterPortSetup) GetVersion() uint8             { return s.Version }

func (s *MBusMasterPortSetup) Access(attr int) core.CosemAttributeAccess {
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
func (s *MBusMasterPortSetup) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.DoubleLongUnsignedData(s.CommSpeed),
		core.UnsignedLongData(s.RecTimeout),
		core.UnsignedLongData(s.SendTimeout),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *MBusMasterPortSetup) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return nil
	}
	if 0 < len(st) {
		if v, ok := st[0].(core.DoubleLongUnsignedData); ok { s.CommSpeed = uint32(v) }
	}
	if 1 < len(st) {
		if v, ok := st[1].(core.UnsignedLongData); ok { s.RecTimeout = uint16(v) }
	}
	if 2 < len(st) {
		if v, ok := st[2].(core.UnsignedLongData); ok { s.SendTimeout = uint16(v) }
	}
	return nil
}


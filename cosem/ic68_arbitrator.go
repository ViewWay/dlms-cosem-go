package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// Arbitrator is the COSEM Arbitrator interface class (IC 68).
type Arbitrator struct {
	LogicalName core.ObisCode
	Version     uint8
	ArbitratorList                      []byte
	Weights                             []byte
	LongestBackupTime                   uint32
}

func (s *Arbitrator) ClassID() uint16               { return core.ClassIDArbitrator }
func (s *Arbitrator) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *Arbitrator) GetVersion() uint8             { return s.Version }

func (s *Arbitrator) Access(attr int) core.CosemAttributeAccess {
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
func (s *Arbitrator) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.OctetStringData(s.ArbitratorList),
		core.OctetStringData(s.Weights),
		core.DoubleLongUnsignedData(s.LongestBackupTime),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *Arbitrator) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return nil
	}
	if 0 < len(st) {
		if v, ok := st[0].(core.OctetStringData); ok { s.ArbitratorList = []byte(v) }
	}
	if 1 < len(st) {
		if v, ok := st[1].(core.OctetStringData); ok { s.Weights = []byte(v) }
	}
	if 2 < len(st) {
		if v, ok := st[2].(core.DoubleLongUnsignedData); ok { s.LongestBackupTime = uint32(v) }
	}
	return nil
}


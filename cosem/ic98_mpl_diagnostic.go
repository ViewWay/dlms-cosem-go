package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// MPLDiagnostic is the COSEM MPLDiagnostic interface class (IC 98).
type MPLDiagnostic struct {
	LogicalName core.ObisCode
	Version     uint8
	MPLDomainID                         uint8
	MPLSeedSetVersion                   uint8
	MPLTrickleTimerExpirations          uint32
	MPLMessagesSent                     uint32
	MPLMessagesReceived                 uint32
	MPLMessagesForwarded                uint32
}

func (s *MPLDiagnostic) ClassID() uint16               { return core.ClassIDMPLDiagnostic }
func (s *MPLDiagnostic) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *MPLDiagnostic) GetVersion() uint8             { return s.Version }

func (s *MPLDiagnostic) Access(attr int) core.CosemAttributeAccess {
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
func (s *MPLDiagnostic) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.UnsignedIntegerData(s.MPLDomainID),
		core.UnsignedIntegerData(s.MPLSeedSetVersion),
		core.DoubleLongUnsignedData(s.MPLTrickleTimerExpirations),
		core.DoubleLongUnsignedData(s.MPLMessagesSent),
		core.DoubleLongUnsignedData(s.MPLMessagesReceived),
		core.DoubleLongUnsignedData(s.MPLMessagesForwarded),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *MPLDiagnostic) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return nil
	}
	if 0 < len(st) {
		if v, ok := st[0].(core.UnsignedIntegerData); ok { s.MPLDomainID = uint8(v) }
	}
	if 1 < len(st) {
		if v, ok := st[1].(core.UnsignedIntegerData); ok { s.MPLSeedSetVersion = uint8(v) }
	}
	if 2 < len(st) {
		if v, ok := st[2].(core.DoubleLongUnsignedData); ok { s.MPLTrickleTimerExpirations = uint32(v) }
	}
	if 3 < len(st) {
		if v, ok := st[3].(core.DoubleLongUnsignedData); ok { s.MPLMessagesSent = uint32(v) }
	}
	if 4 < len(st) {
		if v, ok := st[4].(core.DoubleLongUnsignedData); ok { s.MPLMessagesReceived = uint32(v) }
	}
	if 5 < len(st) {
		if v, ok := st[5].(core.DoubleLongUnsignedData); ok { s.MPLMessagesForwarded = uint32(v) }
	}
	return nil
}


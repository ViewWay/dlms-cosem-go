package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// TokenGateway is the COSEM TokenGateway interface class (IC 115).
type TokenGateway struct {
	LogicalName core.ObisCode
	Version     uint8
	Token                               []byte
	TokenTime                           []byte
	TokenStatus                         uint8
}

func (s *TokenGateway) ClassID() uint16               { return core.ClassIDTokenGateway }
func (s *TokenGateway) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *TokenGateway) GetVersion() uint8             { return s.Version }

func (s *TokenGateway) Access(attr int) core.CosemAttributeAccess {
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
func (s *TokenGateway) MarshalBinary() ([]byte, error) {
	return core.StructureData{
		core.OctetStringData(s.Token),
		core.OctetStringData(s.TokenTime),
		core.UnsignedIntegerData(s.TokenStatus),
	}.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *TokenGateway) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok {
		return nil
	}
	if 0 < len(st) {
		if v, ok := st[0].(core.OctetStringData); ok { s.Token = []byte(v) }
	}
	if 1 < len(st) {
		if v, ok := st[1].(core.OctetStringData); ok { s.TokenTime = []byte(v) }
	}
	if 2 < len(st) {
		if v, ok := st[2].(core.UnsignedIntegerData); ok { s.TokenStatus = uint8(v) }
	}
	return nil
}


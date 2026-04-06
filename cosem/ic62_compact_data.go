package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// CompactData is the COSEM Compact Data interface class (IC 62).
type CompactData struct {
	LogicalName  core.ObisCode
	Buffer       []byte
	CaptureTime  core.CosemDateTime
	CapturePeriod uint32
	Version      uint8
}

func (cd *CompactData) ClassID() uint16               { return core.ClassIDCompactData }
func (cd *CompactData) GetLogicalName() core.ObisCode { return cd.LogicalName }
func (cd *CompactData) GetVersion() uint8             { return cd.Version }
func (cd *CompactData) Access(attr int) core.CosemAttributeAccess {
	switch attr {
	case 1, 2:
		return core.CosemAttributeAccess{Access: core.AccessRead}
	case 3:
		return core.CosemAttributeAccess{Access: core.AccessReadWrite}
	default:
		return core.CosemAttributeAccess{Access: core.AccessNone}
	}
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (cd *CompactData) MarshalBinary() ([]byte, error) {
	return core.OctetStringData(cd.Buffer).ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (cd *CompactData) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	if v, ok := elem.(core.OctetStringData); ok {
		cd.Buffer = v
	}
	return nil
}

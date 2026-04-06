package cosem

import (
	"fmt"

	"github.com/ViewWay/dlms-cosem-go/core"
)

// ProfileGeneric is the COSEM Profile Generic interface class (IC 5).
type ProfileGeneric struct {
	LogicalName    core.ObisCode
	Buffer         []core.StructureData
	CaptureObjects []core.CosemObject
	CapturePeriod  uint16
	EntriesInUse   uint32
	Entries        uint32
	SortMethod     uint8
	SortObject     core.ObisCode
	Version        uint8
}

func (p *ProfileGeneric) ClassID() uint16               { return core.ClassIDProfileGeneric }
func (p *ProfileGeneric) GetLogicalName() core.ObisCode { return p.LogicalName }
func (p *ProfileGeneric) GetVersion() uint8             { return p.Version }
func (p *ProfileGeneric) Access(attr int) core.CosemAttributeAccess {
	switch attr {
	case 2, 3, 4, 7:
		return core.CosemAttributeAccess{Access: core.AccessRead}
	default:
		return core.CosemAttributeAccess{Access: core.AccessNone}
	}
}

func (p *ProfileGeneric) MarshalBinary() ([]byte, error) {
	arr := core.ArrayData{}
	for _, s := range p.Buffer {
		arr = append(arr, s)
	}
	return arr.ToBytes(), nil
}

func (p *ProfileGeneric) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	arr, ok := elem.(core.ArrayData)
	if !ok {
		return fmt.Errorf("profile_generic: expected array")
	}
	p.Buffer = make([]core.StructureData, len(arr))
	for i, item := range arr {
		s, ok := item.(core.StructureData)
		if !ok {
			return fmt.Errorf("profile_generic: expected structure at index %d", i)
		}
		p.Buffer[i] = s
	}
	return nil
}

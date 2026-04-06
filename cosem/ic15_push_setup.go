package cosem

import (
	"fmt"

	"github.com/ViewWay/dlms-cosem-go/core"
)

// PushObject represents an object in push list
type PushObject struct {
	ClassID     uint16
	LogicalName core.ObisCode
	Attribute   uint8
	DataIndex   uint8
}

// PushSetup is the COSEM Push Setup interface class (IC 15).
type PushSetup struct {
	LogicalName                core.ObisCode
	PushObjectList             []PushObject
	Service                    uint16
	Destination                []byte
	CommunicationWindow        core.CosemDateTime
	RandomisationStartInterval uint16
	NumberOfRetries            uint8
	RepetitionDelay            uint16
	Version                    uint8
}

func (ps *PushSetup) ClassID() uint16               { return core.ClassIDPush }
func (ps *PushSetup) GetLogicalName() core.ObisCode { return ps.LogicalName }
func (ps *PushSetup) GetVersion() uint8             { return ps.Version }
func (ps *PushSetup) Access(attr int) core.CosemAttributeAccess {
	switch attr {
	case 1:
		return core.CosemAttributeAccess{Access: core.AccessRead}
	case 2, 3, 4, 5, 6, 7, 8:
		return core.CosemAttributeAccess{Access: core.AccessReadWrite}
	default:
		return core.CosemAttributeAccess{Access: core.AccessNone}
	}
}

func (ps *PushSetup) MarshalBinary() ([]byte, error) {
	arr := core.ArrayData{}
	for _, o := range ps.PushObjectList {
		arr = append(arr, core.StructureData{
			core.UnsignedLongData(o.ClassID),
			core.OctetStringData(o.LogicalName[:]),
			core.IntegerData(int8(o.Attribute)),
			core.UnsignedIntegerData(o.DataIndex),
		})
	}
	return arr.ToBytes(), nil
}

func (ps *PushSetup) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	arr, ok := elem.(core.ArrayData)
	if !ok {
		return fmt.Errorf("push_setup: expected array")
	}
	ps.PushObjectList = make([]PushObject, len(arr))
	for i, item := range arr {
		st, ok := item.(core.StructureData)
		if !ok || len(st) < 4 {
			return fmt.Errorf("push_setup: invalid entry at index %d", i)
		}
		ps.PushObjectList[i].ClassID = uint16(st[0].(core.UnsignedLongData))
		ps.PushObjectList[i].LogicalName = core.ObisCode(st[1].(core.OctetStringData))
		ps.PushObjectList[i].Attribute = uint8(st[2].(core.IntegerData))
		ps.PushObjectList[i].DataIndex = uint8(st[3].(core.UnsignedIntegerData))
	}
	return nil
}

func (ps *PushSetup) AddObject(obj PushObject) {
	ps.PushObjectList = append(ps.PushObjectList, obj)
}

func (ps *PushSetup) ObjectCount() int {
	return len(ps.PushObjectList)
}

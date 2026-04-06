package cosem

import (
	"fmt"

	"github.com/ViewWay/dlms-cosem-go/core"
)

// AssociationSN is the COSEM Association SN interface class (IC 10).
type AssociationSN struct {
	LogicalName        core.ObisCode
	ObjectList         []AssociationSNObject
	AccessRightsList   []AssociationSNAccessRight
	ClientSAP          uint16
	ServerSAP          uint16
	LLSSecret          []byte
	Version            uint8
}

// AssociationSNObject represents an object in the SN object list.
type AssociationSNObject struct {
	BaseAddress uint16
	ClassID     uint16
	LogicalName core.ObisCode
	Version     uint8
}

// AssociationSNAccessRight represents access rights for an object.
type AssociationSNAccessRight struct {
	BaseAddress   uint16
	AccessRight   uint16
}

func (a *AssociationSN) ClassID() uint16               { return core.ClassIDAssociationSN }
func (a *AssociationSN) GetLogicalName() core.ObisCode { return a.LogicalName }
func (a *AssociationSN) GetVersion() uint8             { return a.Version }
func (a *AssociationSN) Access(attr int) core.CosemAttributeAccess {
	switch attr {
	case 1, 2, 3, 4, 5:
		return core.CosemAttributeAccess{Access: core.AccessRead}
	default:
		return core.CosemAttributeAccess{Access: core.AccessNone}
	}
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (a *AssociationSN) MarshalBinary() ([]byte, error) {
	arr := core.ArrayData{}
	for _, o := range a.ObjectList {
		arr = append(arr, core.StructureData{
			core.UnsignedLongData(o.BaseAddress),
			core.UnsignedLongData(o.ClassID),
			core.OctetStringData(o.LogicalName[:]),
			core.UnsignedIntegerData(o.Version),
		})
	}
	return arr.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (a *AssociationSN) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	arr, ok := elem.(core.ArrayData)
	if !ok {
		return fmt.Errorf("association_sn: expected array")
	}
	a.ObjectList = make([]AssociationSNObject, len(arr))
	for i, item := range arr {
		st, ok := item.(core.StructureData)
		if !ok || len(st) < 4 {
			return fmt.Errorf("association_sn: invalid entry at index %d", i)
		}
		a.ObjectList[i].BaseAddress = uint16(st[0].(core.UnsignedLongData))
		a.ObjectList[i].ClassID = uint16(st[1].(core.UnsignedLongData))
		a.ObjectList[i].LogicalName = core.ObisCode(st[2].(core.OctetStringData))
		a.ObjectList[i].Version = uint8(st[3].(core.UnsignedIntegerData))
	}
	return nil
}

// AddObject adds an object to the association.
func (a *AssociationSN) AddObject(obj AssociationSNObject) {
	a.ObjectList = append(a.ObjectList, obj)
}

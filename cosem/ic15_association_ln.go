package cosem

import (
	"fmt"

	"github.com/ViewWay/dlms-cosem-go/core"
)

// ObjectListEntry represents an entry in the object list
type ObjectListEntry struct {
	ClassID     uint16
	LogicalName core.ObisCode
	Version     uint8
}

// AssociationLN is the COSEM Association LN interface class (IC 11).
type AssociationLN struct {
	LogicalName            core.ObisCode
	ObjectList             []ObjectListEntry
	ClientSAP              uint16
	ServerSAP              uint16
	ApplicationContextName string
	AuthenticationMechanism string
	LLSSecret              []byte
	HLSSecret              []byte
	AuthenticationStatus   bool
	Version                uint8
}

func (a *AssociationLN) ClassID() uint16               { return core.ClassIDAssociationLN }
func (a *AssociationLN) GetLogicalName() core.ObisCode { return a.LogicalName }
func (a *AssociationLN) GetVersion() uint8             { return a.Version }
func (a *AssociationLN) Access(attr int) core.CosemAttributeAccess {
	switch attr {
	case 2, 3, 4, 5, 6, 7, 8:
		return core.CosemAttributeAccess{Access: core.AccessRead}
	default:
		return core.CosemAttributeAccess{Access: core.AccessNone}
	}
}

func (a *AssociationLN) MarshalBinary() ([]byte, error) {
	arr := core.ArrayData{}
	for _, o := range a.ObjectList {
		arr = append(arr, core.StructureData{
			core.UnsignedLongData(o.ClassID),
			core.OctetStringData(o.LogicalName[:]),
			core.UnsignedIntegerData(o.Version),
		})
	}
	return arr.ToBytes(), nil
}

func (a *AssociationLN) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	arr, ok := elem.(core.ArrayData)
	if !ok {
		return fmt.Errorf("association_ln: expected array")
	}
	a.ObjectList = make([]ObjectListEntry, len(arr))
	for i, item := range arr {
		st, ok := item.(core.StructureData)
		if !ok || len(st) < 3 {
			return fmt.Errorf("association_ln: invalid entry at index %d", i)
		}
		a.ObjectList[i].ClassID = uint16(st[0].(core.UnsignedLongData))
		a.ObjectList[i].LogicalName = core.ObisCode(st[1].(core.OctetStringData))
		a.ObjectList[i].Version = uint8(st[2].(core.UnsignedIntegerData))
	}
	return nil
}

func (a *AssociationLN) AddObject(entry ObjectListEntry) {
	a.ObjectList = append(a.ObjectList, entry)
}

func (a *AssociationLN) ObjectCount() int {
	return len(a.ObjectList)
}

func (a *AssociationLN) SetAuthenticated(status bool) {
	a.AuthenticationStatus = status
}

func (a *AssociationLN) IsAuthenticated() bool {
	return a.AuthenticationStatus
}

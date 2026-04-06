package cosem

import (
	"fmt"

	"github.com/ViewWay/dlms-cosem-go/core"
)

// RegisterAssignment represents a register assignment with mask list.
type RegisterAssignment struct {
	RegisterReference core.ObisCode
	MaskList          []core.DlmsData
}

// RegisterActivation is the COSEM Register Activation interface class (IC 14).
type RegisterActivation struct {
	LogicalName        core.ObisCode
	RegisterAssignments []RegisterAssignment
	Version            uint8
}

func (ra *RegisterActivation) ClassID() uint16               { return core.ClassIDRegisterActivation }
func (ra *RegisterActivation) GetLogicalName() core.ObisCode { return ra.LogicalName }
func (ra *RegisterActivation) GetVersion() uint8             { return ra.Version }
func (ra *RegisterActivation) Access(attr int) core.CosemAttributeAccess {
	switch attr {
	case 1, 2:
		return core.CosemAttributeAccess{Access: core.AccessRead}
	default:
		return core.CosemAttributeAccess{Access: core.AccessNone}
	}
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (ra *RegisterActivation) MarshalBinary() ([]byte, error) {
	arr := core.ArrayData{}
	for _, a := range ra.RegisterAssignments {
		maskArr := core.ArrayData{}
		for _, m := range a.MaskList {
			maskArr = append(maskArr, m)
		}
		arr = append(arr, core.StructureData{
			core.OctetStringData(a.RegisterReference[:]),
			maskArr,
		})
	}
	return arr.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (ra *RegisterActivation) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	arr, ok := elem.(core.ArrayData)
	if !ok {
		return fmt.Errorf("register_activation: expected array")
	}
	ra.RegisterAssignments = make([]RegisterAssignment, len(arr))
	for i, item := range arr {
		st, ok := item.(core.StructureData)
		if !ok || len(st) < 2 {
			return fmt.Errorf("register_activation: invalid entry at index %d", i)
		}
		ra.RegisterAssignments[i].RegisterReference = core.ObisCode(st[0].(core.OctetStringData))
		if maskArr, ok := st[1].(core.ArrayData); ok {
			ra.RegisterAssignments[i].MaskList = []core.DlmsData(maskArr)
		}
	}
	return nil
}

// AddRegister adds a register assignment.
func (ra *RegisterActivation) AddRegister(assignment RegisterAssignment) {
	ra.RegisterAssignments = append(ra.RegisterAssignments, assignment)
}

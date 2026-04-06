package cosem

import (
	"fmt"

	"github.com/ViewWay/dlms-cosem-go/core"
)

// Script represents a single script entry in ScriptTable
type Script struct {
	ScriptID       uint8
	ScriptSelector uint8
	FileID         uint16
}

// ScriptTable is the COSEM Script Table interface class (IC 7).
type ScriptTable struct {
	LogicalName core.ObisCode
	Scripts     []Script
	Version     uint8
}

func (st *ScriptTable) ClassID() uint16               { return core.ClassIDScriptTable }
func (st *ScriptTable) GetLogicalName() core.ObisCode { return st.LogicalName }
func (st *ScriptTable) GetVersion() uint8             { return st.Version }
func (st *ScriptTable) Access(attr int) core.CosemAttributeAccess {
	switch attr {
	case 1, 2:
		return core.CosemAttributeAccess{Access: core.AccessRead}
	default:
		return core.CosemAttributeAccess{Access: core.AccessNone}
	}
}

func (st *ScriptTable) MarshalBinary() ([]byte, error) {
	arr := core.ArrayData{}
	for _, s := range st.Scripts {
		arr = append(arr, core.StructureData{
			core.UnsignedIntegerData(s.ScriptID),
			core.UnsignedIntegerData(s.ScriptSelector),
			core.UnsignedLongData(s.FileID),
		})
	}
	return arr.ToBytes(), nil
}

func (st *ScriptTable) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	arr, ok := elem.(core.ArrayData)
	if !ok {
		return fmt.Errorf("script_table: expected array")
	}
	st.Scripts = make([]Script, len(arr))
	for i, item := range arr {
		s, ok := item.(core.StructureData)
		if !ok || len(s) < 3 {
			return fmt.Errorf("script_table: invalid script entry at index %d", i)
		}
		st.Scripts[i].ScriptID = uint8(s[0].(core.UnsignedIntegerData))
		st.Scripts[i].ScriptSelector = uint8(s[1].(core.UnsignedIntegerData))
		st.Scripts[i].FileID = uint16(s[2].(core.UnsignedLongData))
	}
	return nil
}

func (st *ScriptTable) AddScript(script Script) {
	st.Scripts = append(st.Scripts, script)
}

func (st *ScriptTable) RemoveScript(scriptID uint8) *Script {
	for i, s := range st.Scripts {
		if s.ScriptID == scriptID {
			st.Scripts = append(st.Scripts[:i], st.Scripts[i+1:]...)
			return &s
		}
	}
	return nil
}

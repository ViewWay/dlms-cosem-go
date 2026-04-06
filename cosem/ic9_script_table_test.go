package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestScriptTable_ClassID(t *testing.T) {
	st := &ScriptTable{}
	if st.ClassID() != 9 {
		t.Errorf("ClassID() = %d, want 9", st.ClassID())
	}
	if st.ClassID() != core.ClassIDScriptTable {
		t.Error("ClassID mismatch with const")
	}
}

func TestScriptTable_New(t *testing.T) {
	st := &ScriptTable{
		LogicalName: core.ObisCode{0, 0, 10, 0, 0, 255},
	}
	if st.Scripts == nil {
		t.Log("Scripts nil by default - ok")
	}
}

func TestScriptTable_MarshalBinary(t *testing.T) {
	st := &ScriptTable{Scripts: []Script{{ScriptID: 1, ScriptSelector: 0}}}
	b, err := st.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	st2 := &ScriptTable{}
	if err := st2.UnmarshalBinary(b); err != nil {
		t.Fatal(err)
	}
}

func TestScriptTable_Fields(t *testing.T) {
	st := &ScriptTable{
		LogicalName: core.ObisCode{0, 0, 10, 0, 0, 255},
		Version:     1,
	}
	if st.GetVersion() != 1 {
		t.Error("GetVersion mismatch")
	}
}

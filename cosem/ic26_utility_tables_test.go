package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestUtilityTables_ClassID(t *testing.T) {
	s := &UtilityTables{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDUtilityTables {
		t.Errorf("expected %d, got %d", core.ClassIDUtilityTables, s.ClassID())
	}
}

func TestUtilityTables_MarshalUnmarshalBinary(t *testing.T) {
	s := &UtilityTables{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &UtilityTables{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

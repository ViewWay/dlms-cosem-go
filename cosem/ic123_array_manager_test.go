package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestArrayManager_ClassID(t *testing.T) {
	s := &ArrayManager{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDArrayManager {
		t.Errorf("expected %d, got %d", core.ClassIDArrayManager, s.ClassID())
	}
}

func TestArrayManager_MarshalUnmarshalBinary(t *testing.T) {
	s := &ArrayManager{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &ArrayManager{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

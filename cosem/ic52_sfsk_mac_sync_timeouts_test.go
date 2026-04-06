package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

func TestSFSKMACSyncTimeouts_ClassID(t *testing.T) {
	s := &SFSKMACSyncTimeouts{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	if s.ClassID() != core.ClassIDSFSKMACSyncTimeouts {
		t.Errorf("expected %d, got %d", core.ClassIDSFSKMACSyncTimeouts, s.ClassID())
	}
}

func TestSFSKMACSyncTimeouts_MarshalUnmarshalBinary(t *testing.T) {
	s := &SFSKMACSyncTimeouts{LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &SFSKMACSyncTimeouts{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
}

package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

// ============================================================================
// Tests for new IC classes
// ============================================================================

func TestAssociationSN_ClassID(t *testing.T) {
	a := &AssociationSN{LogicalName: core.ObisCode{0, 0, 40, 0, 1, 255}}
	if a.ClassID() != core.ClassIDAssociationSN {
		t.Errorf("expected %d, got %d", core.ClassIDAssociationSN, a.ClassID())
	}
}

func TestAssociationSN_AddObject(t *testing.T) {
	a := &AssociationSN{}
	a.AddObject(AssociationSNObject{BaseAddress: 0x1000, ClassID: 1})
	if len(a.ObjectList) != 1 {
		t.Fatalf("expected 1 object, got %d", len(a.ObjectList))
	}
}

func TestAssociationSN_MarshalBinary(t *testing.T) {
	a := &AssociationSN{LogicalName: core.ObisCode{0, 0, 40, 0, 1, 255}}
	a.AddObject(AssociationSNObject{BaseAddress: 0x1000, ClassID: 1, Version: 1})
	data, err := a.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty data")
	}
}

func TestRegisterActivation_ClassID(t *testing.T) {
	ra := &RegisterActivation{LogicalName: core.ObisCode{0, 0, 96, 1, 0, 255}}
	if ra.ClassID() != core.ClassIDRegisterActivation {
		t.Errorf("expected %d, got %d", core.ClassIDRegisterActivation, ra.ClassID())
	}
}

func TestRegisterActivation_AddRegister(t *testing.T) {
	ra := &RegisterActivation{}
	ra.AddRegister(RegisterAssignment{RegisterReference: core.ObisCode{0, 0, 1, 0, 0, 255}})
	if len(ra.RegisterAssignments) != 1 {
		t.Fatalf("expected 1 assignment, got %d", len(ra.RegisterAssignments))
	}
}

func TestRegisterActivation_MarshalBinary(t *testing.T) {
	ra := &RegisterActivation{}
	ra.AddRegister(RegisterAssignment{
		RegisterReference: core.ObisCode{0, 0, 1, 0, 0, 255},
		MaskList:          []core.DlmsData{core.UnsignedIntegerData(1)},
	})
	data, err := ra.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty data")
	}
}

func TestAccount_ClassID(t *testing.T) {
	a := &Account{LogicalName: core.ObisCode{0, 0, 96, 0, 0, 255}}
	if a.ClassID() != core.ClassIDAccount {
		t.Errorf("expected %d, got %d", core.ClassIDAccount, a.ClassID())
	}
}

func TestAccount_MarshalUnmarshalBinary(t *testing.T) {
	a := &Account{VendorInfo: []byte{0x01, 0x02}}
	data, err := a.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	a2 := &Account{}
	if err := a2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
	if string(a2.VendorInfo) != string(a.VendorInfo) {
		t.Error("vendor info mismatch")
	}
}

func TestValueDisplay_ClassID(t *testing.T) {
	vd := &ValueDisplay{LogicalName: core.ObisCode{0, 0, 30, 0, 0, 255}}
	if vd.ClassID() != core.ClassIDValueDisplay {
		t.Errorf("expected %d, got %d", core.ClassIDValueDisplay, vd.ClassID())
	}
}

func TestValueDisplay_MarshalUnmarshalBinary(t *testing.T) {
	vd := &ValueDisplay{ValueToDisplay: core.DoubleLongData(42)}
	data, err := vd.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	vd2 := &ValueDisplay{}
	if err := vd2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
	if vd2.ValueToDisplay == nil {
		t.Error("expected non-nil value")
	}
}

func TestIPv4Setup_ClassID(t *testing.T) {
	s := &IPv4Setup{LogicalName: core.ObisCode{0, 0, 40, 0, 0, 255}}
	if s.ClassID() != core.ClassIDIPv4Setup {
		t.Errorf("expected %d, got %d", core.ClassIDIPv4Setup, s.ClassID())
	}
}

func TestIPv4Setup_SetIPAddress(t *testing.T) {
	s := &IPv4Setup{}
	s.SetIPAddress([4]byte{192, 168, 1, 100})
	if s.IPAddress != [4]byte{192, 168, 1, 100} {
		t.Error("IP address not set correctly")
	}
}

func TestIPv4Setup_MarshalUnmarshalBinary(t *testing.T) {
	s := &IPv4Setup{IPAddress: [4]byte{192, 168, 1, 100}}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	s2 := &IPv4Setup{}
	if err := s2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
	if s2.IPAddress != [4]byte{192, 168, 1, 100} {
		t.Error("IP address mismatch after unmarshal")
	}
}

func TestTCPUDPSetup_ClassID(t *testing.T) {
	s := &TCPUDPSetup{LogicalName: core.ObisCode{0, 0, 41, 0, 0, 255}}
	if s.ClassID() != core.ClassIDIPv4TCPSetup {
		t.Errorf("expected %d, got %d", core.ClassIDIPv4TCPSetup, s.ClassID())
	}
}

func TestTCPUDPSetup_MarshalBinary(t *testing.T) {
	s := &TCPUDPSetup{
		TCPConnections: []TCPUDPSetupConnection{
			{RemoteAddress: [4]byte{10, 0, 0, 1}, RemotePort: 4059},
		},
	}
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty data")
	}
}

func TestPPPSetup_ClassID(t *testing.T) {
	p := &PPPSetup{LogicalName: core.ObisCode{0, 0, 45, 0, 0, 255}}
	if p.ClassID() != core.ClassIDPPPSetup {
		t.Errorf("expected %d, got %d", core.ClassIDPPPSetup, p.ClassID())
	}
}

func TestPPPSetup_MarshalUnmarshalBinary(t *testing.T) {
	p := &PPPSetup{UserName: "admin"}
	data, err := p.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	p2 := &PPPSetup{}
	if err := p2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
	if p2.UserName != "admin" {
		t.Error("username mismatch")
	}
}

func TestStatusMapping_ClassID(t *testing.T) {
	sm := &StatusMapping{LogicalName: core.ObisCode{0, 0, 55, 0, 0, 255}}
	if sm.ClassID() != core.ClassIDStatusMapping {
		t.Errorf("expected %d, got %d", core.ClassIDStatusMapping, sm.ClassID())
	}
}

func TestStatusMapping_BitOperations(t *testing.T) {
	sm := &StatusMapping{}
	sm.SetBit(0, true)
	if !sm.GetBitStatus(0) {
		t.Error("bit 0 should be set")
	}
	sm.SetBit(0, false)
	if sm.GetBitStatus(0) {
		t.Error("bit 0 should be cleared")
	}
}

func TestStatusMapping_MarshalUnmarshalBinary(t *testing.T) {
	sm := &StatusMapping{StatusWord: 0xDEADBEEF}
	data, err := sm.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	sm2 := &StatusMapping{}
	if err := sm2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
	if sm2.StatusWord != 0xDEADBEEF {
		t.Errorf("expected 0xDEADBEEF, got 0x%X", sm2.StatusWord)
	}
}

func TestCompactData_ClassID(t *testing.T) {
	cd := &CompactData{LogicalName: core.ObisCode{0, 0, 62, 0, 0, 255}}
	if cd.ClassID() != core.ClassIDCompactData {
		t.Errorf("expected %d, got %d", core.ClassIDCompactData, cd.ClassID())
	}
}

func TestCompactData_MarshalUnmarshalBinary(t *testing.T) {
	cd := &CompactData{Buffer: []byte{0xDE, 0xAD, 0xBE, 0xEF}}
	data, err := cd.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	cd2 := &CompactData{}
	if err := cd2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
	if string(cd2.Buffer) != string(cd.Buffer) {
		t.Error("buffer mismatch")
	}
}

func TestMBusClient_ClassID(t *testing.T) {
	m := &MBusClient{LogicalName: core.ObisCode{0, 0, 26, 0, 0, 255}}
	if m.ClassID() != core.ClassIDMBusClient {
		t.Errorf("expected %d, got %d", core.ClassIDMBusClient, m.ClassID())
	}
}

func TestMBusClient_MarshalUnmarshalBinary(t *testing.T) {
	m := &MBusClient{PrimaryAddress: 42}
	data, err := m.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	m2 := &MBusClient{}
	if err := m2.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}
	if m2.PrimaryAddress != 42 {
		t.Errorf("expected 42, got %d", m2.PrimaryAddress)
	}
}

// Split IC class tests are in cosem_objects_test.go

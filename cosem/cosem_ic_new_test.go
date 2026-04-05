package cosem

import (
	"testing"

	"github.com/ViewWay/dlms-cosem-go/core"
)

// ============================================================================
// IC10 Script Table Tests
// ============================================================================

func TestScriptTableClassID(t *testing.T) {
	st := &ScriptTable{LogicalName: core.ObisCode{0, 0, 10, 0, 0, 255}}
	if st.ClassID() != core.ClassIDScriptTable {
		t.Errorf("expected ClassID %d, got %d", core.ClassIDScriptTable, st.ClassID())
	}
}

func TestScriptTableAddScript(t *testing.T) {
	st := &ScriptTable{LogicalName: core.ObisCode{0, 0, 10, 0, 0, 255}}
	st.AddScript(Script{ScriptID: 1, ScriptSelector: 0, FileID: 100})
	if len(st.Scripts) != 1 {
		t.Errorf("expected 1 script, got %d", len(st.Scripts))
	}
}

func TestScriptTableRemoveScript(t *testing.T) {
	st := &ScriptTable{LogicalName: core.ObisCode{0, 0, 10, 0, 0, 255}}
	st.AddScript(Script{ScriptID: 1, ScriptSelector: 0, FileID: 100})
	removed := st.RemoveScript(1)
	if removed == nil {
		t.Error("expected script to be removed")
	}
	if removed.ScriptID != 1 {
		t.Errorf("expected ScriptID 1, got %d", removed.ScriptID)
	}
	if len(st.Scripts) != 0 {
		t.Errorf("expected 0 scripts, got %d", len(st.Scripts))
	}
}

func TestScriptTableMarshalBinary(t *testing.T) {
	st := &ScriptTable{LogicalName: core.ObisCode{0, 0, 10, 0, 0, 255}}
	st.AddScript(Script{ScriptID: 1, ScriptSelector: 0, FileID: 100})
	data, err := st.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary failed: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty data")
	}
}

func TestScriptTableAccess(t *testing.T) {
	st := &ScriptTable{LogicalName: core.ObisCode{0, 0, 10, 0, 0, 255}}
	if st.Access(1).Access != core.AccessRead {
		t.Error("expected read access for attr 1")
	}
	if st.Access(2).Access != core.AccessRead {
		t.Error("expected read access for attr 2")
	}
	if st.Access(99).Access != core.AccessNone {
		t.Error("expected no access for attr 99")
	}
}

// ============================================================================
// IC12 Schedule Tests
// ============================================================================

func TestScheduleClassID(t *testing.T) {
	s := &Schedule{LogicalName: core.ObisCode{0, 0, 12, 0, 0, 255}}
	if s.ClassID() != core.ClassIDSchedule {
		t.Errorf("expected ClassID %d, got %d", core.ClassIDSchedule, s.ClassID())
	}
}

func TestScheduleAddEntry(t *testing.T) {
	s := &Schedule{LogicalName: core.ObisCode{0, 0, 12, 0, 0, 255}}
	s.AddEntry(ScheduleEntry{
		Index:          1,
		Enable:         true,
		ScriptSelector: 1,
		ValidForDays:   0x7F,
	})
	if s.EntryCount() != 1 {
		t.Errorf("expected 1 entry, got %d", s.EntryCount())
	}
}

func TestScheduleMarshalBinary(t *testing.T) {
	s := &Schedule{LogicalName: core.ObisCode{0, 0, 12, 0, 0, 255}}
	s.AddEntry(ScheduleEntry{
		Index:          1,
		Enable:         true,
		ScriptSelector: 1,
		ValidForDays:   0x7F,
	})
	data, err := s.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary failed: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty data")
	}
}

func TestScheduleAccess(t *testing.T) {
	s := &Schedule{LogicalName: core.ObisCode{0, 0, 12, 0, 0, 255}}
	if s.Access(1).Access != core.AccessReadWrite {
		t.Error("expected read/write access for attr 1")
	}
}

// ============================================================================
// IC20 Activity Calendar Tests
// ============================================================================

func TestActivityCalendarClassID(t *testing.T) {
	ac := &ActivityCalendar{LogicalName: core.ObisCode{0, 0, 20, 0, 0, 255}}
	if ac.ClassID() != core.ClassIDActivityCalendar {
		t.Errorf("expected ClassID %d, got %d", core.ClassIDActivityCalendar, ac.ClassID())
	}
}

func TestActivityCalendarSetCalendarName(t *testing.T) {
	ac := &ActivityCalendar{LogicalName: core.ObisCode{0, 0, 20, 0, 0, 255}}
	ac.SetCalendarName("Summer")
	if ac.CalendarName != "Summer" {
		t.Errorf("expected calendar name 'Summer', got '%s'", ac.CalendarName)
	}
}

func TestActivityCalendarMarshalBinary(t *testing.T) {
	ac := &ActivityCalendar{LogicalName: core.ObisCode{0, 0, 20, 0, 0, 255}}
	ac.CalendarName = "Test"
	data, err := ac.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary failed: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty data")
	}
}

func TestActivityCalendarAccess(t *testing.T) {
	ac := &ActivityCalendar{LogicalName: core.ObisCode{0, 0, 20, 0, 0, 255}}
	if ac.Access(1).Access != core.AccessRead {
		t.Error("expected read access for attr 1")
	}
	if ac.Access(2).Access != core.AccessReadWrite {
		t.Error("expected read/write access for attr 2")
	}
}

// ============================================================================
// IC21 Register Monitor Tests
// ============================================================================

func TestRegisterMonitorClassID(t *testing.T) {
	rm := &RegisterMonitor{LogicalName: core.ObisCode{0, 0, 21, 0, 0, 255}}
	if rm.ClassID() != core.ClassIDRegisterMonitor {
		t.Errorf("expected ClassID %d, got %d", core.ClassIDRegisterMonitor, rm.ClassID())
	}
}

func TestRegisterMonitorAddThreshold(t *testing.T) {
	rm := &RegisterMonitor{LogicalName: core.ObisCode{0, 0, 21, 0, 0, 255}}
	rm.AddThreshold(core.DoubleLongData(100))
	if rm.ThresholdCount() != 1 {
		t.Errorf("expected 1 threshold, got %d", rm.ThresholdCount())
	}
}

func TestRegisterMonitorMarshalBinary(t *testing.T) {
	rm := &RegisterMonitor{LogicalName: core.ObisCode{0, 0, 21, 0, 0, 255}}
	rm.AddThreshold(core.DoubleLongData(100))
	data, err := rm.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary failed: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty data")
	}
}

func TestRegisterMonitorAccess(t *testing.T) {
	rm := &RegisterMonitor{LogicalName: core.ObisCode{0, 0, 21, 0, 0, 255}}
	if rm.Access(2).Access != core.AccessRead {
		t.Error("expected read access for attr 2")
	}
}

// ============================================================================
// IC22 Single Action Schedule Tests
// ============================================================================

func TestSingleActionScheduleClassID(t *testing.T) {
	sas := &SingleActionSchedule{LogicalName: core.ObisCode{0, 0, 22, 0, 0, 255}}
	if sas.ClassID() != core.ClassIDSingleActionSchedule {
		t.Errorf("expected ClassID %d, got %d", core.ClassIDSingleActionSchedule, sas.ClassID())
	}
}

func TestSingleActionScheduleAddEntry(t *testing.T) {
	sas := &SingleActionSchedule{LogicalName: core.ObisCode{0, 0, 22, 0, 0, 255}}
	sas.AddEntry(ActionScheduleEntry{
		ExecutedScriptLN: core.ObisCode{0, 0, 10, 0, 0, 255},
	})
	if sas.EntryCount() != 1 {
		t.Errorf("expected 1 entry, got %d", sas.EntryCount())
	}
}

func TestSingleActionScheduleMarshalBinary(t *testing.T) {
	sas := &SingleActionSchedule{LogicalName: core.ObisCode{0, 0, 22, 0, 0, 255}}
	sas.AddEntry(ActionScheduleEntry{
		ExecutedScriptLN: core.ObisCode{0, 0, 10, 0, 0, 255},
	})
	data, err := sas.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary failed: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty data")
	}
}

func TestSingleActionScheduleAccess(t *testing.T) {
	sas := &SingleActionSchedule{LogicalName: core.ObisCode{0, 0, 22, 0, 0, 255}}
	if sas.Access(1).Access != core.AccessRead {
		t.Error("expected read access for attr 1")
	}
	if sas.Access(2).Access != core.AccessReadWrite {
		t.Error("expected read/write access for attr 2")
	}
}

// ============================================================================
// IC29 Value Table Tests
// ============================================================================

func TestValueTableClassID(t *testing.T) {
	vt := &ValueTable{LogicalName: core.ObisCode{0, 0, 29, 0, 0, 255}}
	if vt.ClassID() != core.ClassIDValueTable {
		t.Errorf("expected ClassID %d, got %d", core.ClassIDValueTable, vt.ClassID())
	}
}

func TestValueTableAddValue(t *testing.T) {
	vt := &ValueTable{LogicalName: core.ObisCode{0, 0, 29, 0, 0, 255}}
	vt.AddValue(ValueEntry{Index: 1, Value: core.DoubleLongData(100)})
	if vt.ValueCount() != 1 {
		t.Errorf("expected 1 value, got %d", vt.ValueCount())
	}
}

func TestValueTableAddDescriptor(t *testing.T) {
	vt := &ValueTable{LogicalName: core.ObisCode{0, 0, 29, 0, 0, 255}}
	vt.AddDescriptor(ValueDescriptor{Index: 1, Description: "Voltage", Unit: 27})
	if vt.DescriptorCount() != 1 {
		t.Errorf("expected 1 descriptor, got %d", vt.DescriptorCount())
	}
}

func TestValueTableMarshalBinary(t *testing.T) {
	vt := &ValueTable{LogicalName: core.ObisCode{0, 0, 29, 0, 0, 255}}
	vt.AddValue(ValueEntry{Index: 1, Value: core.DoubleLongData(100)})
	data, err := vt.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary failed: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty data")
	}
}

func TestValueTableAccess(t *testing.T) {
	vt := &ValueTable{LogicalName: core.ObisCode{0, 0, 29, 0, 0, 255}}
	if vt.Access(1).Access != core.AccessRead {
		t.Error("expected read access for attr 1")
	}
	if vt.Access(2).Access != core.AccessRead {
		t.Error("expected read access for attr 2")
	}
}

// ============================================================================
// IC70 Disconnect Control Tests
// ============================================================================

func TestDisconnectControlClassID(t *testing.T) {
	dc := &DisconnectControl{LogicalName: core.ObisCode{0, 0, 70, 0, 0, 255}}
	if dc.ClassID() != core.ClassIDDisconnectControl {
		t.Errorf("expected ClassID %d, got %d", core.ClassIDDisconnectControl, dc.ClassID())
	}
}

func TestDisconnectControlDisconnect(t *testing.T) {
	dc := &DisconnectControl{LogicalName: core.ObisCode{0, 0, 70, 0, 0, 255}}
	if err := dc.Disconnect(); err != nil {
		t.Fatalf("Disconnect failed: %v", err)
	}
	if dc.ControlState != DisconnectStateDisconnected {
		t.Errorf("expected state %d, got %d", DisconnectStateDisconnected, dc.ControlState)
	}
}

func TestDisconnectControlReconnect(t *testing.T) {
	dc := &DisconnectControl{LogicalName: core.ObisCode{0, 0, 70, 0, 0, 255}}
	if err := dc.Reconnect(); err != nil {
		t.Fatalf("Reconnect failed: %v", err)
	}
	if dc.ControlState != DisconnectStateConnected {
		t.Errorf("expected state %d, got %d", DisconnectStateConnected, dc.ControlState)
	}
}

func TestDisconnectControlArm(t *testing.T) {
	dc := &DisconnectControl{LogicalName: core.ObisCode{0, 0, 70, 0, 0, 255}}
	if err := dc.Arm(); err != nil {
		t.Fatalf("Arm failed: %v", err)
	}
	if dc.ControlState != DisconnectStateArmed {
		t.Errorf("expected state %d, got %d", DisconnectStateArmed, dc.ControlState)
	}
}

func TestDisconnectControlMarshalBinary(t *testing.T) {
	dc := &DisconnectControl{LogicalName: core.ObisCode{0, 0, 70, 0, 0, 255}}
	data, err := dc.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary failed: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty data")
	}
}

func TestDisconnectControlAccess(t *testing.T) {
	dc := &DisconnectControl{LogicalName: core.ObisCode{0, 0, 70, 0, 0, 255}}
	if dc.Access(2).Access != core.AccessRead {
		t.Error("expected read access for attr 2")
	}
}

// ============================================================================
// IC18 Image Transfer Tests
// ============================================================================

func TestImageTransferClassID(t *testing.T) {
	it := &ImageTransfer{LogicalName: core.ObisCode{0, 0, 44, 0, 0, 255}}
	if it.ClassID() != core.ClassIDImageTransfer {
		t.Errorf("expected ClassID %d, got %d", core.ClassIDImageTransfer, it.ClassID())
	}
}

func TestImageTransferInitiateImage(t *testing.T) {
	it := &ImageTransfer{LogicalName: core.ObisCode{0, 0, 44, 0, 0, 255}}
	it.InitiateImage(200, 100)
	if it.ImageBlockSize != 200 {
		t.Errorf("expected block size 200, got %d", it.ImageBlockSize)
	}
	if it.TransferStatus != ImageTransferStatusInitiated {
		t.Errorf("expected status %d, got %d", ImageTransferStatusInitiated, it.TransferStatus)
	}
}

func TestImageTransferVerifyImage(t *testing.T) {
	it := &ImageTransfer{LogicalName: core.ObisCode{0, 0, 44, 0, 0, 255}}
	if !it.VerifyImage() {
		t.Error("expected VerifyImage to return true")
	}
	if it.TransferStatus != ImageTransferStatusVerificationSuccessful {
		t.Errorf("expected status %d, got %d", ImageTransferStatusVerificationSuccessful, it.TransferStatus)
	}
}

func TestImageTransferActivateImage(t *testing.T) {
	it := &ImageTransfer{LogicalName: core.ObisCode{0, 0, 44, 0, 0, 255}}
	if err := it.ActivateImage(); err != nil {
		t.Fatalf("ActivateImage failed: %v", err)
	}
	if it.TransferStatus != ImageTransferStatusImageActivated {
		t.Errorf("expected status %d, got %d", ImageTransferStatusImageActivated, it.TransferStatus)
	}
}

func TestImageTransferMarshalBinary(t *testing.T) {
	it := &ImageTransfer{LogicalName: core.ObisCode{0, 0, 44, 0, 0, 255}}
	data, err := it.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary failed: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty data")
	}
}

func TestImageTransferAccess(t *testing.T) {
	it := &ImageTransfer{LogicalName: core.ObisCode{0, 0, 44, 0, 0, 255}}
	if it.Access(2).Access != core.AccessRead {
		t.Error("expected read access for attr 2")
	}
}

// ============================================================================
// IC40 Push Setup Tests
// ============================================================================

func TestPushSetupClassID(t *testing.T) {
	ps := &PushSetup{LogicalName: core.ObisCode{0, 0, 40, 0, 0, 255}}
	if ps.ClassID() != core.ClassIDPush {
		t.Errorf("expected ClassID %d, got %d", core.ClassIDPush, ps.ClassID())
	}
}

func TestPushSetupAddObject(t *testing.T) {
	ps := &PushSetup{LogicalName: core.ObisCode{0, 0, 40, 0, 0, 255}}
	ps.AddObject(PushObject{ClassID: 8, LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}, Attribute: 2})
	if ps.ObjectCount() != 1 {
		t.Errorf("expected 1 object, got %d", ps.ObjectCount())
	}
}

func TestPushSetupMarshalBinary(t *testing.T) {
	ps := &PushSetup{LogicalName: core.ObisCode{0, 0, 40, 0, 0, 255}}
	ps.AddObject(PushObject{ClassID: 8, LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}, Attribute: 2})
	data, err := ps.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary failed: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty data")
	}
}

func TestPushSetupAccess(t *testing.T) {
	ps := &PushSetup{LogicalName: core.ObisCode{0, 0, 40, 0, 0, 255}}
	if ps.Access(1).Access != core.AccessRead {
		t.Error("expected read access for attr 1")
	}
	if ps.Access(2).Access != core.AccessReadWrite {
		t.Error("expected read/write access for attr 2")
	}
}

// ============================================================================
// IC15 Association LN Tests
// ============================================================================

func TestAssociationLNClassID(t *testing.T) {
	a := &AssociationLN{LogicalName: core.ObisCode{0, 0, 15, 0, 0, 255}}
	if a.ClassID() != core.ClassIDAssociationLN {
		t.Errorf("expected ClassID %d, got %d", core.ClassIDAssociationLN, a.ClassID())
	}
}

func TestAssociationLNAddObject(t *testing.T) {
	a := &AssociationLN{LogicalName: core.ObisCode{0, 0, 15, 0, 0, 255}}
	a.AddObject(ObjectListEntry{ClassID: 1, LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}})
	if a.ObjectCount() != 1 {
		t.Errorf("expected 1 object, got %d", a.ObjectCount())
	}
}

func TestAssociationLNAuthentication(t *testing.T) {
	a := &AssociationLN{LogicalName: core.ObisCode{0, 0, 15, 0, 0, 255}}
	if a.IsAuthenticated() {
		t.Error("expected not authenticated initially")
	}
	a.SetAuthenticated(true)
	if !a.IsAuthenticated() {
		t.Error("expected authenticated after SetAuthenticated(true)")
	}
}

func TestAssociationLNMarshalBinary(t *testing.T) {
	a := &AssociationLN{LogicalName: core.ObisCode{0, 0, 15, 0, 0, 255}}
	a.AddObject(ObjectListEntry{ClassID: 1, LogicalName: core.ObisCode{0, 0, 0, 0, 0, 255}})
	data, err := a.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary failed: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty data")
	}
}

func TestAssociationLNAccess(t *testing.T) {
	a := &AssociationLN{LogicalName: core.ObisCode{0, 0, 15, 0, 0, 255}}
	if a.Access(2).Access != core.AccessRead {
		t.Error("expected read access for attr 2")
	}
}

// ============================================================================
// IC11 Special Days Table Tests
// ============================================================================

func TestSpecialDaysTableClassID(t *testing.T) {
	sdt := &SpecialDaysTable{LogicalName: core.ObisCode{0, 0, 11, 0, 0, 255}}
	if sdt.ClassID() != core.ClassIDSpecialDaysTable {
		t.Errorf("expected ClassID %d, got %d", core.ClassIDSpecialDaysTable, sdt.ClassID())
	}
}

func TestSpecialDaysTableAddEntry(t *testing.T) {
	sdt := &SpecialDaysTable{LogicalName: core.ObisCode{0, 0, 11, 0, 0, 255}}
	sdt.AddEntry(SpecialDayEntry{Index: 1, DayProfileRef: "Holiday"})
	if sdt.EntryCount() != 1 {
		t.Errorf("expected 1 entry, got %d", sdt.EntryCount())
	}
}

func TestSpecialDaysTableMarshalBinary(t *testing.T) {
	sdt := &SpecialDaysTable{LogicalName: core.ObisCode{0, 0, 11, 0, 0, 255}}
	sdt.AddEntry(SpecialDayEntry{Index: 1, DayProfileRef: "Holiday"})
	data, err := sdt.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary failed: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty data")
	}
}

func TestSpecialDaysTableAccess(t *testing.T) {
	sdt := &SpecialDaysTable{LogicalName: core.ObisCode{0, 0, 11, 0, 0, 255}}
	if sdt.Access(2).Access != core.AccessReadWrite {
		t.Error("expected read/write access for attr 2")
	}
}

// ============================================================================
// IC90 IEC Public Key Tests
// ============================================================================

func TestIECPublicKeyClassID(t *testing.T) {
	pk := &IECPublicKey{LogicalName: core.ObisCode{0, 0, 90, 0, 0, 255}}
	if pk.ClassID() != core.ClassIDIECPublicKey {
		t.Errorf("expected ClassID %d, got %d", core.ClassIDIECPublicKey, pk.ClassID())
	}
}

func TestIECPublicKeySetPublicKey(t *testing.T) {
	pk := &IECPublicKey{LogicalName: core.ObisCode{0, 0, 90, 0, 0, 255}}
	key := []byte{0x01, 0x02, 0x03, 0x04}
	pk.SetPublicKey(IECKeyTypeECDSA_P256, key)
	if pk.KeyType != IECKeyTypeECDSA_P256 {
		t.Errorf("expected key type %d, got %d", IECKeyTypeECDSA_P256, pk.KeyType)
	}
	if len(pk.PublicKeyValue) != 4 {
		t.Errorf("expected public key length 4, got %d", len(pk.PublicKeyValue))
	}
}

func TestIECPublicKeyIsValid(t *testing.T) {
	pk := &IECPublicKey{LogicalName: core.ObisCode{0, 0, 90, 0, 0, 255}}
	if pk.IsValid() {
		t.Error("expected invalid when no key is set")
	}
	pk.SetPublicKey(IECKeyTypeECDSA_P256, []byte{0x01, 0x02, 0x03, 0x04})
	if !pk.IsValid() {
		t.Error("expected valid when key is set")
	}
}

func TestIECPublicKeyMarshalBinary(t *testing.T) {
	pk := &IECPublicKey{LogicalName: core.ObisCode{0, 0, 90, 0, 0, 255}}
	pk.SetPublicKey(IECKeyTypeECDSA_P256, []byte{0x01, 0x02, 0x03, 0x04})
	data, err := pk.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary failed: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty data")
	}
}

func TestIECPublicKeyAccess(t *testing.T) {
	pk := &IECPublicKey{LogicalName: core.ObisCode{0, 0, 90, 0, 0, 255}}
	if pk.Access(2).Access != core.AccessRead {
		t.Error("expected read access for attr 2")
	}
	if pk.Access(6).Access != core.AccessReadWrite {
		t.Error("expected read/write access for attr 6")
	}
}

// ============================================================================
// IC71 Limiter Tests
// ============================================================================

func TestLimiterClassID(t *testing.T) {
	l := &Limiter{LogicalName: core.ObisCode{0, 0, 71, 0, 0, 255}}
	if l.ClassID() != core.ClassIDLimiter {
		t.Errorf("expected ClassID %d, got %d", core.ClassIDLimiter, l.ClassID())
	}
}

func TestLimiterSetThreshold(t *testing.T) {
	l := &Limiter{LogicalName: core.ObisCode{0, 0, 71, 0, 0, 255}}
	l.SetThreshold(core.DoubleLongData(100), core.DoubleLongData(50))
	if l.ThresholdActive == nil || l.ThresholdNormal == nil {
		t.Error("expected thresholds to be set")
	}
}

func TestLimiterAddAction(t *testing.T) {
	l := &Limiter{LogicalName: core.ObisCode{0, 0, 71, 0, 0, 255}}
	l.AddAction(LimiterAction{ScriptSelector: 1})
	if l.ActionCount() != 1 {
		t.Errorf("expected 1 action, got %d", l.ActionCount())
	}
}

func TestLimiterMarshalBinary(t *testing.T) {
	l := &Limiter{LogicalName: core.ObisCode{0, 0, 71, 0, 0, 255}}
	l.MonitoredValue = core.DoubleLongData(100)
	data, err := l.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary failed: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty data")
	}
}

func TestLimiterAccess(t *testing.T) {
	l := &Limiter{LogicalName: core.ObisCode{0, 0, 71, 0, 0, 255}}
	if l.Access(2).Access != core.AccessRead {
		t.Error("expected read access for attr 2")
	}
}

// ============================================================================
// Roundtrip Tests (Marshal/Unmarshal)
// ============================================================================

func TestScriptTableRoundtrip(t *testing.T) {
	st := &ScriptTable{LogicalName: core.ObisCode{0, 0, 10, 0, 0, 255}}
	st.AddScript(Script{ScriptID: 1, ScriptSelector: 0, FileID: 100})
	data, err := st.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary failed: %v", err)
	}
	st2 := &ScriptTable{}
	if err := st2.UnmarshalBinary(data); err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}
	if len(st2.Scripts) != 1 {
		t.Errorf("expected 1 script, got %d", len(st2.Scripts))
	}
	if st2.Scripts[0].ScriptID != 1 {
		t.Errorf("expected ScriptID 1, got %d", st2.Scripts[0].ScriptID)
	}
}

func TestValueTableRoundtrip(t *testing.T) {
	vt := &ValueTable{LogicalName: core.ObisCode{0, 0, 29, 0, 0, 255}}
	vt.AddValue(ValueEntry{Index: 1, Value: core.DoubleLongData(100)})
	data, err := vt.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary failed: %v", err)
	}
	vt2 := &ValueTable{}
	if err := vt2.UnmarshalBinary(data); err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}
	if vt2.ValueCount() != 1 {
		t.Errorf("expected 1 value, got %d", vt2.ValueCount())
	}
	if vt2.Values[0].Index != 1 {
		t.Errorf("expected Index 1, got %d", vt2.Values[0].Index)
	}
}

func TestDisconnectControlRoundtrip(t *testing.T) {
	dc := &DisconnectControl{LogicalName: core.ObisCode{0, 0, 70, 0, 0, 255}}
	dc.ControlState = DisconnectStateArmed
	data, err := dc.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary failed: %v", err)
	}
	dc2 := &DisconnectControl{}
	if err := dc2.UnmarshalBinary(data); err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}
	if dc2.ControlState != DisconnectStateArmed {
		t.Errorf("expected state %d, got %d", DisconnectStateArmed, dc2.ControlState)
	}
}

func TestImageTransferRoundtrip(t *testing.T) {
	it := &ImageTransfer{LogicalName: core.ObisCode{0, 0, 44, 0, 0, 255}}
	it.TransferStatus = ImageTransferStatusInitiated
	data, err := it.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary failed: %v", err)
	}
	it2 := &ImageTransfer{}
	if err := it2.UnmarshalBinary(data); err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}
	if it2.TransferStatus != ImageTransferStatusInitiated {
		t.Errorf("expected status %d, got %d", ImageTransferStatusInitiated, it2.TransferStatus)
	}
}

func TestIECPublicKeyRoundtrip(t *testing.T) {
	pk := &IECPublicKey{LogicalName: core.ObisCode{0, 0, 90, 0, 0, 255}}
	pk.SetPublicKey(IECKeyTypeECDSA_P256, []byte{0x01, 0x02, 0x03, 0x04})
	data, err := pk.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary failed: %v", err)
	}
	pk2 := &IECPublicKey{}
	if err := pk2.UnmarshalBinary(data); err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}
	if len(pk2.PublicKeyValue) != 4 {
		t.Errorf("expected public key length 4, got %d", len(pk2.PublicKeyValue))
	}
}

func TestPushSetupRoundtrip(t *testing.T) {
	ps := &PushSetup{LogicalName: core.ObisCode{0, 0, 40, 0, 0, 255}}
	ps.AddObject(PushObject{ClassID: 8, LogicalName: core.ObisCode{0, 0, 0, 1, 0, 255}, Attribute: 2})
	data, err := ps.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary failed: %v", err)
	}
	ps2 := &PushSetup{}
	if err := ps2.UnmarshalBinary(data); err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}
	if ps2.ObjectCount() != 1 {
		t.Errorf("expected 1 object, got %d", ps2.ObjectCount())
	}
}

func TestAssociationLNRoundtrip(t *testing.T) {
	a := &AssociationLN{LogicalName: core.ObisCode{0, 0, 15, 0, 0, 255}}
	a.AddObject(ObjectListEntry{ClassID: 1, LogicalName: core.ObisCode{0, 0, 0, 2, 0, 255}})
	data, err := a.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary failed: %v", err)
	}
	a2 := &AssociationLN{}
	if err := a2.UnmarshalBinary(data); err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}
	if a2.ObjectCount() != 1 {
		t.Errorf("expected 1 object, got %d", a2.ObjectCount())
	}
}

func TestSpecialDaysTableRoundtrip(t *testing.T) {
	sdt := &SpecialDaysTable{LogicalName: core.ObisCode{0, 0, 11, 0, 0, 255}}
	sdt.AddEntry(SpecialDayEntry{Index: 1, DayProfileRef: "Holiday"})
	data, err := sdt.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary failed: %v", err)
	}
	sdt2 := &SpecialDaysTable{}
	if err := sdt2.UnmarshalBinary(data); err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}
	if sdt2.EntryCount() != 1 {
		t.Errorf("expected 1 entry, got %d", sdt2.EntryCount())
	}
	if sdt2.Entries[0].DayProfileRef != "Holiday" {
		t.Errorf("expected DayProfileRef 'Holiday', got '%s'", sdt2.Entries[0].DayProfileRef)
	}
}

func TestLimiterRoundtrip(t *testing.T) {
	l := &Limiter{LogicalName: core.ObisCode{0, 0, 71, 0, 0, 255}}
	l.MonitoredValue = core.DoubleLongData(100)
	data, err := l.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary failed: %v", err)
	}
	l2 := &Limiter{}
	if err := l2.UnmarshalBinary(data); err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}
	if l2.MonitoredValue == nil {
		t.Error("expected monitored value to be set")
	}
}

// ============================================================================
// GetLogicalName and GetVersion Tests
// ============================================================================

func TestScriptTableGetLogicalName(t *testing.T) {
	ln := core.ObisCode{0, 0, 10, 0, 0, 255}
	st := &ScriptTable{LogicalName: ln}
	if st.GetLogicalName() != ln {
		t.Error("GetLogicalName mismatch")
	}
	if st.GetVersion() != 0 {
		t.Error("GetVersion mismatch")
	}
}

func TestScheduleGetLogicalName(t *testing.T) {
	ln := core.ObisCode{0, 0, 12, 0, 0, 255}
	s := &Schedule{LogicalName: ln, Version: 1}
	if s.GetLogicalName() != ln {
		t.Error("GetLogicalName mismatch")
	}
	if s.GetVersion() != 1 {
		t.Error("GetVersion mismatch")
	}
}

func TestActivityCalendarGetLogicalName(t *testing.T) {
	ln := core.ObisCode{0, 0, 20, 0, 0, 255}
	ac := &ActivityCalendar{LogicalName: ln}
	if ac.GetLogicalName() != ln {
		t.Error("GetLogicalName mismatch")
	}
}

func TestRegisterMonitorGetLogicalName(t *testing.T) {
	ln := core.ObisCode{0, 0, 21, 0, 0, 255}
	rm := &RegisterMonitor{LogicalName: ln}
	if rm.GetLogicalName() != ln {
		t.Error("GetLogicalName mismatch")
	}
}

func TestSingleActionScheduleGetLogicalName(t *testing.T) {
	ln := core.ObisCode{0, 0, 22, 0, 0, 255}
	sas := &SingleActionSchedule{LogicalName: ln}
	if sas.GetLogicalName() != ln {
		t.Error("GetLogicalName mismatch")
	}
}

func TestValueTableGetLogicalName(t *testing.T) {
	ln := core.ObisCode{0, 0, 29, 0, 0, 255}
	vt := &ValueTable{LogicalName: ln}
	if vt.GetLogicalName() != ln {
		t.Error("GetLogicalName mismatch")
	}
}

func TestDisconnectControlGetLogicalName(t *testing.T) {
	ln := core.ObisCode{0, 0, 70, 0, 0, 255}
	dc := &DisconnectControl{LogicalName: ln}
	if dc.GetLogicalName() != ln {
		t.Error("GetLogicalName mismatch")
	}
}

func TestImageTransferGetLogicalName(t *testing.T) {
	ln := core.ObisCode{0, 0, 44, 0, 0, 255}
	it := &ImageTransfer{LogicalName: ln}
	if it.GetLogicalName() != ln {
		t.Error("GetLogicalName mismatch")
	}
}

func TestPushSetupGetLogicalName(t *testing.T) {
	ln := core.ObisCode{0, 0, 40, 0, 0, 255}
	ps := &PushSetup{LogicalName: ln}
	if ps.GetLogicalName() != ln {
		t.Error("GetLogicalName mismatch")
	}
}

func TestAssociationLNGetLogicalName(t *testing.T) {
	ln := core.ObisCode{0, 0, 15, 0, 0, 255}
	a := &AssociationLN{LogicalName: ln}
	if a.GetLogicalName() != ln {
		t.Error("GetLogicalName mismatch")
	}
}

func TestSpecialDaysTableGetLogicalName(t *testing.T) {
	ln := core.ObisCode{0, 0, 11, 0, 0, 255}
	sdt := &SpecialDaysTable{LogicalName: ln}
	if sdt.GetLogicalName() != ln {
		t.Error("GetLogicalName mismatch")
	}
}

func TestIECPublicKeyGetLogicalName(t *testing.T) {
	ln := core.ObisCode{0, 0, 90, 0, 0, 255}
	pk := &IECPublicKey{LogicalName: ln}
	if pk.GetLogicalName() != ln {
		t.Error("GetLogicalName mismatch")
	}
}

func TestLimiterGetLogicalName(t *testing.T) {
	ln := core.ObisCode{0, 0, 71, 0, 0, 255}
	l := &Limiter{LogicalName: ln}
	if l.GetLogicalName() != ln {
		t.Error("GetLogicalName mismatch")
	}
}

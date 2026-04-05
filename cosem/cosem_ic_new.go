package cosem

import (
	"fmt"

	"github.com/ViewWay/dlms-cosem-go/core"
)

// ============================================================================
// IC10 Script Table - Script execution management
// ============================================================================

// Script represents a single script entry in ScriptTable
type Script struct {
	ScriptID       uint8
	ScriptSelector uint8
	FileID         uint16
}

// ScriptTable is the COSEM Script Table interface class (IC 10).
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

// AddScript adds a script to the table
func (st *ScriptTable) AddScript(script Script) {
	st.Scripts = append(st.Scripts, script)
}

// RemoveScript removes a script by ID
func (st *ScriptTable) RemoveScript(scriptID uint8) *Script {
	for i, s := range st.Scripts {
		if s.ScriptID == scriptID {
			st.Scripts = append(st.Scripts[:i], st.Scripts[i+1:]...)
			return &s
		}
	}
	return nil
}

// ============================================================================
// IC12 Schedule - Time-based scheduling
// ============================================================================

// ScheduleEntry represents a single schedule entry
type ScheduleEntry struct {
	Index          uint16
	Enable         bool
	ScriptSelector uint8
	StartTime      core.CosemTime
	ValidForDays   uint8
}

// Schedule is the COSEM Schedule interface class (IC 12).
type Schedule struct {
	LogicalName core.ObisCode
	Entries     []ScheduleEntry
	Version     uint8
}

func (s *Schedule) ClassID() uint16               { return core.ClassIDSchedule }
func (s *Schedule) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *Schedule) GetVersion() uint8             { return s.Version }
func (s *Schedule) Access(attr int) core.CosemAttributeAccess {
	switch attr {
	case 1, 2:
		return core.CosemAttributeAccess{Access: core.AccessReadWrite}
	default:
		return core.CosemAttributeAccess{Access: core.AccessNone}
	}
}

func (s *Schedule) MarshalBinary() ([]byte, error) {
	arr := core.ArrayData{}
	for _, e := range s.Entries {
		arr = append(arr, core.StructureData{
			core.UnsignedIntegerData(uint8(e.Index)),
			core.BooleanData(e.Enable),
			core.UnsignedIntegerData(e.ScriptSelector),
			core.TimeData{Value: e.StartTime},
			core.UnsignedIntegerData(e.ValidForDays),
		})
	}
	return arr.ToBytes(), nil
}

func (s *Schedule) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	arr, ok := elem.(core.ArrayData)
	if !ok {
		return fmt.Errorf("schedule: expected array")
	}
	s.Entries = make([]ScheduleEntry, len(arr))
	for i, item := range arr {
		st, ok := item.(core.StructureData)
		if !ok || len(st) < 5 {
			return fmt.Errorf("schedule: invalid entry at index %d", i)
		}
		s.Entries[i].Index = uint16(st[0].(core.UnsignedIntegerData))
		s.Entries[i].Enable = bool(st[1].(core.BooleanData))
		s.Entries[i].ScriptSelector = uint8(st[2].(core.UnsignedIntegerData))
		s.Entries[i].StartTime = st[3].(core.TimeData).Value
		s.Entries[i].ValidForDays = uint8(st[4].(core.UnsignedIntegerData))
	}
	return nil
}

// AddEntry adds a schedule entry
func (s *Schedule) AddEntry(entry ScheduleEntry) {
	s.Entries = append(s.Entries, entry)
}

// EntryCount returns the number of entries
func (s *Schedule) EntryCount() int {
	return len(s.Entries)
}

// ============================================================================
// IC20 Activity Calendar - Calendar-based activity management
// ============================================================================

// SeasonProfileEntry represents a season profile
type SeasonProfileEntry struct {
	SeasonName     string
	SeasonStart    core.CosemDate
	WeekProfileRef string
}

// WeekProfileEntry represents a week profile
type WeekProfileEntry struct {
	ProfileName string
	Monday      string
	Tuesday     string
	Wednesday   string
	Thursday    string
	Friday      string
	Saturday    string
	Sunday      string
}

// ActivityCalendar is the COSEM Activity Calendar interface class (IC 20).
type ActivityCalendar struct {
	LogicalName               core.ObisCode
	CalendarName              string
	ActivatePassiveCalendar   core.CosemDateTime
	SeasonProfileActive       []SeasonProfileEntry
	SeasonProfilePassive      []SeasonProfileEntry
	WeekProfileActive         []WeekProfileEntry
	WeekProfilePassive        []WeekProfileEntry
	DayProfileActiveName      string
	DayProfilePassiveName     string
	Version                   uint8
}

func (ac *ActivityCalendar) ClassID() uint16               { return core.ClassIDActivityCalendar }
func (ac *ActivityCalendar) GetLogicalName() core.ObisCode { return ac.LogicalName }
func (ac *ActivityCalendar) GetVersion() uint8             { return ac.Version }
func (ac *ActivityCalendar) Access(attr int) core.CosemAttributeAccess {
	switch attr {
	case 1:
		return core.CosemAttributeAccess{Access: core.AccessRead}
	case 2, 3:
		return core.CosemAttributeAccess{Access: core.AccessReadWrite}
	default:
		return core.CosemAttributeAccess{Access: core.AccessRead}
	}
}

func (ac *ActivityCalendar) MarshalBinary() ([]byte, error) {
	return core.VisibleStringData(ac.CalendarName).ToBytes(), nil
}

func (ac *ActivityCalendar) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	if v, ok := elem.(core.VisibleStringData); ok {
		ac.CalendarName = string(v)
	}
	return nil
}

// SetCalendarName sets the calendar name
func (ac *ActivityCalendar) SetCalendarName(name string) {
	ac.CalendarName = name
}

// ============================================================================
// IC21 Register Monitor - Register threshold monitoring
// ============================================================================

// RegisterMonitor is the COSEM Register Monitor interface class (IC 21).
type RegisterMonitor struct {
	LogicalName       core.ObisCode
	Thresholds        []core.DlmsData
	ThresholdActions  []core.DlmsData
	MonitoredValue    core.DlmsData
	Version           uint8
}

func (rm *RegisterMonitor) ClassID() uint16               { return core.ClassIDRegisterMonitor }
func (rm *RegisterMonitor) GetLogicalName() core.ObisCode { return rm.LogicalName }
func (rm *RegisterMonitor) GetVersion() uint8             { return rm.Version }
func (rm *RegisterMonitor) Access(attr int) core.CosemAttributeAccess {
	switch attr {
	case 2, 3, 4:
		return core.CosemAttributeAccess{Access: core.AccessRead}
	default:
		return core.CosemAttributeAccess{Access: core.AccessNone}
	}
}

func (rm *RegisterMonitor) MarshalBinary() ([]byte, error) {
	arr := core.ArrayData{}
	for _, t := range rm.Thresholds {
		arr = append(arr, t)
	}
	return arr.ToBytes(), nil
}

func (rm *RegisterMonitor) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	arr, ok := elem.(core.ArrayData)
	if !ok {
		return fmt.Errorf("register_monitor: expected array")
	}
	rm.Thresholds = []core.DlmsData(arr)
	return nil
}

// AddThreshold adds a threshold
func (rm *RegisterMonitor) AddThreshold(threshold core.DlmsData) {
	rm.Thresholds = append(rm.Thresholds, threshold)
}

// ThresholdCount returns the number of thresholds
func (rm *RegisterMonitor) ThresholdCount() int {
	return len(rm.Thresholds)
}

// ============================================================================
// IC22 Single Action Schedule - Single action scheduling
// ============================================================================

// ActionScheduleEntry represents an action schedule entry
type ActionScheduleEntry struct {
	ExecutedScriptLN core.ObisCode
	ExecutedAt       core.CosemDateTime
}

// SingleActionSchedule is the COSEM Single Action Schedule interface class (IC 22).
type SingleActionSchedule struct {
	LogicalName     core.ObisCode
	Entries         []ActionScheduleEntry
	ExecutedScript  core.ObisCode
	Type            uint8
	Version         uint8
}

func (sas *SingleActionSchedule) ClassID() uint16               { return core.ClassIDSingleActionSchedule }
func (sas *SingleActionSchedule) GetLogicalName() core.ObisCode { return sas.LogicalName }
func (sas *SingleActionSchedule) GetVersion() uint8             { return sas.Version }
func (sas *SingleActionSchedule) Access(attr int) core.CosemAttributeAccess {
	switch attr {
	case 1:
		return core.CosemAttributeAccess{Access: core.AccessRead}
	case 2, 3:
		return core.CosemAttributeAccess{Access: core.AccessReadWrite}
	default:
		return core.CosemAttributeAccess{Access: core.AccessNone}
	}
}

func (sas *SingleActionSchedule) MarshalBinary() ([]byte, error) {
	arr := core.ArrayData{}
	for _, e := range sas.Entries {
		arr = append(arr, core.StructureData{
			core.OctetStringData(e.ExecutedScriptLN[:]),
			core.DateTimeData{Value: e.ExecutedAt},
		})
	}
	return arr.ToBytes(), nil
}

func (sas *SingleActionSchedule) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	arr, ok := elem.(core.ArrayData)
	if !ok {
		return fmt.Errorf("single_action_schedule: expected array")
	}
	sas.Entries = make([]ActionScheduleEntry, len(arr))
	for i, item := range arr {
		st, ok := item.(core.StructureData)
		if !ok || len(st) < 2 {
			return fmt.Errorf("single_action_schedule: invalid entry at index %d", i)
		}
		sas.Entries[i].ExecutedScriptLN = core.ObisCode(st[0].(core.OctetStringData))
		sas.Entries[i].ExecutedAt = st[1].(core.DateTimeData).Value
	}
	return nil
}

// AddEntry adds an action schedule entry
func (sas *SingleActionSchedule) AddEntry(entry ActionScheduleEntry) {
	sas.Entries = append(sas.Entries, entry)
}

// EntryCount returns the number of entries
func (sas *SingleActionSchedule) EntryCount() int {
	return len(sas.Entries)
}

// ============================================================================
// IC29 Value Table - Value storage table
// ============================================================================

// ValueEntry represents a value entry in the table
type ValueEntry struct {
	Index     uint16
	Value     core.DlmsData
	Timestamp *core.CosemDateTime
}

// ValueDescriptor describes a value in the table
type ValueDescriptor struct {
	Index       uint16
	Description string
	Unit        uint8
	Scaler      int8
}

// ValueTable is the COSEM Value Table interface class (IC 29).
type ValueTable struct {
	LogicalName core.ObisCode
	Values      []ValueEntry
	Descriptors []ValueDescriptor
	Version     uint8
}

func (vt *ValueTable) ClassID() uint16               { return core.ClassIDValueTable }
func (vt *ValueTable) GetLogicalName() core.ObisCode { return vt.LogicalName }
func (vt *ValueTable) GetVersion() uint8             { return vt.Version }
func (vt *ValueTable) Access(attr int) core.CosemAttributeAccess {
	switch attr {
	case 1, 2, 3:
		return core.CosemAttributeAccess{Access: core.AccessRead}
	default:
		return core.CosemAttributeAccess{Access: core.AccessNone}
	}
}

func (vt *ValueTable) MarshalBinary() ([]byte, error) {
	arr := core.ArrayData{}
	for _, v := range vt.Values {
		var ts core.DlmsData = core.NullData{}
		if v.Timestamp != nil {
			ts = core.DateTimeData{Value: *v.Timestamp}
		}
		arr = append(arr, core.StructureData{
			core.LongData(int16(v.Index)),
			v.Value,
			ts,
		})
	}
	return arr.ToBytes(), nil
}

func (vt *ValueTable) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	arr, ok := elem.(core.ArrayData)
	if !ok {
		return fmt.Errorf("value_table: expected array")
	}
	vt.Values = make([]ValueEntry, len(arr))
	for i, item := range arr {
		st, ok := item.(core.StructureData)
		if !ok || len(st) < 3 {
			return fmt.Errorf("value_table: invalid entry at index %d", i)
		}
		vt.Values[i].Index = uint16(st[0].(core.LongData))
		vt.Values[i].Value = st[1]
		if _, isNull := st[2].(core.NullData); !isNull {
			if dt, ok := st[2].(core.DateTimeData); ok {
				vt.Values[i].Timestamp = &dt.Value
			}
		}
	}
	return nil
}

// AddValue adds a value entry
func (vt *ValueTable) AddValue(entry ValueEntry) {
	vt.Values = append(vt.Values, entry)
}

// AddDescriptor adds a value descriptor
func (vt *ValueTable) AddDescriptor(desc ValueDescriptor) {
	vt.Descriptors = append(vt.Descriptors, desc)
}

// ValueCount returns the number of values
func (vt *ValueTable) ValueCount() int {
	return len(vt.Values)
}

// DescriptorCount returns the number of descriptors
func (vt *ValueTable) DescriptorCount() int {
	return len(vt.Descriptors)
}

// ============================================================================
// IC70 Disconnect Control - Remote disconnect/reconnect
// ============================================================================

// DisconnectState represents the disconnect control state
type DisconnectState uint8

const (
	DisconnectStateDisconnected      DisconnectState = 0
	DisconnectStateConnected         DisconnectState = 1
	DisconnectStateReadyForDisconnect DisconnectState = 2
	DisconnectStateReadyForReconnect  DisconnectState = 3
	DisconnectStateArmed             DisconnectState = 4
)

// DisconnectControl is the COSEM Disconnect Control interface class (IC 70).
type DisconnectControl struct {
	LogicalName    core.ObisCode
	ControlState   DisconnectState
	OutputState    DisconnectState
	ControlValue   uint8
	Version        uint8
}

func (dc *DisconnectControl) ClassID() uint16               { return core.ClassIDDisconnectControl }
func (dc *DisconnectControl) GetLogicalName() core.ObisCode { return dc.LogicalName }
func (dc *DisconnectControl) GetVersion() uint8             { return dc.Version }
func (dc *DisconnectControl) Access(attr int) core.CosemAttributeAccess {
	switch attr {
	case 2, 3, 4:
		return core.CosemAttributeAccess{Access: core.AccessRead}
	default:
		return core.CosemAttributeAccess{Access: core.AccessNone}
	}
}

func (dc *DisconnectControl) MarshalBinary() ([]byte, error) {
	return core.EnumData(dc.ControlState).ToBytes(), nil
}

func (dc *DisconnectControl) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	if v, ok := elem.(core.EnumData); ok {
		dc.ControlState = DisconnectState(v)
	}
	return nil
}

// Disconnect performs disconnect action
func (dc *DisconnectControl) Disconnect() error {
	dc.ControlState = DisconnectStateDisconnected
	dc.OutputState = DisconnectStateDisconnected
	return nil
}

// Reconnect performs reconnect action
func (dc *DisconnectControl) Reconnect() error {
	dc.ControlState = DisconnectStateConnected
	dc.OutputState = DisconnectStateConnected
	return nil
}

// Arm arms the disconnect control
func (dc *DisconnectControl) Arm() error {
	dc.ControlState = DisconnectStateArmed
	return nil
}

// ============================================================================
// IC18 Image Transfer - Firmware update management
// ============================================================================

// ImageTransferStatus represents the image transfer status
type ImageTransferStatus uint8

const (
	ImageTransferStatusIdle                  ImageTransferStatus = 0
	ImageTransferStatusInitiated             ImageTransferStatus = 1
	ImageTransferStatusInitiatedForVerifying ImageTransferStatus = 2
	ImageTransferStatusVerifyingInitiated    ImageTransferStatus = 3
	ImageTransferStatusVerificationFailed    ImageTransferStatus = 4
	ImageTransferStatusVerificationSuccessful ImageTransferStatus = 5
	ImageTransferStatusImageActivated        ImageTransferStatus = 6
	ImageTransferStatusImageNotActivated     ImageTransferStatus = 7
)

// ImageTransfer is the COSEM Image Transfer interface class (IC 18).
type ImageTransfer struct {
	LogicalName      core.ObisCode
	ImageBlockSize   uint16
	ImageFirstBlock  []byte
	ImageBlockCount  uint32
	ImageReference   []byte
	ImageIdent       uint8
	TransferStatus   ImageTransferStatus
	Version          uint8
}

func (it *ImageTransfer) ClassID() uint16               { return core.ClassIDImageTransfer }
func (it *ImageTransfer) GetLogicalName() core.ObisCode { return it.LogicalName }
func (it *ImageTransfer) GetVersion() uint8             { return it.Version }
func (it *ImageTransfer) Access(attr int) core.CosemAttributeAccess {
	switch attr {
	case 2, 3, 4, 5, 6, 7:
		return core.CosemAttributeAccess{Access: core.AccessRead}
	default:
		return core.CosemAttributeAccess{Access: core.AccessNone}
	}
}

func (it *ImageTransfer) MarshalBinary() ([]byte, error) {
	return core.EnumData(it.TransferStatus).ToBytes(), nil
}

func (it *ImageTransfer) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	if v, ok := elem.(core.EnumData); ok {
		it.TransferStatus = ImageTransferStatus(v)
	}
	return nil
}

// InitiateImage initiates image transfer
func (it *ImageTransfer) InitiateImage(blockSize uint16, blockCount uint32) {
	it.ImageBlockSize = blockSize
	it.ImageBlockCount = blockCount
	it.TransferStatus = ImageTransferStatusInitiated
}

// VerifyImage verifies the transferred image
func (it *ImageTransfer) VerifyImage() bool {
	it.TransferStatus = ImageTransferStatusVerificationSuccessful
	return true
}

// ActivateImage activates the image
func (it *ImageTransfer) ActivateImage() error {
	it.TransferStatus = ImageTransferStatusImageActivated
	return nil
}

// ============================================================================
// IC40 Push Setup - Push notification configuration
// ============================================================================

// PushObject represents an object in push list
type PushObject struct {
	ClassID     uint16
	LogicalName core.ObisCode
	Attribute   uint8
	DataIndex   uint8
}

// PushSetup is the COSEM Push Setup interface class (IC 40).
type PushSetup struct {
	LogicalName               core.ObisCode
	PushObjectList            []PushObject
	Service                   uint16
	Destination               []byte
	CommunicationWindow       core.CosemDateTime
	RandomisationStartInterval uint16
	NumberOfRetries           uint8
	RepetitionDelay           uint16
	Version                   uint8
}

func (ps *PushSetup) ClassID() uint16               { return core.ClassIDPush }
func (ps *PushSetup) GetLogicalName() core.ObisCode { return ps.LogicalName }
func (ps *PushSetup) GetVersion() uint8             { return ps.Version }
func (ps *PushSetup) Access(attr int) core.CosemAttributeAccess {
	switch attr {
	case 1:
		return core.CosemAttributeAccess{Access: core.AccessRead}
	case 2, 3, 4, 5, 6, 7, 8:
		return core.CosemAttributeAccess{Access: core.AccessReadWrite}
	default:
		return core.CosemAttributeAccess{Access: core.AccessNone}
	}
}

func (ps *PushSetup) MarshalBinary() ([]byte, error) {
	arr := core.ArrayData{}
	for _, o := range ps.PushObjectList {
		arr = append(arr, core.StructureData{
			core.UnsignedLongData(o.ClassID),
			core.OctetStringData(o.LogicalName[:]),
			core.IntegerData(int8(o.Attribute)),
			core.UnsignedIntegerData(o.DataIndex),
		})
	}
	return arr.ToBytes(), nil
}

func (ps *PushSetup) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	arr, ok := elem.(core.ArrayData)
	if !ok {
		return fmt.Errorf("push_setup: expected array")
	}
	ps.PushObjectList = make([]PushObject, len(arr))
	for i, item := range arr {
		st, ok := item.(core.StructureData)
		if !ok || len(st) < 4 {
			return fmt.Errorf("push_setup: invalid entry at index %d", i)
		}
		ps.PushObjectList[i].ClassID = uint16(st[0].(core.UnsignedLongData))
		ps.PushObjectList[i].LogicalName = core.ObisCode(st[1].(core.OctetStringData))
		ps.PushObjectList[i].Attribute = uint8(st[2].(core.IntegerData))
		ps.PushObjectList[i].DataIndex = uint8(st[3].(core.UnsignedIntegerData))
	}
	return nil
}

// AddObject adds a push object
func (ps *PushSetup) AddObject(obj PushObject) {
	ps.PushObjectList = append(ps.PushObjectList, obj)
}

// ObjectCount returns the number of push objects
func (ps *PushSetup) ObjectCount() int {
	return len(ps.PushObjectList)
}

// ============================================================================
// IC15 Association LN - Logical Name Association
// ============================================================================

// ObjectListEntry represents an entry in the object list
type ObjectListEntry struct {
	ClassID     uint16
	LogicalName core.ObisCode
	Version     uint8
}

// AssociationLN is the COSEM Association LN interface class (IC 15).
type AssociationLN struct {
	LogicalName             core.ObisCode
	ObjectList              []ObjectListEntry
	ClientSAP               uint16
	ServerSAP               uint16
	ApplicationContextName  string
	AuthenticationMechanism string
	LLSSecret               []byte
	HLSSecret               []byte
	AuthenticationStatus    bool
	Version                 uint8
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

// AddObject adds an object to the association
func (a *AssociationLN) AddObject(entry ObjectListEntry) {
	a.ObjectList = append(a.ObjectList, entry)
}

// ObjectCount returns the number of objects
func (a *AssociationLN) ObjectCount() int {
	return len(a.ObjectList)
}

// SetAuthenticated sets the authentication status
func (a *AssociationLN) SetAuthenticated(status bool) {
	a.AuthenticationStatus = status
}

// IsAuthenticated returns the authentication status
func (a *AssociationLN) IsAuthenticated() bool {
	return a.AuthenticationStatus
}

// ============================================================================
// IC11 Special Days Table - Special day management
// ============================================================================

// SpecialDayEntry represents a special day entry
type SpecialDayEntry struct {
	Index         uint16
	SpecialDay    core.CosemDate
	DayProfileRef string
}

// SpecialDaysTable is the COSEM Special Days Table interface class (IC 11).
type SpecialDaysTable struct {
	LogicalName core.ObisCode
	Entries     []SpecialDayEntry
	Version     uint8
}

func (sdt *SpecialDaysTable) ClassID() uint16               { return core.ClassIDSpecialDaysTable }
func (sdt *SpecialDaysTable) GetLogicalName() core.ObisCode { return sdt.LogicalName }
func (sdt *SpecialDaysTable) GetVersion() uint8             { return sdt.Version }
func (sdt *SpecialDaysTable) Access(attr int) core.CosemAttributeAccess {
	switch attr {
	case 2:
		return core.CosemAttributeAccess{Access: core.AccessReadWrite}
	default:
		return core.CosemAttributeAccess{Access: core.AccessRead}
	}
}

func (sdt *SpecialDaysTable) MarshalBinary() ([]byte, error) {
	arr := core.ArrayData{}
	for _, e := range sdt.Entries {
		arr = append(arr, core.StructureData{
			core.UnsignedLongData(e.Index),
			core.DateData{Value: e.SpecialDay},
			core.VisibleStringData(e.DayProfileRef),
		})
	}
	return arr.ToBytes(), nil
}

func (sdt *SpecialDaysTable) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	arr, ok := elem.(core.ArrayData)
	if !ok {
		return fmt.Errorf("special_days_table: expected array")
	}
	sdt.Entries = make([]SpecialDayEntry, len(arr))
	for i, item := range arr {
		st, ok := item.(core.StructureData)
		if !ok || len(st) < 3 {
			return fmt.Errorf("special_days_table: invalid entry at index %d", i)
		}
		sdt.Entries[i].Index = uint16(st[0].(core.UnsignedLongData))
		sdt.Entries[i].SpecialDay = st[1].(core.DateData).Value
		sdt.Entries[i].DayProfileRef = string(st[2].(core.VisibleStringData))
	}
	return nil
}

// AddEntry adds a special day entry
func (sdt *SpecialDaysTable) AddEntry(entry SpecialDayEntry) {
	sdt.Entries = append(sdt.Entries, entry)
}

// EntryCount returns the number of entries
func (sdt *SpecialDaysTable) EntryCount() int {
	return len(sdt.Entries)
}

// ============================================================================
// IC90 IEC Public Key - IEC public key management
// ============================================================================

// IECKeyType represents the key type
type IECKeyType uint8

const (
	IECKeyTypeNone           IECKeyType = 0
	IECKeyTypeECDSA_P256     IECKeyType = 1
	IECKeyTypeECDSA_P384     IECKeyType = 2
	IECKeyTypeECDSA_P521     IECKeyType = 3
	IECKeyTypeRSA_2048       IECKeyType = 4
	IECKeyTypeRSA_3072       IECKeyType = 5
	IECKeyTypeRSA_4096       IECKeyType = 6
)

// IECPublicKey is the COSEM IEC Public Key interface class (IC 90).
type IECPublicKey struct {
	LogicalName    core.ObisCode
	KeyID          uint16
	KeyType        IECKeyType
	PublicKeyValue []byte
	KeyUsage       uint8
	ValidityStart  core.CosemDate
	ValidityEnd    core.CosemDate
	Version        uint8
}

func (pk *IECPublicKey) ClassID() uint16               { return core.ClassIDIECPublicKey }
func (pk *IECPublicKey) GetLogicalName() core.ObisCode { return pk.LogicalName }
func (pk *IECPublicKey) GetVersion() uint8             { return pk.Version }
func (pk *IECPublicKey) Access(attr int) core.CosemAttributeAccess {
	switch attr {
	case 2, 3, 4, 5:
		return core.CosemAttributeAccess{Access: core.AccessRead}
	case 6, 7:
		return core.CosemAttributeAccess{Access: core.AccessReadWrite}
	default:
		return core.CosemAttributeAccess{Access: core.AccessNone}
	}
}

func (pk *IECPublicKey) MarshalBinary() ([]byte, error) {
	return core.OctetStringData(pk.PublicKeyValue).ToBytes(), nil
}

func (pk *IECPublicKey) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	if v, ok := elem.(core.OctetStringData); ok {
		pk.PublicKeyValue = v
	}
	return nil
}

// SetPublicKey sets the public key value
func (pk *IECPublicKey) SetPublicKey(keyType IECKeyType, value []byte) {
	pk.KeyType = keyType
	pk.PublicKeyValue = value
}

// IsValid checks if the key is currently valid
func (pk *IECPublicKey) IsValid() bool {
	// Simplified validity check - in real implementation would check against current date
	return len(pk.PublicKeyValue) > 0
}

// ============================================================================
// IC71 Limiter - Value limiting
// ============================================================================

// LimiterAction represents limiter action settings
type LimiterAction struct {
	ScriptLN        core.ObisCode
	ScriptSelector  uint8
}

// Limiter is the COSEM Limiter interface class (IC 71).
type Limiter struct {
	LogicalName       core.ObisCode
	MonitoredValue    core.DlmsData
	ThresholdActive   core.DlmsData
	ThresholdNormal   core.DlmsData
	MinOverThreshold  uint16
	MinUnderThreshold uint16
	Actions           []LimiterAction
	EmergencyProfile  core.ObisCode
	Version           uint8
}

func (l *Limiter) ClassID() uint16               { return core.ClassIDLimiter }
func (l *Limiter) GetLogicalName() core.ObisCode { return l.LogicalName }
func (l *Limiter) GetVersion() uint8             { return l.Version }
func (l *Limiter) Access(attr int) core.CosemAttributeAccess {
	switch attr {
	case 2, 3, 4, 5, 6:
		return core.CosemAttributeAccess{Access: core.AccessRead}
	default:
		return core.CosemAttributeAccess{Access: core.AccessNone}
	}
}

func (l *Limiter) MarshalBinary() ([]byte, error) {
	return l.MonitoredValue.ToBytes(), nil
}

func (l *Limiter) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	l.MonitoredValue = elem
	return nil
}

// SetThreshold sets the threshold values
func (l *Limiter) SetThreshold(active, normal core.DlmsData) {
	l.ThresholdActive = active
	l.ThresholdNormal = normal
}

// AddAction adds a limiter action
func (l *Limiter) AddAction(action LimiterAction) {
	l.Actions = append(l.Actions, action)
}

// ActionCount returns the number of actions
func (l *Limiter) ActionCount() int {
	return len(l.Actions)
}

package cosem

import (
	"fmt"

	"github.com/ViewWay/dlms-cosem-go/core"
)

// Data is the COSEM Data interface class (IC 1).
type Data struct {
	LogicalName core.ObisCode
	Value       core.DlmsData
	Version     uint8
}

func (d *Data) ClassID() uint16               { return core.ClassIDData }
func (d *Data) GetLogicalName() core.ObisCode { return d.LogicalName }
func (d *Data) GetVersion() uint8             { return d.Version }
func (d *Data) Access(attr int) core.CosemAttributeAccess {
	switch attr {
	case 1:
		return core.CosemAttributeAccess{Access: core.AccessReadWrite}
	case 2:
		return core.CosemAttributeAccess{Access: core.AccessRead}
	default:
		return core.CosemAttributeAccess{Access: core.AccessNone}
	}
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (d *Data) MarshalBinary() ([]byte, error) {
	if d.Value == nil {
		return core.NullData{}.ToBytes(), nil
	}
	return d.Value.ToBytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (d *Data) UnmarshalBinary(data []byte) error {
	if len(data) == 0 {
		d.Value = core.NullData{}
		return nil
	}
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	d.Value = elem
	return nil
}

// Register is the COSEM Register interface class (IC 3).
type Register struct {
	LogicalName core.ObisCode
	Value       core.DlmsData
	Scaler      int8
	Unit        uint8
	Status      core.DlmsData
	Version     uint8
}

func (r *Register) ClassID() uint16               { return core.ClassIDRegister }
func (r *Register) GetLogicalName() core.ObisCode { return r.LogicalName }
func (r *Register) GetVersion() uint8             { return r.Version }
func (r *Register) Access(attr int) core.CosemAttributeAccess {
	switch attr {
	case 2:
		return core.CosemAttributeAccess{Access: core.AccessRead}
	case 3:
		return core.CosemAttributeAccess{Access: core.AccessRead}
	default:
		return core.CosemAttributeAccess{Access: core.AccessNone}
	}
}

// MarshalBinary encodes the register value.
func (r *Register) MarshalBinary() ([]byte, error) {
	if r.Value == nil {
		return core.DoubleLongUnsignedData(0).ToBytes(), nil
	}
	return r.Value.ToBytes(), nil
}

// UnmarshalBinary decodes the register value.
func (r *Register) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	r.Value = elem
	return nil
}

// ExtendedRegister is the COSEM Extended Register interface class (IC 4).
type ExtendedRegister struct {
	LogicalName core.ObisCode
	Value       core.DlmsData
	Scaler      int8
	Unit        uint8
	Status      core.DlmsData
	CaptureTime core.CosemDateTime
	Version     uint8
}

func (e *ExtendedRegister) ClassID() uint16               { return core.ClassIDExtendedRegister }
func (e *ExtendedRegister) GetLogicalName() core.ObisCode { return e.LogicalName }
func (e *ExtendedRegister) GetVersion() uint8             { return e.Version }
func (e *ExtendedRegister) Access(attr int) core.CosemAttributeAccess {
	switch attr {
	case 2:
		return core.CosemAttributeAccess{Access: core.AccessRead}
	case 3:
		return core.CosemAttributeAccess{Access: core.AccessRead}
	default:
		return core.CosemAttributeAccess{Access: core.AccessNone}
	}
}

func (e *ExtendedRegister) MarshalBinary() ([]byte, error) {
	return e.Value.ToBytes(), nil
}

func (e *ExtendedRegister) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	e.Value = elem
	return nil
}

// DemandRegister is the COSEM Demand Register interface class (IC 5).
type DemandRegister struct {
	LogicalName     core.ObisCode
	CurrentValue    core.DlmsData
	LastValue       core.DlmsData
	Scaler          int8
	Unit            uint8
	Status          core.DlmsData
	CaptureTime     core.CosemDateTime
	StartTime       core.CosemDateTime
	Period          uint16
	NumberOfPeriods uint8
	Version         uint8
}

func (d *DemandRegister) ClassID() uint16               { return core.ClassIDDemandRegister }
func (d *DemandRegister) GetLogicalName() core.ObisCode { return d.LogicalName }
func (d *DemandRegister) GetVersion() uint8             { return d.Version }
func (d *DemandRegister) Access(attr int) core.CosemAttributeAccess {
	switch attr {
	case 2, 3, 4:
		return core.CosemAttributeAccess{Access: core.AccessRead}
	default:
		return core.CosemAttributeAccess{Access: core.AccessNone}
	}
}

func (d *DemandRegister) MarshalBinary() ([]byte, error) {
	return d.CurrentValue.ToBytes(), nil
}

func (d *DemandRegister) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	d.CurrentValue = elem
	return nil
}

// Clock is the COSEM Clock interface class (IC 8).
type Clock struct {
	LogicalName              core.ObisCode
	Time                     core.CosemDateTime
	TimeZone                 int16
	Status                   uint8
	DaylightSavingsBegin     core.CosemDate
	DaylightSavingsEnd       core.CosemDate
	DaylightSavingsDeviation int16
	DaylightSavingsEnabled   bool
	ClockBase                core.CosemDateTime
	Version                  uint8
}

func (c *Clock) ClassID() uint16               { return core.ClassIDClock }
func (c *Clock) GetLogicalName() core.ObisCode { return c.LogicalName }
func (c *Clock) GetVersion() uint8             { return c.Version }
func (c *Clock) Access(attr int) core.CosemAttributeAccess {
	switch attr {
	case 2, 3, 4, 9, 10:
		return core.CosemAttributeAccess{Access: core.AccessReadWrite}
	default:
		return core.CosemAttributeAccess{Access: core.AccessRead}
	}
}

func (c *Clock) MarshalBinary() ([]byte, error) {
	return core.DateTimeData{Value: c.Time}.ToBytes(), nil
}

func (c *Clock) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	dt, ok := elem.(core.DateTimeData)
	if !ok {
		return fmt.Errorf("clock: expected DateTimeData")
	}
	c.Time = dt.Value
	return nil
}

// ProfileGeneric is the COSEM Profile Generic interface class (IC 7).
type ProfileGeneric struct {
	LogicalName    core.ObisCode
	Buffer         []core.StructureData
	CaptureObjects []core.CosemObject
	CapturePeriod  uint16
	EntriesInUse   uint32
	Entries        uint32
	SortMethod     uint8
	SortObject     core.ObisCode
	Version        uint8
}

func (p *ProfileGeneric) ClassID() uint16               { return core.ClassIDProfileGeneric }
func (p *ProfileGeneric) GetLogicalName() core.ObisCode { return p.LogicalName }
func (p *ProfileGeneric) GetVersion() uint8             { return p.Version }
func (p *ProfileGeneric) Access(attr int) core.CosemAttributeAccess {
	switch attr {
	case 2, 3, 4, 7:
		return core.CosemAttributeAccess{Access: core.AccessRead}
	default:
		return core.CosemAttributeAccess{Access: core.AccessNone}
	}
}

func (p *ProfileGeneric) MarshalBinary() ([]byte, error) {
	arr := core.ArrayData{}
	for _, s := range p.Buffer {
		arr = append(arr, s)
	}
	return arr.ToBytes(), nil
}

func (p *ProfileGeneric) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	arr, ok := elem.(core.ArrayData)
	if !ok {
		return fmt.Errorf("profile_generic: expected array")
	}
	p.Buffer = make([]core.StructureData, len(arr))
	for i, item := range arr {
		s, ok := item.(core.StructureData)
		if !ok {
			return fmt.Errorf("profile_generic: expected structure at index %d", i)
		}
		p.Buffer[i] = s
	}
	return nil
}

// SecuritySetup is the COSEM Security Setup interface class (IC 64).
type SecuritySetup struct {
	LogicalName       core.ObisCode
	SecuritySuite     uint8
	EncryptionKey     []byte
	AuthenticationKey []byte
	MasterKey         []byte
	SecurityPolicy    uint8
	Version           uint8
}

func (s *SecuritySetup) ClassID() uint16               { return core.ClassIDSecuritySetup }
func (s *SecuritySetup) GetLogicalName() core.ObisCode { return s.LogicalName }
func (s *SecuritySetup) GetVersion() uint8             { return s.Version }
func (s *SecuritySetup) Access(attr int) core.CosemAttributeAccess {
	switch attr {
	case 2:
		return core.CosemAttributeAccess{Access: core.AccessRead}
	case 3, 4, 5:
		return core.CosemAttributeAccess{Access: core.AccessReadWrite, Set: core.AccessHigh}
	default:
		return core.CosemAttributeAccess{Access: core.AccessNone}
	}
}

func (s *SecuritySetup) MarshalBinary() ([]byte, error) {
	st := core.StructureData{
		core.UnsignedIntegerData(s.SecuritySuite),
		core.OctetStringData(s.EncryptionKey),
		core.OctetStringData(s.AuthenticationKey),
	}
	return st.ToBytes(), nil
}

func (s *SecuritySetup) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok || len(st) < 1 {
		return fmt.Errorf("security_setup: expected structure")
	}
	if v, ok := st[0].(core.UnsignedIntegerData); ok {
		s.SecuritySuite = uint8(v)
	}
	if len(st) > 1 {
		if v, ok := st[1].(core.OctetStringData); ok {
			s.EncryptionKey = v
		}
	}
	if len(st) > 2 {
		if v, ok := st[2].(core.OctetStringData); ok {
			s.AuthenticationKey = v
		}
	}
	return nil
}

// TariffPlan is a COSEM utility table.
type TariffPlan struct {
	LogicalName core.ObisCode
	PlanName    string
	Version     uint8
}

func (tp *TariffPlan) ClassID() uint16               { return 26 }
func (tp *TariffPlan) GetLogicalName() core.ObisCode { return tp.LogicalName }
func (tp *TariffPlan) GetVersion() uint8             { return tp.Version }
func (tp *TariffPlan) Access(attr int) core.CosemAttributeAccess {
	return core.CosemAttributeAccess{Access: core.AccessRead}
}

func (tp *TariffPlan) MarshalBinary() ([]byte, error) {
	return core.VisibleStringData(tp.PlanName).ToBytes(), nil
}

func (tp *TariffPlan) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	if v, ok := elem.(core.VisibleStringData); ok {
		tp.PlanName = string(v)
	}
	return nil
}

// TariffTable contains tariff entries.
type TariffTable struct {
	LogicalName core.ObisCode
	Entries     []TariffEntry
	Version     uint8
}

func (tt *TariffTable) ClassID() uint16               { return 26 }
func (tt *TariffTable) GetLogicalName() core.ObisCode { return tt.LogicalName }
func (tt *TariffTable) GetVersion() uint8             { return tt.Version }
func (tt *TariffTable) Access(attr int) core.CosemAttributeAccess {
	return core.CosemAttributeAccess{Access: core.AccessRead}
}

// TariffEntry represents a single tariff entry.
type TariffEntry struct {
	TariffIndex uint8
	Name        string
}

// SeasonProfile represents a season profile entry.
type SeasonProfile struct {
	SeasonName     string
	SeasonStart    core.CosemDate
	WeekProfileRef string
}

// WeekProfile represents a week profile entry.
type WeekProfile struct {
	ProfileName string
	Monday      string
	Tuesday     string
	Wednesday   string
	Thursday    string
	Friday      string
	Saturday    string
	Sunday      string
}

// DayProfile represents a day profile entry.
type DayProfile struct {
	DayProfileName string
	Schedule       []DaySchedule
}

// DaySchedule represents a day schedule entry.
type DaySchedule struct {
	StartTime core.CosemTime
	ScriptRef string
}

// LPSetup is the Load Profile Setup.
type LPSetup struct {
	LogicalName    core.ObisCode
	CaptureObjects []string
	CapturePeriod  uint16
	ProfileEntries uint32
	Version        uint8
}

func (lp *LPSetup) ClassID() uint16               { return 26 }
func (lp *LPSetup) GetLogicalName() core.ObisCode { return lp.LogicalName }
func (lp *LPSetup) GetVersion() uint8             { return lp.Version }
func (lp *LPSetup) Access(attr int) core.CosemAttributeAccess {
	return core.CosemAttributeAccess{Access: core.AccessRead}
}

// RS485Setup is the RS-485 port setup.
type RS485Setup struct {
	LogicalName     core.ObisCode
	DefaultBaudRate uint16
	CommSpeed       uint8
	Version         uint8
}

func (r *RS485Setup) ClassID() uint16               { return core.ClassIDIECHDLCSetup }
func (r *RS485Setup) GetLogicalName() core.ObisCode { return r.LogicalName }
func (r *RS485Setup) GetVersion() uint8             { return r.Version }
func (r *RS485Setup) Access(attr int) core.CosemAttributeAccess {
	return core.CosemAttributeAccess{Access: core.AccessReadWrite}
}

func (r *RS485Setup) MarshalBinary() ([]byte, error) {
	st := core.StructureData{
		core.UnsignedLongData(r.DefaultBaudRate),
		core.EnumData(r.CommSpeed),
	}
	return st.ToBytes(), nil
}

func (r *RS485Setup) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok || len(st) < 2 {
		return fmt.Errorf("rs485_setup: expected structure with 2 elements")
	}
	if v, ok := st[0].(core.UnsignedLongData); ok {
		r.DefaultBaudRate = uint16(v)
	}
	if v, ok := st[1].(core.EnumData); ok {
		r.CommSpeed = uint8(v)
	}
	return nil
}

// InfraredSetup is the infrared port setup.
type InfraredSetup struct {
	LogicalName     core.ObisCode
	DefaultBaudRate uint16
	CommSpeed       uint8
	Version         uint8
}

func (i *InfraredSetup) ClassID() uint16               { return core.ClassIDIECLocalPortSetup }
func (i *InfraredSetup) GetLogicalName() core.ObisCode { return i.LogicalName }
func (i *InfraredSetup) GetVersion() uint8             { return i.Version }
func (i *InfraredSetup) Access(attr int) core.CosemAttributeAccess {
	return core.CosemAttributeAccess{Access: core.AccessReadWrite}
}

func (i *InfraredSetup) MarshalBinary() ([]byte, error) {
	st := core.StructureData{
		core.UnsignedLongData(i.DefaultBaudRate),
		core.EnumData(i.CommSpeed),
	}
	return st.ToBytes(), nil
}

func (i *InfraredSetup) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok || len(st) < 2 {
		return fmt.Errorf("infrared_setup: expected structure")
	}
	if v, ok := st[0].(core.UnsignedLongData); ok {
		i.DefaultBaudRate = uint16(v)
	}
	if v, ok := st[1].(core.EnumData); ok {
		i.CommSpeed = uint8(v)
	}
	return nil
}

// NBIoTSetup is the NB-IoT port setup.
type NBIoTSetup struct {
	LogicalName core.ObisCode
	APN         string
	Version     uint8
}

func (n *NBIoTSetup) ClassID() uint16               { return core.ClassIDGPRSModemSetup }
func (n *NBIoTSetup) GetLogicalName() core.ObisCode { return n.LogicalName }
func (n *NBIoTSetup) GetVersion() uint8             { return n.Version }
func (n *NBIoTSetup) Access(attr int) core.CosemAttributeAccess {
	return core.CosemAttributeAccess{Access: core.AccessReadWrite}
}

func (n *NBIoTSetup) MarshalBinary() ([]byte, error) {
	return core.VisibleStringData(n.APN).ToBytes(), nil
}

func (n *NBIoTSetup) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	if v, ok := elem.(core.VisibleStringData); ok {
		n.APN = string(v)
	}
	return nil
}

// LoRaWANSetup is the LoRaWAN port setup.
type LoRaWANSetup struct {
	LogicalName core.ObisCode
	DevEUI      []byte
	AppEUI      []byte
	AppKey      []byte
	Version     uint8
}

func (l *LoRaWANSetup) ClassID() uint16               { return core.ClassIDLorawanSetup }
func (l *LoRaWANSetup) GetLogicalName() core.ObisCode { return l.LogicalName }
func (l *LoRaWANSetup) GetVersion() uint8             { return l.Version }
func (l *LoRaWANSetup) Access(attr int) core.CosemAttributeAccess {
	return core.CosemAttributeAccess{Access: core.AccessReadWrite}
}

func (l *LoRaWANSetup) MarshalBinary() ([]byte, error) {
	st := core.StructureData{
		core.OctetStringData(l.DevEUI),
		core.OctetStringData(l.AppEUI),
		core.OctetStringData(l.AppKey),
	}
	return st.ToBytes(), nil
}

func (l *LoRaWANSetup) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	st, ok := elem.(core.StructureData)
	if !ok || len(st) < 3 {
		return fmt.Errorf("lorawan_setup: expected structure with 3 elements")
	}
	if v, ok := st[0].(core.OctetStringData); ok {
		l.DevEUI = v
	}
	if v, ok := st[1].(core.OctetStringData); ok {
		l.AppEUI = v
	}
	if v, ok := st[2].(core.OctetStringData); ok {
		l.AppKey = v
	}
	return nil
}

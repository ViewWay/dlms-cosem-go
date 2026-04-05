package core

import "fmt"

// CosemObject is the interface all COSEM interface classes implement.
type CosemObject interface {
	ClassID() uint16
	LogicalName() ObisCode
	Version() uint8
	Access(attr int) CosemAttributeAccess
}

// CosemAttributeAccess describes access to an attribute.
type CosemAttributeAccess struct {
	Access AccessMode
	Get    AccessLevel
	Set    AccessLevel
}

// AccessMode defines how an attribute can be accessed.
type AccessMode int

const (
	AccessNone      AccessMode = 0
	AccessRead      AccessMode = 1
	AccessWrite     AccessMode = 2
	AccessReadWrite AccessMode = 3
)

// AccessLevel defines the access level.
type AccessLevel int

const (
	AccessPublic AccessLevel = 0
	AccessLow    AccessLevel = 1
	AccessHigh   AccessLevel = 2
)

// AccessSelector defines well-known attribute selectors.
type AccessSelector int

const (
	AccessSelectorNoSelector           AccessSelector = 0
	AccessSelectorRangeDescriptor      AccessSelector = 1
	AccessSelectorEntryDescriptor      AccessSelector = 2
	AccessSelectorProfileGenericBuffer AccessSelector = 3
	AccessSelectorSelectiveAccess      AccessSelector = 4
)

// String returns the selector name.
func (a AccessSelector) String() string {
	switch a {
	case AccessSelectorNoSelector:
		return "no_selector"
	case AccessSelectorRangeDescriptor:
		return "range_descriptor"
	case AccessSelectorEntryDescriptor:
		return "entry_descriptor"
	case AccessSelectorProfileGenericBuffer:
		return "profile_generic_buffer"
	case AccessSelectorSelectiveAccess:
		return "selective_access"
	default:
		return fmt.Sprintf("unknown(%d)", a)
	}
}

// Class IDs for common COSEM interface classes.
const (
	ClassIDData               uint16 = 1
	ClassIDRegister           uint16 = 3
	ClassIDExtendedRegister   uint16 = 4
	ClassIDDemandRegister     uint16 = 5
	ClassIDRegisterActivation uint16 = 6
	ClassIDProfileGeneric     uint16 = 7
	ClassIDClock              uint16 = 8
	ClassIDScriptTable        uint16 = 9
	ClassIDSchedule           uint16 = 10
	ClassIDSpecialDaysTable   uint16 = 11
	ClassIDAssociationSN      uint16 = 12
	ClassIDAssociationLN      uint16 = 15
	ClassIDSAPAssignment      uint16 = 17
	ClassIDImageTransfer      uint16 = 18
	ClassIDIECLocalPortSetup  uint16 = 19
	ClassIDActivityCalendar   uint16 = 20
	ClassIDRegisterMonitor    uint16 = 21
	ClassIDSingleActionSchedule uint16 = 22
	ClassIDIECHDLCSetup       uint16 = 23
	ClassIDUtilityTables      uint16 = 26
	ClassIDValueTable         uint16 = 29
	ClassIDRegisterTable      uint16 = 61
	ClassIDCompactData        uint16 = 62
	ClassIDStatusMapping      uint16 = 63
	ClassIDPush               uint16 = 40
	ClassIDTCPUDPSetup        uint16 = 41
	ClassIDIPv4Setup          uint16 = 42
	ClassIDMACAddressSetup    uint16 = 43
	ClassIDPPPSetup           uint16 = 44
	ClassIDGPRSModemSetup     uint16 = 45
	ClassIDSecuritySetup      uint16 = 64
	ClassIDDisconnectControl  uint16 = 70
	ClassIDLimiter            uint16 = 71
	ClassIDMBusClient         uint16 = 72
	ClassIDSFKPLCMACSetup     uint16 = 82
	ClassIDIECPublicKey       uint16 = 90
	ClassIDNTPSetup           uint16 = 100
	ClassIDZigbeeSASStartup   uint16 = 101
	ClassIDSchcLPWAN          uint16 = 126
	ClassIDLorawanSetup       uint16 = 128
)

// String returns the class name for known class IDs.
func ClassIDName(id uint16) string {
	switch id {
	case ClassIDData:
		return "Data"
	case ClassIDRegister:
		return "Register"
	case ClassIDExtendedRegister:
		return "ExtendedRegister"
	case ClassIDDemandRegister:
		return "DemandRegister"
	case ClassIDProfileGeneric:
		return "ProfileGeneric"
	case ClassIDClock:
		return "Clock"
	case ClassIDScriptTable:
		return "ScriptTable"
	case ClassIDSchedule:
		return "Schedule"
	case ClassIDSpecialDaysTable:
		return "SpecialDaysTable"
	case ClassIDAssociationLN:
		return "AssociationLN"
	case ClassIDAssociationSN:
		return "AssociationSN"
	case ClassIDSecuritySetup:
		return "SecuritySetup"
	case ClassIDImageTransfer:
		return "ImageTransfer"
	case ClassIDIECLocalPortSetup:
		return "IECLocalPortSetup"
	case ClassIDIECHDLCSetup:
		return "IECHDLCSetup"
	case ClassIDTCPUDPSetup:
		return "TCPUDPSetup"
	case ClassIDIPv4Setup:
		return "IPv4Setup"
	case ClassIDMACAddressSetup:
		return "MACAddressSetup"
	case ClassIDPPPSetup:
		return "PPPSetup"
	case ClassIDGPRSModemSetup:
		return "GPRSModemSetup"
	case ClassIDDisconnectControl:
		return "DisconnectControl"
	case ClassIDLimiter:
		return "Limiter"
	case ClassIDMBusClient:
		return "MBusClient"
	case ClassIDLorawanSetup:
		return "LorawanSetup"
	case ClassIDNTPSetup:
		return "NTPSetup"
	case ClassIDPush:
		return "Push"
	case ClassIDUtilityTables:
		return "UtilityTables"
	case ClassIDRegisterTable:
		return "RegisterTable"
	case ClassIDCompactData:
		return "CompactData"
	case ClassIDStatusMapping:
		return "StatusMapping"
	case ClassIDActivityCalendar:
		return "ActivityCalendar"
	case ClassIDRegisterMonitor:
		return "RegisterMonitor"
	case ClassIDSingleActionSchedule:
		return "SingleActionSchedule"
	case ClassIDValueTable:
		return "ValueTable"
	case ClassIDIECPublicKey:
		return "IECPublicKey"
	case ClassIDSAPAssignment:
		return "SAPAssignment"
	default:
		return fmt.Sprintf("Unknown(0x%04x)", id)
	}
}

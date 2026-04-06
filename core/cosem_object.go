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

// Class IDs for common COSEM interface classes (Blue Book Ed.16).
const (
	ClassIDData                       uint16 = 1
	ClassIDRegister                   uint16 = 3
	ClassIDExtendedRegister           uint16 = 4
	ClassIDDemandRegister             uint16 = 5
	ClassIDRegisterActivation         uint16 = 6
	ClassIDProfileGeneric             uint16 = 7
	ClassIDClock                      uint16 = 8
	ClassIDScriptTable                uint16 = 9
	ClassIDSchedule                   uint16 = 10
	ClassIDSpecialDaysTable           uint16 = 11
	ClassIDAssociationSN              uint16 = 12
	ClassIDAssociationLN              uint16 = 15
	ClassIDSAPAssignment              uint16 = 17
	ClassIDImageTransfer              uint16 = 18
	ClassIDActivityCalendar           uint16 = 20
	ClassIDRegisterMonitor            uint16 = 21
	ClassIDSingleActionSchedule       uint16 = 22
	ClassIDIECHDLCSetup               uint16 = 23
	ClassIDIECTwistedPairSetup        uint16 = 24
	ClassIDMBusSlavePortSetup         uint16 = 25
	ClassIDUtilityTables              uint16 = 26
	ClassIDModemConfiguration         uint16 = 27
	ClassIDAutoAnswer                 uint16 = 28
	ClassIDAutoConnect                uint16 = 29
	ClassIDCOSEMDataProtection        uint16 = 30
	ClassIDPushSetup                  uint16 = 40
	ClassIDTCPUDPSetup                uint16 = 41
	ClassIDIPv4Setup                  uint16 = 42
	ClassIDMACAddressSetup            uint16 = 43
	ClassIDPPPSetup                   uint16 = 44
	ClassIDGPRSModemSetup             uint16 = 45
	ClassIDSMTPSetup                  uint16 = 46
	ClassIDGSMDiagnostic              uint16 = 47
	ClassIDIPv6Setup                  uint16 = 48
	ClassIDRegisterTable              uint16 = 61
	ClassIDCompactData                uint16 = 62
	ClassIDStatusMapping              uint16 = 63
	ClassIDSecuritySetup              uint16 = 64
	ClassIDParameterMonitor           uint16 = 65
	ClassIDMeasurementDataMonitoring  uint16 = 66
	ClassIDSensorManager              uint16 = 67
	ClassIDArbitrator                 uint16 = 68
	ClassIDDisconnectControl          uint16 = 70
	ClassIDLimiter                    uint16 = 71
	ClassIDMBusClient                 uint16 = 72
	ClassIDMBusMasterPortSetup        uint16 = 74
	ClassIDMBusDiagnostic             uint16 = 77
	ClassIDNTPSetup                   uint16 = 100
	ClassIDAccount                    uint16 = 111
	ClassIDCredit                     uint16 = 112
	ClassIDCharge                     uint16 = 113
	ClassIDTokenGateway               uint16 = 115
	ClassIDFunctionControl            uint16 = 122
	ClassIDArrayManager               uint16 = 123
	ClassIDCommPortProtection         uint16 = 124
	ClassIDSCHCLPWANSetup             uint16 = 126
	ClassIDLoRaWANSetup               uint16 = 128
	ClassIDLTEMonitoring              uint16 = 151
	ClassIDIEC62055Attributes          uint16 = 116
	ClassIDWirelessModeQChannel        uint16 = 73
	ClassIDDLMSMBusPortSetup           uint16 = 76
	ClassIDCoAPSetup                   uint16 = 152
	ClassIDCoAPDiagnostic              uint16 = 153
	ClassIDSFSKPhyMACSetup             uint16 = 50
	ClassIDSFSKActiveInitiator         uint16 = 51
	ClassIDSFSKMACSyncTimeouts         uint16 = 52
	ClassIDSFSKMACCounters             uint16 = 53
	ClassIDIEC61334LLCSetup            uint16 = 55
	ClassIDSFSKReportingSystemList     uint16 = 56
	ClassIDPrimeLLCSSCSSetup           uint16 = 80
	ClassIDPrimePhysicalCounters       uint16 = 81
	ClassIDPrimeMACSetup               uint16 = 82
	ClassIDPrimeMACFuncParams          uint16 = 83
	ClassIDPrimeMACCounters            uint16 = 84
	ClassIDPrimeMACNetworkAdmin        uint16 = 85
	ClassIDPrimeAppIdentification      uint16 = 86
	ClassIDG3MACCounters               uint16 = 90
	ClassIDG3MACSetup                  uint16 = 91
	ClassIDG36LoWPANSetup              uint16 = 92
	ClassIDG3HybridRFCounters          uint16 = 160
	ClassIDG3HybridRFSetup             uint16 = 161
	ClassIDG3Hybrid6LoWPANSetup        uint16 = 162
	ClassIDHSPLCMACSetup               uint16 = 140
	ClassIDHSPLCCPASSetup              uint16 = 141
	ClassIDHSPLCIPSSASSetup            uint16 = 142
	ClassIDHSPLCHDLCSSASSetup          uint16 = 143
	ClassIDLLCType1Setup               uint16 = 57
	ClassIDLLCType2Setup               uint16 = 58
	ClassIDLLCType3Setup               uint16 = 59
	ClassIDZigBeeSASStartup            uint16 = 101
	ClassIDZigBeeSASJoin               uint16 = 102
	ClassIDZigBeeSASAPSFragmentation    uint16 = 103
	ClassIDZigBeeNetworkControl        uint16 = 104
	ClassIDZigBeeTunnelSetup           uint16 = 105
	ClassIDSCHCLPWANDiagnostic         uint16 = 127
	ClassIDLoRaWANDiagnostic           uint16 = 129
	ClassIDWiSUNSetup                  uint16 = 95
	ClassIDWiSUNDiagnostic             uint16 = 96
	ClassIDRPLDiagnostic               uint16 = 97
	ClassIDMPLDiagnostic               uint16 = 98
	ClassIDIEC14908Identification       uint16 = 130
	ClassIDIEC14908ProtocolSetup        uint16 = 131
	ClassIDIEC14908ProtocolStatus       uint16 = 132
	ClassIDIEC14908Diagnostic          uint16 = 133

	// Non-standard / manufacturer-specific class IDs.
	ClassIDIECLocalPortSetup          uint16 = 19
	ClassIDIECPublicKey               uint16 = 256
	ClassIDValueDisplay                uint16 = 257
	ClassIDValueTable                  uint16 = 258
	ClassIDTariffSchedule              uint16 = 259
	ClassIDLpSetup                     uint16 = 260
	ClassIDIPv4TCPSetup                uint16 = 261

	// Legacy aliases kept for backwards compatibility.
	ClassIDPush               uint16 = ClassIDPushSetup
	ClassIDSmtpSetup          uint16 = ClassIDSMTPSetup
	ClassIDCosemDataProtection uint16 = ClassIDCOSEMDataProtection
	ClassIDMBusSlaveSetup      uint16 = ClassIDMBusSlavePortSetup
	ClassIDSchcLPWAN          uint16 = ClassIDSCHCLPWANSetup
	ClassIDLorawanSetup       uint16 = ClassIDLoRaWANSetup
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

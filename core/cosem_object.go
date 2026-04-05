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
	ClassIDData                  uint16 = 1
	ClassIDRegister              uint16 = 2
	ClassIDExtendedRegister      uint16 = 3
	ClassIDDemandRegister        uint16 = 4
	ClassIDProfileGeneric        uint16 = 5
	ClassIDClock                 uint16 = 6
	ClassIDScriptTable           uint16 = 7
	ClassIDSpecialDaysTable     uint16 = 8
	ClassIDSchedule              uint16 = 9
	ClassIDAssociationSN         uint16 = 10
	ClassIDAssociationLN         uint16 = 11
	ClassIDDemand                uint16 = 12
	ClassIDRegisterMonitor       uint16 = 13
	ClassIDRegisterActivation    uint16 = 14
	ClassIDPushSetup             uint16 = 15
	ClassIDDisconnectControl     uint16 = 16
	ClassIDLimiter               uint16 = 17
	ClassIDActivityCalendar      uint16 = 18
	ClassIDDayProfile            uint16 = 19
	ClassIDWeekProfile           uint16 = 20
	ClassIDTariffSchedule        uint16 = 21
	ClassIDAccount               uint16 = 22
	ClassIDModule                uint16 = 22
	ClassIDIECHDLCSetup          uint16 = 23
	ClassIDModemConfiguration    uint16 = 24
	ClassIDMBusMasterPortSetup   uint16 = 25
	ClassIDMBusClient            uint16 = 26
	ClassIDGPRSModemSetup        uint16 = 27
	ClassIDSmtpSetup             uint16 = 28
	ClassIDSecuritySetup         uint16 = 29
	ClassIDValueDisplay          uint16 = 30
	ClassIDLocalDisplay          uint16 = 31
	ClassIDIECPublicKey          uint16 = 32
	ClassIDCosemDataProtection   uint16 = 33
	ClassIDMBusSlaveSetup        uint16 = 34
	ClassIDImageTransfer         uint16 = 35
	ClassIDFirmwareManagement    uint16 = 36
	ClassIDMultiplier            uint16 = 37
	ClassIDMultiplierSetup       uint16 = 38
	ClassIDTCPUDPSetup           uint16 = 39
	ClassIDIPv4Setup             uint16 = 40
	ClassIDIPv4TCPSetup          uint16 = 41
	ClassIDIPv4UDPSetup          uint16 = 42
	ClassIDIPv6Setup             uint16 = 43
	ClassIDIPv6TCPSetup          uint16 = 44
	ClassIDPPPSetup              uint16 = 45
	ClassIDZigbeeSetup           uint16 = 46
	ClassIDWiSunSetup            uint16 = 47
	ClassIDAutoConnect           uint16 = 48
	ClassIDApplicationContext    uint16 = 49
	ClassIDTransport             uint16 = 50
	ClassIDRoute                 uint16 = 51
	ClassIDRouteSetup            uint16 = 52
	ClassIDSAPAssignment         uint16 = 53
	ClassIDStatusMapping         uint16 = 55
	ClassIDSupplyDisabling       uint16 = 56
	ClassIDSensorManager         uint16 = 57
	ClassIDActuatorSetup         uint16 = 58
	ClassIDSensor                uint16 = 59
	ClassIDSensorSetup           uint16 = 60
	ClassIDActuator              uint16 = 61
	ClassIDCompactData           uint16 = 62
	ClassIDStandardReadout       uint16 = 63
	ClassIDBilling               uint16 = 64
	ClassIDTokenGateway          uint16 = 65
	ClassIDCredit                uint16 = 66
	ClassIDCharge                uint16 = 67
	ClassIDUtilitySubSchedule    uint16 = 68
	ClassIDCluster               uint16 = 69
	ClassIDCalendar              uint16 = 70
	ClassIDSingleActionSchedule  uint16 = 71
	ClassIDScheduledActivity     uint16 = 72
	ClassIDClockControl          uint16 = 73
	ClassIDDataLogger            uint16 = 74
	ClassIDDataStorage           uint16 = 75
	ClassIDEventLog              uint16 = 76
	ClassIDEventLogger           uint16 = 77
	ClassIDFunctionControl       uint16 = 78
	ClassIDParameterMonitor      uint16 = 79
	ClassIDQualityOfService      uint16 = 80
	ClassIDLiftManagement        uint16 = 81
	ClassIDArbitrator            uint16 = 82
	ClassIDMBusDiagnostic         uint16 = 83
	ClassIDMplDiagnostic          uint16 = 84
	ClassIDRplDiagnostic          uint16 = 85
	ClassIDWiSunDiagnostic        uint16 = 86
	ClassIDNTPSetup              uint16 = 87
	ClassIDUps                   uint16 = 88
	ClassIDDirectDisconnect      uint16 = 89
	ClassIDLoadProfileSwitch     uint16 = 90
	ClassIDSerialPort            uint16 = 92
	ClassIDHanSetup              uint16 = 93
	ClassIDStatusDiag            uint16 = 94
	ClassIDExtendedProfile       uint16 = 95
	ClassIDMBusMaster            uint16 = 96
	ClassIDGasFlow               uint16 = 97
	ClassIDGasValve              uint16 = 98
	ClassIDWaterFlow             uint16 = 99
	ClassIDLpSetup               uint16 = 100
	ClassIDRs485Setup            uint16 = 101
	ClassIDInfraredSetup         uint16 = 102
	ClassIDEthernetSetup         uint16 = 103
	ClassIDLteSetup              uint16 = 104
	ClassIDTlsSetup              uint16 = 105
	ClassIDNbiotSetup            uint16 = 106
	ClassIDLorawanSetup          uint16 = 107
	ClassIDPowerQualityMonitor   uint16 = 110
	ClassIDHarmonicMonitor       uint16 = 111
	ClassIDSagSwellMonitor       uint16 = 112
	ClassIDHeatCostAllocator     uint16 = 113
	ClassIDMeasurementDataMonitoring uint16 = 114
	ClassIDCommPortProtection    uint16 = 115
	ClassIDIdentity              uint16 = 116
	ClassIDMACAddressSetup       uint16 = 117
	ClassIDLoadProfile           uint16 = 118
	ClassIDValueTable            uint16 = 119
	ClassIDAlarmHandler          uint16 = 120
	ClassIDUtilityTables         uint16 = 121
	ClassIDActivePowerImport     uint16 = 200
	ClassIDSinglePhaseImport     uint16 = 201
	ClassIDSinglePhaseExport     uint16 = 202
	ClassIDSinglePhaseMq         uint16 = 203
	ClassIDSinglePhase           uint16 = 204
	ClassIDMaximumDemand         uint16 = 205
	ClassIDTotal                 uint16 = 207
	ClassIDCommControl           uint16 = 210
	ClassIDArrayManager          uint16 = 211
	ClassIDSFKPLCMACSetup        uint16 = 126
	ClassIDSchcLPWAN             uint16 = 128
	ClassIDRegisterTable         uint16 = 61
	ClassIDPush                  uint16 = 15
	ClassIDIECLocalPortSetup     uint16 = 102
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

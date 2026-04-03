package core

// Common OBIS codes
var (
	ObisClock                  = MustParseObis("0.0.1.0.0.255")
	ObisActivePowerPlus        = MustParseObis("1.0.1.8.0.255")
	ObisActivePowerMinus       = MustParseObis("1.0.2.8.0.255")
	ObisReactivePowerPlus      = MustParseObis("1.0.3.8.0.255")
	ObisReactivePowerMinus     = MustParseObis("1.0.4.8.0.255")
	ObisVoltageL1              = MustParseObis("1.0.32.7.0.255")
	ObisVoltageL2              = MustParseObis("1.0.52.7.0.255")
	ObisVoltageL3              = MustParseObis("1.0.72.7.0.255")
	ObisCurrentL1              = MustParseObis("1.0.31.7.0.255")
	ObisCurrentL2              = MustParseObis("1.0.51.7.0.255")
	ObisCurrentL3              = MustParseObis("1.0.71.7.0.255")
	ObisPowerFactorL1          = MustParseObis("1.0.33.7.0.255")
	ObisPowerFactorL2          = MustParseObis("1.0.53.7.0.255")
	ObisPowerFactorL3          = MustParseObis("1.0.73.7.0.255")
	ObisFrequency              = MustParseObis("1.0.14.7.0.255")
	ObisLoadProfile            = MustParseObis("1.0.99.1.0.255")
	ObisSecuritySetup          = MustParseObis("0.0.43.0.0.255")
	ObisAssociationLN          = MustParseObis("0.0.40.0.0.255")
	ObisAssociationSN          = MustParseObis("0.0.41.0.0.255")
	ObisLPSetup                = MustParseObis("0.0.44.0.0.255")
	ObisRS485Setup             = MustParseObis("0.0.23.0.0.255")
	ObisInfraredSetup          = MustParseObis("0.0.19.0.0.255")
	ObisNBIoTSetup             = MustParseObis("0.0.45.0.0.255")
	ObisLoRaWANSetup           = MustParseObis("0.0.128.0.0.255")
	ObisTariffPlan             = MustParseObis("0.0.14.0.0.255")
	ObisTariffTable            = MustParseObis("0.0.14.0.1.255")
	ObisSeasonProfile          = MustParseObis("0.0.13.0.0.255")
	ObisWeekProfile            = MustParseObis("0.0.13.0.1.255")
	ObisDayProfile             = MustParseObis("0.0.13.0.2.255")
	ObisTotalActiveEnergy      = MustParseObis("1.0.1.8.0.255")
	ObisTotalReactiveEnergy    = MustParseObis("1.0.3.8.0.255")
	ObisPositiveActiveEnergy   = MustParseObis("1.0.1.8.1.255")
	ObisNegativeActiveEnergy   = MustParseObis("1.0.1.8.2.255")
	ObisPositiveReactiveEnergy = MustParseObis("1.0.3.8.1.255")
	ObisNegativeReactiveEnergy = MustParseObis("1.0.3.8.2.255")
	ObisDemandRegisterMax      = MustParseObis("1.0.13.7.0.255")
	ObisData                   = MustParseObis("0.0.96.1.0.255")
	ObisSerialNumber           = MustParseObis("0.0.96.1.1.255")
	ObisFirmwareVersion        = MustParseObis("0.0.96.1.2.255")
)

// MustParseObis parses an OBIS code string, panicking on error.
func MustParseObis(s string) ObisCode {
	oc, err := ParseObis(s)
	if err != nil {
		panic(err)
	}
	return oc
}

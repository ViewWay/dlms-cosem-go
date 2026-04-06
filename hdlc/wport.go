package hdlc

import "fmt"

// wPort constants for DLMS/COSEM (Green Book 7.3.3.4)
//
// Wrapper port numbers (wPort) are used for addressing DLMS/COSEM
// Application Entities in UDP and TCP transport layers.

const (
	// Reserved wrapper port numbers
	WPortDlmsCosemUDP = 4059 // DLMS/COSEM UDP
	WPortDlmsCosemTCP = 4059 // DLMS/COSEM TCP

	// Client side reserved addresses
	WPortNoStation           = 0x0000 // No-station
	WPortClientMgmtProcess   = 0x0001 // Client Management Process
	WPortPublicClient       = 0x0010 // Public Client

	// Server side reserved addresses
	WPortMgmtLogicalDevice = 0x0001 // Management Logical Device
	WPortAllStation        = 0x007F // All-station (Broadcast)
)

// IsReservedWPort checks if a wPort number is reserved
func IsReservedWPort(wport int) bool {
	reservedPorts := []int{
		WPortNoStation,
		WPortClientMgmtProcess,
		WPortPublicClient,
		WPortMgmtLogicalDevice,
		WPortAllStation,
	}

	for _, port := range reservedPorts {
		if port == wport {
			return true
		}
	}

	return false
}

// GetWPortDescription returns a description of a wPort number
func GetWPortDescription(wport int) string {
	// Special handling for 4059 (used for both UDP and TCP)
	if wport == 4059 {
		return "DLMS/COSEM UDP/TCP"
	}

	// Special handling for 0x0001 (used by both client and server)
	if wport == 0x0001 {
		return "Client Management Process / Management Logical Device"
	}

	descriptions := map[int]string{
		WPortNoStation:   "No-station",
		WPortPublicClient: "Public Client",
		WPortAllStation:  "All-station (Broadcast)",
	}

	if desc, ok := descriptions[wport]; ok {
		return desc
	}

	return fmt.Sprintf("Unknown (0x%04X)", wport)
}

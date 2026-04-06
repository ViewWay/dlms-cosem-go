package core

// DLMS/COSEM protocol constants

const (
	// DLMSUDP_PORT is the DLMS/COSEM UDP reserved port (Green Book 7.3.3.4)
	// Port 4059 is the IANA-assigned port for DLMS/COSEM over UDP.
	DLMSUDP_PORT uint16 = 4059

	// DLMSTCP_PORT is the DLMS/COSEM TCP default port
	// While not officially reserved, port 4059 is commonly used for
	// DLMS/COSEM over TCP as well.
	DLMSTCP_PORT uint16 = 4059
)

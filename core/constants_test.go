package core

import "testing"

func TestDLMSConstants(t *testing.T) {
	if DLMSUDP_PORT != 4059 {
		t.Errorf("DLMSUDP_PORT should be 4059, got %d", DLMSUDP_PORT)
	}

	if DLMSTCP_PORT != 4059 {
		t.Errorf("DLMSTCP_PORT should be 4059, got %d", DLMSTCP_PORT)
	}
}

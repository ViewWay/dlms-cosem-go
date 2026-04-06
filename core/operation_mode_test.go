package core

import (
	"testing"
)

func TestOperationMode_AllowsWrite(t *testing.T) {
	tests := []struct {
		mode     OperationMode
		expected bool
	}{
		{ModeReadout, false},
		{ModeProgramming, true},
	}

	for _, tt := range tests {
		t.Run(tt.mode.String(), func(t *testing.T) {
			if got := tt.mode.AllowsWrite(); got != tt.expected {
				t.Errorf("OperationMode.AllowsWrite() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestOperationMode_AllowsRead(t *testing.T) {
	tests := []struct {
		mode     OperationMode
		expected bool
	}{
		{ModeReadout, true},
		{ModeProgramming, true},
	}

	for _, tt := range tests {
		t.Run(tt.mode.String(), func(t *testing.T) {
			if got := tt.mode.AllowsRead(); got != tt.expected {
				t.Errorf("OperationMode.AllowsRead() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCheckWritePermission(t *testing.T) {
	// Test readout mode (should fail)
	err := CheckWritePermission(ModeReadout)
	if err == nil {
		t.Error("CheckWritePermission(ModeReadout) should return error")
	}

	// Test programming mode (should succeed)
	err = CheckWritePermission(ModeProgramming)
	if err != nil {
		t.Errorf("CheckWritePermission(ModeProgramming) should not return error, got %v", err)
	}
}

func TestSwitchMode(t *testing.T) {
	mode := ModeReadout
	SwitchMode(&mode, ModeProgramming)
	if mode != ModeProgramming {
		t.Errorf("SwitchMode failed, got %v, want %v", mode, ModeProgramming)
	}
}

func TestValidateOperation(t *testing.T) {
	// Test read operations (should always succeed)
	if err := ValidateOperation(ModeReadout, "read"); err != nil {
		t.Errorf("ValidateOperation(ModeReadout, \"read\") failed: %v", err)
	}
	if err := ValidateOperation(ModeProgramming, "read"); err != nil {
		t.Errorf("ValidateOperation(ModeProgramming, \"read\") failed: %v", err)
	}

	// Test write operations
	if err := ValidateOperation(ModeReadout, "write"); err == nil {
		t.Error("ValidateOperation(ModeReadout, \"write\") should fail")
	}
	if err := ValidateOperation(ModeProgramming, "write"); err != nil {
		t.Errorf("ValidateOperation(ModeProgramming, \"write\") failed: %v", err)
	}
}

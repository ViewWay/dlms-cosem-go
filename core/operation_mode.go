package core

import (
	"errors"
	"fmt"
)

// OperationMode represents the DLMS/COSEM operation mode (Green Book 6.2-6.3)
type OperationMode int

const (
	// ModeReadout: Read-only access with HLS password authentication
	// In this mode, only read operations are allowed.
	ModeReadout OperationMode = iota

	// ModeProgramming: Read-write access with HLS-SLS signature authentication
	// In this mode, both read and write operations are allowed.
	ModeProgramming
)

// String returns the string representation of the operation mode
func (m OperationMode) String() string {
	return [...]string{
		"Readout",
		"Programming",
	}[m]
}

// AllowsWrite checks if write operations are allowed in this mode
func (m OperationMode) AllowsWrite() bool {
	return m == ModeProgramming
}

// AllowsRead checks if read operations are allowed in this mode
func (m OperationMode) AllowsRead() bool {
	return true // Both modes allow read operations
}

// OperationModeError represents an error related to operation mode violations
type OperationModeError struct {
	Message string
}

func (e *OperationModeError) Error() string {
	return fmt.Sprintf("operation mode error: %s", e.Message)
}

// NewOperationModeError creates a new operation mode error
func NewOperationModeError(message string) *OperationModeError {
	return &OperationModeError{Message: message}
}

// CheckWritePermission checks if write operations are allowed in the current mode
// Returns an error if write operations are not allowed
func CheckWritePermission(mode OperationMode) error {
	if !mode.AllowsWrite() {
		return NewOperationModeError(
			fmt.Sprintf("write operations not allowed in %s mode", mode),
		)
	}
	return nil
}

// SwitchMode switches the operation mode to a new mode
func SwitchMode(currentMode *OperationMode, newMode OperationMode) {
	*currentMode = newMode
}

// ValidateOperation validates an operation against the current operation mode
// operation should be "read" or "write"
func ValidateOperation(mode OperationMode, operation string) error {
	switch operation {
	case "read":
		if !mode.AllowsRead() {
			return NewOperationModeError("read operations not allowed")
		}
		return nil
	case "write":
		return CheckWritePermission(mode)
	default:
		return errors.New("unknown operation type")
	}
}

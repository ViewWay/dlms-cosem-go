package core

import (
	"fmt"
	"strconv"
	"strings"
)

// ObisCode represents a 6-byte OBIS (Object Identification System) code.
type ObisCode [6]byte

// ParseObis parses a string like "1.0.0.1.0.255" into an ObisCode.
func ParseObis(s string) (ObisCode, error) {
	var oc ObisCode
	parts := strings.Split(s, ".")
	if len(parts) != 5 && len(parts) != 6 {
		return oc, fmt.Errorf("obis: expected 5 or 6 parts, got %d", len(parts))
	}
	if len(parts) == 5 {
		parts = append(parts, "255")
	}
	for i := 0; i < 6; i++ {
		v, err := strconv.Atoi(parts[i])
		if err != nil || v < 0 || v > 255 {
			return oc, fmt.Errorf("obis: invalid part %q at position %d", parts[i], i)
		}
		oc[i] = byte(v)
	}
	return oc, nil
}

// ObisFromBytes creates an ObisCode from 6 bytes.
func ObisFromBytes(b []byte) (ObisCode, error) {
	if len(b) != 6 {
		return ObisCode{}, fmt.Errorf("obis: expected 6 bytes, got %d", len(b))
	}
	var oc ObisCode
	copy(oc[:], b)
	return oc, nil
}

// String returns the OBIS code as "A.B.C.D.E.F".
func (o ObisCode) String() string {
	return fmt.Sprintf("%d.%d.%d.%d.%d.%d", o[0], o[1], o[2], o[3], o[4], o[5])
}

// Bytes returns the 6-byte representation.
func (o ObisCode) Bytes() []byte {
	return o[:]
}

// Compare returns -1, 0, or 1.
func (o ObisCode) Compare(other ObisCode) int {
	for i := 0; i < 6; i++ {
		if o[i] < other[i] {
			return -1
		}
		if o[i] > other[i] {
			return 1
		}
	}
	return 0
}

// A-F accessors
func (o ObisCode) A() byte { return o[0] }
func (o ObisCode) B() byte { return o[1] }
func (o ObisCode) C() byte { return o[2] }
func (o ObisCode) D() byte { return o[3] }
func (o ObisCode) E() byte { return o[4] }
func (o ObisCode) F() byte { return o[5] }

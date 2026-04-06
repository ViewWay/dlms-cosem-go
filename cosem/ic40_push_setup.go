package cosem

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/ViewWay/dlms-cosem-go/core"
)

// PushObject represents an object in push list
type PushObject struct {
	ClassID     uint16
	LogicalName core.ObisCode
	Attribute   uint8
	DataIndex   uint8
}

// PushSetup is the COSEM Push Setup interface class (IC 15).
type PushSetup struct {
	LogicalName                core.ObisCode
	PushObjectList             []PushObject
	Service                    uint16
	Destination                []byte
	CommunicationWindow        core.CosemDateTime
	RandomisationStartInterval uint16
	NumberOfRetries            uint8
	RepetitionDelay            uint16
	Version                    uint8
}

func (ps *PushSetup) ClassID() uint16               { return core.ClassIDPush }
func (ps *PushSetup) GetLogicalName() core.ObisCode { return ps.LogicalName }
func (ps *PushSetup) GetVersion() uint8             { return ps.Version }
func (ps *PushSetup) Access(attr int) core.CosemAttributeAccess {
	switch attr {
	case 1:
		return core.CosemAttributeAccess{Access: core.AccessRead}
	case 2, 3, 4, 5, 6, 7, 8:
		return core.CosemAttributeAccess{Access: core.AccessReadWrite}
	default:
		return core.CosemAttributeAccess{Access: core.AccessNone}
	}
}

func (ps *PushSetup) MarshalBinary() ([]byte, error) {
	arr := core.ArrayData{}
	for _, o := range ps.PushObjectList {
		arr = append(arr, core.StructureData{
			core.UnsignedLongData(o.ClassID),
			core.OctetStringData(o.LogicalName[:]),
			core.IntegerData(int8(o.Attribute)),
			core.UnsignedIntegerData(o.DataIndex),
		})
	}
	return arr.ToBytes(), nil
}

func (ps *PushSetup) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	arr, ok := elem.(core.ArrayData)
	if !ok {
		return fmt.Errorf("push_setup: expected array")
	}
	ps.PushObjectList = make([]PushObject, len(arr))
	for i, item := range arr {
		st, ok := item.(core.StructureData)
		if !ok || len(st) < 4 {
			return fmt.Errorf("push_setup: invalid entry at index %d", i)
		}
		ps.PushObjectList[i].ClassID = uint16(st[0].(core.UnsignedLongData))
		ps.PushObjectList[i].LogicalName = core.ObisCode(st[1].(core.OctetStringData))
		ps.PushObjectList[i].Attribute = uint8(st[2].(core.IntegerData))
		ps.PushObjectList[i].DataIndex = uint8(st[3].(core.UnsignedIntegerData))
	}
	return nil
}

func (ps *PushSetup) AddObject(obj PushObject) {
	ps.PushObjectList = append(ps.PushObjectList, obj)
}

func (ps *PushSetup) ObjectCount() int {
	return len(ps.PushObjectList)
}

// SendDestinationAndMethod combines Destination (SAP) and Service.
// In a real implementation this would resolve to a transport endpoint.
type SendDestinationAndMethod struct {
	SAP      []byte
	Service  uint16
}

// PushResult captures the outcome of a push attempt.
type PushResult struct {
	Success bool
	Attempt int
	Error   error
}

// ExecutePush triggers a push: serializes the PushObjectList and returns
// the encoded payload along with the SendDestinationAndMethod metadata.
// The caller is responsible for actual transport delivery.
func (ps *PushSetup) ExecutePush() ([]byte, *SendDestinationAndMethod, error) {
	if len(ps.PushObjectList) == 0 {
		return nil, nil, fmt.Errorf("push_setup: empty object list")
	}
	payload, err := ps.MarshalBinary()
	if err != nil {
		return nil, nil, fmt.Errorf("push_setup: serialize failed: %w", err)
	}
	dest := &SendDestinationAndMethod{
		SAP:     ps.Destination,
		Service: ps.Service,
	}
	return payload, dest, nil
}

// IsWithinWindow checks whether the current time falls within the
// RandomisationStartInterval window before the scheduled push time.
// Returns true if now is within [windowEnd - interval, windowEnd].
func (ps *PushSetup) IsWithinWindow(now, windowEnd time.Time) bool {
	interval := time.Duration(ps.RandomisationStartInterval) * time.Second
	if interval <= 0 {
		return true
	}
	start := windowEnd.Add(-interval)
	return !now.Before(start) && !now.After(windowEnd)
}

// ShouldRetry returns true if the attempt count has not exceeded NumberOfRetries.
// attempt is 1-based (first call = attempt 1).
func (ps *PushSetup) ShouldRetry(attempt int) bool {
	if attempt < 1 {
		return false
	}
	return attempt <= int(ps.NumberOfRetries)
}

// GetRetryDelay calculates the delay before the next retry attempt.
// Uses RepetitionDelay with a random jitter in [0, RepetitionDelay) seconds.
func (ps *PushSetup) GetRetryDelay(attempt int) time.Duration {
	base := time.Duration(ps.RepetitionDelay) * time.Second
	if base <= 0 {
		return time.Second // default 1s
	}
	jitter := time.Duration(rand.Intn(int(base))) * time.Second
	return base + jitter
}

package core

import (
	"encoding/binary"
	"fmt"
	"time"
)

// CosemDateTime represents a DLMS/COSEM DateTime (12 bytes).
type CosemDateTime struct {
	Year        uint16
	Month       uint8
	Day         uint8
	DayOfWeek   uint8
	Hour        uint8
	Minute      uint8
	Second      uint8
	Hundredths  uint8
	Deviation   int32
	ClockStatus uint8
}

// CosemDateTimeFromBytes parses 12 bytes into a CosemDateTime.
func CosemDateTimeFromBytes(data []byte) (CosemDateTime, error) {
	if len(data) < 12 {
		return CosemDateTime{}, fmt.Errorf("datetime: expected 12 bytes, got %d", len(data))
	}
	return CosemDateTime{
		Year:        binary.BigEndian.Uint16(data[0:2]),
		Month:       data[2],
		Day:         data[3],
		DayOfWeek:   data[4],
		Hour:        data[5],
		Minute:      data[6],
		Second:      data[7],
		Hundredths:  data[8],
		Deviation:   int32(binary.BigEndian.Uint16(data[9:11])),
		ClockStatus: data[11],
	}, nil
}

// ToBytes encodes the CosemDateTime to 12 bytes.
func (dt CosemDateTime) ToBytes() []byte {
	buf := make([]byte, 12)
	binary.BigEndian.PutUint16(buf[0:2], dt.Year)
	buf[2] = dt.Month
	buf[3] = dt.Day
	buf[4] = dt.DayOfWeek
	buf[5] = dt.Hour
	buf[6] = dt.Minute
	buf[7] = dt.Second
	buf[8] = dt.Hundredths
	binary.BigEndian.PutUint16(buf[9:11], uint16(dt.Deviation))
	buf[11] = dt.ClockStatus
	return buf
}

// IsInvalid returns true if the date/time is not specified.
func (dt CosemDateTime) IsInvalid() bool {
	return dt.Year == 0xFFFF && dt.Month == 0xFF && dt.Day == 0xFF
}

// ToTime converts to time.Time (zero time if invalid).
func (dt CosemDateTime) ToTime() time.Time {
	if dt.IsInvalid() {
		return time.Time{}
	}
	loc := time.UTC
	if dt.Deviation != 0x8000 {
		offset := int(dt.Deviation)
		loc = time.FixedZone("", offset*60)
	}
	return time.Date(int(dt.Year), time.Month(dt.Month), int(dt.Day),
		int(dt.Hour), int(dt.Minute), int(dt.Second), int(dt.Hundredths)*1e7, loc)
}

// CosemDate represents a DLMS/COSEM Date (5 bytes).
type CosemDate struct {
	Year      uint16
	Month     uint8
	Day       uint8
	DayOfWeek uint8
	Deviation int32
}

// CosemDateFromBytes parses 5 bytes into a CosemDate.
func CosemDateFromBytes(data []byte) CosemDate {
	if len(data) < 5 {
		return CosemDate{Year: 0xFFFF, Month: 0xFF, Day: 0xFF, Deviation: 0x8000}
	}
	return CosemDate{
		Year:      binary.BigEndian.Uint16(data[0:2]),
		Month:     data[2],
		Day:       data[3],
		DayOfWeek: data[4],
	}
}

// ToBytes encodes the CosemDate to 5 bytes.
func (d CosemDate) ToBytes() []byte {
	buf := make([]byte, 5)
	binary.BigEndian.PutUint16(buf[0:2], d.Year)
	buf[2] = d.Month
	buf[3] = d.Day
	buf[4] = d.DayOfWeek
	return buf
}

// IsInvalid returns true if the date is not specified.
func (d CosemDate) IsInvalid() bool {
	return d.Year == 0xFFFF && d.Month == 0xFF && d.Day == 0xFF
}

// ToTime converts to time.Time.
func (d CosemDate) ToTime() time.Time {
	if d.IsInvalid() {
		return time.Time{}
	}
	return time.Date(int(d.Year), time.Month(d.Month), int(d.Day), 0, 0, 0, 0, time.UTC)
}

// CosemTime represents a DLMS/COSEM Time (4 bytes).
type CosemTime struct {
	Hour       uint8
	Minute     uint8
	Second     uint8
	Hundredths uint8
}

// CosemTimeFromBytes parses 4 bytes into a CosemTime.
func CosemTimeFromBytes(data []byte) CosemTime {
	if len(data) < 4 {
		return CosemTime{Hour: 0xFF, Minute: 0xFF, Second: 0xFF, Hundredths: 0xFF}
	}
	return CosemTime{
		Hour:       data[0],
		Minute:     data[1],
		Second:     data[2],
		Hundredths: data[3],
	}
}

// ToBytes encodes the CosemTime to 4 bytes.
func (t CosemTime) ToBytes() []byte {
	return []byte{t.Hour, t.Minute, t.Second, t.Hundredths}
}

// IsInvalid returns true if the time is not specified.
func (t CosemTime) IsInvalid() bool {
	return t.Hour == 0xFF && t.Minute == 0xFF && t.Second == 0xFF
}

// ToTimeDuration converts to a time.Duration from midnight.
func (t CosemTime) ToTimeDuration() time.Duration {
	if t.IsInvalid() {
		return 0
	}
	return time.Duration(t.Hour)*time.Hour +
		time.Duration(t.Minute)*time.Minute +
		time.Duration(t.Second)*time.Second +
		time.Duration(t.Hundredths)*10*time.Millisecond
}

package dpt

import (
	"fmt"
	"time"
)

// DPT_19001 represents DPT 19.001 / DateTime p. 43.
type DPT_19001 struct {
	Year       uint16 // 0 = 1900, 255 = 2155
	Month      uint8  // 1 ... 12
	DayOfMonth uint8  // 1 ... 31
	DayOfWeek  uint8  // 0 = any day, 1 = Monday ... 7 = Sunday
	HourOfDay  uint8  // 0 ... 24
	Minutes    uint8  // 0 ... 59
	Seconds    uint8  // 0 ... 59
	F          bool   // Fault (0 = no fault, 1 = fault)
	WD         bool   // Working Day (0 = bank day, 1 = working day)
	NWD        bool   // No Working Day (0 = WD field valid, 1 = WD field not valid)
	NY         bool   // No Year (0 = Year field valid, 1 = Year field not valid)
	ND         bool   // No Date (0 = Month and Day of Month fields valid, 1 = Month and Day of Month fields not valid)
	NDoW       bool   // No Day of Week (0 = Day of week field valid, 1 = Day of week field not valid)
	NT         bool   // No Time (0 = Hour of day, Minutes and Seconds fields valid, 1 = Hour of day, Minutes and Seconds fields not valid)
	SUTI       bool   // Standard Summer Time (0 = UT+X, 1 = UT+X+1)
	CLQ        bool   // Quality of Clock (0 = clock without ext. sync signal, 1 = clock with ext. sync signal)
}

func boolTOuint8(value bool) uint8 {
	if value {
		return 1
	}
	return 0
}

func unit8TObool(value uint8) bool {
	return value == 1
}

func (d DPT_19001) Pack() []byte {
	var buf = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0}
	if d.IsValid() {
		buf[1] = uint8(d.Year - 1900)
		buf[2] = d.Month & 0xF
		buf[3] = d.DayOfMonth & 0x1F
		buf[4] = d.DayOfWeek<<5 | (d.HourOfDay & 0x1F)
		buf[5] = d.Minutes & 0x3F
		buf[6] = d.Seconds & 0x3F
		buf[7] = boolTOuint8(d.F) & 0x1 << 7
		buf[7] |= boolTOuint8(d.WD) & 0x1 << 6
		buf[7] |= boolTOuint8(d.NWD) & 0x1 << 5
		buf[7] |= boolTOuint8(d.NY) & 0x1 << 4
		buf[7] |= boolTOuint8(d.ND) & 0x1 << 3
		buf[7] |= boolTOuint8(d.NDoW) & 0x1 << 2
		buf[7] |= boolTOuint8(d.NT) & 0x1 << 1
		buf[7] |= boolTOuint8(d.SUTI) & 0x1
		buf[8] = boolTOuint8(d.CLQ) & 0x1 << 7
	}
	return buf
}

func (d *DPT_19001) Unpack(data []byte) error {
	if len(data) != 9 {
		return ErrInvalidLength
	}

	d.Year = uint16(data[1]&0xFF) + 1900
	d.Month = data[2] & 0xF
	d.DayOfMonth = data[3] & 0x1F
	d.DayOfWeek = uint8(data[4] >> 5 & 0x07)
	d.HourOfDay = uint8(data[4] & 0x1F)
	d.Minutes = uint8(data[5] & 0x3F)
	d.Seconds = uint8(data[6] & 0x3F)
	d.F = unit8TObool(data[7] >> 7 & 0x1)
	d.WD = unit8TObool(data[7] >> 6 & 0x1)
	d.NWD = unit8TObool(data[7] >> 5 & 0x1)
	d.NY = unit8TObool(data[7] >> 4 & 0x1)
	d.ND = unit8TObool(data[7] >> 3 & 0x1)
	d.NDoW = unit8TObool(data[7] >> 2 & 0x1)
	d.NT = unit8TObool(data[7] >> 1 & 0x1)
	d.SUTI = unit8TObool(data[7] & 0x1)
	d.CLQ = unit8TObool(data[8] >> 7 & 0x1)

	if !d.IsValid() {
		return fmt.Errorf("payload is out of range")
	}

	return nil
}

func (d DPT_19001) Unit() string {
	return ""
}

func (d DPT_19001) IsValid() bool {
	tm := time.Date(int(d.Year), time.Month(d.Month), int(d.DayOfMonth), int(d.HourOfDay), int(d.Minutes), int(d.Seconds), 0, time.UTC)
	if tm.Year() < 1900 || tm.Year() > 2155 {
		// fmt.Println("year not valid")
		return false
	}
	if int(d.Year) != tm.Year() {
		// fmt.Println("year not valid 2")
		return false
	}
	if time.Month(d.Month) != tm.Month() {
		// fmt.Println("month not valid")
		return false
	}
	if int(d.DayOfMonth) != tm.Day() {
		// fmt.Println("day of month not valid")
		return false
	}
	if d.DayOfWeek > 7 {
		// fmt.Println("day of week not valid")
		return false
	}
	if d.HourOfDay > 24 {
		// fmt.Println("hour of day not valid")
		return false
	}
	// if hour is 24, minutes and second need to be 0
	if d.HourOfDay == 24 && (d.Minutes != 0 || d.Seconds != 0) {
		// fmt.Println("hour of day not valid 2")
		return false
	}
	if d.Minutes > 59 {
		// fmt.Println("minutes not valid")
		return false
	}
	if d.Seconds > 59 {
		// fmt.Println("seconds not valid")
		return false
	}
	// TODO finish validation of flags

	return true
}

func (d DPT_19001) String() string {
	timeString := ""
	weekday := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	if 0 < d.DayOfWeek && d.DayOfWeek <= 7 {
		timeString = fmt.Sprintf("%s %02d:%02d:%02d", weekday[d.DayOfWeek-1], d.HourOfDay, d.Minutes, d.Seconds)
	} else {
		timeString = fmt.Sprintf("%02d:%02d:%02d", d.HourOfDay, d.Minutes, d.Seconds)
	}
	return fmt.Sprintf("%04d-%02d-%02d %s", d.Year, d.Month, d.DayOfMonth, timeString)
}

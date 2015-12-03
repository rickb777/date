// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package date

import (
	"time"
	"fmt"
	"strconv"
)

const zero time.Duration = 0

// Clock specifies a time of day. It extends the existing time.Duration, applying
// that to the time since midnight on some arbitrary day.
//
// It is not intended that Clock be used to represent periods greater than 24 hours nor
// negative values. However, for such lengths of time, a fixed 24 hours per day
// is assumed and a modulo operation Mod() is provided to discard whole multiples of 24 hours.
//
// See https://en.wikipedia.org/wiki/ISO_8601#Times
type Clock time.Duration

const (
// ClockDay is a fixed period of 24 hours. This does not take account of daylight savings, so is not fully general.
	ClockDay    Clock = Clock(time.Hour * 24)
	ClockHour   Clock = Clock(time.Hour)
	ClockMinute Clock = Clock(time.Minute)
	ClockSecond Clock = Clock(time.Second)
)

// HhMmSs returns a new Clock with specified hour, minute, second.
func HhMmSs(h, m, s int) Clock {
	hns := Clock(h) * ClockHour
	mns := Clock(m) * ClockMinute
	sns := Clock(s) * ClockSecond
	return Clock(hns + mns + sns)
}

// Add returns a new Clock offset from this clock specified hour, minute, second. The parameters can be negative.
// If required, use Mod() to correct any overflow or underflow.
func (c Clock) Add(h, m, s int) Clock {
	hns := Clock(h) * ClockHour
	mns := Clock(m) * ClockMinute
	sns := Clock(s) * ClockSecond
	return c + hns + mns + sns
}

// ParseClock converts a string representation to a Clock. Acceptable representations
// are as per ISO-8601 - see https://en.wikipedia.org/wiki/ISO_8601#Times
func ParseClock(hms string) (clock Clock, err error) {
	switch len(hms) {
	case 2: // HH
		return parseClockParts(hms, hms, "", "", "")

	case 4: // HHMM
		return parseClockParts(hms, hms[:2], hms[2:], "", "")

	case 5: // HH:MM
		if hms[2] != ':' {
			return 0, fmt.Errorf("date.ParseClock: cannot parse %s", hms)
		}
		return parseClockParts(hms, hms[:2], hms[3:], "", "")

	case 6: // HHMMSS
		return parseClockParts(hms, hms[:2], hms[2:4], hms[4:], "")

	case 8: // HH:MM:SS
		if hms[2] != ':' || hms[5] != ':' {
			return 0, fmt.Errorf("date.ParseClock: cannot parse %s", hms)
		}
		return parseClockParts(hms, hms[:2], hms[3:5], hms[6:], "")

	default:
		if hms[2] != ':' || hms[5] != ':' || hms[8] != '.' {
			return 0, fmt.Errorf("date.ParseClock: cannot parse %s", hms)
		}
		return parseClockParts(hms, hms[:2], hms[3:5], hms[6:8], hms[9:])
	}
	return 0, fmt.Errorf("date.ParseClock: cannot parse %s", hms)
}

func parseClockParts(hms, hh, mm, ss, nnnns string) (clock Clock, err error) {
	h := 0
	m := 0
	s := 0
	ns := 0
	if hh != "" {
		h, err = strconv.Atoi(hh)
		if err != nil {
			return 0, fmt.Errorf("date.ParseClock: cannot parse %s: %v", hms, err)
		}
	}
	if mm != "" {
		m, err = strconv.Atoi(mm)
		if err != nil {
			return 0, fmt.Errorf("date.ParseClock: cannot parse %s: %v", hms, err)
		}
	}
	if ss != "" {
		s, err = strconv.Atoi(ss)
		if err != nil {
			return 0, fmt.Errorf("date.ParseClock: cannot parse %s: %v", hms, err)
		}
	}
	if nnnns != "" {
		ns, err = strconv.Atoi(nnnns)
		if err != nil {
			return 0, fmt.Errorf("date.ParseClock: cannot parse %s: %v", hms, err)
		}
	}
	return HhMmSs(h, m, s) + Clock(ns), nil
}

// IsInOneDay tests whether a clock time is in the range 0 to 24 hours, inclusive. Inside this
// range, a Clock is generally well-behaved. But outside it, there may be errors due to daylight
// savings. Note that 24:00:00 is included as a special case as per ISO-8601 definition of midnight.
func (c Clock) IsInOneDay() bool {
	return 0 <= c && c <= ClockDay
}

// Days gets the number of whole days represented by the Clock, assuming that each day is a fixed
// 24 hour period. Negative values are treated so that the range -23h59m59s to -1s is fully
// enclosed in a day numbered -1, and so on. This means that the result is zero only for the
// clock range 0s to 23h59m59s, for which IsInOneDay() returns true.
func (c Clock) Days() int {
	if c < 0 {
		return int(c / ClockDay) - 1
	} else {
		return int(c / ClockDay)
	}
}

// Mod24 calculates the remainder vs 24 hours using Euclidean division, in which the result
// will be less than 24 hours and is never negative.
// https://en.wikipedia.org/wiki/Modulo_operation
func (c Clock) Mod24() Clock {
	if 0 <= c && c < ClockDay {
		return c
	}
	if c < 0 {
		q := 1 - c / ClockDay
		return c + (q * ClockDay)
	}
	q := c / ClockDay
	return c - (q * ClockDay)
}

// Hours gets the clock-face number of hours (calculated from the modulo time, see Mod24).
func (c Clock) Hours() int {
	return int(clockHours(c.Mod24()))
}

// Minutes gets the clock-face number of minutes (calculated from the modulo time, see Mod24).
// For example, for 22:35 this will return 35.
func (c Clock) Minutes() int {
	return int(clockMinutes(c.Mod24()))
}

// Seconds gets the clock-face number of seconds (calculated from the modulo time, see Mod24).
// For example, for 10:20:30 this will return 30.
func (c Clock) Seconds() int {
	return int(clockSeconds(c.Mod24()))
}

// Nanosec gets the clock-face number of nanoseconds (calculated from the modulo time, see Mod24).
// For example, for 10:20:30.456111222 this will return 456111222.
func (c Clock) Nanosec() int64 {
	return int64(clockNanosec(c.Mod24()))
}

func clockHours(cm Clock) Clock {
	return (cm / ClockHour)
}

func clockMinutes(cm Clock) Clock {
	return (cm - clockHours(cm) * ClockHour) / ClockMinute
}

func clockSeconds(cm Clock) Clock {
	return (cm - clockHours(cm) * ClockHour - clockMinutes(cm) * ClockMinute) / ClockSecond
}

func clockNanosec(cm Clock) Clock {
	return cm - clockHours(cm) * ClockHour - clockMinutes(cm) * ClockMinute - clockSeconds(cm) * ClockSecond
}

// Hh gets the clock-face number of hours as a two-digit string (calculated from the modulo time, see Mod24).
func (c Clock) Hh() string {
	cm := c.Mod24()
	return fmt.Sprintf("%02d", clockHours(cm))
}

// HhMm gets the clock-face number of hours and minutes as a five-digit string (calculated from the
// modulo time, see Mod24).
func (c Clock) HhMm() string {
	cm := c.Mod24()
	return fmt.Sprintf("%02d:%02d", clockHours(cm), clockMinutes(cm))
}

// HhMmSs gets the clock-face number of hours, minutes, seconds as an eight-digit string
// (calculated from the modulo time, see Mod24).
func (c Clock) HhMmSs() string {
	cm := c.Mod24()
	return fmt.Sprintf("%02d:%02d:%02d", clockHours(cm), clockMinutes(cm), clockSeconds(cm))
}

// String gets the clock-face number of hours, minutes, seconds and nanoseconds as an 18-digit string
// (calculated from the modulo time, see Mod24).
func (c Clock) String() string {
	cm := c.Mod24()
	return fmt.Sprintf("%02d:%02d:%02d.%09d", clockHours(cm), clockMinutes(cm), clockSeconds(cm), clockNanosec(cm))
}

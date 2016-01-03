// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package clock

import (
	"math"
	"time"
)

// Clock specifies a time of day. It complements the existing time.Duration, applying
// that to the time since midnight (on some arbitrary day in some arbitrary timezone).
// The resolution is to the nearest millisecond, unlike time.Duration (which has nanosecond
// resolution).
//
// It is not intended that Clock be used to represent periods greater than 24 hours nor
// negative values. However, for such lengths of time, a fixed 24 hours per day
// is assumed and a modulo operation Mod24 is provided to discard whole multiples of 24 hours.
//
// See https://en.wikipedia.org/wiki/ISO_8601#Times
type Clock int32

// Common durations.
const (
	// ClockDay is a fixed period of 24 hours. This does not take account of daylight savings, so is not fully general.
	ClockDay Clock = Clock(time.Hour * 24 / time.Millisecond)

	// ClockHour is one hour; it has a similar meaning to time.Hour.
	ClockHour Clock = Clock(time.Hour / time.Millisecond)

	// ClockMinute is one minute; it has a similar meaning to time.Minute.
	ClockMinute Clock = Clock(time.Minute / time.Millisecond)

	// ClockSecond is one second; it has a similar meaning to time.Second.
	ClockSecond Clock = Clock(time.Second / time.Millisecond)
)

// Undefined is provided because the zero value of a Clock *is* defined (i.e. midnight).
// A special value is chosen, which is math.MinInt32.
const Undefined Clock = Clock(math.MinInt32)

// New returns a new Clock with specified hour, minute, second and millisecond.
func New(hour, minute, second, millisec int) Clock {
	hx := Clock(hour) * ClockHour
	mx := Clock(minute) * ClockMinute
	sx := Clock(second) * ClockSecond
	return Clock(hx + mx + sx + Clock(millisec))
}

// NewAt returns a new Clock with specified hour, minute, second and millisecond.
func NewAt(t time.Time) Clock {
	hour, minute, second := t.Clock()
	hx := Clock(hour) * ClockHour
	mx := Clock(minute) * ClockMinute
	sx := Clock(second) * ClockSecond
	ms := Clock(t.Nanosecond() / int(time.Millisecond))
	return Clock(hx + mx + sx + ms)
}

// SinceMidnight returns a new Clock based on a duration since some arbitrary midnight.
func SinceMidnight(d time.Duration) Clock {
	return Clock(d / time.Millisecond)
}

// DurationSinceMidnight convert a clock to a time.Duration since some arbitrary midnight.
func (c Clock) DurationSinceMidnight() time.Duration {
	return time.Duration(c) * time.Millisecond
}

// Add returns a new Clock offset from this clock specified hour, minute, second and millisecond.
// The parameters can be negative.
// If required, use Mod() to correct any overflow or underflow.
func (c Clock) Add(h, m, s, ms int) Clock {
	hx := Clock(h) * ClockHour
	mx := Clock(m) * ClockMinute
	sx := Clock(s) * ClockSecond
	return c + hx + mx + sx + Clock(ms)
}

// IsInOneDay tests whether a clock time is in the range 0 to 24 hours, inclusive. Inside this
// range, a Clock is generally well-behaved. But outside it, there may be errors due to daylight
// savings. Note that 24:00:00 is included as a special case as per ISO-8601 definition of midnight.
func (c Clock) IsInOneDay() bool {
	return 0 <= c && c <= ClockDay
}

// IsMidnight tests whether a clock time is midnight. This is shorthand for c.Mod24() == 0.
// For large values, this assumes that every day has 24 hours.
func (c Clock) IsMidnight() bool {
	return c.Mod24() == 0
}

// Mod24 calculates the remainder vs 24 hours using Euclidean division, in which the result
// will be less than 24 hours and is never negative. Note that this imposes the assumption that
// every day has 24 hours (not correct when daylight saving changes in any timezone).
//
// https://en.wikipedia.org/wiki/Modulo_operation
func (c Clock) Mod24() Clock {
	if 0 <= c && c < ClockDay {
		return c
	}
	if c < 0 {
		q := 1 - c/ClockDay
		m := c + (q * ClockDay)
		if m == ClockDay {
			m = 0
		}
		return m
	}
	q := c / ClockDay
	return c - (q * ClockDay)
}

// Days gets the number of whole days represented by the Clock, assuming that each day is a fixed
// 24 hour period. Negative values are treated so that the range -23h59m59s to -1s is fully
// enclosed in a day numbered -1, and so on. This means that the result is zero only for the
// clock range 0s to 23h59m59s, for which IsInOneDay() returns true.
func (c Clock) Days() int {
	if c < 0 {
		return int(c/ClockDay) - 1
	} else {
		return int(c / ClockDay)
	}
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

// Millisec gets the clock-face number of milliseconds (calculated from the modulo time, see Mod24).
// For example, for 10:20:30.456 this will return 456.
func (c Clock) Millisec() int {
	return int(clockMillisec(c.Mod24()))
}

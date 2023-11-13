// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package clock specifies a time of day with resolution to the nearest millisecond.
package clock

import (
	"math"
	"time"

	"github.com/rickb777/period"
)

// Clock specifies a time of day. It complements the existing time.Duration, applying
// that to the time since midnight (on some arbitrary day in some arbitrary timezone).
// The resolution is to the nearest nanosecond, like time.Duration.
//
// It is not intended that Clock be used to represent periods greater than 24 hours nor
// negative values. However, for such lengths of time, a fixed 24 hours per day
// is assumed and a modulo operation Mod24 is provided to discard whole multiples of 24 hours.
//
// Clock is a type of integer (actually int64), so values can be compared and sorted as
// per other integers. The constants Second, Minute, Hour and Day can be added and subtracted
// with obvious outcomes.
//
// See https://en.wikipedia.org/wiki/ISO_8601#Times
type Clock int64

// Common durations - second, minute, hour and day.
const (
	// Millisecond is one millisecond; it has a similar meaning to time.Millisecond.
	Millisecond Clock = Clock(time.Millisecond)

	// Second is one second; it has a similar meaning to time.Second.
	Second Clock = Clock(time.Second)

	// Minute is one minute; it has a similar meaning to time.Minute.
	Minute Clock = Clock(time.Minute)

	// Hour is one hour; it has a similar meaning to time.Hour.
	Hour Clock = Clock(time.Hour)

	// Day is a fixed period of 24 hours. This does not take account of daylight savings,
	// so is not fully general.
	Day Clock = Clock(time.Hour * 24)
)

// Midnight is the zero value of a Clock.
const Midnight Clock = 0

// Noon is at 12pm.
const Noon Clock = Hour * 12

// Undefined is provided because the zero value of a Clock *is* defined (i.e. Midnight).
// So a special value is chosen, which is math.MinInt64.
const Undefined Clock = Clock(math.MinInt64)

//-------------------------------------------------------------------------------------------------

// New returns a new Clock with specified hour, minute, second and millisecond.
// To set sub-millisecond digits, chain this with AddDuration.
func New(hour, minute, second, millisec int) Clock {
	hx := Clock(hour) * Hour
	mx := Clock(minute) * Minute
	sx := Clock(second) * Second
	ms := Clock(millisec) * Millisecond
	return hx + mx + sx + ms
}

// NewAt returns a new Clock with specified hour, minute, seconds (to nanosecond resolution).
func NewAt(t time.Time) Clock {
	hour, minute, second := t.Clock()
	hx := Clock(hour) * Hour
	mx := Clock(minute) * Minute
	sx := Clock(second) * Second
	ns := Clock(t.Nanosecond())
	return hx + mx + sx + ns
}

// SinceMidnight returns a new Clock based on a duration since some arbitrary midnight.
func SinceMidnight(d time.Duration) Clock {
	return Clock(d)
}

// DurationSinceMidnight convert a clock to a time.Duration since some arbitrary midnight.
func (c Clock) DurationSinceMidnight() time.Duration {
	return time.Duration(c)
}

// Add returns a new Clock offset from this clock specified hour, minute, second and millisecond.
// The parameters can be negative.
//
// If required, use Mod24() to correct any overflow or underflow.
func (c Clock) Add(h, m, s, ms int) Clock {
	hx := Clock(h) * Hour
	mx := Clock(m) * Minute
	sx := Clock(s) * Second
	ns := Clock(ms) * Millisecond
	return c + hx + mx + sx + ns
}

// AddDuration returns a new Clock offset from this clock by a duration.
// The parameter can be negative.
//
// If required, use Mod24() to correct any overflow or underflow.
func (c Clock) AddDuration(d time.Duration) Clock {
	return c + Clock(d)
}

// AddPeriod returns a new Clock offset from this clock by a time period.
// The parameter can be negative.
//
// If required, use Mod24() to correct any overflow or underflow.
//
// The boolean flag is true when the result is precise and false if an
// approximation.
func (c Clock) AddPeriod(p period.Period) (Clock, bool) {
	d, precise := p.Duration()
	return c.AddDuration(d), precise
}

// ModSubtract returns the duration between two clock times.
//
// If c2 is before c (i.e. c2 < c), the result is the duration computed from c - c2.
//
// But if c is before c2, it is assumed that c is after midnight and c2 is before midnight. The
// result is the sum of the evening time from c2 to midnight with the morning time from midnight to c.
// This is the same as Mod24(c - c2).
func (c Clock) ModSubtract(c2 Clock) time.Duration {
	ms := c - c2
	return ms.Mod24().DurationSinceMidnight()
}

//-------------------------------------------------------------------------------------------------

// IsInOneDay tests whether a clock time is in the range 0 to 24 hours, inclusive. Inside this
// range, a Clock is generally well-behaved. But outside it, there may be errors due to daylight
// savings. Note that 24:00:00 is included as a special case as per ISO-8601 definition of midnight.
func (c Clock) IsInOneDay() bool {
	return Midnight <= c && c <= Day
}

// IsMidnight tests whether a clock time is midnight. This is shorthand for c.Mod24() == 0.
// For large values, this assumes that every day has 24 hours.
func (c Clock) IsMidnight() bool {
	return c.Mod24() == Midnight
}

// TruncateMillisecond discards any fractional digits within the millisecond represented by c.
// For example, for 10:20:30.456111222 this will return 10:20:30.456.
// This method will force the String method to limit its output to three decimal places.
func (c Clock) TruncateMillisecond() Clock {
	return (c / Millisecond) * Millisecond
}

// Mod24 calculates the remainder vs 24 hours using Euclidean division, in which the result
// will be less than 24 hours and is never negative. Note that this imposes the assumption that
// every day has 24 hours (not correct when daylight saving changes in any timezone).
//
// https://en.wikipedia.org/wiki/Modulo_operation
func (c Clock) Mod24() Clock {
	if Midnight <= c && c < Day {
		return c
	}
	if c < Midnight {
		q := 1 - c/Day
		m := c + (q * Day)
		if m == Day {
			m = Midnight
		}
		return m
	}
	q := c / Day
	return c - (q * Day)
}

//-------------------------------------------------------------------------------------------------

// Days gets the number of whole days represented by the Clock, assuming that each day is a fixed
// 24 hour period. Negative values are treated so that the range -23h59m59s to -1s is fully
// enclosed in a day numbered -1, and so on. This means that the result is zero only for the
// clock range 0s to 23h59m59s, for which IsInOneDay() returns true.
func (c Clock) Days() int {
	if c < Midnight {
		return int(c/Day) - 1
	}
	return int(c / Day)
}

// Hour gets the clock-face number of hours (calculated from the modulo time, see Mod24).
func (c Clock) Hour() int {
	return int(clockHour(c.Mod24()))
}

// Minute gets the clock-face number of minutes (calculated from the modulo time, see Mod24).
// For example, for 22:35 this will return 35.
func (c Clock) Minute() int {
	return int(clockMinute(c.Mod24()))
}

// Second gets the clock-face number of seconds (calculated from the modulo time, see Mod24).
// For example, for 10:20:30 this will return 30.
func (c Clock) Second() int {
	return int(clockSecond(c.Mod24()))
}

// Millisecond gets the clock-face number of milliseconds within the second specified by c
// (calculated from the modulo time, see Mod24), in the range [0, 999].
// For example, for 10:20:30.456 this will return 456.
func (c Clock) Millisecond() int {
	return int(clockNanosecond(c.Mod24()) / 1_000_000)
}

// Nanosecond gets the clock-face number of nanoseconds within the second specified by c
// (calculated from the modulo time, see Mod24), in the range [0, 999999999].
// For example, for 10:20:30.456111222 this will return 456111222.
func (c Clock) Nanosecond() int {
	return int(clockNanosecond(c.Mod24()))
}

// HourMinuteSecond gets the hours, minutes and seconds values (calculated from the modulo time, see Mod24).
func (c Clock) HourMinuteSecond() (int, int, int) {
	c2 := c.Mod24()
	return int(clockHour(c2)), int(clockMinute(c2)), int(clockSecond(c2))
}

//-------------------------------------------------------------------------------------------------

func clockHour(cm Clock) Clock {
	return cm / Hour
}

func clockHour12(cm Clock) (Clock, string) {
	h := clockHour(cm)
	if h < 1 {
		return 12, "am"
	} else if h > 12 {
		return h - 12, "pm"
	} else if h == 12 {
		return 12, "pm"
	}
	return h, "am"
}

func clockMinute(cm Clock) Clock {
	return (cm % Hour) / Minute
}

func clockSecond(cm Clock) Clock {
	return (cm % Minute) / Second
}

func clockMillisecond(cm Clock) Clock {
	return (cm % Second) / Millisecond
}

func clockNanosecond(cm Clock) Clock {
	return cm % Second
}

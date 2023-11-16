// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package date

import (
	"math"
	"time"

	"github.com/rickb777/date/v2/clock"
	"github.com/rickb777/date/v2/gregorian"
	"github.com/rickb777/period"
)

// A Date represents a date under the proleptic Gregorian calendar as
// used by ISO 8601. This calendar uses astronomical year numbering,
// so it includes a year 0 and represents earlier years as negative numbers
// (i.e. year 0 is 1 BC; year -1 is 2 BC, and so on).
//
// On 32-bit architectures, dates range from about year -5 million to +5 million.
// On 64-bit architecturew, the range is huge.
//
// Programs using dates should typically store and pass them as values,
// not pointers.  That is, date variables and struct fields should be of
// type date.Date, not *date.Date unless the pointer indicates an optional
// value.  A Date value can be used by multiple goroutines simultaneously.
//
// Date values can be compared using the ==, !=, >, >=, <, and <= operators.
//
// Because a Date is a number of days since Zero, + and - operations
// add or subtract some number of days.
//
// The first official date of the Gregorian calendar was Friday, October 15th
// 1582, quite unrelated to the Unix epoch or the year 0 used here. The Date
// type does not distinguish between official Gregorian dates and earlier
// proleptic dates, which can also be represented when needed.
type Date int

const (
	// Zero is the named zero value for Date and corresponds to Saturday, January 1,
	// year 0 in the proleptic Gregorian calendar using astronomical year numbering.
	Zero Date = 0
)

// New returns the Date value corresponding to the given year, month, and day.
//
// The month and day may be outside their usual ranges and will be normalized
// during the conversion.
func New(year int, month time.Month, day int) Date {
	t := time.Date(year, month, day, 12, 0, 0, 0, time.UTC)
	return encode(t)
}

// NewAt returns the Date value corresponding to the given time.
// Note that the date is computed relative to the time zone specified by
// the given Time value.
func NewAt(t time.Time) Date {
	return encode(t)
}

// Today returns today's date according to the current local time.
func Today() Date {
	return encode(time.Now())
}

// TodayUTC returns today's date according to the current UTC time.
func TodayUTC() Date {
	return encode(time.Now().UTC())
}

// TodayIn returns today's date according to the current time relative to
// the specified location.
func TodayIn(loc *time.Location) Date {
	t := time.Now().In(loc)
	return encode(t)
}

// Min returns the smallest representable date, which is nearly 6 million years in the past.
func Min() Date {
	return Date(math.MinInt32 + 1)
}

// Max returns the largest representable date, which is nearly 6 million years in the future.
func Max() Date {
	return Date(math.MaxInt32 - zeroOffset)
}

// MidnightUTC returns a Time value corresponding to midnight on the given date d,
// UTC time.  Note that midnight is the beginning of the day rather than the end.
func (d Date) MidnightUTC() time.Time {
	return decode(d)
}

// Midnight returns a Time value corresponding to midnight on the given date d,
// local time.  Note that midnight is the beginning of the day rather than the end.
func (d Date) Midnight() time.Time {
	return d.MidnightIn(time.Local)
}

// MidnightIn returns a Time value corresponding to midnight on the given date d,
// relative to the specified time zone.  Note that midnight is the beginning
// of the day rather than the end.
func (d Date) MidnightIn(loc *time.Location) time.Time {
	return d.Time(0, loc)
}

// Time returns a Time value corresponding to a clock time on the given date d,
// relative to the specified time zone. A common use-case is to obtain the midnight
// time, for which the clock value is simply zero.
func (d Date) Time(clock clock.Clock, loc *time.Location) time.Time {
	t := decode(d).In(loc)
	_, offset := t.Zone()
	return t.Add(time.Duration(-offset) * time.Second).Add(time.Duration(clock))
}

// Date returns the year, month, and day of d.
// The first day of the month is 1.
func (d Date) Date() (year int, month time.Month, day int) {
	return decode(d).Date()
}

// LastDayOfMonth returns the last day of the month specified by d.
// The first day of the month is 1.
func (d Date) LastDayOfMonth() int {
	y, m, _ := d.Date()
	return gregorian.DaysIn(y, m)
}

// Day returns the day of the month specified by d.
// The first day of the month is 1.
func (d Date) Day() int {
	return decode(d).Day()
}

// Month returns the month of the year specified by d.
func (d Date) Month() time.Month {
	return decode(d).Month()
}

// Year returns the year specified by d.
func (d Date) Year() int {
	return decode(d).Year()
}

// YearDay returns the day of the year specified by d, in the range [1,365] for
// non-leap years, and [1,366] in leap years.
func (d Date) YearDay() int {
	return decode(d).YearDay()
}

// Weekday returns the day of the week specified by d.
func (d Date) Weekday() time.Weekday {
	// Date zero, January 1, 0000, fell on a Saturday
	const weekdayZero = time.Saturday
	// Taking into account potential for overflow and negative offset
	return time.Weekday((int(weekdayZero) + int(d)%7 + 7) % 7)
}

// ISOWeek returns the ISO 8601 year and week number in which d occurs.
// Week ranges from 1 to 53. Jan 01 to Jan 03 of year n might belong to
// week 52 or 53 of year n-1, and Dec 29 to Dec 31 might belong to week 1
// of year n+1.
func (d Date) ISOWeek() (year, week int) {
	return decode(d).ISOWeek()
}

// AddDate returns the date corresponding to adding the given number of years,
// months, and days to d. For example, AddData(-1, 2, 3) applied to
// January 1, 2011 returns March 4, 2010.
//
// AddDate normalizes its result in the same way that Date does,
// so, for example, adding one month to October 31 yields
// December 1, the normalized form for November 31.
//
// The addition of all fields is performed before normalisation of any; this can affect
// the result. For example, adding 0y 1m 3d to September 28 gives October 31 (not
// November 1).
func (d Date) AddDate(years, months, days int) Date {
	t := decode(d).AddDate(years, months, days)
	return encode(t)
}

// AddPeriod returns the date corresponding to adding the given period. If the
// period's fields are be negative, this results in an earlier date.
//
// Any time component only affects the result for periods containing
// more that 24 hours in the hours/minutes/seconds fields
func (d Date) AddPeriod(delta period.Period) Date {
	t1 := decode(d)
	t2, _ := delta.AddTo(t1)
	return encode(t2)
}

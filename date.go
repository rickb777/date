// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package date provides functionality for working with dates.
//
// This package introduces a light-weight Date type that is storage-efficient
// and covenient for calendrical calculations and date parsing and formatting
// (including years outside the [0,9999] interval).
//
// Credits
//
// This package follows very closely the design of package time
// (http://golang.org/pkg/time/) in the standard library, many of the Date
// methods are implemented using the corresponding methods of the time.Time
// type, and much of the documentation is copied directly from that package.
//
// References
//
// https://golang.org/src/time/time.go
//
// https://en.wikipedia.org/wiki/Gregorian_calendar
//
// https://en.wikipedia.org/wiki/Proleptic_Gregorian_calendar
//
// https://en.wikipedia.org/wiki/Astronomical_year_numbering
//
// https://en.wikipedia.org/wiki/ISO_8601
//
// https://tools.ietf.org/html/rfc822
//
// https://tools.ietf.org/html/rfc850
//
// https://tools.ietf.org/html/rfc1123
//
// https://tools.ietf.org/html/rfc3339
//
package date

import (
	"errors"
	"fmt"
	"math"
	"time"
)

// A Date represents a date under the (proleptic) Gregorian calendar as
// used by ISO 8601. This calendar uses astronomical year numbering,
// so it includes a year 0 and represents earlier years as negative numbers
// (i.e. year 0 is 1 BC; year -1 is 2 BC, and so on).
//
// A Date value requires 4 bytes of storage and can represent dates from
// Tue, 23 Jun -5,877,641 (5,877,642 BC) to Fri, 11 Jul 5,881,580.
// Dates outside that range will "wrap around".
//
// Programs using dates should typically store and pass them as values,
// not pointers.  That is, date variables and struct fields should be of
// type date.Date, not *date.Date.  A Date value can be used by
// multiple goroutines simultaneously.
//
// Date values can be compared using the Before, After, and Equal methods
// as well as the == and != operators.
// The Sub method subtracts two dates, returning the number of days between
// them.
// The Add method adds a Date and a number of days, producing a Date.
//
// The zero value of type Date is Thursday, January 1, 1970.
// As this date is unlikely to come up in practice, the IsZero method gives
// a simple way of detecting a date that has not been initialized explicitly.
//
type Date struct {
	// day gives the number of days elapsed since date zero.
	day int32
}

// New returns the Date value corresponding to the given year, month, and day.
//
// The month and day may be outside their usual ranges and will be normalized
// during the conversion.
func New(year int, month time.Month, day int) Date {
	t := time.Date(year, month, day, 12, 0, 0, 0, time.UTC)
	return Date{encode(t)}
}

// NewAt returns the Date value corresponding to the given time.
// Note that the date is computed relative to the time zone specified by
// the given Time value.
func NewAt(t time.Time) Date {
	return Date{encode(t)}
}

// Today returns today's date according to the current local time.
func Today() Date {
	t := time.Now()
	return Date{encode(t)}
}

// TodayUTC returns today's date according to the current UTC time.
func TodayUTC() Date {
	t := time.Now().UTC()
	return Date{encode(t)}
}

// TodayIn returns today's date according to the current time relative to
// the specified location.
func TodayIn(loc *time.Location) Date {
	t := time.Now().In(loc)
	return Date{encode(t)}
}

// Min returns the smallest representable date.
func Min() Date {
	return Date{math.MinInt32}
}

// Max returns the largest representable date.
func Max() Date {
	return Date{math.MaxInt32}
}

// UTC returns a Time value corresponding to midnight on the given date,
// UTC time.  Note that midnight is the beginning of the day rather than the end.
func (d Date) UTC() time.Time {
	return decode(d.day)
}

// Local returns a Time value corresponding to midnight on the given date,
// local time.  Note that midnight is the beginning of the day rather than the end.
func (d Date) Local() time.Time {
	return d.In(time.Local)
}

// In returns a Time value corresponding to midnight on the given date,
// relative to the specified time zone.  Note that midnight is the beginning
// of the day rather than the end.
func (d Date) In(loc *time.Location) time.Time {
	t := decode(d.day).In(loc)
	_, offset := t.Zone()
	return t.Add(time.Duration(-offset) * time.Second)
}

// Date returns the year, month, and day of d.
func (d Date) Date() (year int, month time.Month, day int) {
	t := decode(d.day)
	return t.Date()
}

// Day returns the day of the month specified by d.
// The first day of the month is 1.
func (d Date) Day() int {
	t := decode(d.day)
	return t.Day()
}

// Month returns the month of the year specified by d.
func (d Date) Month() time.Month {
	t := decode(d.day)
	return t.Month()
}

// Year returns the year specified by d.
func (d Date) Year() int {
	t := decode(d.day)
	return t.Year()
}

// YearDay returns the day of the year specified by d, in the range [1,365] for
// non-leap years, and [1,366] in leap years.
func (d Date) YearDay() int {
	t := decode(d.day)
	return t.YearDay()
}

// Weekday returns the day of the week specified by d.
func (d Date) Weekday() time.Weekday {
	// Date zero, January 1, 1970, fell on a Thursday
	wdayZero := time.Thursday
	// Taking into account potential for overflow and negative offset
	return time.Weekday((int32(wdayZero) + d.day%7 + 7) % 7)
}

// ISOWeek returns the ISO 8601 year and week number in which d occurs.
// Week ranges from 1 to 53. Jan 01 to Jan 03 of year n might belong to
// week 52 or 53 of year n-1, and Dec 29 to Dec 31 might belong to week 1
// of year n+1.
func (d Date) ISOWeek() (year, week int) {
	t := decode(d.day)
	return t.ISOWeek()
}

// IsZero reports whether t represents the zero date.
func (d Date) IsZero() bool {
	return d.day == 0
}

// Equal reports whether d and u represent the same date.
func (d Date) Equal(u Date) bool {
	return d.day == u.day
}

// Before reports whether the date d is before u.
func (d Date) Before(u Date) bool {
	return d.day < u.day
}

// After reports whether the date d is after u.
func (d Date) After(u Date) bool {
	return d.day > u.day
}

// Add returns the date d plus the given number of days.
func (d Date) Add(days int) Date {
	return Date{d.day + int32(days)}
}

// AddDate returns the date corresponding to adding the given number of years,
// months, and days to d. For example, AddData(-1, 2, 3) applied to
// January 1, 2011 returns March 4, 2010.
func (d Date) AddDate(years, months, days int) Date {
	t := decode(d.day)
	t = t.AddDate(years, months, days)
	return Date{encode(t)}
}

// Sub returns d-u as the number of days between the two dates.
func (d Date) Sub(u Date) (days int) {
	return int(d.day - u.day)
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (d Date) MarshalBinary() ([]byte, error) {
	enc := []byte{
		byte(d.day >> 24),
		byte(d.day >> 16),
		byte(d.day >> 8),
		byte(d.day),
	}
	return enc, nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (d *Date) UnmarshalBinary(data []byte) error {
	if len(data) == 0 {
		return errors.New("Date.UnmarshalBinary: no data")
	}
	if len(data) != 4 {
		return errors.New("Date.UnmarshalBinary: invalid length")
	}

	d.day = int32(data[3]) | int32(data[2])<<8 | int32(data[1])<<16 | int32(data[0])<<24

	return nil
}

// GobEncode implements the gob.GobEncoder interface.
func (d Date) GobEncode() ([]byte, error) {
	return d.MarshalBinary()
}

// GobDecode implements the gob.GobDecoder interface.
func (d *Date) GobDecode(data []byte) error {
	return d.UnmarshalBinary(data)
}

// MarshalJSON implements the json.Marshaler interface.
// The date is a quoted string in ISO 8601 extended format (e.g. "2006-01-02").
// If the year of the date falls outside the [0,9999] range, this format
// produces an expanded year representation with possibly extra year digits
// beyond the prescribed four-digit minimum and with a + or - sign prefix
// (e.g. , "+12345-06-07", "-0987-06-05").
func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.String() + `"`), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// The date is expected to be a quoted string in ISO 8601 extended format
// (e.g. "2006-01-02", "+12345-06-07", "-0987-06-05");
// the year must use at least 4 digits and if outside the [0,9999] range
// must be prefixed with a + or - sign.
func (d *Date) UnmarshalJSON(data []byte) (err error) {
	value := string(data)
	n := len(value)
	if n < 2 || value[0] != '"' || value[n-1] != '"' {
		return fmt.Errorf("Date.UnmarshalJSON: missing double quotes (%s)", value)
	}
	u, err := ParseISO(value[1 : n-1])
	if err != nil {
		return err
	}
	d.day = u.day
	return nil
}

// MarshalText implements the encoding.TextMarshaler interface.
// The date is given in ISO 8601 extended format (e.g. "2006-01-02").
// If the year of the date falls outside the [0,9999] range, this format
// produces an expanded year representation with possibly extra year digits
// beyond the prescribed four-digit minimum and with a + or - sign prefix
// (e.g. , "+12345-06-07", "-0987-06-05").
func (d Date) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// The date is expected to be in ISO 8601 extended format
// (e.g. "2006-01-02", "+12345-06-07", "-0987-06-05");
// the year must use at least 4 digits and if outside the [0,9999] range
// must be prefixed with a + or - sign.
func (d *Date) UnmarshalText(data []byte) error {
	u, err := ParseISO(string(data))
	if err != nil {
		return err
	}
	d.day = u.day
	return nil
}

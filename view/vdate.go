// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package view provides a simple API for formatting dates as strings in a manner that is easy to use in view-models,
// especially when using Go templates.
package view

import (
	"github.com/rickb777/date/v2"
)

const (
	// DMYFormat is a typical British representation.
	DMYFormat = "02/01/2006"
	// MDYFormat is a typical American representation.
	MDYFormat = "01/02/2006"
	// ISOFormat is ISO-8601 YYYY-MM-DD.
	ISOFormat = "2006-02-01"
	// DefaultFormat is used by Format() unless a different format is set.
	DefaultFormat = DMYFormat
)

// A VDate holds a Date and provides easy ways to render it, e.g. in Go templates.
type VDate struct {
	d      date.Date
	format string
}

// NewVDate wraps a Date.
func NewVDate(d date.Date) VDate {
	return VDate{d, DefaultFormat}
}

// Date returns the underlying date.
func (v VDate) Date() date.Date {
	return v.d
}

// IsYesterday returns true if the date is yesterday's date.
func (v VDate) IsYesterday() bool {
	return v.d+1 == date.Today()
}

// IsToday returns true if the date is today's date.
func (v VDate) IsToday() bool {
	return v.d == date.Today()
}

// IsTomorrow returns true if the date is tomorrow's date.
func (v VDate) IsTomorrow() bool {
	return v.d-1 == date.Today()
}

// IsOdd returns true if the date is an odd number. This is useful for
// zebra striping etc.
func (v VDate) IsOdd() bool {
	return v.d%2 == 0
}

// String formats the date in basic ISO8601 format YYYY-MM-DD.
func (v VDate) String() string {
	if v.d == date.Zero {
		return ""
	}
	return v.d.String()
}

// WithFormat creates a new instance containing the specified format string.
func (v VDate) WithFormat(f string) VDate {
	return VDate{v.d, f}
}

// Format formats the date using the specified format string, or "02/01/2006" by default.
// Use WithFormat to set this up.
func (v VDate) Format() string {
	return v.d.Format(v.format)
}

// Mon returns the day name as three letters.
func (v VDate) Mon() string {
	return v.d.Format("Mon")
}

// Monday returns the full day name.
func (v VDate) Monday() string {
	return v.d.Format("Monday")
}

// Day2 returns the day number without a leading zero.
func (v VDate) Day2() string {
	return v.d.Format("2")
}

// Day02 returns the day number with a leading zero if necessary.
func (v VDate) Day02() string {
	return v.d.Format("02")
}

// Day2nd returns the day number without a leading zero but with the appropriate
// "st", "nd", "rd", "th" suffix.
func (v VDate) Day2nd() string {
	return v.d.Format("2nd")
}

// Month1 returns the month number without a leading zero.
func (v VDate) Month1() string {
	return v.d.Format("1")
}

// Month01 returns the month number with a leading zero if necessary.
func (v VDate) Month01() string {
	return v.d.Format("01")
}

// Jan returns the month name abbreviated to three letters.
func (v VDate) Jan() string {
	return v.d.Format("Jan")
}

// January returns the full month name.
func (v VDate) January() string {
	return v.d.Format("January")
}

// Year returns the four-digit year.
func (v VDate) Year() string {
	return v.d.Format("2006")
}

// Next returns a fluent generator for later dates.
func (v VDate) Next() VDateDelta {
	return VDateDelta{d: v.d, format: v.format, sign: 1}
}

// Previous returns a fluent generator for earlier dates.
func (v VDate) Previous() VDateDelta {
	return VDateDelta{d: v.d, format: v.format, sign: -1}
}

//-------------------------------------------------------------------------------------------------
// Only lossy transcoding is supported here because the intention is that data exchange should be
// via the main Date type; VDate is only intended for output through view layers.

// MarshalText implements the encoding.TextMarshaler interface.
func (v VDate) MarshalText() ([]byte, error) {
	return v.d.MarshalText()
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// Note that the format value gets lost.
func (v *VDate) UnmarshalText(data []byte) (err error) {
	var u date.Date
	err = u.UnmarshalText(data)
	if err == nil {
		v.d = u
		v.format = DefaultFormat
	}
	return err
}

//-------------------------------------------------------------------------------------------------

// VDateDelta is a VDate with the ability to add or subtract days, weeks, months or years.
type VDateDelta struct {
	d      date.Date
	format string
	sign   int
}

// Day adds or subtracts one day.
func (dd VDateDelta) Day() VDate {
	return VDate{d: dd.d + date.Date(dd.sign), format: dd.format}
}

// Week adds or subtracts one week.
func (dd VDateDelta) Week() VDate {
	return VDate{d: dd.d + date.Date(dd.sign*7), format: dd.format}
}

// Month adds or subtracts one month.
func (dd VDateDelta) Month() VDate {
	return VDate{d: dd.d.AddDate(0, dd.sign, 0), format: dd.format}
}

// Year adds or subtracts one year.
func (dd VDateDelta) Year() VDate {
	return VDate{d: dd.d.AddDate(dd.sign, 0, 0), format: dd.format}
}

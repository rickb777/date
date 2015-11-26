// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package daterange

import (
	. "github.com/rickb777/date"
	"time"
//	"fmt"
	"fmt"
)

// DateRange carries a pair of dates encompassing a range. In operations on the range,
// the start and end are both considered to be inclusive.
type DateRange struct {
	Start Date
	End   Date
}

// NewDateRangeOf assembles a new date range from a start time and a duration, discarding
// the precise time-of-day information. The start time includes a location, which is not
// necessarily UTC. The duration can be negative; the result is
// normalised so that the end date is not before the start date.
func NewDateRangeOf(start time.Time, duration time.Duration) DateRange {
	sd := NewAt(start)
	ed := NewAt(start.Add(duration))
	return DateRange{sd, ed}.Normalise()
}

// NewDateRange assembles a new date range from two dates, normalising them so that the
// end date is not before the start date.
func NewDateRange(start, end Date) DateRange {
	return DateRange{start, end}.Normalise()
}

// NewYearOf constructs the range encompassing the whole year specified.
func NewYearOf(year int) DateRange {
	start := New(year, time.January, 1)
	end := New(year, time.December, 31)
	return DateRange{start, end}
}

// NewMonthOf constructs the range encompassing the whole month specified for a given year.
// It handles leap years correctly.
func NewMonthOf(year int, month time.Month) DateRange {
	start := New(year, month, 1)
	endT := time.Date(year, month + 1, 1, 0, 0, 0, 0, time.UTC)
	end := NewAt(endT.Add(-1))
	return DateRange{start, end}
}

// OneDayRange constructs a range of exactly one day. This is often a useful basis for
// further operations.
func OneDayRange(day Date) DateRange {
	return NewDateRange(day, day)
}

// Normalise ensures that the start date is before (or equal to) the end date.
// They are swapped if necessary. The normalised date range is returned.
func (dateRange DateRange) Normalise() DateRange {
	if dateRange.End.Before(dateRange.Start) {
		return DateRange{dateRange.End, dateRange.Start}
	}
	return dateRange
}

// ShiftBy moves the date range by moving both the start and end dates similarly.
// A negative parameter is allowed.
func (dateRange DateRange) ShiftBy(days int) DateRange {
	if days == 0 {
		return dateRange
	}
	newStart := dateRange.Start.Add(days)
	newEnd := dateRange.End.Add(days)
	return DateRange{newStart, newEnd}
}

// ExtendBy extends (or reduces) the date range by moving the end date.
// A negative parameter is allowed and the result is normalised.
func (dateRange DateRange) ExtendBy(days int) DateRange {
	if days == 0 {
		return dateRange
	}
	// this relies on normalisation provided by the function
	newEnd := dateRange.End.Add(days)
	return DateRange{dateRange.Start, newEnd}.Normalise()
}

func (dateRange DateRange) String() string {
	return fmt.Sprintf("%s to %s", dateRange.Start, dateRange.End)
}

// Contains tests whether the date range contains a specified date. The start and end of
// the range are both treated inclusively.
//
// If a calculation needs to be 'half-open' (i.e. the end is exclusive), simply use the
// expression 'dateRange.ExtendBy(-1).Contains(d)'
func (dateRange DateRange) Contains(d Date) bool {
	return !(d.Before(dateRange.Start) || d.After(dateRange.End))
}

// StartUTC assumes that the start date is a UTC date and gets the start time of that date, as UTC.
// It returns midnight on the first day of the range.
func (dateRange DateRange) StartUTC() time.Time {
	return dateRange.Start.UTC()
}

// EndUTC assumes that the end date is a UTC date and gets the end time of that date, as UTC.
// It returns the very last nanosecond before midnight the following day.
func (dateRange DateRange) EndUTC() time.Time {
	return dateRange.End.Add(1).UTC().Add(minusOneNano)
}

const minusOneNano time.Duration = -1

// ContainsTime tests whether a given local time is within the date range. The time range is
// from midnight on the start day to one nanosecond before midnight on the day after the end date.
//
// If a calculation needs to be 'half-open' (i.e. the end is exclusive), simply use the
// expression 'dateRange.ExtendBy(-1).ContainsTime(t)'
func (dateRange DateRange) ContainsTime(t time.Time) bool {
	utc := t.In(time.UTC)
	return !(utc.Before(dateRange.StartUTC()) || dateRange.EndUTC().Before(utc))
}

// Merge conjoins two date ranges. As a special case, if one range is entirely contained within
// the other range, the larger of the two is returned. Otherwise, the result is the start of the
// earlier one to the end of the later one, even if the two ranges don't overlap.
func (dateRange DateRange) Merge(other DateRange) DateRange {
	if dateRange.Start.After(other.Start) {
		// swap the ranges to simplify the logic
		return other.Merge(dateRange)

	} else if dateRange.End.After(other.End) {
		// other is a proper subrange of dateRange
		return dateRange

	} else {
		return DateRange{dateRange.Start, other.End}
	}
}

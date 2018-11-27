// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package timespan

import (
	"fmt"
	"time"

	"github.com/rickb777/date"
	"github.com/rickb777/date/period"
)

const minusOneNano time.Duration = -1

// DateRange carries a date and a number of days and describes a range between two dates.
type DateRange struct {
	mark date.Date
	days date.PeriodOfDays
}

// NewDateRangeOf assembles a new date range from a start time and a duration, discarding
// the precise time-of-day information. The start time includes a location, which is not
// necessarily UTC. The duration can be negative.
func NewDateRangeOf(start time.Time, duration time.Duration) DateRange {
	sd := date.NewAt(start)
	ed := date.NewAt(start.Add(duration))
	return DateRange{sd, date.PeriodOfDays(ed.Sub(sd))}
}

// NewDateRange assembles a new date range from two dates. These are half-open, so
// if start and end are the same, the range spans zero (not one) day. Similarly, if they
// are on subsequent days, the range is one date (not two).
// The result is normalised.
func NewDateRange(start, end date.Date) DateRange {
	if end.Before(start) {
		return DateRange{end, date.PeriodOfDays(start.Sub(end))}
	}
	return DateRange{start, date.PeriodOfDays(end.Sub(start))}
}

// NewYearOf constructs the range encompassing the whole year specified.
func NewYearOf(year int) DateRange {
	start := date.New(year, time.January, 1)
	end := date.New(year+1, time.January, 1)
	return DateRange{start, date.PeriodOfDays(end.Sub(start))}
}

// NewMonthOf constructs the range encompassing the whole month specified for a given year.
// It handles leap years correctly.
func NewMonthOf(year int, month time.Month) DateRange {
	start := date.New(year, month, 1)
	endT := time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)
	end := date.NewAt(endT)
	return DateRange{start, date.PeriodOfDays(end.Sub(start))}
}

// EmptyRange constructs an empty range. This is often a useful basis for
// further operations but note that the end date is undefined.
func EmptyRange(day date.Date) DateRange {
	return DateRange{day, 0}
}

// OneDayRange constructs a range of exactly one day. This is often a useful basis for
// further operations. Note that the last date is the same as the start date.
func OneDayRange(day date.Date) DateRange {
	return DateRange{day, 1}
}

// DayRange constructs a range of n days.
//
// Note that n can be negative. In this case, the specified day will be the end day,
// which is outside of the half-open range; the last day will be the day before the
// day specified.
func DayRange(day date.Date, n date.PeriodOfDays) DateRange {
	if n < 0 {
		return DateRange{day.Add(n), -n}
	}
	return DateRange{day, n}
}

// Days returns the period represented by this range. This will never be negative.
func (dateRange DateRange) Days() date.PeriodOfDays {
	if dateRange.days < 0 {
		return -dateRange.days
	}
	return dateRange.days
}

// IsZero returns true if this has a zero start date and the the range is empty.
// Usually this is because the range was created via the zero value.
func (dateRange DateRange) IsZero() bool {
	return dateRange.days == 0 && dateRange.mark.IsZero()
}

// IsEmpty returns true if this has a starting date but the range is empty (zero days).
func (dateRange DateRange) IsEmpty() bool {
	return dateRange.days == 0
}

// Start returns the earliest date represented by this range.
func (dateRange DateRange) Start() date.Date {
	if dateRange.days < 0 {
		return dateRange.mark.Add(date.PeriodOfDays(1 + dateRange.days))
	}
	return dateRange.mark
}

// Last returns the last date (inclusive) represented by this range. Be careful because
// if the range is empty (i.e. has zero days), then the last is undefined so an empty date
// is returned. Therefore it is often more useful to use End() instead of Last().
// See also IsEmpty().
func (dateRange DateRange) Last() date.Date {
	if dateRange.days < 0 {
		return dateRange.mark // because mark is at the end
	} else if dateRange.days == 0 {
		return date.Date{}
	}
	return dateRange.mark.Add(dateRange.days - 1)
}

// End returns the date following the last date of the range. End can be considered to
// be the exclusive end, i.e. the final value of a half-open range.
//
// If the range is empty (i.e. has zero days), then the start date is returned, this being
// also the (half-open) end value in that case. This is more useful than the undefined result
// returned by Last() for empty ranges.
func (dateRange DateRange) End() date.Date {
	if dateRange.days < 0 {
		return dateRange.mark.Add(1) // because mark is at the end
	}
	return dateRange.mark.Add(dateRange.days)
}

// Normalise ensures that the number of days is zero or positive.
// The normalised date range is returned;
// in this value, the mark date is the same as the start date.
func (dateRange DateRange) Normalise() DateRange {
	if dateRange.days < 0 {
		return DateRange{dateRange.mark.Add(dateRange.days), -dateRange.days}
	}
	return dateRange
}

// ShiftBy moves the date range by moving both the start and end dates similarly.
// A negative parameter is allowed.
func (dateRange DateRange) ShiftBy(days date.PeriodOfDays) DateRange {
	if days == 0 {
		return dateRange
	}
	newMark := dateRange.mark.Add(days)
	return DateRange{newMark, dateRange.days}
}

// ExtendBy extends (or reduces) the date range by moving the end date.
// A negative parameter is allowed and this may cause the range to become inverted
// (i.e. the mark date becomes the end date instead of the start date).
func (dateRange DateRange) ExtendBy(days date.PeriodOfDays) DateRange {
	if days == 0 {
		return dateRange
	}
	return DateRange{dateRange.mark, dateRange.days + days}.Normalise()
}

// ShiftByPeriod moves the date range by moving both the start and end dates similarly.
// A negative parameter is allowed.
//
// Any time component is ignored. Therefore, be careful with periods containing
// more that 24 hours in the hours/minutes/seconds fields. These will not be
// normalised for you; if you want this behaviour, call delta.Normalise(false)
// on the input parameter.
//
// For example, PT24H adds nothing, whereas P1D adds one day as expected. To
// convert a period such as PT24H to its equivalent P1D, use
// delta.Normalise(false) as the input.
func (dateRange DateRange) ShiftByPeriod(delta period.Period) DateRange {
	if delta.IsZero() {
		return dateRange
	}
	newMark := dateRange.mark.AddPeriod(delta)
	//fmt.Printf("mark + %v : %v -> %v", delta, dateRange.mark, newMark)
	return DateRange{newMark, dateRange.days}
}

// ExtendByPeriod extends (or reduces) the date range by moving the end date.
// A negative parameter is allowed and this may cause the range to become inverted
// (i.e. the mark date becomes the end date instead of the start date).
func (dateRange DateRange) ExtendByPeriod(delta period.Period) DateRange {
	if delta.IsZero() {
		return dateRange
	}
	newEnd := dateRange.End().AddPeriod(delta)
	//fmt.Printf("%v, end + %v : %v -> %v", dateRange.mark, delta, dateRange.End(), newEnd)
	return NewDateRange(dateRange.Start(), newEnd)
}

// String describes the date range in human-readable form.
func (dateRange DateRange) String() string {
	norm := dateRange.Normalise()
	switch norm.days {
	case 0:
		return fmt.Sprintf("0 days at %s", norm.mark)
	case 1:
		return fmt.Sprintf("1 day on %s", norm.mark)
	default:
		return fmt.Sprintf("%d days from %s to %s", norm.days, norm.Start(), norm.Last())
	}
}

// Contains tests whether the date range contains a specified date.
// Empty date ranges (i.e. zero days) never contain anything.
func (dateRange DateRange) Contains(d date.Date) bool {
	if dateRange.days == 0 {
		return false
	}
	return !(d.Before(dateRange.Start()) || d.After(dateRange.Last()))
}

// StartUTC assumes that the start date is a UTC date and gets the start time of that date, as UTC.
// It returns midnight on the first day of the range.
func (dateRange DateRange) StartUTC() time.Time {
	return dateRange.Start().UTC()
}

// EndUTC assumes that the end date is a UTC date and returns the time a nanosecond after the end time
// in a specified location. Along with StartUTC, this gives a 'half-open' range where the start
// is inclusive and the end is exclusive.
func (dateRange DateRange) EndUTC() time.Time {
	return dateRange.End().UTC()
}

// ContainsTime tests whether a given local time is within the date range. The time range is
// from midnight on the start day to one nanosecond before midnight on the day after the end date.
// Empty date ranges (i.e. zero days) never contain anything.
//
// If a calculation needs to be 'half-open' (i.e. the end date is exclusive), simply use the
// expression 'dateRange.ExtendBy(-1).ContainsTime(t)'
func (dateRange DateRange) ContainsTime(t time.Time) bool {
	if dateRange.days == 0 {
		return false
	}
	utc := t.In(time.UTC)
	return !(utc.Before(dateRange.StartUTC()) || dateRange.EndUTC().Add(minusOneNano).Before(utc))
}

// Merge combines two date ranges by calculating a date range that just encompasses them both.
// There are two special cases.
//
// Firstly, if one range is entirely contained within the other range, the larger of the two is
// returned. Otherwise, the result is from the start of the earlier one to the end of the later
// one, even if the two ranges don't overlap.
//
// Secondly, if either range is the zero value (see IsZero), it is excluded from the merge and
// the other range is returned unchanged.
func (dateRange DateRange) Merge(otherRange DateRange) DateRange {
	if otherRange.IsZero() {
		return dateRange
	}
	if dateRange.IsZero() {
		return otherRange
	}
	minStart := dateRange.Start().Min(otherRange.Start())
	maxEnd := dateRange.End().Max(otherRange.End())
	return NewDateRange(minStart, maxEnd)
}

// Duration computes the duration (in nanoseconds) from midnight at the start of the date
// range up to and including the very last nanosecond before midnight on the end day.
// The calculation is for UTC, which does not have daylight saving and every day has 24 hours.
//
// If the range is greater than approximately 290 years, the result will hard-limit to the
// minimum or maximum possible duration (see time.Sub(t)).
func (dateRange DateRange) Duration() time.Duration {
	return dateRange.End().UTC().Sub(dateRange.Start().UTC())
}

// DurationIn computes the duration (in nanoseconds) from midnight at the start of the date
// range up to and including the very last nanosecond before midnight on the end day.
// The calculation is for the specified location, which may have daylight saving, so not every day
// necessarily has 24 hours. If the date range spans the day the clocks are changed, this is
// taken into account.
//
// If the range is greater than approximately 290 years, the result will hard-limit to the
// minimum or maximum possible duration (see time.Sub(t)).
func (dateRange DateRange) DurationIn(loc *time.Location) time.Duration {
	return dateRange.EndTimeIn(loc).Sub(dateRange.StartTimeIn(loc))
}

// StartTimeIn returns the start time in a specified location.
func (dateRange DateRange) StartTimeIn(loc *time.Location) time.Time {
	return dateRange.Start().In(loc)
}

// EndTimeIn returns the nanosecond after the end time in a specified location. Along with
// StartTimeIn, this gives a 'half-open' range where the start is inclusive and the end is
// exclusive.
func (dateRange DateRange) EndTimeIn(loc *time.Location) time.Time {
	return dateRange.End().In(loc)
}

// TimeSpanIn obtains the time span corresponding to the date range in a specified location.
// The result is normalised.
func (dateRange DateRange) TimeSpanIn(loc *time.Location) TimeSpan {
	s := dateRange.StartTimeIn(loc)
	d := dateRange.DurationIn(loc)
	return TimeSpan{s, d}
}

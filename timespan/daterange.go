// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package timespan

import (
	. "github.com/rickb777/date"
	"time"
	"fmt"
)

const minusOneNano time.Duration = -1

// DateRange carries a date and a number of days and describes a range between two dates.
type DateRange struct {
	mark Date
	days PeriodOfDays
}

// NewDateRangeOf assembles a new date range from a start time and a duration, discarding
// the precise time-of-day information. The start time includes a location, which is not
// necessarily UTC. The duration can be negative.
func NewDateRangeOf(start time.Time, duration time.Duration) DateRange {
	sd := NewAt(start)
	ed := NewAt(start.Add(duration))
	return DateRange{sd, PeriodOfDays(ed.Sub(sd))}
}

// NewDateRange assembles a new date range from two dates.
func NewDateRange(start, end Date) DateRange {
	if end.Before(start) {
		return DateRange{start, PeriodOfDays(end.Sub(start) - 1)}
	}
	return DateRange{start, PeriodOfDays(end.Sub(start) + 1)}
}

// NewYearOf constructs the range encompassing the whole year specified.
func NewYearOf(year int) DateRange {
	start := New(year, time.January, 1)
	end := New(year + 1, time.January, 1)
	return DateRange{start, PeriodOfDays(end.Sub(start))}
}

// NewMonthOf constructs the range encompassing the whole month specified for a given year.
// It handles leap years correctly.
func NewMonthOf(year int, month time.Month) DateRange {
	start := New(year, month, 1)
	endT := time.Date(year, month + 1, 1, 0, 0, 0, 0, time.UTC)
	end := NewAt(endT)
	return DateRange{start, PeriodOfDays(end.Sub(start))}
}

// ZeroRange constructs an empty range. This is often a useful basis for
// further operations but note that the end date is undefined.
func ZeroRange(day Date) DateRange {
	return DateRange{day, 0}
}

// OneDayRange constructs a range of exactly one day. This is often a useful basis for
// further operations. Note that the end date is the same as the start date.
func OneDayRange(day Date) DateRange {
	return DateRange{day, 1}
}

// Days returns the period represented by this range.
func (dateRange DateRange) Days() PeriodOfDays {
	return dateRange.days
}

// Start returns the earliest date represented by this range.
func (dateRange DateRange) Start() Date {
	if dateRange.days < 0 {
		return dateRange.mark.Add(PeriodOfDays(1 + dateRange.days))
	}
	return dateRange.mark
}

// End returns the latest date (inclusive) represented by this range. If the range is empty (i.e.
// has zero days), then an empty date is returned.
func (dateRange DateRange) End() Date {
	if dateRange.days < 0 {
		return dateRange.mark
	} else if dateRange.days == 0 {
		return Date{}
	}
	return dateRange.mark.Add(dateRange.days - 1)
}

// Next returns the date that follows the end date of the range. If the range is empty (i.e.
// has zero days), then an empty date is returned.
func (dateRange DateRange) Next() Date {
	if dateRange.days < 0 {
		return dateRange.mark.Add(1)
	} else if dateRange.days == 0 {
		return Date{}
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
func (dateRange DateRange) ShiftBy(days PeriodOfDays) DateRange {
	if days == 0 {
		return dateRange
	}
	newMark := dateRange.mark.Add(days)
	return DateRange{newMark, dateRange.days}
}

// ExtendBy extends (or reduces) the date range by moving the end date.
// A negative parameter is allowed and this may cause the range to become inverted
// (i.e. the mark date becomes the end date instead of the start date).
func (dateRange DateRange) ExtendBy(days PeriodOfDays) DateRange {
	if days == 0 {
		return dateRange
	}
	return DateRange{dateRange.mark, dateRange.days + days}
}

func (dateRange DateRange) String() string {
	switch dateRange.days {
	case 0:
		return fmt.Sprintf("0 days from %s", dateRange.mark)
	case 1, -1:
		return fmt.Sprintf("1 day on %s", dateRange.mark)
	default:
		if dateRange.days < 0 {
			return fmt.Sprintf("%d days from %s to %s", -dateRange.days, dateRange.Start(), dateRange.End())
		}
		return fmt.Sprintf("%d days from %s to %s", dateRange.days, dateRange.Start(), dateRange.End())
	}
}

// Contains tests whether the date range contains a specified date.
// Empty date ranges (i.e. zero days) never contain anything.
func (dateRange DateRange) Contains(d Date) bool {
	if dateRange.days == 0 {
		return false
	}
	return !(d.Before(dateRange.Start()) || d.After(dateRange.End()))
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
	return dateRange.Next().UTC()
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

// Merge combines two date 	ranges by calculating a date range that just encompasses them both.
// As a special case, if one range is entirely contained within the other range, the larger of
// the two is returned. Otherwise, the result is the start of the earlier one to the end of the
// later one, even if the two ranges don't overlap.
func (dateRange DateRange) Merge(other DateRange) DateRange {
	start := dateRange.Start()
	if start.After(other.Start()) {
		// swap the ranges to simplify the logic
		return other.Merge(dateRange)

	} else {
		oEnd := other.End()
		if dateRange.End().After(oEnd) {
			// other is a proper subrange of dateRange
			return dateRange

		} else {
			return NewDateRange(start, oEnd)
		}
	}
}

// Duration computes the duration (in nanoseconds) from midnight at the start of the date
// range up to and including the very last nanosecond before midnight the following day after the end.
// The calculation is for UTC, which does not have daylight saving and every day has 24 hours.
//
// If the range is greater than approximately 290 years, the result will hard-limit to the
// minimum or maximum possible duration (see time.Sub(t)).
func (dateRange DateRange) Duration() time.Duration {
	return dateRange.Next().UTC().Sub(dateRange.Start().UTC())
}

// DurationIn computes the duration (in nanoseconds) from midnight at the start of the date
// range up to and including the very last nanosecond before midnight the following day after the end.
// The calculation is for the specified location, which may have daylight saving, so not every day has
// 24 hours. If the date range spans the day the clocks are changed, this is taken into account.
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
	return dateRange.Next().In(loc)
}

// TimeSpanIn obtains the time span corresponding to the date range in a specified location.
// The result is normalised.
func (dateRange DateRange) TimeSpanIn(loc *time.Location) TimeSpan {
	dr := dateRange.Normalise()
	s := dr.StartTimeIn(loc)
	d := dr.DurationIn(loc)
	return TimeSpan{s, d}
}


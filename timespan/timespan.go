// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package timespan

import (
	"time"
	"fmt"
	"github.com/rickb777/date"
)

const TimestampFormat = "2006-01-02 15:04:05"
//const ISOFormat = "2006-01-02T15:04:05"

type TimeSpan struct {
	Mark     time.Time
	Duration time.Duration
}

// ZeroTimeSpan creates a new zero-duration time span from a time.
func ZeroTimeSpan(start time.Time) TimeSpan {
	return TimeSpan{start, 0}
}

// NewTimeSpan creates a new time span from two times. The start and end can be in either order;
// the result will be normalised.
func NewTimeSpan(t1, t2 time.Time) TimeSpan {
	if t2.Before(t1) {
		return TimeSpan{t2, t1.Sub(t2)}
	}
	return TimeSpan{t1, t2.Sub(t1)}
}

// End gets the end time of the time span.
func (ts TimeSpan) Start() time.Time {
	if ts.Duration < 0 {
		return ts.Mark.Add(ts.Duration)
	}
	return ts.Mark
}

// End gets the end time of the time span.
func (ts TimeSpan) End() time.Time {
	if ts.Duration < 0 {
		return ts.Mark
	}
	return ts.Mark.Add(ts.Duration)
}

// Normalise ensures that the mark time is at the start time and the duration is positive.
// The normalised timespan is returned.
func (ts TimeSpan) Normalise() TimeSpan {
	if ts.Duration < 0 {
		return TimeSpan{ts.Mark.Add(ts.Duration), -ts.Duration}
	}
	return ts
}

// ShiftBy moves the date range by moving both the start and end times similarly.
// A negative parameter is allowed.
func (ts TimeSpan) ShiftBy(d time.Duration) TimeSpan {
	return TimeSpan{ts.Mark.Add(d), ts.Duration}
}

// ExtendBy lengthens the time span by a specified amount. The parameter may be negative,
// in which case it is possible that the end of the time span will appear to be before the
// start. However, the result is normalised so that the resulting start is the lesser value
// and the duration is always non-negative.
func (ts TimeSpan) ExtendBy(d time.Duration) TimeSpan {
	return TimeSpan{ts.Mark, ts.Duration + d}.Normalise()
}

// ExtendWithoutWrapping lengthens the time span by a specified amount. The parameter may be
// negative, but if its magnitude is large than the time span's duration, it will be truncated
// so that the result has zero duration in that case. The start time is never altered.
func (ts TimeSpan) ExtendWithoutWrapping(d time.Duration) TimeSpan {
	if d < 0 && -d > ts.Duration {
		return TimeSpan{ts.Mark, 0}
	}
	return TimeSpan{ts.Mark, ts.Duration + d}
}

func (ts TimeSpan) String() string {
	return fmt.Sprintf("%s from %s to %s", ts.Duration, ts.Mark.Format(TimestampFormat), ts.End().Format(TimestampFormat))
}

// In returns a TimeSpan adjusted from its current location to a new location. Because
// location is considered to be a presentational attribute, the actual time itself is not
// altered by this function. This matches the behaviour of time.Time.In(loc).
func (ts TimeSpan) In(loc *time.Location) TimeSpan {
	t := ts.Mark.In(loc)
	return TimeSpan{t, ts.Duration}
}

// DateRangeIn obtains the date range corresponding to the time span in a specified location.
// The result is normalised.
func (ts TimeSpan) DateRangeIn(loc *time.Location) DateRange {
	no := ts.Normalise()
	startDate := date.NewAt(no.Mark.In(loc))
	endDate := date.NewAt(no.End().In(loc))
	return NewDateRange(startDate, endDate)
}

// Contains tests whether a given moment of time is enclosed within the time span. The
// start time is inclusive; the end time is exclusive.
// If t has a different locality to the time-span, it is adjusted accordingly.
func (ts TimeSpan) Contains(t time.Time) bool {
	tl := t.In(ts.Mark.Location())
	return ts.Mark.Equal(tl) || ts.Mark.Before(tl) && ts.End().After(tl)
}

// Merge combines two time spans by calculating a time span that just encompasses them both.
// As a special case, if one span is entirely contained within the other span, the larger of
// the two is returned. Otherwise, the result is the start of the earlier one to the end of the
// later one, even if the two spans don't overlap.
func (ts TimeSpan) Merge(other TimeSpan) TimeSpan {
	if ts.Mark.After(other.Mark) {
		// swap the ranges to simplify the logic
		return other.Merge(ts)

	} else if ts.End().After(other.End()) {
		// other is a proper subrange of ts
		return ts

	} else {
		return NewTimeSpan(ts.Mark, other.End())
	}
}

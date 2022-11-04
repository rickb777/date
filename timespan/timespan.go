// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package timespan

import (
	"fmt"
	"strings"
	"time"

	"github.com/rickb777/date"
	"github.com/rickb777/date/period"
)

// TimestampFormat is a simple format for date & time, "2006-01-02 15:04:05".
const TimestampFormat = "2006-01-02 15:04:05"

//const ISOFormat = "2006-01-02T15:04:05"

// TimeSpan holds a span of time between two instants with a 1 nanosecond resolution.
// It is implemented using a time.Duration, therefore is limited to a maximum span of 290 years.
type TimeSpan struct {
	mark     time.Time
	duration time.Duration
}

// ZeroTimeSpan creates a new zero-duration time span at a specified time.
func ZeroTimeSpan(start time.Time) TimeSpan {
	return TimeSpan{start, 0}
}

// TimeSpanOf creates a new time span at a specified time and duration.
func TimeSpanOf(start time.Time, d time.Duration) TimeSpan {
	return TimeSpan{start, d}
}

// NewTimeSpan creates a new time span from two times. The start and end can be in either
// order; the result will be normalised. The inputs are half-open: the start is included and
// the end is excluded.
func NewTimeSpan(t1, t2 time.Time) TimeSpan {
	if t2.Before(t1) {
		return TimeSpan{t2, t1.Sub(t2)}
	}
	return TimeSpan{t1, t2.Sub(t1)}
}

// Start gets the end time of the time span.
func (ts TimeSpan) Start() time.Time {
	if ts.duration < 0 {
		return ts.mark.Add(ts.duration)
	}
	return ts.mark
}

// End gets the end time of the time span. Strictly, this is one nanosecond after the
// range of time included in the time span; this implements the half-open model.
func (ts TimeSpan) End() time.Time {
	if ts.duration < 0 {
		return ts.mark
	}
	return ts.mark.Add(ts.duration)
}

// Duration gets the duration of the time span.
func (ts TimeSpan) Duration() time.Duration {
	return ts.duration
}

// IsEmpty returns true if this is an empty time span (zero duration).
func (ts TimeSpan) IsEmpty() bool {
	return ts.duration == 0
}

// Normalise ensures that the mark time is at the start time and the duration is positive.
// The normalised time span is returned.
func (ts TimeSpan) Normalise() TimeSpan {
	if ts.duration < 0 {
		return TimeSpan{ts.mark.Add(ts.duration), -ts.duration}
	}
	return ts
}

// ShiftBy moves the time span by moving both the start and end times similarly.
// A negative parameter is allowed.
func (ts TimeSpan) ShiftBy(d time.Duration) TimeSpan {
	return TimeSpan{ts.mark.Add(d), ts.duration}
}

// ExtendBy lengthens the time span by a specified amount. The parameter may be negative,
// in which case it is possible that the end of the time span will appear to be before the
// start. However, the result is normalised so that the resulting start is the lesser value.
func (ts TimeSpan) ExtendBy(d time.Duration) TimeSpan {
	return TimeSpan{ts.mark, ts.duration + d}.Normalise()
}

// ExtendWithoutWrapping lengthens the time span by a specified amount. The parameter may be
// negative, but if its magnitude is large than the time span's duration, it will be truncated
// so that the result has zero duration in that case. The start time is never altered.
func (ts TimeSpan) ExtendWithoutWrapping(d time.Duration) TimeSpan {
	tsn := ts.Normalise()
	if d < 0 && -d > tsn.duration {
		return TimeSpan{tsn.mark, 0}
	}
	return TimeSpan{tsn.mark, tsn.duration + d}
}

// String produces a human-readable description of a time span.
func (ts TimeSpan) String() string {
	return fmt.Sprintf("%s from %s to %s", ts.duration, ts.mark.Format(TimestampFormat), ts.End().Format(TimestampFormat))
}

// In returns a TimeSpan adjusted from its current location to a new location. Because
// location is considered to be a presentational attribute, the actual time itself is not
// altered by this function. This matches the behaviour of time.Time.In(loc).
func (ts TimeSpan) In(loc *time.Location) TimeSpan {
	t := ts.mark.In(loc)
	return TimeSpan{t, ts.duration}
}

// DateRangeIn obtains the date range corresponding to the time span in a specified location.
// The result is normalised.
func (ts TimeSpan) DateRangeIn(loc *time.Location) DateRange {
	no := ts.Normalise()
	startDate := date.NewAt(no.mark.In(loc))
	endDate := date.NewAt(no.End().In(loc))
	return NewDateRange(startDate, endDate)
}

// Contains tests whether a given moment of time is enclosed within the time span. The
// start time is inclusive; the end time is exclusive.
// If t has a different locality to the time-span, it is adjusted accordingly.
func (ts TimeSpan) Contains(t time.Time) bool {
	tl := t.In(ts.mark.Location())
	return ts.mark.Equal(tl) || ts.mark.Before(tl) && ts.End().After(tl)
}

// Merge combines two time spans by calculating a time span that just encompasses them both.
// As a special case, if one span is entirely contained within the other span, the larger of
// the two is returned. Otherwise, the result is the start of the earlier one to the end of the
// later one, even if the two spans don't overlap.
func (ts TimeSpan) Merge(other TimeSpan) TimeSpan {
	if ts.mark.After(other.mark) {
		// swap the ranges to simplify the logic
		return other.Merge(ts)

	} else if ts.End().After(other.End()) {
		// other is a proper subrange of ts
		return ts

	} else {
		return NewTimeSpan(ts.mark, other.End())
	}
}

// RFC5545DateTimeLayout is the format string used by iCalendar (RFC5545). Note
// that "Z" is to be appended when the time is UTC.
const RFC5545DateTimeLayout = "20060102T150405"

// RFC5545DateTimeZulu is the UTC format string used by iCalendar (RFC5545). Note
// that this cannot be used for parsing with time.Parse.
const RFC5545DateTimeZulu = RFC5545DateTimeLayout + "Z"

func layoutHasTimezone(layout string) bool {
	return strings.IndexByte(layout, 'Z') >= 0 || strings.Contains(layout, "-07")
}

// Equal reports whether ts and us represent the same time start and duration.
// Two times can be equal even if they are in different locations.
// For example, 6:00 +0200 CEST and 4:00 UTC are Equal.
func (ts TimeSpan) Equal(us TimeSpan) bool {
	return ts.Duration() == us.Duration() && ts.Start().Equal(us.Start())
}

// Format returns a textual representation of the time value formatted according to layout.
// It produces a string containing the start and end time. Or, if useDuration is true,
// it returns a string containing the start time and the duration.
//
// The layout string is as specified for time.Format. If it doesn't have a timezone element
// ("07" or "Z") and the times in the timespan are UTC, the "Z" zulu indicator is added.
// This is as required by iCalendar (RFC5545).
//
// Also, if the layout is blank, it defaults to RFC5545DateTimeLayout.
//
// The separator between the two parts of the result would be "/" for RFC5545, but can be
// anything.
func (ts TimeSpan) Format(layout, separator string, useDuration bool) string {
	if layout == "" {
		layout = RFC5545DateTimeLayout
	}

	// if the time is UTC and the format doesn't contain zulu field ("Z") or timezone field ("07")
	if ts.mark.Location().String() == "UTC" && !layoutHasTimezone(layout) {
		layout = RFC5545DateTimeZulu
	}

	s := ts.Start()
	e := ts.End()

	if useDuration {
		p := period.Between(s, e)
		return fmt.Sprintf("%s%s%s", s.Format(layout), separator, p)
	}

	return fmt.Sprintf("%s%s%s", s.Format(layout), separator, e.Format(layout))
}

// FormatRFC5545 formats the timespan as a string containing the start time and end time, or the
// start time and duration, if useDuration is true. The two parts are separated by slash.
// The time(s) is expressed as UTC zulu.
// This is as required by iCalendar (RFC5545).
func (ts TimeSpan) FormatRFC5545(useDuration bool) string {
	return ts.Format(RFC5545DateTimeZulu, "/", useDuration)
}

// MarshalText formats the timespan as a string using, using RFC5545 layout.
// This implements the encoding.TextMarshaler interface.
func (ts TimeSpan) MarshalText() (text []byte, err error) {
	s := ts.Format(RFC5545DateTimeZulu, "/", true)
	return []byte(s), nil
}

// ParseRFC5545InLocation parses a string as a timespan. The string must contain either of
//
//	time "/" time
//	time "/" period
//
// If the input time(s) ends in "Z", the location is UTC (as per RFC5545). Otherwise, the
// specified location will be used for the resulting times; this behaves the same as
// time.ParseInLocation.
func ParseRFC5545InLocation(text string, loc *time.Location) (TimeSpan, error) {
	slash := strings.IndexByte(text, '/')
	if slash < 0 {
		return TimeSpan{}, fmt.Errorf("cannot parse %q because there is no separator '/'", text)
	}

	start := text[:slash]
	rest := text[slash+1:]

	st, err := parseTimeInLocation(start, loc)
	if err != nil {
		return TimeSpan{}, fmt.Errorf("cannot parse start time in %q: %s", text, err.Error())
	}

	//fmt.Printf("got %20s %s\n", st.Location(), st.Format(RFC5545DateTimeLayout))

	if rest == "" {
		return TimeSpan{}, fmt.Errorf("cannot parse %q because there is end time or duration", text)
	}

	if rest[0] == 'P' {
		pe, e2 := period.Parse(rest)
		if e2 != nil {
			return TimeSpan{}, fmt.Errorf("cannot parse period in %q: %s", text, e2.Error())
		}

		du, precise := pe.Duration()
		if precise {
			return TimeSpan{st, du}, nil
		}

		et := st.AddDate(pe.Years(), pe.Months(), pe.Days())
		return NewTimeSpan(st, et), nil
	}

	et, err := parseTimeInLocation(rest, loc)
	return NewTimeSpan(st, et), err
}

func parseTimeInLocation(text string, loc *time.Location) (time.Time, error) {
	if strings.HasSuffix(text, "Z") {
		text = text[:len(text)-1]
		return time.ParseInLocation(RFC5545DateTimeLayout, text, time.UTC)
	}
	return time.ParseInLocation(RFC5545DateTimeLayout, text, loc)
}

// UnmarshalText parses a string as a timespan. It expects RFC5545 layout.
//
// If the receiver timespan is non-nil and has a time with a location,
// this location is used for parsing. Otherwise time.Local is used.
//
// This implements the encoding.TextUnmarshaler interface.
func (ts *TimeSpan) UnmarshalText(text []byte) (err error) {
	loc := time.Local
	if ts != nil {
		loc = ts.mark.Location()
	}
	*ts, err = ParseRFC5545InLocation(string(text), loc)
	return
}

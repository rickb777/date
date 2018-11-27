// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package timespan

import (
	"testing"
	"time"

	"github.com/rickb777/date"
)

const zero time.Duration = 0

var t0327 = time.Date(2015, 3, 27, 0, 0, 0, 0, time.UTC)
var t0328 = time.Date(2015, 3, 28, 0, 0, 0, 0, time.UTC)
var t0329 = time.Date(2015, 3, 29, 0, 0, 0, 0, time.UTC) // n.b. clocks go forward (UK)
var t0330 = time.Date(2015, 3, 30, 0, 0, 0, 0, time.UTC)

func TestZeroTimeSpan(t *testing.T) {
	ts := ZeroTimeSpan(t0327)
	isEq(t, 0, ts.mark, t0327)
	isEq(t, 0, ts.Duration(), zero)
	isEq(t, 0, ts.End(), t0327)
}

func TestNewTimeSpan(t *testing.T) {
	ts1 := NewTimeSpan(t0327, t0327)
	isEq(t, 0, ts1.mark, t0327)
	isEq(t, 0, ts1.Duration(), zero)
	isEq(t, 0, ts1.IsEmpty(), true)
	isEq(t, 0, ts1.End(), t0327)

	ts2 := NewTimeSpan(t0327, t0328)
	isEq(t, 0, ts2.mark, t0327)
	isEq(t, 0, ts2.Duration(), time.Hour*24)
	isEq(t, 0, ts2.IsEmpty(), false)
	isEq(t, 0, ts2.End(), t0328)

	ts3 := NewTimeSpan(t0329, t0327)
	isEq(t, 0, ts3.mark, t0327)
	isEq(t, 0, ts3.Duration(), time.Hour*48)
	isEq(t, 0, ts3.IsEmpty(), false)
	isEq(t, 0, ts3.End(), t0329)
}

func TestTSEnd(t *testing.T) {
	ts1 := TimeSpan{t0328, time.Hour * 24}
	isEq(t, 0, ts1.Start(), t0328)
	isEq(t, 0, ts1.End(), t0329)

	// not normalised, deliberately
	ts2 := TimeSpan{t0328, -time.Hour * 24}
	isEq(t, 0, ts2.Start(), t0327)
	isEq(t, 0, ts2.End(), t0328)
}

func TestTSShiftBy(t *testing.T) {
	ts1 := NewTimeSpan(t0327, t0328).ShiftBy(time.Hour * 24)
	isEq(t, 0, ts1.mark, t0328)
	isEq(t, 0, ts1.Duration(), time.Hour*24)
	isEq(t, 0, ts1.End(), t0329)

	ts2 := NewTimeSpan(t0328, t0329).ShiftBy(-time.Hour * 24)
	isEq(t, 0, ts2.mark, t0327)
	isEq(t, 0, ts2.Duration(), time.Hour*24)
	isEq(t, 0, ts2.End(), t0328)
}

func TestTSExtendBy(t *testing.T) {
	ts1 := NewTimeSpan(t0327, t0328).ExtendBy(time.Hour * 24)
	isEq(t, 0, ts1.mark, t0327)
	isEq(t, 0, ts1.Duration(), time.Hour*48)
	isEq(t, 0, ts1.End(), t0329)

	ts2 := NewTimeSpan(t0328, t0329).ExtendBy(-time.Hour * 48)
	isEq(t, 0, ts2.mark, t0327)
	isEq(t, 0, ts2.Duration(), time.Hour*24)
	isEq(t, 0, ts2.End(), t0328)
}

func TestTSExtendWithoutWrapping(t *testing.T) {
	ts1 := NewTimeSpan(t0327, t0328).ExtendWithoutWrapping(time.Hour * 24)
	isEq(t, 0, ts1.mark, t0327)
	isEq(t, 0, ts1.Duration(), time.Hour*48)
	isEq(t, 0, ts1.End(), t0329)

	ts2 := NewTimeSpan(t0328, t0329).ExtendWithoutWrapping(-time.Hour * 48)
	isEq(t, 0, ts2.mark, t0328)
	isEq(t, 0, ts2.Duration(), zero)
	isEq(t, 0, ts2.End(), t0328)
}

func TestTSString(t *testing.T) {
	s := NewTimeSpan(t0327, t0328).String()
	isEq(t, 0, s, "24h0m0s from 2015-03-27 00:00:00 to 2015-03-28 00:00:00")
}

func TestTSEqual(t *testing.T) {
	// use Berlin, which is UTC+1/+2
	berlin, _ := time.LoadLocation("Europe/Berlin")
	t0 := time.Date(2015, 2, 20, 10, 13, 25, 0, time.UTC)
	t1 := t0.Add(time.Hour)
	z0 := ZeroTimeSpan(t0)
	ts1 := z0.ExtendBy(time.Hour)

	cases := []struct {
		a, b TimeSpan
	}{
		{z0, NewTimeSpan(t0, t0)},
		{z0, z0.In(berlin)},
		{ts1, ts1},
		{ts1, NewTimeSpan(t0, t1)},
		{ts1, ts1.In(berlin)},
		{ts1, ZeroTimeSpan(t1).ExtendBy(-time.Hour)},
	}

	for i, c := range cases {
		if !c.a.Equal(c.b) {
			t.Errorf("%d: %v is not equal to %v", i, c.a, c.b)
		}
	}
}

func TestTSNotEqual(t *testing.T) {
	t0 := time.Date(2015, 2, 20, 10, 13, 25, 0, time.UTC)
	t1 := t0.Add(time.Hour)

	cases := []struct {
		a, b TimeSpan
	}{
		{ZeroTimeSpan(t0), TimeSpanOf(t0, time.Hour)},
		{ZeroTimeSpan(t0), ZeroTimeSpan(t1)},
	}

	for i, c := range cases {
		if c.a.Equal(c.b) {
			t.Errorf("%d: %v is not equal to %v", i, c.a, c.b)
		}
	}
}

func TestTSFormat(t *testing.T) {
	// use Berlin, which is UTC-1
	berlin, _ := time.LoadLocation("Europe/Berlin")
	t0 := time.Date(2015, 3, 27, 10, 13, 14, 0, time.UTC)

	cases := []struct {
		start                  time.Time
		duration               time.Duration
		useDuration            bool
		layout, separator, exp string
	}{
		{t0, time.Hour, true, "", " for ", "20150327T101314Z for PT1H"},
		{t0, time.Hour, true, "", "/", "20150327T101314Z/PT1H"},
		{t0.In(berlin), time.Minute, true, "", "/", "20150327T111314/PT1M"},
		{t0.In(berlin), time.Hour, true, "2006-01-02T15:04:05", "/", "2015-03-27T11:13:14/PT1H"},
		{t0.In(berlin), time.Hour, true, "2006-01-02T15:04:05-07", "/", "2015-03-27T11:13:14+01/PT1H"},
		{t0, time.Hour, true, "2006-01-02T15:04:05-07", "/", "2015-03-27T10:13:14+00/PT1H"},
		{t0, time.Hour, true, "2006-01-02T15:04:05Z07", "/", "2015-03-27T10:13:14Z/PT1H"},

		{t0, time.Hour, false, "", " to ", "20150327T101314Z to 20150327T111314Z"},
		{t0, time.Hour, false, "", "/", "20150327T101314Z/20150327T111314Z"},
		{t0.In(berlin), time.Minute, false, "", "/", "20150327T111314/20150327T111414"},
		{t0.In(berlin), time.Hour, false, "2006-01-02T15:04:05", "/", "2015-03-27T11:13:14/2015-03-27T12:13:14"},
		{t0.In(berlin), time.Hour, false, "2006-01-02T15:04:05-07", "/", "2015-03-27T11:13:14+01/2015-03-27T12:13:14+01"},
		{t0, time.Hour, false, "2006-01-02T15:04:05-07", "/", "2015-03-27T10:13:14+00/2015-03-27T11:13:14+00"},
		{t0, time.Hour, false, "2006-01-02T15:04:05Z07", "/", "2015-03-27T10:13:14Z/2015-03-27T11:13:14Z"},
	}

	for _, c := range cases {
		ts := TimeSpan{c.start, c.duration}
		isEq(t, 0, ts.Format(c.layout, c.separator, c.useDuration), c.exp)
	}
}

func TestTSMarshalText(t *testing.T) {
	// use Berlin, which is UTC+1 or +2 in summer
	berlin, _ := time.LoadLocation("Europe/Berlin")
	t0 := time.Date(2015, 2, 14, 10, 13, 14, 0, time.UTC)
	t1 := time.Date(2015, 6, 27, 10, 13, 15, 0, time.UTC)

	cases := []struct {
		start    time.Time
		duration time.Duration
		exp      string
	}{
		{t0, time.Hour, "20150214T101314Z/PT1H"},
		{t1, 2 * time.Hour, "20150627T101315Z/PT2H"},
		{t0.In(berlin), time.Minute, "20150214T111314Z/PT1M"}, // UTC+1
		{t1.In(berlin), time.Second, "20150627T121315Z/PT1S"}, // UTC+2
	}

	for i, c := range cases {
		ts := TimeSpan{c.start, c.duration}

		s := ts.FormatRFC5545(true)
		isEq(t, i, s, c.exp)

		b, err := ts.MarshalText()
		isEq(t, i, err, nil)
		isEq(t, i, string(b), c.exp)
	}
}

func TestTSParseInLocation(t *testing.T) {
	// use Berlin, which is UTC-1
	berlin, _ := time.LoadLocation("Europe/Berlin")
	t0120 := time.Date(2015, 1, 20, 10, 13, 14, 0, time.UTC)
	// just before start of daylight savings
	t0325 := time.Date(2015, 3, 25, 10, 13, 14, 0, time.UTC)

	cases := []struct {
		start    time.Time
		duration time.Duration
		text     string
	}{
		{t0325, time.Hour, "20150325T101314Z/PT1H"},
		{t0325, 2 * time.Second, "20150325T101314Z/PT2S"},
		{t0120.In(berlin), time.Minute, "20150120T111314/PT1M"},
		{t0325, 336 * time.Hour, "20150325T101314Z/P2W"},
		{t0120.In(berlin), 72 * time.Hour, "20150120T111314/P3D"},
		// This case has the daylight-savings clock shift
		{t0325.In(berlin), 167 * time.Hour, "20150325T111314/P1W"},
	}

	for i, c := range cases {
		ts1, err := ParseRFC5545InLocation(c.text, c.start.Location())
		if err != nil {
			t.Errorf("%d: %s %v %v", i, c.text, ts1.String(), err)
		}

		if !ts1.Start().Equal(c.start) {
			t.Errorf("%d: %s", i, ts1)
		}

		if ts1.Duration() != c.duration {
			t.Errorf("%d: %s", i, ts1)
		}

		ts2 := TimeSpan{}.In(c.start.Location())
		err = ts2.UnmarshalText([]byte(c.text))
		if err != nil {
			t.Errorf("%d: %s: %v %v", i, c.text, ts2.String(), err)
		}

		if !ts1.Equal(ts2) {
			t.Errorf("%d: %s: %v is not equal to %v", i, c.text, ts1, ts2)
		}
	}
}

func TestTSParseInLocationErrors(t *testing.T) {
	cases := []struct {
		text string
	}{
		{"20150327T101314Z PT1H"},
		{"2015XX27T101314/PT1H"},
		{"20150127T101314/2016XX27T101314"},
		{"20150127T101314/P1Z"},
		{"20150327T101314Z/"},
		{"/PT1H"},
	}

	for _, c := range cases {
		ts, err := ParseRFC5545InLocation(c.text, time.UTC)
		if err == nil {
			t.Errorf(ts.String())
		}
	}
}

func TestTSContains(t *testing.T) {
	ts := NewTimeSpan(t0327, t0329)
	isEq(t, 0, ts.Contains(t0327.Add(minusOneNano)), false)
	isEq(t, 0, ts.Contains(t0327), true)
	isEq(t, 0, ts.Contains(t0328), true)
	isEq(t, 0, ts.Contains(t0329.Add(minusOneNano)), true)
	isEq(t, 0, ts.Contains(t0329), false)
}

func TestTSIn(t *testing.T) {
	ts := ZeroTimeSpan(t0327).In(time.FixedZone("Test", 7200))
	isEq(t, 0, ts.mark.Equal(t0327), true)
	isEq(t, 0, ts.Duration(), zero)
	isEq(t, 0, ts.End().Equal(t0327), true)
}

func TestTSMerge1(t *testing.T) {
	ts1 := NewTimeSpan(t0327, t0328)
	ts2 := NewTimeSpan(t0327, t0330)
	m1 := ts1.Merge(ts2)
	m2 := ts2.Merge(ts1)
	isEq(t, 0, m1.mark, t0327)
	isEq(t, 0, m1.End(), t0330)
	isEq(t, 0, m1, m2)
}

func TestTSMerge2(t *testing.T) {
	ts1 := NewTimeSpan(t0328, t0329)
	ts2 := NewTimeSpan(t0327, t0330)
	m1 := ts1.Merge(ts2)
	m2 := ts2.Merge(ts1)
	isEq(t, 0, m1.mark, t0327)
	isEq(t, 0, m1.End(), t0330)
	isEq(t, 0, m1, m2)
}

func TestTSMerge3(t *testing.T) {
	ts1 := NewTimeSpan(t0329, t0330)
	ts2 := NewTimeSpan(t0327, t0330)
	m1 := ts1.Merge(ts2)
	m2 := ts2.Merge(ts1)
	isEq(t, 0, m1.mark, t0327)
	isEq(t, 0, m1.End(), t0330)
	isEq(t, 0, m1, m2)
}

func TestTSMergeOverlapping(t *testing.T) {
	ts1 := NewTimeSpan(t0327, t0329)
	ts2 := NewTimeSpan(t0328, t0330)
	m1 := ts1.Merge(ts2)
	m2 := ts2.Merge(ts1)
	isEq(t, 0, m1.mark, t0327)
	isEq(t, 0, m1.End(), t0330)
	isEq(t, 0, m1, m2)
}

func TestTSMergeNonOverlapping(t *testing.T) {
	ts1 := NewTimeSpan(t0327, t0328)
	ts2 := NewTimeSpan(t0329, t0330)
	m1 := ts1.Merge(ts2)
	m2 := ts2.Merge(ts1)
	isEq(t, 0, m1.mark, t0327)
	isEq(t, 0, m1.End(), t0330)
	isEq(t, 0, m1, m2)
}

func TestConversion1(t *testing.T) {
	ts1 := ZeroTimeSpan(t0327)
	dr := ts1.DateRangeIn(time.UTC)
	ts2 := dr.TimeSpanIn(time.UTC)
	isEq(t, 0, dr.Start(), d0327)
	isEq(t, 0, dr.IsEmpty(), true)
	isEq(t, 0, ts1.Start(), ts1.End())
	isEq(t, 0, ts1.Duration(), zero)
	isEq(t, 0, dr.Days(), date.PeriodOfDays(0))
	isEq(t, 0, ts2.Duration(), zero)
	isEq(t, 0, ts1, ts2)
}

func TestConversion2(t *testing.T) {
	ts1 := NewTimeSpan(t0327, t0328)
	dr := ts1.DateRangeIn(time.UTC)
	ts2 := dr.TimeSpanIn(time.UTC)
	isEq(t, 0, dr.Start(), d0327)
	isEq(t, 0, dr.End(), d0328)
	isEq(t, 0, ts1, ts2)
	isEq(t, 0, ts1.Duration(), time.Hour*24)
}

func TestConversion3(t *testing.T) {
	dr1 := NewDateRange(d0327, d0330) // weekend of clocks changing
	ts1 := dr1.TimeSpanIn(london)
	dr2 := ts1.DateRangeIn(london)
	ts2 := dr2.TimeSpanIn(london)
	isEq(t, 0, dr1.Start(), d0327)
	isEq(t, 0, dr1.End(), d0330)
	isEq(t, 0, dr1, dr2)
	isEq(t, 0, ts1, ts2)
	isEq(t, 0, ts1.Duration(), time.Hour*71)
}

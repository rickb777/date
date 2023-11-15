// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package date

import (
	"runtime/debug"
	"testing"
	"time"

	"github.com/rickb777/period"
)

func same(d Date, t time.Time) bool {
	yd, wd := d.ISOWeek()
	yt, wt := t.ISOWeek()
	return d.Year() == t.Year() &&
		d.Month() == t.Month() &&
		d.Day() == t.Day() &&
		d.Weekday() == t.Weekday() &&
		d.YearDay() == t.YearDay() &&
		yd == yt && wd == wt
}

func TestDate_New(t *testing.T) {
	cases := []string{
		"0000-01-01T00:00:00+00:00",
		"0001-01-01T00:00:00+00:00",
		"1614-01-01T01:02:03+04:00",
		"1970-01-01T00:00:00+00:00",
		"1815-12-10T05:06:07+00:00",
		"1900-01-01T00:00:00+00:00",
		"1901-09-10T00:00:00-05:00",
		"1998-09-01T00:00:00-08:00",
		"2000-01-01T00:00:00+00:00",
		"9999-12-31T00:00:00+00:00",
	}
	for i, c := range cases {
		tIn, err := time.Parse(time.RFC3339, c)
		if err != nil {
			t.Errorf("%d: New(%v) cannot parse input: %v", i, c, err)
			continue
		}
		dOut := New(tIn.Year(), tIn.Month(), tIn.Day())
		if !same(dOut, tIn) {
			t.Errorf("%d: New(%v) == %v, want date of %v", i, c, dOut, tIn)
		}
		dOut = NewAt(tIn)
		if !same(dOut, tIn) {
			t.Errorf("%d: NewAt(%v) == %v, want date of %v", i, c, dOut, tIn)
		}
	}
}

func TestDate_Today(t *testing.T) {
	today := Today()
	now := time.Now()
	if !same(today, now) {
		t.Errorf("Today == %v, want date of %v", today, now)
	}
	today = TodayUTC()
	now = time.Now().UTC()
	if !same(today, now) {
		t.Errorf("TodayUTC == %v, want date of %v", today, now)
	}
	cases := []int{-10, -5, -3, 0, 1, 4, 8, 12}
	for _, c := range cases {
		location := time.FixedZone("zone", c*60*60)
		today = TodayIn(location)
		now = time.Now().In(location)
		if !same(today, now) {
			t.Errorf("TodayIn(%v) == %v, want date of %v", c, today, now)
		}
	}
}

func TestDate_Time(t *testing.T) {
	cases := []struct {
		d Date
	}{
		{New(-1234, time.February, 5)},
		{New(-1, time.January, 1)},
		{New(0, time.April, 12)},
		{New(1, time.January, 1)},
		{New(1946, time.February, 4)},
		{New(1970, time.January, 1)},
		{New(1976, time.April, 1)},
		{New(1999, time.December, 1)},
		{New(1111111, time.June, 21)},
	}
	zones := []int{-12, -10, -5, -4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 8, 12}
	for i, c := range cases {
		d := c.d
		tUTC := d.UTC()
		if !same(d, tUTC) {
			t.Errorf("%d: TimeUTC(%v) == %v, want date part %v", i, d, tUTC, d)
		}
		if tUTC.Location() != time.UTC {
			t.Errorf("%d: TimeUTC(%v) == %v, want %v", i, d, tUTC.Location(), time.UTC)
		}
		tLocal := d.Local()
		if !same(d, tLocal) {
			t.Errorf("%d: TimeLocal(%v) == %v, want date part %v", i, d, tLocal, d)
		}
		if tLocal.Location() != time.Local {
			t.Errorf("%d: TimeLocal(%v) == %v, want %v", i, d, tLocal.Location(), time.Local)
		}
		for _, z := range zones {
			location := time.FixedZone("zone", z*60*60)
			tInLoc := d.In(location)
			if !same(d, tInLoc) {
				t.Errorf("%d: TimeIn(%v) == %v, want date part %v", i, d, tInLoc, d)
			}
			if tInLoc.Location() != location {
				t.Errorf("%d: TimeIn(%v) == %v, want %v", i, d, tInLoc.Location(), location)
			}
		}
	}
}

func TestPredicates(t *testing.T) {
	// The list of case dates must be sorted in ascending order
	cases := []struct {
		d Date
	}{
		{New(-1234, time.February, 5)},
		{New(0, time.April, 12)},
		{New(1, time.January, 1)},
		{New(1946, time.February, 4)},
		{New(1970, time.January, 1)},
		{New(1976, time.April, 1)},
		{New(1999, time.December, 1)},
		{New(1111111, time.June, 21)},
	}
	for i, ci := range cases {
		di := ci.d
		for j, cj := range cases {
			dj := cj.d
			testPredicate(t, di, dj, di == dj, i == j, "==")
			testPredicate(t, di, dj, di != dj, i != j, "!=")
		}
	}
}

func testPredicate(t *testing.T, di, dj Date, p, q bool, m string) {
	if p != q {
		t.Errorf("%s(%v, %v) == %v, want %v\n%v", m, di, dj, p, q, debug.Stack())
	}
}

func TestDate_AddDate(t *testing.T) {
	cases := []struct {
		d                   Date
		years, months, days int
		expected            Date
	}{
		{New(1970, time.January, 1), 1, 2, 3, New(1971, time.March, 4)},
		{New(1999, time.September, 28), 6, 4, 2, New(2006, time.January, 30)},
		{New(1999, time.September, 28), 0, 0, 3, New(1999, time.October, 1)},
		{New(1999, time.September, 28), 0, 1, 3, New(1999, time.October, 31)},
	}
	for _, c := range cases {
		di := c.d
		dj := di.AddDate(c.years, c.months, c.days)
		if dj != c.expected {
			t.Errorf("%v AddDate(%v,%v,%v) == %v, want %v", di, c.years, c.months, c.days, dj, c.expected)
		}
		dk := dj.AddDate(-c.years, -c.months, -c.days)
		if dk != di {
			t.Errorf("%v AddDate(%v,%v,%v) == %v, want %v", dj, -c.years, -c.months, -c.days, dk, di)
		}
	}
}

func TestDate_AddPeriod(t *testing.T) {
	cases := []struct {
		in       Date
		delta    period.Period
		expected Date
	}{
		{New(1970, time.January, 1), period.NewYMWD(0, 0, 0, 0), New(1970, time.January, 1)},
		{New(1971, time.January, 1), period.NewYMWD(10, 0, 0, 0), New(1981, time.January, 1)},
		{New(1972, time.January, 1), period.NewYMWD(0, 10, 0, 0), New(1972, time.November, 1)},
		{New(1972, time.January, 1), period.NewYMWD(0, 24, 0, 0), New(1974, time.January, 1)},
		{New(1973, time.January, 1), period.NewYMWD(0, 0, 1, 0), New(1973, time.January, 8)},
		{New(1973, time.January, 1), period.NewYMWD(0, 0, 0, 10), New(1973, time.January, 11)},
		{New(1973, time.January, 1), period.NewYMWD(0, 0, 0, 365), New(1974, time.January, 1)},
		{New(1974, time.January, 1), period.NewHMS(1, 2, 3), New(1974, time.January, 1)},
		// note: the period is not normalised so the HMS is ignored even though it's more than one day
		{New(1975, time.January, 1), period.NewHMS(24, 2, 3), New(1975, time.January, 2)},
	}
	for i, c := range cases {
		out := c.in.AddPeriod(c.delta)
		if out != c.expected {
			t.Errorf("%d: %v.AddPeriod(%v) == %v, want %v", i, c.in, c.delta, out, c.expected)
		}
	}
}

// See main testin in period_test.go
func TestIsLeap(t *testing.T) {
	cases := []struct {
		year     int
		expected bool
	}{
		{2000, true},
		{2001, false},
	}
	for _, c := range cases {
		got := IsLeap(c.year)
		if got != c.expected {
			t.Errorf("TestIsLeap(%d) == %v, want %v", c.year, got, c.expected)
		}
	}
}

func TestDaysIn(t *testing.T) {
	cases := []struct {
		year     int
		month    time.Month
		expected int
	}{
		{2000, time.January, 31},
		{2000, time.February, 29},
		{2001, time.February, 28},
		{2001, time.April, 30},
	}
	for _, c := range cases {
		got1 := DaysIn(c.year, c.month)
		if got1 != c.expected {
			t.Errorf("DaysIn(%d, %d) == %v, want %v", c.year, c.month, got1, c.expected)
		}
		d := New(c.year, c.month, 1)
		got2 := d.LastDayOfMonth()
		if got2 != c.expected {
			t.Errorf("DaysIn(%d) == %v, want %v", c.year, got2, c.expected)
		}
	}
}

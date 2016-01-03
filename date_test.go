// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package date

import (
	"runtime/debug"
	"testing"
	"time"
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

func TestNew(t *testing.T) {
	cases := []string{
		"0000-01-01T00:00:00+00:00",
		"0001-01-01T00:00:00+00:00",
		"1614-01-01T01:02:03+04:00",
		"1970-01-01T00:00:00+00:00",
		"1815-12-10T05:06:07+00:00",
		"1901-09-10T00:00:00-05:00",
		"1998-09-01T00:00:00-08:00",
		"2000-01-01T00:00:00+00:00",
		"9999-12-31T00:00:00+00:00",
	}
	for _, c := range cases {
		tIn, err := time.Parse(time.RFC3339, c)
		if err != nil {
			t.Errorf("New(%v) cannot parse input: %v", c, err)
			continue
		}
		dOut := New(tIn.Year(), tIn.Month(), tIn.Day())
		if !same(dOut, tIn) {
			t.Errorf("New(%v) == %v, want date of %v", c, dOut, tIn)
		}
		dOut = NewAt(tIn)
		if !same(dOut, tIn) {
			t.Errorf("NewAt(%v) == %v, want date of %v", c, dOut, tIn)
		}
	}
}

func TestDaysSinceEpoch(t *testing.T) {
	zero := Date{}.DaysSinceEpoch()
	if zero != 0 {
		t.Errorf("Non zero %v", zero)
	}
	today := Today()
	days := today.DaysSinceEpoch()
	copy := NewOfDays(days)
	if today != copy || days == 0 {
		t.Errorf("Today == %v, want date of %v", today, copy)
	}
}

func TestToday(t *testing.T) {
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

func TestTime(t *testing.T) {
	cases := []struct {
		year  int
		month time.Month
		day   int
	}{
		{-1234, time.February, 5},
		{0, time.April, 12},
		{1, time.January, 1},
		{1946, time.February, 4},
		{1970, time.January, 1},
		{1976, time.April, 1},
		{1999, time.December, 1},
		{1111111, time.June, 21},
	}
	zones := []int{-12, -10, -5, -4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 8, 12}
	for _, c := range cases {
		d := New(c.year, c.month, c.day)
		tUTC := d.UTC()
		if !same(d, tUTC) {
			t.Errorf("TimeUTC(%v) == %v, want date part %v", d, tUTC, d)
		}
		if tUTC.Location() != time.UTC {
			t.Errorf("TimeUTC(%v) == %v, want %v", d, tUTC.Location(), time.UTC)
		}
		tLocal := d.Local()
		if !same(d, tLocal) {
			t.Errorf("TimeLocal(%v) == %v, want date part %v", d, tLocal, d)
		}
		if tLocal.Location() != time.Local {
			t.Errorf("TimeLocal(%v) == %v, want %v", d, tLocal.Location(), time.Local)
		}
		for _, z := range zones {
			location := time.FixedZone("zone", z*60*60)
			tInLoc := d.In(location)
			if !same(d, tInLoc) {
				t.Errorf("TimeIn(%v) == %v, want date part %v", d, tInLoc, d)
			}
			if tInLoc.Location() != location {
				t.Errorf("TimeIn(%v) == %v, want %v", d, tInLoc.Location(), location)
			}
		}
	}
}

func TestPredicates(t *testing.T) {
	// The list of case dates must be sorted in ascending order
	cases := []struct {
		year  int
		month time.Month
		day   int
	}{
		{-1234, time.February, 5},
		{0, time.April, 12},
		{1, time.January, 1},
		{1946, time.February, 4},
		{1970, time.January, 1},
		{1976, time.April, 1},
		{1999, time.December, 1},
		{1111111, time.June, 21},
	}
	for i, ci := range cases {
		di := New(ci.year, ci.month, ci.day)
		for j, cj := range cases {
			dj := New(cj.year, cj.month, cj.day)
			testPredicate(t, di, dj, di.Equal(dj), i == j, "Equal")
			testPredicate(t, di, dj, di.Before(dj), i < j, "Before")
			testPredicate(t, di, dj, di.After(dj), i > j, "After")
			testPredicate(t, di, dj, di == dj, i == j, "==")
			testPredicate(t, di, dj, di != dj, i != j, "!=")
		}
	}

	// Test IsZero
	zero := Date{}
	if !zero.IsZero() {
		t.Errorf("IsZero(%v) == false, want true", zero)
	}
	today := Today()
	if today.IsZero() {
		t.Errorf("IsZero(%v) == true, want false", today)
	}
}

func testPredicate(t *testing.T, di, dj Date, p, q bool, m string) {
	if p != q {
		t.Errorf("%s(%v, %v) == %v, want %v\n%v", m, di, dj, p, q, debug.Stack())
	}
}

func TestArithmetic(t *testing.T) {
	cases := []struct {
		year  int
		month time.Month
		day   int
	}{
		{-1234, time.February, 5},
		{0, time.April, 12},
		{1, time.January, 1},
		{1946, time.February, 4},
		{1970, time.January, 1},
		{1976, time.April, 1},
		{1999, time.December, 1},
		{1111111, time.June, 21},
	}
	offsets := []PeriodOfDays{-1000000, -9999, -555, -99, -22, -1, 0, 1, 22, 99, 555, 9999, 1000000}
	for _, c := range cases {
		di := New(c.year, c.month, c.day)
		for _, days := range offsets {
			dj := di.Add(days)
			days2 := dj.Sub(di)
			if days2 != days {
				t.Errorf("AddSub(%v,%v) == %v, want %v", di, days, days2, days)
			}
			d3 := dj.Add(-days)
			if d3 != di {
				t.Errorf("AddNeg(%v,%v) == %v, want %v", di, days, d3, di)
			}
			eMin1 := min(di.day, dj.day)
			aMin1 := di.Min(dj)
			if aMin1.day != eMin1 {
				t.Errorf("%v.Max(%v) is %s", di, dj, aMin1)
			}
			eMax1 := max(di.day, dj.day)
			aMax1 := di.Max(dj)
			if aMax1.day != eMax1 {
				t.Errorf("%v.Max(%v) is %s", di, dj, aMax1)
			}
		}
	}
}

func min(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func max(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func TestIsLeap(t *testing.T) {
	cases := []struct {
		year     int
		expected bool
	}{
		{2000, true},
		{2400, true},
		{2001, false},
		{2002, false},
		{2003, false},
		{2003, false},
		{2004, true},
		{2005, false},
		{1800, false},
		{1900, false},
		{2200, false},
		{2300, false},
		{2500, false},
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
			t.Errorf("DaysIn(%d) == %v, want %v", c.year, got1, c.expected)
		}
		d := New(c.year, c.month, 1)
		got2 := d.LastDayOfMonth()
		if got2 != c.expected {
			t.Errorf("DaysIn(%d) == %v, want %v", c.year, got2, c.expected)
		}
	}
}

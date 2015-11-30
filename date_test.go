// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package date

import (
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
		location := time.FixedZone("zone", c * 60 * 60)
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
			location := time.FixedZone("zone", z * 60 * 60)
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
			p := di.Equal(dj)
			q := i == j
			if p != q {
				t.Errorf("Equal(%v, %v) == %v, want %v", di, dj, p, q)
			}
			p = di.Before(dj)
			q = i < j
			if p != q {
				t.Errorf("Before(%v, %v) == %v, want %v", di, dj, p, q)
			}
			p = di.After(dj)
			q = i > j
			if p != q {
				t.Errorf("After(%v, %v) == %v, want %v", di, dj, p, q)
			}
			p = di == dj
			q = i == j
			if p != q {
				t.Errorf("Equal(%v, %v) == %v, want %v", di, dj, p, q)
			}
			p = di != dj
			q = i != j
			if p != q {
				t.Errorf("Equal(%v, %v) == %v, want %v", di, dj, p, q)
			}
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
		d := New(c.year, c.month, c.day)
		for _, days := range offsets {
			d2 := d.Add(days)
			days2 := d2.Sub(d)
			if days2 != days {
				t.Errorf("AddSub(%v,%v) == %v, want %v", d, days, days2, days)
			}
			d3 := d2.Add(-days)
			if d3 != d {
				t.Errorf("AddNeg(%v,%v) == %v, want %v", d, days, d3, d)
			}
		}
	}
}

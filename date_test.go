// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package date_test

import (
	"testing"
	"time"

	"github.com/fxtlabs/date"
)

func same(d date.Date, t time.Time) bool {
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
			t.Errorf("New(%q) cannot parse input: %q", c, err)
			continue
		}
		dOut := date.New(tIn.Year(), tIn.Month(), tIn.Day())
		if !same(dOut, tIn) {
			t.Errorf("New(%q) == %q, want date of %q", c, dOut, tIn)
		}
		dOut = date.NewAt(tIn)
		if !same(dOut, tIn) {
			t.Errorf("NewAt(%q) == %q, want date of %q", c, dOut, tIn)
		}
	}
}

func TestToday(t *testing.T) {
	today := date.Today()
	now := time.Now()
	if !same(today, now) {
		t.Errorf("Today == %q, want date of %q", today, now)
	}
	today = date.TodayUTC()
	now = time.Now().UTC()
	if !same(today, now) {
		t.Errorf("TodayUTC == %q, want date of %q", today, now)
	}
	cases := []int{-10, -5, -3, 0, 1, 4, 8, 12}
	for _, c := range cases {
		location := time.FixedZone("zone", c*60*60)
		today = date.TodayIn(location)
		now = time.Now().In(location)
		if !same(today, now) {
			t.Errorf("TodayIn(%q) == %q, want date of %q", c, today, now)
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
		di := date.New(ci.year, ci.month, ci.day)
		for j, cj := range cases {
			dj := date.New(cj.year, cj.month, cj.day)
			p := di.Equal(dj)
			q := i == j
			if p != q {
				t.Errorf("Equal(%q, %q) == %q, want %q", di, dj, p, q)
			}
			p = di.Before(dj)
			q = i < j
			if p != q {
				t.Errorf("Before(%q, %q) == %q, want %q", di, dj, p, q)
			}
			p = di.After(dj)
			q = i > j
			if p != q {
				t.Errorf("After(%q, %q) == %q, want %q", di, dj, p, q)
			}
			p = di == dj
			q = i == j
			if p != q {
				t.Errorf("Equal(%q, %q) == %q, want %q", di, dj, p, q)
			}
			p = di != dj
			q = i != j
			if p != q {
				t.Errorf("Equal(%q, %q) == %q, want %q", di, dj, p, q)
			}
		}
	}

	// Test IsZero
	zero := date.Date{}
	if !zero.IsZero() {
		t.Errorf("IsZero(%q) == false, want true", zero)
	}
	today := date.Today()
	if today.IsZero() {
		t.Errorf("IsZero(%q) == true, want false", today)
	}
}

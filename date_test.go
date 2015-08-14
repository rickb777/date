// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package date_test

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
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
		d := date.New(c.year, c.month, c.day)
		tUTC := d.UTC()
		if !same(d, tUTC) {
			t.Errorf("TimeUTC(%q) == %q, want date part %q", d, tUTC, d)
		}
		if tUTC.Location() != time.UTC {
			t.Errorf("TimeUTC(%q) == %q, want %q", d, tUTC.Location(), time.UTC)
		}
		tLocal := d.Local()
		if !same(d, tLocal) {
			t.Errorf("TimeLocal(%q) == %q, want date part %q", d, tLocal, d)
		}
		if tLocal.Location() != time.Local {
			t.Errorf("TimeLocal(%q) == %q, want %q", d, tLocal.Location(), time.Local)
		}
		for _, z := range zones {
			location := time.FixedZone("zone", z*60*60)
			tInLoc := d.In(location)
			if !same(d, tInLoc) {
				t.Errorf("TimeIn(%q) == %q, want date part %q", d, tInLoc, d)
			}
			if tInLoc.Location() != location {
				t.Errorf("TimeIn(%q) == %q, want %q", d, tInLoc.Location(), location)
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
	offsets := []int{-1000000, -9999, -555, -99, -22, -1, 0, 1, 22, 99, 555, 9999, 1000000}
	for _, c := range cases {
		d := date.New(c.year, c.month, c.day)
		for _, days := range offsets {
			d2 := d.Add(days)
			days2 := d2.Sub(d)
			if days2 != days {
				t.Errorf("AddSub(%q,%q) == %q, want %q", d, days, days2, days)
			}
			d3 := d2.Add(-days)
			if d3 != d {
				t.Errorf("AddNeg(%q,%q) == %q, want %q", d, days, d3, d)
			}
		}
	}
}

func TestGobEncoding(t *testing.T) {
	var b bytes.Buffer
	encoder := gob.NewEncoder(&b)
	decoder := gob.NewDecoder(&b)
	cases := []date.Date{
		date.New(-11111, time.February, 3),
		date.New(-1, time.December, 31),
		date.New(0, time.January, 1),
		date.New(1, time.January, 1),
		date.New(1970, time.January, 1),
		date.New(2012, time.June, 25),
		date.New(12345, time.June, 7),
	}
	for _, c := range cases {
		var d date.Date
		err := encoder.Encode(&c)
		if err != nil {
			t.Errorf("Gob(%q) encode error %q", c, err)
		} else {
			err = decoder.Decode(&d)
			if err != nil {
				t.Errorf("Gob(%q) decode error %q", c, err)
			}
		}
	}
}

func TestInvalidGob(t *testing.T) {
	cases := []struct {
		bytes []byte
		want  string
	}{
		{[]byte{}, "Date.UnmarshalBinary: no data"},
		{[]byte{1, 2, 3}, "Date.UnmarshalBinary: invalid length"},
	}
	for _, c := range cases {
		var ignored date.Date
		err := ignored.GobDecode(c.bytes)
		if err == nil || err.Error() != c.want {
			t.Errorf("InvalidGobDecode(%q) == %q, want %q", c.bytes, err, c.want)
		}
		err = ignored.UnmarshalBinary(c.bytes)
		if err == nil || err.Error() != c.want {
			t.Errorf("InvalidUnmarshalBinary(%q) == %q, want %q", c.bytes, err, c.want)
		}
	}
}

func TestJSONMarshalling(t *testing.T) {
	var d date.Date
	cases := []struct {
		value date.Date
		want  string
	}{
		{date.New(-11111, time.February, 3), `"-11111-02-03"`},
		{date.New(-1, time.December, 31), `"-0001-12-31"`},
		{date.New(0, time.January, 1), `"0000-01-01"`},
		{date.New(1, time.January, 1), `"0001-01-01"`},
		{date.New(1970, time.January, 1), `"1970-01-01"`},
		{date.New(2012, time.June, 25), `"2012-06-25"`},
		{date.New(12345, time.June, 7), `"+12345-06-07"`},
	}
	for _, c := range cases {
		bytes, err := json.Marshal(c.value)
		if err != nil {
			t.Errorf("JSON(%q) marshal error %q", c, err)
		} else if string(bytes) != c.want {
			t.Errorf("%v JSON(%q) == %q, want %q", c.value, string(bytes), c.want)
		} else {
			err = json.Unmarshal(bytes, &d)
			if err != nil {
				t.Errorf("JSON(%q) unmarshal error %q", c.value, err)
			}
		}
	}
}

func TestInvalidJSON(t *testing.T) {
	cases := []struct {
		value string
		want  string
	}{
		{`"not-a-date"`, `Date.ParseISO: cannot parse not-a-date`},
		{`2015-08-15"`, `Date.UnmarshalJSON: missing double quotes (2015-08-15")`},
		{`"2015-08-15`, `Date.UnmarshalJSON: missing double quotes ("2015-08-15)`},
		{`"215-08-15"`, `Date.ParseISO: cannot parse 215-08-15`},
	}
	for _, c := range cases {
		var d date.Date
		err := d.UnmarshalJSON([]byte(c.value))
		if err == nil || err.Error() != c.want {
			t.Errorf("InvalidJSON(%q) == %q, want %q", c.value, err, c.want)
		}
	}
}

func TestTextMarshalling(t *testing.T) {
	var d date.Date
	cases := []struct {
		value date.Date
		want  string
	}{
		{date.New(-11111, time.February, 3), "-11111-02-03"},
		{date.New(-1, time.December, 31), "-0001-12-31"},
		{date.New(0, time.January, 1), "0000-01-01"},
		{date.New(1, time.January, 1), "0001-01-01"},
		{date.New(1970, time.January, 1), "1970-01-01"},
		{date.New(2012, time.June, 25), "2012-06-25"},
		{date.New(12345, time.June, 7), "+12345-06-07"},
	}
	for _, c := range cases {
		bytes, err := c.value.MarshalText()
		if err != nil {
			t.Errorf("Text(%q) marshal error %q", c, err)
		} else if string(bytes) != c.want {
			t.Errorf("%v Text(%q) == %q, want %q", c.value, string(bytes), c.want)
		} else {
			err = d.UnmarshalText(bytes)
			if err != nil {
				t.Errorf("Text(%q) unmarshal error %q", c.value, err)
			}
		}
	}
}

func TestInvalidText(t *testing.T) {
	cases := []struct {
		value string
		want  string
	}{
		{`not-a-date`, `Date.ParseISO: cannot parse not-a-date`},
		{`215-08-15`, `Date.ParseISO: cannot parse 215-08-15`},
	}
	for _, c := range cases {
		var d date.Date
		err := d.UnmarshalText([]byte(c.value))
		if err == nil || err.Error() != c.want {
			t.Errorf("InvalidText(%q) == %q, want %q", c.value, err, c.want)
		}
	}
}

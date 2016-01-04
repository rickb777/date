// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package date

import (
	"math/rand"
	"testing"
	"time"
)

func TestEncode(t *testing.T) {
	cases := []int{
		0, 1, 28, 30, 31, 32, 364, 365, 366, 367, 500, 1000, 10000, 100000,
	}
	tBase := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)
	for i, c := range cases {
		d := encode(tBase.AddDate(0, 0, c))
		if d != PeriodOfDays(c) {
			t.Errorf("Encode(%v) == %v, want %v", i, d, c)
		}
		d = encode(tBase.AddDate(0, 0, -c))
		if d != PeriodOfDays(-c) {
			t.Errorf("Encode(%v) == %v, want %v", i, d, c)
		}
	}
}

func TestEncodeDecode(t *testing.T) {
	cases := []struct {
		year  int
		month time.Month
		day   int
	}{
		{1969, time.December, 31},
		{1970, time.January, 1},
		{1970, time.January, 2},
		{2000, time.February, 28},
		{2000, time.February, 29},
		{2000, time.March, 1},
		{2004, time.February, 28},
		{2004, time.February, 29},
		{2004, time.March, 1},
		{2100, time.February, 28},
		{2100, time.February, 29},
		{2100, time.March, 1},
		{0, time.January, 1},
		{1, time.February, 3},
		{19, time.March, 4},
		{100, time.April, 5},
		{2000, time.May, 6},
		{30000, time.June, 7},
		{400000, time.July, 8},
		{5000000, time.August, 9},
		{-1, time.September, 11},
		{-19, time.October, 12},
		{-100, time.November, 13},
		{-2000, time.December, 14},
		{-30000, time.February, 15},
		{-400000, time.May, 16},
		{-5000000, time.September, 17},
	}
	for _, c := range cases {
		tIn := time.Date(c.year, c.month, c.day, 0, 0, 0, 0, time.UTC)
		d := encode(tIn)
		tOut := decode(d)
		if !tIn.Equal(tOut) {
			t.Errorf("EncodeDecode(%v) == %v, want %v", c, tOut, tIn)
		}
	}
}

func TestDecodeEncode(t *testing.T) {
	for i := 0; i < 1000; i++ {
		c := PeriodOfDays(rand.Int31())
		d := encode(decode(c))
		if d != c {
			t.Errorf("DecodeEncode(%v) == %v, want %v", i, d, c)
		}
	}
	for i := 0; i < 1000; i++ {
		c := -PeriodOfDays(rand.Int31())
		d := encode(decode(c))
		if d != c {
			t.Errorf("DecodeEncode(%v) == %v, want %v", i, d, c)
		}
	}
}

// TestZone checks that the conversions between a time.Time value and the
// internal representation of a Date value correctly handle time zones other
// than UTC, especially in cases where the local date at a given time is
// different from the UTC date for that same time.
func TestZone(t *testing.T) {
	cases := []string{
		"2015-07-29 15:12:34 +0000",
		"2015-07-29 15:12:34 -0500",
		"2015-07-29 15:12:34 +0500",
		"2015-07-29 21:12:34 -0500",
		"2015-07-29 21:12:34 -0500",
		"2015-07-29 03:12:34 +0500",
		"2015-07-29 03:12:34 +0500",
	}
	for _, c := range cases {
		tIn, err := time.Parse("2006-01-02 15:04:05 -0700", c)
		if err != nil {
			t.Errorf("Zone(%v) cannot parse %v", c, c)
		}
		d := encode(tIn)
		tOut := decode(d)
		yIn, mIn, dIn := tIn.Date()
		yOut, mOut, dOut := tOut.Date()
		if yIn != yOut {
			t.Errorf("Zone(%v).y == %v, want %v", c, yOut, yIn)
		}
		if mIn != mOut {
			t.Errorf("Zone(%v).m == %v, want %v", c, mOut, mIn)
		}
		if dIn != dOut {
			t.Errorf("Zone(%v).d == %v, want %v", c, dOut, dIn)
		}
	}
}

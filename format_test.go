// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package date

import (
	"testing"
)

func TestDate_String(t *testing.T) {
	cases := []struct {
		value string
	}{
		{"-0001-01-01"},
		{"0000-01-01"},
		{"1000-01-01"},
		{"1970-01-01"},
		{"2000-11-22"},
		{"+10000-01-01"},
	}
	for _, c := range cases {
		d := MustParseISO(c.value)
		value := d.String()
		if value != c.value {
			t.Errorf("String() == %v, want %v", value, c.value)
		}
	}
}

func TestDate_FormatISO(t *testing.T) {
	cases := []struct {
		value string
		n     int
	}{
		{"-5000-02-03", 4},
		{"-05000-02-03", 5},
		{"-005000-02-03", 6},
		{"+0000-01-01", 4},
		{"+00000-01-01", 5},
		{"+1000-01-01", 4},
		{"+01000-01-01", 5},
		{"+1970-01-01", 4},
		{"+001999-12-31", 6},
		{"+999999-12-31", 6},
	}
	for _, c := range cases {
		d := MustParseISO(c.value)
		value := d.FormatISO(c.n)
		if value != c.value {
			t.Errorf("FormatISO(%v) == %v, want %v", c, value, c.value)
		}
	}
}

func TestDate_Format(t *testing.T) {
	cases := []struct {
		value    string
		format   string
		expected string
	}{
		{"1970-01-01", "2 Jan 2006", "1 Jan 1970"},
		{"1970-01-01", "Jan 02 2006", "Jan 01 1970"},
		{"1970-01-01", "Jan 2nd 2006", "Jan 1st 1970"},
		{"2016-01-01", "2nd Jan 2006", "1st Jan 2016"},
		{"2016-02-02", "Jan 2nd 2006", "Feb 2nd 2016"},
		{"2016-03-03", "Jan 2nd 2006", "Mar 3rd 2016"},
		{"2016-04-04", "2nd Jan 2006", "4th Apr 2016"},
		{"2016-05-20", "Jan 2nd 2006", "May 20th 2016"},
		{"2016-06-21", "Jan 2nd 2006", "Jun 21st 2016"},
		{"2016-07-22", "Jan 2nd 2006", "Jul 22nd 2016"},
		{"2016-08-23", "Jan 2nd 2006", "Aug 23rd 2016"},
		{"2016-09-30", "Jan 2nd 2006", "Sep 30th 2016"},
		{"2016-10-31", "Jan 2nd 2006", "Oct 31st 2016"},
		{"2016-01-07", "Monday January 2nd 2006", "Thursday January 7th 2016"},
		{"2016-01-07", "Monday 2nd Monday 2nd", "Thursday 7th Thursday 7th"},
		{"2016-11-01", "2nd 2nd 2nd", "1st 1st 1st"},
	}
	for _, c := range cases {
		d := MustParseISO(c.value)
		actual := d.Format(c.format)
		if actual != c.expected {
			t.Errorf("Format(%v) == %v, want %v", c, actual, c.expected)
		}
	}
}

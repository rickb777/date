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
		{value: "-0001-01-01"},
		{value: "0000-01-01"},
		{value: "0001-01-01"},
		{value: "1000-01-01"},
		{value: "1970-01-01"},
		{value: "2000-11-22"},
		{value: "+10000-01-01"},
	}
	for _, c := range cases {
		d := MustParseISO(c.value)
		value := d.String()
		if value != c.value {
			t.Errorf("String() == %v, want %v", value, c.value)
		}
	}
}

func TestDate_FormatOrdinal(t *testing.T) {
	cases := []struct {
		value, expected string
	}{
		{value: "-5000-001", expected: "-5000-001"},
		{value: "-05000-001", expected: "-5000-001"},
		{value: "-005000-001", expected: "-5000-001"},
		{value: "0001-001", expected: "0001-001"},
		{value: "00000-001", expected: "0000-001"},
		{value: "1000-001", expected: "1000-001"},
		{value: "01000-001", expected: "1000-001"},
		{value: "1970-001", expected: "1970-001"},
		{value: "001999-365", expected: "1999-365"},
		{value: "999999-365", expected: "999999-365"},
	}
	for i, c := range cases {
		d := MustParseISO(c.value)
		value := d.FormatOrdinal()
		if value != c.expected {
			t.Errorf("%d: FormatOrdinal(%v) == %v, want %v", i, c, value, c.value)
		}
	}
}

func TestDate_FormatISO(t *testing.T) {
	cases := []struct {
		value string
		n     int
	}{
		{value: "-5000-02-03", n: 4},
		{value: "-05000-02-03", n: 5},
		{value: "-005000-02-03", n: 6},
		{value: "+0000-01-01", n: 4},
		{value: "+00000-01-01", n: 5},
		{value: "+1000-01-01", n: 4},
		{value: "+01000-01-01", n: 5},
		{value: "+1970-01-01", n: 4},
		{value: "+001999-12-31", n: 6},
		{value: "+999999-12-31", n: 6},
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
		{value: "1970-01-01", format: "2 Jan 2006", expected: "1 Jan 1970"},
		{value: "1970-01-01", format: "Jan 02 2006", expected: "Jan 01 1970"},
		{value: "1970-01-01", format: "Jan 2nd 2006", expected: "Jan 1st 1970"},
		{value: "2016-01-01", format: "2nd Jan 2006", expected: "1st Jan 2016"},
		{value: "2016-02-02", format: "Jan 2nd 2006", expected: "Feb 2nd 2016"},
		{value: "2016-03-03", format: "Jan 2nd 2006", expected: "Mar 3rd 2016"},
		{value: "2016-04-04", format: "2nd Jan 2006", expected: "4th Apr 2016"},
		{value: "2016-05-20", format: "Jan 2nd 2006", expected: "May 20th 2016"},
		{value: "2016-06-21", format: "Jan 2nd 2006", expected: "Jun 21st 2016"},
		{value: "2016-07-22", format: "Jan 2nd 2006", expected: "Jul 22nd 2016"},
		{value: "2016-08-23", format: "Jan 2nd 2006", expected: "Aug 23rd 2016"},
		{value: "2016-09-30", format: "Jan 2nd 2006", expected: "Sep 30th 2016"},
		{value: "2016-10-31", format: "Jan 2nd 2006", expected: "Oct 31st 2016"},
		{value: "2016-01-07", format: "Monday January 2nd 2006", expected: "Thursday January 7th 2016"},
		{value: "2016-01-07", format: "Monday 2nd Monday 2nd", expected: "Thursday 7th Thursday 7th"},
		{value: "2016-11-01", format: "2nd 2nd 2nd", expected: "1st 1st 1st"},
	}
	for _, c := range cases {
		d := MustParseISO(c.value)
		actual := d.Format(c.format)
		if actual != c.expected {
			t.Errorf("Format(%v) == %v, want %v", c, actual, c.expected)
		}
	}
}

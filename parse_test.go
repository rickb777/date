// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package date

import (
	"testing"
	time "time"
)

func TestAutoParse(t *testing.T) {
	cases := []struct {
		value string
		year  int
		month time.Month
		day   int
	}{
		{value: "01-01-1970", year: 1970, month: time.January, day: 1},
		{value: "+1970-01-01", year: 1970, month: time.January, day: 1},
		{value: "+01970-01-02", year: 1970, month: time.January, day: 2},
		{value: " 31/12/1969 ", year: 1969, month: time.December, day: 31},
		{value: "1969/12/31", year: 1969, month: time.December, day: 31},
		{value: "1969.12.31", year: 1969, month: time.December, day: 31},
		{value: "1969-12-31", year: 1969, month: time.December, day: 31},
		{value: "2000-02-28", year: 2000, month: time.February, day: 28},
		{value: "+2000-02-29", year: 2000, month: time.February, day: 29},
		{value: "+02000-03-01", year: 2000, month: time.March, day: 1},
		{value: "+002004-02-28", year: 2004, month: time.February, day: 28},
		{value: "2004-02-29", year: 2004, month: time.February, day: 29},
		{value: "2004-03-01", year: 2004, month: time.March, day: 1},
		{value: "0000-01-01", month: time.January, day: 1},
		{value: "+0001-02-03", year: 1, month: time.February, day: 3},
		{value: " +00019-03-04 ", year: 19, month: time.March, day: 4},
		{value: "0100-04-05", year: 100, month: time.April, day: 5},
		{value: "2000-05-06", year: 2000, month: time.May, day: 6},
		{value: "+5000000-08-09", year: 5000000, month: time.August, day: 9},
		{value: "-0001-09-11", year: -1, month: time.September, day: 11},
		{value: " -0019-10-12 ", year: -19, month: time.October, day: 12},
		{value: "-00100-11-13", year: -100, month: time.November, day: 13},
		{value: "-02000-12-14", year: -2000, month: time.December, day: 14},
		{value: "-30000-02-15", year: -30000, month: time.February, day: 15},
		{value: "-0400000-05-16", year: -400000, month: time.May, day: 16},
		{value: "-5000000-09-17", year: -5000000, month: time.September, day: 17},
		{value: "12340506", year: 1234, month: time.May, day: 6},
		{value: "+12340506", year: 1234, month: time.May, day: 6},
		{value: "-00191012", year: -19, month: time.October, day: 12},
		{value: " -00191012 ", year: -19, month: time.October, day: 12},
	}
	for _, c := range cases {
		d := MustAutoParse(c.value)
		year, month, day := d.Date()
		if year != c.year || month != c.month || day != c.day {
			t.Errorf("ParseISO(%v) == %v, want (%v, %v, %v)", c.value, d, c.year, c.month, c.day)
		}
	}

	badCases := []string{
		"1234-05",
		"1234-5-6",
		"1234-05-6",
		"1234-5-06",
		"1234-0A-06",
		"1234-05-0B",
		"1234-05-06trailing",
		"padding1234-05-06",
		"1-02-03",
		"10-11-12",
		"100-02-03",
		"+1-02-03",
		"+10-11-12",
		"+100-02-03",
		"-123-05-06",
		"--",
		"",
		"  ",
	}
	for _, c := range badCases {
		d, err := AutoParse(c)
		if err == nil {
			t.Errorf("ParseISO(%v) == %v", c, d)
		}
	}
}

func TestParseISO(t *testing.T) {
	cases := []struct {
		value string
		year  int
		month time.Month
		day   int
	}{
		{"1969-12-31", 1969, time.December, 31},
		{"+1970-01-01", 1970, time.January, 1},
		{"+01970-01-02", 1970, time.January, 2},
		{"2000-02-28", 2000, time.February, 28},
		{"+2000-02-29", 2000, time.February, 29},
		{"+02000-03-01", 2000, time.March, 1},
		{"+002004-02-28", 2004, time.February, 28},
		{"2004-02-29", 2004, time.February, 29},
		{"2004-03-01", 2004, time.March, 1},
		{"0000-01-01", 0, time.January, 1},
		{"+0001-02-03", 1, time.February, 3},
		{"+00019-03-04", 19, time.March, 4},
		{"0100-04-05", 100, time.April, 5},
		{"2000-05-06", 2000, time.May, 6},
		{"+30000-06-07", 30000, time.June, 7},
		{"+400000-07-08", 400000, time.July, 8},
		{"+5000000-08-09", 5000000, time.August, 9},
		{"-0001-09-11", -1, time.September, 11},
		{"-0019-10-12", -19, time.October, 12},
		{"-00100-11-13", -100, time.November, 13},
		{"-02000-12-14", -2000, time.December, 14},
		{"-30000-02-15", -30000, time.February, 15},
		{"-0400000-05-16", -400000, time.May, 16},
		{"-5000000-09-17", -5000000, time.September, 17},
		{"12340506", 1234, time.May, 6},
		{"+12340506", 1234, time.May, 6},
		{"-00191012", -19, time.October, 12},
	}
	for _, c := range cases {
		d := MustParseISO(c.value)
		year, month, day := d.Date()
		if year != c.year || month != c.month || day != c.day {
			t.Errorf("ParseISO(%v) == %v, want (%v, %v, %v)", c.value, d, c.year, c.month, c.day)
		}
	}

	badCases := []string{
		"1234-05",
		"1234-5-6",
		"1234-05-6",
		"1234-5-06",
		"1234/05/06",
		"1234-0A-06",
		"1234-05-0B",
		"1234-05-06trailing",
		"padding1234-05-06",
		"1-02-03",
		"10-11-12",
		"100-02-03",
		"+1-02-03",
		"+10-11-12",
		"+100-02-03",
		"-123-05-06",
	}
	for _, c := range badCases {
		d, err := ParseISO(c)
		if err == nil {
			t.Errorf("ParseISO(%v) == %v", c, d)
		}
	}
}

func BenchmarkParseISO(b *testing.B) {
	cases := []struct {
		layout string
		value  string
		year   int
		month  time.Month
		day    int
	}{
		{ISO8601, "1969-12-31", 1969, time.December, 31},
		{ISO8601, "2000-02-28", 2000, time.February, 28},
		{ISO8601, "2004-02-29", 2004, time.February, 29},
		{ISO8601, "2004-03-01", 2004, time.March, 1},
		{ISO8601, "0000-01-01", 0, time.January, 1},
		{ISO8601, "0001-02-03", 1, time.February, 3},
		{ISO8601, "0100-04-05", 100, time.April, 5},
		{ISO8601, "2000-05-06", 2000, time.May, 6},
	}
	for n := 0; n < b.N; n++ {
		c := cases[n%len(cases)]
		_, err := ParseISO(c.value)
		if err != nil {
			b.Errorf("ParseISO(%v) == %v", c.value, err)
		}
	}
}

func TestParse(t *testing.T) {
	// Test ability to parse a few common date formats
	cases := []struct {
		layout string
		value  string
		year   int
		month  time.Month
		day    int
	}{
		{layout: ISO8601, value: "1969-12-31", year: 1969, month: time.December, day: 31},
		{layout: ISO8601B, value: "19700101", year: 1970, month: time.January, day: 1},
		{layout: RFC822, value: "29-Feb-00", year: 2000, month: time.February, day: 29},
		{layout: RFC822W, value: "Mon, 01-Mar-04", year: 2004, month: time.March, day: 1},
		{layout: RFC850, value: "Wednesday, 12-Aug-15", year: 2015, month: time.August, day: 12},
		{layout: RFC1123, value: "05 Dec 1928", year: 1928, month: time.December, day: 5},
		{layout: RFC1123W, value: "Mon, 05 Dec 1928", year: 1928, month: time.December, day: 5},
		{layout: RFC3339, value: "2345-06-07", year: 2345, month: time.June, day: 7},
		{layout: time.RFC3339Nano, value: "2020-04-01T12:11:10.101+09:00", year: 2020, month: time.April, day: 1},
		{layout: "20060102", value: "20190619", year: 2019, month: time.June, day: 19},
	}
	for _, c := range cases {
		d := MustParse(c.layout, c.value)
		year, month, day := d.Date()
		if year != c.year || month != c.month || day != c.day {
			t.Errorf("Parse(%v) == %v, want (%v, %v, %v)", c.value, d, c.year, c.month, c.day)
		}
	}

	// Test inability to parse ISO 8601 expanded year format
	badCases := []string{
		"+1234-05-06",
		"+12345-06-07",
		"-1234-05-06",
		"-12345-06-07",
	}
	for _, c := range badCases {
		d, err := Parse(ISO8601, c)
		if err == nil {
			t.Errorf("Parse(%v) == %v", c, d)
		}
	}
}

func BenchmarkParse(b *testing.B) {
	// Test ability to parse a few common date formats
	cases := []struct {
		layout string
		value  string
		year   int
		month  time.Month
		day    int
	}{
		{layout: ISO8601, value: "1969-12-31", year: 1969, month: time.December, day: 31},
		{layout: ISO8601, value: "2000-02-28", year: 2000, month: time.February, day: 28},
		{layout: ISO8601, value: "2004-02-29", year: 2004, month: time.February, day: 29},
		{layout: ISO8601, value: "2004-03-01", year: 2004, month: time.March, day: 1},
		{layout: ISO8601, value: "0000-01-01", month: time.January, day: 1},
		{layout: ISO8601, value: "0001-02-03", year: 1, month: time.February, day: 3},
		{layout: ISO8601, value: "0100-04-05", year: 100, month: time.April, day: 5},
		{layout: ISO8601, value: "2000-05-06", year: 2000, month: time.May, day: 6},
	}
	for n := 0; n < b.N; n++ {
		c := cases[n%len(cases)]
		_, err := Parse(c.layout, c.value)
		if err != nil {
			b.Errorf("Parse(%v) == %v", c.value, err)
		}
	}
}

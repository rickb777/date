// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package date

import (
	"fmt"
	"testing"
	time "time"
)

func TestAutoParse_both(t *testing.T) {
	cases := []struct {
		value string
		year  int
		month time.Month
		day   int
	}{
		{value: "01-01-1970", year: 1970, month: time.January, day: 1},
		{value: "+1970-01-01", year: 1970, month: time.January, day: 1},
		{value: "+01970-01-02", year: 1970, month: time.January, day: 2},
		{value: "1969/12/31", year: 1969, month: time.December, day: 31},
		{value: "1969.12.31", year: 1969, month: time.December, day: 31},
		{value: "1969-12-31", year: 1969, month: time.December, day: 31},
		{value: "2000-02-28", year: 2000, month: time.February, day: 28},
		{value: "+2000-02-29", year: 2000, month: time.February, day: 29},
		{value: "+02000-03-01", year: 2000, month: time.March, day: 1},
		{value: "+002004-02-28", year: 2004, month: time.February, day: 28},
		{value: "2004-02-29", year: 2004, month: time.February, day: 29},
		{value: "2004-060", year: 2004, month: time.February, day: 29},
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
		// yyyy-ooo ordinal cases
		{value: "2004-001", year: 2004, month: time.January, day: 1},
		{value: "2004-060", year: 2004, month: time.February, day: 29},
		{value: "2004-366", year: 2004, month: time.December, day: 31},
		{value: "2003-365", year: 2003, month: time.December, day: 31},
		// basic format is only supported for yyyymmdd (yyyyooo ordinal is not supported)
		{value: "12340506", year: 1234, month: time.May, day: 6},
		{value: "+12340506", year: 1234, month: time.May, day: 6},
		{value: "-00191012", year: -19, month: time.October, day: 12},
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("%d %s", i, c.value), func(t *testing.T) {
			d := MustAutoParse(c.value)
			year, month, day := d.Date()
			if year != c.year || month != c.month || day != c.day {
				t.Errorf("ParseISO(%v) == %v, want (%v, %v, %v)", c.value, d, c.year, c.month, c.day)
			}

			d = MustAutoParseUS(c.value)
			year, month, day = d.Date()
			if year != c.year || month != c.month || day != c.day {
				t.Errorf("ParseISO(%v) == %v, want (%v, %v, %v)", c.value, d, c.year, c.month, c.day)
			}
		})
	}
}

func TestAutoParse(t *testing.T) {
	cases := []struct {
		value string
		year  int
		month time.Month
		day   int
	}{
		{value: " 31/12/1969 ", year: 1969, month: time.December, day: 31},
		{value: " 5/6/1905 ", year: 1905, month: time.June, day: 5},
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("%d %s", i, c.value), func(t *testing.T) {
			d := MustAutoParse(c.value)
			year, month, day := d.Date()
			if year != c.year || month != c.month || day != c.day {
				t.Errorf("ParseISO(%v) == %v, want (%v, %v, %v)", c.value, d, c.year, c.month, c.day)
			}
		})
	}
}

func TestAutoParseUS(t *testing.T) {
	cases := []struct {
		value string
		year  int
		month time.Month
		day   int
	}{
		{value: " 12/31/1969 ", year: 1969, month: time.December, day: 31},
		{value: " 6/5/1905 ", year: 1905, month: time.June, day: 5},
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("%d %s", i, c.value), func(t *testing.T) {
			d := MustAutoParseUS(c.value)
			year, month, day := d.Date()
			if year != c.year || month != c.month || day != c.day {
				t.Errorf("ParseISO(%v) == %v, want (%v, %v, %v)", c.value, d, c.year, c.month, c.day)
			}
		})
	}
}

func TestAutoParse_errors(t *testing.T) {
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
		"20210506T0Z",
		"2021-05-06T0:0:0Z",
		"--",
		"",
		"  ",
	}
	for _, c := range badCases {
		d, err := AutoParse(c)
		if err == nil {
			t.Errorf("ParseISO(%v) == %v", c, d)
		}

		d, err = AutoParseUS(c)
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
		{value: "1969-12-31", year: 1969, month: time.December, day: 31},
		{value: "+1970-01-01", year: 1970, month: time.January, day: 1},
		{value: "+01970-01-02", year: 1970, month: time.January, day: 2},
		{value: "2000-02-28", year: 2000, month: time.February, day: 28},
		{value: "2018-02-03T00:00:00Z", year: 2018, month: time.February, day: 3},
		{value: "+2000-02-29", year: 2000, month: time.February, day: 29},
		{value: "+02000-03-01", year: 2000, month: time.March, day: 1},
		{value: "+002004-02-28", year: 2004, month: time.February, day: 28},
		{value: "2004-02-29", year: 2004, month: time.February, day: 29},
		{value: "2004-03-01", year: 2004, month: time.March, day: 1},
		{value: "0000-01-01", month: time.January, day: 1},
		{value: "+0001-02-03", year: 1, month: time.February, day: 3},
		{value: "+00019-03-04", year: 19, month: time.March, day: 4},
		{value: "0100-04-05", year: 100, month: time.April, day: 5},
		{value: "2000-05-06", year: 2000, month: time.May, day: 6},
		{value: "+30000-06-07", year: 30000, month: time.June, day: 7},
		{value: "+400000-07-08", year: 400000, month: time.July, day: 8},
		{value: "+5000000-08-09", year: 5000000, month: time.August, day: 9},
		{value: "-0001-09-11", year: -1, month: time.September, day: 11},
		{value: "-0019-10-12", year: -19, month: time.October, day: 12},
		{value: "-00100-11-13", year: -100, month: time.November, day: 13},
		{value: "-02000-12-14", year: -2000, month: time.December, day: 14},
		{value: "-30000-02-15", year: -30000, month: time.February, day: 15},
		{value: "-0400000-05-16", year: -400000, month: time.May, day: 16},
		{value: "-5000000-09-17", year: -5000000, month: time.September, day: 17},
		// yyyy-ooo ordinal cases
		{value: "2004-001", year: 2004, month: time.January, day: 1},
		{value: "2004-060", year: 2004, month: time.February, day: 29},
		{value: "2004-366", year: 2004, month: time.December, day: 31},
		{value: "2003-365", year: 2003, month: time.December, day: 31},
		// basic format is only supported for yyyymmdd (yyyyooo ordinal is not supported)
		{value: "12340506", year: 1234, month: time.May, day: 6},
		{value: "+12340506", year: 1234, month: time.May, day: 6},
		{value: "-00191012", year: -19, month: time.October, day: 12},
		{value: "20210506T010203Z", year: 2021, month: time.May, day: 6},
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("%d %s", i, c.value), func(t *testing.T) {
			d := MustParseISO(c.value)
			year, month, day := d.Date()
			if year != c.year || month != c.month || day != c.day {
				t.Errorf("ParseISO(%v) == %v, want (%v, %v, %v)", c.value, d, c.year, c.month, c.day)
			}
		})
	}
}

func TestParseISO_errors(t *testing.T) {
	cases := []struct {
		value string
		want  string
	}{
		{value: ``, want: `date.ParseISO: cannot parse "": ` + "too short"},
		{value: `-`, want: `date.ParseISO: cannot parse "-": ` + "too short"},
		{value: `z`, want: `date.ParseISO: cannot parse "z": ` + "too short"},
		{value: `z--`, want: `date.ParseISO: cannot parse "z--": ` + "year has wrong length\nmonth has wrong length\nday has wrong length"},
		{value: `not-a-date`, want: `date.ParseISO: cannot parse "not-a-date": ` + "year has wrong length\nmonth has wrong length\nday has wrong length"},
		{value: `foot-of-og`, want: `date.ParseISO: cannot parse "foot-of-og": ` + "invalid year\ninvalid month\ninvalid day"},
		{value: `215-08-15`, want: `date.ParseISO: cannot parse "215-08-15": year has wrong length`},
		{value: "1234-05", want: `date.ParseISO: cannot parse "1234-05": incorrect length for ordinal date yyyy-ooo`},
		{value: "1234-5-6", want: `date.ParseISO: cannot parse "1234-5-6": ` + "month has wrong length\nday has wrong length"},
		{value: "1234-05-6", want: `date.ParseISO: cannot parse "1234-05-6": day has wrong length`},
		{value: "1234-5-06", want: `date.ParseISO: cannot parse "1234-5-06": month has wrong length`},
		{value: "1234/05/06", want: `date.ParseISO: cannot parse "1234/05/06": ` + "invalid year\ninvalid month"},
		{value: "1234-0A-06", want: `date.ParseISO: cannot parse "1234-0A-06": invalid month`},
		{value: "1234-05-0B", want: `date.ParseISO: cannot parse "1234-05-0B": invalid day`},
		{value: "1234-05-06trailing", want: `date.ParseISO: cannot parse "1234-05-06trailing": day has wrong length`},
		{value: "padding1234-05-06", want: `date.ParseISO: cannot parse "padding1234-05-06": invalid year`},
		{value: "1-02-03", want: `date.ParseISO: cannot parse "1-02-03": year has wrong length`},
		{value: "10-11-12", want: `date.ParseISO: cannot parse "10-11-12": year has wrong length`},
		{value: "100-02-03", want: `date.ParseISO: cannot parse "100-02-03": year has wrong length`},
		{value: "+1-02-03", want: `date.ParseISO: cannot parse "+1-02-03": year has wrong length`},
		{value: "+10-11-12", want: `date.ParseISO: cannot parse "+10-11-12": year has wrong length`},
		{value: "+100-02-03", want: `date.ParseISO: cannot parse "+100-02-03": year has wrong length`},
		{value: "-123-05-06", want: `date.ParseISO: cannot parse "-123-05-06": year has wrong length`},
		{value: "2018-02-03T0:0:0Z", want: `date.ParseISO: date-time "2018-02-03T0:0:0Z": not a time`},
		{value: "2018-02-03T0Z", want: `date.ParseISO: date-time "2018-02-03T0Z": not a time`},
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("%d %s", i, c.value), func(t *testing.T) {
			d, err := ParseISO(c.value)
			if err == nil {
				t.Errorf("ParseISO(%v) == %v", c, d)
			}
			if err.Error() != c.want {
				t.Errorf("got %s\nwant %s", err.Error(), c.want)
			}
		})
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
	for i, c := range cases {
		t.Run(fmt.Sprintf("%d %s", i, c.value), func(t *testing.T) {
			d := MustParse(c.layout, c.value)
			year, month, day := d.Date()
			if year != c.year || month != c.month || day != c.day {
				t.Errorf("Parse(%v) == %v, want (%v, %v, %v)", c.value, d, c.year, c.month, c.day)
			}
		})
	}
}

func TestParse_errors(t *testing.T) {
	// Test inability to parse ISO 8601 expanded year format
	badCases := []string{
		"+1234-05-06", // plus sign is not allowed
		"+12345-06-07",
		"12345-06-07", // five digits are not allowed
		"-1234-05-06", // negative sign is not allowed
		"-12345-06-07",
	}
	for i, c := range badCases {
		t.Run(fmt.Sprintf("%d %s", i, c), func(t *testing.T) {
			d, err := Parse(ISO8601, c)
			if err == nil {
				t.Errorf("Parse(%v) == %v", c, d)
			}
		})
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

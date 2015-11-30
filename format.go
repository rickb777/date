// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package date

import (
	"fmt"
	"strconv"
	"time"
	"strings"
)

// These are predefined layouts for use in Date.Format and Date.Parse.
// The reference date used in the layouts is the same date used by the
// time package in the standard library:
//     Monday, Jan 2, 2006
// To define your own format, write down what the reference date would look
// like formatted your way; see the values of the predefined layouts for
// examples. The model is to demonstrate what the reference date looks like
// so that the Parse function and Format method can apply the same
// transformation to a general date value.
const (
	ISO8601 = "2006-01-02" // ISO 8601 extended format
	ISO8601B = "20060102"   // ISO 8601 basic format
	RFC822 = "02-Jan-06"
	RFC822W = "Mon, 02-Jan-06" // RFC822 with day of the week
	RFC850 = "Monday, 02-Jan-06"
	RFC1123 = "02 Jan 2006"
	RFC1123W = "Mon, 02 Jan 2006" // RFC1123 with day of the week
	RFC3339 = "2006-01-02"
)

// ParseISO parses an ISO 8601 formatted string and returns the date value it represents.
// In addition to the common formats (e.g. 2006-01-02 and 20060102), this function
// accepts date strings using the expanded year representation
// with possibly extra year digits beyond the prescribed four-digit minimum
// and with a + or - sign prefix (e.g. , "+12345-06-07", "-0987-06-05").
//
// Note that ParseISO is a little looser than the ISO 8601 standard and will
// be happy to parse dates with a year longer in length than the four-digit minimum even
// if they are missing the + sign prefix.
//
// Function Date.Parse can be used to parse date strings in other formats, but it
// is currently not able to parse ISO 8601 formatted strings that use the
// expanded year format.
//
// Background: https://en.wikipedia.org/wiki/ISO_8601#Dates
func ParseISO(value string) (Date, error) {
	if len(value) < 8 {
		return Date{}, fmt.Errorf("Date.ParseISO: cannot parse %s: incorrect length", value)
	}

	abs := value
	if value[0] == '+' || value[0] == '-' {
		abs = value[1:]
	}

	dash1 := strings.IndexByte(abs, '-')
	fm1 := dash1 + 1
	fm2 := dash1 + 3
	fd1 := dash1 + 4
	fd2 := dash1 + 6

	if dash1 < 0 {
		// switch to YYYYMMDD format
		dash1 = 4
		fm1 = 4
		fm2 = 6
		fd1 = 6
		fd2 = 8
	} else if abs[fm2] != '-' {
		return Date{}, fmt.Errorf("Date.ParseISO: cannot parse %s: incorrect syntax", value)
	}
	//fmt.Printf("%s %d %d %d %d %d\n", value, dash1, fm1, fm2, fd1, fd2)

	if len(abs) != fd2 {
		return Date{}, fmt.Errorf("Date.ParseISO: cannot parse %s: incorrect length", value)
	}

	year, err := parseField(value, abs[:dash1], "year", 4, -1)
	if err != nil {
		return Date{}, err
	}

	month, err := parseField(value, abs[fm1 : fm2], "month", -1, 2)
	if err != nil {
		return Date{}, err
	}

	day, err := parseField(value, abs[fd1:], "day", -1, 2)
	if err != nil {
		return Date{}, err
	}

	if value[0] == '-' {
		year = -year
	}

	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	return Date{encode(t)}, nil
}

func parseField(value, field, name string, minLength, requiredLength int) (int, error) {
	if (minLength > 0 && len(field) < minLength) || (requiredLength > 0 && len(field) != requiredLength) {
		return 0, fmt.Errorf("Date.ParseISO: cannot parse %s: invalid %s", value, name)
	}
	number, err := strconv.Atoi(field)
	if err != nil {
		return 0, fmt.Errorf("Date.ParseISO: cannot parse %s: invalid %s", value, name)
	}
	return number, nil
}

// Parse parses a formatted string and returns the Date value it represents.
// The layout defines the format by showing how the reference date, defined
// to be
//     Monday, Jan 2, 2006
// would be interpreted if it were the value; it serves as an example of the
// input format. The same interpretation will then be made to the input string.
//
// This function actually uses time.Parse to parse the input and can use any
// layout accepted by time.Parse, but returns only the date part of the
// parsed Time value.
//
// This function cannot currently parse ISO 8601 strings that use the expanded
// year format; you should use Date.ParseISO to parse those strings correctly.
func Parse(layout, value string) (Date, error) {
	t, err := time.Parse(layout, value)
	if err != nil {
		return Date{}, err
	}
	return Date{encode(t)}, nil
}

// String returns the time formatted in ISO 8601 extended format
// (e.g. "2006-01-02").  If the year of the date falls outside the
// [0,9999] range, this format produces an expanded year representation
// with possibly extra year digits beyond the prescribed four-digit minimum
// and with a + or - sign prefix (e.g. , "+12345-06-07", "-0987-06-05").
func (d Date) String() string {
	year, month, day := d.Date()
	if 0 <= year && year < 10000 {
		return fmt.Sprintf("%04d-%02d-%02d", year, month, day)
	}
	return fmt.Sprintf("%+05d-%02d-%02d", year, month, day)
}

// FormatISO returns a textual representation of the date value formatted
// according to the expanded year variant of the ISO 8601 extended format;
// the year of the date is represented as a signed integer using the
// specified number of digits (ignored if less than four).
// The string representation of the year will take more than the specified
// number of digits if the magnitude of the year is too large to fit.
//
// Function Date.Format can be used to format Date values in other formats,
// but it is currently not able to format dates according to the expanded
// year variant of the ISO 8601 format.
func (d Date) FormatISO(yearDigits int) string {
	n := 5 // four-digit minimum plus sign
	if yearDigits > 4 {
		n += yearDigits - 4
	}
	year, month, day := d.Date()
	return fmt.Sprintf("%+0*d-%02d-%02d", n, year, month, day)
}

// Format returns a textual representation of the date value formatted according
// to layout, which defines the format by showing how the reference date,
// defined to be
//     Mon, Jan 2, 2006
// would be displayed if it were the value; it serves as an example of the
// desired output.
//
// This function actually uses time.Format to format the input and can use any
// layout accepted by time.Format by extending its date to a time at
// 00:00:00.000 UTC.
//
// Additionally, it is able to insert the day-number suffix into the output string.
// This is done by including "nd" in the format string, which will become
//     Mon, Jan 2nd, 2006
// For example, New Year's Day might be rendered as "Fri, Jan 1st, 2016". To alter
// the suffix strings for a different locale, change DaySuffixes or use FormatWithSuffixes
// instead.
//
// This function cannot currently format Date values according to the expanded
// year variant of ISO 8601; you should use Date.FormatISO to that effect.
func (d Date) Format(layout string) string {
	return d.FormatWithSuffixes(layout, DaySuffixes)
}

// FormatWithSuffixes is the same as Format, except the suffix strings can be specified
// explicitly, which allows multiple locales to be supported. The suffixes slice should
// contain 31 strings covering the days 1 (index 0) to 31 (index 30).
func (d Date) FormatWithSuffixes(layout string, suffixes []string) string {
	t := decode(d.day)
	parts := strings.Split(layout, "nd")
	switch len(parts) {
	case 1:
		return t.Format(layout)

	default:
		a := make([]string, 0, 2 * len(parts) - 1)
		for i, p := range parts {
			if i > 0 {
				a = append(a, suffixes[d.Day() - 1])
			}
			a = append(a, t.Format(p))
		}
		return strings.Join(a, "")
	}
}

// DaySuffixes is the default array of strings used as suffixes when a format string
// contains "nd" (as in "second"). This can be altered at startup in order to change
// the default locale strings used for formatting dates. It supports every locale that
// uses the Gregorian calendar and has a suffix after the day-of-month number.
var DaySuffixes = []string{
	"st", "nd", "rd", "th", "th", // 1 - 5
	"th", "th", "th", "th", "th", // 6 - 10
	"th", "th", "th", "th", "th", // 11 - 15
	"th", "th", "th", "th", "th", // 16 - 20
	"st", "nd", "rd", "th", "th", // 21 - 25
	"th", "th", "th", "th", "th", // 26 - 30
	"st", // 31
}

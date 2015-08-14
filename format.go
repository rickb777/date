// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package date

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
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
	ISO8601  = "2006-01-02" // ISO 8601 extended format
	ISO8601B = "20060102"   // ISO 8601 basic format
	RFC822   = "02-Jan-06"
	RFC822W  = "Mon, 02-Jan-06" // RFC822 with day of the week
	RFC850   = "Monday, 02-Jan-06"
	RFC1123  = "02 Jan 2006"
	RFC1123W = "Mon, 02 Jan 2006" // RFC1123 with day of the week
	RFC3339  = "2006-01-02"
)

// reISO8601 is the regular expression used to parse date strings in the
// ISO 8601 extended format, with or without an expanded year representation.
var reISO8601 = regexp.MustCompile(`^([-+]?\d{4,})-(\d{2})-(\d{2})$`)

// ParseISO parses an ISO 8601 formatted string and returns the date value it represents.
// In addition to the common extended format (e.g. 2006-01-02), this function
// accepts date strings using the expanded year representation
// with possibly extra year digits beyond the prescribed four-digit minimum
// and with a + or - sign prefix (e.g. , "+12345-06-07", "-0987-06-05").
//
// Note that ParseISO is a little looser than the ISO 8601 standard and will
// be happy to parse dates with a year longer than the four-digit minimum even
// if they are missing the + sign prefix.
//
// Function Date.Parse can be used to parse date strings in other formats, but it
// is currently not able to parse ISO 8601 formatted strings that use the
// expanded year format.
func ParseISO(value string) (Date, error) {
	m := reISO8601.FindStringSubmatch(value)
	if len(m) != 4 {
		return Date{}, fmt.Errorf("Date.ParseISO: cannot parse %s", value)
	}
	// No need to check for errors since the regexp guarantees the matches
	// are valid integers
	year, _ := strconv.Atoi(m[1])
	month, _ := strconv.Atoi(m[2])
	day, _ := strconv.Atoi(m[3])

	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	return Date{encode(t)}, nil
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
		return Date{0}, err
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
// This function cannot currently format Date values according to the expanded
// year variant of ISO 8601; you should use Date.FormatISO to that effect.
func (d Date) Format(layout string) string {
	t := decode(d.day)
	return t.Format(layout)
}

// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package date

import (
	"fmt"
	"io"
	"strings"
)

// These are predefined layouts for use in Date.Format and Date.Parse.
// The reference date used in the layouts is the same date used by the
// time package in the standard library:
//
//	Monday, Jan 2, 2006
//
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

// String returns the time formatted in ISO 8601 extended format
// (e.g. "2006-01-02").  If the year of the date falls outside the
// [0,9999] range, this format produces an expanded year representation
// with possibly extra year digits beyond the prescribed four-digit minimum
// and with a + or - sign prefix (e.g. , "+12345-06-07", "-0987-06-05").
func (d Date) String() string {
	buf := &strings.Builder{}
	buf.Grow(12)
	d.WriteTo(buf)
	return buf.String()
}

// WriteTo is as per String, albeit writing to an io.Writer.
func (d Date) WriteTo(w io.Writer) (n64 int64, err error) {
	var n int
	year, month, day := d.Date()
	if 0 <= year && year < 10000 {
		n, err = fmt.Fprintf(w, "%04d-%02d-%02d", year, month, day)
	} else {
		n, err = fmt.Fprintf(w, "%+05d-%02d-%02d", year, month, day)
	}
	return int64(n), err
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
//
//	Mon, Jan 2, 2006
//
// would be displayed if it were the value; it serves as an example of the
// desired output.
//
// This function actually uses time.Format to format the input and can use any
// layout accepted by time.Format by extending its date to a time at
// 00:00:00.000 UTC.
//
// Additionally, it is able to insert the day-number suffix into the output string.
// This is done by including "nd" in the format string, which will become
//
//	Mon, Jan 2nd, 2006
//
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
		// If the format contains "Monday", it has been split so repair it.
		i := 1
		for i < len(parts) {
			if i > 0 && strings.HasSuffix(parts[i-1], "Mo") && strings.HasPrefix(parts[i], "ay") {
				parts[i-1] = parts[i-1] + "nd" + parts[i]
				copy(parts[i:], parts[i+1:])
				parts = parts[:len(parts)-1]
			} else {
				i++
			}
		}
		a := make([]string, 0, 2*len(parts)-1)
		for i, p := range parts {
			if i > 0 {
				a = append(a, suffixes[d.Day()-1])
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

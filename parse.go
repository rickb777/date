// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package date

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// MustAutoParse is as per AutoParse except that it panics if the string cannot be parsed.
// This is intended for setup code; don't use it for user inputs.
func MustAutoParse(value string) Date {
	d, err := AutoParse(value)
	if err != nil {
		panic(err)
	}
	return d
}

// AutoParse is like ParseISO, except that it automatically adapts to a variety of date formats
// provided that they can be detected unambiguously. Specifically, this includes the "European"
// and "British" date formats but not the common US format. Surrounding whitespace is ignored.
// The supported formats are:
//
// * all formats supported by ParseISO
//
// * yyyy/mm/dd | yyyy.mm.dd (or any similar pattern)
//
// * dd/mm/yyyy | dd.mm.yyyy (or any similar pattern)
//
// * surrounding whitespace is ignored
func AutoParse(value string) (Date, error) {
	abs := strings.TrimSpace(value)
	if len(abs) == 0 {
		return Date{}, errors.New("Date.AutoParse: cannot parse a blank string")
	}

	sign := ""
	if abs[0] == '+' || abs[0] == '-' {
		sign = abs[:1]
		abs = abs[1:]
	}

	if len(abs) >= 10 {
		i1 := -1
		i2 := -1
		for i, r := range abs {
			if unicode.IsPunct(r) {
				if i1 < 0 {
					i1 = i
				} else {
					i2 = i
				}
			}
		}
		if i1 >= 4 && i2 > i1 && abs[i1] == abs[i2] {
			// just normalise the punctuation
			a := []byte(abs)
			a[i1] = '-'
			a[i2] = '-'
			abs = string(a)
		} else if i1 >= 2 && i2 > i1 && abs[i1] == abs[i2] {
			// harder case - need to swap the field order
			dd := abs[0:i1]
			mm := abs[i1+1 : i2]
			yyyy := abs[i2+1:]
			abs = fmt.Sprintf("%s-%s-%s", yyyy, mm, dd)
		}
	}
	return ParseISO(sign + abs)
}

// MustParseISO is as per ParseISO except that it panics if the string cannot be parsed.
// This is intended for setup code; don't use it for user inputs.
func MustParseISO(value string) Date {
	d, err := ParseISO(value)
	if err != nil {
		panic(err)
	}
	return d
}

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
// Function date.Parse can be used to parse date strings in other formats, but it
// is currently not able to parse ISO 8601 formatted strings that use the
// expanded year format.
//
// Background: https://en.wikipedia.org/wiki/ISO_8601#Dates
func ParseISO(value string) (Date, error) {
	if len(value) < 8 {
		return Date{}, fmt.Errorf("Date.ParseISO: cannot parse %q: incorrect length", value)
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
		return Date{}, fmt.Errorf("Date.ParseISO: cannot parse %q: incorrect syntax", value)
	}
	//fmt.Printf("%s %d %d %d %d %d\n", value, dash1, fm1, fm2, fd1, fd2)

	if len(abs) != fd2 {
		return Date{}, fmt.Errorf("Date.ParseISO: cannot parse %q: incorrect length", value)
	}

	year, err := parseField(value, abs[:dash1], "year", 4, -1)
	if err != nil {
		return Date{}, err
	}

	month, err := parseField(value, abs[fm1:fm2], "month", -1, 2)
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
		return 0, fmt.Errorf("Date.ParseISO: cannot parse %q: invalid %s", value, name)
	}
	number, err := strconv.Atoi(field)
	if err != nil {
		return 0, fmt.Errorf("Date.ParseISO: cannot parse %q: invalid %s", value, name)
	}
	return number, nil
}

// MustParse is as per Parse except that it panics if the string cannot be parsed.
// This is intended for setup code; don't use it for user inputs.
func MustParse(layout, value string) Date {
	d, err := Parse(layout, value)
	if err != nil {
		panic(err)
	}
	return d
}

// Parse parses a formatted string of a known layout and returns the Date value it represents.
// The layout defines the format by showing how the reference date, defined
// to be
//
//	Monday, Jan 2, 2006
//
// would be interpreted if it were the value; it serves as an example of the
// input format. The same interpretation will then be made to the input string.
//
// This function actually uses time.Parse to parse the input and can use any
// layout accepted by time.Parse, but returns only the date part of the
// parsed Time value.
//
// This function cannot currently parse ISO 8601 strings that use the expanded
// year format; you should use date.ParseISO to parse those strings correctly.
func Parse(layout, value string) (Date, error) {
	t, err := time.Parse(layout, value)
	if err != nil {
		return Date{}, err
	}
	return Date{encode(t)}, nil
}

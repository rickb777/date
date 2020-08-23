// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package period

import (
	"fmt"
	"strconv"
	"strings"
)

// MustParse is as per Parse except that it panics if the string cannot be parsed.
// This is intended for setup code; don't use it for user inputs.
func MustParse(value string) Period {
	d, err := Parse(value)
	if err != nil {
		panic(err)
	}
	return d
}

// Parse parses strings that specify periods using ISO-8601 rules.
//
// In addition, a plus or minus sign can precede the period, e.g. "-P10D"
//
// The value is normalised, e.g. multiple of 12 months become years so "P24M"
// is the same as "P2Y". However, this is done without loss of precision, so
// for example whole numbers of days do not contribute to the months tally
// because the number of days per month is variable.
//
// The zero value can be represented in several ways: all of the following
// are equivalent: "P0Y", "P0M", "P0W", "P0D", "PT0H", PT0M", PT0S", and "P0".
// The canonical zero is "P0D".
func Parse(period string) (Period, error) {
	return ParseWithNormalise(period, true)
}

// ParseWithNormalise parses strings that specify periods using ISO-8601 rules
// with an option to specify whether to normalise parsed period components.
//
// In addition, a plus or minus sign can precede the period, e.g. "-P10D"

// The returned value is only normalised when normalise is set to `true`, and
// normalisation will convert e.g. multiple of 12 months into years so "P24M"
// is the same as "P2Y". However, this is done without loss of precision, so
// for example whole numbers of days do not contribute to the months tally
// because the number of days per month is variable.
//
// The zero value can be represented in several ways: all of the following
// are equivalent: "P0Y", "P0M", "P0W", "P0D", "PT0H", PT0M", PT0S", and "P0".
// The canonical zero is "P0D".
func ParseWithNormalise(period string, normalise bool) (Period, error) {
	if period == "" || period == "-" || period == "+" {
		return Period{}, fmt.Errorf("cannot parse a blank string as a period")
	}

	if period == "P0" {
		return Period{}, nil
	}

	p64, err := parse(period, normalise)
	if err != nil {
		return Period{}, err
	}
	return p64.toPeriod()
}

func parse(period string, normalise bool) (*period64, error) {
	neg := false
	remaining := period
	if remaining[0] == '-' {
		neg = true
		remaining = remaining[1:]
	} else if remaining[0] == '+' {
		remaining = remaining[1:]
	}

	if remaining[0] != 'P' {
		return nil, fmt.Errorf("%s: expected 'P' period mark at the start", period)
	}
	remaining = remaining[1:]

	var whole int64
	var fraction int8
	result := &period64{input: period, neg: neg}
	var weekValue int64
	var years, months, weeks, days, hours, minutes, seconds itemState
	var des byte
	var err error
	nComponents := 0

	years, months, weeks, days = Armed, Armed, Armed, Armed

	isHMS := false
	for len(remaining) > 0 {
		if remaining[0] == 'T' {
			if isHMS {
				return nil, fmt.Errorf("%s: 'T' designator cannot occur more than once", period)
			}
			isHMS = true

			years, months, weeks, days = Unready, Unready, Unready, Unready
			hours, minutes, seconds = Armed, Armed, Armed

			remaining = remaining[1:]

		} else {
			whole, fraction, des, remaining, err = parseNextField(remaining, period)
			if err != nil {
				return nil, err
			}

			switch des {
			case 'Y':
				years, err = years.testAndSet(whole, fraction, Year, result, &result.years)
			case 'W':
				weeks, err = weeks.testAndSet(whole, fraction, Week, result, &weekValue)
			case 'D':
				days, err = days.testAndSet(whole, fraction, Day, result, &result.days)
			case 'H':
				hours, err = hours.testAndSet(whole, fraction, Hour, result, &result.hours)
			case 'S':
				seconds, err = seconds.testAndSet(whole, fraction, Second, result, &result.seconds)
			case 'M':
				if isHMS {
					minutes, err = minutes.testAndSet(whole, fraction, Minute, result, &result.minutes)
				} else {
					months, err = months.testAndSet(whole, fraction, Month, result, &result.months)
				}
			default:
				return nil, fmt.Errorf("%s: expected a number not '%c'", period, des)
			}
			nComponents++

			if err != nil {
				return nil, err
			}
		}
	}

	if nComponents == 0 {
		return nil, fmt.Errorf("%s: expected 'Y', 'M', 'W', 'D', 'H', 'M', or 'S' designator", period)
	}

	result.days += weekValue * 7

	if normalise {
		result = result.normalise64(true)
	}

	return result, nil
}

//-------------------------------------------------------------------------------------------------

type itemState int

const (
	Unready itemState = iota
	Armed
	Set
)

func (i itemState) testAndSet(whole int64, fraction int8, des designator, result *period64, value *int64) (itemState, error) {
	switch i {
	case Unready:
		return i, fmt.Errorf("%s: '%c' designator cannot occur here", result.input, des.Byte())
	case Set:
		return i, fmt.Errorf("%s: '%c' designator cannot occur more than once", result.input, des)
	}
	if fraction != 0 && result.fraction != 0 {
		return i, fmt.Errorf("%s: '%c' & '%c' only the last field can have a fraction", result.input, result.fpart, des.Byte())
	}

	if fraction != 0 {
		if des == Week {
			result.fraction = fraction * 7
			result.fpart = Day
		} else {
			result.fraction = fraction
			result.fpart = des
		}
	}

	*value = whole
	return Set, nil
}

//-------------------------------------------------------------------------------------------------

func parseNextField(str, original string) (int64, int8, byte, string, error) {
	i := scanDigits(str)
	if i < 0 {
		return 0, 0, 0, "", fmt.Errorf("%s: missing designator at the end", original)
	}

	des := str[i]
	whole, frac, err := parseDecimalNumber(str[:i], original, des)
	return whole, frac, des, str[i+1:], err
}

// Fixed-point three decimal places
func parseDecimalNumber(number, original string, des byte) (whole int64, frac int8, err error) {
	dec := strings.IndexByte(number, '.')
	if dec < 0 {
		dec = strings.IndexByte(number, ',')
	}

	if dec >= 0 {
		whole, err = strconv.ParseInt(number[:dec], 10, 64)
		if err == nil {
			var n int64
			number = number[dec+1:]
			switch len(number) {
			case 0: // skip
			case 1:
				n, err = strconv.ParseInt(number, 10, 64)
				frac = int8(n * 10)
			default:
				n, err = strconv.ParseInt(number[:2], 10, 64)
				frac = int8(n)
			}
		}
	} else {
		whole, err = strconv.ParseInt(number, 10, 64)
	}

	if err != nil {
		return whole, 0, fmt.Errorf("%s: expected a number but found '%c'", original, des)
	}

	return whole, frac, err
}

// scanDigits finds the first non-digit byte after a given starting point.
// Note that it does not care about runes or UTF-8 encoding; it assumes that
// a period string is always valid ASCII as well as UTF-8.
func scanDigits(s string) int {
	for i, c := range s {
		if !isDigit(c) {
			return i
		}
	}
	return -1
}

func isDigit(c rune) bool {
	return ('0' <= c && c <= '9') || c == '.' || c == ','
}

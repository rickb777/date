// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package period

import (
	"fmt"
	"math"
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
	if period == "" {
		return Period{}, fmt.Errorf("cannot parse a blank string as a period")
	}

	if period == "P0" {
		return Period{}, nil
	}

	neg := false
	remaining := period
	if remaining[0] == '-' {
		neg = true
		remaining = remaining[1:]
	} else if remaining[0] == '+' {
		remaining = remaining[1:]
	}

	if remaining[0] != 'P' {
		return Period{}, fmt.Errorf("expected 'P' period designator at the start: %s", period)
	}
	remaining = remaining[1:]

	var n cent64
	var years, months, weeks, days, hours, minutes, seconds item
	var designator byte
	var err error
	nComponents := 0

	years.armed = true
	months.armed = true
	weeks.armed = true
	days.armed = true

	isHMS := false
	for len(remaining) > 0 {
		if remaining[0] == 'T' {
			if isHMS {
				return Period{}, fmt.Errorf("'T' designator cannot occur more than once: %s", period)
			}
			isHMS = true

			years.armed = false
			months.armed = false
			weeks.armed = false
			days.armed = false
			hours.armed = true
			minutes.armed = true
			seconds.armed = true

			remaining = remaining[1:]

		} else {
			n, designator, remaining, err = parseNextField(remaining, period)
			if err != nil {
				return Period{}, err
			}

			switch designator {
			case 'T':
				if isHMS {
					return Period{}, fmt.Errorf("'T' designator cannot occur more than once: %s", period)
				}
				isHMS = true
				remaining = remaining[1:]
			case 'Y':
				nComponents++
				years, err = years.testAndSet(n, designator, period)
			case 'W':
				nComponents++
				weeks, err = weeks.testAndSet(n, designator, period)
			case 'D':
				nComponents++
				days, err = days.testAndSet(n, designator, period)
			case 'H':
				nComponents++
				hours, err = hours.testAndSet(n, designator, period)
			case 'S':
				nComponents++
				seconds, err = seconds.testAndSet(n, designator, period)
			case 'M':
				nComponents++
				if isHMS {
					minutes, err = minutes.testAndSet(n, designator, period)
				} else {
					months, err = months.testAndSet(n, designator, period)
				}
			default:
				return Period{}, fmt.Errorf("expected a number not '%c': %s", designator, period)
			}

			if err != nil {
				return Period{}, err
			}
		}
	}

	if nComponents == 0 {
		return Period{}, fmt.Errorf("expected 'Y', 'M', 'W', 'D', 'H', 'M', or 'S' designator: %s", period)
	}

	result := period64{
		centiMonths:  years.value*12 + months.value,
		centiDays:    weeks.value*7 + days.value,
		centiSeconds: (hours.value * 3600) + (minutes.value * 60) + seconds.value,
		neg:          neg,
		showAs:       period,
	}

	if normalise {
		result = result.normalise64(true)
	}

	err = nil
	m := result.overflowedFields()
	if len(m) > 0 {
		kind := "non-normalised"
		if normalise {
			kind = "normalised"
		}
		err = fmt.Errorf("%s period overflows %s: %s", kind, strings.Join(m, ","), period)
	}

	return result.toPeriod(), err
}

//-------------------------------------------------------------------------------------------------

type item struct {
	value      cent64
	armed, set bool
}

func (i item) overflows() bool {
	return i.value > math.MaxInt32
}

func (i item) testAndSet(v cent64, designator byte, original string) (item, error) {
	if !i.armed {
		return i, fmt.Errorf("'%c' designator cannot occur here: %s", designator, original)
	}
	if i.set {
		return i, fmt.Errorf("'%c' designator cannot occur more than once: %s", designator, original)
	}
	i.value = v
	i.set = true
	return i, nil
}

//-------------------------------------------------------------------------------------------------

func parseNextField(s, original string) (cent64, byte, string, error) {
	i := scanDigits(s)
	if i < 0 {
		return 0, 0, "", fmt.Errorf("missing designator at the end: %s", original)
	}

	designator := s[i]
	n, err := parseDecimalFixedPoint(s[:i], original, designator)
	return n, designator, s[i+1:], err
}

// Fixed-point three decimal places
func parseDecimalFixedPoint(s, original string, designator byte) (cent64, error) {
	dec := strings.IndexByte(s, '.')
	if dec < 0 {
		dec = strings.IndexByte(s, ',')
	}

	if dec >= 0 {
		dp := len(s) - dec
		if dp > 2 {
			s = s[:dec] + s[dec+1:dec+3]
		} else if dp > 1 {
			s = s[:dec] + s[dec+1:dec+2] + "0"
		} else {
			s = s[:dec] + s[dec+1:] + "00"
		}
	} else {
		s = s + "00"
	}

	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return cent64(n), fmt.Errorf("expected a number before the '%c' designator: %s", designator, original)
	}
	return cent64(n), nil
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

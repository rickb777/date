// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package period

import (
	"fmt"
	"strconv"
	"strings"
)

type NormalisationMode int

const (
	// Verbatim is the mode that disables normalisation during parsing. Parsing
	// will fail for inputs that do not fit the numeric range (int16 using one
	// fixed decimal place), which means ± 2^16 / 10 (i.e. ±3276.7).
	Verbatim NormalisationMode = iota

	// Constrained is the mode that only allows normalisation if it is essential to
	// avoid numeric overflow. Otherwise it is like Verbatim. This is the default mode.
	Constrained

	// Normalised is the main normalisation mode: this is done without loss of
	// precision, so for example multiples of 60 seconds become minutes, and multiples
	// of sixty minutes become hours.
	//
	// However, whole multiples of days or weeks do not contribute to the months
	// tally because the number of days per month is variable. Also note that the
	// number of hours per day is not always 24, due to daylight savings changes,
	// so multiples of 24 hours do not contribute to days.
	//
	// Normalisation also has heuristics to minimise fractions where they can be
	// carried right into a less-significant field.
	Normalised

	// Imprecise is the normalisation mode that aggressively normalises input values
	// making assumptions that days are all 24 hours and that multiples of days and
	// weeks can be carried to months and years according to the Gregorian rule of
	// 365.2425 days per year.
	Imprecise
)

// DefaultNormalisation is Constrained but you can change this.
var DefaultNormalisation = Constrained

//-------------------------------------------------------------------------------------------------

// MustParse is as per Parse except that it panics if the string cannot be parsed.
// This is intended for setup code; don't use it for user inputs.
// By default, the value is normalised.
// Normalisation can be disabled using the optional flag.
func MustParse(value string, normalise ...NormalisationMode) Period {
	d, err := Parse(value, normalise...)
	if err != nil {
		panic(err)
	}
	return d
}

// Parse parses strings that specify periods using ISO-8601 rules.
//
// In addition, a plus or minus sign can precede the period, e.g. "-P10D"
//
// Normalisation is controlled by the optional parameter and the value of
// DefaultNormalisation.
//
// The zero value can be represented in several ways: all of the following
// are equivalent: "P0Y", "P0M", "P0W", "P0D", "PT0H", PT0M", PT0S", and "P0".
// The canonical zero is "P0D".
func Parse(period string, normaliseOpt ...NormalisationMode) (Period, error) {
	normalise := DefaultNormalisation
	if len(normaliseOpt) > 0 {
		normalise = normaliseOpt[0]
	}

	if period == "" || period == "-" || period == "+" {
		return Period{}, fmt.Errorf("cannot parse a blank string as a period")
	}

	if period == "P0" {
		return Period{}, nil
	}

	p64, err := parse(period)
	if err != nil {
		return Period{}, err
	}

	if normalise == Constrained && p64.checkOverflow() != nil {
		normalise = Normalised // bump it up
	}

	if normalise >= Normalised {
		p64 = p64.normalise64(normalise < Imprecise)
	}

	return p64.toPeriod(), p64.checkOverflow()
}

func parse(period string) (*period64, error) {
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

	p64 := &period64{input: period, neg: neg}

	var number, prevFraction int64
	var years, months, weeks, days, hours, minutes, seconds itemState
	var designator, prevDesignator byte
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
			number, designator, remaining, err = parseNextField(remaining, period)
			if err != nil {
				return nil, err
			}

			fraction := number % 10
			if prevFraction != 0 && fraction != 0 {
				return nil, fmt.Errorf("%s: '%c' & '%c' only the last field can have a fraction", period, prevDesignator, designator)
			}

			switch designator {
			case 'Y':
				years, err = years.testAndSet(number, 'Y', p64, &p64.years)
			case 'W':
				weeks, err = weeks.testAndSet(number, 'W', p64, &p64.weeks)
			case 'D':
				days, err = days.testAndSet(number, 'D', p64, &p64.days)
			case 'H':
				hours, err = hours.testAndSet(number, 'H', p64, &p64.hours)
			case 'S':
				seconds, err = seconds.testAndSet(number, 'S', p64, &p64.seconds)
			case 'M':
				if isHMS {
					minutes, err = minutes.testAndSet(number, 'M', p64, &p64.minutes)
				} else {
					months, err = months.testAndSet(number, 'M', p64, &p64.months)
				}
			default:
				return nil, fmt.Errorf("%s: expected a designator Y, M, W, D, H, or S not '%c'", period, designator)
			}
			nComponents++

			if err != nil {
				return nil, err
			}

			prevFraction = fraction
			prevDesignator = designator
		}
	}

	if nComponents == 0 {
		return nil, fmt.Errorf("%s: expected 'Y', 'M', 'W', 'D', 'H', 'M', or 'S' designator", period)
	}

	p64.denormal = p64.months >= 120 || p64.weeks >= 520 || p64.days >= 70 ||
		p64.hours >= 240 || p64.minutes >= 600 || p64.seconds >= 600

	return p64, nil
}

//-------------------------------------------------------------------------------------------------

type itemState int

const (
	Unready itemState = iota
	Armed
	Set
)

func (i itemState) testAndSet(number int64, designator byte, result *period64, value *int64) (itemState, error) {
	switch i {
	case Unready:
		return i, fmt.Errorf("%s: '%c' designator cannot occur here", result.input, designator)
	case Set:
		return i, fmt.Errorf("%s: '%c' designator cannot occur more than once", result.input, designator)
	}

	*value = number
	return Set, nil
}

//-------------------------------------------------------------------------------------------------

func parseNextField(str, original string) (int64, byte, string, error) {
	i := scanDigits(str)
	if i < 0 {
		return 0, 0, "", fmt.Errorf("%s: missing designator at the end", original)
	}

	des := str[i]
	number, err := parseDecimalNumber(str[:i], original, des)
	return number, des, str[i+1:], err
}

// Fixed-point one decimal place
func parseDecimalNumber(number, original string, des byte) (int64, error) {
	dec := strings.IndexByte(number, '.')
	if dec < 0 {
		dec = strings.IndexByte(number, ',')
	}

	var integer, fraction int64
	var err error
	if dec >= 0 {
		integer, err = strconv.ParseInt(number[:dec], 10, 64)
		if err == nil {
			number = number[dec+1:]
			switch len(number) {
			case 0: // skip
			case 1:
				fraction, err = strconv.ParseInt(number, 10, 64)
			default:
				fraction, err = strconv.ParseInt(number[:1], 10, 64)
			}
		}
	} else {
		integer, err = strconv.ParseInt(number, 10, 64)
	}

	if err != nil {
		return 0, fmt.Errorf("%s: expected a number but found '%c'", original, des)
	}

	return integer*10 + fraction, err
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

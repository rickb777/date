// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package period

import (
	"time"
)

// Add adds two periods together. Use this method along with Negate in order to subtract periods.
//
// The result is not normalised and may overflow arithmetically (to make this unlikely, use Normalise on
// the inputs before adding them).
func (period Period) Add(that Period) Period {
	years := period.years + that.years
	months := period.months + that.months
	weeks := period.weeks + that.weeks
	days := period.days + that.days
	hours := period.hours + that.hours
	minutes := period.minutes + that.minutes
	seconds := period.seconds + that.seconds

	denormal1 := months >= 120 || days > 300 || hours >= 240 || minutes >= 600 || seconds >= 600
	denormal2 := months <= -120 || days < -300 || hours <= -240 || minutes <= -600 || seconds <= -600

	return Period{
		years: years, months: months, weeks: weeks, days: days,
		hours: hours, minutes: minutes, seconds: seconds,
		denormal: denormal1 || denormal2,
	}
}

//-------------------------------------------------------------------------------------------------

// AddTo adds the period to a time, returning the result.
// A flag is also returned that is true when the conversion was precise, and false otherwise.
//
// When the period specifies hours, minutes and seconds only, the result is precise.
//
// Similarly, when the period specifies whole years, months, weeks and days (i.e. without fractions),
// the result is precise.
//
// However, when years, months or days contains fractions, the result is only an approximation (it
// assumes that all days are 24 hours and every year is 365.2425 days, as per Gregorian calendar rules).
func (period Period) AddTo(t time.Time) (time.Time, bool) {
	wholeYears := (period.years % 10) == 0
	wholeMonths := (period.months % 10) == 0
	wholeWeeks := (period.weeks % 10) == 0
	wholeDays := (period.days % 10) == 0

	if wholeYears && wholeMonths && wholeWeeks && wholeDays {
		// in this case, time.AddDate provides an exact solution
		stE3 := totalSecondsE3(period)
		t1 := t.AddDate(int(period.years/10), int(period.months/10), 7*int(period.weeks/10)+int(period.days/10))
		return t1.Add(stE3 * time.Millisecond), true
	}

	d, precise := period.Duration()
	return t.Add(d), precise
}

//-------------------------------------------------------------------------------------------------

// Scale a period by a multiplication factor. Obviously, this can both enlarge and shrink it,
// and change the sign if the factor is negative. The result is normalised, but integer overflows
// are silently ignored.
//
// Bear in mind that the internal representation is limited by fixed-point arithmetic with two
// decimal places; each field is only int16.
//
// Known issue: scaling by a large reduction factor (i.e. much less than one) doesn't work properly.
func (period Period) Scale(factor float32) Period {
	result, _ := period.ScaleWithOverflowCheck(factor)
	return result
}

// ScaleWithOverflowCheck scales a period by a multiplication factor. Obviously, this can both
// enlarge and shrink it, and change the sign if negative. The result is normalised. An error
// is returned if integer overflow happened.
//
// Bear in mind that the internal representation is limited by fixed-point arithmetic with one
// decimal place; each field is only int16.
//
// Known issue: scaling by a large reduction factor (i.e. much less than one) doesn't work properly.
func (period Period) ScaleWithOverflowCheck(factor float32) (Period, error) {
	ap, neg := period.absNeg()

	if -0.5 < factor && factor < 0.5 {
		d, pr1 := ap.Duration()
		mul := float64(d) * float64(factor)
		p2, pr2 := NewOf(time.Duration(mul))
		return p2.Normalise(pr1 && pr2), nil
	}

	y := int64(float32(ap.years) * factor)
	m := int64(float32(ap.months) * factor)
	w := int64(float32(ap.weeks) * factor)
	d := int64(float32(ap.days) * factor)
	hh := int64(float32(ap.hours) * factor)
	mm := int64(float32(ap.minutes) * factor)
	ss := int64(float32(ap.seconds) * factor)

	p64 := &period64{years: y, months: m, weeks: w, days: d, hours: hh, minutes: mm, seconds: ss, neg: neg, denormal: true}
	n64 := p64.normalise64(true)
	return n64.toPeriod(), n64.checkOverflow()
}

// RationalScale scales a period by a rational multiplication factor. Obviously, this can both enlarge and shrink it,
// and change the sign if negative. The result is normalised. An error is returned if integer overflow
// happened.
//
// If the divisor is zero, a panic will arise.
//
// Bear in mind that the internal representation is limited by fixed-point arithmetic with two
// decimal places; each field is only int16.
//func (period Period) RationalScale(multiplier, divisor int) (Period, error) {
//	return period.rationalScale64(int64(multiplier), int64(divisor))
//}

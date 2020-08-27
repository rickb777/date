// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package period

import (
	"fmt"
	"math/big"
	"time"
)

// Add adds two periods together. Use this method along with Negate in order to subtract periods.
//
// The result is not normalised and may overflow arithmetically (to make this unlikely, use Normalise on
// the inputs before adding them).
func (period Period) Add(that Period) Period {
	fpart := period.fpart
	if period.fpart != that.fpart {
		//TODO
	}

	fraction := period.fraction + that.fraction
	if fraction == 0 {
		fpart = NoFraction
	}

	return Period{
		years:    period.years + that.years,
		months:   period.months + that.months,
		days:     period.days + that.days,
		hours:    period.hours + that.hours,
		minutes:  period.minutes + that.minutes,
		seconds:  period.seconds + that.seconds,
		fraction: fraction,
		fpart:    fpart,
	}
}

//-------------------------------------------------------------------------------------------------

// AddTo adds the period to a time, returning the result.
// A flag is also returned that is true when the conversion was precise and false otherwise.
//
// When the period specifies hours, minutes and seconds only, the result is precise.
// Also, when the period specifies whole years, months and days (i.e. without fractions), the
// result is precise. However, when years, months or days contains fractions, the result
// is only an approximation (it assumes that all days are 24 hours and every year is 365.2425
// days, as per Gregorian calendar rules).
func (period Period) AddTo(t time.Time) (time.Time, bool) {
	wholeYears := period.fpart != Year
	wholeMonths := period.fpart != Month
	wholeDays := period.fpart != Day

	if wholeYears && wholeMonths && wholeDays {
		// in this case, time.AddDate provides an exact solution
		t1 := t.AddDate(int(period.years), int(period.months), int(period.days))
		return t1.Add(period.hmsDuration()), true
	}

	d, precise := period.Duration()
	return t.Add(d), precise
}

//-------------------------------------------------------------------------------------------------

// Scale a period by a multiplication factor. Obviously, this can both enlarge and shrink it,
// and change the sign if negative. The result is normalised, but integer overflows are silently
// ignored.
//
// Bear in mind that the internal representation is limited by fixed-point arithmetic with two
// decimal places; each field is only int16.
//
// Known issue: scaling by a large reduction factor (i.e. much less than one) doesn't work properly.
func (period Period) Scale(factor float32) Period {
	result, _ := period.ScaleWithOverflowCheck(factor)
	return result
}

// ScaleWithOverflowCheck a period by a multiplication factor. Obviously, this can both enlarge and shrink it,
// and change the sign if negative. The result is normalised. An error is returned if integer overflow
// happened.
//
// Bear in mind that the internal representation is limited by fixed-point arithmetic with two
// decimal places; each field is only int16.
func (period Period) ScaleWithOverflowCheck(factor float32) (Period, error) {
	str := fmt.Sprintf("%f", factor)
	bigRat, ok := new(big.Rat).SetString(str)
	if !ok {
		return Period{}, fmt.Errorf("unable to scale period %s using %f", period, factor)
	}

	multiplier64 := bigRat.Num().Int64()
	divisor64 := bigRat.Denom().Int64()
	return period.rationalScale64(multiplier64, divisor64)
}

// RationalScale scales a period by a rational multiplication factor. Obviously, this can both enlarge and shrink it,
// and change the sign if negative. The result is normalised. An error is returned if integer overflow
// happened.
//
// If the divisor is zero, a panic will arise.
//
// Bear in mind that the internal representation is limited by fixed-point arithmetic with two
// decimal places; each field is only int16.
func (period Period) RationalScale(multiplier, divisor int) (Period, error) {
	return period.rationalScale64(int64(multiplier), int64(divisor))
}

func (period Period) rationalScale64(m, d int64) (Period, error) {
	ap, neg := period.absNeg()

	cy := ap.centiYears()
	cm := ap.centiMonths()
	cd := ap.centiDays()
	chh := ap.centiHours()
	cmm := ap.centiMinutes()
	css := ap.centiSeconds()

	cym := cy*12 + cm
	chms := (chh * 3600) + (cmm * 60) + css

	mcym := cym * m
	mcd := cd * m
	mchms := chms * m

	cymr := mcym % d
	cdr := mcd % d
	chmsr := mchms % d

	if cymr == 0 && cdr == 0 && chmsr == 0 {
		// special case: scaled result is integral
		scd := mcd / d
		if d > m && scd*d != mcd {
			mchms = mcd * 24
			mcd = 0
		}
		return integralScaledResult(mcym/d, scd, mchms/d, neg)
	}

	// fall back on reliable but approximate algorithm
	ymdDuration := time.Duration(cym*daysPerMonthE6+cd*oneE6) * 864 * time.Microsecond
	hmsDuration := time.Duration(chms) * 10 * time.Millisecond
	duration := ymdDuration + hmsDuration
	pr1 := ymdDuration == 0
	mul := (int64(duration) * m) / d
	// add 5ms to round half-up
	p2, pr2 := NewOf(time.Duration(mul) + 5*time.Millisecond)
	precise := pr1 && pr2
	return p2.condNegate(neg).Normalise(precise).Simplify(precise), nil
}

func integralScaledResult(ymi, di, hmsi int64, neg bool) (Period, error) {
	p64 := &period64{
		months:  ymi / 100,
		days:    di / 100,
		seconds: hmsi / 100,
		neg:     neg,
	}

	ymf := ymi % 100
	if ymf != 0 {
		p64.fraction = int8(ymf)
		p64.fpart = Month
	}

	df := di % 100
	if df != 0 {
		p64.fraction = int8(df)
		p64.fpart = Day
	}

	sf := hmsi % 100
	if sf != 0 {
		p64.fraction = int8(sf)
		p64.fpart = Second
	}

	return p64.normalise64(true).toPeriod()
}

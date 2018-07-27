// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package period

import (
	"fmt"
	"time"
)

const daysPerYearE4 int64 = 3652425   // 365.2425 days by the Gregorian rule
const daysPerMonthE4 int64 = 304375   // 30.4375 days per month
const daysPerMonthE6 int64 = 30436875 // 30.436875 days per month

const oneE4 = 10000
const oneE5 = 100000
const oneE6 = 1000000
const oneE7 = 10000000

const hundredMs = 100 * time.Millisecond

// reminder: int64 overflow is after 9,223,372,036,854,775,807 (math.MaxInt64)

// Period holds a period of time and provides conversion to/from ISO-8601 representations.
// Therefore there are six fields: years, months, days, hours, minutes, and seconds.
//
// In the ISO representation, decimal fractions are supported, although only the last non-zero
// component is allowed to have a fraction according to the Standard. For example "P2.5Y"
// is 2.5 years.
//
// However, in this implementation, the precision is limited to one decimal place only, by
// means of integers with fixed point arithmetic. (This avoids using float32 in the struct,
// so there are no problems testing equality using ==.)
//
// The implementation limits the range of possible values to ± 2^16 / 10 in each field.
// Note in particular that the range of years is limited to approximately ± 3276.
//
// The concept of weeks exists in string representations of periods, but otherwise weeks
// are unimportant. The period contains a number of days from which the number of weeks can
// be calculated when needed.
//
// Note that although fractional weeks can be parsed, they will never be returned via String().
// This is because the number of weeks is always inferred from the number of days.
//
type Period struct {
	years, months, days, hours, minutes, seconds int16
}

// NewYMD creates a simple period without any fractional parts. The fields are initialised verbatim
// without any normalisation; e.g. 12 months will not become 1 year. Use the Normalise method if you
// need to.
//
// All the parameters must have the same sign (otherwise a panic occurs).
func NewYMD(years, months, days int) Period {
	return New(years, months, days, 0, 0, 0)
}

// NewHMS creates a simple period without any fractional parts. The fields are initialised verbatim
// without any normalisation; e.g. 120 seconds will not become 2 minutes. Use the Normalise method
// if you need to.
//
// All the parameters must have the same sign (otherwise a panic occurs).
func NewHMS(hours, minutes, seconds int) Period {
	return New(0, 0, 0, hours, minutes, seconds)
}

// New creates a simple period without any fractional parts. The fields are initialised verbatim
// without any normalisation; e.g. 120 seconds will not become 2 minutes. Use the Normalise method
// if you need to.
//
// All the parameters must have the same sign (otherwise a panic occurs).
func New(years, months, days, hours, minutes, seconds int) Period {
	if (years >= 0 && months >= 0 && days >= 0 && hours >= 0 && minutes >= 0 && seconds >= 0) ||
		(years <= 0 && months <= 0 && days <= 0 && hours <= 0 && minutes <= 0 && seconds <= 0) {
		return Period{
			int16(years) * 10, int16(months) * 10, int16(days) * 10,
			int16(hours) * 10, int16(minutes) * 10, int16(seconds) * 10,
		}
	}
	panic(fmt.Sprintf("Periods must have homogeneous signs; got P%dY%dM%dDT%dH%dM%dS",
		years, months, days, hours, minutes, seconds))
}

func newE1(years, months, days, hours, minutes, seconds int64) Period {
	return Period{
		int16(years), int16(months), int16(days),
		int16(hours), int16(minutes), int16(seconds),
	}
}

// TODO NewFloat

// NewOf converts a time duration to a Period, and also indicates whether the conversion is precise.
// Any time duration that spans more than ± 3276 hours will be approximated by assuming that there
// are 24 hours per day, 30.4375 per month and 365.25 days per year.
func NewOf(duration time.Duration) (p Period, precise bool) {
	var sign int16 = 1
	d := duration
	if duration < 0 {
		sign = -1
		d = -duration
	}

	sign10 := sign * 10

	totalHours := int64(d / time.Hour)

	// check for 16-bit overflow - occurs near the 4.5 month mark
	if totalHours < 3277 {
		// simple HMS case
		minutes := d % time.Hour / time.Minute
		seconds := d % time.Minute / hundredMs
		return Period{0, 0, 0, sign10 * int16(totalHours), sign10 * int16(minutes), sign * int16(seconds)}, true
	}

	totalDays := totalHours / 24 // ignoring daylight savings adjustments

	if totalDays < 3277 {
		hours := totalHours - totalDays*24
		minutes := d % time.Hour / time.Minute
		seconds := d % time.Minute / hundredMs
		return Period{0, 0, sign10 * int16(totalDays), sign10 * int16(hours), sign10 * int16(minutes), sign * int16(seconds)}, false
	}

	// TODO it is uncertain whether this is too imprecise and should be improved
	years := (oneE4 * totalDays) / daysPerYearE4
	months := ((oneE4 * totalDays) / daysPerMonthE4) - (12 * years)
	hours := totalHours - totalDays*24
	totalDays = ((totalDays * oneE4) - (daysPerMonthE4 * months) - (daysPerYearE4 * years)) / oneE4
	return Period{sign10 * int16(years), sign10 * int16(months), sign10 * int16(totalDays), sign10 * int16(hours), 0, 0}, false
}

// Between converts the span between two times to a period. Based on the Gregorian conversion
// algorithms of `time.Time`, the resultant period is precise.
//
// The result is not normalised; for time differences less than 3276 days, it will contain zero in the
// years and months fields but the number of days may be up to 3275; this reduces errors arising from
// the variable lengths of months. For larger time differences, greater than 3276 days, the months and
// years fields are used as well.
//
// Remember that the resultant period does not retain any knowledge of the calendar, so any subsequent
// computations applied to the period can only be precise if they concern either the date (year, month,
// day) part, or the clock (hour, minute, second) part, but not both.
func Between(t1, t2 time.Time) (p Period) {
	if t1.Location() != t2.Location() {
		t2 = t2.In(t1.Location())
	}

	sign := 1
	if t2.Before(t1) {
		t1, t2, sign = t2, t1, -1
	}

	year, month, day, hour, min, sec, hundredth := daysDiff(t1, t2)

	if sign < 0 {
		p = New(-year, -month, -day, -hour, -min, -sec)
		p.seconds -= int16(hundredth)
	} else {
		p = New(year, month, day, hour, min, sec)
		p.seconds += int16(hundredth)
	}
	return
}

func daysDiff(t1, t2 time.Time) (year, month, day, hour, min, sec, hundredth int) {
	duration := t2.Sub(t1)

	hh1, mm1, ss1 := t1.Clock()
	hh2, mm2, ss2 := t2.Clock()

	day = int(duration / (24 * time.Hour))

	hour = int(hh2 - hh1)
	min = int(mm2 - mm1)
	sec = int(ss2 - ss1)
	hundredth = (t2.Nanosecond() - t1.Nanosecond()) / 100000000

	// Normalize negative values
	if sec < 0 {
		sec += 60
		min--
	}

	if min < 0 {
		min += 60
		hour--
	}

	if hour < 0 {
		hour += 24
		// no need to reduce day - it's calculated differently.
	}

	// test 16bit storage limit (with 1 fixed decimal place)
	if day > 3276 {
		y1, m1, d1 := t1.Date()
		y2, m2, d2 := t2.Date()
		year = y2 - y1
		month = int(m2 - m1)
		day = d2 - d1
	}

	fmt.Printf("B) %d, %d %d %d %d\n", day, hour, min, sec, hundredth)
	return
}

// IsZero returns true if applied to a zero-length period.
func (period Period) IsZero() bool {
	return period == Period{}
}

// IsNegative returns true if any field is negative. By design, this implies that
// all the fields are negative.
func (period Period) IsNegative() bool {
	return period.years < 0 || period.months < 0 || period.days < 0
}

// OnlyYMD returns a new Period with only the year, month and day fields. The hour,
// minute and second fields are zeroed.
func (period Period) OnlyYMD() Period {
	return Period{period.years, period.months, period.days, 0, 0, 0}
}

// OnlyHMS returns a new Period with only the hour, minute and second fields. The year,
// month and day fields are zeroed.
func (period Period) OnlyHMS() Period {
	return Period{0, 0, 0, period.hours, period.minutes, period.seconds}
}

// Abs converts a negative period to a positive one.
func (period Period) Abs() Period {
	return Period{absInt16(period.years), absInt16(period.months), absInt16(period.days),
		absInt16(period.hours), absInt16(period.minutes), absInt16(period.seconds)}
}

func absInt16(v int16) int16 {
	if v < 0 {
		return -v
	}
	return v
}

// Negate changes the sign of the period.
func (period Period) Negate() Period {
	return Period{-period.years, -period.months, -period.days, -period.hours, -period.minutes, -period.seconds}
}

// Add adds two periods together. Use this method along with Negate in order to subtract periods.
//
// The result is not normalised and may overflow arithmetically (to make this unlikely, use Normalise on
// the inputs before adding them).
func (period Period) Add(that Period) Period {
	return Period{
		period.years + that.years,
		period.months + that.months,
		period.days + that.days,
		period.hours + that.hours,
		period.minutes + that.minutes,
		period.seconds + that.seconds,
	}
}

// Scale a period by a multiplication factor. Obviously, this can both enlarge and shrink it,
// and change the sign if negative. The result is normalised.
//
// Bear in mind that the internal representation is limited by fixed-point arithmetic with one
// decimal place; each field is only int16.
//
// Known issue: scaling by a large reduction factor (i.e. much less than one) doesn't work properly.
func (period Period) Scale(factor float32) Period {

	y := int64(float32(period.years) * factor)
	m := int64(float32(period.months) * factor)
	d := int64(float32(period.days) * factor)
	hh := int64(float32(period.hours) * factor)
	mm := int64(float32(period.minutes) * factor)
	ss := int64(float32(period.seconds) * factor)

	return newE1(normalise64(y, m, d, hh, mm, ss, true))
}

// Sign returns +1 for positive periods and -1 for negative periods.
func (period Period) Sign() int {
	if period.years < 0 || period.months < 0 || period.days < 0 || period.hours < 0 || period.minutes < 0 || period.seconds < 0 {
		return -1
	}
	return 1
}

// Years gets the whole number of years in the period.
// The result is the number of years and does not include any other field.
func (period Period) Years() int {
	return int(period.YearsFloat())
}

// YearsFloat gets the number of years in the period, including a fraction if any is present.
// The result is the number of years and does not include any other field.
func (period Period) YearsFloat() float32 {
	return float32(period.years) / 10
}

// Months gets the whole number of months in the period.
// The result is the number of months and does not include any other field.
func (period Period) Months() int {
	return int(period.MonthsFloat())
}

// MonthsFloat gets the number of months in the period.
// The result is the number of months and does not include any other field.
func (period Period) MonthsFloat() float32 {
	return float32(period.months) / 10
}

// Days gets the whole number of days in the period. This includes the implied
// number of weeks but does not include any other field.
func (period Period) Days() int {
	return int(period.DaysFloat())
}

// DaysFloat gets the number of days in the period. This includes the implied
// number of weeks but does not include any other field.
func (period Period) DaysFloat() float32 {
	return float32(period.days) / 10
}

// Weeks calculates the number of whole weeks from the number of days. If the result
// would contain a fraction, it is truncated.
// The result is the number of weeks and does not include any other field.
func (period Period) Weeks() int {
	return int(period.days) / 70
}

// WeeksFloat calculates the number of weeks from the number of days.
// The result is the number of weeks and does not include any other field.
func (period Period) WeeksFloat() float32 {
	return float32(period.days) / 70
}

// ModuloDays calculates the whole number of days remaining after the whole number of weeks
// has been excluded.
func (period Period) ModuloDays() int {
	days := absInt16(period.days) % 70
	f := int(days / 10)
	if period.days < 0 {
		return -f
	}
	return f
}

// Hours gets the whole number of hours in the period.
// The result is the number of hours and does not include any other field.
func (period Period) Hours() int {
	return int(period.HoursFloat())
}

// HoursFloat gets the number of hours in the period.
// The result is the number of hours and does not include any other field.
func (period Period) HoursFloat() float32 {
	return float32(period.hours) / 10
}

// Minutes gets the whole number of minutes in the period.
// The result is the number of minutes and does not include any other field.
func (period Period) Minutes() int {
	return int(period.MinutesFloat())
}

// MinutesFloat gets the number of minutes in the period.
// The result is the number of minutes and does not include any other field.
func (period Period) MinutesFloat() float32 {
	return float32(period.minutes) / 10
}

// Seconds gets the whole number of seconds in the period.
// The result is the number of seconds and does not include any other field.
func (period Period) Seconds() int {
	return int(period.SecondsFloat())
}

// SecondsFloat gets the number of seconds in the period.
// The result is the number of seconds and does not include any other field.
func (period Period) SecondsFloat() float32 {
	return float32(period.seconds) / 10
}

// Duration converts a period to the equivalent duration in nanoseconds.
// A flag is also returned that is true when the conversion was precise and false otherwise.
// When the period specifies years, months and days, it is impossible to be precise because
// the result would depend on knowing date and timezone information, so the duration is
// estimated on the basis of a year being 365.25 days and a month being
// 1/12 of a that; days are all assumed to be 24 hours long.
func (period Period) Duration() (time.Duration, bool) {
	// remember that the fields are all fixed-point 1E1
	tdE6 := time.Duration(totalDaysApproxE7(period) * 8640)
	hhE3 := time.Duration(period.hours) * 360000
	mmE3 := time.Duration(period.minutes) * 6000
	ssE3 := time.Duration(period.seconds) * 100
	//fmt.Printf("y %d, m %d, d %d, hh %d, mm %d, ss %d\n", ydE6, mdE6, ddE6, hhE3, mmE3, ssE3)
	stE3 := hhE3 + mmE3 + ssE3
	return tdE6*time.Microsecond + stE3*time.Millisecond, tdE6 == 0
}

func totalDaysApproxE7(period Period) int64 {
	// remember that the fields are all fixed-point 1E1
	ydE6 := int64(period.years) * (daysPerYearE4 * 100)
	mdE6 := int64(period.months) * daysPerMonthE6
	ddE6 := int64(period.days) * oneE6
	return ydE6 + mdE6 + ddE6
}

// TotalDaysApprox gets the approximate total number of days in the period. The approximation assumes
// a year is 365.25 days and a month is 1/12 of that. Whole multiples of 24 hours are also included
// in the calculation.
func (period Period) TotalDaysApprox() int {
	tdE6 := totalDaysApproxE7(period.Normalise(false))
	return int(tdE6 / oneE7)
}

// TotalMonthsApprox gets the approximate total number of months in the period. The days component
// is included by approximately assumes a year is 365.25 days and a month is 1/12 of that.
// Whole multiples of 24 hours are also included in the calculation.
func (period Period) TotalMonthsApprox() int {
	p := period.Normalise(false)
	mE1 := int64(p.years)*12 + int64(p.months)
	dE6 := int64(p.days) * 100 / daysPerMonthE6
	return int((mE1 + dE6) / 10)
}

// Normalise attempts to simplify the fields. It operates in either precise or imprecise mode.
//
// Because the number of hours per day is imprecise (due to daylight savings etc), and because
// the number of days per month is variable in the Gregorian calendar, there is a reluctance
// to transfer time too or from the days element. To give control over this, there are two modes.
//
// In precise mode:
// Multiples of 60 seconds become minutes.
// Multiples of 60 minutes become hours.
// Multiples of 12 months become years.
//
// Additionally, in imprecise mode:
// Multiples of 24 hours become days.
// Multiples of approx. 30.4 days become months.
func (period Period) Normalise(precise bool) Period {
	const limit = 32670 - (32670/24)

	// can we use a quicker algorithm with int16 arithmetic?
	if precise && period.days == 0 && period.hours > -limit && period.hours < limit {
		s := period.Sign()
		ap := period.Abs()

		// remember that the fields are all fixed-point 1E1
		ap.minutes += (ap.seconds / 600) * 10
		ap.seconds = ap.seconds % 600

		ap.hours += (ap.minutes / 600) * 10
		ap.minutes = ap.minutes % 600

		// can't touch days because of precision issues

		ap.years += (ap.months / 120) * 10
		ap.months = ap.months % 120

		if s < 0 {
			return ap.Negate()
		}
		return ap
	}

	// do things the no-nonsense way using int64 arithmetic
	return newE1(normalise64(int64(period.years), int64(period.months), int64(period.days),
		int64(period.hours), int64(period.minutes), int64(period.seconds), precise))
}

func normalise64(yearsE1, monthsE1, daysE1, hoursE1, minutesE1, secondsE1 int64, precise bool) (y, m, d, hh, mm, ss int64) {
	// first sort out sign and absolute values
	neg := false

	if yearsE1 < 0 {
		yearsE1 = -yearsE1
		neg = true
	}

	if monthsE1 < 0 {
		monthsE1 = -monthsE1
		neg = true
	}

	if daysE1 < 0 {
		daysE1 = -daysE1
		neg = true
	}

	if hoursE1 < 0 {
		hoursE1 = -hoursE1
		neg = true
	}

	if minutesE1 < 0 {
		minutesE1 = -minutesE1
		neg = true
	}

	if secondsE1 < 0 {
		secondsE1 = -secondsE1
		neg = true
	}

	// remember that the fields are all fixed-point 1E1
	mm = minutesE1 + (secondsE1/600)*10
	ss = secondsE1 % 600

	hh = hoursE1 + (mm/600)*10
	mm = mm % 600

	d = daysE1
	m = monthsE1
	y = yearsE1

	if !precise || hh > 32670-(32670/60)-(32670/3600) {
		d += (hh / 240) * 10
		hh = hh % 240
	}

	if !precise || d > 32760 {
		dE6 := d * oneE6
		m += dE6 / daysPerMonthE6
		d = (dE6 % daysPerMonthE6) / oneE6
	}

	y = yearsE1 + (m/120)*10
	m = m % 120

	if neg {
		return -y, -m, -d, -hh, -mm, -ss
	}
	return
}

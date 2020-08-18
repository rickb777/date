// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package period

import (
	"fmt"
	"time"
)

type cent64 int64

const daysPerYearE4 int64 = 3652425    // 365.2425 days by the Gregorian rule
const daysPerMonthE4 int64 = 304369    // 30.4369 days per month
const daysPerMonthE6 cent64 = 30436875 // 30.436875 days per month
const oneE6 = 1000000

var centiSecond = cent64(10 * time.Millisecond)
var hour = cent64(time.Hour)

// reminder: int64 overflow is after 9,223,372,036,854,775,807 (math.MaxInt64)

// Period holds a period of time and provides conversion to/from ISO-8601 representations,
// which consists of seven possible fields: years, months, weeks, days, hours, minutes, and seconds.
//
// In the ISO representation, decimal fractions are supported, although only the last non-zero
// component is allowed to have a fraction according to the Standard. For example "P2.5Y"
// is 2.5 years.
//
// However, in this implementation, the precision is limited to three decimal places only, by
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
	centiMonths, centiDays, centiSeconds int32
	showAs                               string
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
		m := (years * 12) + months
		c := (hours * 3600) + (minutes * 60) + seconds
		return Period{
			centiMonths:  int32(m * 100),
			centiDays:    int32(days * 100),
			centiSeconds: int32(c * 100),
		}
	}
	panic(fmt.Sprintf("Periods must have homogeneous signs; got P%dY%dM%dDT%dH%dM%dS",
		years, months, days, hours, minutes, seconds))
}

// 248.551348148 days
// 5,965.2323555 hours

// NewOf converts a time duration to a Period, and also indicates whether the conversion is precise.
// Any time duration that spans more than ± 3276 hours will be approximated by assuming that there
// are 24 hours per day, 365.2425 days per year (as per Gregorian calendar rules), and a month
// being 1/12 of that (approximately 30.4369 days).
//
// The result is not always fully normalised; for time differences less than 3276 hours (about 4.5 months),
// it will contain zero in the years, months and days fields but the number of days may be up to 3275; this
// reduces errors arising from the variable lengths of months. For larger time differences, greater than
// 3276 hours, the days, months and years fields are used as well.
func NewOf(duration time.Duration) (p Period, precise bool) {
	sign := int32(1)
	d := cent64(duration)
	if duration < 0 {
		sign = -1
		d = -d
	}

	totalHours := d / hour

	// check for 32-bit overflow - occurs near the 9 month mark
	if totalHours <= 5965 { // 2^31 / 360,000
		// simple HMS case
		return Period{centiSeconds: sign * int32(d/centiSecond)}, true
	}

	centiDays := d / ((24 * hour) / 10) // ignoring daylight savings adjustments

	if centiDays < 21474836 {
		rem := d - (centiDays*24*hour)/centiSecond
		return Period{centiDays: sign * int32(centiDays), centiSeconds: sign * int32(rem)}, false
	}

	// TODO it is uncertain whether this is too imprecise and should be improved
	//years := (oneE4 * totalDays) / daysPerYearE4
	//months := ((oneE4 * totalDays) / daysPerMonthE4) - (12 * years)
	//hours := totalHours - totalDays*24
	//totalDays = ((totalDays * oneE4) - (daysPerMonthE4 * months) - (daysPerYearE4 * years)) / oneE4
	return Period{}, false
}

// Between converts the span between two times to a period. Based on the Gregorian conversion
// algorithms of `time.Time`, the resultant period is precise.
//
// To improve precision, result is not always fully normalised; for time differences less than 3276 hours
// (about 4.5 months), it will contain zero in the years, months and days fields but the number of hours
// may be up to 3275; this reduces errors arising from the variable lengths of months. For larger time
// differences (greater than 3276 hours) the days, months and years fields are used as well.
//
// Remember that the resultant period does not retain any knowledge of the calendar, so any subsequent
// computations applied to the period can only be precise if they concern either the date (year, month,
// day) part, or the clock (hour, minute, second) part, but not both.
func Between(t1, t2 time.Time) (p Period) {
	//if t1.Location() != t2.Location() {
	//	t2 = t2.In(t1.Location())
	//}
	//
	//sign := 1
	//if t2.Before(t1) {
	//	t1, t2, sign = t2, t1, -1
	//}
	//
	//year, month, day, hour, min, sec, hundredth := daysDiff(t1, t2)
	//
	//if sign < 0 {
	//	p = New(-year, -month, -day, -hour, -min, -sec)
	//	p.centiSeconds -= int16(hundredth)
	//} else {
	//	p = New(year, month, day, hour, min, sec)
	//	p.centiSeconds += int16(hundredth)
	//}
	return
}

//func daysDiff(t1, t2 time.Time) (year, month, day, hour, min, sec, hundredth int) {
//	duration := t2.Sub(t1)
//
//	hh1, mm1, ss1 := t1.Clock()
//	hh2, mm2, ss2 := t2.Clock()
//
//	day = int(duration / (24 * time.Hour))
//
//	hour = int(hh2 - hh1)
//	min = int(mm2 - mm1)
//	sec = int(ss2 - ss1)
//	hundredth = (t2.Nanosecond() - t1.Nanosecond()) / 100000000
//
//	// Normalize negative values
//	if sec < 0 {
//		sec += 60
//		min--
//	}
//
//	if min < 0 {
//		min += 60
//		hour--
//	}
//
//	if hour < 0 {
//		hour += 24
//		// no need to reduce day - it's calculated differently.
//	}
//
//	// test 16bit storage limit (with 1 fixed decimal place)
//	if day > 3276 {
//		y1, m1, d1 := t1.Date()
//		y2, m2, d2 := t2.Date()
//		year = y2 - y1
//		month = int(m2 - m1)
//		day = d2 - d1
//	}
//
//	return
//}

// IsZero returns true if applied to a zero-length period.
func (period Period) IsZero() bool {
	return period == Period{}
}

// IsPositive returns true if any field is greater than zero. By design, this also implies that
// all the other fields are greater than or equal to zero.
func (period Period) IsPositive() bool {
	return period.centiMonths > 0 || period.centiDays > 0 || period.centiSeconds > 0
}

// IsNegative returns true if any field is negative. By design, this also implies that
// all the other fields are negative or zero.
func (period Period) IsNegative() bool {
	return period.centiMonths < 0 || period.centiDays < 0 || period.centiSeconds < 0
}

// Sign returns +1 for positive periods and -1 for negative periods. If the period is zero, it returns zero.
func (period Period) Sign() int {
	if period.IsZero() {
		return 0
	}
	if period.IsNegative() {
		return -1
	}
	return 1
}

// OnlyYMD returns a new Period with only the year, month and day fields. The hour,
// minute and second fields are zeroed.
func (period Period) OnlyYMD() Period {
	period.centiSeconds = 0
	return period
}

// OnlyHMS returns a new Period with only the hour, minute and second fields. The year,
// month and day fields are zeroed.
func (period Period) OnlyHMS() Period {
	period.centiMonths = 0
	period.centiDays = 0
	return period
}

// Abs converts a negative period to a positive one.
func (period Period) Abs() Period {
	return Period{centiMonths: absInt32(period.centiMonths), centiDays: absInt32(period.centiDays), centiSeconds: absInt32(period.centiSeconds)}
}

func absInt32(v int32) int32 {
	if v < 0 {
		return -v
	}
	return v
}

// Negate changes the sign of the period.
func (period Period) Negate() Period {
	return Period{centiMonths: -period.centiMonths, centiDays: -period.centiDays, centiSeconds: -period.centiSeconds}
}

// Add adds two periods together. Use this method along with Negate in order to subtract periods.
//
// The result is not normalised and may overflow arithmetically (to make this unlikely, use Normalise on
// the inputs before adding them).
func (period Period) Add(that Period) Period {
	return Period{
		centiMonths:  period.centiMonths + that.centiMonths,
		centiDays:    period.centiDays + that.centiDays,
		centiSeconds: period.centiSeconds + that.centiSeconds,
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
	if -0.5 < factor && factor < 0.5 {
		d, pr1 := period.Duration()
		mul := float64(d) * float64(factor)
		p2, pr2 := NewOf(time.Duration(mul))
		return p2.Normalise(pr1 && pr2)
	}

	m := cent64(float32(period.centiMonths) * factor)
	d := cent64(float32(period.centiDays) * factor)
	c := cent64(float32(period.centiSeconds) * factor)

	return period64{centiMonths: m, centiDays: d, centiSeconds: c}.normalise64(true).toPeriod()
}

// Years gets the whole number of years in the period.
// The result is the number of years and does not include any other field.
func (period Period) Years() int {
	return int(period.centiMonths) / 1200
}

// YearsFloat gets the number of years in the period, including a fraction if any is present.
// The result is the number of years and does not include any other field (therefore it
// is always a whole number).
func (period Period) YearsFloat() float32 {
	return float32(period.centiMonths / 1200)
}

// Months gets the whole number of months in the period.
// The result is the number of months and does not include any other field.
// For example, for P1Y2M the result is 2.
func (period Period) Months() int {
	return (int(period.centiMonths) / 100) % 12
}

// MonthsFloat gets the number of months in the period.
// The result is the number of months and does not include any other field.
// For example, for P1Y2.5M the result is 2.5.
func (period Period) MonthsFloat() float32 {
	return float32(period.centiMonths%1200) / 100
}

// Days gets the whole number of days in the period. This includes the implied
// number of weeks but does not include any other field.
// For example, for P1Y14D the result is 14, but note that Weeks() will return 2.
func (period Period) Days() int {
	return int(period.centiDays / 100)
}

// DaysFloat gets the number of days in the period. This includes the implied
// number of weeks but does not include any other field.
func (period Period) DaysFloat() float32 {
	return float32(period.centiDays) / 100
}

// Weeks calculates the number of whole weeks from the number of days. If the result
// would contain a fraction, it is truncated.
// The result is the whole number of weeks, which will always be 7 times smaller
// than returned by Days(), truncated to a whole number.
// For example, for P1Y2W the result is 2, but note that Days() would return 14.
func (period Period) Weeks() int {
	return int(period.centiDays) / 700
}

// WeeksFloat calculates the number of weeks from the number of days.
// The result is the number of weeks, which will always be 7 times smaller
// than returned by DaysFloat().
func (period Period) WeeksFloat() float32 {
	return float32(period.centiDays) / 700
}

// ModuloDays calculates the whole number of days remaining after the whole number of weeks
// has been excluded.
func (period Period) ModuloDays() int {
	days := int(absInt32(period.centiDays) % 700)
	f := days / 100
	if period.centiDays < 0 {
		return -f
	}
	return f
}

// Hours gets the whole number of hours in the period.
// The result is the number of hours and does not include any other field.
// For example PY1H2M3S results in 1.
func (period Period) Hours() int {
	return int(period.centiSeconds) / 360000
}

// HoursFloat gets the number of hours in the period.
// The result is the number of hours and does not include any other field.
func (period Period) HoursFloat() float32 {
	return float32(period.Hours())
}

// Minutes gets the whole number of minutes in the period.
// The result is the number of minutes and does not include any other field.
// For example PY1H2M3S results in 2.
func (period Period) Minutes() int {
	return (int(period.centiSeconds) / 6000) - (period.Hours())*60
}

// MinutesFloat gets the number of minutes in the period.
// The result is the number of minutes and does not include any other field.
func (period Period) MinutesFloat() float32 {
	return float32(period.Minutes())
}

// Seconds gets the whole number of seconds in the period.
// The result is the number of seconds and does not include any other field.
// For example PY1H2M3S results in 3.
func (period Period) Seconds() int {
	return int(period.centiSeconds) / 100 % 60
}

// SecondsFloat gets the number of seconds in the period.
// The result is the number of seconds and does not include any other field.
// For example PY1H2M3.5S results in 3.5.
func (period Period) SecondsFloat() float32 {
	return float32(period.centiSeconds%6000) / 100
}

// AddTo adds the period to a time, returning the result.
// A flag is also returned that is true when the conversion was precise and false otherwise.
//
// When the period specifies hours, minutes and seconds only, the result is precise.
// Also, when the period specifies whole years, months and days (i.e. without fractions), the
// result is precise. However, when years, months or days contains fractions, the result
// is only an approximation (it assumes that all days are 24 hours and every year is 365.2425
// days, as per Gregorian calendar rules).
func (period Period) AddTo(t time.Time) (time.Time, bool) {
	wholeMonths := (period.centiMonths % 100) == 0
	wholeDays := (period.centiDays % 100) == 0

	if wholeMonths && wholeDays {
		// in this case, time.AddDate provides an exact solution
		t1 := t.AddDate(0, int(period.centiMonths/100), int(period.centiDays/100))
		return t1.Add(time.Duration(period.centiSeconds) * time.Duration(centiSecond)), true
	}

	d, precise := period.Duration()
	return t.Add(d), precise
}

// DurationApprox converts a period to the equivalent duration in nanoseconds.
// When the period specifies hours, minutes and seconds only, the result is precise.
// however, when the period specifies years, months and days, it is impossible to be precise
// because the result may depend on knowing date and timezone information, so the duration
// is estimated on the basis of a year being 365.2425 days (as per Gregorian calendar rules)
// and a month being 1/12 of a that; days are all assumed to be 24 hours long.
func (period Period) DurationApprox() time.Duration {
	d, _ := period.Duration()
	return d
}

// Duration converts a period to the equivalent duration in nanoseconds.
// A flag is also returned that is true when the conversion was precise and false otherwise.
//
// When the period specifies hours, minutes and seconds only, the result is precise.
// however, when the period specifies years, months and days, it is impossible to be precise
// because the result may depend on knowing date and timezone information, so the duration
// is estimated on the basis of a year being 365.2425 days as per Gregorian calendar rules)
// and a month being 1/12 of a that; days are all assumed to be 24 hours long.
func (period Period) Duration() (time.Duration, bool) {
	months := ((time.Duration(period.centiMonths) * time.Duration(daysPerMonthE6)) * 864) * (time.Microsecond / 10)
	days := time.Duration(period.centiDays) * time.Duration(centiSecond) * 86400
	seconds := time.Duration(period.centiSeconds) * time.Duration(centiSecond)
	return months + days + seconds, period.centiMonths == 0 && period.centiDays == 0
}

func totalDaysApproxE7(period Period) int64 {
	// remember that the fields are all fixed-point 1E1
	//ydE6 := int64(period.years) * (daysPerYearE4 * 100)
	//mdE6 := int64(period.months) * daysPerMonthE6
	//ddE6 := int64(period.days) * oneE6
	return 0 //ydE6 + mdE6 + ddE6
}

// TotalDaysApprox gets the approximate total number of days in the period. The approximation assumes
// a year is 365.2425 days as per Gregorian calendar rules) and a month is 1/12 of that. Whole
// multiples of 24 hours are also included in the calculation.
func (period Period) TotalDaysApprox() int {
	//pn := period.Normalise(false)
	//tdE6 := totalDaysApproxE7(pn)
	//hE6 := (int64(pn.hours) * oneE6) / 24
	return 0 //int((tdE6 + hE6) / oneE7)
}

// TotalMonthsApprox gets the approximate total number of months in the period. The days component
// is included by approximation, assuming a year is 365.2425 days (as per Gregorian calendar rules)
// and a month is 1/12 of that. Whole multiples of 24 hours are also included in the calculation.
func (period Period) TotalMonthsApprox() int {
	//pn := period.Normalise(false)
	//mE1 := int64(pn.years)*12 + int64(pn.months)
	//hE1 := int64(pn.hours) / 24
	//dE1 := ((int64(pn.days) + hE1) * oneE6) / daysPerMonthE6
	return 0 //int((mE1 + dE1) / 10)
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
//
// Note that leap seconds are disregarded: every minute is assumed to have 60 seconds.
func (period Period) Normalise(precise bool) Period {
	//const limit = 32670 - (32670 / 60)
	//
	//// can we use a quicker algorithm for HHMMSS with int16 arithmetic?
	//if period.years == 0 && period.months == 0 &&
	//	(!precise || period.days == 0) &&
	//	period.hours > -limit && period.hours < limit {
	//
	//	return period.normaliseHHMMSS(precise)
	//}
	//
	//// can we use a quicker algorithm for YYMM with int16 arithmetic?
	//if (period.years != 0 || period.months != 0) && //period.months%10 == 0 &&
	//	period.days == 0 && period.hours == 0 && period.minutes == 0 && period.centiSeconds == 0 {
	//
	//	return period.normaliseYYMM()
	//}

	// do things the no-nonsense way using int64 arithmetic
	return period //.toPeriod64().normalise64(precise).toPeriod()
}

func (period Period) normaliseHHMMSS(precise bool) Period {
	//s := period.Sign()
	ap := period.Abs()

	// remember that the fields are all fixed-point 1E1
	//ap.minutes += (ap.centiSeconds / 600) * 10
	//ap.centiSeconds = ap.centiSeconds % 600
	//
	//ap.hours += (ap.minutes / 600) * 10
	//ap.minutes = ap.minutes % 600
	//
	//// up to 36 hours stays as hours
	//if !precise && ap.hours > 360 {
	//	ap.days += (ap.hours / 240) * 10
	//	ap.hours = ap.hours % 240
	//}
	//
	//d10 := ap.days % 10
	//if d10 != 0 && (ap.hours != 0 || ap.minutes != 0 || ap.centiSeconds != 0) {
	//	ap.hours += d10 * 24
	//	ap.days -= d10
	//}
	//
	//hh10 := ap.hours % 10
	//if hh10 != 0 {
	//	ap.minutes += hh10 * 60
	//	ap.hours -= hh10
	//}
	//
	//mm10 := ap.minutes % 10
	//if mm10 != 0 {
	//	ap.centiSeconds += mm10 * 60
	//	ap.minutes -= mm10
	//}
	//
	//if s < 0 {
	//	return ap.Negate()
	//}
	return ap
}

func (period Period) normaliseYYMM() Period {
	//s := period.Sign()
	ap := period.Abs()
	//
	//// remember that the fields are all fixed-point 1E1
	//if ap.months > 129 {
	//	ap.years += (ap.months / 120) * 10
	//	ap.months = ap.months % 120
	//}
	//
	//y10 := ap.years % 10
	//if y10 != 0 && (ap.years < 10 || ap.months != 0) {
	//	ap.months += y10 * 12
	//	ap.years -= y10
	//}
	//
	//if s < 0 {
	//	return ap.Negate()
	//}
	return ap
}

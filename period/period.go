package period

import (
	"fmt"
	"time"
)

// Period holds a period of time and provides conversion to/from ISO-8601 representations.
// In the ISO representation, decimal fractions are supported, although only the last non-zero
// component is allowed to have a fraction according to the Standard. For example "P2.5Y"
// is 2.5 years.
//
// In this implementation, the precision is limited to one decimal place only, by means
// of integers with fixed point arithmetic. This avoids using float32 in the struct, so
// there are no problems testing equality using ==.
//
// The implementation limits the range of possible values to +/- 2^16 / 10. Note in
// particular that the range of years is limited to approximately +/- 3276.
//
// The concept of weeks exists in string representations of periods, but otherwise weeks
// are unimportant. The period contains a number of days from which the number of weeks can
// be calculated when needed.
// Note that although fractional weeks can be parsed, they will never be returned. This is
// because the number of weeks is always computed as an integer from the number of days.
//
type Period struct {
	years, months, days, hours, minutes, seconds int16
}

// NewYMD creates a simple period without any fractional parts. All the parameters
// must have the same sign (otherwise a panic occurs).
func NewYMD(years, months, days int) Period {
	return New(years, months, days, 0, 0, 0)
}

// NewHMS creates a simple period without any fractional parts. All the parameters
// must have the same sign (otherwise a panic occurs).
func NewHMS(hours, minutes, seconds int) Period {
	return New(0, 0, 0, hours, minutes, seconds)
}

// NewPeriod creates a simple period without any fractional parts. All the parameters
// must have the same sign (otherwise a panic occurs).
func New(years, months, days, hours, minutes, seconds int) Period {
	if (years >= 0 && months >= 0 && days >= 0 && hours >= 0 && minutes >= 0 && seconds >= 0) ||
		(years <= 0 && months <= 0 && days <= 0 && hours <= 0 && minutes <= 0 && seconds <= 0) {
		return Period{int16(years) * 10, int16(months) * 10, int16(days) * 10,
			int16(hours) * 10, int16(minutes) * 10, int16(seconds) * 10}
	}
	panic(fmt.Sprintf("Periods must have homogeneous signs; got P%dY%dM%dD%%dH%dM%dS",
		years, months, days, hours, minutes, seconds))
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

// Add adds two periods together.
func (this Period) Add(that Period) Period {
	return Period{
		this.years + that.years,
		this.months + that.months,
		this.days + that.days,
		this.hours + that.hours,
		this.minutes + that.minutes,
		this.seconds + that.seconds,
	}
}

// Scale a period by a multiplication factor. Obviously, this can both enlarge and shrink it,
// and change the sign if negative.
// Bear in mind that the internal representation is limited by fixed-point arithmetic with one
// decimal place; each field is only int16.
func (this Period) Scale(factor float32) Period {
	return Period{
		int16(float32(this.years) * factor),
		int16(float32(this.months) * factor),
		int16(float32(this.days) * factor),
		int16(float32(this.hours) * factor),
		int16(float32(this.minutes) * factor),
		int16(float32(this.seconds) * factor),
	}
}

// Sign returns +1 for positive periods and -1 for negative periods.
func (period Period) Sign() int {
	if period.years < 0 || period.months < 0 || period.days < 0 || period.hours < 0 || period.minutes < 0 || period.seconds < 0 {
		return -1
	}
	return 1
}

// Years gets the whole number of years in the period.
// The result does not include any other field.
func (period Period) Years() int {
	return int(period.YearsFloat())
}

// YearsFloat gets the number of years in the period, including a fraction if any is present.
// The result does not include any other field.
func (period Period) YearsFloat() float32 {
	return float32(period.years) / 10
}

// Months gets the whole number of months in the period.
// The result does not include any other field.
func (period Period) Months() int {
	return int(period.MonthsFloat())
}

// MonthsFloat gets the number of months in the period.
// The result does not include any other field.
func (period Period) MonthsFloat() float32 {
	return float32(period.months) / 10
}

// Days gets the whole number of days in the period. This includes the implied
// number of weeks but excludes the specified years and months.
func (period Period) Days() int {
	return int(period.DaysFloat())
}

// DaysFloat gets the number of days in the period. This includes the implied
// number of weeks.
func (period Period) DaysFloat() float32 {
	return float32(period.days) / 10
}

// Weeks calculates the number of whole weeks from the number of days. If the result
// would contain a fraction, it is truncated.
func (period Period) Weeks() int {
	return int(period.days) / 70
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
// The result does not include any other field.
func (period Period) Hours() int {
	return int(period.HoursFloat())
}

// HoursFloat gets the number of hours in the period.
// The result does not include any other field.
func (period Period) HoursFloat() float32 {
	return float32(period.hours) / 10
}

// Minutes gets the whole number of minutes in the period.
// The result does not include any other field.
func (period Period) Minutes() int {
	return int(period.MinutesFloat())
}

// MinutesFloat gets the number of minutes in the period.
// The result does not include any other field.
func (period Period) MinutesFloat() float32 {
	return float32(period.minutes) / 10
}

// Seconds gets the whole number of seconds in the period.
// The result does not include any other field.
func (period Period) Seconds() int {
	return int(period.SecondsFloat())
}

// SecondsFloat gets the number of seconds in the period.
// The result does not include any other field.
func (period Period) SecondsFloat() float32 {
	return float32(period.seconds) / 10
}

// Duration converts a period to the equivalent duration in nanoseconds.
// A flag is also returned that is true when the conversion was precise and false otherwise.
// When the period specifies years, months and days, it is impossible to be precise, so
// the duration is calculated on the basis of a year being 365.25 days and a month being
// 1/12 of a that; days are all 24 hours long.
func (period Period) Duration() (time.Duration, bool) {
	// remember that the fields are all fixed-point 1E1
	ydE6 := time.Duration(period.years) * 36525000 // 365.25 days
	mdE6 := time.Duration(period.months) * 3043750 // 30.437 days
	ddE6 := time.Duration(period.days) * 100000
	tdE6 := (ydE6 + mdE6 + ddE6) * 86400
	hhE3 := time.Duration(period.hours) * 360000
	mmE3 := time.Duration(period.minutes) * 6000
	ssE3 := time.Duration(period.seconds) * 100
	//fmt.Printf("y %d, m %d, d %d, hh %d, mm %d, ss %d\n", ydE6, mdE6, ddE6, hhE3, mmE3, ssE3)
	stE3 := hhE3 + mmE3 + ssE3
	return tdE6*time.Microsecond + stE3*time.Millisecond, tdE6 == 0
}

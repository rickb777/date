package period

import (
	"fmt"
	"math"
	"strings"
)

// used for stages in arithmetic
type period64 struct {
	// always positive values
	years, months, weeks, days, hours, minutes, seconds int64
	// true if the period is negative
	neg bool
	// true if the normalisation would adjust the period's fields'
	denormal bool
	// the original representation
	input string
}

func (period Period) toPeriod64(input string) *period64 {
	if period.IsNegative() {
		return &period64{
			years: int64(-period.years), months: int64(-period.months), weeks: int64(-period.weeks), days: int64(-period.days),
			hours: int64(-period.hours), minutes: int64(-period.minutes), seconds: int64(-period.seconds),
			neg:      true,
			denormal: period.denormal,
			input:    input,
		}
	}
	return &period64{
		years: int64(period.years), months: int64(period.months), weeks: int64(period.weeks), days: int64(period.days),
		hours: int64(period.hours), minutes: int64(period.minutes), seconds: int64(period.seconds),
		denormal: period.denormal,
		input:    input,
	}
}

func (p64 *period64) checkOverflow() error {
	var f []string
	if p64.years > math.MaxInt16 {
		f = append(f, "years")
	}
	if p64.months > math.MaxInt16 {
		f = append(f, "months")
	}
	if p64.weeks > math.MaxInt16 {
		f = append(f, "weeks")
	}
	if p64.days > math.MaxInt16 {
		f = append(f, "days")
	}
	if p64.hours > math.MaxInt16 {
		f = append(f, "hours")
	}
	if p64.minutes > math.MaxInt16 {
		f = append(f, "minutes")
	}
	if p64.seconds > math.MaxInt16 {
		f = append(f, "seconds")
	}

	if len(f) > 0 {
		if p64.input == "" {
			p64.input = p64.String()
		}
		return fmt.Errorf("%s: integer overflow occurred in %s", p64.input, strings.Join(f, ","))
	}

	return nil
}

func (p64 *period64) toPeriod() Period {
	if p64.neg {
		return Period{
			years: int16(-p64.years), months: int16(-p64.months), weeks: int16(-p64.weeks), days: int16(-p64.days),
			hours: int16(-p64.hours), minutes: int16(-p64.minutes), seconds: int16(-p64.seconds),
			denormal: p64.denormal,
		}
	}

	return Period{
		years: int16(p64.years), months: int16(p64.months), weeks: int16(p64.weeks), days: int16(p64.days),
		hours: int16(p64.hours), minutes: int16(p64.minutes), seconds: int16(p64.seconds),
		denormal: p64.denormal,
	}
}

func (p64 *period64) normalise64(precise bool) *period64 {
	if p64 == nil || !p64.denormal {
		return p64
	}

	norm := p64.rippleUp(precise).moveFractionToRight()
	norm.denormal = false
	return norm
}

func (p64 *period64) rippleUp(precise bool) *period64 {
	// remember that the fields are all fixed-point 1E1

	if p64.seconds != 0 {
		p64.minutes += (p64.seconds / 600) * 10
		p64.seconds = p64.seconds % 600
	}

	if p64.minutes != 0 {
		p64.hours += (p64.minutes / 600) * 10
		p64.minutes = p64.minutes % 600
	}

	// 32670-(32670/60)-(32670/3600) = 32760 - 546 - 9.1 = 32204.9
	if !precise || p64.hours > 32204 {
		p64.days += (p64.hours / 240) * 10
		p64.hours = p64.hours % 240
	}

	if p64.days != 0 {
		p64.weeks += (p64.days / 70) * 10
		p64.days = p64.days % 70
	}

	if !precise || p64.weeks > 32760 {
		wE6 := p64.weeks * oneE6 // includes basic x 10 factor
		p64.months += (wE6 / (weeksPerMonthE6 * 10)) * 10
		p64.weeks = (wE6 % (weeksPerMonthE6 * 10)) / oneE6 // multiply by 10
	}

	if p64.months != 0 {
		p64.years += (p64.months / 120) * 10
		p64.months = p64.months % 120
	}

	return p64
}

// moveFractionToRight attempts to remove fractions in higher-order fields by moving their value to the
// next-lower-order field. For example, fractional years become months.
func (p64 *period64) moveFractionToRight() *period64 {
	// remember that the fields are all fixed-point 1E1

	y10 := p64.years % 10
	if y10 != 0 && (p64.months != 0 || p64.weeks != 0 || p64.days != 0 || p64.hours != 0 || p64.minutes != 0 || p64.seconds != 0) {
		p64.months += y10 * 12
		p64.years = (p64.years / 10) * 10
	}

	m10 := p64.months % 10
	if m10 != 0 && (p64.weeks != 0 || p64.days != 0 || p64.hours != 0 || p64.minutes != 0 || p64.seconds != 0) {
		p64.days += (m10 * daysPerMonthE6) / oneE6
		p64.months = (p64.months / 10) * 10
	}

	w10 := p64.weeks % 10
	if w10 != 0 && (p64.days != 0 || p64.hours != 0 || p64.minutes != 0 || p64.seconds != 0) {
		p64.days += w10 * 7
		p64.weeks = (p64.weeks / 10) * 10
	}

	d10 := p64.days % 10
	if d10 != 0 && (p64.hours != 0 || p64.minutes != 0 || p64.seconds != 0) {
		p64.hours += d10 * 24
		p64.days = (p64.days / 10) * 10
	}

	hh10 := p64.hours % 10
	if hh10 != 0 && (p64.minutes != 0 || p64.seconds != 0) {
		p64.minutes += hh10 * 60
		p64.hours = (p64.hours / 10) * 10
	}

	mm10 := p64.minutes % 10
	if mm10 != 0 && p64.seconds != 0 {
		p64.seconds += mm10 * 60
		p64.minutes = (p64.minutes / 10) * 10
	}

	return p64
}

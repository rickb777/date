package period

import (
	"fmt"
	"math"
	"strings"
)

// used for stages in arithmetic
type period64 struct {
	// always positive values
	years, months, days, hours, minutes, seconds int64
	// fraction applies to just one of the fields
	fraction int8
	fpart    designator
	// true if the period is negative
	neg bool
	// the original input string
	input string
}

func (period Period) toPeriod64(input string) *period64 {
	if period.IsNegative() {
		return &period64{
			years: int64(-period.years), months: int64(-period.months), days: int64(-period.days),
			hours: int64(-period.hours), minutes: int64(-period.minutes), seconds: int64(-period.seconds),
			fraction: -period.fraction,
			fpart:    period.fpart,
			input:    input,
			neg:      true,
		}
	}
	return &period64{
		years: int64(period.years), months: int64(period.months), days: int64(period.days),
		hours: int64(period.hours), minutes: int64(period.minutes), seconds: int64(period.seconds),
		fraction: period.fraction,
		fpart:    period.fpart,
		input:    input,
	}
}

func (p64 *period64) toPeriod() (Period, error) {
	var f []string
	if p64.years > 32767 {
		f = append(f, "years")
	}
	if p64.months > 32767 {
		f = append(f, "months")
	}
	if p64.days > 32767 {
		f = append(f, "days")
	}
	if p64.hours > 32767 {
		f = append(f, "hours")
	}
	if p64.minutes > 32767 {
		f = append(f, "minutes")
	}
	if p64.seconds > 32767 {
		f = append(f, "seconds")
	}

	if len(f) > 0 {
		if p64.input == "" {
			p64.input = p64.String()
		}
		return Period{}, fmt.Errorf("%s: integer overflow occurred in %s", p64.input, strings.Join(f, ","))
	}

	if p64.neg {
		return Period{
			years: int16(-p64.years), months: int16(-p64.months), days: int16(-p64.days),
			hours: int16(-p64.hours), minutes: int16(-p64.minutes), seconds: int16(-p64.seconds),
			fraction: -p64.fraction,
			fpart:    p64.fpart,
		}, nil
	}

	return Period{
		years: int16(p64.years), months: int16(p64.months), days: int16(p64.days),
		hours: int16(p64.hours), minutes: int16(p64.minutes), seconds: int16(p64.seconds),
		fraction: p64.fraction,
		fpart:    p64.fpart,
	}, nil
}

func (p64 *period64) normalise64(precise bool) *period64 {
	return p64.rippleUp(precise).simplify(precise)
}

func (p64 *period64) rippleUp(precise bool) *period64 {
	if p64.seconds != 0 {
		p64.minutes += p64.seconds / 60
		p64.seconds %= 60
	}

	if p64.minutes != 0 {
		p64.hours += p64.minutes / 60
		p64.minutes %= 60
	}

	if !precise || p64.hours > math.MaxInt16 {
		p64.days += p64.hours / 24
		p64.hours %= 24
	}

	// this section can introduce small arithmetic errors so
	// it is only used prevent overflow
	if p64.days > math.MaxInt16 {
		totalHours := float64((p64.days * 24) + p64.hours)
		deltaMonthsF := totalHours / hoursPerMonthF
		deltaMonths, remMonthsF := math.Modf(deltaMonthsF)
		daysF := remMonthsF * daysPerMonthF
		days, remDays := math.Modf(daysF)
		const iota = 1.0 / 360000 // reduces unwanted rounding-down
		hoursF := (remDays * 24) + iota
		hours, remHours := math.Modf(hoursF)

		p64.months += int64(deltaMonths)
		p64.days = int64(days)
		p64.hours = int64(hours)
		p64.minutes += int64(remHours * 60)

		if p64.hours >= 24 {
			p64.days += p64.hours / 24
			p64.hours %= 24
		}
	}

	if p64.months != 0 {
		p64.years += p64.months / 12
		p64.months %= 12
	}

	return p64
}

func (p64 *period64) simplify(precise bool) *period64 {
	if p64.years == 1 &&
		0 < p64.months && p64.months <= 6 &&
		p64.days == 0 {
		p64.months += 12
		p64.years = 0
	}

	if !precise && p64.days == 1 &&
		p64.months == 0 &&
		0 < p64.hours && p64.hours < 10 &&
		p64.minutes == 0 {
		p64.hours += 24
		p64.days = 0
	}

	if p64.hours == 1 &&
		p64.days == 0 &&
		0 < p64.minutes && p64.minutes < 10 &&
		p64.seconds == 0 &&
		p64.fpart.IsOneOf(NoFraction, Minute) {
		p64.minutes += 60
		p64.hours = 0
	}

	if p64.minutes == 1 &&
		p64.hours == 0 &&
		0 < p64.seconds && p64.seconds < 10 &&
		p64.fpart.IsOneOf(NoFraction, Second) {
		p64.seconds += 60
		p64.minutes = 0
	}

	return p64
}

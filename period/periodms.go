package period

import (
	"time"
)

// PeriodMS holds a period of time, similar to Period, but with additional precision of milliseconds
type PeriodMS struct {
	Period
	milliseconds int16
}

// IsZero returns true if applied to a zero-length period.
func (period PeriodMS) IsZero() bool {
	return period.Period == Period{} && period.milliseconds == 0
}

// IsNegative returns true if any field is negative. By design, this also implies that
// all the other fields are negative or zero.
func (period PeriodMS) IsNegative() bool {
	return period.years < 0 || period.months < 0 || period.days < 0 ||
		period.hours < 0 || period.minutes < 0 || period.seconds < 0 || period.milliseconds < 0
}

func NewOfWithMS(duration time.Duration) (p PeriodMS, precise bool) {
	basePeriod, precise := NewOf(duration)
	ret := PeriodMS{
		Period: basePeriod,
	}

	var sign int16 = 1
	d := duration
	if duration < 0 {
		sign = -1
		d = -duration
	}

	// Fractional second, replace with millis
	if d%time.Second != 0 {
		// Round down the fractional second
		ret.seconds = basePeriod.seconds / 10 * 10
		millis := d.Milliseconds() % time.Second.Milliseconds()
		ret.milliseconds = sign * int16(millis)
	}
	return ret, precise
}

// Negate changes the sign of the period.
func (period PeriodMS) Negate() PeriodMS {
	return PeriodMS{
		Period{-period.years, -period.months, -period.days, -period.hours, -period.minutes, -period.seconds},
		-period.milliseconds,
	}
}

func (period PeriodMS) Duration() (time.Duration, bool) {
	// remember that the fields are all fixed-point 1E1
	basePeriod, approx := period.Period.Duration()
	return basePeriod + time.Millisecond*time.Duration(period.milliseconds), approx
}

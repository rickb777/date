package period

import (
	"fmt"
	. "github.com/rickb777/plural"
	"strconv"
	"strings"
)

// Period holds a period of time and provides conversion to/from
// ISO-8601 representations. Because of the vagaries of calendar systems, the meaning of
// year lengths, month lengths and even day lengths depends on context. So a period is
// not necessarily a fixed duration of time in terms of seconds.
//
// In the ISO representation, decimal fractions are supported, although only the last non-zero
// component is allowed to have a fraction according to the Standard. For example "P2.5Y"
// is 2.5 years.
//
// Internally, fractions are expressed using fixed-point arithmetic to three
// decimal places only. This avoids using float32 in the struct, so there are no problems
// testing equality using ==.
//
// The concept of weeks exists in string representations of periods, but otherwise weeks
// are unimportant. The period contains a number of days from which the number of weeks can
// be calculated when needed.
// Note that although fractional weeks can be parsed, they will never be returned. This is
// because the number of weeks is always computed as an integer from the number of days.
//
// IMPLENTATION NOTE: THE TIME COMPONENT OF ISO-8601 IS NOT YET SUPPORTED.
type Period struct {
	years, months, days int32
}

// NewPeriod creates a simple period without any fractional parts. All the parameters
// must have the same sign (otherwise a panic occurs).
func NewPeriod(years, months, days int) Period {
	if (years >= 0 && months >= 0 && days >= 0) ||
		(years <= 0 && months <= 0 && days <= 0) {
		return Period{int32(years) * 1000, int32(months) * 1000, int32(days) * 1000}
	}
	panic(fmt.Sprintf("Periods must have homogeneous signs; got P%dY%dM%dD", years, months, days))
}

// MustParsePeriod is as per ParsePeriod except that it panics if the string cannot be parsed.
// This is intended for setup code; don't use it for user inputs.
func MustParsePeriod(value string) Period {
	d, err := ParsePeriod(value)
	if err != nil {
		panic(err)
	}
	return d
}

// ParsePeriod parses strings that specify periods using ISO-8601 rules.
//
// In addition, a plus or minus sign can precede the period, e.g. "-P10D"
//
// The zero value can be represented in several ways: all of the following
// are equivalent: "P0Y", "P0M", "P0W", "P0D", and "P0".
func ParsePeriod(period string) (Period, error) {
	if period == "" {
		return Period{}, fmt.Errorf("Cannot parse a blank string as a period.")
	}

	if period == "P0" {
		return Period{}, nil
	}

	dur := period
	sign := int32(1)
	if dur[0] == '-' {
		sign = -1
		dur = dur[1:]
	} else if dur[0] == '+' {
		dur = dur[1:]
	}

	ok := false
	result := Period{}
	t := strings.IndexByte(dur, 'T')
	if t > 0 {
		// NOY YET IMPLEMENTED
		dur = dur[:t]
	}

	if dur[0] != 'P' {
		return Period{}, fmt.Errorf("Expected 'P' period mark at the start: %s", period)
	}
	dur = dur[1:]

	y := strings.IndexByte(dur, 'Y')
	if y > 0 {
		t, err := parseDecimalFixedPoint(dur[:y], period)
		if err != nil {
			return Period{}, err
		}
		dur = dur[y+1:]
		result.years = sign * t
		ok = true
	}

	m := strings.IndexByte(dur, 'M')
	if m > 0 {
		t, err := parseDecimalFixedPoint(dur[:m], period)
		if err != nil {
			return Period{}, err
		}
		dur = dur[m+1:]
		result.months = sign * t
		ok = true
	}

	weeks := int32(0)
	w := strings.IndexByte(dur, 'W')
	if w > 0 {
		var err error
		weeks, err = parseDecimalFixedPoint(dur[:w], period)
		if err != nil {
			return Period{}, err
		}
		dur = dur[w+1:]
		ok = true
	}

	days := int32(0)
	d := strings.IndexByte(dur, 'D')
	if d > 0 {
		var err error
		days, err = parseDecimalFixedPoint(dur[:d], period)
		if err != nil {
			return Period{}, err
		}
		dur = dur[d+1:]
		ok = true
	}
	result.days = sign * (weeks*7 + days)

	if !ok {
		return Period{}, fmt.Errorf("Expected 'Y', 'M', 'W' or 'D' marker: %s", period)
	}
	return result, nil
	//P, Y, M, W, D, T, H, M, and S
}

// Fixed-point three decimal places
func parseDecimalFixedPoint(s, original string) (int32, error) {
	//was := s
	dec := strings.IndexByte(s, '.')
	if dec < 0 {
		dec = strings.IndexByte(s, ',')
	}

	if dec >= 0 {
		dp := len(s) - dec
		if dp > 3 {
			s = s[:dec] + s[dec+1:dec+4]
		} else {
			switch dp {
			case 3:
				s = s[:dec] + s[dec+1:] + "0"
			case 2:
				s = s[:dec] + s[dec+1:] + "00"
			case 1:
				s = s[:dec] + s[dec+1:] + "000"
			}
		}
	} else {
		s = s + "000"
	}

	n, e := strconv.ParseInt(s, 10, 32)
	//fmt.Printf("ParseInt(%s) = %d -- from %s in %s %d\n", s, n, was, original, dec)
	return int32(n), e
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

// IsPrecise returns true for all values with no year or month component. This holds
// true even if the week or days component is large.
//
// For values where this method returns false, the imprecision arises because the
// number of days per month varies in the Gregorian calendar and the number of
// days per year is different for leap years.
//func (d Period) IsPrecise() bool {
//	return d.years == 0 && d.months == 0
//}

// Format converts the period to human-readable form using the default localisation.
func (period Period) Format() string {
	return period.FormatWithPeriodNames(PeriodYearNames, PeriodMonthNames, PeriodWeekNames, PeriodDayNames)
}

// FormatWithPeriodNames converts the period to human-readable form in a localisable way.
func (period Period) FormatWithPeriodNames(yearNames Plurals, monthNames Plurals, weekNames Plurals, dayNames Plurals) string {
	period = period.Abs()

	parts := make([]string, 0)
	parts = appendNonBlank(parts, yearNames.FormatFloat(absFloat1000(period.years)))
	parts = appendNonBlank(parts, monthNames.FormatFloat(absFloat1000(period.months)))

	if (period.years == 0 && period.months == 0) || period.days > 0 {
		if len(weekNames) > 0 {
			weeks := period.days / 7000
			mdays := period.days % 7000
			//fmt.Printf("%v %#v - %d %d\n", period, period, weeks, mdays)
			if weeks > 0 {
				parts = appendNonBlank(parts, weekNames.FormatInt(int(weeks)))
			}
			if mdays > 0 || weeks == 0 {
				parts = appendNonBlank(parts, dayNames.FormatFloat(absFloat1000(mdays)))
			}
		} else {
			parts = appendNonBlank(parts, dayNames.FormatFloat(absFloat1000(period.days)))
		}
	}

	return strings.Join(parts, ", ")
}

func appendNonBlank(parts []string, s string) []string {
	if s == "" {
		return parts
	}
	return append(parts, s)
}

// PeriodDayNames provides the English default format names for the days part of the period.
// This is a sequence of plurals where the first match is used, otherwise the last one is used.
// The last one must include a "%g" placeholder for the number.
var PeriodDayNames = Plurals{Case{0, "%v days"}, Case{1, "%v day"}, Case{2, "%v days"}}

// PeriodWeekNames is as for PeriodDayNames but for weeks.
var PeriodWeekNames = Plurals{Case{0, ""}, Case{1, "%v week"}, Case{2, "%v weeks"}}

// PeriodMonthNames is as for PeriodDayNames but for months.
var PeriodMonthNames = Plurals{Case{0, ""}, Case{1, "%g month"}, Case{2, "%g months"}}

// PeriodYearNames is as for PeriodDayNames but for years.
var PeriodYearNames = Plurals{Case{0, ""}, Case{1, "%g year"}, Case{2, "%g years"}}

// String converts the period to -8601 form.
func (period Period) String() string {
	if period.IsZero() {
		return "P0D"
	}

	s := ""
	if period.years < 0 || period.months < 0 || period.days < 0 {
		s = "-"
	}

	y, m, w, d := "", "", "", ""

	if period.years != 0 {
		y = fmt.Sprintf("%gY", absFloat1000(period.years))
	}
	if period.months != 0 {
		m = fmt.Sprintf("%gM", absFloat1000(period.months))
	}
	if period.days != 0 {
		//days := absInt32(period.days)
		//weeks := days / 7
		//if (weeks >= 1000) {
		//	w = fmt.Sprintf("%gW", absFloat(weeks))
		//}
		//mdays := days % 7
		if period.days != 0 {
			d = fmt.Sprintf("%gD", absFloat1000(period.days))
		}
	}

	return fmt.Sprintf("%sP%s%s%s%s", s, y, m, w, d)
}

func absFloat1000(v int32) float32 {
	f := float32(v) / 1000
	if v < 0 {
		return -f
	}
	return f
}

// Abs converts a negative period to a positive one.
func (period Period) Abs() Period {
	return Period{absInt32(period.years), absInt32(period.months), absInt32(period.days)}
}

func absInt32(v int32) int32 {
	if v < 0 {
		return -v
	}
	return v
}

// Negate changes the sign of the period.
func (period Period) Negate() Period {
	return Period{-period.years, -period.months, -period.days}
}

// Sign returns +1 for positive periods and -1 for negative periods.
func (period Period) Sign() int {
	if period.years < 0 {
		return -1
	}
	return 1
}

// Years gets the whole number of years in the period.
func (period Period) Years() int {
	return int(period.YearsFloat())
}

// YearsFloat gets the number of years in the period, including a fraction if any is present.
func (period Period) YearsFloat() float32 {
	return float32(period.years) / 1000
}

// Months gets the whole number of months in the period.
func (period Period) Months() int {
	return int(period.MonthsFloat())
}

// MonthsFloat gets the number of months in the period.
func (period Period) MonthsFloat() float32 {
	return float32(period.months) / 1000
}

// Days gets the whole number of days in the period. This includes the implied
// number of weeks but excludes the specified years and months.
func (period Period) Days() int {
	return int(period.DaysFloat())
}

// DaysFloat gets the number of days in the period. This includes the implied
// number of weeks.
func (period Period) DaysFloat() float32 {
	return float32(period.days) / 1000
}

// Weeks calculates the number of whole weeks from the number of days. If the result
// would contain a fraction, it is truncated.
func (period Period) Weeks() int {
	return int(period.days) / 7000
}

// ModuloDays calculates the whole number of days remaining after the whole number of weeks
// has been excluded.
func (period Period) ModuloDays() int {
	days := absInt32(period.days) % 7000
	f := int(days / 1000)
	if period.days < 0 {
		return -f
	}
	return f
}

//TODO gobencode

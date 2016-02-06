package period

import (
	"fmt"
	. "github.com/rickb777/plural"
	"strconv"
	"strings"
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

// MustParse is as per Parse except that it panics if the string cannot be parsed.
// This is intended for setup code; don't use it for user inputs.
func MustParse(value string) Period {
	d, err := Parse(value)
	if err != nil {
		panic(err)
	}
	return d
}

// Parse parses strings that specify periods using ISO-8601 rules.
//
// In addition, a plus or minus sign can precede the period, e.g. "-P10D"
//
// The zero value can be represented in several ways: all of the following
// are equivalent: "P0Y", "P0M", "P0W", "P0D", and "P0".
func Parse(period string) (Period, error) {
	if period == "" {
		return Period{}, fmt.Errorf("Cannot parse a blank string as a period.")
	}

	if period == "P0" {
		return Period{}, nil
	}

	pcopy := period
	negate := false
	if pcopy[0] == '-' {
		negate = true
		pcopy = pcopy[1:]
	} else if pcopy[0] == '+' {
		pcopy = pcopy[1:]
	}

	if pcopy[0] != 'P' {
		return Period{}, fmt.Errorf("Expected 'P' period mark at the start: %s", period)
	}
	pcopy = pcopy[1:]

	result := Period{}

	st := parseState{period, pcopy, false, nil}
	t := strings.IndexByte(pcopy, 'T')
	if t >= 0 {
		st.pcopy = pcopy[t+1:]

		result.hours, st = parseField(st, 'H')
		if st.err != nil {
			return Period{}, st.err
		}

		result.minutes, st = parseField(st, 'M')
		if st.err != nil {
			return Period{}, st.err
		}

		result.seconds, st = parseField(st, 'S')
		if st.err != nil {
			return Period{}, st.err
		}

		st.pcopy = pcopy[:t]
	}

	result.years, st = parseField(st, 'Y')
	if st.err != nil {
		return Period{}, st.err
	}

	result.months, st = parseField(st, 'M')
	if st.err != nil {
		return Period{}, st.err
	}

	weeks, st := parseField(st, 'W')
	if st.err != nil {
		return Period{}, st.err
	}

	days, st := parseField(st, 'D')
	if st.err != nil {
		return Period{}, st.err
	}

	result.days = weeks*7 + days
	//fmt.Printf("%#v\n", st)

	if !st.ok {
		return Period{}, fmt.Errorf("Expected 'Y', 'M', 'W', 'D', 'H', 'M', or 'S' marker: %s", period)
	}
	if negate {
		return result.Negate(), nil
	}
	return result, nil
}

type parseState struct {
	period, pcopy string
	ok            bool
	err           error
}

func parseField(st parseState, mark byte) (int16, parseState) {
	//fmt.Printf("%c %#v\n", mark, st)
	r := int16(0)
	m := strings.IndexByte(st.pcopy, mark)
	if m > 0 {
		r, st.err = parseDecimalFixedPoint(st.pcopy[:m], st.period)
		if st.err != nil {
			return 0, st
		}
		st.pcopy = st.pcopy[m+1:]
		st.ok = true
	}
	return r, st
}

// Fixed-point three decimal places
func parseDecimalFixedPoint(s, original string) (int16, error) {
	//was := s
	dec := strings.IndexByte(s, '.')
	if dec < 0 {
		dec = strings.IndexByte(s, ',')
	}

	if dec >= 0 {
		dp := len(s) - dec
		if dp > 1 {
			s = s[:dec] + s[dec+1:dec+2]
		} else {
			s = s[:dec] + s[dec+1:] + "0"
		}
	} else {
		s = s + "0"
	}

	n, e := strconv.ParseInt(s, 10, 16)
	//fmt.Printf("ParseInt(%s) = %d -- from %s in %s %d\n", s, n, was, original, dec)
	return int16(n), e
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
	return period.FormatWithPeriodNames(PeriodYearNames, PeriodMonthNames, PeriodWeekNames, PeriodDayNames, PeriodHourNames, PeriodMinuteNames, PeriodSecondNames)
}

// FormatWithPeriodNames converts the period to human-readable form in a localisable way.
func (period Period) FormatWithPeriodNames(yearNames, monthNames, weekNames, dayNames, hourNames, minNames, secNames Plurals) string {
	period = period.Abs()

	parts := make([]string, 0)
	parts = appendNonBlank(parts, yearNames.FormatFloat(absFloat10(period.years)))
	parts = appendNonBlank(parts, monthNames.FormatFloat(absFloat10(period.months)))

	if period.days > 0 || (period.IsZero()) {
		if len(weekNames) > 0 {
			weeks := period.days / 70
			mdays := period.days % 70
			//fmt.Printf("%v %#v - %d %d\n", period, period, weeks, mdays)
			if weeks > 0 {
				parts = appendNonBlank(parts, weekNames.FormatInt(int(weeks)))
			}
			if mdays > 0 || weeks == 0 {
				parts = appendNonBlank(parts, dayNames.FormatFloat(absFloat10(mdays)))
			}
		} else {
			parts = appendNonBlank(parts, dayNames.FormatFloat(absFloat10(period.days)))
		}
	}
	parts = appendNonBlank(parts, hourNames.FormatFloat(absFloat10(period.hours)))
	parts = appendNonBlank(parts, minNames.FormatFloat(absFloat10(period.minutes)))
	parts = appendNonBlank(parts, secNames.FormatFloat(absFloat10(period.seconds)))

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

// PeriodHourNames is as for PeriodDayNames but for hours.
var PeriodHourNames = Plurals{Case{0, ""}, Case{1, "%v hour"}, Case{2, "%v hours"}}

// PeriodMinuteNames is as for PeriodDayNames but for minutes.
var PeriodMinuteNames = Plurals{Case{0, ""}, Case{1, "%g minute"}, Case{2, "%g minutes"}}

// PeriodSecondNames is as for PeriodDayNames but for seconds.
var PeriodSecondNames = Plurals{Case{0, ""}, Case{1, "%g second"}, Case{2, "%g seconds"}}

// String converts the period to -8601 form.
func (period Period) String() string {
	if period.IsZero() {
		return "P0D"
	}

	s := ""
	if period.Sign() < 0 {
		s = "-"
	}

	y, m, w, d, t, hh, mm, ss := "", "", "", "", "", "", "", ""

	if period.years != 0 {
		y = fmt.Sprintf("%gY", absFloat10(period.years))
	}
	if period.months != 0 {
		m = fmt.Sprintf("%gM", absFloat10(period.months))
	}
	if period.days != 0 {
		//days := absInt32(period.days)
		//weeks := days / 7
		//if (weeks >= 10) {
		//	w = fmt.Sprintf("%gW", absFloat(weeks))
		//}
		//mdays := days % 7
		if period.days != 0 {
			d = fmt.Sprintf("%gD", absFloat10(period.days))
		}
	}
	if period.hours != 0 || period.minutes != 0 || period.seconds != 0 {
		t = "T"
	}
	if period.hours != 0 {
		hh = fmt.Sprintf("%gH", absFloat10(period.hours))
	}
	if period.minutes != 0 {
		mm = fmt.Sprintf("%gM", absFloat10(period.minutes))
	}
	if period.seconds != 0 {
		ss = fmt.Sprintf("%gS", absFloat10(period.seconds))
	}

	return fmt.Sprintf("%sP%s%s%s%s%s%s%s%s", s, y, m, w, d, t, hh, mm, ss)
}

func absFloat10(v int16) float32 {
	f := float32(v) / 10
	if v < 0 {
		return -f
	}
	return f
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
// decimal place; each field only is int16.
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
// the duration is calculated on the basis of a year being 365.2 days and a month being
// 1/12 of a year; days are all 24 hours long.
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

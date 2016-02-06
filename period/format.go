package period

import (
	"fmt"
	. "github.com/rickb777/plural"
	"strings"
)

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
// The last one must include a "%v" placeholder for the number.
var PeriodDayNames = Plurals{Case{0, "%v days"}, Case{1, "%v day"}, Case{2, "%v days"}}

// PeriodWeekNames is as for PeriodDayNames but for weeks.
var PeriodWeekNames = Plurals{Case{0, ""}, Case{1, "%v week"}, Case{2, "%v weeks"}}

// PeriodMonthNames is as for PeriodDayNames but for months.
var PeriodMonthNames = Plurals{Case{0, ""}, Case{1, "%v month"}, Case{2, "%g months"}}

// PeriodYearNames is as for PeriodDayNames but for years.
var PeriodYearNames = Plurals{Case{0, ""}, Case{1, "%v year"}, Case{2, "%v years"}}

// PeriodHourNames is as for PeriodDayNames but for hours.
var PeriodHourNames = Plurals{Case{0, ""}, Case{1, "%v hour"}, Case{2, "%v hours"}}

// PeriodMinuteNames is as for PeriodDayNames but for minutes.
var PeriodMinuteNames = Plurals{Case{0, ""}, Case{1, "%v minute"}, Case{2, "%v minutes"}}

// PeriodSecondNames is as for PeriodDayNames but for seconds.
var PeriodSecondNames = Plurals{Case{0, ""}, Case{1, "%v second"}, Case{2, "%v seconds"}}

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

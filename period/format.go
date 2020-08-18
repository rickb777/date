// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package period

import (
	"fmt"
	"github.com/rickb777/plural"
	"strings"
)

// Format converts the period to human-readable form using the default localisation.
func (period Period) Format() string {
	return period.FormatWithPeriodNames(PeriodYearNames, PeriodMonthNames, PeriodWeekNames, PeriodDayNames, PeriodHourNames, PeriodMinuteNames, PeriodSecondNames)
}

// FormatWithPeriodNames converts the period to human-readable form in a localisable way.
func (period Period) FormatWithPeriodNames(yearNames, monthNames, weekNames, dayNames, hourNames, minNames, secNames plural.Plurals) string {
	period = period.Abs()

	parts := make([]string, 0)
	years, months := period.unpackYM()
	parts = appendNonBlank(parts, yearNames.FormatFloat(absFloat1i32(years)))
	parts = appendNonBlank(parts, monthNames.FormatFloat(absFloat100i32(months)))

	if period.centiDays > 0 || (period.IsZero()) {
		if len(weekNames) > 0 {
			weeks := period.centiDays / 700
			mdays := period.centiDays % 700
			//fmt.Printf("%v %#v - %d %d\n", period, period, weeks, centiDays)
			if weeks > 0 {
				parts = appendNonBlank(parts, weekNames.FormatInt(int(weeks)))
			}
			if mdays > 0 || weeks == 0 {
				parts = appendNonBlank(parts, dayNames.FormatFloat(absFloat100i32(mdays)))
			}
		} else {
			parts = appendNonBlank(parts, dayNames.FormatFloat(absFloat100i32(period.centiDays)))
		}
	}

	parts = appendNonBlank(parts, hourNames.FormatFloat(absFloat1i32(period.centiSeconds/360000)))
	parts = appendNonBlank(parts, minNames.FormatFloat(absFloat1i32((period.centiSeconds%360000)/6000)))
	parts = appendNonBlank(parts, secNames.FormatFloat(absFloat100i32(period.centiSeconds%6000)))

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
var PeriodDayNames = plural.FromZero("%v days", "%v day", "%v days")

// PeriodWeekNames is as for PeriodDayNames but for weeks.
var PeriodWeekNames = plural.FromZero("", "%v week", "%v weeks")

// PeriodMonthNames is as for PeriodDayNames but for months.
var PeriodMonthNames = plural.FromZero("", "%v month", "%v months")

// PeriodYearNames is as for PeriodDayNames but for years.
var PeriodYearNames = plural.FromZero("", "%v year", "%v years")

// PeriodHourNames is as for PeriodDayNames but for hours.
var PeriodHourNames = plural.FromZero("", "%v hour", "%v hours")

// PeriodMinuteNames is as for PeriodDayNames but for minutes.
var PeriodMinuteNames = plural.FromZero("", "%v minute", "%v minutes")

// PeriodSecondNames is as for PeriodDayNames but for seconds.
var PeriodSecondNames = plural.FromZero("", "%v second", "%v seconds")

// String converts the period to ISO-8601 form.
func (period Period) String() string {
	return period.showAs
}

func (period Period) toString() string {
	if period.IsZero() {
		return "P0D"
	}

	buf := &strings.Builder{}
	if period.Sign() < 0 {
		buf.WriteByte('-')
		period = period.Negate()
	}

	buf.WriteByte('P')

	if period.centiMonths != 0 {
		years, months := period.unpackYM()
		if years != 0 {
			fmt.Fprintf(buf, "%dY", years)
		}
		if months != 0 {
			fmt.Fprintf(buf, "%gM", absFloat100i32(months))
		}
	}

	if period.centiDays != 0 {
		if period.centiDays%700 == 0 {
			fmt.Fprintf(buf, "%gW", absFloat100i32(period.centiDays/7))
		} else {
			fmt.Fprintf(buf, "%gD", absFloat100i32(period.centiDays))
		}
	}

	if period.centiSeconds != 0 {
		hours, minutes, seconds := period.unpackHMS()
		buf.WriteByte('T')
		if hours != 0 {
			fmt.Fprintf(buf, "%dH", hours)
		}
		if minutes != 0 {
			fmt.Fprintf(buf, "%dM", minutes)
		}
		if seconds != 0 {
			fmt.Fprintf(buf, "%gS", absFloat100i32(seconds))
		}
	}

	return buf.String()
}

func (period Period) unpackYM() (int32, int32) {
	years := period.centiMonths / 1200
	months := period.centiMonths - (years * 1200)
	return years, months
}

func (period Period) unpackHMS() (int32, int32, int32) {
	hours := period.centiSeconds / 360000
	seconds := period.centiSeconds - hours*360000

	minutes := seconds / 6000
	seconds -= minutes * 6000
	return hours, minutes, seconds
}

func absFloat1i32(v int32) float32 {
	f := float32(v)
	if v < 0 {
		return -f
	}
	return f
}

func absFloat100i32(v int32) float32 {
	return absFloat1i32(v) / 100
}

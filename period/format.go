// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package period

import (
	"fmt"
	"io"
	"strings"

	"github.com/rickb777/plural"
)

// Format converts the period to human-readable form using the default localisation.
// Multiples of 7 days are shown as weeks.
func (period Period) Format() string {
	return period.FormatWithPeriodNames(PeriodYearNames, PeriodMonthNames, PeriodWeekNames, PeriodDayNames, PeriodHourNames, PeriodMinuteNames, PeriodSecondNames)
}

// FormatWithoutWeeks converts the period to human-readable form using the default localisation.
// Multiples of 7 days are not shown as weeks.
func (period Period) FormatWithoutWeeks() string {
	return period.FormatWithPeriodNames(PeriodYearNames, PeriodMonthNames, plural.Plurals{}, PeriodDayNames, PeriodHourNames, PeriodMinuteNames, PeriodSecondNames)
}

// FormatWithPeriodNames converts the period to human-readable form in a localisable way.
func (period Period) FormatWithPeriodNames(yearNames, monthNames, weekNames, dayNames, hourNames, minNames, secNames plural.Plurals) string {
	period = period.Abs()

	parts := make([]string, 0)
	parts = appendNonBlank(parts, yearNames.FormatFloat(period.YearsFloat()))
	parts = appendNonBlank(parts, monthNames.FormatFloat(period.MonthsFloat()))

	if period.days > 0 || period.fpart == Day || period.IsZero() {
		if len(weekNames) > 0 {
			weeks := period.days / 7
			mdays := period.days % 7
			//fmt.Printf("%v %#v - %d %d\n", period, period, weeks, mdays)
			if weeks > 0 {
				parts = appendNonBlank(parts, weekNames.FormatInt(int(weeks)))
			}
			if mdays > 0 || weeks == 0 || period.fpart == Day {
				period.days = mdays
				parts = appendNonBlank(parts, dayNames.FormatFloat(period.DaysFloat()))
			}
		} else {
			parts = appendNonBlank(parts, dayNames.FormatFloat(period.DaysFloat()))
		}
	}
	parts = appendNonBlank(parts, hourNames.FormatFloat(period.HoursFloat()))
	parts = appendNonBlank(parts, minNames.FormatFloat(period.MinutesFloat()))
	parts = appendNonBlank(parts, secNames.FormatFloat(period.SecondsFloat()))

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
	return period.toPeriod64("").String()
}

func (p64 period64) String() string {
	if p64 == (period64{}) {
		return "P0D"
	}

	buf := &strings.Builder{}
	if p64.neg {
		buf.WriteByte('-')
	}

	buf.WriteByte('P')

	writeField64(buf, p64.years, p64.fraction, p64.fpart, Year)
	writeField64(buf, p64.months, p64.fraction, p64.fpart, Month)

	if p64.days != 0 && p64.days%7 == 0 {
		writeField64(buf, p64.days/7, 0, 0, Week)
	} else {
		writeField64(buf, p64.days, p64.fraction, p64.fpart, Day)
	}

	if p64.hours != 0 || p64.minutes != 0 || p64.seconds != 0 || (p64.fraction != 0 && p64.fpart.IsOneOf(Hour, Minute, Second)) {
		buf.WriteByte('T')
	}

	writeField64(buf, p64.hours, p64.fraction, p64.fpart, Hour)
	writeField64(buf, p64.minutes, p64.fraction, p64.fpart, Minute)
	writeField64(buf, p64.seconds, p64.fraction, p64.fpart, Second)

	return buf.String()
}

func writeField64(w io.Writer, field int64, fraction int8, fpart, designator designator) {
	if field != 0 || (fraction != 0 && fpart == designator) {
		fmt.Fprintf(w, "%d", field)
		if fpart == designator {
			if fraction%10 == 0 {
				fmt.Fprintf(w, ".%d", fraction/10)
			} else {
				fmt.Fprintf(w, ".%02d", fraction)
			}
		}
		w.(io.ByteWriter).WriteByte(designator.Byte())
	}
}

// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package period

import (
	"fmt"
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

var oneDay = 24 * time.Hour
var oneMonthApprox = 2629746 * time.Second // 30.436875 days
var oneYearApprox = 31556952 * time.Second // 365.2425 days

func TestParseErrors(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		value     string
		normalise bool
		expected  string
		expvalue  string
	}{
		{"", false, "cannot parse a blank string as a period", ""},
		{`P000`, false, `: missing designator at the end`, "P000"},
		{"XY", false, ": expected 'P' period mark at the start", "XY"},
		{"PxY", false, ": expected a number but found 'x'", "PxY"},
		{"PxW", false, ": expected a number but found 'x'", "PxW"},
		{"PxD", false, ": expected a number but found 'x'", "PxD"},
		{"PTxH", false, ": expected a number but found 'x'", "PTxH"},
		{"PTxM", false, ": expected a number but found 'x'", "PTxM"},
		{"PTxS", false, ": expected a number but found 'x'", "PTxS"},
		{"P1HT1M", false, ": 'H' designator cannot occur here", "P1HT1M"},
		{"PT1Y", false, ": 'Y' designator cannot occur here", "PT1Y"},
		{"P1S", false, ": 'S' designator cannot occur here", "P1S"},
		{"P1D2D", false, ": 'D' designator cannot occur more than once", "P1D2D"},
		{"PT1HT1S", false, ": 'T' designator cannot occur more than once", "PT1HT1S"},
		{"P0.1YT0.1S", false, ": 'Y' & 'S' only the last field can have a fraction", "P0.1YT0.1S"},
		{"P", false, ": expected 'Y', 'M', 'W', 'D', 'H', 'M', or 'S' designator", "P"},
		// integer overflow
		{"P32768Y", false, ": integer overflow occurred in years", "P32768Y"},
		{"P32768M", false, ": integer overflow occurred in months", "P32768M"},
		{"P32768W", false, ": integer overflow occurred in days", "P32768W"},
		{"P32768D", false, ": integer overflow occurred in days", "P32768D"},
		{"PT32768H", false, ": integer overflow occurred in hours", "PT32768H"},
		{"PT32768M", false, ": integer overflow occurred in minutes", "PT32768M"},
		{"PT32768S", false, ": integer overflow occurred in seconds", "PT32768S"},
		{"PT32768H32768M32768S", false, ": integer overflow occurred in hours,minutes,seconds", "PT32768H32768M32768S"},
		{"PT103412160000S", false, ": integer overflow occurred in seconds", "PT103412160000S"},
	}
	for i, c := range cases {
		_, ep := Parse(c.value, c.normalise)
		g.Expect(ep).To(HaveOccurred(), info(i, c.value))
		g.Expect(ep.Error()).To(Equal(c.expvalue+c.expected), info(i, c.value))

		_, en := Parse("-"+c.value, c.normalise)
		g.Expect(en).To(HaveOccurred(), info(i, c.value))
		if c.expvalue != "" {
			g.Expect(en.Error()).To(Equal("-"+c.expvalue+c.expected), info(i, c.value))
		} else {
			g.Expect(en.Error()).To(Equal(c.expected), info(i, c.value))
		}
	}
}

//-------------------------------------------------------------------------------------------------

func TestParsePeriodWithNormalise(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		value    string
		reversed string
		period   Period
	}{
		// all rollovers
		{"PT1234.5S", "PT20M34.5S", Period{minutes: 20, seconds: 34, fraction: 50, fpart: Second}},
		{"PT1234.5M", "PT20H34M30S", Period{hours: 20, minutes: 34, seconds: 30}},
		{"PT12345.6H", "PT12345H36M", Period{hours: 12345, minutes: 36}},
		//TODO {"P32768.1D", "P89Y8M17DT22H8M.1D", Period{years: 89, months: 8, days: 17, hours: 24, minutes: 32, fraction: 10, fpart: Day}},
		{"P1234.5M", "P102Y10.5M", Period{years: 102, months: 10, fraction: 50, fpart: Month}},
		// largest possible number of seconds normalised only in hours, mins, sec
		{"PT11592000S", "PT3220H", Period{hours: 3220}},
		{"-PT11592000S", "-PT3220H", Period{hours: -3220}},
		{"PT11595599S", "PT3220H59M59S", Period{hours: 3220, minutes: 59, seconds: 59}},
		// largest possible number of seconds normalised only in days, hours, mins, sec
		{"PT283046400S", "P468W", Period{days: 3276}},
		{"-PT283046400S", "-P468W", Period{days: -3276}},
		{"PT43084443590S", "P1365Y3M15DT4H73M50S", Period{years: 1365, months: 3, days: 15, hours: 4, minutes: 73, seconds: 50}},
		{"PT103412159999S", "P3277YT6H110M59S", Period{years: 3277, months: 0, days: 0, hours: 6, minutes: 110, seconds: 59}},
		{"PT283132799S", "P468WT23H59M59S", Period{days: 3276, hours: 23, minutes: 59, seconds: 59}},
		// other examples are in TestNormalise
	}
	for i, c := range cases {
		p, err := Parse(c.value)
		s := info(i, c.value)
		g.Expect(err).NotTo(HaveOccurred(), s)
		expectValid(t, p, s)
		g.Expect(p).To(Equal(c.period), s)
		// reversal is expected not to be an identity
		g.Expect(p.String()).To(Equal(c.reversed), s+" reversed")
	}
}

//-------------------------------------------------------------------------------------------------

func TestParsePeriodWithoutNormalise(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		value    string
		reversed string
		period   Period
	}{
		// zero
		{"P0D", "P0D", Period{}},
		// special zero cases: parse is not identity when reversed
		{"P0", "P0D", Period{}},
		{"P0Y", "P0D", Period{}},
		{"P0M", "P0D", Period{}},
		{"P0W", "P0D", Period{}},
		{"PT0H", "P0D", Period{}},
		{"PT0M", "P0D", Period{}},
		{"PT0S", "P0D", Period{}},
		// ones
		{"P1Y", "P1Y", Period{years: 1}},
		{"P1M", "P1M", Period{months: 1}},
		{"P1W", "P1W", Period{days: 7}},
		{"P1D", "P1D", Period{days: 1}},
		{"PT1H", "PT1H", Period{hours: 1}},
		{"PT1M", "PT1M", Period{minutes: 1}},
		{"PT1S", "PT1S", Period{seconds: 1}},
		// smallest
		{"P0.01Y", "P0.01Y", Period{fraction: 1, fpart: Year}},
		{"-P0.01Y", "-P0.01Y", Period{fraction: -1, fpart: Year}},
		{"P0.01M", "P0.01M", Period{fraction: 1, fpart: Month}},
		{"-P0.01M", "-P0.01M", Period{fraction: -1, fpart: Month}},
		{"P0.01D", "P0.01D", Period{fraction: 1, fpart: Day}},
		{"-P0.01D", "-P0.01D", Period{fraction: -1, fpart: Day}},
		{"PT0.01H", "PT0.01H", Period{fraction: 1, fpart: Hour}},
		{"-PT0.01H", "-PT0.01H", Period{fraction: -1, fpart: Hour}},
		{"PT0.01M", "PT0.01M", Period{fraction: 1, fpart: Minute}},
		{"-PT0.01M", "-PT0.01M", Period{fraction: -1, fpart: Minute}},
		{"PT0.01S", "PT0.01S", Period{fraction: 1, fpart: Second}},
		{"-PT0.01S", "-PT0.01S", Period{fraction: -1, fpart: Second}},
		// week special case: also not identity when reversed
		{"P0.01W", "P0.07D", Period{fraction: 7, fpart: Day}},
		{"-P0.01W", "-P0.07D", Period{fraction: -7, fpart: Day}},
		// largest
		{"PT32767.99S", "PT32767.99S", Period{seconds: 32767, fraction: 99, fpart: Second}},
		{"PT32767.99M", "PT32767.99M", Period{minutes: 32767, fraction: 99, fpart: Minute}},
		{"PT32767.99H", "PT32767.99H", Period{hours: 32767, fraction: 99, fpart: Hour}},
		{"P32766.99D", "P32766.99D", Period{days: 32766, fraction: 99, fpart: Day}},
		{"P32767.99M", "P32767.99M", Period{months: 32767, fraction: 99, fpart: Month}},
		{"P32767.99Y", "P32767.99Y", Period{years: 32767, fraction: 99, fpart: Year}},

		{"P3Y", "P3Y", Period{years: 3}},
		{"P6M", "P6M", Period{months: 6}},
		{"P5W", "P5W", Period{days: 35}},
		{"P4D", "P4D", Period{days: 4}},
		{"PT12H", "PT12H", Period{hours: 12}},
		{"PT30M", "PT30M", Period{minutes: 30}},
		{"PT25S", "PT25S", Period{seconds: 25}},
		{"PT30M67.6S", "PT30M67.6S", Period{minutes: 30, seconds: 67, fraction: 60, fpart: Second}},
		{"P2.Y", "P2Y", Period{years: 2}},
		{"P2.5Y", "P2.5Y", Period{years: 2, fraction: 50, fpart: Year}},
		{"P2.15Y", "P2.15Y", Period{years: 2, fraction: 15, fpart: Year}},
		{"P2.125Y", "P2.12Y", Period{years: 2, fraction: 12, fpart: Year}},
		{"P1Y2.M", "P1Y2M", Period{years: 1, months: 2}},
		{"P1Y2.5M", "P1Y2.5M", Period{years: 1, months: 2, fraction: 50, fpart: Month}},
		{"P1Y2.15M", "P1Y2.15M", Period{years: 1, months: 2, fraction: 15, fpart: Month}},
		{"P1Y2.125M", "P1Y2.12M", Period{years: 1, months: 2, fraction: 12, fpart: Month}},
		// others
		{"P3Y6M5W4DT12H40M5S", "P3Y6M39DT12H40M5S", Period{years: 3, months: 6, days: 39, hours: 12, minutes: 40, seconds: 5}},
		{"+P3Y6M5W4DT12H40M5S", "P3Y6M39DT12H40M5S", Period{years: 3, months: 6, days: 39, hours: 12, minutes: 40, seconds: 5}},
		{"-P3Y6M5W4DT12H40M5S", "-P3Y6M39DT12H40M5S", Period{years: -3, months: -6, days: -39, hours: -12, minutes: -40, seconds: -5}},
		{"P1Y14M35DT48H125M800S", "P1Y14M5WT48H125M800S", Period{years: 1, months: 14, days: 35, hours: 48, minutes: 125, seconds: 800}},
	}
	for i, c := range cases {
		p, err := Parse(c.value, false)
		s := info(i, c.value)
		g.Expect(err).NotTo(HaveOccurred(), s)
		expectValid(t, p, s)
		g.Expect(p).To(Equal(c.period), s)
		// reversal is usually expected to be an identity
		g.Expect(p.String()).To(Equal(c.reversed), s+" reversed")
	}
}

//-------------------------------------------------------------------------------------------------

func TestPeriodString(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		value  string
		period Period
	}{
		// note: the negative cases are also covered (see below)

		{"P0D", Period{}},
		// ones
		{"P1Y", Period{years: 1}},
		{"P1M", Period{months: 1}},
		{"P1W", Period{days: 7}},
		{"P1D", Period{days: 1}},
		{"PT1H", Period{hours: 1}},
		{"PT1M", Period{minutes: 1}},
		{"PT1S", Period{seconds: 1}},
		// smallest
		{"P0.01Y", Period{fraction: 1, fpart: Year}},
		{"P0.01M", Period{fraction: 1, fpart: Month}},
		{"P0.07D", Period{fraction: 7, fpart: Day}},
		{"P0.01D", Period{fraction: 1, fpart: Day}},
		{"PT0.01H", Period{fraction: 1, fpart: Hour}},
		{"PT0.01M", Period{fraction: 1, fpart: Minute}},
		{"PT0.01S", Period{fraction: 1, fpart: Second}},

		{"P3Y", Period{years: 3}},
		{"P6M", Period{months: 6}},
		{"P5W", Period{days: 35}},
		{"P4W", Period{days: 28}},
		{"P4D", Period{days: 4}},
		{"PT12H", Period{hours: 12}},
		{"PT30M", Period{minutes: 30}},
		{"PT5S", Period{seconds: 5}},
		{"P3Y6M39DT1H2M4.09S", Period{years: 3, months: 6, days: 39, hours: 1, minutes: 2, seconds: 4, fraction: 9, fpart: Second}},

		{"P2.5Y", Period{years: 2, fraction: 50, fpart: Year}},
		{"P2.49Y", Period{years: 2, fraction: 49, fpart: Year}},
		{"P2.5M", Period{months: 2, fraction: 50, fpart: Month}},
		{"P2.5D", Period{days: 2, fraction: 50, fpart: Day}},
		{"PT2.5H", Period{hours: 2, fraction: 50, fpart: Hour}},
		{"PT2.5M", Period{minutes: 2, fraction: 50, fpart: Minute}},
		{"PT2.5S", Period{seconds: 2, fraction: 50, fpart: Second}},
	}
	for i, c := range cases {
		sp := c.period.String()
		g.Expect(sp).To(Equal(c.value), info(i, c.value))

		if !c.period.IsZero() {
			sn := c.period.Negate().String()
			g.Expect(sn).To(Equal("-"+c.value), info(i, c.value))
		}
	}
}

//-------------------------------------------------------------------------------------------------

func TestPeriodIntComponents(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		value                      string
		y, m, w, d, dx, hh, mm, ss int
	}{
		// note: the negative cases are also covered (see below)

		{value: "P0D"},
		{value: "P1Y", y: 1},
		{value: "P1W", w: 1, d: 7},
		{value: "P6M", m: 6},
		{value: "P12M", y: 1},
		{value: "P39D", w: 5, d: 39, dx: 4},
		{value: "P4D", d: 4, dx: 4},
		{value: "PT12H", hh: 12},
		{value: "PT60M", hh: 1},
		{value: "PT30M", mm: 30},
		{value: "PT5S", ss: 5},
	}
	for i, c := range cases {
		pp := MustParse(c.value)
		g.Expect(pp.Years()).To(Equal(c.y), info(i, pp))
		g.Expect(pp.Months()).To(Equal(c.m), info(i, pp))
		g.Expect(pp.Weeks()).To(Equal(c.w), info(i, pp))
		g.Expect(pp.Days()).To(Equal(c.d), info(i, pp))
		g.Expect(pp.ModuloDays()).To(Equal(c.dx), info(i, pp))
		g.Expect(pp.Hours()).To(Equal(c.hh), info(i, pp))
		g.Expect(pp.Minutes()).To(Equal(c.mm), info(i, pp))
		g.Expect(pp.Seconds()).To(Equal(c.ss), info(i, pp))

		pn := pp.Negate()
		g.Expect(pn.Years()).To(Equal(-c.y), info(i, pn))
		g.Expect(pn.Months()).To(Equal(-c.m), info(i, pn))
		g.Expect(pn.Weeks()).To(Equal(-c.w), info(i, pn))
		g.Expect(pn.Days()).To(Equal(-c.d), info(i, pn))
		g.Expect(pn.ModuloDays()).To(Equal(-c.dx), info(i, pn))
		g.Expect(pn.Hours()).To(Equal(-c.hh), info(i, pn))
		g.Expect(pn.Minutes()).To(Equal(-c.mm), info(i, pn))
		g.Expect(pn.Seconds()).To(Equal(-c.ss), info(i, pn))
	}
}

//-------------------------------------------------------------------------------------------------

func TestPeriodFloatComponents(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		value                      string
		y, m, w, d, dx, hh, mm, ss float32
	}{
		// note: the negative cases are also covered (see below)

		{value: "P0"}, // zero case

		// YMD cases
		{value: "P1Y", y: 1},
		{value: "P1.5Y", y: 1.5},
		{value: "P1.01Y", y: 1.01},
		{value: "P1M", m: 1},
		{value: "P1.5M", m: 1.5},
		{value: "P1.01M", m: 1.01},
		{value: "P6M", m: 6},
		{value: "P12M", m: 12},
		{value: "P7D", w: 1, d: 7},
		{value: "P7.7D", w: 1.1, d: 7.7},
		{value: "P7.01D", w: 7.01 / 7, d: 7.01},
		{value: "P1D", w: 1.0 / 7, d: 1},
		{value: "P1.1D", w: 1.1 / 7, d: 1.1},
		{value: "P1.01D", w: 1.01 / 7, d: 1.01},
		{value: "P39D", w: 5.571429, d: 39, dx: 4},
		{value: "P4D", w: 0.5714286, d: 4, dx: 4},

		// HMS cases
		{value: "PT1.1H", hh: 1.1},
		{value: "PT1.01H", hh: 1.01},
		{value: "PT1H6M", hh: 1, mm: 6},
		{value: "PT12H", hh: 12},
		{value: "PT1.1M", mm: 1.1},
		{value: "PT1.01M", mm: 1.01},
		{value: "PT1M6S", mm: 1, ss: 6},
		{value: "PT30M", mm: 30},
		{value: "PT1.1S", ss: 1.1},
		{value: "PT1.01S", ss: 1.01},
		{value: "PT5S", ss: 5},
	}
	for i, c := range cases {
		pp, _ := Parse(c.value, false)
		g.Expect(pp.YearsFloat()).To(Equal(c.y), info(i, pp))
		g.Expect(pp.MonthsFloat()).To(Equal(c.m), info(i, pp))
		g.Expect(pp.WeeksFloat()).To(Equal(c.w), info(i, pp))
		g.Expect(pp.DaysFloat()).To(Equal(c.d), info(i, pp))
		g.Expect(pp.HoursFloat()).To(Equal(c.hh), info(i, pp))
		g.Expect(pp.MinutesFloat()).To(Equal(c.mm), info(i, pp))
		g.Expect(pp.SecondsFloat()).To(Equal(c.ss), info(i, pp))

		pn := pp.Negate()
		g.Expect(pn.YearsFloat()).To(Equal(-c.y), info(i, pn))
		g.Expect(pn.MonthsFloat()).To(Equal(-c.m), info(i, pn))
		g.Expect(pn.WeeksFloat()).To(Equal(-c.w), info(i, pn))
		g.Expect(pn.DaysFloat()).To(Equal(-c.d), info(i, pn))
		g.Expect(pn.HoursFloat()).To(Equal(-c.hh), info(i, pn))
		g.Expect(pn.MinutesFloat()).To(Equal(-c.mm), info(i, pn))
		g.Expect(pn.SecondsFloat()).To(Equal(-c.ss), info(i, pn))
	}
}

//-------------------------------------------------------------------------------------------------

func TestPeriodToDuration(t *testing.T) {
	cases := []struct {
		value    string
		duration time.Duration
		precise  bool
	}{
		// note: the negative cases are also covered (see below)

		{"P0D", time.Duration(0), true},
		{"PT1S", 1 * time.Second, true},
		{"PT0.1S", 100 * time.Millisecond, true},
		{"PT3276S", 3276 * time.Second, true},
		{"PT1M", 60 * time.Second, true},
		{"PT0.1M", 6 * time.Second, true},
		{"PT3276M", 3276 * time.Minute, true},
		{"PT1H", 3600 * time.Second, true},
		{"PT0.1H", 360 * time.Second, true},
		{"PT3220H", 3220 * time.Hour, true},
		// days, months and years conversions are never precise
		{"P1D", 24 * time.Hour, false},
		{"P0.1D", 144 * time.Minute, false},
		{"P3276D", 3276 * 24 * time.Hour, false},
		{"P1M", oneMonthApprox, false},
		{"P0.1M", oneMonthApprox / 10, false},
		{"P3276M", 3276 * oneMonthApprox, false},
		{"P1Y", oneYearApprox, false},
		{"P3276Y", 3276 * oneYearApprox, false}, // near the upper limit of range
	}
	for i, c := range cases {
		testPeriodToDuration(t, i, c.value, c.duration, c.precise)
		testPeriodToDuration(t, i, "-"+c.value, -c.duration, c.precise)
	}
}

func testPeriodToDuration(t *testing.T, i int, value string, duration time.Duration, precise bool) {
	t.Helper()
	g := NewGomegaWithT(t)
	hint := info(i, "%s %s %v", value, duration, precise)
	pp := MustParse(value)
	d1, prec := pp.Duration()
	g.Expect(d1).To(Equal(duration), hint)
	g.Expect(prec).To(Equal(precise), hint)
	d2 := pp.DurationApprox()
	if precise {
		g.Expect(d2).To(Equal(duration), hint)
	}
}

//-------------------------------------------------------------------------------------------------

func TestSignPositiveNegative(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		value    string
		positive bool
		negative bool
		sign     int
	}{
		{"P0D", false, false, 0},
		{"PT1S", true, false, 1},
		{"PT0.1S", true, false, 1},
		{"-PT1S", false, true, -1},
		{"-PT0.1S", false, true, -1},
		{"PT1M", true, false, 1},
		{"PT0.1M", true, false, 1},
		{"-PT1M", false, true, -1},
		{"-PT0.1M", false, true, -1},
		{"PT1H", true, false, 1},
		{"PT0.1H", true, false, 1},
		{"-PT1H", false, true, -1},
		{"-PT0.1H", false, true, -1},
		{"P1D", true, false, 1},
		{"P10.D", true, false, 1},
		{"-P1D", false, true, -1},
		{"-P0.1D", false, true, -1},
		{"P1M", true, false, 1},
		{"P0.1M", true, false, 1},
		{"-P1M", false, true, -1},
		{"-P0.1M", false, true, -1},
		{"P1Y", true, false, 1},
		{"P0.1Y", true, false, 1},
		{"-P1Y", false, true, -1},
		{"-P0.1Y", false, true, -1},
	}
	for i, c := range cases {
		p := MustParse(c.value)
		g.Expect(p.IsPositive()).To(Equal(c.positive), info(i, c.value))
		g.Expect(p.IsNegative()).To(Equal(c.negative), info(i, c.value))
		g.Expect(p.Sign()).To(Equal(c.sign), info(i, c.value))
	}
}

//-------------------------------------------------------------------------------------------------

func TestPeriodApproxDays(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		value      string
		approxDays int
	}{
		// note: the negative cases are also covered (see below)

		{"P0D", 0},
		{"PT24H", 1},
		{"PT49H", 2},
		{"P1D", 1},
		{"P1M", 30},
		{"P1Y", 365},
	}
	for i, c := range cases {
		p := MustParse(c.value)
		td1 := p.TotalDaysApprox()
		g.Expect(td1).To(Equal(c.approxDays), info(i, c.value))

		td2 := p.Negate().TotalDaysApprox()
		g.Expect(td2).To(Equal(-c.approxDays), info(i, c.value))
	}
}

//-------------------------------------------------------------------------------------------------

func TestPeriodApproxMonths(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		value        string
		approxMonths int
	}{
		// note: the negative cases are also covered (see below)

		{"P0D", 0},
		{"P1D", 0},
		{"P30D", 0},
		{"P31D", 1},
		{"P60D", 1},
		{"P62D", 2},
		{"P1M", 1},
		{"P12M", 12},
		{"P2M31D", 3},
		{"P1Y", 12},
		{"P2Y3M", 27},
		{"PT24H", 0},
		{"PT744H", 1},
	}
	for i, c := range cases {
		p := MustParse(c.value)
		td1 := p.TotalMonthsApprox()
		g.Expect(td1).To(Equal(c.approxMonths), info(i, c.value))

		td2 := p.Negate().TotalMonthsApprox()
		g.Expect(td2).To(Equal(-c.approxMonths), info(i, c.value))
	}
}

//-------------------------------------------------------------------------------------------------

func TestNewPeriod(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		period                                       string
		years, months, days, hours, minutes, seconds int
	}{
		// note: the negative cases are also covered (see below)

		{period: "P0"}, // zero case

		{period: "PT1S", seconds: 1},
		{period: "PT1M", minutes: 1},
		{period: "PT1H", hours: 1},
		{period: "P1D", days: 1},
		{period: "P1M", months: 1},
		{period: "P1Y", years: 1},
		{period: "P100Y222M700D", years: 100, months: 222, days: 700},
	}
	for i, c := range cases {
		ep, _ := Parse(c.period, false)
		pp := New(c.years, c.months, c.days, c.hours, c.minutes, c.seconds)
		expectValid(t, pp, info(i, c.period))
		g.Expect(pp).To(Equal(ep), info(i, c.period))
		g.Expect(pp.Years()).To(Equal(c.years), info(i, c.period))
		g.Expect(pp.Months()).To(Equal(c.months), info(i, c.period))
		g.Expect(pp.Days()).To(Equal(c.days), info(i, c.period))

		pn := New(-c.years, -c.months, -c.days, -c.hours, -c.minutes, -c.seconds)
		en := ep.Negate()
		expectValid(t, pn, info(i, en))
		g.Expect(pn).To(Equal(en), info(i, en))
		g.Expect(pn.Years()).To(Equal(-c.years), info(i, en))
		g.Expect(pn.Months()).To(Equal(-c.months), info(i, en))
		g.Expect(pn.Days()).To(Equal(-c.days), info(i, en))
	}
}

//-------------------------------------------------------------------------------------------------

func TestNewHMS(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		period                  Period
		hours, minutes, seconds int
	}{
		// note: the negative cases are also covered (see below)

		{}, // zero case

		{period: Period{seconds: 1}, seconds: 1},
		{period: Period{minutes: 1}, minutes: 1},
		{period: Period{hours: 1}, hours: 1},
		{period: Period{hours: 3, minutes: 4, seconds: 5}, hours: 3, minutes: 4, seconds: 5},
		{period: Period{hours: 32767, minutes: 32767, seconds: 32767}, hours: 32767, minutes: 32767, seconds: 32767},
	}
	for i, c := range cases {
		pp := NewHMS(c.hours, c.minutes, c.seconds)
		expectValid(t, pp, info(i, c.period))
		g.Expect(pp).To(Equal(c.period), info(i, c.period))
		g.Expect(pp.Hours()).To(Equal(c.hours), info(i, c.period))
		g.Expect(pp.Minutes()).To(Equal(c.minutes), info(i, c.period))
		g.Expect(pp.Seconds()).To(Equal(c.seconds), info(i, c.period))

		pn := NewHMS(-c.hours, -c.minutes, -c.seconds)
		en := c.period.Negate()
		expectValid(t, pn, info(i, en))
		g.Expect(pn).To(Equal(en), info(i, en))
		g.Expect(pn.Hours()).To(Equal(-c.hours), info(i, en))
		g.Expect(pn.Minutes()).To(Equal(-c.minutes), info(i, en))
		g.Expect(pn.Seconds()).To(Equal(-c.seconds), info(i, en))
	}
}

//-------------------------------------------------------------------------------------------------

func TestNewYMD(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		period              Period
		years, months, days int
	}{
		// note: the negative cases are also covered (see below)

		{}, // zero case

		{period: Period{days: 1}, days: 1},
		{period: Period{months: 1}, months: 1},
		{period: Period{years: 1}, years: 1},
		{period: Period{years: 100, months: 222, days: 700}, years: 100, months: 222, days: 700},
		{period: Period{years: 32767, months: 32767, days: 32767}, years: 32767, months: 32767, days: 32767},
	}
	for i, c := range cases {
		pp := NewYMD(c.years, c.months, c.days)
		expectValid(t, pp, info(i, c.period))
		g.Expect(pp).To(Equal(c.period), info(i, c.period))
		g.Expect(pp.Years()).To(Equal(c.years), info(i, c.period))
		g.Expect(pp.Months()).To(Equal(c.months), info(i, c.period))
		g.Expect(pp.Days()).To(Equal(c.days), info(i, c.period))

		pn := NewYMD(-c.years, -c.months, -c.days)
		en := c.period.Negate()
		expectValid(t, pn, info(i, en))
		g.Expect(pn).To(Equal(en), info(i, en))
		g.Expect(pn.Years()).To(Equal(-c.years), info(i, en))
		g.Expect(pn.Months()).To(Equal(-c.months), info(i, en))
		g.Expect(pn.Days()).To(Equal(-c.days), info(i, en))
	}
}

//-------------------------------------------------------------------------------------------------

func TestNewOf(t *testing.T) {
	// note: the negative cases are also covered (see below)

	ms123 := time.Minute + 2*time.Second + 30*time.Millisecond
	// HMS tests
	testNewOf(t, 1, 10*time.Millisecond, Period{fraction: 1, fpart: Second}, true)
	testNewOf(t, 2, time.Second, Period{seconds: 1}, true)
	testNewOf(t, 3, time.Minute, Period{minutes: 1}, true)
	testNewOf(t, 4, time.Hour, Period{hours: 1}, true)
	testNewOf(t, 5, time.Hour+ms123, Period{hours: 1, minutes: 1, seconds: 2, fraction: 3, fpart: Second}, true)
	testNewOf(t, 6, 24*time.Hour+time.Minute+time.Second, Period{hours: 24, minutes: 1, seconds: 1}, true)
	testNewOf(t, 7, 32767*time.Hour+59*time.Minute+59*time.Second+990*time.Millisecond, Period{hours: 32767, minutes: 59, seconds: 59, fraction: 99, fpart: Second}, true)
	testNewOf(t, 8, 30*time.Minute+67*time.Second+450*time.Millisecond, Period{minutes: 31, seconds: 7, fraction: 45, fpart: Second}, true)

	// YMD tests: must be over 32767 hours (approx 45 months), otherwise HMS will take care of it
	// first rollover: >32767 hours
	testNewOf(t, 9, 32768*time.Hour+ms123, Period{days: 1365, hours: 8, minutes: 1, seconds: 2, fraction: 3, fpart: Second}, false)

	// second rollover: >32767 days
	testNewOf(t, 10, 24*32768*time.Hour+ms123, Period{years: 89, months: 8, days: 17}, false)
	testNewOf(t, 11, 36525*oneDay, Period{years: 100}, false)
}

func testNewOf(t *testing.T, i int, source time.Duration, expected Period, precise bool) {
	t.Helper()
	testNewOf1(t, i, source, expected, precise)
	testNewOf1(t, i, -source, expected.Negate(), precise)
}

func testNewOf1(t *testing.T, i int, source time.Duration, expected Period, precise bool) {
	t.Helper()
	g := NewGomegaWithT(t)

	n, p := NewOf(source)
	rev, _ := expected.Duration()
	info := fmt.Sprintf("%d: source %v expected %+v precise %v rev %v", i, source, expected, precise, rev)
	expectValid(t, n, info)
	g.Expect(n).To(Equal(expected), info)
	g.Expect(p).To(Equal(precise), info)
	if precise {
		g.Expect(rev).To(Equal(source), info)
	}
}

//-------------------------------------------------------------------------------------------------

func TestBetween(t *testing.T) {
	g := NewGomegaWithT(t)
	now := time.Now()

	cases := []struct {
		a, b     time.Time
		expected Period
	}{
		// note: the negative cases are also covered (see below)

		{now, now, Period{}},

		// simple positive date calculations
		{utc(2015, 1, 1, 0, 0, 0, 0), utc(2015, 1, 1, 0, 0, 0, 10), Period{seconds: 0, fraction: 1, fpart: Second}},
		{utc(2015, 1, 1, 0, 0, 0, 0), utc(2015, 2, 2, 1, 1, 1, 10), Period{days: 32, hours: 1, minutes: 1, seconds: 1, fraction: 1, fpart: Second}},
		{utc(2015, 2, 1, 0, 0, 0, 0), utc(2015, 3, 2, 1, 1, 1, 0), Period{days: 29, hours: 1, minutes: 1, seconds: 1}},
		{utc(2015, 3, 1, 0, 0, 0, 0), utc(2015, 4, 2, 1, 1, 1, 0), Period{days: 32, hours: 1, minutes: 1, seconds: 1}},
		{utc(2015, 4, 1, 0, 0, 0, 0), utc(2015, 5, 2, 1, 1, 1, 0), Period{days: 31, hours: 1, minutes: 1, seconds: 1}},
		{utc(2015, 5, 1, 0, 0, 0, 0), utc(2015, 6, 2, 1, 1, 1, 0), Period{days: 32, hours: 1, minutes: 1, seconds: 1}},
		{utc(2015, 6, 1, 0, 0, 0, 0), utc(2015, 7, 2, 1, 1, 1, 0), Period{days: 31, hours: 1, minutes: 1, seconds: 1}},
		{utc(2015, 1, 1, 0, 0, 0, 0), utc(2015, 7, 2, 1, 1, 1, 0), Period{days: 182, hours: 1, minutes: 1, seconds: 1}},

		//// less than one month
		{utc(2016, 1, 2, 0, 0, 0, 0), utc(2016, 2, 1, 0, 0, 0, 0), Period{days: 30}},
		{utc(2015, 2, 2, 0, 0, 0, 0), utc(2015, 3, 1, 0, 0, 0, 0), Period{days: 27}}, // non-leap
		{utc(2016, 2, 2, 0, 0, 0, 0), utc(2016, 3, 1, 0, 0, 0, 0), Period{days: 28}}, // leap year
		{utc(2016, 3, 2, 0, 0, 0, 0), utc(2016, 4, 1, 0, 0, 0, 0), Period{days: 30}},
		{utc(2016, 4, 2, 0, 0, 0, 0), utc(2016, 5, 1, 0, 0, 0, 0), Period{days: 29}},
		{utc(2016, 5, 2, 0, 0, 0, 0), utc(2016, 6, 1, 0, 0, 0, 0), Period{days: 30}},
		{utc(2016, 6, 2, 0, 0, 0, 0), utc(2016, 7, 1, 0, 0, 0, 0), Period{days: 29}},

		// BST drops an hour at the daylight-saving transition
		{utc(2015, 1, 1, 0, 0, 0, 0), bst(2015, 7, 2, 1, 1, 1, 10), Period{days: 182, minutes: 1, seconds: 1, fraction: 1, fpart: Second}},

		// daytime only
		{utc(2015, 1, 1, 2, 3, 4, 0), utc(2015, 1, 1, 2, 3, 4, 500), Period{fraction: 50, fpart: Second}},
		{utc(2015, 1, 1, 2, 3, 4, 0), utc(2015, 1, 1, 4, 4, 7, 500), Period{hours: 2, minutes: 1, seconds: 3, fraction: 50, fpart: Second}},
		{utc(2015, 1, 1, 2, 3, 4, 500), utc(2015, 1, 1, 4, 4, 7, 0), Period{hours: 2, minutes: 1, seconds: 2, fraction: 50, fpart: Second}},

		// different dates and times
		{utc(2015, 2, 1, 1, 0, 0, 0), utc(2015, 5, 30, 5, 6, 7, 0), Period{days: 118, hours: 4, minutes: 6, seconds: 7}},
		{utc(2015, 2, 1, 1, 0, 0, 0), bst(2015, 5, 30, 5, 6, 7, 0), Period{days: 118, hours: 3, minutes: 6, seconds: 7}},

		// earlier month in later year
		{utc(2015, 12, 22, 0, 0, 0, 0), utc(2016, 1, 10, 5, 6, 7, 0), Period{days: 19, hours: 5, minutes: 6, seconds: 7}},
		{utc(2015, 2, 11, 5, 6, 7, 500), utc(2016, 1, 10, 0, 0, 0, 0), Period{days: 332, hours: 18, minutes: 53, seconds: 52, fraction: 50, fpart: Second}},

		// larger ranges
		{utc(2009, 1, 1, 0, 0, 1, 0), utc(2016, 12, 31, 0, 0, 2, 0), Period{days: 2921, seconds: 1}},
		{utc(2000, 1, 1, 0, 0, 1, 0), utc(2016, 12, 31, 0, 0, 2, 0), Period{years: 0, months: 0, days: 6209, seconds: 1}},
		{utc(1900, 1, 1, 0, 0, 1, 0), utc(2009, 12, 31, 0, 0, 2, 0), Period{years: 109, months: 11, days: 30, seconds: 1}},
	}
	for i, c := range cases {
		pp := Between(c.a, c.b)
		g.Expect(pp).To(Equal(c.expected), info(i, c.expected))

		pn := Between(c.b, c.a)
		en := c.expected.Negate()
		g.Expect(pn).To(Equal(en), info(i, en))
	}
}

//-------------------------------------------------------------------------------------------------

func TestNormaliseUnchanged(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		source period64
	}{
		// note: the negative cases are also covered (see below)

		// zero case
		{period64{}},

		{period64{years: 1}},
		{period64{months: 1}},
		{period64{days: 1}},
		{period64{hours: 1}},
		{period64{minutes: 1}},
		{period64{seconds: 1}},

		{period64{years: 1, months: 1, days: 1, hours: 1, minutes: 1, seconds: 1, fraction: 1, fpart: Second}},

		{period64{days: 1, hours: 7}},
		{period64{days: 1, hours: 1, minutes: 1}},
		{period64{days: 1, hours: 1, seconds: 1}},
		{period64{months: 1, days: 1, hours: 1}},

		{period64{minutes: 1, seconds: 10}},
		{period64{hours: 1, minutes: 10}},
		{period64{years: 1, months: 7}},

		{period64{months: 1, fraction: 1, fpart: Month}},
		{period64{days: 1, fraction: 1, fpart: Day}},
		{period64{hours: 1, fraction: 1, fpart: Hour}},
		{period64{minutes: 1, fraction: 1, fpart: Minute}},
		{period64{seconds: 1, fraction: 1, fpart: Second}},

		// don't carry days to months...
		{period64{days: 32}},
		{period64{days: 32767}},

		// don't carry MaxInt16 - 1 where it would cause small arithmetic errors
		{period64{years: 32767}},
		{period64{days: 32767}},
	}
	for i, c := range cases {
		p, err := c.source.toPeriod()
		g.Expect(err).NotTo(HaveOccurred())

		testNormaliseBothSigns(t, i, c.source, p, true)
		testNormaliseBothSigns(t, i, c.source, p, false)
	}
}

//-------------------------------------------------------------------------------------------------

func TestNormaliseChanged(t *testing.T) {
	cases := []struct {
		source          period64
		precise, approx Period
	}{
		// note: the negative cases are also covered (see below)

		// carry seconds to minutes
		{period64{seconds: 70}, Period{minutes: 1, seconds: 10}, Period{minutes: 1, seconds: 10}},
		{period64{seconds: 699}, Period{minutes: 11, seconds: 39}, Period{minutes: 11, seconds: 39}},

		// carry minutes to hours
		{period64{minutes: 70}, Period{hours: 1, minutes: 10}, Period{hours: 1, minutes: 10}},
		{period64{minutes: 699}, Period{hours: 11, minutes: 39}, Period{hours: 11, minutes: 39}},

		// simplify 1 hour to minutes
		{period64{hours: 1, fraction: 25, fpart: Hour}, Period{hours: 1, minutes: 15}, Period{hours: 1, minutes: 15}},
		{period64{hours: 1, fraction: 75, fpart: Hour}, Period{hours: 1, minutes: 45}, Period{hours: 1, minutes: 45}},

		// carry hours to days
		{period64{hours: 48}, Period{hours: 48}, Period{days: 2}},
		{period64{hours: 49}, Period{hours: 49}, Period{days: 2, hours: 1}},
		{period64{hours: 32767}, Period{hours: 32767}, Period{days: 1365, hours: 7}},
		{period64{years: 1, months: 2, days: 3, hours: 32767}, Period{years: 1, months: 2, days: 3, hours: 32767}, Period{years: 1, months: 2, days: 1368, hours: 7}},
		{period64{hours: 32768}, Period{days: 1365, hours: 8}, Period{days: 1365, hours: 8}},
		{period64{years: 1, months: 2, days: 3, hours: 32768}, Period{years: 1, months: 2, days: 1368, hours: 8}, Period{years: 1, months: 2, days: 1368, hours: 8}},

		// carry months to years
		{period64{months: 12}, Period{years: 1}, Period{years: 1}},
		{period64{months: 13}, Period{years: 1, months: 1}, Period{years: 1, months: 1}},
		{period64{months: 25}, Period{years: 2, months: 1}, Period{years: 2, months: 1}},

		// carry days to prevent overflow
		{period64{days: 32768}, Period{years: 89, months: 8, days: 17, hours: 22, minutes: 8}, Period{years: 89, months: 8, days: 17, hours: 22, minutes: 8}},

		// full ripple up
		{period64{months: 121, days: 305, hours: 240, minutes: 60, seconds: 61}, Period{years: 10, months: 1, days: 305, hours: 241, minutes: 1, seconds: 1}, Period{years: 10, months: 1, days: 315, hours: 1, minutes: 1, seconds: 1}},

		// carry years to months
		{period64{years: 1}, Period{years: 1}, Period{years: 1}},
		{period64{years: 1, fraction: 75, fpart: Year}, Period{years: 1, months: 9}, Period{years: 1, months: 9}},
		{period64{years: 1, months: 7}, Period{years: 1, months: 7}, Period{years: 1, months: 7}},
	}
	for i, c := range cases {
		testNormaliseBothSigns(t, i, c.source, c.precise, true)
		testNormaliseBothSigns(t, i, c.source, c.approx, false)
	}
}

func testNormaliseBothSigns(t *testing.T, i int, source period64, expected Period, precise bool) {
	g := NewGomegaWithT(t)
	t.Helper()

	sstr := source.String()
	n1, err := source.normalise64(precise).toPeriod()
	info1 := fmt.Sprintf("%d: %s.Normalise(%v) expected %s to equal %s", i, sstr, precise, n1, expected)
	expectValid(t, n1, info1)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(n1).To(Equal(expected), info1)

	source.neg = !source.neg
	eneg := expected.Negate()
	n2, err := source.normalise64(precise).toPeriod()
	info2 := fmt.Sprintf("%d: %s.Normalise(%v) expected %s to equal %s", i, sstr, precise, n2, eneg)
	expectValid(t, n2, info2)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(n2).To(Equal(eneg), info2)
}

//-------------------------------------------------------------------------------------------------

// FIXME
func TestNormaliseWithBorrow(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		source          period64
		precise, approx Period
	}{
		// borrow seconds from minutes
		//{period64{days: 2, hours: -250}, Period{days: 1, hours: 23, minutes: 59, seconds: 10}, Period{days: 1, hours: 23, minutes: 59, seconds: 10}},
		//{period64{days: 2, seconds: -50}, Period{days: 1, hours: 23, minutes: 59, seconds: 10}, Period{days: 1, hours: 23, minutes: 59, seconds: 10}},
		//{period64{minutes: 2, seconds: -50}, Period{minutes: 1, seconds: 10}, Period{minutes: 1, seconds: 10}},
		//{period64{minutes: 2, seconds: -70}, Period{seconds: 50}, Period{seconds: 50}},
		//{period64{hours: 2, seconds: -50}, Period{hours: 1, minutes: 59, seconds: 10}, Period{hours: 1, minutes: 59, seconds: 10}},
	}
	for i, c := range cases {
		p1 := c.source // copy before normalise - note the pointer receiver
		n1, err := p1.normalise64(true).toPeriod()
		info1 := fmt.Sprintf("%d: %s.Normalise(true) expected %s to equal %s", i, c.source, n1, c.precise)
		expectValid(t, n1, info1)
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(n1).To(Equal(c.precise), info1)

		p2 := c.source
		n2, err := p2.normalise64(false).toPeriod()
		info2 := fmt.Sprintf("%d: %s.Normalise(false) expected %s to equal %s", i, c.source, n2, c.approx)
		expectValid(t, n2, info2)
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(n2).To(Equal(c.approx), info2)
	}
}

//-------------------------------------------------------------------------------------------------

func TestSimplify(t *testing.T) {
	cases := []struct {
		source          Period
		precise, approx Period
	}{
		// note: the negative cases are also covered (see below)

		// simplify 1 minute to seconds
		{Period{minutes: 1}, Period{minutes: 1}, Period{minutes: 1}},
		{Period{minutes: 1, seconds: 31}, Period{minutes: 1, seconds: 31}, Period{minutes: 1, seconds: 31}},
		{Period{minutes: 1, seconds: 30}, Period{seconds: 90}, Period{seconds: 90}},
		{Period{minutes: 1, seconds: 30, fraction: 1, fpart: Second}, Period{seconds: 90, fraction: 1, fpart: Second}, Period{seconds: 90, fraction: 1, fpart: Second}},

		// simplify 1 hour to minutes
		{Period{hours: 1}, Period{hours: 1}, Period{hours: 1}},
		{Period{hours: 1, minutes: 10}, Period{minutes: 70}, Period{minutes: 70}},
		{Period{hours: 1, minutes: 11}, Period{hours: 1, minutes: 11}, Period{hours: 1, minutes: 11}},
		{Period{hours: 1, minutes: 10, fraction: 1, fpart: Minute}, Period{minutes: 70, fraction: 1, fpart: Minute}, Period{minutes: 70, fraction: 1, fpart: Minute}},

		// simplify days
		{Period{days: 1, hours: 6}, Period{days: 1, hours: 6}, Period{hours: 30}},
		{Period{days: 1, hours: 7}, Period{days: 1, hours: 7}, Period{days: 1, hours: 7}},
		{Period{days: 1, hours: 6, fraction: 1, fpart: Hour}, Period{days: 1, hours: 6, fraction: 1, fpart: Hour}, Period{hours: 30, fraction: 1, fpart: Hour}},

		// simplify months
		{Period{years: 1}, Period{years: 1}, Period{years: 1}},
		{Period{years: 1, months: 9}, Period{months: 21}, Period{months: 21}},
		{Period{years: 1, months: 10}, Period{years: 1, months: 10}, Period{years: 1, months: 10}},
		{Period{years: 1, months: 9, fraction: 1, fpart: Month}, Period{months: 21, fraction: 1, fpart: Month}, Period{months: 21, fraction: 1, fpart: Month}},

		// fractional years don't simplify
		{Period{years: 1, fraction: 1, fpart: Year}, Period{years: 1, fraction: 1, fpart: Year}, Period{years: 1, fraction: 1, fpart: Year}},

		// discard proper fractions
		{Period{years: 10, fraction: 1, fpart: Month}, Period{years: 10, fraction: 1, fpart: Month}, Period{years: 10}},

		{Period{years: 1, fraction: 1, fpart: Day}, Period{years: 1, fraction: 1, fpart: Day}, Period{years: 1}},
		{Period{months: 12, fraction: 1, fpart: Day}, Period{months: 12, fraction: 1, fpart: Day}, Period{months: 12}},

		{Period{years: 1, fraction: 1, fpart: Hour}, Period{years: 1, fraction: 1, fpart: Hour}, Period{years: 1}},
		{Period{months: 1, fraction: 1, fpart: Hour}, Period{months: 1, fraction: 1, fpart: Hour}, Period{months: 1}},
		{Period{days: 30, fraction: 1, fpart: Hour}, Period{days: 30, fraction: 1, fpart: Hour}, Period{days: 30}},

		{Period{years: 1, fraction: 1, fpart: Minute}, Period{years: 1, fraction: 1, fpart: Minute}, Period{years: 1}},
		{Period{months: 1, fraction: 1, fpart: Minute}, Period{months: 1, fraction: 1, fpart: Minute}, Period{months: 1}},
		{Period{days: 1, fraction: 1, fpart: Minute}, Period{days: 1, fraction: 1, fpart: Minute}, Period{days: 1}},

		{Period{years: 1, fraction: 1, fpart: Second}, Period{years: 1, fraction: 1, fpart: Second}, Period{years: 1}},
		{Period{months: 1, fraction: 1, fpart: Second}, Period{months: 1, fraction: 1, fpart: Second}, Period{months: 1}},
		{Period{days: 1, fraction: 1, fpart: Second}, Period{days: 1, fraction: 1, fpart: Second}, Period{days: 1}},
		{Period{hours: 1, fraction: 1, fpart: Second}, Period{hours: 1, fraction: 1, fpart: Second}, Period{hours: 1}},
	}
	for i, c := range cases {
		testSimplifyBothSigns(t, i, c.source, c.precise, true)
		testSimplifyBothSigns(t, i, c.source, c.approx, false)
	}

	g := NewGomegaWithT(t)
	g.Expect(Period{days: 1, hours: 7}.Simplify(false, 6, 7, 30)).To(Equal(Period{hours: 31}))
	g.Expect(Period{hours: 1, minutes: 30}.Simplify(true, 6, 30)).To(Equal(Period{minutes: 90}))
	g.Expect(Period{years: 1, months: 11}.Simplify(true, 11)).To(Equal(Period{months: 23}))
	g.Expect(Period{years: 1, months: 6}.Simplify(true)).To(Equal(Period{months: 18}))
}

func testSimplifyBothSigns(t *testing.T, i int, source Period, expected Period, precise bool) {
	g := NewGomegaWithT(t)
	t.Helper()

	sstr := source.String()
	n1 := source.Simplify(precise, 9, 6, 10, 30)
	info1 := fmt.Sprintf("%d: %s.Simplify(%v) expected %s to equal %s", i, sstr, precise, n1, expected)
	expectValid(t, n1, info1)
	g.Expect(n1).To(Equal(expected), info1)

	eneg := expected.Negate()
	n2 := source.Negate().Simplify(precise, 9, 6, 10, 30)
	info2 := fmt.Sprintf("%d: %s.Simplify(%v) expected %s to equal %s", i, sstr, precise, n2, eneg)
	expectValid(t, n2, info2)
	g.Expect(n2).To(Equal(eneg), info2)
}

func TestPeriodFormat(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		period  string
		expectW string
		expectD string
	}{
		// note: the negative cases are also covered (see below)

		{"P0D", "0 days", ""},

		{"P1Y1M7D", "1 year, 1 month, 1 week", "1 year, 1 month, 7 days"},
		{"P1Y1M1W1D", "1 year, 1 month, 1 week, 1 day", "1 year, 1 month, 8 days"},
		{"PT1H1M1S", "1 hour, 1 minute, 1 second", ""},
		{"P1Y1M1W1DT1H1M1S", "1 year, 1 month, 1 week, 1 day, 1 hour, 1 minute, 1 second", ""},
		{"P3Y6M39DT2H7M9S", "3 years, 6 months, 5 weeks, 4 days, 2 hours, 7 minutes, 9 seconds", ""},
		{"P365D", "52 weeks, 1 day", ""},

		{"P1Y", "1 year", ""},
		{"P3Y", "3 years", ""},
		{"P1.1Y", "1.1 years", ""},
		{"P2.5Y", "2 years, 6 months", ""},
		{"P2.6Y", "2.6 years", ""},
		{"P2.15Y", "2.15 years", ""},
		{"P2.125Y", "2.12 years", ""},

		{"P1M", "1 month", ""},
		{"P6M", "6 months", ""},
		{"P1.1M", "1.1 months", ""},
		{"P2.5M", "2.5 months", ""},
		{"P2.15M", "2.15 months", ""},
		{"P2.125M", "2.12 months", ""},

		{"P1W", "1 week", "7 days"},
		{"P1.1W", "1 week, 0.7 day", "7.7 days"},
		{"P7D", "1 week", "7 days"},
		{"P35D", "5 weeks", "35 days"},
		{"P1D", "1 day", "1 day"},
		{"P4D", "4 days", "4 days"},
		{"P1.1D", "1.1 days", ""},
		{"P2.5D", "2.5 days", ""},
		{"P2.15D", "2.15 days", ""},
		{"P2.125D", "2.12 days", ""},

		{"PT1H", "1 hour", ""},
		{"PT1.1H", "1 hour, 6 minutes", ""},
		{"PT2.5H", "2 hours, 30 minutes", ""},
		{"PT2.15H", "2 hours, 9 minutes", ""},
		{"PT2.125H", "2.12 hours", ""},

		{"PT1M", "1 minute", ""},
		{"PT1.1M", "1 minute, 6 seconds", ""},
		{"PT2.5M", "2 minutes, 30 seconds", ""},
		{"PT2.15M", "2 minutes, 9 seconds", ""},
		{"PT2.125M", "2.12 minutes", ""},

		{"PT1S", "1 second", ""},
		{"PT1.1S", "1.1 seconds", ""},
		{"PT2.5S", "2.5 seconds", ""},
		{"PT2.15S", "2.15 seconds", ""},
		{"PT2.125S", "2.12 seconds", ""},
	}
	for i, c := range cases {
		p := MustParse(c.period)
		sp := p.Format()
		g.Expect(sp).To(Equal(c.expectW), info(i, "%s -> %s", p, c.expectW))

		en := p.Negate()
		sn := en.Format()
		g.Expect(sn).To(Equal(c.expectW), info(i, "%s -> %s", en, c.expectW))

		if c.expectD != "" {
			s := MustParse(c.period).FormatWithoutWeeks()
			g.Expect(s).To(Equal(c.expectD), info(i, "%s -> %s", p, c.expectD))
		}
	}
}

//-------------------------------------------------------------------------------------------------

func TestPeriodOnlyYMD(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		one    string
		expect string
	}{
		{"P1Y2M3DT4H5M6S", "P1Y2M3D"},
		{"-P6Y5M4DT3H2M1S", "-P6Y5M4D"},
	}
	for i, c := range cases {
		s := MustParse(c.one).OnlyYMD()
		g.Expect(s).To(Equal(MustParse(c.expect)), info(i, c.expect))
	}
}

func TestPeriodOnlyHMS(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		one    string
		expect string
	}{
		{"P1Y2M3DT4H5M6S", "PT4H5M6S"},
		{"-P6Y5M4DT3H2M1S", "-PT3H2M1S"},
	}
	for i, c := range cases {
		s := MustParse(c.one).OnlyHMS()
		g.Expect(s).To(Equal(MustParse(c.expect)), info(i, c.expect))
	}
}

func utc(year int, month time.Month, day, hour, min, sec, msec int) time.Time {
	return time.Date(year, month, day, hour, min, sec, msec*int(time.Millisecond), time.UTC)
}

func bst(year int, month time.Month, day, hour, min, sec, msec int) time.Time {
	return time.Date(year, month, day, hour, min, sec, msec*int(time.Millisecond), london)
}

var london *time.Location // UTC + 1 hour during summer

func init() {
	london, _ = time.LoadLocation("Europe/London")
}

func info(i int, m ...interface{}) string {
	if s, ok := m[0].(string); ok {
		m[0] = i
		return fmt.Sprintf("%d "+s, m...)
	}
	return fmt.Sprintf("%d %v", i, m[0])
}

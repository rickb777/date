// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package period

import (
	"fmt"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	"github.com/rickb777/plural"
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
		_, ep := ParseWithNormalise(c.value, c.normalise)
		g.Expect(ep).To(HaveOccurred(), info(i, c.value))
		g.Expect(ep.Error()).To(Equal(c.expvalue+c.expected), info(i, c.value))

		_, en := ParseWithNormalise("-"+c.value, c.normalise)
		g.Expect(en).To(HaveOccurred(), info(i, c.value))
		if c.expvalue != "" {
			g.Expect(en.Error()).To(Equal("-"+c.expvalue+c.expected), info(i, c.value))
		} else {
			g.Expect(en.Error()).To(Equal(c.expected), info(i, c.value))
		}
	}
}

func TestParsePeriodWithNormalise(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		value    string
		reversed string
		period   Period
	}{
		// all rollovers
		{"PT1234.5S", "PT20M34.5S", Period{minutes: 20, seconds: 34, fraction: 50, fpart: Second}},
		{"PT1234.5M", "PT20H34.5M", Period{hours: 20, minutes: 34, fraction: 50, fpart: Minute}},
		{"PT12345.6H", "PT12345.6H", Period{hours: 12345, fraction: 60, fpart: Hour}},
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
		g.Expect(err).NotTo(HaveOccurred(), info(i, c.value))
		g.Expect(p).To(Equal(c.period), info(i, c.value))
		// reversal is expected not to be an identity
		g.Expect(p.String()).To(Equal(c.reversed), info(i, c.value)+" reversed")
	}
}

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
		p, err := ParseWithNormalise(c.value, false)
		g.Expect(err).NotTo(HaveOccurred(), info(i, c.value))
		g.Expect(p).To(Equal(c.period), info(i, c.value))
		// reversal is usually expected to be an identity
		g.Expect(p.String()).To(Equal(c.reversed), info(i, c.value)+" reversed")
	}
}

func TestPeriodString(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		value  string
		period Period
	}{
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
		// negative
		{"-P0.01Y", Period{fraction: -1, fpart: Year}},
		{"-P0.01M", Period{fraction: -1, fpart: Month}},
		{"-P0.07D", Period{fraction: -7, fpart: Day}},
		{"-P0.01D", Period{fraction: -1, fpart: Day}},
		{"-PT0.01H", Period{fraction: -1, fpart: Hour}},
		{"-PT0.01M", Period{fraction: -1, fpart: Minute}},
		{"-PT0.01S", Period{fraction: -1, fpart: Second}},

		{"P3Y", Period{years: 3}},
		{"-P3Y", Period{years: -3}},
		{"P6M", Period{months: 6}},
		{"-P6M", Period{months: -6}},
		{"P5W", Period{days: 35}},
		{"-P5W", Period{days: -35}},
		{"P4W", Period{days: 28}},
		{"-P4W", Period{days: -28}},
		{"P4D", Period{days: 4}},
		{"-P4D", Period{days: -4}},
		{"PT12H", Period{hours: 12}},
		{"PT30M", Period{minutes: 30}},
		{"PT5S", Period{seconds: 5}},
		{"P3Y6M39DT1H2M4.09S", Period{years: 3, months: 6, days: 39, hours: 1, minutes: 2, seconds: 4, fraction: 9, fpart: Second}},
		{"-P3Y6M39DT1H2M4.09S", Period{years: -3, months: -6, days: -39, hours: -1, minutes: -2, seconds: -4, fraction: -9, fpart: Second}},

		{"P2.5Y", Period{years: 2, fraction: 50, fpart: Year}},
		{"P2.49Y", Period{years: 2, fraction: 49, fpart: Year}},
		{"P2.5M", Period{months: 2, fraction: 50, fpart: Month}},
		{"P2.5D", Period{days: 2, fraction: 50, fpart: Day}},
		{"PT2.5H", Period{hours: 2, fraction: 50, fpart: Hour}},
		{"PT2.5M", Period{minutes: 2, fraction: 50, fpart: Minute}},
		{"PT2.5S", Period{seconds: 2, fraction: 50, fpart: Second}},
	}
	for i, c := range cases {
		s := c.period.String()
		g.Expect(s).To(Equal(c.value), info(i, c.value))
	}
}

func TestPeriodIntComponents(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		value                      string
		y, m, w, d, dx, hh, mm, ss int
	}{
		{value: "P0D"},
		{value: "P1Y", y: 1},
		{value: "-P1Y", y: -1},
		{value: "P1W", w: 1, d: 7},
		{value: "-P1W", w: -1, d: -7},
		{value: "P6M", m: 6},
		{value: "-P6M", m: -6},
		{value: "P12M", y: 1},
		{value: "-P12M", y: -1, m: -0},
		{value: "P39D", w: 5, d: 39, dx: 4},
		{value: "-P39D", w: -5, d: -39, dx: -4},
		{value: "P4D", d: 4, dx: 4},
		{value: "-P4D", d: -4, dx: -4},
		{value: "PT12H", hh: 12},
		{value: "PT60M", hh: 1},
		{value: "PT30M", mm: 30},
		{value: "PT5S", ss: 5},
	}
	for i, c := range cases {
		p := MustParse(c.value)
		g.Expect(p.Years()).To(Equal(c.y), info(i, c.value))
		g.Expect(p.Months()).To(Equal(c.m), info(i, c.value))
		g.Expect(p.Weeks()).To(Equal(c.w), info(i, c.value))
		g.Expect(p.Days()).To(Equal(c.d), info(i, c.value))
		g.Expect(p.ModuloDays()).To(Equal(c.dx), info(i, c.value))
		g.Expect(p.Hours()).To(Equal(c.hh), info(i, c.value))
		g.Expect(p.Minutes()).To(Equal(c.mm), info(i, c.value))
		g.Expect(p.Seconds()).To(Equal(c.ss), info(i, c.value))
	}
}

func TestPeriodFloatComponents(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		value                      Period
		y, m, w, d, dx, hh, mm, ss float32
	}{
		// note: the negative cases are also covered (see below)

		{}, // zero case

		// YMD cases
		{value: Period{years: 1}, y: 1},
		{value: Period{years: 1, fraction: 50, fpart: Year}, y: 1.5},
		{value: Period{years: 1, fraction: 1, fpart: Year}, y: 1.01},
		{value: Period{months: 1}, m: 1},
		{value: Period{months: 1, fraction: 50, fpart: Month}, m: 1.5},
		{value: Period{months: 1, fraction: 1, fpart: Month}, m: 1.01},
		{value: Period{months: 6}, m: 6},
		{value: Period{months: 12}, m: 12},
		{value: Period{days: 7}, w: 1, d: 7},
		{value: Period{days: 7, fraction: 70, fpart: Day}, w: 1.1, d: 7.7},
		{value: Period{days: 7, fraction: 1, fpart: Day}, w: 7.01 / 7, d: 7.01},
		{value: Period{days: 1}, w: 1.0 / 7, d: 1},
		{value: Period{days: 1, fraction: 10, fpart: Day}, w: 1.1 / 7, d: 1.1},
		{value: Period{days: 1, fraction: 1, fpart: Day}, w: 1.01 / 7, d: 1.01},
		{value: Period{days: 39}, w: 5.571429, d: 39, dx: 4},
		{value: Period{days: 4}, w: 0.5714286, d: 4, dx: 4},

		// HMS cases
		{value: Period{hours: 1, fraction: 10, fpart: Hour}, hh: 1.1},
		{value: Period{hours: 1, fraction: 1, fpart: Hour}, hh: 1.01},
		{value: Period{hours: 1, minutes: 6}, hh: 1, mm: 6},
		{value: Period{hours: 12}, hh: 12},
		{value: Period{minutes: 1, fraction: 10, fpart: Minute}, mm: 1.1},
		{value: Period{minutes: 1, fraction: 1, fpart: Minute}, mm: 1.01},
		{value: Period{minutes: 1, seconds: 6}, mm: 1, ss: 6},
		{value: Period{minutes: 30}, mm: 30},
		{value: Period{seconds: 1, fraction: 10, fpart: Second}, ss: 1.1},
		{value: Period{seconds: 1, fraction: 1, fpart: Second}, ss: 1.01},
		{value: Period{seconds: 5}, ss: 5},
	}
	for i, c := range cases {
		pp := c.value
		g.Expect(pp.YearsFloat()).To(Equal(c.y), info(i, c.value))
		g.Expect(pp.MonthsFloat()).To(Equal(c.m), info(i, c.value))
		g.Expect(pp.WeeksFloat()).To(Equal(c.w), info(i, c.value))
		g.Expect(pp.DaysFloat()).To(Equal(c.d), info(i, c.value))
		g.Expect(pp.HoursFloat()).To(Equal(c.hh), info(i, c.value))
		g.Expect(pp.MinutesFloat()).To(Equal(c.mm), info(i, c.value))
		g.Expect(pp.SecondsFloat()).To(Equal(c.ss), info(i, c.value))

		pn := c.value.Negate()
		g.Expect(pn.YearsFloat()).To(Equal(-c.y), info(i, c.value))
		g.Expect(pn.MonthsFloat()).To(Equal(-c.m), info(i, c.value))
		g.Expect(pn.WeeksFloat()).To(Equal(-c.w), info(i, c.value))
		g.Expect(pn.DaysFloat()).To(Equal(-c.d), info(i, c.value))
		g.Expect(pn.HoursFloat()).To(Equal(-c.hh), info(i, c.value))
		g.Expect(pn.MinutesFloat()).To(Equal(-c.mm), info(i, c.value))
		g.Expect(pn.SecondsFloat()).To(Equal(-c.ss), info(i, c.value))
	}
}

//TODO
func xTestPeriodAddToTime(t *testing.T) {
	g := NewGomegaWithT(t)

	const ms = 1000000
	const sec = 1000 * ms
	const min = 60 * sec
	const hr = 60 * min

	// A conveniently round number (14 July 2017 @ 2:40am UTC)
	var t0 = time.Unix(1500000000, 0).UTC()

	cases := []struct {
		value   string
		result  time.Time
		precise bool
	}{
		// precise cases
		{"P0D", t0, true},
		{"PT1S", t0.Add(sec), true},
		{"PT0.1S", t0.Add(100 * ms), true},
		{"-PT0.1S", t0.Add(-100 * ms), true},
		{"PT3276S", t0.Add(3276 * sec), true},
		{"PT1M", t0.Add(60 * sec), true},
		{"PT0.1M", t0.Add(6 * sec), true},
		{"PT3276M", t0.Add(3276 * min), true},
		{"PT1H", t0.Add(hr), true},
		{"PT0.1H", t0.Add(6 * min), true},
		{"PT3276H", t0.Add(3276 * hr), true},
		{"P1D", t0.AddDate(0, 0, 1), true},
		{"P3276D", t0.AddDate(0, 0, 3276), true},
		{"P1M", t0.AddDate(0, 1, 0), true},
		{"P3276M", t0.AddDate(0, 3276, 0), true},
		{"P1Y", t0.AddDate(1, 0, 0), true},
		{"-P1Y", t0.AddDate(-1, 0, 0), true},
		{"P3276Y", t0.AddDate(3276, 0, 0), true},   // near the upper limit of range
		{"-P3276Y", t0.AddDate(-3276, 0, 0), true}, // near the lower limit of range
		// approximate cases
		{"P0.1D", t0.Add(144 * min), false},
		{"-P0.1D", t0.Add(-144 * min), false},
		{"P0.1M", t0.Add(oneMonthApprox / 10), false},
		{"P0.1Y", t0.Add(oneYearApprox / 10), false},
		// after normalisation, this period is one month and 9.2 days
		{"-P0.1Y0.1M0.1D", t0.Add(-oneMonthApprox - (13248 * min)), false},
	}
	for i, c := range cases {
		p := MustParse(c.value)
		t1, prec := p.AddTo(t0)
		g.Expect(t1).To(Equal(c.result), info(i, c.value))
		g.Expect(prec).To(Equal(c.precise), info(i, c.value))
	}
}

func TestPeriodToDuration(t *testing.T) {
	cases := []struct {
		value    string
		duration time.Duration
		precise  bool
	}{
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

func TestSignPotisitveNegative(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		value    string
		positive bool
		negative bool
		sign     int
	}{
		{"P0D", false, false, 0},
		{"PT1S", true, false, 1},
		{"-PT1S", false, true, -1},
		{"PT1M", true, false, 1},
		{"-PT1M", false, true, -1},
		{"PT1H", true, false, 1},
		{"-PT1H", false, true, -1},
		{"P1D", true, false, 1},
		{"-P1D", false, true, -1},
		{"P1M", true, false, 1},
		{"-P1M", false, true, -1},
		{"P1Y", true, false, 1},
		{"-P1Y", false, true, -1},
	}
	for i, c := range cases {
		p := MustParse(c.value)
		g.Expect(p.IsPositive()).To(Equal(c.positive), info(i, c.value))
		g.Expect(p.IsNegative()).To(Equal(c.negative), info(i, c.value))
		g.Expect(p.Sign()).To(Equal(c.sign), info(i, c.value))
	}
}

func TestPeriodApproxDays(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		value      string
		approxDays int
	}{
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

func TestPeriodApproxMonths(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		value        string
		approxMonths int
	}{
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

func TestNewPeriod(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		period                                       Period
		years, months, days, hours, minutes, seconds int
	}{
		{}, // zero case

		{period: Period{seconds: 1}, seconds: 1},
		{period: Period{minutes: 1}, minutes: 1},
		{period: Period{hours: 1}, hours: 1},
		{period: Period{days: 1}, days: 1},
		{period: Period{months: 1}, months: 1},
		{period: Period{years: 1}, years: 1},
		{period: Period{years: 100, months: 222, days: 700}, years: 100, months: 222, days: 700},
	}
	for i, c := range cases {
		pp := New(c.years, c.months, c.days, c.hours, c.minutes, c.seconds)
		g.Expect(pp).To(Equal(c.period), info(i, c.period))
		g.Expect(pp.Years()).To(Equal(c.years), info(i, c.period))
		g.Expect(pp.Months()).To(Equal(c.months), info(i, c.period))
		g.Expect(pp.Days()).To(Equal(c.days), info(i, c.period))

		pn := New(-c.years, -c.months, -c.days, -c.hours, -c.minutes, -c.seconds)
		en := c.period.Negate()
		g.Expect(pn).To(Equal(en), info(i, en))
		g.Expect(pn.Years()).To(Equal(-c.years), info(i, en))
		g.Expect(pn.Months()).To(Equal(-c.months), info(i, en))
		g.Expect(pn.Days()).To(Equal(-c.days), info(i, en))
	}
}

func TestNewHMS(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		period                  Period
		hours, minutes, seconds int
	}{
		{}, // zero case

		{period: Period{seconds: 1}, seconds: 1},
		{period: Period{minutes: 1}, minutes: 1},
		{period: Period{hours: 1}, hours: 1},
		{period: Period{hours: 3, minutes: 4, seconds: 5}, hours: 3, minutes: 4, seconds: 5},
		{period: Period{hours: 32767, minutes: 32767, seconds: 32767}, hours: 32767, minutes: 32767, seconds: 32767},
	}
	for i, c := range cases {
		pp := NewHMS(c.hours, c.minutes, c.seconds)
		g.Expect(pp).To(Equal(c.period), info(i, c.period))
		g.Expect(pp.Hours()).To(Equal(c.hours), info(i, c.period))
		g.Expect(pp.Minutes()).To(Equal(c.minutes), info(i, c.period))
		g.Expect(pp.Seconds()).To(Equal(c.seconds), info(i, c.period))

		pn := NewHMS(-c.hours, -c.minutes, -c.seconds)
		en := c.period.Negate()
		g.Expect(pn).To(Equal(en), info(i, en))
		g.Expect(pn.Hours()).To(Equal(-c.hours), info(i, en))
		g.Expect(pn.Minutes()).To(Equal(-c.minutes), info(i, en))
		g.Expect(pn.Seconds()).To(Equal(-c.seconds), info(i, en))
	}
}

func TestNewYMD(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		period              Period
		years, months, days int
	}{
		{}, // zero case

		{period: Period{days: 1}, days: 1},
		{period: Period{months: 1}, months: 1},
		{period: Period{years: 1}, years: 1},
		{period: Period{years: 100, months: 222, days: 700}, years: 100, months: 222, days: 700},
		{period: Period{years: 32767, months: 32767, days: 32767}, years: 32767, months: 32767, days: 32767},
	}
	for i, c := range cases {
		pp := NewYMD(c.years, c.months, c.days)
		g.Expect(pp).To(Equal(c.period), info(i, c.period))
		g.Expect(pp.Years()).To(Equal(c.years), info(i, c.period))
		g.Expect(pp.Months()).To(Equal(c.months), info(i, c.period))
		g.Expect(pp.Days()).To(Equal(c.days), info(i, c.period))

		pn := NewYMD(-c.years, -c.months, -c.days)
		en := c.period.Negate()
		g.Expect(pn).To(Equal(en), info(i, en))
		g.Expect(pn.Years()).To(Equal(-c.years), info(i, en))
		g.Expect(pn.Months()).To(Equal(-c.months), info(i, en))
		g.Expect(pn.Days()).To(Equal(-c.days), info(i, en))
	}
}

//TODO
func xTestNewOf(t *testing.T) {
	// HMS tests
	testNewOf(t, 1, 10*time.Millisecond, Period{fraction: 1, fpart: Second}, true)
	testNewOf(t, 2, time.Second, Period{seconds: 1}, true)
	testNewOf(t, 3, time.Minute, Period{minutes: 1}, true)
	testNewOf(t, 4, time.Hour, Period{hours: 1}, true)
	testNewOf(t, 5, time.Hour+time.Minute+time.Second, Period{hours: 1, minutes: 1, seconds: 1}, true)
	testNewOf(t, 6, 24*time.Hour+time.Minute+time.Second, Period{hours: 24, minutes: 1, seconds: 1}, true)
	testNewOf(t, 7, 32767*time.Hour+59*time.Minute+59*time.Second+990*time.Millisecond, Period{hours: 32767, minutes: 59, seconds: 59, fraction: 99, fpart: Second}, true)
	testNewOf(t, 8, 30*time.Minute+67*time.Second+450*time.Millisecond, Period{minutes: 31, seconds: 7, fraction: 45, fpart: Second}, true)

	// YMD tests: must be over 32767 hours (approx 45 months), otherwise HMS will take care of it
	// first rollover: >32767 hours
	testNewOf(t, 9, 32768*time.Hour, Period{days: 1365, hours: 8}, false)

	// second rollover: >32767 days
	testNewOf(t, 10, 24*32768*time.Hour, Period{years: 89, months: 8, days: 17}, false)
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
	g.Expect(n).To(Equal(expected), info)
	g.Expect(p).To(Equal(precise), info)
	if precise {
		g.Expect(rev).To(Equal(source), info)
	}
}

//TODO
func xTestBetween(t *testing.T) {
	g := NewGomegaWithT(t)
	now := time.Now()

	cases := []struct {
		a, b     time.Time
		expected Period
	}{
		{now, now, Period{}},

		// simple positive date calculations
		{utc(2015, 1, 1, 0, 0, 0, 0), utc(2015, 1, 1, 0, 0, 0, 100), Period{seconds: 1}},
		{utc(2015, 1, 1, 0, 0, 0, 0), utc(2015, 2, 2, 1, 1, 1, 1), Period{days: 32, hours: 1, minutes: 1, seconds: 1}},
		{utc(2015, 2, 1, 0, 0, 0, 0), utc(2015, 3, 2, 1, 1, 1, 1), Period{days: 29, hours: 1, minutes: 1, seconds: 1}},
		{utc(2015, 3, 1, 0, 0, 0, 0), utc(2015, 4, 2, 1, 1, 1, 1), Period{days: 32, hours: 1, minutes: 1, seconds: 1}},
		{utc(2015, 4, 1, 0, 0, 0, 0), utc(2015, 5, 2, 1, 1, 1, 1), Period{days: 31, hours: 1, minutes: 1, seconds: 1}},
		{utc(2015, 5, 1, 0, 0, 0, 0), utc(2015, 6, 2, 1, 1, 1, 1), Period{days: 32, hours: 1, minutes: 1, seconds: 1}},
		{utc(2015, 6, 1, 0, 0, 0, 0), utc(2015, 7, 2, 1, 1, 1, 1), Period{days: 31, hours: 1, minutes: 1, seconds: 1}},
		{utc(2015, 1, 1, 0, 0, 0, 0), utc(2015, 7, 2, 1, 1, 1, 1), Period{days: 182, hours: 1, minutes: 1, seconds: 1}},

		// less than one month
		{utc(2016, 1, 2, 0, 0, 0, 0), utc(2016, 2, 1, 0, 0, 0, 0), Period{days: 300}},
		{utc(2015, 2, 2, 0, 0, 0, 0), utc(2015, 3, 1, 0, 0, 0, 0), Period{days: 270}}, // non-leap
		{utc(2016, 2, 2, 0, 0, 0, 0), utc(2016, 3, 1, 0, 0, 0, 0), Period{days: 280}}, // leap year
		{utc(2016, 3, 2, 0, 0, 0, 0), utc(2016, 4, 1, 0, 0, 0, 0), Period{days: 300}},
		{utc(2016, 4, 2, 0, 0, 0, 0), utc(2016, 5, 1, 0, 0, 0, 0), Period{days: 290}},
		{utc(2016, 5, 2, 0, 0, 0, 0), utc(2016, 6, 1, 0, 0, 0, 0), Period{days: 300}},
		{utc(2016, 6, 2, 0, 0, 0, 0), utc(2016, 7, 1, 0, 0, 0, 0), Period{days: 290}},

		// BST drops an hour at the daylight-saving transition
		{utc(2015, 1, 1, 0, 0, 0, 0), bst(2015, 7, 2, 1, 1, 1, 1), Period{days: 1820, minutes: 10, seconds: 10}},

		// negative date calculation
		{utc(2015, 1, 1, 0, 0, 0, 100), utc(2015, 1, 1, 0, 0, 0, 0), Period{seconds: -1}},
		{utc(2015, 6, 2, 0, 0, 0, 0), utc(2015, 5, 1, 0, 0, 0, 0), Period{days: -320}},
		{utc(2015, 6, 2, 1, 1, 1, 1), utc(2015, 5, 1, 0, 0, 0, 0), Period{days: -320, hours: -10, minutes: -10, seconds: -10}},

		// daytime only
		{utc(2015, 1, 1, 2, 3, 4, 0), utc(2015, 1, 1, 2, 3, 4, 500), Period{seconds: 5}},
		{utc(2015, 1, 1, 2, 3, 4, 0), utc(2015, 1, 1, 4, 4, 7, 500), Period{hours: 20, minutes: 10, seconds: 35}},
		{utc(2015, 1, 1, 2, 3, 4, 500), utc(2015, 1, 1, 4, 4, 7, 0), Period{hours: 20, minutes: 10, seconds: 25}},

		// different dates and times
		{utc(2015, 2, 1, 1, 0, 0, 0), utc(2015, 5, 30, 5, 6, 7, 0), Period{days: 1180, hours: 40, minutes: 60, seconds: 70}},
		{utc(2015, 2, 1, 1, 0, 0, 0), bst(2015, 5, 30, 5, 6, 7, 0), Period{days: 1180, hours: 30, minutes: 60, seconds: 70}},

		// earlier month in later year
		{utc(2015, 12, 22, 0, 0, 0, 0), utc(2016, 1, 10, 5, 6, 7, 0), Period{days: 190, hours: 50, minutes: 60, seconds: 70}},
		{utc(2015, 2, 11, 5, 6, 7, 500), utc(2016, 1, 10, 0, 0, 0, 0), Period{days: 3320, hours: 180, minutes: 530, seconds: 525}},

		// larger ranges
		{utc(2009, 1, 1, 0, 0, 1, 0), utc(2016, 12, 31, 0, 0, 2, 0), Period{days: 29210, seconds: 10}},
		{utc(2008, 1, 1, 0, 0, 1, 0), utc(2016, 12, 31, 0, 0, 2, 0), Period{years: 80, months: 110, days: 300, seconds: 10}},
	}
	for i, c := range cases {
		n := Between(c.a, c.b)
		g.Expect(n).To(Equal(c.expected), info(i, c.expected))
	}
}

func TestNormaliseUnchanged(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		source period64
	}{
		// zero case
		{period64{}},

		{period64{years: 1}},
		{period64{months: 1}},
		{period64{days: 1}},
		{period64{hours: 1}},
		{period64{minutes: 1}},
		{period64{seconds: 1}},

		{period64{years: 1, months: 1, days: 1, hours: 1, minutes: 1, seconds: 1, fraction: 1, fpart: Second}},

		{period64{minutes: 1, seconds: 10}},
		{period64{hours: 1, minutes: 10}},
		{period64{years: 1, months: 7}},

		{period64{years: 1, fraction: 1, fpart: Year}},
		{period64{months: 1, fraction: 1, fpart: Month}},
		{period64{days: 1, fraction: 1, fpart: Day}},
		{period64{hours: 1, fraction: 1, fpart: Hour}},
		{period64{minutes: 1, fraction: 1, fpart: Minute}},
		{period64{seconds: 1, fraction: 1, fpart: Second}},

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

func TestNormaliseChanged(t *testing.T) {
	cases := []struct {
		source          period64
		precise, approx Period
	}{
		// carry seconds to minutes
		{period64{seconds: 70}, Period{minutes: 1, seconds: 10}, Period{minutes: 1, seconds: 10}},
		{period64{seconds: 699}, Period{minutes: 11, seconds: 39}, Period{minutes: 11, seconds: 39}},

		// carry minutes to hours
		{period64{minutes: 70}, Period{hours: 1, minutes: 10}, Period{hours: 1, minutes: 10}},
		{period64{minutes: 699}, Period{hours: 11, minutes: 39}, Period{hours: 11, minutes: 39}},

		// simplify 1 minute to seconds
		{period64{minutes: 1, seconds: 9}, Period{seconds: 69}, Period{seconds: 69}},
		{period64{minutes: 1, seconds: 9, fraction: 1, fpart: Second}, Period{seconds: 69, fraction: 1, fpart: Second}, Period{seconds: 69, fraction: 1, fpart: Second}},

		// simplify 1 hour to minutes
		{period64{hours: 1, minutes: 9}, Period{minutes: 69}, Period{minutes: 69}},
		{period64{hours: 1, minutes: 9, fraction: 1, fpart: Minute}, Period{minutes: 69, fraction: 1, fpart: Minute}, Period{minutes: 69, fraction: 1, fpart: Minute}},

		// carry hours to days
		{period64{hours: 48}, Period{hours: 48}, Period{days: 2}},
		{period64{hours: 49}, Period{hours: 49}, Period{days: 2, hours: 1}},
		{period64{hours: 32767}, Period{hours: 32767}, Period{days: 1365, hours: 7}},
		{period64{years: 1, months: 2, days: 3, hours: 32767}, Period{years: 1, months: 2, days: 3, hours: 32767}, Period{years: 1, months: 2, days: 1368, hours: 7}},
		{period64{hours: 32768}, Period{days: 1365, hours: 8}, Period{days: 1365, hours: 8}},
		{period64{years: 1, months: 2, days: 3, hours: 32768}, Period{years: 1, months: 2, days: 1368, hours: 8}, Period{years: 1, months: 2, days: 1368, hours: 8}},

		// carry months to years
		{period64{months: 12}, Period{years: 1}, Period{years: 1}},
		{period64{months: 13}, Period{months: 13}, Period{months: 13}},
		{period64{months: 25}, Period{years: 2, months: 1}, Period{years: 2, months: 1}},

		// don't carry days to months...
		{period64{days: 32}, Period{days: 32}, Period{days: 32}},
		{period64{days: 32767}, Period{days: 32767}, Period{days: 32767}},

		// ...except to prevent overflow
		{period64{days: 32768}, Period{years: 89, months: 8, days: 17, hours: 22, minutes: 8}, Period{years: 89, months: 8, days: 17, hours: 22, minutes: 8}},

		// full ripple up
		{period64{months: 121, days: 305, hours: 239, minutes: 591, seconds: 601}, Period{years: 10, months: 1, days: 305, hours: 249, minutes: 1, seconds: 1}, Period{years: 10, months: 1, days: 315, hours: 9, minutes: 1, seconds: 1}},

		// carry years to months
		{period64{years: 1}, Period{years: 1}, Period{years: 1}},
		{period64{years: 1, months: 6}, Period{months: 18}, Period{months: 18}},
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
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(n1).To(Equal(expected), info1)

	source.neg = !source.neg
	eneg := expected.Negate()
	n2, err := source.normalise64(precise).toPeriod()
	info2 := fmt.Sprintf("%d: %s.Normalise(%v) expected %s to equal %s", i, sstr, precise, n2, eneg)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(n2).To(Equal(eneg), info2)
}

//TODO
func xTestPeriodFormat(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		period string
		expect string
	}{
		{"P0D", "0 days"},
		{"P1Y", "1 year"},
		{"P3Y", "3 years"},
		{"-P3Y", "3 years"},
		{"P1M", "1 month"},
		{"P6M", "6 months"},
		{"-P6M", "6 months"},
		{"P1W", "1 week"},
		{"-P1W", "1 week"},
		{"P7D", "1 week"},
		{"P35D", "5 weeks"},
		{"-P35D", "5 weeks"},
		{"P1D", "1 day"},
		{"P4D", "4 days"},
		{"-P4D", "4 days"},
		{"P1Y1M8D", "1 year, 1 month, 1 week, 1 day"},
		{"PT1H1M1S", "1 hour, 1 minute, 1 second"},
		{"P1Y1M8DT1H1M1S", "1 year, 1 month, 1 week, 1 day, 1 hour, 1 minute, 1 second"},
		{"P3Y6M39DT2H7M9S", "3 years, 6 months, 5 weeks, 4 days, 2 hours, 7 minutes, 9 seconds"},
		{"-P3Y6M39DT2H7M9S", "3 years, 6 months, 5 weeks, 4 days, 2 hours, 7 minutes, 9 seconds"},
		{"P1.1Y", "1.1 years"},
		{"P2.5Y", "2.5 years"},
		{"P2.15Y", "2.1 years"},
		{"P2.125Y", "2.1 years"},
	}
	for i, c := range cases {
		s := MustParse(c.period).Format()
		g.Expect(s).To(Equal(c.expect), info(i, c.expect))
	}
}

//TODO
func xTestPeriodScale(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		one    string
		m      float32
		expect string
	}{
		{"P0D", 2, "P0D"},
		{"P1D", 2, "P2D"},
		{"P1D", 0, "P0D"},
		{"P1D", 365, "P365D"},
		{"P1M", 2, "P2M"},
		{"P1M", 12, "P1Y"},
		//TODO {"P1Y3M", 1.0/15, "P1M"},
		{"P1Y", 2, "P2Y"},
		{"PT1H", 2, "PT2H"},
		{"PT1M", 2, "PT2M"},
		{"PT1S", 2, "PT2S"},
		{"P1D", 0.5, "P0.5D"},
		{"P1M", 0.5, "P0.5M"},
		{"P1Y", 0.5, "P0.5Y"},
		{"PT1H", 0.5, "PT0.5H"},
		{"PT1H", 0.1, "PT6M"},
		//TODO {"PT1H", 0.01, "PT36S"},
		{"PT1M", 0.5, "PT0.5M"},
		{"PT1S", 0.5, "PT0.5S"},
		{"PT1H", 1.0 / 3600, "PT1S"},
		{"P1Y2M3DT4H5M6S", 2, "P2Y4M6DT8H10M12S"},
		{"P2Y4M6DT8H10M12S", -0.5, "-P1Y2M3DT4H5M6S"},
		{"-P2Y4M6DT8H10M12S", 0.5, "-P1Y2M3DT4H5M6S"},
		{"-P2Y4M6DT8H10M12S", -0.5, "P1Y2M3DT4H5M6S"},
		{"PT1M", 60, "PT1H"},
		{"PT1S", 60, "PT1M"},
		{"PT1S", 86400, "PT24H"},
		{"PT1S", 86400000, "P1000D"},
		{"P365.5D", 10, "P10Y2.5D"},
		//{"P365.5D", 0.1, "P36DT12H"},
	}
	for i, c := range cases {
		s := MustParse(c.one).Scale(c.m)
		g.Expect(s).To(Equal(MustParse(c.expect)), info(i, c.expect))
	}
}

//TODO
func xTestPeriodAdd(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		one, two string
		expect   string
	}{
		{"P0D", "P0D", "P0D"},
		{"P1D", "P1D", "P2D"},
		{"P1M", "P1M", "P2M"},
		{"P1Y", "P1Y", "P2Y"},
		{"PT1H", "PT1H", "PT2H"},
		{"PT1M", "PT1M", "PT2M"},
		{"PT1S", "PT1S", "PT2S"},
		{"P1Y2M3DT4H5M6S", "P6Y5M4DT3H2M1S", "P7Y7M7DT7H7M7S"},
		{"P7Y7M7DT7H7M7S", "-P7Y7M7DT7H7M7S", "P0D"},
	}
	for i, c := range cases {
		s := MustParse(c.one).Add(MustParse(c.two))
		g.Expect(s).To(Equal(MustParse(c.expect)), info(i, c.expect))
	}
}

//TODO
func xTestPeriodFormatWithoutWeeks(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		period string
		expect string
	}{
		{"P0D", "0 days"},
		{"P1Y", "1 year"},
		{"P3Y", "3 years"},
		{"-P3Y", "3 years"},
		{"P1M", "1 month"},
		{"P6M", "6 months"},
		{"-P6M", "6 months"},
		{"P7D", "7 days"},
		{"P35D", "35 days"},
		{"-P35D", "35 days"},
		{"P1D", "1 day"},
		{"P4D", "4 days"},
		{"-P4D", "4 days"},
		{"P1Y1M1DT1H1M1S", "1 year, 1 month, 1 day, 1 hour, 1 minute, 1 second"},
		{"P3Y6M39DT2H7M9S", "3 years, 6 months, 39 days, 2 hours, 7 minutes, 9 seconds"},
		{"-P3Y6M39DT2H7M9S", "3 years, 6 months, 39 days, 2 hours, 7 minutes, 9 seconds"},
		{"P1.1Y", "1.1 years"},
		{"P2.5Y", "2.5 years"},
		{"P2.15Y", "2.1 years"},
		{"P2.125Y", "2.1 years"},
	}
	for i, c := range cases {
		s := MustParse(c.period).FormatWithPeriodNames(PeriodYearNames, PeriodMonthNames, plural.Plurals{}, PeriodDayNames,
			PeriodHourNames, PeriodMinuteNames, PeriodSecondNames)
		g.Expect(s).To(Equal(c.expect), info(i, c.expect))
	}
}

//TODO
func xTestPeriodParseOnlyYMD(t *testing.T) {
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

//TODO
func xTestPeriodParseOnlyHMS(t *testing.T) {
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

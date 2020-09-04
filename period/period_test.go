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
		{"XY", false, "expected 'P' period mark at the start: ", "XY"},
		{"PxY", false, "expected a number before the 'Y' designator: ", "PxY"},
		{"PxW", false, "expected a number before the 'W' designator: ", "PxW"},
		{"PxD", false, "expected a number before the 'D' designator: ", "PxD"},
		{"PTxH", false, "expected a number before the 'H' designator: ", "PTxH"},
		{"PTxM", false, "expected a number before the 'M' designator: ", "PTxM"},
		{"PTxS", false, "expected a number before the 'S' designator: ", "PTxS"},
		{"P1HT1M", false, "unexpected remaining components 1H: ", "P1HT1M"},
		{"PT1Y", false, "unexpected remaining components 1Y: ", "PT1Y"},
		{"P1S", false, "unexpected remaining components 1S: ", "P1S"},
		// integer overflow
		{"P32768Y", false, "integer overflow occurred in years: ", "P32768Y"},
		{"P32768M", false, "integer overflow occurred in months: ", "P32768M"},
		{"P32768D", false, "integer overflow occurred in days: ", "P32768D"},
		{"PT32768H", false, "integer overflow occurred in hours: ", "PT32768H"},
		{"PT32768M", false, "integer overflow occurred in minutes: ", "PT32768M"},
		{"PT32768S", false, "integer overflow occurred in seconds: ", "PT32768S"},
		{"PT32768H32768M32768S", false, "integer overflow occurred in hours,minutes,seconds: ", "PT32768H32768M32768S"},
		{"PT103412160000S", false, "integer overflow occurred in seconds: ", "PT103412160000S"},
		{"P39324M", true, "integer overflow occurred in years: ", "P39324M"},
		{"P1196900D", true, "integer overflow occurred in years: ", "P1196900D"},
		{"PT28725600H", true, "integer overflow occurred in years: ", "PT28725600H"},
		{"PT1723536000M", true, "integer overflow occurred in years: ", "PT1723536000M"},
		{"PT103412160000S", true, "integer overflow occurred in years: ", "PT103412160000S"},
	}
	for i, c := range cases {
		_, ep := ParseWithNormalise(c.value, c.normalise)
		g.Expect(ep).To(HaveOccurred(), info(i, c.value))
		g.Expect(ep.Error()).To(Equal(c.expected+c.expvalue), info(i, c.value))

		_, en := ParseWithNormalise("-"+c.value, c.normalise)
		g.Expect(en).To(HaveOccurred(), info(i, c.value))
		if c.expvalue != "" {
			g.Expect(en.Error()).To(Equal(c.expected+"-"+c.expvalue), info(i, c.value))
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
		{"PT1234.5S", "PT20M34.5S", Period{minutes: 200, seconds: 345}},
		{"PT1234.5M", "PT20H34.5M", Period{hours: 200, minutes: 345}},
		{"PT12345.6H", "P514DT9.6H", Period{days: 5140, hours: 96}},
		{"P3276.1D", "P8Y11M19.2D", Period{years: 80, months: 110, days: 192}},
		{"P1234.5M", "P102Y10.5M", Period{years: 1020, months: 105}},
		// largest possible number of seconds normalised only in hours, mins, sec
		{"PT11592000S", "PT3220H", Period{hours: 32200}},
		{"-PT11592000S", "-PT3220H", Period{hours: -32200}},
		{"PT11595599S", "PT3220H59M59S", Period{hours: 32200, minutes: 590, seconds: 590}},
		// largest possible number of seconds normalised only in days, hours, mins, sec
		{"PT283046400S", "P468W", Period{days: 32760}},
		{"-PT283046400S", "-P468W", Period{days: -32760}},
		{"PT43084443590S", "P1365Y3M2WT26H83M50S", Period{years: 13650, months: 30, days: 140, hours: 260, minutes: 830, seconds: 500}},
		{"PT103412159999S", "P3276Y11M29DT37H83M59S", Period{years: 32760, months: 110, days: 290, hours: 370, minutes: 830, seconds: 590}},
		{"PT283132799S", "P468WT23H59M59S", Period{days: 32760, hours: 230, minutes: 590, seconds: 590}},
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
		{"P1Y", "P1Y", Period{years: 10}},
		{"P1M", "P1M", Period{months: 10}},
		{"P1W", "P1W", Period{days: 70}},
		{"P1D", "P1D", Period{days: 10}},
		{"PT1H", "PT1H", Period{hours: 10}},
		{"PT1M", "PT1M", Period{minutes: 10}},
		{"PT1S", "PT1S", Period{seconds: 10}},
		// smallest
		{"P0.1Y", "P0.1Y", Period{years: 1}},
		{"-P0.1Y", "-P0.1Y", Period{years: -1}},
		{"P0.1M", "P0.1M", Period{months: 1}},
		{"-P0.1M", "-P0.1M", Period{months: -1}},
		{"P0.1D", "P0.1D", Period{days: 1}},
		{"-P0.1D", "-P0.1D", Period{days: -1}},
		{"PT0.1H", "PT0.1H", Period{hours: 1}},
		{"-PT0.1H", "-PT0.1H", Period{hours: -1}},
		{"PT0.1M", "PT0.1M", Period{minutes: 1}},
		{"-PT0.1M", "-PT0.1M", Period{minutes: -1}},
		{"PT0.1S", "PT0.1S", Period{seconds: 1}},
		{"-PT0.1S", "-PT0.1S", Period{seconds: -1}},
		// week special case: also not identity when reversed
		{"P0.1W", "P0.7D", Period{days: 7}},
		{"-P0.1W", "-P0.7D", Period{days: -7}},
		// largest
		{"PT3276.7S", "PT3276.7S", Period{seconds: 32767}},
		{"PT3276.7M", "PT3276.7M", Period{minutes: 32767}},
		{"PT3276.7H", "PT3276.7H", Period{hours: 32767}},
		{"P3276.7D", "P3276.7D", Period{days: 32767}},
		{"P3276.7M", "P3276.7M", Period{months: 32767}},
		{"P3276.7Y", "P3276.7Y", Period{years: 32767}},

		{"P3Y", "P3Y", Period{years: 30}},
		{"P6M", "P6M", Period{months: 60}},
		{"P5W", "P5W", Period{days: 350}},
		{"P4D", "P4D", Period{days: 40}},
		{"PT12H", "PT12H", Period{hours: 120}},
		{"PT30M", "PT30M", Period{minutes: 300}},
		{"PT25S", "PT25S", Period{seconds: 250}},
		{"PT30M67.6S", "PT30M67.6S", Period{minutes: 300, seconds: 676}},
		{"P2.Y", "P2Y", Period{years: 20}},
		{"P2.5Y", "P2.5Y", Period{years: 25}},
		{"P2.15Y", "P2.1Y", Period{years: 21}},
		{"P2.125Y", "P2.1Y", Period{years: 21}},
		{"P1Y2.M", "P1Y2M", Period{years: 10, months: 20}},
		{"P1Y2.5M", "P1Y2.5M", Period{years: 10, months: 25}},
		{"P1Y2.15M", "P1Y2.1M", Period{years: 10, months: 21}},
		{"P1Y2.125M", "P1Y2.1M", Period{years: 10, months: 21}},
		{"P3276.7Y", "P3276.7Y", Period{years: 32767}},
		{"-P3276.7Y", "-P3276.7Y", Period{years: -32767}},
		// others
		{"P3Y6M5W4DT12H40M5S", "P3Y6M39DT12H40M5S", Period{years: 30, months: 60, days: 390, hours: 120, minutes: 400, seconds: 50}},
		{"+P3Y6M5W4DT12H40M5S", "P3Y6M39DT12H40M5S", Period{years: 30, months: 60, days: 390, hours: 120, minutes: 400, seconds: 50}},
		{"-P3Y6M5W4DT12H40M5S", "-P3Y6M39DT12H40M5S", Period{years: -30, months: -60, days: -390, hours: -120, minutes: -400, seconds: -50}},
		{"P1Y14M35DT48H125M800S", "P1Y14M5WT48H125M800S", Period{years: 10, months: 140, days: 350, hours: 480, minutes: 1250, seconds: 8000}},
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
		{"P1Y", Period{years: 10}},
		{"P1M", Period{months: 10}},
		{"P1W", Period{days: 70}},
		{"P1D", Period{days: 10}},
		{"PT1H", Period{hours: 10}},
		{"PT1M", Period{minutes: 10}},
		{"PT1S", Period{seconds: 10}},
		// smallest
		{"P0.1Y", Period{years: 1}},
		{"P0.1M", Period{months: 1}},
		{"P0.7D", Period{days: 7}},
		{"P0.1D", Period{days: 1}},
		{"PT0.1H", Period{hours: 1}},
		{"PT0.1M", Period{minutes: 1}},
		{"PT0.1S", Period{seconds: 1}},
		// negative
		{"-P0.1Y", Period{years: -1}},
		{"-P0.1M", Period{months: -1}},
		{"-P0.7D", Period{days: -7}},
		{"-P0.1D", Period{days: -1}},
		{"-PT0.1H", Period{hours: -1}},
		{"-PT0.1M", Period{minutes: -1}},
		{"-PT0.1S", Period{seconds: -1}},

		{"P3Y", Period{years: 30}},
		{"-P3Y", Period{years: -30}},
		{"P6M", Period{months: 60}},
		{"-P6M", Period{months: -60}},
		{"P5W", Period{days: 350}},
		{"-P5W", Period{days: -350}},
		{"P4W", Period{days: 280}},
		{"-P4W", Period{days: -280}},
		{"P4D", Period{days: 40}},
		{"-P4D", Period{days: -40}},
		{"PT12H", Period{hours: 120}},
		{"PT30M", Period{minutes: 300}},
		{"PT5S", Period{seconds: 50}},
		{"P3Y6M39DT1H2M4S", Period{years: 30, months: 60, days: 390, hours: 10, minutes: 20, seconds: 40}},
		{"-P3Y6M39DT1H2M4S", Period{years: -30, months: -60, days: -390, hours: -10, minutes: -20, seconds: -40}},
		{"P2.5Y", Period{years: 25}},
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
		{value: Period{years: 10}, y: 1},
		{value: Period{years: 15}, y: 1.5},
		{value: Period{months: 10}, m: 1},
		{value: Period{months: 15}, m: 1.5},
		{value: Period{months: 60}, m: 6},
		{value: Period{months: 120}, m: 12},
		{value: Period{days: 70}, w: 1, d: 7},
		{value: Period{days: 77}, w: 1.1, d: 7.7},
		{value: Period{days: 10}, w: 1.0 / 7, d: 1},
		{value: Period{days: 11}, w: 1.1 / 7, d: 1.1},
		{value: Period{days: 390}, w: 5.571429, d: 39, dx: 4},
		{value: Period{days: 40}, w: 0.5714286, d: 4, dx: 4},

		// HMS cases
		{value: Period{hours: 11}, hh: 1.1},
		{value: Period{hours: 10, minutes: 60}, hh: 1, mm: 6},
		{value: Period{hours: 120}, hh: 12},
		{value: Period{minutes: 11}, mm: 1.1},
		{value: Period{minutes: 10, seconds: 60}, mm: 1, ss: 6},
		{value: Period{minutes: 300}, mm: 30},
		{value: Period{seconds: 11}, ss: 1.1},
		{value: Period{seconds: 50}, ss: 5},
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

func TestPeriodToDuration(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		value    string
		duration time.Duration
		precise  bool
	}{
		{"P0D", time.Duration(0), true},
		{"PT1S", 1 * time.Second, true},
		{"PT0.1S", 100 * time.Millisecond, true},
		{"-PT0.1S", -100 * time.Millisecond, true},
		{"PT3276S", 3276 * time.Second, true},
		{"PT1M", 60 * time.Second, true},
		{"PT0.1M", 6 * time.Second, true},
		{"PT3276M", 3276 * time.Minute, true},
		{"PT1H", 3600 * time.Second, true},
		{"PT0.1H", 360 * time.Second, true},
		{"PT3220H", 3220 * time.Hour, true},
		{"PT3221H", 3221 * time.Hour, false}, // threshold of normalisation wrapping
		// days, months and years conversions are never precise
		{"P1D", 24 * time.Hour, false},
		{"P0.1D", 144 * time.Minute, false},
		{"P3276D", 3276 * 24 * time.Hour, false},
		{"P1M", oneMonthApprox, false},
		{"P0.1M", oneMonthApprox / 10, false},
		{"P3276M", 3276 * oneMonthApprox, false},
		{"P1Y", oneYearApprox, false},
		{"-P1Y", -oneYearApprox, false},
		{"P3276Y", 3276 * oneYearApprox, false},   // near the upper limit of range
		{"-P3276Y", -3276 * oneYearApprox, false}, // near the lower limit of range
	}
	for i, c := range cases {
		p := MustParse(c.value)
		d1, prec := p.Duration()
		g.Expect(d1).To(Equal(c.duration), info(i, c.value))
		g.Expect(prec).To(Equal(c.precise), info(i, c.value))
		d2 := p.DurationApprox()
		if c.precise {
			g.Expect(d2).To(Equal(c.duration), info(i, c.value))
		}
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
		{"-P1Y", -365},
	}
	for i, c := range cases {
		p := MustParse(c.value)
		td := p.TotalDaysApprox()
		g.Expect(td).To(Equal(c.approxDays), info(i, c.value))
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
		{"-P1Y", -12},
		{"PT24H", 0},
		{"PT744H", 1},
	}
	for i, c := range cases {
		p := MustParse(c.value)
		td := p.TotalMonthsApprox()
		g.Expect(td).To(Equal(c.approxMonths), info(i, c.value))
	}
}

func TestNewPeriod(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		period                                       Period
		years, months, days, hours, minutes, seconds int
	}{
		{}, // zero case

		// positives
		{period: Period{seconds: 10}, seconds: 1},
		{period: Period{minutes: 10}, minutes: 1},
		{period: Period{hours: 10}, hours: 1},
		{period: Period{days: 10}, days: 1},
		{period: Period{months: 10}, months: 1},
		{period: Period{years: 10}, years: 1},
		{period: Period{1000, 2220, 7000, 0, 0, 0}, years: 100, months: 222, days: 700},
		// negatives
		{period: Period{seconds: -10}, seconds: -1},
		{period: Period{minutes: -10}, minutes: -1},
		{period: Period{hours: -10}, hours: -1},
		{period: Period{days: -10}, days: -1},
		{period: Period{months: -10}, months: -1},
		{period: Period{years: -10}, years: -1},
	}
	for i, c := range cases {
		p := New(c.years, c.months, c.days, c.hours, c.minutes, c.seconds)
		g.Expect(p).To(Equal(c.period), info(i, c.period))
		g.Expect(p.Years()).To(Equal(c.years), info(i, c.period))
		g.Expect(p.Months()).To(Equal(c.months), info(i, c.period))
		g.Expect(p.Days()).To(Equal(c.days), info(i, c.period))
	}
}

func TestNewHMS(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		period                  Period
		hours, minutes, seconds int
	}{
		{}, // zero case
		// postives
		{period: Period{seconds: 10}, seconds: 1},
		{period: Period{minutes: 10}, minutes: 1},
		{period: Period{hours: 10}, hours: 1},
		// negatives
		{period: Period{seconds: -10}, seconds: -1},
		{period: Period{minutes: -10}, minutes: -1},
		{period: Period{hours: -10}, hours: -1},
	}
	for i, c := range cases {
		p := NewHMS(c.hours, c.minutes, c.seconds)
		g.Expect(p).To(Equal(c.period), info(i, c.period))
		g.Expect(p.Hours()).To(Equal(c.hours), info(i, c.period))
		g.Expect(p.Minutes()).To(Equal(c.minutes), info(i, c.period))
		g.Expect(p.Seconds()).To(Equal(c.seconds), info(i, c.period))
	}
}

func TestNewYMD(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		period              Period
		years, months, days int
	}{
		{}, // zero case
		// positives
		{period: Period{days: 10}, days: 1},
		{period: Period{months: 10}, months: 1},
		{period: Period{years: 10}, years: 1},
		{period: Period{years: 1000, months: 2220, days: 7000}, years: 100, months: 222, days: 700},
		// negatives
		{period: Period{days: -10}, days: -1},
		{period: Period{months: -10}, months: -1},
		{period: Period{years: -10}, years: -1},
	}
	for i, c := range cases {
		p := NewYMD(c.years, c.months, c.days)
		g.Expect(p).To(Equal(c.period), info(i, c.period))
		g.Expect(p.Years()).To(Equal(c.years), info(i, c.period))
		g.Expect(p.Months()).To(Equal(c.months), info(i, c.period))
		g.Expect(p.Days()).To(Equal(c.days), info(i, c.period))
	}
}

func TestNewOf(t *testing.T) {
	// HMS tests
	testNewOf(t, 100*time.Millisecond, Period{seconds: 1}, true)
	testNewOf(t, time.Second, Period{seconds: 10}, true)
	testNewOf(t, time.Minute, Period{minutes: 10}, true)
	testNewOf(t, time.Hour, Period{hours: 10}, true)
	testNewOf(t, time.Hour+time.Minute+time.Second, Period{hours: 10, minutes: 10, seconds: 10}, true)
	testNewOf(t, 24*time.Hour+time.Minute+time.Second, Period{hours: 240, minutes: 10, seconds: 10}, true)
	testNewOf(t, 3276*time.Hour+59*time.Minute+59*time.Second, Period{hours: 32760, minutes: 590, seconds: 590}, true)
	testNewOf(t, 30*time.Minute+67*time.Second+600*time.Millisecond, Period{minutes: 310, seconds: 76}, true)

	// YMD tests: must be over 3276 hours (approx 4.5 months), otherwise HMS will take care of it
	// first rollover: >3276 hours
	testNewOf(t, 3277*time.Hour, Period{days: 1360, hours: 130}, false)
	testNewOf(t, 3288*time.Hour, Period{days: 1370}, false)
	testNewOf(t, 3289*time.Hour, Period{days: 1370, hours: 10}, false)
	testNewOf(t, 24*3276*time.Hour, Period{days: 32760}, false)

	// second rollover: >3276 days
	testNewOf(t, 24*3277*time.Hour, Period{years: 80, months: 110, days: 200}, false)
	testNewOf(t, 3277*oneDay, Period{years: 80, months: 110, days: 200}, false)
	testNewOf(t, 3277*oneDay+time.Hour+time.Minute+time.Second, Period{years: 80, months: 110, days: 200, hours: 10}, false)
	testNewOf(t, 36525*oneDay, Period{years: 1000}, false)
}

func testNewOf(t *testing.T, source time.Duration, expected Period, precise bool) {
	t.Helper()
	testNewOf1(t, source, expected, precise)
	testNewOf1(t, -source, expected.Negate(), precise)
}

func testNewOf1(t *testing.T, source time.Duration, expected Period, precise bool) {
	t.Helper()
	g := NewGomegaWithT(t)

	n, p := NewOf(source)
	rev, _ := expected.Duration()
	info := fmt.Sprintf("%v %+v %v %v", source, expected, precise, rev)
	g.Expect(n).To(Equal(expected), info)
	g.Expect(p).To(Equal(precise), info)
	//g.Expect(rev).To(Equal(source), info)
}

func TestBetween(t *testing.T) {
	g := NewGomegaWithT(t)
	now := time.Now()

	cases := []struct {
		a, b     time.Time
		expected Period
	}{
		{now, now, Period{0, 0, 0, 0, 0, 0}},

		// simple positive date calculations
		{utc(2015, 1, 1, 0, 0, 0, 0), utc(2015, 1, 1, 0, 0, 0, 100), Period{seconds: 1}},
		{utc(2015, 1, 1, 0, 0, 0, 0), utc(2015, 2, 2, 1, 1, 1, 1), Period{days: 320, hours: 10, minutes: 10, seconds: 10}},
		{utc(2015, 2, 1, 0, 0, 0, 0), utc(2015, 3, 2, 1, 1, 1, 1), Period{days: 290, hours: 10, minutes: 10, seconds: 10}},
		{utc(2015, 3, 1, 0, 0, 0, 0), utc(2015, 4, 2, 1, 1, 1, 1), Period{days: 320, hours: 10, minutes: 10, seconds: 10}},
		{utc(2015, 4, 1, 0, 0, 0, 0), utc(2015, 5, 2, 1, 1, 1, 1), Period{days: 310, hours: 10, minutes: 10, seconds: 10}},
		{utc(2015, 5, 1, 0, 0, 0, 0), utc(2015, 6, 2, 1, 1, 1, 1), Period{days: 320, hours: 10, minutes: 10, seconds: 10}},
		{utc(2015, 6, 1, 0, 0, 0, 0), utc(2015, 7, 2, 1, 1, 1, 1), Period{days: 310, hours: 10, minutes: 10, seconds: 10}},
		{utc(2015, 1, 1, 0, 0, 0, 0), utc(2015, 7, 2, 1, 1, 1, 1), Period{days: 1820, hours: 10, minutes: 10, seconds: 10}},

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

func TestNormalise(t *testing.T) {
	cases := []struct {
		source, precise, approx Period
	}{
		// zero cases
		{New(0, 0, 0, 0, 0, 0), New(0, 0, 0, 0, 0, 0), New(0, 0, 0, 0, 0, 0)},

		// carry seconds to minutes
		{Period{seconds: 699}, Period{minutes: 10, seconds: 99}, Period{minutes: 10, seconds: 99}},

		// carry minutes to seconds
		{Period{minutes: 5}, Period{seconds: 300}, Period{seconds: 300}},
		{Period{minutes: 1}, Period{seconds: 60}, Period{seconds: 60}},
		{Period{minutes: 55}, Period{minutes: 50, seconds: 300}, Period{minutes: 50, seconds: 300}},

		// carry minutes to hours
		{Period{minutes: 699}, Period{hours: 10, minutes: 90, seconds: 540}, Period{hours: 10, minutes: 90, seconds: 540}},

		// carry hours to minutes
		{Period{hours: 5}, Period{minutes: 300}, Period{minutes: 300}},

		// carry hours to days
		{Period{hours: 249}, Period{hours: 240, minutes: 540}, Period{hours: 240, minutes: 540}},
		{Period{hours: 249}, Period{hours: 240, minutes: 540}, Period{hours: 240, minutes: 540}},
		{Period{hours: 369}, Period{hours: 360, minutes: 540}, Period{days: 10, hours: 120, minutes: 540}},
		{Period{hours: 249, seconds: 10}, Period{hours: 240, minutes: 540, seconds: 10}, Period{hours: 240, minutes: 540, seconds: 10}},

		// carry days to hours
		{Period{days: 5, hours: 30}, Period{hours: 150}, Period{hours: 150}},

		// carry months to years
		{Period{months: 125}, Period{months: 125}, Period{months: 125}},
		{Period{months: 131}, Period{years: 10, months: 11}, Period{years: 10, months: 11}},

		// carry days to months
		{Period{days: 323}, Period{days: 323}, Period{days: 323}},

		// carry months to days
		{Period{months: 5, days: 203}, Period{days: 355}, Period{months: 10, days: 50}},

		// full ripple up
		{Period{months: 121, days: 305, hours: 239, minutes: 591, seconds: 601}, Period{years: 10, days: 330, hours: 360, minutes: 540, seconds: 61}, Period{years: 10, months: 10, days: 40, minutes: 540, seconds: 61}},

		// carry years to months
		{Period{years: 5}, Period{months: 60}, Period{months: 60}},
		{Period{years: 5, months: 25}, Period{months: 85}, Period{months: 85}},
		{Period{years: 5, months: 20, days: 10}, Period{months: 80, days: 10}, Period{months: 80, days: 10}},
	}
	for i, c := range cases {
		testNormaliseBothSigns(t, i, c.source, c.precise, true)
		testNormaliseBothSigns(t, i, c.source, c.approx, false)
	}
}

func testNormaliseBothSigns(t *testing.T, i int, source, expected Period, precise bool) {
	g := NewGomegaWithT(t)
	t.Helper()

	n1 := source.Normalise(precise)
	if n1 != expected {
		t.Errorf("%d: %v.Normalise(%v) %s\n   gives %-22s %#v %s,\n    want %-22s %#v %s",
			i, source, precise, source.DurationApprox(),
			n1, n1, n1.DurationApprox(),
			expected, expected, expected.DurationApprox())
	}

	sneg := source.Negate()
	eneg := expected.Negate()
	n2 := sneg.Normalise(precise)
	g.Expect(n2).To(Equal(eneg))
}

func TestPeriodFormat(t *testing.T) {
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

func TestPeriodFormatWithoutWeeks(t *testing.T) {
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

func TestPeriodParseOnlyYMD(t *testing.T) {
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

func TestPeriodParseOnlyHMS(t *testing.T) {
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

func info(i int, m interface{}) string {
	return fmt.Sprintf("%d %v", i, m)
}

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
var oneMonthApprox = 2629746 * time.Millisecond // 30.436875 days
var oneYearApprox = 31556952 * time.Second      // 365.2425 days

func TestParseErrors(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		value    string
		expected string
	}{
		{"", "cannot parse a blank string as a period"},
		{"XY", "expected 'P' period designator at the start: XY"},
		{"PxY", "expected a number not 'x': PxY"},
		{"PxW", "expected a number not 'x': PxW"},
		{"PxD", "expected a number not 'x': PxD"},
		{"PTxH", "expected a number not 'x': PTxH"},
		{"PTxS", "expected a number not 'x': PTxS"},
		{"P1HT1M", "'H' designator cannot occur here: P1HT1M"},
		{"PT1Y", "'Y' designator cannot occur here: PT1Y"},
		{"P1S", "'S' designator cannot occur here: P1S"},
	}
	for i, c := range cases {
		_, err := Parse(c.value)
		g.Expect(err).To(HaveOccurred(), info(i, c.value))
		g.Expect(err.Error()).To(Equal(c.expected), info(i, c.value))
	}
}

func TestParsePeriodWithNormalise(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		value  string
		period Period
	}{
		//{"P3Y", Period{3600, 0, 0}},
		//{"P6M", Period{600, 0, 0}},
		//{"P5W", Period{0, 3500, 0}},
		//{"P4D", Period{0, 400, 0}},
		//{"PT12H", Period{0, 0, 4320000}},
		//{"PT30M", Period{0, 0, 180000}},
		//{"PT25S", Period{0, 0, 2500}},
		//{"PT30M67.621S", Period{0, 0, 186762}},
		//{"P3Y6M5W4DT12H40M5S", Period{4200, 3900, 4560500}},
		//{"+P3Y6M5W4DT12H40M5S", Period{4200, 3900, 4560500}},
		//{"-P3Y6M5W4DT12H40M5S", Period{-4200, -3900, -4560500}},
		//{"P2.Y", Period{2400, 0, 0}},
		//{"P2.5Y", Period{3000, 0, 0}},
		//{"P2.15Y", Period{2580, 0, 0}},
		//{"P1.234Y", Period{1476, 0, 0}},
		//{"P1Y2.M", Period{1400, 0, 0}},
		//{"P1Y2.5M", Period{1450, 0, 0}},
		//{"P1Y2.15M", Period{1415, 0, 0}},
		//{"P1Y1.234M", Period{1323, 0, 0}},

		// days rollover
		{"P21474836.48D", Period{centiDays: 2147483647}},
	}
	for i, c := range cases {
		p, err := Parse(c.value)
		g.Expect(err).NotTo(HaveOccurred(), info(i, c.value))
		g.Expect(p).To(Equal(c.period), info(i, c.value))
	}
}

func TestParsePeriodWithoutNormalise(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		value  string
		period Period
	}{
		// zeroes
		{"P0", Period{}},
		{"P0Y", Period{}},
		{"P0M", Period{}},
		{"P0W", Period{}},
		{"P0D", Period{}},
		{"PT0H", Period{}},
		{"PT0M", Period{}},
		{"PT0S", Period{}},
		// ones
		{"P1Y", Period{1200, 0, 0, "P1Y"}},
		{"P1M", Period{100, 0, 0, "P1M"}},
		{"P1W", Period{0, 700, 0, "P1W"}},
		{"P1D", Period{0, 100, 0, "P1D"}},
		{"PT1H", Period{0, 0, 360000, "PT1H"}},
		{"PT1M", Period{0, 0, 6000, "PT1M"}},
		{"PT1S", Period{0, 0, 100, "PT1S"}},
		// smallest
		{"P0.01Y", Period{12, 0, 0, "P0.01Y"}},
		{"-P0.01Y", Period{-12, 0, 0, "-P0.01Y"}},
		{"P0.01M", Period{1, 0, 0, "P0.01M"}},
		{"-P0.01M", Period{-1, 0, 0, "-P0.01M"}},
		{"P0.01W", Period{0, 7, 0, "P0.01W"}},
		{"-P0.01W", Period{0, -7, 0, "-P0.01W"}},
		{"P0.01D", Period{0, 1, 0, "P0.01D"}},
		{"-P0.01D", Period{0, -1, 0, "-P0.01D"}},
		{"PT0.01H", Period{0, 0, 3600, "PT0.01H"}},
		{"-PT0.01H", Period{0, 0, -3600, "-PT0.01H"}},
		{"PT0.01M", Period{0, 0, 60, "PT0.01M"}},
		{"-PT0.01M", Period{0, 0, -60, "-PT0.01M"}},
		{"PT0.01S", Period{0, 0, 1, "PT0.01S"}},
		{"-PT0.01S", Period{0, 0, -1, "-PT0.01S"}},

		{"P1Y14M35DT48H125M800S", Period{2600, 3500, 18110000, "P1Y14M35DT48H125M800S"}},
		// largest possible number of seconds
		{"PT21474836.47S", Period{0, 0, 2147483647, "PT21474836.47S"}},
		{"-PT21474836.47S", Period{0, 0, -2147483647, "-PT21474836.47S"}},
		// largest possible number of days
		{"P21474836.47D", Period{0, 2147483647, 0, "P21474836.47D"}},
		{"-P21474836.47D", Period{0, -2147483647, 0, "-P21474836.47D"}},
		// largest possible number of months
		{"P21474836.47M", Period{2147483647, 0, 0, "P21474836.47M"}},
		{"-P21474836.47M", Period{-2147483647, 0, 0, "-P21474836.47M"}},
		// largest possible number of years
		{"P1789569.70Y", Period{2147483640, 0, 0, "P1789569.70Y"}},
		{"-P1789569.70Y", Period{-2147483640, 0, 0, "-P1789569.70Y"}},
	}
	for i, c := range cases {
		p, err := ParseWithNormalise(c.value, false)
		g.Expect(err).NotTo(HaveOccurred(), info(i, c.value))
		g.Expect(p).To(Equal(c.period), info(i, c.value))
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
		{"P1Y", Period{1200, 0, 0, "P1Y"}},
		{"P1M", Period{100, 0, 0, "P1M"}},
		{"P1W", Period{0, 700, 0, "P1W"}},
		{"P1D", Period{0, 100, 0, "P1D"}},
		{"PT1H", Period{0, 0, 360000, "PT1H"}},
		{"PT1M", Period{0, 0, 6000, "PT1M"}},
		{"PT1S", Period{0, 0, 100, "PT1S"}},
		// smallest
		{"P0.01M", Period{1, 0, 0, "P0.01M"}},
		{"-P0.01M", Period{-1, 0, 0, "-P0.01M"}},
		{"P0.12M", Period{12, 0, 0, "P0.12M"}},
		{"-P0.12M", Period{-12, 0, 0, "-P0.12M"}},
		{"P0.01D", Period{0, 1, 0, "P0.01D"}},
		{"-P0.01D", Period{0, -1, 0, "-P0.01D"}},
		{"P0.07D", Period{0, 7, 0, "P0.07D"}},
		{"-P0.07D", Period{0, -7, 0, "-P0.07D"}},
		{"PT0.01S", Period{0, 0, 1, "PT0.01S"}},
		{"-PT0.01S", Period{0, 0, -1, "-PT0.01S"}},
		{"PT3.6S", Period{0, 0, 360, "PT3.6S"}},
		{"-PT3.6S", Period{0, 0, -360, "-PT3.6S"}},
		{"PT0.6S", Period{0, 0, 60, "PT0.6S"}},
		{"-PT0.6S", Period{0, 0, -60, "-PT0.6S"}},
		// months
		{"P3Y", Period{3600, 0, 0, "P3Y"}},
		{"-P3Y", Period{-3600, 0, 0, "-P3Y"}},
		{"P6M", Period{600, 0, 0, "P6M"}},
		{"-P6M", Period{-600, 0, 0, "-P6M"}},
		{"P0.01M", Period{1, 0, 0, "P0.01M"}},
		{"-P0.01M", Period{-1, 0, 0, "-P0.01M"}},
		{"P1789569Y8.47M", Period{2147483647, 0, 0, "P1789569Y8.47M"}},
		{"-P1789569Y8.47M", Period{-2147483647, 0, 0, "-P1789569Y8.47M"}},
		// days
		{"P0.01D", Period{0, 1, 0, "P0.01D"}},
		{"-P0.01D", Period{0, -1, 0, "-P0.01D"}},
		{"P5W", Period{0, 3500, 0, "P5W"}},
		{"-P5W", Period{0, -3500, 0, "-P5W"}},
		{"P4D", Period{0, 400, 0, "P4D"}},
		{"-P4D", Period{0, -400, 0, "-P4D"}},
		// centiSeconds
		{"PT1S", Period{0, 0, 100, "PT1S"}},
		{"PT1M", Period{0, 0, 6000, "PT1M"}},
		{"PT1H", Period{0, 0, 360000, "PT1H"}},
		{"PT24H", Period{0, 0, 8640000, "PT24H"}},
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
		{"P0D", 0, 0, 0, 0, 0, 0, 0, 0},
		{"P1Y", 1, 0, 0, 0, 0, 0, 0, 0},
		{"-P1Y", -1, 0, 0, 0, 0, 0, 0, 0},
		{"P1W", 0, 0, 1, 7, 0, 0, 0, 0},
		{"-P1W", 0, 0, -1, -7, 0, 0, 0, 0},
		{"P6M", 0, 6, 0, 0, 0, 0, 0, 0},
		{"-P6M", 0, -6, 0, 0, 0, 0, 0, 0},
		{"P12M", 1, 0, 0, 0, 0, 0, 0, 0},
		{"-P12M", -1, -0, 0, 0, 0, 0, 0, 0},
		{"P39D", 0, 0, 5, 39, 4, 0, 0, 0},
		{"-P39D", 0, 0, -5, -39, -4, 0, 0, 0},
		{"P4D", 0, 0, 0, 4, 4, 0, 0, 0},
		{"-P4D", 0, 0, 0, -4, -4, 0, 0, 0},
		{"PT12H", 0, 0, 0, 0, 0, 12, 0, 0},
		{"PT60M", 0, 0, 0, 0, 0, 1, 0, 0},
		{"PT30M", 0, 0, 0, 0, 0, 0, 30, 0},
		{"PT5S", 0, 0, 0, 0, 0, 0, 0, 5},
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
		value                      string
		y, m, w, d, dx, hh, mm, ss float32
	}{
		// note: the negative cases are also covered (see below)

		{"P0D", 0, 0, 0, 0, 0, 0, 0, 0},

		// YMD cases
		{"P1Y", 1, 0, 0, 0, 0, 0, 0, 0},
		{"P1.5Y", 1, 6, 0, 0, 0, 0, 0, 0},
		{"P1M", 0, 1, 0, 0, 0, 0, 0, 0},
		{"P1.5M", 0, 1.5, 0, 0, 0, 0, 0, 0},
		{"P6M", 0, 6, 0, 0, 0, 0, 0, 0},
		{"P12M", 1, 0, 0, 0, 0, 0, 0, 0},
		{"P1W", 0, 0, 1, 7, 0, 0, 0, 0},
		{"P1.1W", 0, 0, 1.1, 7.7, 0, 0, 0, 0},
		{"P1D", 0, 0, 1.0 / 7, 1, 0, 0, 0, 0},
		{"P1.1D", 0, 0, 1.1 / 7, 1.1, 0, 0, 0, 0},
		{"P39D", 0, 0, 5.571429, 39, 4, 0, 0, 0},
		{"P4D", 0, 0, 0.5714286, 4, 4, 0, 0, 0},

		// HMS cases
		{"PT1.1H", 0, 0, 0, 0, 0, 1, 6, 0},
		{"PT12H", 0, 0, 0, 0, 0, 12, 0, 0},
		{"PT1.1M", 0, 0, 0, 0, 0, 0, 1, 6},
		{"PT30M", 0, 0, 0, 0, 0, 0, 30, 0},
		{"PT1.1S", 0, 0, 0, 0, 0, 0, 0, 1.1},
		{"PT5S", 0, 0, 0, 0, 0, 0, 0, 5},
	}
	for i, c := range cases {
		pp := MustParse(c.value)
		g.Expect(pp.YearsFloat()).To(Equal(c.y), info(i, c.value))
		g.Expect(pp.MonthsFloat()).To(Equal(c.m), info(i, c.value))
		g.Expect(pp.WeeksFloat()).To(Equal(c.w), info(i, c.value))
		g.Expect(pp.DaysFloat()).To(Equal(c.d), info(i, c.value))
		g.Expect(pp.HoursFloat()).To(Equal(c.hh), info(i, c.value))
		g.Expect(pp.MinutesFloat()).To(Equal(c.mm), info(i, c.value))
		g.Expect(pp.SecondsFloat()).To(Equal(c.ss), info(i, c.value))

		pn := MustParse("-" + c.value)
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
		{"P0.001M", t0.Add(oneMonthApprox), false},
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

//TODO
func xTestPeriodToDuration(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		value    string
		duration time.Duration
		precise  bool
	}{
		// hours, minutes seconds conversions are always precise
		{"P0D", time.Duration(0), true},
		{"PT1S", 1 * time.Second, true},
		{"PT0.001S", time.Millisecond, true},
		{"PT1M", 60 * time.Second, true},
		{"PT0.001M", 60 * time.Millisecond, true},
		{"PT1H", 3600 * time.Second, true},
		{"PT0.001H", 3600 * time.Millisecond, true},
		{"PT2147483.647S", 2147483647 * time.Millisecond, true},
		{"PT35791.394M", 2147483640 * time.Millisecond, true},
		{"PT596.523H", 2147482800 * time.Millisecond, true},
		// days, months and years conversions are never precise
		{"P1D", 24 * time.Hour, false},
		{"P0.001D", 86400 * time.Millisecond, false},
		{"P1M", 1000 * oneMonthApprox, false},
		{"P0.001M", oneMonthApprox, false},
		{"P1Y", oneYearApprox, false},
		{"P0.001Y", oneYearApprox / 1000, false},
		{"P106751.991D", 9223372022400 * time.Millisecond, false}, // duration just less than 2^63-1
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

//TODO
func xTestPeriodApproxDays(t *testing.T) {
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

//TODO
func xTestPeriodApproxMonths(t *testing.T) {
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
		years, months, days, hours, minutes, seconds int
		period                                       Period
	}{
		{0, 0, 0, 0, 0, 0, Period{}},

		{0, 0, 0, 0, 0, 1, Period{centiSeconds: 100}},
		{0, 0, 0, 0, 1, 0, Period{centiSeconds: 6000}},
		{0, 0, 0, 1, 0, 0, Period{centiSeconds: 360000}},
		{0, 0, 1, 0, 0, 0, Period{centiMonths: 100}},
		{0, 1, 0, 0, 0, 0, Period{centiMonths: 100}},
		{1, 0, 0, 0, 0, 0, Period{centiMonths: 1200}},
		//{100, 222, 700, 0, 0, 0, Period{1422000, 700000, 0}},
		{0, 0, 0, 0, 0, -1, Period{centiSeconds: -100}},
		{0, 0, 0, 0, -1, 0, Period{centiSeconds: -6000}},
		{0, 0, 0, -1, 0, 0, Period{centiSeconds: -360000}},
		{0, 0, -1, 0, 0, 0, Period{centiMonths: -100}},
		{0, -1, 0, 0, 0, 0, Period{centiMonths: -100}},
		{-1, 0, 0, 0, 0, 0, Period{centiMonths: -1200}},
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
		hours, minutes, seconds int
		period                  Period
	}{
		{0, 0, 0, Period{}},
		{0, 0, 1, Period{centiSeconds: 100}},
		{0, 1, 0, Period{centiSeconds: 6000}},
		{1, 0, 0, Period{centiSeconds: 360000}},
		{0, 0, -1, Period{centiSeconds: -100}},
		{0, -1, 0, Period{centiSeconds: -6000}},
		{-1, 0, 0, Period{centiSeconds: -360000}},
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
		years, months, days int
		period              Period
	}{
		{0, 0, 0, Period{}},
		{0, 0, 1, Period{centiDays: 100}},
		{0, 1, 0, Period{centiMonths: 100}},
		{1, 0, 0, Period{centiMonths: 1200}},
		{118, 6, 700, Period{centiMonths: 142200, centiDays: 70000}},
		{0, 0, -1, Period{centiDays: -100}},
		{0, -1, 0, Period{centiMonths: -100}},
		{-1, 0, 0, Period{centiMonths: -1200}},
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
	testNewOf(t, 10*time.Millisecond, Period{centiDays: 1}, true)
	testNewOf(t, time.Second, Period{centiDays: 100}, true)
	testNewOf(t, time.Minute, Period{centiDays: 6000}, true)
	testNewOf(t, time.Hour, Period{centiDays: 360000}, true)
	testNewOf(t, time.Hour+time.Minute+time.Second, Period{centiDays: 366100}, true)
	testNewOf(t, 24*time.Hour+time.Minute+time.Second, Period{centiDays: 8646100}, true)
	testNewOf(t, 5965*time.Hour+13*time.Minute+56*time.Second+time.Duration(47*centiSecond), Period{centiDays: 2147483647}, true)

	// YMD tests: must be over 5965 hours (approx 8.2 months), otherwise HMS will take care of it
	// first rollover: >3276 hours
	//testNewOf(t, 5966*time.Hour, Period{0, 2485, 5040000}, false)
	//testNewOf(t, 3288*time.Hour, Period{0, 0, 0}, false)
	//testNewOf(t, 3289*time.Hour, Period{0, 0, 0}, false)
	//testNewOf(t, 24*3276*time.Hour, Period{0, 0, 0}, false)

	// second rollover: >3276 days
	//testNewOf(t, 24*3277*time.Hour, Period{80, 110, 200, 0, 0, 0}, false)
	//testNewOf(t, 3277*oneDay, Period{80, 110, 200, 0, 0, 0}, false)
	//testNewOf(t, 3277*oneDay+time.Hour+time.Minute+time.Second, Period{80, 110, 200, 10, 0, 0}, false)
	//testNewOf(t, 36525*oneDay, Period{1000, 0, 0, 0, 0, 0}, false)
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

	info := fmt.Sprintf("%d (%s) -> expected %s, actual %s %d", int64(source), source, expected, n, rev)
	g.Expect(n).To(Equal(expected), info)
	g.Expect(p).To(Equal(precise), info)
	g.Expect(rev).To(Equal(source), info)
}

//TODO
func TestBetween(t *testing.T) {
	g := NewGomegaWithT(t)
	now := time.Now()

	cases := []struct {
		a, b     time.Time
		expected Period
	}{
		{now, now, Period{}},

		// simple positive date calculations
		//{utc(2015, 1, 1, 0, 0, 0, 0), utc(2015, 1, 1, 0, 0, 0, 100), Period{0, 0, 0, 0, 1}},
		//{utc(2015, 1, 1, 0, 0, 0, 0), utc(2015, 2, 2, 1, 1, 1, 1), Period{0, 320, 10, 10, 10}},
		//{utc(2015, 2, 1, 0, 0, 0, 0), utc(2015, 3, 2, 1, 1, 1, 1), Period{0, 290, 10, 10, 10}},
		//{utc(2015, 3, 1, 0, 0, 0, 0), utc(2015, 4, 2, 1, 1, 1, 1), Period{0, 320, 10, 10, 10}},
		//{utc(2015, 4, 1, 0, 0, 0, 0), utc(2015, 5, 2, 1, 1, 1, 1), Period{0, 310, 10, 10, 10}},
		//{utc(2015, 5, 1, 0, 0, 0, 0), utc(2015, 6, 2, 1, 1, 1, 1), Period{0, 320, 10, 10, 10}},
		//{utc(2015, 6, 1, 0, 0, 0, 0), utc(2015, 7, 2, 1, 1, 1, 1), Period{0, 310, 10, 10, 10}},
		//{utc(2015, 1, 1, 0, 0, 0, 0), utc(2015, 7, 2, 1, 1, 1, 1), Period{0, 1820, 10, 10, 10}},

		// less than one month
		//{utc(2016, 1, 2, 0, 0, 0, 0), utc(2016, 2, 1, 0, 0, 0, 0), Period{0, 300, 0, 0, 0}},
		//{utc(2015, 2, 2, 0, 0, 0, 0), utc(2015, 3, 1, 0, 0, 0, 0), Period{0, 270, 0, 0, 0}}, // non-leap
		//{utc(2016, 2, 2, 0, 0, 0, 0), utc(2016, 3, 1, 0, 0, 0, 0), Period{0, 280, 0, 0, 0}}, // leap year
		//{utc(2016, 3, 2, 0, 0, 0, 0), utc(2016, 4, 1, 0, 0, 0, 0), Period{0, 300, 0, 0, 0}},
		//{utc(2016, 4, 2, 0, 0, 0, 0), utc(2016, 5, 1, 0, 0, 0, 0), Period{0, 290, 0, 0, 0}},
		//{utc(2016, 5, 2, 0, 0, 0, 0), utc(2016, 6, 1, 0, 0, 0, 0), Period{0, 300, 0, 0, 0}},
		//{utc(2016, 6, 2, 0, 0, 0, 0), utc(2016, 7, 1, 0, 0, 0, 0), Period{0, 290, 0, 0, 0}},

		// BST drops an hour at the daylight-saving transition
		//{utc(2015, 1, 1, 0, 0, 0, 0), bst(2015, 7, 2, 1, 1, 1, 1), Period{0, 1820, 0, 10, 10}},

		// negative date calculation
		//{utc(2015, 1, 1, 0, 0, 0, 100), utc(2015, 1, 1, 0, 0, 0, 0), Period{0, 0, 0, 0, -1}},
		//{utc(2015, 6, 2, 0, 0, 0, 0), utc(2015, 5, 1, 0, 0, 0, 0), Period{0, -320, 0, 0, 0}},
		//{utc(2015, 6, 2, 1, 1, 1, 1), utc(2015, 5, 1, 0, 0, 0, 0), Period{0, -320, -10, -10, -10}},

		// daytime only
		//{utc(2015, 1, 1, 2, 3, 4, 0), utc(2015, 1, 1, 2, 3, 4, 500), Period{0, 0, 0, 0, 5}},
		//{utc(2015, 1, 1, 2, 3, 4, 0), utc(2015, 1, 1, 4, 4, 7, 500), Period{0, 0, 20, 10, 35}},
		//{utc(2015, 1, 1, 2, 3, 4, 500), utc(2015, 1, 1, 4, 4, 7, 0), Period{0, 0, 20, 10, 25}},

		// different dates and times
		//{utc(2015, 2, 1, 1, 0, 0, 0), utc(2015, 5, 30, 5, 6, 7, 0), Period{0, 1180, 40, 60, 70}},
		//{utc(2015, 2, 1, 1, 0, 0, 0), bst(2015, 5, 30, 5, 6, 7, 0), Period{0, 1180, 30, 60, 70}},

		// earlier month in later year
		//{utc(2015, 12, 22, 0, 0, 0, 0), utc(2016, 1, 10, 5, 6, 7, 0), Period{0, 190, 50, 60, 70}},
		//{utc(2015, 2, 11, 5, 6, 7, 500), utc(2016, 1, 10, 0, 0, 0, 0), Period{0, 3320, 180, 530, 525}},

		// larger ranges
		//{utc(2009, 1, 1, 0, 0, 1, 0), utc(2016, 12, 31, 0, 0, 2, 0), Period{0, 29210, 0, 0, 10}},
		//{utc(2008, 1, 1, 0, 0, 1, 0), utc(2016, 12, 31, 0, 0, 2, 0), Period{80, 110, 300, 0, 0, 10}},
	}
	for i, c := range cases {
		n := Between(c.a, c.b)
		g.Expect(n).To(Equal(c.expected), info(i, c.expected))
	}
}

//TODO
func TestNormalise(t *testing.T) {
	// zero cases
	testNormalise(t, New(0, 0, 0, 0, 0, 0), New(0, 0, 0, 0, 0, 0), New(0, 0, 0, 0, 0, 0))

	// carry centiSeconds to minutes
	//testNormalise(t, Period{0, 0, 0, 0, 699}, Period{0, 0, 0, 10, 99}, Period{0, 0, 0, 10, 99})
	//
	//// carry minutes to centiSeconds
	//testNormalise(t, Period{0, 0, 0, 5, 0}, Period{0, 0, 0, 0, 300}, Period{0, 0, 0, 0, 300})
	//testNormalise(t, Period{0, 0, 0, 1, 0}, Period{0, 0, 0, 0, 60}, Period{0, 0, 0, 0, 60})
	//testNormalise(t, Period{0, 0, 0, 55, 0}, Period{0, 0, 0, 50, 300}, Period{0, 0, 0, 50, 300})
	//
	//// carry minutes to hours
	//testNormalise(t, Period{0, 0, 0, 699, 0}, Period{0, 0, 10, 90, 540}, Period{0, 0, 10, 90, 540})
	//
	//// carry hours to minutes
	//testNormalise(t, Period{0, 0, 5, 0, 0}, Period{0, 0, 0, 300, 0}, Period{0, 0, 0, 300, 0})
	//
	//// carry hours to days
	//testNormalise(t, Period{0, 0, 249, 0, 0}, Period{0, 0, 240, 540, 0}, Period{0, 0, 240, 540, 0})
	//testNormalise(t, Period{0, 0, 249, 0, 0}, Period{0, 0, 240, 540, 0}, Period{0, 0, 240, 540, 0})
	//testNormalise(t, Period{0, 0, 369, 0, 0}, Period{0, 0, 360, 540, 0}, Period{0, 10, 120, 540, 0})
	//testNormalise(t, Period{0, 0, 249, 0, 10}, Period{0, 0, 240, 540, 10}, Period{0, 0, 240, 540, 10})
	//
	//// carry days to hours
	//testNormalise(t, Period{0, 5, 30, 0, 0}, Period{0, 0, 150, 00, 0}, Period{0, 0, 150, 0, 0})
	//
	//// carry months to years
	//testNormalise(t, Period{125, 0, 0, 0, 0}, Period{125, 0, 0, 0, 0}, Period{125, 0, 0, 0, 0})
	//testNormalise(t, Period{131, 0, 0, 0, 0}, Period{10, 11, 0, 0, 0, 0}, Period{10, 11, 0, 0, 0, 0})
	//
	//// carry days to months
	//testNormalise(t, Period{0, 323, 0, 0, 0}, Period{0, 323, 0, 0, 0}, Period{0, 323, 0, 0, 0})
	//
	//// carry months to days
	//testNormalise(t, Period{5, 203, 0, 0, 0}, Period{0, 355, 0, 0, 0}, Period{10, 50, 0, 0, 0})
	//
	//// full ripple up
	//testNormalise(t, Period{121, 305, 239, 591, 601}, Period{10, 0, 330, 360, 540, 61}, Period{10, 10, 40, 0, 540, 61})
	//
	//// carry years to months
	//testNormalise(t, Period{5, 0, 0, 0, 0, 0}, Period{60, 0, 0, 0, 0}, Period{60, 0, 0, 0, 0})
	//testNormalise(t, Period{5, 25, 0, 0, 0, 0}, Period{85, 0, 0, 0, 0}, Period{85, 0, 0, 0, 0})
	//testNormalise(t, Period{5, 20, 10, 0, 0, 0}, Period{80, 10, 0, 0, 0}, Period{80, 10, 0, 0, 0})
}

func testNormalise(t *testing.T, source, precise, approx Period) {
	t.Helper()

	testNormaliseBothSigns(t, source, precise, true)
	testNormaliseBothSigns(t, source, approx, false)
}

func testNormaliseBothSigns(t *testing.T, source, expected Period, precise bool) {
	g := NewGomegaWithT(t)
	t.Helper()

	n1 := source.Normalise(precise)
	if n1 != expected {
		t.Errorf("%v.Normalise(%v) %s\n   gives %-22s %#v %s,\n    want %-22s %#v %s",
			source, precise, source.DurationApprox(),
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
		{"P1Y3M9D", "1 year, 3 months, 1 week, 2 days"},
		{"PT1H2M3.45S", "1 hour, 2 minutes, 3.45 seconds"},
		{"P1Y1M8DT1H1M1S", "1 year, 1 month, 1 week, 1 day, 1 hour, 1 minute, 1 second"},
		{"P3Y6M39DT2H7M9.12S", "3 years, 6 months, 5 weeks, 4 days, 2 hours, 7 minutes, 9.12 seconds"},
		{"-P3Y6M39DT2H7M9.12S", "3 years, 6 months, 5 weeks, 4 days, 2 hours, 7 minutes, 9.12 seconds"},
		{"P1.1Y", "1 year, 1.2 months"},
		{"P2.5Y", "2 years, 6 months"},
		{"P2.15Y", "2 years, 1.8 months"},
	}
	for i, c := range cases {
		s := MustParse(c.period).Format()
		g.Expect(s).To(Equal(c.expect), info(i, c.expect))
	}
}

func TestPeriodScale(t *testing.T) {
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
		//{"PT1S", 86400000, "PT86400000S"},
		{"P365.5D", 10, "P3655D"},
		//{"P365.5D", 0.1, "P36DT12H"},
	}
	for i, c := range cases {
		s := MustParse(c.one).Scale(c.m)
		g.Expect(s).To(Equal(MustParse(c.expect)), info(i, c.expect))
	}
}

func TestPeriodAdd(t *testing.T) {
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
		{"P1.1Y", "1 year, 1.2 months"},
		{"P2.5Y", "2 years, 6 months"},
		{"P2.15Y", "2 years, 1.8 months"},
		{"P2.125Y", "2 years, 1.44 months"}, // note truncation of last digit
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

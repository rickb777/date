// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package period

import (
	"fmt"
	"strings"
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

var oneDay = 24 * time.Hour
var oneMonthApprox = 2629746 * time.Second // 30.436875 days
var oneYearApprox = 31556952 * time.Second // 365.2425 days

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
		{"P1Y", Period{years: 10}},
		{"P1M", Period{months: 10}},
		{"P1W", Period{weeks: 10}},
		{"P1D", Period{days: 10}},
		{"PT1H", Period{hours: 10}},
		{"PT1M", Period{minutes: 10}},
		{"PT1S", Period{seconds: 10}},
		// smallest
		{"P0.1Y", Period{years: 1}},
		{"P0.1M", Period{months: 1}},
		{"P0.1W", Period{weeks: 1}},
		{"P0.7D", Period{days: 7}},
		{"P0.1D", Period{days: 1}},
		{"PT0.1H", Period{hours: 1}},
		{"PT0.1M", Period{minutes: 1}},
		{"PT0.1S", Period{seconds: 1}},

		{"P3Y", Period{years: 30}},
		{"P6M", Period{months: 60}},
		{"P5W", Period{weeks: 50}},
		{"P4D", Period{days: 40}},
		{"PT12H", Period{hours: 120}},
		{"PT30M", Period{minutes: 300}},
		{"PT5S", Period{seconds: 50}},
		{"P3Y6M39DT1H2M4.9S", Period{years: 30, months: 60, days: 390, hours: 10, minutes: 20, seconds: 49}},

		{"P2.5Y", Period{years: 25}},
		{"P2.5M", Period{months: 25}},
		{"P2.5D", Period{days: 25}},
		{"PT2.5H", Period{hours: 25}},
		{"PT2.5M", Period{minutes: 25}},
		{"PT2.5S", Period{seconds: 25}},
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
		value                          string
		y, m, w, dw, d, dx, hh, mm, ss int
	}{
		// note: the negative cases are also covered (see below)

		{value: "P0D"},
		{value: "P1Y", y: 1},
		{value: "P1M", m: 1},
		{value: "P6M", m: 6},
		{value: "P12M", m: 12},
		{value: "P1W", w: 1, dw: 7},
		{value: "P6W", w: 6, dw: 42},
		{value: "P52W", w: 52, dw: 364},
		{value: "P39D", dw: 39, d: 39, dx: 4},
		{value: "P4D", dw: 4, d: 4, dx: 4},

		{value: "PT1H", hh: 1},
		{value: "PT24H", hh: 24},
		{value: "PT1M", mm: 1},
		{value: "PT60M", mm: 60},
		{value: "PT1M", mm: 1},
		{value: "PT30M", mm: 30},
		{value: "PT60S", ss: 60},
		{value: "PT1S", ss: 1},

		{value: "P1Y2M3W4DT5H6M7S", y: 1, m: 2, w: 3, dw: 25, d: 4, dx: 4, hh: 5, mm: 6, ss: 7},
	}
	for i, c := range cases {
		pp := MustParse(c.value, Verbatim)
		g.Expect(pp.Years()).To(Equal(c.y), info(i, pp))
		g.Expect(pp.Months()).To(Equal(c.m), info(i, pp))
		g.Expect(pp.Weeks()).To(Equal(c.w), info(i, pp))
		g.Expect(pp.DaysAndWeeks()).To(Equal(c.dw), info(i, pp))
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
		{value: "P1.1Y", y: 1.1},
		{value: "P1M", m: 1},
		{value: "P1.5M", m: 1.5},
		{value: "P1.1M", m: 1.1},
		{value: "P6M", m: 6},
		{value: "P12M", m: 12},
		{value: "P1W", w: 1},
		{value: "P7D", d: 7},
		{value: "P1.1W", w: 1.1},
		{value: "P7.7D", d: 7.7},
		{value: "P7.1D", d: 7.1},
		{value: "P1D", d: 1},
		{value: "P1.1D", d: 1.1},
		{value: "P1.1D", d: 1.1},
		{value: "P39D", d: 39, dx: 4},
		{value: "P4D", d: 4, dx: 4},

		// HMS cases
		{value: "PT1.1H", hh: 1.1},
		{value: "PT1H6M", hh: 1, mm: 6},
		{value: "PT12H", hh: 12},
		{value: "PT1.1M", mm: 1.1},
		{value: "PT1M6S", mm: 1, ss: 6},
		{value: "PT30M", mm: 30},
		{value: "PT1.1S", ss: 1.1},
		{value: "PT5S", ss: 5},
	}
	for i, c := range cases {
		pp := MustParse(c.value, Verbatim)
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
		{"P1D", 24 * time.Hour, false},
		// days, months and years conversions are never precise
		{"P0.1D", 144 * time.Minute, false},
		{"P3276D", 3276 * 24 * time.Hour, false},
		{"P1W", 168 * time.Hour, false},
		{"P0.1W", 16*time.Hour + 48*time.Minute, false},
		{"P3276W", 3276 * 7 * 24 * time.Hour, false},
		{"P1M", oneMonthApprox, false},
		{"P0.1M", oneMonthApprox / 10, false},
		{"P3276M", 3276 * oneMonthApprox, false},
		{"P1Y", oneYearApprox, false},
		{"P3276Y", 3276 * oneYearApprox, false}, // near the upper limit of range
		// long second spans
		{"PT86400S", 86400 * time.Second, true},
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
	pp := MustParse(value, Constrained)
	d1, prec := pp.Duration()
	g.Expect(d1).To(Equal(duration), hint)
	g.Expect(prec).To(Equal(precise), hint)
	d2 := pp.DurationApprox()
	g.Expect(d2).To(Equal(duration), hint)
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
		{"P1M", true, false, 1},
		{"P0.1M", true, false, 1},
		{"-P1M", false, true, -1},
		{"-P0.1M", false, true, -1},
		{"P1D", true, false, 1},
		{"P10.D", true, false, 1},
		{"-P1D", false, true, -1},
		{"-P0.1D", false, true, -1},
		{"P1W", true, false, 1},
		{"P10.W", true, false, 1},
		{"-P1W", false, true, -1},
		{"-P0.1W", false, true, -1},
		{"P1Y", true, false, 1},
		{"P0.1Y", true, false, 1},
		{"-P1Y", false, true, -1},
		{"-P0.1Y", false, true, -1},
	}
	for i, c := range cases {
		p := MustParse(c.value, Verbatim)
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
		{"P1W", 7},
		{"P1M", 30},
		{"P1Y", 365},
	}
	for i, c := range cases {
		p := MustParse(c.value, Verbatim)
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
		{"P4.4W", 1},
		{"P1M", 1},
		{"P12M", 12},
		{"P2M31D", 3},
		{"P1Y", 12},
		{"P2Y3M", 27},
		{"PT24H", 0},
		{"PT744H", 1},
	}
	for i, c := range cases {
		p := MustParse(c.value, Verbatim)
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
		years, months, days, hours, minutes, seconds int
		ep                                           Period
	}{
		// note: the negative cases are also covered (see below)

		{ep: Period{}}, // zero case

		{seconds: 1, ep: Period{seconds: 10}},
		{minutes: 1, ep: Period{minutes: 10}},
		{hours: 1, ep: Period{hours: 10}},
		{days: 1, ep: Period{days: 10}},
		{months: 1, ep: Period{months: 10}},
		{years: 1, ep: Period{years: 10}},
		{years: 100, months: 222, days: 700, ep: Period{years: 1000, months: 2220, days: 7000, denormal: true}},
	}
	for i, c := range cases {
		pp := New(c.years, c.months, c.days, c.hours, c.minutes, c.seconds)
		expectValid(t, pp, info(i, c.ep))
		g.Expect(pp).To(Equal(c.ep), info(i, c.ep))
		g.Expect(pp.Years()).To(Equal(c.years), info(i, c.ep))
		g.Expect(pp.Months()).To(Equal(c.months), info(i, c.ep))
		g.Expect(pp.Days()).To(Equal(c.days), info(i, c.ep))

		pn := New(-c.years, -c.months, -c.days, -c.hours, -c.minutes, -c.seconds)
		en := c.ep.Negate()
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

		{period: Period{seconds: 10}, seconds: 1},
		{period: Period{minutes: 10}, minutes: 1},
		{period: Period{hours: 10}, hours: 1},
		{period: Period{hours: 30, minutes: 40, seconds: 50}, hours: 3, minutes: 4, seconds: 5},
		{period: Period{hours: 32760, minutes: 32760, seconds: 32760, denormal: true}, hours: 3276, minutes: 3276, seconds: 3276},
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

		{period: Period{days: 10}, days: 1},
		{period: Period{months: 10}, months: 1},
		{period: Period{years: 10}, years: 1},
		{period: Period{years: 1000, months: 2220, days: 7000, denormal: true}, years: 100, months: 222, days: 700},
		{period: Period{years: 32760, months: 32760, days: 32760, denormal: true}, years: 3276, months: 3276, days: 3276},
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

	// HMS tests
	testNewOf(t, 1, 100*time.Millisecond, Period{seconds: 1}, true)
	testNewOf(t, 2, time.Second, Period{seconds: 10}, true)
	testNewOf(t, 3, time.Minute, Period{minutes: 10}, true)
	testNewOf(t, 4, time.Hour, Period{hours: 10}, true)
	testNewOf(t, 5, time.Hour+time.Minute+time.Second, Period{hours: 10, minutes: 10, seconds: 10}, true)
	testNewOf(t, 6, 24*time.Hour+time.Minute+time.Second, Period{hours: 240, minutes: 10, seconds: 10}, true)
	testNewOf(t, 7, 3276*time.Hour+59*time.Minute+59*time.Second, Period{hours: 32760, minutes: 590, seconds: 590}, true)
	testNewOf(t, 8, 30*time.Minute+67*time.Second+600*time.Millisecond, Period{minutes: 310, seconds: 76}, true)

	// YMD tests: must be over 3276 hours (approx 4.5 months), otherwise HMS will take care of it
	// first rollover: >3276 hours
	testNewOf(t, 10, 3277*time.Hour, Period{days: 1360, hours: 130}, false)
	testNewOf(t, 11, 3288*time.Hour, Period{days: 1370}, false)
	testNewOf(t, 12, 3289*time.Hour, Period{days: 1370, hours: 10}, false)
	testNewOf(t, 13, 24*3276*time.Hour, Period{days: 32760}, false)

	// second rollover: >3276 days
	testNewOf(t, 14, 24*3277*time.Hour, Period{years: 80, months: 110, days: 200}, false)
	testNewOf(t, 15, 3277*oneDay, Period{years: 80, months: 110, days: 200}, false)
	testNewOf(t, 16, 3277*oneDay+time.Hour+time.Minute+time.Second, Period{years: 80, months: 110, days: 200, hours: 10}, false)
	testNewOf(t, 17, 36525*oneDay, Period{years: 1000}, false)
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
		{utc(2015, 1, 1, 0, 0, 0, 0), utc(2015, 1, 1, 0, 0, 0, 100), Period{seconds: 1}},
		{utc(2015, 1, 1, 0, 0, 0, 0), utc(2015, 2, 2, 1, 1, 1, 1), Period{days: 320, hours: 10, minutes: 10, seconds: 10, denormal: true}},
		{utc(2015, 2, 1, 0, 0, 0, 0), utc(2015, 3, 2, 1, 1, 1, 1), Period{days: 290, hours: 10, minutes: 10, seconds: 10}},
		{utc(2015, 3, 1, 0, 0, 0, 0), utc(2015, 4, 2, 1, 1, 1, 1), Period{days: 320, hours: 10, minutes: 10, seconds: 10, denormal: true}},
		{utc(2015, 4, 1, 0, 0, 0, 0), utc(2015, 5, 2, 1, 1, 1, 1), Period{days: 310, hours: 10, minutes: 10, seconds: 10, denormal: true}},
		{utc(2015, 5, 1, 0, 0, 0, 0), utc(2015, 6, 2, 1, 1, 1, 1), Period{days: 320, hours: 10, minutes: 10, seconds: 10, denormal: true}},
		{utc(2015, 6, 1, 0, 0, 0, 0), utc(2015, 7, 2, 1, 1, 1, 1), Period{days: 310, hours: 10, minutes: 10, seconds: 10, denormal: true}},
		{utc(2015, 1, 1, 0, 0, 0, 0), utc(2015, 7, 2, 1, 1, 1, 1), Period{days: 1820, hours: 10, minutes: 10, seconds: 10, denormal: true}},

		// less than one month
		{utc(2016, 1, 2, 0, 0, 0, 0), utc(2016, 2, 1, 0, 0, 0, 0), Period{days: 300}},
		{utc(2015, 2, 2, 0, 0, 0, 0), utc(2015, 3, 1, 0, 0, 0, 0), Period{days: 270}}, // non-leap
		{utc(2016, 2, 2, 0, 0, 0, 0), utc(2016, 3, 1, 0, 0, 0, 0), Period{days: 280}}, // leap year
		{utc(2016, 3, 2, 0, 0, 0, 0), utc(2016, 4, 1, 0, 0, 0, 0), Period{days: 300}},
		{utc(2016, 4, 2, 0, 0, 0, 0), utc(2016, 5, 1, 0, 0, 0, 0), Period{days: 290}},
		{utc(2016, 5, 2, 0, 0, 0, 0), utc(2016, 6, 1, 0, 0, 0, 0), Period{days: 300}},
		{utc(2016, 6, 2, 0, 0, 0, 0), utc(2016, 7, 1, 0, 0, 0, 0), Period{days: 290}},

		// BST drops an hour at the daylight-saving transition
		{utc(2015, 1, 1, 0, 0, 0, 0), bst(2015, 7, 2, 1, 1, 1, 1), Period{days: 1820, minutes: 10, seconds: 10, denormal: true}},

		// daytime only
		{utc(2015, 1, 1, 2, 3, 4, 0), utc(2015, 1, 1, 2, 3, 4, 500), Period{seconds: 5}},
		{utc(2015, 1, 1, 2, 3, 4, 0), utc(2015, 1, 1, 4, 4, 7, 500), Period{hours: 20, minutes: 10, seconds: 35}},
		{utc(2015, 1, 1, 2, 3, 4, 500), utc(2015, 1, 1, 4, 4, 7, 0), Period{hours: 20, minutes: 10, seconds: 25}},

		// different dates and times
		{utc(2015, 2, 1, 1, 0, 0, 0), utc(2015, 5, 30, 5, 6, 7, 0), Period{days: 1180, hours: 40, minutes: 60, seconds: 70, denormal: true}},
		{utc(2015, 2, 1, 1, 0, 0, 0), bst(2015, 5, 30, 5, 6, 7, 0), Period{days: 1180, hours: 30, minutes: 60, seconds: 70, denormal: true}},

		// earlier month in later year
		{utc(2015, 12, 22, 0, 0, 0, 0), utc(2016, 1, 10, 5, 6, 7, 0), Period{days: 190, hours: 50, minutes: 60, seconds: 70}},
		{utc(2015, 2, 11, 5, 6, 7, 500), utc(2016, 1, 10, 0, 0, 0, 0), Period{days: 3320, hours: 180, minutes: 530, seconds: 525, denormal: true}},

		// larger ranges
		{utc(2009, 1, 1, 0, 0, 1, 0), utc(2016, 12, 31, 0, 0, 2, 0), Period{days: 29210, seconds: 10, denormal: true}},
		{utc(2008, 1, 1, 0, 0, 1, 0), utc(2016, 12, 31, 0, 0, 2, 0), Period{years: 80, months: 110, days: 300, seconds: 10}},
		{utc(1900, 1, 1, 0, 0, 1, 0), utc(2009, 12, 31, 0, 0, 2, 0), Period{years: 1090, months: 110, days: 300, seconds: 10}},
	}
	for i, c := range cases {
		pp := Between(c.a, c.b)
		g.Expect(pp).To(Equal(c.expected), info(i, c.expected))

		pn := Between(c.b, c.a)
		expectValid(t, pn, info(i, c.expected))
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

		{period64{years: 10, months: 10, days: 10, hours: 10, minutes: 10, seconds: 11}},

		{period64{days: 10, hours: 70}},
		{period64{days: 10, hours: 10, minutes: 10}},
		{period64{days: 10, hours: 10, seconds: 10}},
		{period64{months: 10, days: 10, hours: 10}},

		{period64{minutes: 10, seconds: 10}},
		{period64{hours: 10, minutes: 10}},
		{period64{years: 10, months: 7}},

		{period64{months: 11}},
		{period64{days: 11}},
		{period64{hours: 11}},
		{period64{minutes: 11}},
		{period64{seconds: 11}},

		// don't carry days to months
		// don't carry months to years
	}
	for i, c := range cases {
		g.Expect(c.source.checkOverflow()).NotTo(HaveOccurred())
		p := c.source.toPeriod()

		testNormalise(t, i, c.source, p, true)
		testNormalise(t, i, c.source, p, false)
		c.source.neg = true
		testNormalise(t, i, c.source, p.Negate(), true)
		testNormalise(t, i, c.source, p.Negate(), false)
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
		{source: period64{seconds: 600}, precise: Period{minutes: 10}},
		{source: period64{seconds: 700}, precise: Period{minutes: 10, seconds: 100}},
		{source: period64{seconds: 6990}, precise: Period{minutes: 110, seconds: 390}},

		// carry minutes to hours
		{source: period64{minutes: 700}, precise: Period{hours: 10, minutes: 100}},
		{source: period64{minutes: 6990}, precise: Period{hours: 110, minutes: 390}},

		// carry hours to days
		{source: period64{hours: 480}, precise: Period{hours: 480}, approx: Period{days: 20}},
		{source: period64{hours: 490}, precise: Period{hours: 490}, approx: Period{days: 20, hours: 10}},
		//{source: period64{hours: 32761}, precise: Period{hours: 32761}, approx: Period{months: 40, days: 140, minutes: 169}},
		//{source: period64{years: 10, months: 20, days: 30, hours: 32767}, precise: Period{years: 10, months: 20, days: 30, hours: 32767}, approx: Period{years: 10, months: 60, days: 170, hours: 175}},
		//{source: period64{hours: 32768}, precise: "P136DT12.8H", approx: "P4M 14D T17.6H"},
		//{source: period64{years: 10, months: 20, days: 30, hours: 32768}, precise: "P1Y 2M 139D T12.8H", approx: "P1Y 6M 17D T17.6H"},

		// carry days to months
		{source: period64{days: 310}, precise: Period{weeks: 40, days: 30}, approx: Period{weeks: 40, days: 30}},
		//{source: period64{days: 32760}, precise: "P3276D", approx: "P8Y 11M 19.2D"},
		//{source: period64{days: 32761}, precise: "P8Y 11M 19.3D"},

		// carry months to years
		{source: period64{months: 120}, precise: Period{years: 10}},
		{source: period64{months: 132}, precise: Period{years: 10, months: 12}},
		{source: period64{months: 250}, precise: Period{years: 20, months: 10}},

		// full ripple up
		{
			source:  period64{months: 130, days: 70, hours: 240, minutes: 600, seconds: 611},
			precise: Period{years: 10, months: 10, weeks: 10, days: 0, hours: 250, minutes: 10, seconds: 11},
			approx:  Period{years: 10, months: 10, weeks: 10, days: 10, hours: 10, minutes: 10, seconds: 11},
		},
	}
	for i, c := range cases {
		if c.approx.IsZero() {
			c.approx = c.precise
		}
		testNormalise(t, i, c.source, c.precise, true)
		testNormalise(t, i, c.source, c.approx, false)
		c.source.neg = true
		testNormalise(t, i, c.source, c.precise.Negate(), true)
		testNormalise(t, i, c.source, c.approx.Negate(), false)
	}
}

func testNormalise(t *testing.T, i int, source period64, expected Period, precise bool) {
	g := NewGomegaWithT(t)
	t.Helper()

	sstr := source.String()
	source.denormal = true
	sn := source.normalise64(precise)
	g.Expect(sn.checkOverflow()).NotTo(HaveOccurred())
	n := sn.toPeriod()
	info := fmt.Sprintf("%d: %s.Normalise(%v) expected %s to equal %s", i, sstr, precise, n, expected)
	expectValid(t, n, info)
	g.Expect(n).To(Equal(expected), info)

	if !precise {
		p1 := source.toPeriod()
		d1, pr1 := p1.Duration()
		d2, pr2 := expected.Duration()
		g.Expect(pr1).To(Equal(pr2), info)
		g.Expect(d1).To(Equal(d2), info)
	}
}

//-------------------------------------------------------------------------------------------------

func TestPeriodFormat(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		period          string
		expectWithWeeks string // with weeks
		expectDaysOnly  string // without weeks
	}{
		// note: the negative cases are also covered (see below)

		{"P0D", "0 days", ""},

		{"P1Y1M7D", "1 year, 1 month, 7 days", "1 year, 1 month, 7 days"},
		{"P1Y1M1W1D", "1 year, 1 month, 1 week, 1 day", "1 year, 1 month, 8 days"},
		{"PT1H1M1S", "1 hour, 1 minute, 1 second", ""},
		{"P1Y1M1W1DT1H1M1S", "1 year, 1 month, 1 week, 1 day, 1 hour, 1 minute, 1 second", ""},
		{"P3Y6M39DT2H7M9S", "3 years, 6 months, 39 days, 2 hours, 7 minutes, 9 seconds", ""},
		{"P365D", "365 days", ""},

		{"P1Y", "1 year", ""},
		{"P3Y", "3 years", ""},
		{"P1.1Y", "1.1 years", ""},
		{"P2.5Y", "2.5 years", ""},

		{"P1M", "1 month", ""},
		{"P6M", "6 months", ""},
		{"P1.1M", "1.1 months", ""},
		{"P2.5M", "2.5 months", ""},

		{"P1W", "1 week", "7 days"},
		{"P1.1W", "1.1 weeks", "7.7 days"},
		{"P7D", "7 days", "7 days"},
		{"P35D", "35 days", "35 days"},
		{"P1D", "1 day", "1 day"},
		{"P4D", "4 days", "4 days"},
		{"P1.1D", "1.1 days", ""},

		{"PT1H", "1 hour", ""},
		{"PT1.1H", "1.1 hours", ""},

		{"PT1M", "1 minute", ""},
		{"PT1.1M", "1.1 minutes", ""},

		{"PT1S", "1 second", ""},
		{"PT1.1S", "1.1 seconds", ""},
	}
	for i, c := range cases {
		p := MustParse(c.period, Verbatim)
		sp := p.Format()
		g.Expect(sp).To(Equal(c.expectWithWeeks), info(i, "a. %s -> %s", p, c.expectWithWeeks))

		en := p.Negate()
		sn := en.Format()
		g.Expect(sn).To(Equal(c.expectWithWeeks), info(i, "b. %s -> %s", en, c.expectWithWeeks))

		if c.expectDaysOnly != "" {
			s := MustParse(c.period, Verbatim).FormatWithoutWeeks()
			g.Expect(s).To(Equal(c.expectDaysOnly), info(i, "c. %s -> %s", p, c.expectDaysOnly))
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
		s := MustParse(c.one, Verbatim).OnlyYMD()
		g.Expect(s).To(Equal(MustParse(c.expect, Verbatim)), info(i, c.expect))
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
		s := MustParse(c.one, Verbatim).OnlyHMS()
		g.Expect(s).To(Equal(MustParse(c.expect, Verbatim)), info(i, c.expect))
	}
}

//-------------------------------------------------------------------------------------------------

func TestSimplify(t *testing.T) {
	cases := []struct {
		source          string
		precise, approx Period
	}{
		// note: the negative cases are also covered (see below)

		// simplify 1 year to months (a = 9)
		{source: "P1Y"},
		{source: "P1Y10M"},
		{source: "P1Y9M", precise: Period{months: 210, denormal: true}},
		{source: "P1Y8.9M", precise: Period{months: 209, denormal: true}},

		// simplify 1 week to days (approx only) (b = 6)
		{source: "P1W2D", precise: Period{days: 90, denormal: true}},
		{source: "P2W1D"},
		{source: "P1W5.9D", precise: Period{days: 129, denormal: true}},

		// simplify 1 day to hours (approx only) (b = 6)
		{source: "P1DT6H", precise: Period{days: 10, hours: 60}, approx: Period{hours: 300, denormal: true}},
		{source: "P1DT7H"},
		{source: "P1DT5.9H", precise: Period{days: 10, hours: 59}, approx: Period{hours: 299, denormal: true}},

		// simplify 1 hour to minutes (c = 10)
		{source: "PT1H"},
		{source: "PT1H21M"},
		{source: "PT1H10M", precise: Period{minutes: 700, denormal: true}},
		{source: "PT1H9.9M", precise: Period{minutes: 699, denormal: true}},

		// simplify 1 minute to seconds (d = 30)
		{source: "PT1M"},    // unchanged
		{source: "PT1M31S"}, // ditto
		{source: "PT1M30S", precise: Period{seconds: 900, denormal: true}},
		{source: "PT1M29.9S", precise: Period{seconds: 899, denormal: true}},

		// fractional years don't simplify
		{source: "P1.1Y"},

		// retained proper fractions
		{source: "P1Y0.1D"},
		{source: "P12M0.1D"},
		{source: "P1YT0.1H"},
		{source: "P1MT0.1H"},
		{source: "P1Y0.1M", precise: Period{months: 121, denormal: true}},
		{source: "P1DT0.1H", precise: Period{days: 10, hours: 1}, approx: Period{hours: 241, denormal: true}},
		{source: "P1YT0.1M"},
		{source: "P1MT0.1M"},
		{source: "P1DT0.1M"},

		// discard proper fractions - months
		{source: "P10Y0.1M", precise: Period{years: 100, months: 1}, approx: Period{years: 100}},
		// discard proper fractions - days
		{source: "P1Y0.1D", precise: Period{years: 10, days: 1}, approx: Period{years: 10}},
		{source: "P12M0.1D", precise: Period{months: 120, days: 1, denormal: true}, approx: Period{months: 120, denormal: true}},
		//TODO
		//// discard proper fractions - hours
		//{source: "P1YT0.1H", precise: "P1YT0.1H", approx: "P1Y"},
		//{source: "P1MT0.1H", precise: "P1MT0.1H", approx: "P1M"},
		//{source: "P30DT0.1H", precise: "P30DT0.1H", approx: "P30D"},
		//// discard proper fractions - minutes
		//{source: "P1YT0.1M", precise: "P1YT0.1M", approx: "P1Y"},
		//{source: "P1MT0.1M", precise: "P1MT0.1M", approx: "P1M"},
		//{source: "P1DT0.1M", precise: "P1DT0.1M", approx: "P1D"},
		//{source: "PT24H0.1M", precise: "PT24H0.1M", approx: "PT24H"},
		//// discard proper fractions - seconds
		//{source: "P1YT0.1S", precise: "P1YT0.1S", approx: "P1Y"},
		//{source: "P1MT0.1S", precise: "P1MT0.1S", approx: "P1M"},
		//{source: "P1DT0.1S", precise: "P1DT0.1S", approx: "P1D"},
		//{source: "PT1H0.1S", precise: "PT1H0.1S", approx: "PT1H"},
		//{source: "PT60M0.1S", precise: "PT60M0.1S", approx: "PT60M"},
	}
	for i, c := range cases {
		p := MustParse(nospace(c.source), Verbatim)
		if c.precise.IsZero() {
			// unchanged cases
			testSimplify(t, i, p, p, true)
			testSimplify(t, i, p.Negate(), p.Negate(), true)

		} else if c.approx.IsZero() {
			// changed but precise/approx has same result
			testSimplify(t, i, p, c.precise, true)
			testSimplify(t, i, p.Negate(), c.precise.Negate(), true)

		} else {
			// changed and precise/approx have different results
			testSimplify(t, i, p, c.precise, true)
			testSimplify(t, i, p.Negate(), c.precise.Negate(), true)
			testSimplify(t, i, p, c.approx, false)
			testSimplify(t, i, p.Negate(), c.approx.Negate(), false)
		}
	}

	g := NewGomegaWithT(t)
	g.Expect(Period{days: 10, hours: 70}.Simplify(false, 6, 7, 30)).To(Equal(Period{hours: 310, denormal: true}))
	g.Expect(Period{hours: 10, minutes: 300}.Simplify(true, 6, 30)).To(Equal(Period{minutes: 900, denormal: true}))
	g.Expect(Period{years: 10, months: 110}.Simplify(true, 11)).To(Equal(Period{months: 230, denormal: true}))
	g.Expect(Period{days: 10, hours: 60}.Simplify(false)).To(Equal(Period{hours: 300, denormal: true}))
}

func testSimplify(t *testing.T, i int, source Period, expected Period, precise bool) {
	g := NewGomegaWithT(t)
	t.Helper()

	sstr := source.String()
	n := source.Simplify(precise, 9, 6, 10, 30)
	info := fmt.Sprintf("%d: %s.Simplify(%v) expected %s to equal %s", i, sstr, precise, n, expected)
	expectValid(t, n, info)
	g.Expect(n).To(Equal(expected), info)
}

//-------------------------------------------------------------------------------------------------

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

func nospace(s string) string {
	b := new(strings.Builder)
	for _, r := range s {
		if r != ' ' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

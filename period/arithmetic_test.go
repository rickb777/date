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

func TestPeriodScale(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		one    string
		m      float32
		expect Period
	}{
		{"P0D", 2, Period{}},
		{"P1D", 2, Period{days: 20}},
		{"P1D", 0, Period{}},
		{"P1D", 365, Period{weeks: 520, days: 10}},
		{"P1W", 2, Period{weeks: 20}},
		{"P1M", 2, Period{months: 20}},
		{"P1M", 12, Period{years: 10}},
		//TODO {"P1Y3M", 1.0/15, "P1M"},
		{"P1Y", 2, Period{years: 20}},
		{"PT1H", 2, Period{hours: 20}},
		{"PT1M", 2, Period{minutes: 20}},
		{"PT1S", 2, Period{seconds: 20}},
		{"P1D", 0.5, Period{days: 5}},
		{"P1W", 0.5, Period{weeks: 5}},
		{"P1M", 0.5, Period{months: 5}},
		{"P1Y", 0.5, Period{years: 5}},
		{"PT1H", 0.5, Period{hours: 5}},
		{"PT1H", 0.1, Period{minutes: 60}},
		{"PT1H", 0.01, Period{seconds: 359}}, // rounding error
		{"PT1M", 0.5, Period{minutes: 5}},
		{"PT1S", 0.5, Period{seconds: 5}},
		{"PT1H", 1.0 / 3600, Period{seconds: 10}},
		{"P1Y2M3W1DT5H6M7S", 2, Period{years: 20, months: 40, weeks: 60, days: 20, hours: 100, minutes: 120, seconds: 140}},
		{"P2Y4M6W2DT10H12M14S", -0.5, Period{years: -10, months: -20, weeks: -30, days: -10, hours: -50, minutes: -60, seconds: -70}},
		{"-P2Y4M6W2DT10H12M14S", 0.5, Period{years: -10, months: -20, weeks: -30, days: -10, hours: -50, minutes: -60, seconds: -70}},
		{"-P2Y4M6W2DT10H12M14S", -0.5, Period{years: 10, months: 20, weeks: 30, days: 10, hours: 50, minutes: 60, seconds: 70}},
		{"PT1M", 60, Period{hours: 10}},
		{"PT1S", 60, Period{minutes: 10}},
		{"PT1S", 86400, Period{hours: 240}},
		{"PT1S", 86400000, Period{weeks: 1420, days: 60}},
		{"P365.5D", 10, Period{weeks: 5220, days: 10}},
		{"P365.5D", 0.1, Period{hours: 8770, minutes: 120}},
	}
	for i, c := range cases {
		p := MustParse(c.one, Verbatim)
		s := p.Scale(c.m)
		g.Expect(s).To(Equal(c.expect), info(i, "%s x %v = %s", c.one, c.m, c.expect))

		_, err := p.ScaleWithOverflowCheck(c.m)
		g.Expect(err).NotTo(HaveOccurred(), info(i, "%s x %v = %s", c.one, c.m, c.expect))
	}
}
func TestPeriodScale_overflow(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		one string
		m   float32
	}{
		{"PT1S", 1e12},
		{"PT1M", 1e10},
		{"PT1H", 1e8},
		{"P1D", 1e7},
		{"P1W", 1e6},
		{"P1M", 1e5},
		{"P1Y", 1e4},
	}
	for i, c := range cases {
		s, err := MustParse(c.one, Verbatim).ScaleWithOverflowCheck(c.m)
		g.Expect(err).To(HaveOccurred(), info(i, "%s x %v = %s", c.one, c.m, s))
	}
}

func TestPeriodAdd(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		one, two string
		expect   Period
	}{
		{"P0D", "P0D", Period{}},
		{"P1D", "P1D", Period{days: 20}},
		{"P1W", "P1W", Period{weeks: 20}},
		{"P1M", "P1M", Period{months: 20}},
		{"P1Y", "P1Y", Period{years: 20}},
		{"PT1H", "PT1H", Period{hours: 20}},
		{"PT1M", "PT1M", Period{minutes: 20}},
		{"PT1S", "PT1S", Period{seconds: 20}},
		{"P1Y2M3W4DT5H6M7S", "P7Y6M5W4DT3H2M1S", Period{years: 80, months: 80, weeks: 80, days: 80, hours: 80, minutes: 80, seconds: 80}},
		{"P7Y7M7W7DT7H7M7S", "-P7Y7M7W7DT7H7M7S", Period{}},
	}
	for i, c := range cases {
		s := MustParse(c.one, Verbatim).Add(MustParse(c.two, Verbatim))
		expectValid(t, s, info(i, c.expect))
		g.Expect(s).To(Equal(c.expect), info(i, c.expect))
	}
}

func TestPeriodAddToTime(t *testing.T) {
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
	}
	for i, c := range cases {
		p := MustParse(c.value)
		t1, prec := p.AddTo(t0)
		g.Expect(t1).To(Equal(c.result), info(i, c.value))
		g.Expect(prec).To(Equal(c.precise), info(i, c.value))
	}
}

func expectValid(t *testing.T, period Period, hint interface{}) Period {
	t.Helper()
	g := NewGomegaWithT(t)
	info := fmt.Sprintf("%v: invalid: %#v", hint, period)

	// check all the signs are consistent
	nPoz := pos(period.years) + pos(period.months) + pos(period.days) + pos(period.hours) + pos(period.minutes) + pos(period.seconds)
	nNeg := neg(period.years) + neg(period.months) + neg(period.days) + neg(period.hours) + neg(period.minutes) + neg(period.seconds)
	g.Expect(nPoz == 0 || nNeg == 0).To(BeTrue(), info+" inconsistent signs")

	// only one field must have a fraction
	yearsFraction := fraction(period.years)
	monthsFraction := fraction(period.months)
	daysFraction := fraction(period.days)
	hoursFraction := fraction(period.hours)
	minutesFraction := fraction(period.minutes)

	if yearsFraction != 0 {
		g.Expect(period.months).To(BeZero(), info+" year fraction exists")
		g.Expect(period.days).To(BeZero(), info+" year fraction exists")
		g.Expect(period.hours).To(BeZero(), info+" year fraction exists")
		g.Expect(period.minutes).To(BeZero(), info+" year fraction exists")
		g.Expect(period.seconds).To(BeZero(), info+" year fraction exists")
	}

	if monthsFraction != 0 {
		g.Expect(period.days).To(BeZero(), info+" month fraction exists")
		g.Expect(period.hours).To(BeZero(), info+" month fraction exists")
		g.Expect(period.minutes).To(BeZero(), info+" month fraction exists")
		g.Expect(period.seconds).To(BeZero(), info+" month fraction exists")
	}

	if daysFraction != 0 {
		g.Expect(period.hours).To(BeZero(), info+" day fraction exists")
		g.Expect(period.minutes).To(BeZero(), info+" day fraction exists")
		g.Expect(period.seconds).To(BeZero(), info+" day fraction exists")
	}

	if hoursFraction != 0 {
		g.Expect(period.minutes).To(BeZero(), info+" hour fraction exists")
		g.Expect(period.seconds).To(BeZero(), info+" hour fraction exists")
	}

	if minutesFraction != 0 {
		g.Expect(period.seconds).To(BeZero(), info+" minute fraction exists")
	}

	return period
}

func fraction(i int16) int {
	return int(i) % 10
}

func pos(i int16) int {
	if i > 0 {
		return 1
	}
	return 0
}

func neg(i int16) int {
	if i < 0 {
		return 1
	}
	return 0
}

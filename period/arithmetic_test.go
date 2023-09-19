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
		s := MustParse(c.one, false).Scale(c.m)
		g.Expect(s).To(Equal(MustParse(c.expect, false)), info(i, c.expect))
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
		s := MustParse(c.one, false).Add(MustParse(c.two, false))
		expectValid(t, s, info(i, c.expect))
		g.Expect(s).To(Equal(MustParse(c.expect, false)), info(i, c.expect))
	}
}

func TestPeriodAddToTime(t *testing.T) {
	g := NewGomegaWithT(t)

	const ms = 1000000
	const sec = 1000 * ms
	const min = 60 * sec
	const hr = 60 * min

	est, err := time.LoadLocation("America/New_York")
	g.Expect(err).NotTo(HaveOccurred())

	times := []time.Time{
		// A conveniently round number but with non-zero nanoseconds (14 July 2017 @ 2:40am UTC)
		time.Unix(1500000000, 1).UTC(),
		// This specific time fails for EST due behaviour of Time.AddDate
		time.Date(2020, 11, 1, 1, 0, 0, 0, est),
	}

	for _, t0 := range times {
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
			p, err := ParseWithNormalise(c.value, false)
			g.Expect(err).NotTo(HaveOccurred())

			t1, prec := p.AddTo(t0)
			g.Expect(t1).To(Equal(c.result), info(i, c.value, t0))
			g.Expect(prec).To(Equal(c.precise), info(i, c.value, t0))
		}
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

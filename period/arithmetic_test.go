// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package period

import (
	"fmt"
	"math"
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

func TestPeriodScale_errors(t *testing.T) {
	g := NewGomegaWithT(t)
	_, err := Period{}.ScaleWithOverflowCheck(float32(math.NaN()))
	g.Expect(err).To(HaveOccurred())
	_, err = Period{}.ScaleWithOverflowCheck(float32(math.Inf(1)))
	g.Expect(err).To(HaveOccurred())
}

//-------------------------------------------------------------------------------------------------

func TestPeriodScale_simpleCases(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		ymdDesignators string
		hmsDesignators string
		one            string
		m              float32
		expect         string
	}{
		// note: the negative cases are also covered (see below)

		{"YMDW", "HMS", "0", 2, "0"},
		{"YMDW", "HMS", "1", 0, "0"},
		{"YMDW", "HMS", "1", 1, "1"},
		{"YMDW", "HMS", "1", 2, "2"},
		{"YMD", "HMS", "1", 0.5, "0.5"},
		{"MD", "HMS", "1", 0.1, "0.1"},
		{"YMDW", "HMS", "10", 2, "20"},
		{"YMDW", "HMS", "400", 10, "4000"},
		{"YMDW", "HMS", "1", 500, "500"},
	}
	for i, c := range cases {
		for _, des := range c.ymdDesignators {
			pp := MustParse(fmt.Sprintf("P%s%c", c.one, des))
			ep := MustParse(fmt.Sprintf("P%s%c", c.expect, des))
			en := ep.Negate()

			g.Expect(pp.ScaleWithOverflowCheck(c.m)).To(Equal(ep), info(i, "%s x %g", pp, c.m))
			g.Expect(pp.ScaleWithOverflowCheck(-c.m)).To(Equal(en), info(i, "%s x %g", pp, c.m))

			pn := pp.Negate()
			g.Expect(pn.ScaleWithOverflowCheck(c.m)).To(Equal(en), info(i, "%s x %g", en, c.m))
			g.Expect(pn.ScaleWithOverflowCheck(-c.m)).To(Equal(ep), info(i, "%s x %g", en, c.m))
		}

		for _, des := range c.hmsDesignators {
			pp := MustParse(fmt.Sprintf("PT%s%c", c.one, des))
			ep := MustParse(fmt.Sprintf("PT%s%c", c.expect, des))
			g.Expect(pp.ScaleWithOverflowCheck(c.m)).To(Equal(ep), info(i, "%s x %g", pp, c.m))

			en := ep.Negate()
			pn := pp.Negate()
			g.Expect(pn.ScaleWithOverflowCheck(c.m)).To(Equal(en), info(i, "%s x %g", en, c.m))
			g.Expect(pn.ScaleWithOverflowCheck(-c.m)).To(Equal(ep), info(i, "%s x %g", en, c.m))
		}
	}
}

//-------------------------------------------------------------------------------------------------

func TestPeriodScale_complexCases(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		one    string
		m, d   int
		expect string
	}{
		// note: the negative cases are also covered (see below)

		{"PT1S", 1, 100, "PT0.01S"},
		{"PT1S", 1, 10, "PT0.1S"},
		{"PT1M", 1, 2, "PT30S"},
		{"PT1H", 1, 2, "PT30M"},
		{"PT1M", 1, 60, "PT1S"},
		{"PT1H", 1, 60, "PT1M"},
		{"PT1H", 1, 7, "PT8M34.29S"},
		{"PT1M", 1, 7, "PT8.57S"},

		{"PT1M", 60, 1, "PT1H"},
		{"PT1S", 60, 1, "PT1M"},
		{"PT1S", 86400, 1, "PT24H"},

		{"P1D", 1, 2, "P0.5D"},
		{"P1D", 1, 10, "P0.1D"},
		{"P1D", 1, 24, "PT1H"},
		{"P1D", 1, 1440, "PT1M"},
		{"P1D", 1, 86400, "PT1S"},

		{"P2M", 1, 2, "P1M"},
		{"P1M", 1, 2, "P0.5M"},
		{"P2Y", 1, 2, "P1Y"},
		{"P1Y", 1, 2, "P6M"},
		{"P1Y", 1, 365, "P1DT57.4S"},

		{"P1Y2M3DT4H5M6S", 2, 1, "P2Y4M6DT8H10M12S"},
		{"P2Y4M6DT8H10M12S", -1, 2, "-P1Y2M3DT4H5M6S"},
		{"-P2Y4M6DT8H10M12S", 1, 2, "-P1Y2M3DT4H5M6S"},
		{"-P2Y4M6DT8H10M12S", -1, 2, "P1Y2M3DT4H5M6S"},

		{"PT1S", 86400000, 1, "PT24000H"},
		{"PT1H", 24 * 32768, 1, "P89Y8M17DT22H8M"},
		{"P365.5D", 10, 1, "P3655D"},
		{"P365D", 1, 2, "P182.5D"},
		{"P3650D", 1, 10, "P365D"},

		// cases with acceptable small rounding errors
		{"P18262D", 1, 100, "P182.62D"},
	}
	for i, c := range cases {
		pp := MustParse(c.one)
		ep := MustParse(c.expect)
		g.Expect(pp.RationalScale(c.m, c.d)).To(Equal(ep), info(i, "%s x %d/%d", c.one, c.m, c.d))
		g.Expect(pp.Negate().RationalScale(c.m, c.d)).To(Equal(ep.Negate()), info(i, "-%s x %d/%d", c.one, c.m, c.d))
	}
}

//-------------------------------------------------------------------------------------------------

func TestPeriodAdd(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		one, two string
		expect   string
	}{
		// simple cases
		{"P0D", "P0D", "P0D"},
		{"P1D", "P1D", "P2D"},
		{"P1M", "P1M", "P2M"},
		{"P1Y", "P1Y", "P2Y"},
		{"PT1H", "PT1H", "PT2H"},
		{"PT1M", "PT1M", "PT2M"},
		{"PT1S", "PT1S", "PT2S"},
		{"P1Y2M3DT4H5M6.70S", "P6Y5M4DT3H2M1.07S", "P7Y7M7DT7H7M7.77S"},
		{"P7Y7M7DT7H7M7.77S", "-P7Y7M7DT7H7M7.77S", "P0D"},
		{"PT1.5M", "PT1M", "PT2.5M"},

		// non-trivial cases
		//{"PT1.5M", "PT32.5S", "PT2M2.5S"},

		// fraction handling - carry one
		{"P1.7Y", "P1.8Y", "P3.5Y"},
		{"P1.7M", "P1.8M", "P3.5M"},
		{"P1.7D", "P1.8D", "P3.5D"},
		{"PT1.7H", "PT1.8H", "PT3.5H"},
		{"PT1.7M", "PT1.8M", "PT3.5M"},
		{"PT1.7S", "PT1.8S", "PT3.5S"},

		// fraction handling - integer result
		{"P1.7Y", "P1.3Y", "P3Y"},
		{"P1.7M", "P1.3M", "P3M"},
		{"P1.7D", "P1.3D", "P3D"},
		{"PT1.7H", "PT1.3H", "PT3H"},
		{"PT1.7M", "PT1.3M", "PT3M"},
		{"PT1.7S", "PT1.3S", "PT3S"},
	}
	for i, c := range cases {
		p1 := MustParse(c.one, false)
		p2 := MustParse(c.two, false)
		pe := MustParse(c.expect, false)

		s := info(i, c.expect)
		g.Expect(expectValid(t, p1.Add(p2), s)).To(Equal(pe), s)
		g.Expect(expectValid(t, p1.Negate().Add(p2.Negate()), s)).To(Equal(pe.Negate()), s+" neg")
	}
}

//-------------------------------------------------------------------------------------------------

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
		{"PT32767S", t0.Add(32767 * sec), true},
		{"PT1M", t0.Add(60 * sec), true},
		{"PT0.1M", t0.Add(6 * sec), true},
		{"PT32767M", t0.Add(32767 * min), true},
		{"PT1H", t0.Add(hr), true},
		{"PT0.1H", t0.Add(6 * min), true},
		{"PT32767H", t0.Add(32767 * hr), true},
		{"P1D", t0.AddDate(0, 0, 1), true},
		{"P32767D", t0.AddDate(0, 0, 32767), true},
		{"P1M", t0.AddDate(0, 1, 0), true},
		{"P32767M", t0.AddDate(0, 32767, 0), true},
		{"P1Y", t0.AddDate(1, 0, 0), true},
		{"-P1Y", t0.AddDate(-1, 0, 0), true},
		{"P32767Y", t0.AddDate(32767, 0, 0), true},   // near the upper limit of range
		{"-P32767Y", t0.AddDate(-32767, 0, 0), true}, // near the lower limit of range
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
	nPoz := pos(period.years) + pos(period.months) + pos(period.days) + pos(period.hours) + pos(period.minutes) + pos(period.seconds) + pos(int16(period.fraction))
	nNeg := neg(period.years) + neg(period.months) + neg(period.days) + neg(period.hours) + neg(period.minutes) + neg(period.seconds) + neg(int16(period.fraction))
	if nPoz > 0 && nNeg > 0 {
		t.Errorf("%s: inconsistent signs in\n%#v", info, period)
	}

	// check no intermediate fraction is present
	switch period.fpart {
	case NoFraction:
		g.Expect(period.fraction).To(BeZero(), info)
	case Minute:
		g.Expect(period.seconds).To(BeZero(), info)
	case Hour:
		g.Expect(period.seconds).To(BeZero(), info)
		g.Expect(period.minutes).To(BeZero(), info)
	case Day:
		g.Expect(period.seconds).To(BeZero(), info)
		g.Expect(period.minutes).To(BeZero(), info)
		g.Expect(period.hours).To(BeZero(), info)
	case Month:
		g.Expect(period.seconds).To(BeZero(), info)
		g.Expect(period.minutes).To(BeZero(), info)
		g.Expect(period.hours).To(BeZero(), info)
		g.Expect(period.days).To(BeZero(), info)
	case Year:
		g.Expect(period.seconds).To(BeZero(), info)
		g.Expect(period.minutes).To(BeZero(), info)
		g.Expect(period.hours).To(BeZero(), info)
		g.Expect(period.days).To(BeZero(), info)
		g.Expect(period.months).To(BeZero(), info)
	}

	// the fraction must be in the range -99 to +99
	g.Expect(period.fraction).To(BeNumerically("<", 100), info)
	g.Expect(period.fraction).To(BeNumerically(">", -100), info)
	return period
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

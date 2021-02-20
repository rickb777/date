// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package period

import (
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

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
		{`PT1`, false, `: missing designator at the end`, "PT1"},
		{"XY", false, ": expected 'P' period mark at the start", "XY"},
		{"PxY", false, ": expected a number but found 'x'", "PxY"},
		{"PxW", false, ": expected a number but found 'x'", "PxW"},
		{"PxD", false, ": expected a number but found 'x'", "PxD"},
		{"PTxH", false, ": expected a number but found 'x'", "PTxH"},
		{"PTxM", false, ": expected a number but found 'x'", "PTxM"},
		{"PTxS", false, ": expected a number but found 'x'", "PTxS"},
		{"PT1A", false, ": expected a designator Y, M, W, D, H, or S not 'A'", "PT1A"},
		{"P1HT1M", false, ": 'H' designator cannot occur here", "P1HT1M"},
		{"PT1Y", false, ": 'Y' designator cannot occur here", "PT1Y"},
		{"P1S", false, ": 'S' designator cannot occur here", "P1S"},
		{"P1D2D", false, ": 'D' designator cannot occur more than once", "P1D2D"},
		{"PT1HT1S", false, ": 'T' designator cannot occur more than once", "PT1HT1S"},
		{"P0.1YT0.1S", false, ": 'Y' & 'S' only the last field can have a fraction", "P0.1YT0.1S"},
		{"P", false, ": expected 'Y', 'M', 'W', 'D', 'H', 'M', or 'S' designator", "P"},
		// integer overflow
		{"P32768Y", false, ": integer overflow occurred in years", "P32768Y"},
		{"P393216M", false, ": integer overflow occurred in months", "P393216M"},
		{"P1703936W", false, ": integer overflow occurred in weeks", "P1703936W"},
		{"P10657546D", false, ": integer overflow occurred in days", "P10657546D"},
		//{"PT32768H", false, ": integer overflow occurred in hours", "PT32768H"},
		//{"PT32768M", false, ": integer overflow occurred in minutes", "PT32768M"},
		//{"PT32768S", false, ": integer overflow occurred in seconds", "PT32768S"},
		//{"PT32768H32768M32768S", false, ": integer overflow occurred in hours,minutes,seconds", "PT32768H32768M32768S"},
		//{"PT103412160000S", false, ": integer overflow occurred in seconds", "PT103412160000S"},
	}
	for i, c := range cases {
		_, ep := Parse(c.value, Verbatim)
		g.Expect(ep).To(HaveOccurred(), info(i, c.value))
		g.Expect(ep.Error()).To(Equal(c.expvalue+c.expected), info(i, c.value))

		_, en := Parse("-"+c.value, Verbatim)
		g.Expect(en).To(HaveOccurred(), info(i, c.value))
		if c.expvalue != "" {
			g.Expect(en.Error()).To(Equal("-"+c.expvalue+c.expected), info(i, c.value))
		} else {
			g.Expect(en.Error()).To(Equal(c.expected), info(i, c.value))
		}
	}
}

//-------------------------------------------------------------------------------------------------

func TestParsePeriodNormalised(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		value    string
		reversed string
		period   Period
	}{
		// all rollovers
		{"PT1234.5S", "PT20M34.5S", Period{minutes: 200, seconds: 345}},
		{"PT1234.5M", "PT20H34.5M", Period{hours: 200, minutes: 345}},
		{"PT12345.6H", "P73W3DT9.6H", Period{weeks: 730, days: 30, hours: 96}},
		{"P3277W", "P62Y9M2.8W", Period{years: 620, months: 90, weeks: 28}},
		{"P3277D", "P468W1D", Period{years: 0, months: 0, weeks: 4680, days: 10}},
		{"P22939D", "P62Y9M2.8W", Period{years: 620, months: 90, weeks: 28}},
		{"P1234.5M", "P102Y10.5M", Period{years: 1020, months: 105}},
		// largest possible number of seconds normalised only in hours, mins, sec
		{"PT11592000S", "PT3220H", Period{hours: 32200}},
		{"-PT11592000S", "-PT3220H", Period{hours: -32200}},
		{"PT11595599S", "PT3220H59M59S", Period{hours: 32200, minutes: 590, seconds: 590}},
		// largest possible number of seconds normalised only in weeks, days, hours, mins, sec
		{"PT1981324800S", "P3276W", Period{weeks: 32760}},
		{"-PT1981324800S", "-P3276W", Period{weeks: -32760}},

		{"PT11793600S", "P19W3DT12H", Period{weeks: 190, days: 30, hours: 120}},
		// other examples are in TestNormalise
	}
	for i, c := range cases {
		p, err := Parse(c.value, Normalised)
		s := info(i, "%s %d", c.value, c.period.DurationApprox()/time.Second)
		g.Expect(err).NotTo(HaveOccurred(), s)
		expectValid(t, p, s)
		g.Expect(p).To(Equal(c.period), s)
		// reversal is expected not to be an identity
		g.Expect(p.String()).To(Equal(c.reversed), s+" reversed")
	}
}

//-------------------------------------------------------------------------------------------------

func TestParsePeriodVerbatim(t *testing.T) {
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
		{"P1W", "P1W", Period{weeks: 10}},
		{"P1D", "P1D", Period{days: 10}},
		{"PT1H", "PT1H", Period{hours: 10}},
		{"PT1M", "PT1M", Period{minutes: 10}},
		{"PT1S", "PT1S", Period{seconds: 10}},
		// smallest
		{"P0.1Y", "P0.1Y", Period{years: 1}},
		{"-P0.1Y", "-P0.1Y", Period{years: -1}},
		{"P0.1M", "P0.1M", Period{months: 1}},
		{"-P0.1M", "-P0.1M", Period{months: -1}},
		{"P0.1W", "P0.1W", Period{weeks: 1}},
		{"-P0.1W", "-P0.1W", Period{weeks: -1}},
		{"P0.1D", "P0.1D", Period{days: 1}},
		{"-P0.1D", "-P0.1D", Period{days: -1}},
		{"PT0.1H", "PT0.1H", Period{hours: 1}},
		{"-PT0.1H", "-PT0.1H", Period{hours: -1}},
		{"PT0.1M", "PT0.1M", Period{minutes: 1}},
		{"-PT0.1M", "-PT0.1M", Period{minutes: -1}},
		{"PT0.1S", "PT0.1S", Period{seconds: 1}},
		{"-PT0.1S", "-PT0.1S", Period{seconds: -1}},
		// largest
		{"PT3276.7S", "PT3276.7S", Period{seconds: 32767, denormal: true}},
		{"PT3276.7M", "PT3276.7M", Period{minutes: 32767, denormal: true}},
		{"PT3276.7H", "PT3276.7H", Period{hours: 32767, denormal: true}},
		{"P3276.7D", "P3276.7D", Period{days: 32767, denormal: true}},
		{"P3276.7W", "P3276.7W", Period{weeks: 32767, denormal: true}},
		{"P3276.7M", "P3276.7M", Period{months: 32767, denormal: true}},
		{"P3276.7Y", "P3276.7Y", Period{years: 32767, denormal: false}},

		{"P3Y", "P3Y", Period{years: 30}},
		{"P6M", "P6M", Period{months: 60}},
		{"P5W", "P5W", Period{weeks: 50}},
		{"P4D", "P4D", Period{days: 40}},
		{"PT12H", "PT12H", Period{hours: 120}},
		{"PT30M", "PT30M", Period{minutes: 300}},
		{"PT25S", "PT25S", Period{seconds: 250}},
		{"PT30M67.6S", "PT30M67.6S", Period{minutes: 300, seconds: 676, denormal: true}},
		{"P2.Y", "P2Y", Period{years: 20}},
		{"P2.5Y", "P2.5Y", Period{years: 25}},
		{"P2.15Y", "P2.1Y", Period{years: 21}},
		{"P1Y2.M", "P1Y2M", Period{years: 10, months: 20}},
		{"P1Y2.5M", "P1Y2.5M", Period{years: 10, months: 25}},
		{"P1Y2.15M", "P1Y2.1M", Period{years: 10, months: 21}},
		// others
		{"P3Y6M5W4DT12H40M5S", "P3Y6M5W4DT12H40M5S", Period{years: 30, months: 60, weeks: 50, days: 40, hours: 120, minutes: 400, seconds: 50}},
		{"+P3Y6M5W4DT12H40M5S", "P3Y6M5W4DT12H40M5S", Period{years: 30, months: 60, weeks: 50, days: 40, hours: 120, minutes: 400, seconds: 50}},
		{"-P3Y6M5W4DT12H40M5S", "-P3Y6M5W4DT12H40M5S", Period{years: -30, months: -60, weeks: -50, days: -40, hours: -120, minutes: -400, seconds: -50}},
		{"P1Y14M35DT48H125M800S", "P1Y14M35DT48H125M800S", Period{years: 10, months: 140, weeks: 0, days: 350, hours: 480, minutes: 1250, seconds: 8000, denormal: true}},
	}
	for i, c := range cases {
		p, err := Parse(c.value, Verbatim)
		s := info(i, c.value)
		g.Expect(err).NotTo(HaveOccurred(), s)
		expectValid(t, p, s)
		g.Expect(p).To(Equal(c.period), s)
		// reversal is usually expected to be an identity
		g.Expect(p.String()).To(Equal(c.reversed), s+" reversed")
	}
}

package period

import (
	"fmt"
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

func TestNewOfMS(t *testing.T) {
	// note: the negative cases are also covered (see below)

	// HMS tests
	testNewOfMS(t, 1, 100*time.Millisecond, PeriodMS{Period: Period{}, milliseconds: 100}, true)
	testNewOfMS(t, 2, time.Second, PeriodMS{Period: Period{seconds: 10}}, true)
	testNewOfMS(t, 3, time.Minute, PeriodMS{Period: Period{minutes: 10}}, true)
	testNewOfMS(t, 4, time.Hour, PeriodMS{Period: Period{hours: 10}}, true)
	testNewOfMS(t, 5, time.Hour+time.Minute+time.Second, PeriodMS{Period: Period{hours: 10, minutes: 10, seconds: 10}}, true)
	testNewOfMS(t, 6, 24*time.Hour+time.Minute+time.Second, PeriodMS{Period: Period{hours: 240, minutes: 10, seconds: 10}}, true)
	testNewOfMS(t, 7, 3276*time.Hour+59*time.Minute+59*time.Second, PeriodMS{Period: Period{hours: 32760, minutes: 590, seconds: 590}}, true)
	testNewOfMS(t, 8, 30*time.Minute+67*time.Second+600*time.Millisecond, PeriodMS{Period: Period{minutes: 310, seconds: 70}, milliseconds: 600}, true)

	// YMD tests: must be over 3276 hours (approx 4.5 months), otherwise HMS will take care of it
	// first rollover: >3276 hours
	testNewOfMS(t, 10, 3277*time.Hour, PeriodMS{Period: Period{days: 1360, hours: 130}}, false)
	testNewOfMS(t, 11, 3288*time.Hour, PeriodMS{Period: Period{days: 1370}}, false)
	testNewOfMS(t, 12, 3289*time.Hour, PeriodMS{Period: Period{days: 1370, hours: 10}}, false)
	testNewOfMS(t, 13, 24*3276*time.Hour, PeriodMS{Period: Period{days: 32760}}, false)

	// second rollover: >3276 days
	testNewOfMS(t, 14, 24*3277*time.Hour, PeriodMS{Period: Period{years: 80, months: 110, days: 200}}, false)
	testNewOfMS(t, 15, 3277*oneDay, PeriodMS{Period: Period{years: 80, months: 110, days: 200}}, false)
	testNewOfMS(t, 16, 3277*oneDay+time.Hour+time.Minute+time.Second, PeriodMS{Period: Period{years: 80, months: 110, days: 200, hours: 10}}, false)
	testNewOfMS(t, 17, 36525*oneDay, PeriodMS{Period: Period{years: 1000}}, false)
}

func testNewOfMS(t *testing.T, i int, source time.Duration, expected PeriodMS, precise bool) {
	t.Helper()
	testNewOf1MS(t, i, source, expected, precise)
	testNewOf1MS(t, i, -source, expected.Negate(), precise)
}

func testNewOf1MS(t *testing.T, i int, source time.Duration, expected PeriodMS, precise bool) {
	t.Helper()
	g := NewGomegaWithT(t)

	n, p := NewOfWithMS(source)
	rev, _ := expected.Duration()
	info := fmt.Sprintf("%d: source %v expected %+v precise %v rev %v", i, source, expected, precise, rev)
	expectValid(t, n.Period, info)
	if fraction(n.Period.seconds) != 0 {
		g.Expect(n.milliseconds).To(BeZero(), info+" seconds fraction exists")
	}
	g.Expect(n).To(Equal(expected), info)
	g.Expect(p).To(Equal(precise), info)
	if precise {
		g.Expect(rev).To(Equal(source), info)
	}
}

func TestPeriodMSString(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		value  string
		period PeriodMS
	}{
		// note: the negative cases are also covered (see below)

		{"P0D", PeriodMS{}},
		// ones
		{"P1Y", PeriodMS{Period: Period{years: 10}}},
		{"P1M", PeriodMS{Period: Period{months: 10}}},
		{"P1W", PeriodMS{Period: Period{days: 70}}},
		{"P1D", PeriodMS{Period: Period{days: 10}}},
		{"PT1H", PeriodMS{Period: Period{hours: 10}}},
		{"PT1M", PeriodMS{Period: Period{minutes: 10}}},
		{"PT1S", PeriodMS{Period: Period{seconds: 10}}},
		// smallest
		{"P0.1Y", PeriodMS{Period: Period{years: 1}}},
		{"P0.1M", PeriodMS{Period: Period{months: 1}}},
		{"P0.7D", PeriodMS{Period: Period{days: 7}}},
		{"P0.1D", PeriodMS{Period: Period{days: 1}}},
		{"PT0.1H", PeriodMS{Period: Period{hours: 1}}},
		{"PT0.1M", PeriodMS{Period: Period{minutes: 1}}},
		{"PT0.001S", PeriodMS{milliseconds: 1}},
		{"PT0.010S", PeriodMS{milliseconds: 10}},
		{"PT0.100S", PeriodMS{milliseconds: 100}},

		{"P3Y", PeriodMS{Period: Period{years: 30}}},
		{"P6M", PeriodMS{Period: Period{months: 60}}},
		{"P5W", PeriodMS{Period: Period{days: 350}}},
		{"P4W", PeriodMS{Period: Period{days: 280}}},
		{"P4D", PeriodMS{Period: Period{days: 40}}},
		{"PT12H", PeriodMS{Period: Period{hours: 120}}},
		{"PT30M", PeriodMS{Period: Period{minutes: 300}}},
		{"PT5S", PeriodMS{Period: Period{seconds: 50}}},
		{"P3Y6M39DT1H2M4.900S", PeriodMS{
			Period: Period{
				years: 30, months: 60, days: 390, hours: 10, minutes: 20, seconds: 40,
			},
			milliseconds: 900,
		}},

		{"P2.5Y", PeriodMS{Period: Period{years: 25}}},
		{"P2.5M", PeriodMS{Period: Period{months: 25}}},
		{"P2.5D", PeriodMS{Period: Period{days: 25}}},
		{"PT2.5H", PeriodMS{Period: Period{hours: 25}}},
		{"PT2.5M", PeriodMS{Period: Period{minutes: 25}}},
		{"PT2.500S", PeriodMS{Period: Period{seconds: 20}, milliseconds: 500}},
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

func TestPeriodMSToDuration(t *testing.T) {
	cases := []struct {
		value    string
		duration time.Duration
		precise  bool
	}{
		// note: the negative cases are also covered (see below)

		{"P0D", time.Duration(0), true},
		{"PT1S", 1 * time.Second, true},
		{"PT0.100S", 100 * time.Millisecond, true},
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
		testPeriodMSToDuration(t, i, c.value, c.duration, c.precise)
		testPeriodMSToDuration(t, i, "-"+c.value, -c.duration, c.precise)
	}
}

func testPeriodMSToDuration(t *testing.T, i int, value string, duration time.Duration, precise bool) {
	t.Helper()
	g := NewGomegaWithT(t)
	hint := info(i, "%s %s %v", value, duration, precise)
	pp := MustParse(value, false)
	d1, prec := pp.Duration()
	g.Expect(d1).To(Equal(duration), hint)
	g.Expect(prec).To(Equal(precise), hint)
	d2 := pp.DurationApprox()
	if precise {
		g.Expect(d2).To(Equal(duration), hint)
	}
}

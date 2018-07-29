// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package period

import (
	"github.com/rickb777/plural"
	"testing"
	"time"
)

var oneDay = 24 * time.Hour
var oneMonthApprox = 2629746 * time.Second // 30.436875 days
var oneYearApprox = 31556952 * time.Second // 365.2425 days

func TestParsePeriod(t *testing.T) {
	cases := []struct {
		value  string
		period Period
	}{
		{"P0", Period{}},
		{"P0Y", Period{}},
		{"P0M", Period{}},
		{"P0D", Period{}},
		{"PT0H", Period{}},
		{"PT0M", Period{}},
		{"PT0S", Period{}},
		{"P3Y", Period{30, 0, 0, 0, 0, 0}},
		{"P6M", Period{0, 60, 0, 0, 0, 0}},
		{"P5W", Period{0, 0, 350, 0, 0, 0}},
		{"P4D", Period{0, 0, 40, 0, 0, 0}},
		{"PT12H", Period{0, 0, 0, 120, 0, 0}},
		{"PT30M", Period{0, 0, 0, 0, 300, 0}},
		{"PT25S", Period{0, 0, 0, 0, 0, 250}},
		{"P3Y6M5W4DT12H40M5S", Period{30, 60, 390, 120, 400, 50}},
		{"+P3Y6M5W4DT12H40M5S", Period{30, 60, 390, 120, 400, 50}},
		{"-P3Y6M5W4DT12H40M5S", Period{-30, -60, -390, -120, -400, -50}},
		{"P2.Y", Period{20, 0, 0, 0, 0, 0}},
		{"P2.5Y", Period{25, 0, 0, 0, 0, 0}},
		{"P2.15Y", Period{21, 0, 0, 0, 0, 0}},
		{"P2.125Y", Period{21, 0, 0, 0, 0, 0}},
		{"P1Y2.M", Period{10, 20, 0, 0, 0, 0}},
		{"P1Y2.5M", Period{10, 25, 0, 0, 0, 0}},
		{"P1Y2.15M", Period{10, 21, 0, 0, 0, 0}},
		{"P1Y2.125M", Period{10, 21, 0, 0, 0, 0}},
	}
	for i, c := range cases {
		d := MustParse(c.value)
		if d != c.period {
			t.Errorf("%d: MustParsePeriod(%v) == %#v, want (%#v)", i, c.value, d, c.period)
		}
	}

	badCases := []string{
		"13M",
		"P",
	}
	for i, c := range badCases {
		d, err := Parse(c)
		if err == nil {
			t.Errorf("%d: ParsePeriod(%v) == %v", i, c, d)
		}
	}
}

func TestPeriodString(t *testing.T) {
	cases := []struct {
		value  string
		period Period
	}{
		{"P0D", Period{}},
		{"P3Y", Period{30, 0, 0, 0, 0, 0}},
		{"-P3Y", Period{-30, 0, 0, 0, 0, 0}},
		{"P6M", Period{0, 60, 0, 0, 0, 0}},
		{"-P6M", Period{0, -60, 0, 0, 0, 0}},
		{"P5W", Period{0, 0, 350, 0, 0, 0}},
		{"-P5W", Period{0, 0, -350, 0, 0, 0}},
		{"P4W", Period{0, 0, 280, 0, 0, 0}},
		{"-P4W", Period{0, 0, -280, 0, 0, 0}},
		{"P4D", Period{0, 0, 40, 0, 0, 0}},
		{"-P4D", Period{0, 0, -40, 0, 0, 0}},
		{"PT12H", Period{0, 0, 0, 120, 0, 0}},
		{"PT30M", Period{0, 0, 0, 0, 300, 0}},
		{"PT5S", Period{0, 0, 0, 0, 0, 50}},
		{"P3Y6M39DT1H2M4S", Period{30, 60, 390, 10, 20, 40}},
		{"-P3Y6M39DT1H2M4S", Period{-30, -60, -390, 10, 20, 40}},
		{"P2.5Y", Period{25, 0, 0, 0, 0, 0}},
	}
	for i, c := range cases {
		s := c.period.String()
		if s != c.value {
			t.Errorf("%d: String() == %s, want %s for %+v", i, s, c.value, c.period)
		}
	}
}

func TestPeriodIntComponents(t *testing.T) {
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
		{"P39D", 0, 0, 5, 39, 4, 0, 0, 0},
		{"-P39D", 0, 0, -5, -39, -4, 0, 0, 0},
		{"P4D", 0, 0, 0, 4, 4, 0, 0, 0},
		{"-P4D", 0, 0, 0, -4, -4, 0, 0, 0},
		{"PT12H", 0, 0, 0, 0, 0, 12, 0, 0},
		{"PT30M", 0, 0, 0, 0, 0, 0, 30, 0},
		{"PT5S", 0, 0, 0, 0, 0, 0, 0, 5},
	}
	for i, c := range cases {
		p := MustParse(c.value)
		if p.Years() != c.y {
			t.Errorf("%d: %s.Years() == %d, want %d", i, c.value, p.Years(), c.y)
		}
		if p.Months() != c.m {
			t.Errorf("%d: %s.Months() == %d, want %d", i, c.value, p.Months(), c.m)
		}
		if p.Weeks() != c.w {
			t.Errorf("%d: %s.Weeks() == %d, want %d", i, c.value, p.Weeks(), c.w)
		}
		if p.Days() != c.d {
			t.Errorf("%d: %s.Days() == %d, want %d", i, c.value, p.Days(), c.d)
		}
		if p.ModuloDays() != c.dx {
			t.Errorf("%d: %s.ModuloDays() == %d, want %d", i, c.value, p.ModuloDays(), c.dx)
		}
		if p.Hours() != c.hh {
			t.Errorf("%d: %s.Hours() == %d, want %d", i, c.value, p.Hours(), c.hh)
		}
		if p.Minutes() != c.mm {
			t.Errorf("%d: %s.Minutes() == %d, want %d", i, c.value, p.Minutes(), c.mm)
		}
		if p.Seconds() != c.ss {
			t.Errorf("%d: %s.Seconds() == %d, want %d", i, c.value, p.Seconds(), c.ss)
		}
	}
}

func TestPeriodFloatComponents(t *testing.T) {
	cases := []struct {
		value                      string
		y, m, w, d, dx, hh, mm, ss float32
	}{
		{"P0D", 0, 0, 0, 0, 0, 0, 0, 0},

		// YMD cases
		{"P1Y", 1, 0, 0, 0, 0, 0, 0, 0},
		{"P1.1Y", 1.1, 0, 0, 0, 0, 0, 0, 0},
		{"-P1Y", -1, 0, 0, 0, 0, 0, 0, 0},
		{"P1W", 0, 0, 1, 7, 0, 0, 0, 0},
		{"P1.1W", 0, 0, 1.1, 7.7, 0, 0, 0, 0},
		{"-P1W", 0, 0, -1, -7, 0, 0, 0, 0},
		{"P1.1M", 0, 1.1, 0, 0, 0, 0, 0, 0},
		{"P6M", 0, 6, 0, 0, 0, 0, 0, 0},
		{"-P6M", 0, -6, 0, 0, 0, 0, 0, 0},
		{"P39D", 0, 0, 5.571429, 39, 4, 0, 0, 0},
		{"-P39D", 0, 0, -5.571429, -39, -4, 0, 0, 0},
		{"P4D", 0, 0, 0.5714286, 4, 4, 0, 0, 0},
		{"-P4D", 0, 0, -0.5714286, -4, -4, 0, 0, 0},

		// HMS cases
		{"PT1.1H", 0, 0, 0, 0, 0, 1.1, 0, 0},
		{"PT12H", 0, 0, 0, 0, 0, 12, 0, 0},
		{"PT1.1M", 0, 0, 0, 0, 0, 0, 1.1, 0},
		{"PT30M", 0, 0, 0, 0, 0, 0, 30, 0},
		{"PT1.1S", 0, 0, 0, 0, 0, 0, 0, 1.1},
		{"PT5S", 0, 0, 0, 0, 0, 0, 0, 5},
	}
	for i, c := range cases {
		p := MustParse(c.value)
		if p.YearsFloat() != c.y {
			t.Errorf("%d: %s.YearsFloat() == %g, want %g", i, c.value, p.YearsFloat(), c.y)
		}
		if p.MonthsFloat() != c.m {
			t.Errorf("%d: %s.MonthsFloat() == %g, want %g", i, c.value, p.MonthsFloat(), c.m)
		}
		if p.WeeksFloat() != c.w {
			t.Errorf("%d: %s.WeeksFloat() == %g, want %g", i, c.value, p.WeeksFloat(), c.w)
		}
		if p.DaysFloat() != c.d {
			t.Errorf("%d: %s.DaysFloat() == %g, want %g", i, c.value, p.DaysFloat(), c.d)
		}
		if p.HoursFloat() != c.hh {
			t.Errorf("%d: %s.HoursFloat() == %g, want %g", i, c.value, p.HoursFloat(), c.hh)
		}
		if p.MinutesFloat() != c.mm {
			t.Errorf("%d: %s.MinutesFloat() == %g, want %g", i, c.value, p.MinutesFloat(), c.mm)
		}
		if p.SecondsFloat() != c.ss {
			t.Errorf("%d: %s.SecondsFloat() == %g, want %g", i, c.value, p.SecondsFloat(), c.ss)
		}
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
		{"PT3276H", 3276 * time.Hour, true},
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
		if d1 != c.duration {
			t.Errorf("%d: Duration() == %s %v, want %s for %s", i, d1, prec, c.duration, c.value)
		}
		if prec != c.precise {
			t.Errorf("%d: Duration() == %s %v, want %v for %s", i, d1, prec, c.precise, c.value)
		}
		d2 := p.DurationApprox()
		if c.precise && d2 != c.duration {
			t.Errorf("%d: DurationApprox() == %s %v, want %s for %s", i, d2, prec, c.duration, c.value)
		}
	}
}

func TestIsNegative(t *testing.T) {
	cases := []struct {
		value    string
		expected bool
	}{
		{"P0D", false},
		{"PT1S", false},
		{"-PT1S", true},
		{"PT1M", false},
		{"-PT1M", true},
		{"PT1H", false},
		{"-PT1H", true},
		{"P1D", false},
		{"-P1D", true},
		{"P1M", false},
		{"-P1M", true},
		{"P1Y", false},
		{"-P1Y", true},
	}
	for i, c := range cases {
		p := MustParse(c.value)
		got := p.IsNegative()
		if got != c.expected {
			t.Errorf("%d: %v.IsNegative() == %v, want %v", i, p, got, c.expected)
		}
	}
}

func TestPeriodApproxDays(t *testing.T) {
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
		if td != c.approxDays {
			t.Errorf("%d: %v.TotalDaysApprox() == %v, want %v", i, p, td, c.approxDays)
		}
	}
}

func TestPeriodApproxMonths(t *testing.T) {
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
		if td != c.approxMonths {
			t.Errorf("%d: %v.TotalMonthsApprox() == %v, want %v", i, p, td, c.approxMonths)
		}
	}
}

func TestNewPeriod(t *testing.T) {
	cases := []struct {
		years, months, days, hours, minutes, seconds int
		period                                       Period
	}{
		{0, 0, 0, 0, 0, 0, Period{0, 0, 0, 0, 0, 0}},
		{0, 0, 0, 0, 0, 1, Period{0, 0, 0, 0, 0, 10}},
		{0, 0, 0, 0, 1, 0, Period{0, 0, 0, 0, 10, 0}},
		{0, 0, 0, 1, 0, 0, Period{0, 0, 0, 10, 0, 0}},
		{0, 0, 1, 0, 0, 0, Period{0, 0, 10, 0, 0, 0}},
		{0, 1, 0, 0, 0, 0, Period{0, 10, 0, 0, 0, 0}},
		{1, 0, 0, 0, 0, 0, Period{10, 0, 0, 0, 0, 0}},
		{100, 222, 700, 0, 0, 0, Period{1000, 2220, 7000, 0, 0, 0}},
		{0, 0, 0, 0, 0, -1, Period{0, 0, 0, 0, 0, -10}},
		{0, 0, 0, 0, -1, 0, Period{0, 0, 0, 0, -10, 0}},
		{0, 0, 0, -1, 0, 0, Period{0, 0, 0, -10, 0, 0}},
		{0, 0, -1, 0, 0, 0, Period{0, 0, -10, 0, 0, 0}},
		{0, -1, 0, 0, 0, 0, Period{0, -10, 0, 0, 0, 0}},
		{-1, 0, 0, 0, 0, 0, Period{-10, 0, 0, 0, 0, 0}},
	}
	for i, c := range cases {
		p := New(c.years, c.months, c.days, c.hours, c.minutes, c.seconds)
		if p != c.period {
			t.Errorf("%d: %d,%d,%d gives %#v, want %#v", i, c.years, c.months, c.days, p, c.period)
		}
		if p.Years() != c.years {
			t.Errorf("%d: %#v, got %d want %d", i, p, p.Years(), c.years)
		}
		if p.Months() != c.months {
			t.Errorf("%d: %#v, got %d want %d", i, p, p.Months(), c.months)
		}
		if p.Days() != c.days {
			t.Errorf("%d: %#v, got %d want %d", i, p, p.Days(), c.days)
		}
	}
}

func TestNewHMS(t *testing.T) {
	cases := []struct {
		hours, minutes, seconds int
		period                  Period
	}{
		{0, 0, 0, Period{0, 0, 0, 0, 0, 0}},
		{0, 0, 1, Period{0, 0, 0, 0, 0, 10}},
		{0, 1, 0, Period{0, 0, 0, 0, 10, 0}},
		{1, 0, 0, Period{0, 0, 0, 10, 0, 0}},
		{0, 0, -1, Period{0, 0, 0, 0, 0, -10}},
		{0, -1, 0, Period{0, 0, 0, 0, -10, 0}},
		{-1, 0, 0, Period{0, 0, 0, -10, 0, 0}},
	}
	for i, c := range cases {
		p := NewHMS(c.hours, c.minutes, c.seconds)
		if p != c.period {
			t.Errorf("%d: gives %#v, want %#v", i, p, c.period)
		}
		if p.Hours() != c.hours {
			t.Errorf("%d: %#v, got %d want %d", i, p, p.Years(), c.hours)
		}
		if p.Minutes() != c.minutes {
			t.Errorf("%d: %#v, got %d want %d", i, p, p.Months(), c.minutes)
		}
		if p.Seconds() != c.seconds {
			t.Errorf("%d: %#v, got %d want %d", i, p, p.Days(), c.seconds)
		}
	}
}

func TestNewYMD(t *testing.T) {
	cases := []struct {
		years, months, days int
		period              Period
	}{
		{0, 0, 0, Period{0, 0, 0, 0, 0, 0}},
		{0, 0, 1, Period{0, 0, 10, 0, 0, 0}},
		{0, 1, 0, Period{0, 10, 0, 0, 0, 0}},
		{1, 0, 0, Period{10, 0, 0, 0, 0, 0}},
		{100, 222, 700, Period{1000, 2220, 7000, 0, 0, 0}},
		{0, 0, -1, Period{0, 0, -10, 0, 0, 0}},
		{0, -1, 0, Period{0, -10, 0, 0, 0, 0}},
		{-1, 0, 0, Period{-10, 0, 0, 0, 0, 0}},
	}
	for i, c := range cases {
		p := NewYMD(c.years, c.months, c.days)
		if p != c.period {
			t.Errorf("%d: %d,%d,%d gives %#v, want %#v", i, c.years, c.months, c.days, p, c.period)
		}
		if p.Years() != c.years {
			t.Errorf("%d: %#v, got %d want %d", i, p, p.Years(), c.years)
		}
		if p.Months() != c.months {
			t.Errorf("%d: %#v, got %d want %d", i, p, p.Months(), c.months)
		}
		if p.Days() != c.days {
			t.Errorf("%d: %#v, got %d want %d", i, p, p.Days(), c.days)
		}
	}
}

func TestNewOf(t *testing.T) {
	// HMS tests
	testNewOf(t, 100*time.Millisecond, Period{0, 0, 0, 0, 0, 1}, true)
	testNewOf(t, time.Second, Period{0, 0, 0, 0, 0, 10}, true)
	testNewOf(t, time.Minute, Period{0, 0, 0, 0, 10, 0}, true)
	testNewOf(t, time.Hour, Period{0, 0, 0, 10, 0, 0}, true)
	testNewOf(t, time.Hour+time.Minute+time.Second, Period{0, 0, 0, 10, 10, 10}, true)
	testNewOf(t, 24*time.Hour+time.Minute+time.Second, Period{0, 0, 0, 240, 10, 10}, true)
	testNewOf(t, 3276*time.Hour+59*time.Minute+59*time.Second, Period{0, 0, 0, 32760, 590, 590}, true)

	// YMD tests: must be over 3276 hours (approx 4.5 months), otherwise HMS will take care of it
	// first rollover: 3276 hours
	testNewOf(t, 3288*time.Hour, Period{0, 0, 1370, 0, 0, 0}, false)
	testNewOf(t, 3289*time.Hour, Period{0, 0, 1370, 10, 0, 0}, false)
	testNewOf(t, 3277*time.Hour, Period{0, 0, 1360, 130, 0, 0}, false)

	// second rollover: 3276 days
	testNewOf(t, 3277*oneDay, Period{80, 110, 200, 0, 0, 0}, false)
	testNewOf(t, 3277*oneDay+time.Hour+time.Minute+time.Second, Period{80, 110, 200, 10, 0, 0}, false)
	testNewOf(t, 36525*oneDay, Period{1000, 0, 0, 0, 0, 0}, false)
}

func testNewOf(t *testing.T, source time.Duration, expected Period, precise bool) {
	t.Helper()
	testNewOf1(t, source, expected, precise)
	testNewOf1(t, -source, expected.Negate(), precise)
}

func testNewOf1(t *testing.T, source time.Duration, expected Period, precise bool) {
	t.Helper()
	ms := time.Millisecond

	n, p := NewOf(source)
	rev, _ := expected.Duration()
	if n != expected {
		t.Errorf("NewOf(%s) (%dms)\n    gives %-20s %#v,\n     want %-20s (%dms)", source, source/ms, n, n, expected, rev/ms)
	}
	if p != precise {
		t.Errorf("NewOf(%s) (%dms)\n    gives %v,\n     want %v for %v (%dms)", source, source/ms, p, precise, expected, rev/ms)
	}
	//if rev != source {
	//	t.Logf("%d: NewOf(%s) input %dms differs from expected %dms", i, source, source/ms, rev/ms)
	//}
}

func TestBetween(t *testing.T) {
	now := time.Now()

	cases := []struct {
		a, b     time.Time
		expected Period
	}{
		{now, now, Period{0, 0, 0, 0, 0, 0}},

		// simple positive date calculations
		{utc(2015, 1, 1, 0, 0, 0, 0), utc(2015, 1, 1, 0, 0, 0, 100), Period{0, 0, 0, 0, 0, 1}},
		{utc(2015, 1, 1, 0, 0, 0, 0), utc(2015, 2, 2, 1, 1, 1, 1), Period{0, 0, 320, 10, 10, 10}},
		{utc(2015, 2, 1, 0, 0, 0, 0), utc(2015, 3, 2, 1, 1, 1, 1), Period{0, 0, 290, 10, 10, 10}},
		{utc(2015, 3, 1, 0, 0, 0, 0), utc(2015, 4, 2, 1, 1, 1, 1), Period{0, 0, 320, 10, 10, 10}},
		{utc(2015, 4, 1, 0, 0, 0, 0), utc(2015, 5, 2, 1, 1, 1, 1), Period{0, 0, 310, 10, 10, 10}},
		{utc(2015, 5, 1, 0, 0, 0, 0), utc(2015, 6, 2, 1, 1, 1, 1), Period{0, 0, 320, 10, 10, 10}},
		{utc(2015, 6, 1, 0, 0, 0, 0), utc(2015, 7, 2, 1, 1, 1, 1), Period{0, 0, 310, 10, 10, 10}},
		{utc(2015, 1, 1, 0, 0, 0, 0), utc(2015, 7, 2, 1, 1, 1, 1), Period{0, 0, 1820, 10, 10, 10}},

		// less than one month
		{utc(2016, 1, 2, 0, 0, 0, 0), utc(2016, 2, 1, 0, 0, 0, 0), Period{0, 0, 300, 0, 0, 0}},
		{utc(2015, 2, 2, 0, 0, 0, 0), utc(2015, 3, 1, 0, 0, 0, 0), Period{0, 0, 270, 0, 0, 0}}, // non-leap
		{utc(2016, 2, 2, 0, 0, 0, 0), utc(2016, 3, 1, 0, 0, 0, 0), Period{0, 0, 280, 0, 0, 0}}, // leap year
		{utc(2016, 3, 2, 0, 0, 0, 0), utc(2016, 4, 1, 0, 0, 0, 0), Period{0, 0, 300, 0, 0, 0}},
		{utc(2016, 4, 2, 0, 0, 0, 0), utc(2016, 5, 1, 0, 0, 0, 0), Period{0, 0, 290, 0, 0, 0}},
		{utc(2016, 5, 2, 0, 0, 0, 0), utc(2016, 6, 1, 0, 0, 0, 0), Period{0, 0, 300, 0, 0, 0}},
		{utc(2016, 6, 2, 0, 0, 0, 0), utc(2016, 7, 1, 0, 0, 0, 0), Period{0, 0, 290, 0, 0, 0}},

		// BST drops an hour at the daylight-saving transition
		{utc(2015, 1, 1, 0, 0, 0, 0), bst(2015, 7, 2, 1, 1, 1, 1), Period{0, 0, 1820, 0, 10, 10}},

		// negative date calculation
		{utc(2015, 1, 1, 0, 0, 0, 100), utc(2015, 1, 1, 0, 0, 0, 0), Period{0, 0, 0, 0, 0, -1}},
		{utc(2015, 6, 2, 0, 0, 0, 0), utc(2015, 5, 1, 0, 0, 0, 0), Period{0, 0, -320, 0, 0, 0}},
		{utc(2015, 6, 2, 1, 1, 1, 1), utc(2015, 5, 1, 0, 0, 0, 0), Period{0, 0, -320, -10, -10, -10}},

		// daytime only
		{utc(2015, 1, 1, 2, 3, 4, 0), utc(2015, 1, 1, 2, 3, 4, 500), Period{0, 0, 0, 0, 0, 5}},
		{utc(2015, 1, 1, 2, 3, 4, 0), utc(2015, 1, 1, 4, 4, 7, 500), Period{0, 0, 0, 20, 10, 35}},
		{utc(2015, 1, 1, 2, 3, 4, 500), utc(2015, 1, 1, 4, 4, 7, 0), Period{0, 0, 0, 20, 10, 25}},

		// different dates and times
		{utc(2015, 2, 1, 1, 0, 0, 0), utc(2015, 5, 30, 5, 6, 7, 0), Period{0, 0, 1180, 40, 60, 70}},
		{utc(2015, 2, 1, 1, 0, 0, 0), bst(2015, 5, 30, 5, 6, 7, 0), Period{0, 0, 1180, 30, 60, 70}},

		// earlier month in later year
		{utc(2015, 12, 22, 0, 0, 0, 0), utc(2016, 1, 10, 5, 6, 7, 0), Period{0, 0, 190, 50, 60, 70}},
		{utc(2015, 2, 11, 5, 6, 7, 500), utc(2016, 1, 10, 0, 0, 0, 0), Period{0, 0, 3320, 180, 530, 525}},

		// larger ranges
		{utc(2009, 1, 1, 0, 0, 1, 0), utc(2016, 12, 31, 0, 0, 2, 0), Period{0, 0, 29210, 0, 0, 10}},
		{utc(2008, 1, 1, 0, 0, 1, 0), utc(2016, 12, 31, 0, 0, 2, 0), Period{80, 110, 300, 0, 0, 10}},
	}
	for i, c := range cases {
		n := Between(c.a, c.b)
		if n != c.expected {
			t.Errorf("%d: Between(%v, %v)\n  gives %-20s %#v,\n   want %-20s %#v", i, c.a, c.b, n, n, c.expected, c.expected)
		}
	}
}

func TestNormalise(t *testing.T) {
	// zero cases
	testNormalise(t, New(0, 0, 0, 0, 0, 0), New(0, 0, 0, 0, 0, 0), New(0, 0, 0, 0, 0, 0))

	// carry seconds to minutes
	testNormalise(t, Period{0, 0, 0, 0, 0, 699}, Period{0, 0, 0, 0, 10, 99}, Period{0, 0, 0, 0, 10, 99})

	// carry minutes to seconds
	testNormalise(t, Period{0, 0, 0, 0, 5, 0}, Period{0, 0, 0, 0, 0, 300}, Period{0, 0, 0, 0, 0, 300})
	testNormalise(t, Period{0, 0, 0, 0, 1, 0}, Period{0, 0, 0, 0, 0, 60}, Period{0, 0, 0, 0, 0, 60})
	testNormalise(t, Period{0, 0, 0, 0, 55, 0}, Period{0, 0, 0, 0, 50, 300}, Period{0, 0, 0, 0, 50, 300})

	// carry minutes to hours
	testNormalise(t, Period{0, 0, 0, 0, 699, 0}, Period{0, 0, 0, 10, 90, 540}, Period{0, 0, 0, 10, 90, 540})

	// carry hours to minutes
	testNormalise(t, Period{0, 0, 0, 5, 0, 0}, Period{0, 0, 0, 0, 300, 0}, Period{0, 0, 0, 0, 300, 0})

	// carry hours to days
	testNormalise(t, Period{0, 0, 0, 249, 0, 0}, Period{0, 0, 0, 240, 540, 0}, Period{0, 0, 0, 240, 540, 0})
	testNormalise(t, Period{0, 0, 0, 249, 0, 0}, Period{0, 0, 0, 240, 540, 0}, Period{0, 0, 0, 240, 540, 0})
	testNormalise(t, Period{0, 0, 0, 369, 0, 0}, Period{0, 0, 0, 360, 540, 0}, Period{0, 0, 10, 120, 540, 0})
	testNormalise(t, Period{0, 0, 0, 249, 0, 10}, Period{0, 0, 0, 240, 540, 10}, Period{0, 0, 0, 240, 540, 10})

	// carry days to hours
	testNormalise(t, Period{0, 0, 5, 30, 0, 0}, Period{0, 0, 0, 150, 00, 0}, Period{0, 0, 0, 150, 0, 0})

	// carry months to years
	testNormalise(t, Period{0, 125, 0, 0, 0, 0}, Period{0, 125, 0, 0, 0, 0}, Period{0, 125, 0, 0, 0, 0})
	testNormalise(t, Period{0, 131, 0, 0, 0, 0}, Period{10, 11, 0, 0, 0, 0}, Period{10, 11, 0, 0, 0, 0})

	// carry days to months
	testNormalise(t, Period{0, 0, 323, 0, 0, 0}, Period{0, 0, 323, 0, 0, 0}, Period{0, 0, 323, 0, 0, 0})

	// carry months to days
	testNormalise(t, Period{0, 5, 203, 0, 0, 0}, Period{0, 0, 355, 0, 0, 0}, Period{0, 10, 50, 0, 0, 0})

	// full ripple up
	testNormalise(t, Period{0, 121, 305, 239, 591, 601}, Period{10, 0, 330, 360, 540, 61}, Period{10, 10, 40, 0, 540, 61})

	// carry years to months
	testNormalise(t, Period{5, 0, 0, 0, 0, 0}, Period{0, 60, 0, 0, 0, 0}, Period{0, 60, 0, 0, 0, 0})
	testNormalise(t, Period{5, 25, 0, 0, 0, 0}, Period{0, 85, 0, 0, 0, 0}, Period{0, 85, 0, 0, 0, 0})
	testNormalise(t, Period{5, 20, 10, 0, 0, 0}, Period{0, 80, 10, 0, 0, 0}, Period{0, 80, 10, 0, 0, 0})
}

func testNormalise(t *testing.T, source, precise, approx Period) {
	t.Helper()

	testNormaliseBothSigns(t, source, precise, true)
	testNormaliseBothSigns(t, source, approx, false)
}

func testNormaliseBothSigns(t *testing.T, source, expected Period, precise bool) {
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
	if n2 != eneg {
		t.Errorf("%v.Normalise(%v) %s\n   gives %-22s %#v %s,\n    want %-22s %#v %s",
			sneg, precise, sneg.DurationApprox(),
			n2, n2, n2.DurationApprox(),
			eneg, eneg, eneg.DurationApprox())
	}
}

func TestPeriodFormat(t *testing.T) {
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
		if s != c.expect {
			t.Errorf("%d: Format() == %s, want %s for %+v", i, s, c.expect, c.period)
		}
	}
}

func TestPeriodScale(t *testing.T) {
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
		if s != MustParse(c.expect) {
			t.Errorf("%d: %s.Scale(%g) == %v, want %s", i, c.one, c.m, s, c.expect)
		}
	}
}

func TestPeriodAdd(t *testing.T) {
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
		if s != MustParse(c.expect) {
			t.Errorf("%d: %s.Add(%s) == %v, want %s", i, c.one, c.two, s, c.expect)
		}
	}
}

func TestPeriodFormatWithoutWeeks(t *testing.T) {
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
		if s != c.expect {
			t.Errorf("%d: Format() == %s, want %s for %+v", i, s, c.expect, c.period)
		}
	}
}

func TestPeriodParseOnlyYMD(t *testing.T) {
	cases := []struct {
		one    string
		expect string
	}{
		{"P1Y2M3DT4H5M6S", "P1Y2M3D"},
		{"-P6Y5M4DT3H2M1S", "-P6Y5M4D"},
	}
	for i, c := range cases {
		s := MustParse(c.one).OnlyYMD()
		if s != MustParse(c.expect) {
			t.Errorf("%d: %s.OnlyYMD() == %v, want %s", i, c.one, s, c.expect)
		}
	}
}

func TestPeriodParseOnlyHMS(t *testing.T) {
	cases := []struct {
		one    string
		expect string
	}{
		{"P1Y2M3DT4H5M6S", "PT4H5M6S"},
		{"-P6Y5M4DT3H2M1S", "-PT3H2M1S"},
	}
	for i, c := range cases {
		s := MustParse(c.one).OnlyHMS()
		if s != MustParse(c.expect) {
			t.Errorf("%d: %s.OnlyHMS() == %v, want %s", i, c.one, s, c.expect)
		}
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

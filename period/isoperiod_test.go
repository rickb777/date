package period

import (
	"github.com/rickb777/plural"
	"testing"
	"time"
)

func TestParsePeriod(t *testing.T) {
	cases := []struct {
		value  string
		period Period
	}{
		{"P0", Period{}},
		{"P0D", Period{}},
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
	for _, c := range cases {
		d := MustParse(c.value)
		if d != c.period {
			t.Errorf("MustParsePeriod(%v) == %#v, want (%#v)", c.value, d, c.period)
		}
	}

	badCases := []string{
		"13M",
		"P",
	}
	for _, c := range badCases {
		d, err := Parse(c)
		if err == nil {
			t.Errorf("ParsePeriod(%v) == %v", c, d)
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
		{"P35D", Period{0, 0, 350, 0, 0, 0}},
		{"-P35D", Period{0, 0, -350, 0, 0, 0}},
		{"P4D", Period{0, 0, 40, 0, 0, 0}},
		{"-P4D", Period{0, 0, -40, 0, 0, 0}},
		{"PT12H", Period{0, 0, 0, 120, 0, 0}},
		{"PT30M", Period{0, 0, 0, 0, 300, 0}},
		{"PT5S", Period{0, 0, 0, 0, 0, 50}},
		{"P3Y6M39DT1H2M4S", Period{30, 60, 390, 10, 20, 40}},
		{"-P3Y6M39DT1H2M4S", Period{-30, -60, -390, 10, 20, 40}},
		{"P2.5Y", Period{25, 0, 0, 0, 0, 0}},
	}
	for _, c := range cases {
		s := c.period.String()
		if s != c.value {
			t.Errorf("String() == %s, want %s for %+v", s, c.value, c.period)
		}
	}
}

func TestPeriodComponents(t *testing.T) {
	cases := []struct {
		value                      string
		y, m, w, d, dx, hh, mm, ss int
	}{
		{"P0D", 0, 0, 0, 0, 0, 0, 0, 0},
		{"P1Y", 1, 0, 0, 0, 0, 0, 0, 0},
		{"-P1Y", -1, 0, 0, 0, 0, 0, 0, 0},
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
	for _, c := range cases {
		p := MustParse(c.value)
		if p.Years() != c.y {
			t.Errorf("%s.Years() == %d, want %d", c.value, p.Years(), c.y)
		}
		if p.Months() != c.m {
			t.Errorf("%s.Months() == %d, want %d", c.value, p.Months(), c.m)
		}
		if p.Weeks() != c.w {
			t.Errorf("%s.Weeks() == %d, want %d", c.value, p.Weeks(), c.w)
		}
		if p.Days() != c.d {
			t.Errorf("%s.Days() == %d, want %d", c.value, p.Days(), c.d)
		}
		if p.ModuloDays() != c.dx {
			t.Errorf("%s.ModuloDays() == %d, want %d", c.value, p.ModuloDays(), c.dx)
		}
		if p.Hours() != c.hh {
			t.Errorf("%s.Hours() == %d, want %d", c.value, p.Hours(), c.hh)
		}
		if p.Minutes() != c.mm {
			t.Errorf("%s.Minutes() == %d, want %d", c.value, p.Minutes(), c.mm)
		}
		if p.Seconds() != c.ss {
			t.Errorf("%s.Seconds() == %d, want %d", c.value, p.Seconds(), c.ss)
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
		{"PT1S", time.Duration(1 * time.Second), true},
		{"PT1M", time.Duration(60 * time.Second), true},
		{"PT1H", time.Duration(3600 * time.Second), true},
		{"P1D", time.Duration(24 * time.Hour), false},
		{"P1M", time.Duration(2629800 * time.Second), false},
		{"P1Y", time.Duration(31557600 * time.Second), false},
		{"-P1Y", -time.Duration(31557600 * time.Second), false},
	}
	for _, c := range cases {
		s, p := MustParse(c.value).Duration()
		if s != c.duration {
			t.Errorf("Duration() == %s %v, want %s for %+v", s, p, c.duration, c.value)
		}
		if p != c.precise {
			t.Errorf("Duration() == %s %v, want %v for %+v", s, p, c.precise, c.value)
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
	for _, c := range cases {
		p := New(c.years, c.months, c.days, c.hours, c.minutes, c.seconds)
		if p != c.period {
			t.Errorf("%d,%d,%d gives %#v, want %#v", c.years, c.months, c.days, p, c.period)
		}
		if p.Years() != c.years {
			t.Errorf("%#v, got %d want %d", p, p.Years(), c.years)
		}
		if p.Months() != c.months {
			t.Errorf("%#v, got %d want %d", p, p.Months(), c.months)
		}
		if p.Days() != c.days {
			t.Errorf("%#v, got %d want %d", p, p.Days(), c.days)
		}
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
	for _, c := range cases {
		s := MustParse(c.period).Format()
		if s != c.expect {
			t.Errorf("Format() == %s, want %s for %+v", s, c.expect, c.period)
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
	for _, c := range cases {
		s := MustParse(c.period).FormatWithPeriodNames(PeriodYearNames, PeriodMonthNames, plural.Plurals{}, PeriodDayNames,
			PeriodHourNames, PeriodMinuteNames, PeriodSecondNames)
		if s != c.expect {
			t.Errorf("Format() == %s, want %s for %+v", s, c.expect, c.period)
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
	for _, c := range cases {
		s := MustParse(c.one).Add(MustParse(c.two))
		if s != MustParse(c.expect) {
			t.Errorf("%s.Add(%s) == %v, want %s", c.one, c.two, s, c.expect)
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
		{"P1M", 2, "P2M"},
		{"P1Y", 2, "P2Y"},
		{"PT1H", 2, "PT2H"},
		{"PT1M", 2, "PT2M"},
		{"PT1S", 2, "PT2S"},
		{"P1D", 0.5, "P0.5D"},
		{"P1M", 0.5, "P0.5M"},
		{"P1Y", 0.5, "P0.5Y"},
		{"PT1H", 0.5, "PT0.5H"},
		{"PT1M", 0.5, "PT0.5M"},
		{"PT1S", 0.5, "PT0.5S"},
		{"P1Y2M3DT4H5M6S", 2, "P2Y4M6DT8H10M12S"},
		{"P2Y4M6DT8H10M12S", -0.5, "-P1Y2M3DT4H5M6S"},
	}
	for _, c := range cases {
		s := MustParse(c.one).Scale(c.m)
		if s != MustParse(c.expect) {
			t.Errorf("%s.Scale(%g) == %v, want %s", c.one, c.m, s, c.expect)
		}
	}
}

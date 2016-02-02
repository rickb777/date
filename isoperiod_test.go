package date

import (
	"testing"
)

func TestParsePeriod(t *testing.T) {
	cases := []struct {
		value  string
		period Period
	}{
		{"P0", Period{}},
		{"P0D", Period{}},
		{"P3Y", Period{3000, 0, 0}},
		{"P6M", Period{0, 6000, 0}},
		{"P5W", Period{0, 0, 35000}},
		{"P4D", Period{0, 0, 4000}},
		//{"PT12H", Period{}},
		//{"PT30M", Period{}},
		//{"PT5S", Period{}},
		{"P3Y6M5W4DT12H30M5S", Period{3000, 6000, 39000}},
		{"+P3Y6M5W4DT12H30M5S", Period{3000, 6000, 39000}},
		{"-P3Y6M5W4DT12H30M5S", Period{-3000, -6000, -39000}},
		{"P2.Y", Period{2000, 0, 0}},
		{"P2.5Y", Period{2500, 0, 0}},
		{"P2.15Y", Period{2150, 0, 0}},
		{"P2.125Y", Period{2125, 0, 0}},
	}
	for _, c := range cases {
		d := MustParsePeriod(c.value)
		if d != c.period {
			t.Errorf("MustParsePeriod(%v) == %#v, want (%#v)", c.value, d, c.period)
		}
	}

	badCases := []string{
		"13M",
		"P",
	}
	for _, c := range badCases {
		d, err := ParsePeriod(c)
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
		{"P3Y", Period{3000, 0, 0}},
		{"-P3Y", Period{-3000, 0, 0}},
		{"P6M", Period{0, 6000, 0}},
		{"-P6M", Period{0, -6000, 0}},
		{"P35D", Period{0, 0, 35000}},
		{"-P35D", Period{0, 0, -35000}},
		{"P4D", Period{0, 0, 4000}},
		{"-P4D", Period{0, 0, -4000}},
		//{"PT12H", Period{}},
		//{"PT30M", Period{}},
		//{"PT5S", Period{}},
		{"P3Y6M39D", Period{3000, 6000, 39000}},
		{"-P3Y6M39D", Period{-3000, -6000, -39000}},
		{"P2.5Y", Period{2500, 0, 0}},
		{"P2.15Y", Period{2150, 0, 0}},
		{"P2.125Y", Period{2125, 0, 0}},
	}
	for _, c := range cases {
		s := c.period.String()
		if s != c.value {
			t.Errorf("String() == %s, want %s for %+v", s, c.value, c.period)
		}
	}
}

func TestNewPeriod(t *testing.T) {
	cases := []struct {
		years, months, days int
		period              Period
	}{
		{0, 0, 0, Period{0, 0, 0}},
		{0, 0, 1, Period{0, 0, 1000}},
		{0, 1, 0, Period{0, 1000, 0}},
		{1, 0, 0, Period{1000, 0, 0}},
		{100, 222, 700, Period{100000, 222000, 700000}},
		{0, 0, -1, Period{0, 0, -1000}},
		{0, -1, 0, Period{0, -1000, 0}},
		{-1, 0, 0, Period{-1000, 0, 0}},
	}
	for _, c := range cases {
		p := NewPeriod(c.years, c.months, c.days)
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

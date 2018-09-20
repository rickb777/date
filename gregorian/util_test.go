package gregorian

import (
	"testing"
	"time"
)

func TestIsLeap(t *testing.T) {
	cases := []struct {
		year     int
		expected bool
	}{
		{0, true}, // year zero is not defined under some conventions but is in ISO8601
		{2000, true},
		{2400, true},
		{2001, false},
		{2002, false},
		{2003, false},
		{2003, false},
		{2004, true},
		{2005, false},
		{1800, false},
		{1900, false},
		{2200, false},
		{2300, false},
		{2500, false},
	}
	for _, c := range cases {
		got := IsLeap(c.year)
		if got != c.expected {
			t.Errorf("TestIsLeap(%d) == %v, want %v", c.year, got, c.expected)
		}
	}
}

func TestDaysInYear(t *testing.T) {
	cases := []struct {
		year     int
		expected int
	}{
		{2000, 366},
		{2001, 365},
	}
	for _, c := range cases {
		got1 := DaysInYear(c.year)
		if got1 != c.expected {
			t.Errorf("DaysInYear(%d) == %v, want %v", c.year, got1, c.expected)
		}
	}
}

func TestDaysIn(t *testing.T) {
	cases := []struct {
		year     int
		month    time.Month
		expected int
	}{
		{2000, time.January, 31},
		{2000, time.February, 29},
		{2001, time.February, 28},
		{2001, time.April, 30},
		{2001, time.May, 31},
		{2001, time.June, 30},
		{2001, time.July, 31},
		{2001, time.August, 31},
		{2001, time.September, 30},
		{2001, time.October, 31},
		{2001, time.November, 30},
		{2001, time.December, 31},
	}
	for _, c := range cases {
		got1 := DaysIn(c.year, c.month)
		if got1 != c.expected {
			t.Errorf("DaysIn(%d, %d) == %v, want %v", c.year, c.month, got1, c.expected)
		}
	}
}

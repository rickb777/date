// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package clock

import (
	"testing"
	"time"
)

func TestClockHoursMinutesSeconds(t *testing.T) {
	cases := []struct {
		in      Clock
		h, m, s int
	}{
		{New(0, 0, 0), 0, 0, 0},
		{New(1, 2, 3), 1, 2, 3},
		{New(23, 59, 59), 23, 59, 59},
		{New(0, 0, -1), 23, 59, 59},
	}
	for _, c := range cases {
		h := c.in.Hours()
		m := c.in.Minutes()
		s := c.in.Seconds()
		if h != c.h || m != c.m || s != c.s {
			t.Errorf("got %d %d %d, want %v", h, m, s, c.in)
		}
	}
}

func TestClockIsInOneDay(t *testing.T) {
	cases := []struct {
		in   Clock
		want bool
	}{
		{New(0, 0, 0), true},
		{New(24, 0, 0), true},
		{New(-24, 0, 0), false},
		{New(48, 0, 0), false},
		{New(0, 0, 1), true},
		{New(2, 0, 1), true},
		{New(-1, 0, 0), false},
		{New(0, 0, -1), false},
	}
	for _, c := range cases {
		got := c.in.IsInOneDay()
		if got != c.want {
			t.Errorf("%v got %v, want %v", c.in, c.in.IsInOneDay(), c.want)
		}
	}
}

func TestClockAdd(t *testing.T) {
	cases := []struct {
		h, m, s  int
		in, want Clock
	}{
		{0, 0, 0, 2 * ClockHour, New(2, 0, 0)},
		{0, 0, 1, 2 * ClockHour, New(2, 0, 1)},
		{0, 0, -1, 2 * ClockHour, New(1, 59, 59)},
		{0, 1, 0, 2 * ClockHour, New(2, 1, 0)},
		{0, -1, 0, 2 * ClockHour, New(1, 59, 0)},
		{1, 0, 0, 2 * ClockHour, New(3, 0, 0)},
		{-1, 0, 0, 2 * ClockHour, New(1, 0, 0)},
		{-2, 0, 0, 2 * ClockHour, New(0, 0, 0)},
		{-2, 0, -1, 2 * ClockHour, New(0, 0, -1)},
	}
	for _, c := range cases {
		got := c.in.Add(c.h, c.m, c.s)
		if got != c.want {
			t.Errorf("%d %d %d: got %v, want %v", c.h, c.m, c.s, got, c.want)
		}
	}
}

func TestClockIsMidnight(t *testing.T) {
	cases := []struct {
		in   Clock
		want bool
	}{
		{New(0, 0, 0), true},
		{ClockDay, true},
		{24 * ClockHour, true},
		{New(24, 0, 0), true},
		{New(-24, 0, 0), true},
		{New(-48, 0, 0), true},
		{New(48, 0, 0), true},
		{New(0, 0, 1), false},
		{New(2, 0, 1), false},
		{New(-1, 0, 0), false},
		{New(0, 0, -1), false},
	}
	for i, c := range cases {
		got := c.in.IsMidnight()
		if got != c.want {
			t.Errorf("%d: %v got %v, want %v, %d", i, c.in, c.in.IsMidnight(), c.want, c.in.Mod24())
		}
	}
}

func TestClockMod(t *testing.T) {
	cases := []struct {
		h, want Clock
	}{
		{0, 0},
		{1 * ClockHour, 1 * ClockHour},
		{2 * ClockHour, 2 * ClockHour},
		{23 * ClockHour, 23 * ClockHour},
		{24 * ClockHour, 0},
		{-24 * ClockHour, 0},
		{-48 * ClockHour, 0},
		{25 * ClockHour, ClockHour},
		{49 * ClockHour, ClockHour},
		{-1 * ClockHour, 23 * ClockHour},
		{-23 * ClockHour, ClockHour},
		{New(0, 0, 1), ClockSecond},
		{New(0, 0, -1), New(23, 59, 59)},
	}
	for i, c := range cases {
		clock := c.h
		got := clock.Mod24()
		if got != c.want {
			t.Errorf("%d: %dh: got %#v, want %#v", i, c.h, got, c.want)
		}
	}
}

func TestClockDays(t *testing.T) {
	cases := []struct {
		h, days int
	}{
		{0, 0},
		{1, 0},
		{23, 0},
		{24, 1},
		{25, 1},
		{48, 2},
		{49, 2},
		{-1, -1},
		{-23, -1},
		{-24, -2},
	}
	for i, c := range cases {
		clock := Clock(c.h) * ClockHour
		if clock.Days() != c.days {
			t.Errorf("%d: %dh: got %v, want %v", i, c.h, clock.Days(), c.days)
		}
	}
}

func TestClockString(t *testing.T) {
	cases := []struct {
		h, m, s, ns           time.Duration
		hh, hhmm, hhmmss, str string
	}{
		{0, 0, 0, 0, "00", "00:00", "00:00:00", "00:00:00.000000000"},
		{0, 0, 0, 1, "00", "00:00", "00:00:00", "00:00:00.000000001"},
		{0, 0, 1, 0, "00", "00:00", "00:00:01", "00:00:01.000000000"},
		{0, 1, 0, 0, "00", "00:01", "00:01:00", "00:01:00.000000000"},
		{1, 0, 0, 0, "01", "01:00", "01:00:00", "01:00:00.000000000"},
		{1, 2, 3, 4, "01", "01:02", "01:02:03", "01:02:03.000000004"},
		{-1, -1, -1, -1, "22", "22:58", "22:58:58", "22:58:58.999999999"},
	}
	for _, c := range cases {
		d := Clock(c.h * time.Hour + c.m * time.Minute + c.s * time.Second + c.ns)
		if d.Hh() != c.hh {
			t.Errorf("%d, %d, %d, %d, got %v, want %v", c.h, c.m, c.s, c.ns, d.Hh(), c.hh)
		}
		if d.HhMm() != c.hhmm {
			t.Errorf("%d, %d, %d, %d, got %v, want %v", c.h, c.m, c.s, c.ns, d.HhMm(), c.hhmm)
		}
		if d.HhMmSs() != c.hhmmss {
			t.Errorf("%d, %d, %d, %d, got %v, want %v", c.h, c.m, c.s, c.ns, d.HhMmSs(), c.hhmmss)
		}
		if d.String() != c.str {
			t.Errorf("%d, %d, %d, %d, got %v, want %v", c.h, c.m, c.s, c.ns, d.String(), c.str)
		}
	}
}

func TestClockParse(t *testing.T) {
	cases := []struct {
		str  string
		want Clock
	}{
		{"00", New(0, 0, 0)},
		{"01", New(1, 0, 0)},
		{"23", New(23, 0, 0)},
		{"00:00", New(0, 0, 0)},
		{"00:01", New(0, 1, 0)},
		{"01:00", New(1, 0, 0)},
		{"01:02", New(1, 2, 0)},
		{"23:59", New(23, 59, 0)},
		{"00:00:00", New(0, 0, 0)},
		{"00:00:01", New(0, 0, 1)},
		{"00:01:00", New(0, 1, 0)},
		{"01:00:00", New(1, 0, 0)},
		{"01:02:03", New(1, 2, 3)},
		{"23:59:59", New(23, 59, 59)},
		{"00:00:00.000000000", New(0, 0, 0)},
		{"00:00:00.000000001", New(0, 0, 0) + 1},
		{"00:00:01.000000000", New(0, 0, 1)},
		{"00:01:00.000000000", New(0, 1, 0)},
		{"01:00:00.000000000", New(1, 0, 0)},
		{"01:02:03.000000004", New(1, 2, 3) + 4},
		{"23:59:59.999999999", New(23, 59, 59) + 999999999},
	}
	for _, c := range cases {
		str, err := Parse(c.str)
		if err != nil {
			t.Errorf("%s, error %v", c.str, err)
		}
		if str != c.want {
			t.Errorf("%s, got %v, want %v", c.str, str, c.want)
		}
	}
}

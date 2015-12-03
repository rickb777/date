// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package date

import (
	"time"
	"testing"
)

func TestClockHoursMinutesSeconds(t *testing.T) {
	cases := []struct {
		in      Clock
		h, m, s int
	}{
		{HhMmSs(0, 0, 0), 0, 0, 0},
		{HhMmSs(1, 2, 3), 1, 2, 3},
		{HhMmSs(23, 59, 59), 23, 59, 59},
		{HhMmSs(0, 0, -1), 23, 59, 59},
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

func TestClockAdd(t *testing.T) {
	cases := []struct {
		h, m, s  int
		in, want Clock
	}{
		{0, 0, 0, 2 * ClockHour, HhMmSs(2, 0, 0)},
		{0, 0, 1, 2 * ClockHour, HhMmSs(2, 0, 1)},
		{0, 0, -1, 2 * ClockHour, HhMmSs(1, 59, 59)},
		{0, 1, 0, 2 * ClockHour, HhMmSs(2, 1, 0)},
		{0, -1, 0, 2 * ClockHour, HhMmSs(1, 59, 0)},
		{1, 0, 0, 2 * ClockHour, HhMmSs(3, 0, 0)},
		{-1, 0, 0, 2 * ClockHour, HhMmSs(1, 0, 0)},
		{-2, 0, 0, 2 * ClockHour, HhMmSs(0, 0, 0)},
		{-2, 0, -1, 2 * ClockHour, HhMmSs(0, 0, -1)},
	}
	for _, c := range cases {
		got := c.in.Add(c.h, c.m, c.s)
		if got != c.want {
			t.Errorf("%d %d %d: got %v, want %v", c.h, c.m, c.s, got, c.want)
		}
	}
}

func TestClockMod(t *testing.T) {
	cases := []struct {
		h, mod Clock
	}{
		{0, 0},
		{1, 1 * ClockHour},
		{2, 2 * ClockHour},
		{23, 23 * ClockHour},
		{24, 0},
		{25, ClockHour},
		{49, ClockHour},
		{-1, 23 * ClockHour},
		{-23, ClockHour},
	}
	for _, c := range cases {
		clock := c.h * ClockHour
		if clock.Mod24() != c.mod {
			t.Errorf("%dh: got %v, want %v", c.h, clock.Mod24(), c.mod)
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
	for _, c := range cases {
		clock := Clock(c.h) * ClockHour
		if clock.Days() != c.days {
			t.Errorf("%dh: got %v, want %v", c.h, clock.Days(), c.days)
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
		{"00", HhMmSs(0, 0, 0)},
		{"01", HhMmSs(1, 0, 0)},
		{"23", HhMmSs(23, 0, 0)},
		{"00:00", HhMmSs(0, 0, 0)},
		{"00:01", HhMmSs(0, 1, 0)},
		{"01:00", HhMmSs(1, 0, 0)},
		{"01:02", HhMmSs(1, 2, 0)},
		{"23:59", HhMmSs(23, 59, 0)},
		{"00:00:00", HhMmSs(0, 0, 0)},
		{"00:00:01", HhMmSs(0, 0, 1)},
		{"00:01:00", HhMmSs(0, 1, 0)},
		{"01:00:00", HhMmSs(1, 0, 0)},
		{"01:02:03", HhMmSs(1, 2, 3)},
		{"23:59:59", HhMmSs(23, 59, 59)},
		{"00:00:00.000000000", HhMmSs(0, 0, 0)},
		{"00:00:00.000000001", HhMmSs(0, 0, 0) + 1},
		{"00:00:01.000000000", HhMmSs(0, 0, 1)},
		{"00:01:00.000000000", HhMmSs(0, 1, 0)},
		{"01:00:00.000000000", HhMmSs(1, 0, 0)},
		{"01:02:03.000000004", HhMmSs(1, 2, 3) + 4},
		{"23:59:59.999999999", HhMmSs(23, 59, 59) + 999999999},
	}
	for _, c := range cases {
		str, err := ParseClock(c.str)
		if err != nil {
			t.Errorf("%s, error %v", c.str, err)
		}
		if str != c.want {
			t.Errorf("%s, got %v, want %v", c.str, str, c.want)
		}
	}
}

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
		in          Clock
		h, m, s, ms int
	}{
		{New(0, 0, 0, 0), 0, 0, 0, 0},
		{New(1, 2, 3, 4), 1, 2, 3, 4},
		{New(23, 59, 59, 999), 23, 59, 59, 999},
		{New(0, 0, 0, -1), 23, 59, 59, 999},
		{NewAt(time.Date(2015, 12, 4, 18, 50, 42, 173444111, time.UTC)), 18, 50, 42, 173},
	}
	for i, x := range cases {
		h := x.in.Hours()
		m := x.in.Minutes()
		s := x.in.Seconds()
		ms := x.in.Millisec()
		if h != x.h || m != x.m || s != x.s || ms != x.ms {
			t.Errorf("%d: got %02d:%02d:%02d.%03d, want %v (%d)", i, h, m, s, ms, x.in, x.in)
		}
	}
}

func TestClockSinceMidnight(t *testing.T) {
	cases := []struct {
		in Clock
		d  time.Duration
	}{
		{New(1, 2, 3, 4), time.Hour + 2*time.Minute + 3*time.Second + 4},
		{New(23, 59, 59, 999), 23*time.Hour + 59*time.Minute + 59*time.Second + 999},
	}
	for i, x := range cases {
		d := x.in.DurationSinceMidnight()
		c2 := SinceMidnight(d)
		if c2 != x.in {
			t.Errorf("%d: got %v, want %v (%d)", i, c2, x.in, x.in)
		}
	}
}

func TestClockIsInOneDay(t *testing.T) {
	cases := []struct {
		in   Clock
		want bool
	}{
		{New(0, 0, 0, 0), true},
		{New(24, 0, 0, 0), true},
		{New(-24, 0, 0, 0), false},
		{New(48, 0, 0, 0), false},
		{New(0, 0, 0, 1), true},
		{New(2, 0, 0, 1), true},
		{New(-1, 0, 0, 0), false},
		{New(0, 0, 0, -1), false},
	}
	for _, x := range cases {
		got := x.in.IsInOneDay()
		if got != x.want {
			t.Errorf("%v got %v, want %v", x.in, x.in.IsInOneDay(), x.want)
		}
	}
}

func TestClockAdd(t *testing.T) {
	cases := []struct {
		h, m, s, ms int
		in, want    Clock
	}{
		{0, 0, 0, 0, 2 * Hour, New(2, 0, 0, 0)},
		{0, 0, 0, 1, 2 * Hour, New(2, 0, 0, 1)},
		{0, 0, 0, -1, 2 * Hour, New(1, 59, 59, 999)},
		{0, 0, 1, 0, 2 * Hour, New(2, 0, 1, 0)},
		{0, 0, -1, 0, 2 * Hour, New(1, 59, 59, 0)},
		{0, 1, 0, 0, 2 * Hour, New(2, 1, 0, 0)},
		{0, -1, 0, 0, 2 * Hour, New(1, 59, 0, 0)},
		{1, 0, 0, 0, 2 * Hour, New(3, 0, 0, 0)},
		{-1, 0, 0, 0, 2 * Hour, New(1, 0, 0, 0)},
		{-2, 0, 0, 0, 2 * Hour, New(0, 0, 0, 0)},
		{-2, 0, -1, -1, 2 * Hour, New(0, 0, -1, -1)},
	}
	for i, x := range cases {
		got := x.in.Add(x.h, x.m, x.s, x.ms)
		if got != x.want {
			t.Errorf("%d: %d %d %d.%d: got %v, want %v", i, x.h, x.m, x.s, x.ms, got, x.want)
		}
	}
}

func TestClockAddDuration(t *testing.T) {
	cases := []struct {
		d        time.Duration
		in, want Clock
	}{
		{0, 2 * Hour, New(2, 0, 0, 0)},
		{1, 2 * Hour, New(2, 0, 0, 0)},
		{time.Millisecond, 2 * Hour, New(2, 0, 0, 1)},
		{-time.Second, 2 * Hour, New(1, 59, 59, 0)},
		{7 * time.Minute, 2 * Hour, New(2, 7, 0, 0)},
	}
	for i, x := range cases {
		got := x.in.AddDuration(x.d)
		if got != x.want {
			t.Errorf("%d: %d: got %v, want %v", i, x.d, got, x.want)
		}
	}
}

func TestClockSubtract(t *testing.T) {
	cases := []struct {
		c1, c2 Clock
		want   time.Duration
	}{
		{New(1, 2, 3, 4), New(1, 2, 3, 4), 0 * time.Hour},
		{New(2, 0, 0, 0), New(0, 0, 0, 0), 2 * time.Hour},
		{New(0, 0, 0, 0), New(2, 0, 0, 0), 22 * time.Hour},
		{New(1, 0, 0, 0), New(23, 0, 0, 0), 2 * time.Hour},
		{New(23, 0, 0, 0), New(1, 0, 0, 0), 22 * time.Hour},
		{New(1, 2, 3, 5), New(1, 2, 3, 4), 1 * time.Millisecond},
		{New(1, 2, 3, 4), New(1, 2, 3, 5), 24*time.Hour - 1*time.Millisecond},
	}
	for i, x := range cases {
		got := x.c1.ModSubtract(x.c2)
		if got != x.want {
			t.Errorf("%d: %v - %v: got %v, want %v", i, x.c1, x.c2, got, x.want)
		}
	}
}

func TestClockIsMidnight(t *testing.T) {
	cases := []struct {
		in   Clock
		want bool
	}{
		{New(0, 0, 0, 0), true},
		{Day, true},
		{24 * Hour, true},
		{New(24, 0, 0, 0), true},
		{New(-24, 0, 0, 0), true},
		{New(-48, 0, 0, 0), true},
		{New(48, 0, 0, 0), true},
		{New(0, 0, 0, 1), false},
		{New(2, 0, 0, 1), false},
		{New(-1, 0, 0, 0), false},
		{New(0, 0, 0, -1), false},
	}
	for i, x := range cases {
		got := x.in.IsMidnight()
		if got != x.want {
			t.Errorf("%d: %v got %v, want %v, %d", i, x.in, x.in.IsMidnight(), x.want, x.in.Mod24())
		}
	}
}

func TestClockMod(t *testing.T) {
	cases := []struct {
		h, want Clock
	}{
		{0, 0},
		{1 * Hour, 1 * Hour},
		{2 * Hour, 2 * Hour},
		{23 * Hour, 23 * Hour},
		{24 * Hour, 0},
		{-24 * Hour, 0},
		{-48 * Hour, 0},
		{25 * Hour, Hour},
		{49 * Hour, Hour},
		{-1 * Hour, 23 * Hour},
		{-23 * Hour, Hour},
		{New(0, 0, 0, 1), 1},
		{New(0, 0, 1, 0), Second},
		{New(0, 0, 0, -1), New(23, 59, 59, 999)},
	}
	for i, x := range cases {
		clock := x.h
		got := clock.Mod24()
		if got != x.want {
			t.Errorf("%d: %dh: got %#v, want %#v", i, x.h, got, x.want)
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
	for i, x := range cases {
		clock := Clock(x.h) * Hour
		if clock.Days() != x.days {
			t.Errorf("%d: %dh: got %v, want %v", i, x.h, clock.Days(), x.days)
		}
	}
}

func TestClockString(t *testing.T) {
	cases := []struct {
		h, m, s, ms                                Clock
		hh, hhmm, hhmmss, str, h12, hmm12, hmmss12 string
	}{
		{0, 0, 0, 0, "00", "00:00", "00:00:00", "00:00:00.000", "12am", "12:00am", "12:00:00am"},
		{0, 0, 0, 1, "00", "00:00", "00:00:00", "00:00:00.001", "12am", "12:00am", "12:00:00am"},
		{0, 0, 1, 0, "00", "00:00", "00:00:01", "00:00:01.000", "12am", "12:00am", "12:00:01am"},
		{0, 1, 0, 0, "00", "00:01", "00:01:00", "00:01:00.000", "12am", "12:01am", "12:01:00am"},
		{1, 0, 0, 0, "01", "01:00", "01:00:00", "01:00:00.000", "1am", "1:00am", "1:00:00am"},
		{1, 2, 3, 4, "01", "01:02", "01:02:03", "01:02:03.004", "1am", "1:02am", "1:02:03am"},
		{11, 0, 0, 0, "11", "11:00", "11:00:00", "11:00:00.000", "11am", "11:00am", "11:00:00am"},
		{12, 0, 0, 0, "12", "12:00", "12:00:00", "12:00:00.000", "12pm", "12:00pm", "12:00:00pm"},
		{13, 0, 0, 0, "13", "13:00", "13:00:00", "13:00:00.000", "1pm", "1:00pm", "1:00:00pm"},
		{24, 0, 0, 0, "24", "24:00", "24:00:00", "24:00:00.000", "12am", "12:00am", "12:00:00am"},
		{24, 0, 0, 1, "00", "00:00", "00:00:00", "00:00:00.001", "12am", "12:00am", "12:00:00am"},
		{-1, 0, 0, 0, "23", "23:00", "23:00:00", "23:00:00.000", "11pm", "11:00pm", "11:00:00pm"},
		{-1, -1, -1, -1, "22", "22:58", "22:58:58", "22:58:58.999", "10pm", "10:58pm", "10:58:58pm"},
	}
	for _, x := range cases {
		d := Clock(x.h*Hour + x.m*Minute + x.s*Second + x.ms)
		if d.Hh() != x.hh {
			t.Errorf("%d, %d, %d, %d, got %v, want %v (%d)", x.h, x.m, x.s, x.ms, d.Hh(), x.hh, d)
		}
		if d.HhMm() != x.hhmm {
			t.Errorf("%d, %d, %d, %d, got %v, want %v (%d)", x.h, x.m, x.s, x.ms, d.HhMm(), x.hhmm, d)
		}
		if d.HhMmSs() != x.hhmmss {
			t.Errorf("%d, %d, %d, %d, got %v, want %v (%d)", x.h, x.m, x.s, x.ms, d.HhMmSs(), x.hhmmss, d)
		}
		if d.String() != x.str {
			t.Errorf("%d, %d, %d, %d, got %v, want %v (%d)", x.h, x.m, x.s, x.ms, d.String(), x.str, d)
		}
		if d.Hh12() != x.h12 {
			t.Errorf("%d, %d, %d, %d, got %v, want %v (%d)", x.h, x.m, x.s, x.ms, d.Hh12(), x.h12, d)
		}
		if d.HhMm12() != x.hmm12 {
			t.Errorf("%d, %d, %d, %d, got %v, want %v (%d)", x.h, x.m, x.s, x.ms, d.HhMm12(), x.hmm12, d)
		}
		if d.HhMmSs12() != x.hmmss12 {
			t.Errorf("%d, %d, %d, %d, got %v, want %v (%d)", x.h, x.m, x.s, x.ms, d.HhMmSs12(), x.hmmss12, d)
		}
	}
}

func TestClockParseGoods(t *testing.T) {
	cases := []struct {
		str  string
		want Clock
	}{
		{"00", New(0, 0, 0, 0)},
		{"01", New(1, 0, 0, 0)},
		{"23", New(23, 0, 0, 0)},
		{"00:00", New(0, 0, 0, 0)},
		{"00:01", New(0, 1, 0, 0)},
		{"01:00", New(1, 0, 0, 0)},
		{"01:02", New(1, 2, 0, 0)},
		{"23:59", New(23, 59, 0, 0)},
		{"0911", New(9, 11, 0, 0)},
		{"1024", New(10, 24, 0, 0)},
		{"2359", New(23, 59, 0, 0)},
		{"00:00:00", New(0, 0, 0, 0)},
		{"00:00:01", New(0, 0, 1, 0)},
		{"00:01:00", New(0, 1, 0, 0)},
		{"01:00:00", New(1, 0, 0, 0)},
		{"01:02:03", New(1, 2, 3, 0)},
		{"23:59:59", New(23, 59, 59, 0)},
		{"235959", New(23, 59, 59, 0)},
		{"00:00:00.000", New(0, 0, 0, 0)},
		{"00:00:00.001", New(0, 0, 0, 1)},
		{"00:00:01.000", New(0, 0, 1, 0)},
		{"00:01:00.000", New(0, 1, 0, 0)},
		{"01:00:00.000", New(1, 0, 0, 0)},
		{"01:02:03.004", New(1, 2, 3, 4)},
		{"01:02:03.04", New(1, 2, 3, 40)},
		{"01:02:03.4", New(1, 2, 3, 400)},
		{"23:59:59.999", New(23, 59, 59, 999)},
		{"0am", New(0, 0, 0, 0)},
		{"00am", New(0, 0, 0, 0)},
		{"12am", New(0, 0, 0, 0)},
		{"12pm", New(12, 0, 0, 0)},
		{"12:01am", New(0, 1, 0, 0)},
		{"12:01pm", New(12, 1, 0, 0)},
		{"12:01:02am", New(0, 1, 2, 0)},
		{"12:01:02pm", New(12, 1, 2, 0)},
		{"1am", New(1, 0, 0, 0)},
		{"1pm", New(13, 0, 0, 0)},
		{"1:00am", New(1, 0, 0, 0)},
		{"1:23am", New(1, 23, 0, 0)},
		{"1:23pm", New(13, 23, 0, 0)},
		{"01:23pm", New(13, 23, 0, 0)},
		{"1:00:00am", New(1, 0, 0, 0)},
		{"1:02:03pm", New(13, 2, 3, 0)},
		{"01:02:03pm", New(13, 2, 3, 0)},
		{"1:02:03.004pm", New(13, 2, 3, 4)},
		{"01:02:03.004pm", New(13, 2, 3, 4)},
		{"1:20:30.04pm", New(13, 20, 30, 40)},
		{"1:20:30.4pm", New(13, 20, 30, 400)},
		{"1:20:30.pm", New(13, 20, 30, 0)},
	}
	for _, x := range cases {
		str := MustParse(x.str)
		if str != x.want {
			t.Errorf("%s, got %v, want %v", x.str, str, x.want)
		}
	}
}

func TestClockParseBads(t *testing.T) {
	cases := []struct {
		str string
	}{
		{"0"},
		{"0:01"},
		{"0:00:01"},
		{"hh"},
		{"00-00"},
		{"00:00-00"},
		{"00:00:00-"},
		{"00:00:00-0"},
		{"00:00:00-00"},
		{"00:00:00-000"},
		{"00:mm"},
		{"00:00:ss"},
		{"00:00:00.xxx"},
		{"01-02:03.004"},
		{"01:02-03.04"},
		{"01:02:03-4"},
		{"12xm"},
		{"12-01am"},
		{"12:01-02am"},
		{"ham"},
		{"hham"},
		{"1xm"},
		{"1-00am"},
		{"1:00-00am"},
		{"1:02:03-4pm"},
		{"1:02:03-04pm"},
		{"1:02:03-004pm"},
		{"1:02:03.0045pm"},
	}
	for _, x := range cases {
		c, err := Parse(x.str)
		if err == nil {
			t.Errorf("%s, got %#v, want err", x.str, c)
			//		} else {
			//			println(err.Error())
		}
	}
}

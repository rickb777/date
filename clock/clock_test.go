// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package clock

import (
	"fmt"
	"testing"
	time "time"

	"github.com/govalues/decimal"
	"github.com/rickb777/period"
)

func TestClockHoursMinutesSeconds(t *testing.T) {
	cases := []struct {
		in              Clock
		h, m, s, ms, ns int
	}{
		{in: New(0, 0, 0, 0), h: 0, m: 0, s: 0, ms: 0, ns: 0},
		{in: New(1, 2, 3, 4), h: 1, m: 2, s: 3, ms: 4, ns: 4_000_000},
		{in: New(23, 59, 59, 999), h: 23, m: 59, s: 59, ms: 999, ns: 999_000_000},
		{in: New(0, 0, 0, -1), h: 23, m: 59, s: 59, ms: 999, ns: 999_000_000},
		{in: NewAt(time.Date(2015, 12, 4, 18, 50, 42, 101202303, time.UTC)), h: 18, m: 50, s: 42, ms: 101, ns: 101202303},
	}
	for i, x := range cases {
		t.Run(fmt.Sprintf("%d %s", i, x.in), func(t *testing.T) {
			h1 := x.in.Hour()
			m1 := x.in.Minute()
			s1 := x.in.Second()
			if h1 != x.h || m1 != x.m || s1 != x.s {
				t.Errorf("%d: got %02d:%02d:%02d, want %s (%d)", i, h1, m1, s1, x.in, x.in)
			}

			ms := x.in.Millisecond()
			ns := x.in.Nanosecond()
			if ms != x.ms || ns != x.ns {
				t.Errorf("%d: got %02d:%02d:%02d.%03d (%d), want %s (%d)", i, h1, m1, s1, ms, ns, x.in, x.in)
			}

			h2, m2, s2 := x.in.HourMinuteSecond()
			if h2 != x.h || m2 != x.m || s2 != x.s {
				t.Errorf("%d: got %02d:%02d:%02d, want %s (%d)", i, h2, m2, s2, x.in, x.in)
			}
		})
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
		t.Run(fmt.Sprintf("%d %s", i, x.in), func(t *testing.T) {
			d := x.in.DurationSinceMidnight()
			c2 := SinceMidnight(d)
			if c2 != x.in {
				t.Errorf("%d: got %v, want %v (%d)", i, c2, x.in, x.in)
			}
		})
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
	for i, x := range cases {
		t.Run(fmt.Sprintf("%d %s", i, x.in), func(t *testing.T) {
			got := x.in.IsInOneDay()
			if got != x.want {
				t.Errorf("%v got %v, want %v", x.in, x.in.IsInOneDay(), x.want)
			}
		})
	}
}

func TestClockAdd(t *testing.T) {
	in := 2 * Hour
	cases := []struct {
		h, m, s, ms int
		want        Clock
	}{
		{h: 0, m: 0, s: 0, ms: 0, want: New(2, 0, 0, 0)},
		{h: 0, m: 0, s: 0, ms: 1, want: New(2, 0, 0, 1)},
		{h: 0, m: 0, s: 0, ms: -1, want: New(1, 59, 59, 999)},
		{h: 0, m: 0, s: 1, ms: 0, want: New(2, 0, 1, 0)},
		{h: 0, m: 0, s: -1, ms: 0, want: New(1, 59, 59, 0)},
		{h: 0, m: 1, s: 0, ms: 0, want: New(2, 1, 0, 0)},
		{h: 0, m: -1, s: 0, ms: 0, want: New(1, 59, 0, 0)},
		{h: 1, m: 0, s: 0, ms: 0, want: New(3, 0, 0, 0)},
		{h: -1, m: 0, s: 0, ms: 0, want: New(1, 0, 0, 0)},
		{h: -2, m: 0, s: 0, ms: 0, want: New(0, 0, 0, 0)},
		{h: -2, m: 0, s: -1, ms: -1, want: New(0, 0, -1, -1)},
	}
	for i, x := range cases {
		t.Run(fmt.Sprintf("%d %s", i, x.want), func(t *testing.T) {
			got := in.Add(x.h, x.m, x.s, x.ms)
			if got != x.want {
				t.Errorf("%d: %d %d %d.%d: got %v, want %v", i, x.h, x.m, x.s, x.ms, got, x.want)
			}
		})
	}
}

func TestClockAddDuration(t *testing.T) {
	in := 2 * Hour
	cases := []struct {
		d    time.Duration
		want Clock
		ns   int
	}{
		{d: 0, want: 2 * Hour},
		{d: 1, want: 2*Hour + 1, ns: 1},
		{d: time.Millisecond, want: New(2, 0, 0, 1), ns: 1_000_000},
		{d: -time.Second, want: New(1, 59, 59, 0)},
		{d: 7 * time.Minute, want: New(2, 7, 0, 0)},
	}
	for i, x := range cases {
		t.Run(fmt.Sprintf("%d %s", i, x.want), func(t *testing.T) {
			got := in.AddDuration(x.d)
			if got != x.want {
				t.Errorf("%d: %d: got %v, want %v", i, x.d, got, x.want)
			}
			ns := got.Nanosecond()
			if ns != x.ns {
				t.Errorf("%d: got %dns, want %dns", i, ns, x.ns)
			}
		})
	}
}

func TestClockAddPeriod(t *testing.T) {
	cases := []struct {
		p        period.Period
		in, want Clock
	}{
		{period.Zero, 2 * Hour, New(2, 0, 0, 0)},
		{period.NewHMS(0, 0, 1), 2 * Hour, New(2, 0, 1, 0)},
		{period.NewHMS(0, 0, -1), 2 * Hour, New(1, 59, 59, 0)},
		{period.MustNewDecimal(decimal.Zero, decimal.Zero, decimal.Zero, decimal.Zero, decimal.Zero, decimal.Zero, decimal.MustNew(1, 3)), 2 * Hour, New(2, 0, 0, 1)},
		{period.NewHMS(0, 7, 0), 2 * Hour, New(2, 7, 0, 0)},
	}
	for i, x := range cases {
		t.Run(fmt.Sprintf("%d %s", i, x.want), func(t *testing.T) {
			got, ok := x.in.AddPeriod(x.p)
			if !ok {
				t.Errorf("%d: %s: not ok", i, x.p)
			}
			if got != x.want {
				t.Errorf("%d: %s: got %v, want %v", i, x.p, got, x.want)
			}
		})
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
		t.Run(fmt.Sprintf("%d %s", i, x.want), func(t *testing.T) {
			got := x.c1.ModSubtract(x.c2)
			if got != x.want {
				t.Errorf("%d: %v - %v: got %v, want %v", i, x.c1, x.c2, got, x.want)
			}
		})
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
		t.Run(fmt.Sprintf("%d %s", i, x.in), func(t *testing.T) {
			got := x.in.IsMidnight()
			if got != x.want {
				t.Errorf("%d: %v got %v, want %v, %d", i, x.in, x.in.IsMidnight(), x.want, x.in.Mod24())
			}
		})
	}
}

func TestClockMod24(t *testing.T) {
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
		{New(0, 0, 0, 1), Millisecond},
		{New(0, 0, 1, 0), Second},
		{New(0, 0, 0, -1), New(23, 59, 59, 999)},
	}
	for i, x := range cases {
		t.Run(fmt.Sprintf("%d %s", i, x.want), func(t *testing.T) {
			clock := x.h
			got := clock.Mod24()
			if got != x.want {
				t.Errorf("%d: %dh: got %#v, want %#v", i, x.h, got, x.want)
			}
		})
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
		t.Run(fmt.Sprintf("%d %d %x", i, x.h, x.days), func(t *testing.T) {
			clock := Clock(x.h) * Hour
			if clock.Days() != x.days {
				t.Errorf("%d: %dh: got %v, want %v", i, x.h, clock.Days(), x.days)
			}
		})
	}
}

func TestClockFormat(t *testing.T) {
	cases := []struct {
		h, m, s, ms, ns                       Clock
		hh, hhmm, hhmmss, h12, hmm12, hmmss12 string
	}{
		{0, 0, 0, 0, 0, "00", "00:00", "00:00:00", "12am", "12:00am", "12:00:00am"},
		{0, 0, 0, 1, 0, "00", "00:00", "00:00:00", "12am", "12:00am", "12:00:00am"},
		{0, 0, 0, 0, 1, "00", "00:00", "00:00:00", "12am", "12:00am", "12:00:00am"},
		{0, 0, 1, 0, 0, "00", "00:00", "00:00:01", "12am", "12:00am", "12:00:01am"},
		{0, 1, 0, 0, 0, "00", "00:01", "00:01:00", "12am", "12:01am", "12:01:00am"},
		{1, 0, 0, 0, 0, "01", "01:00", "01:00:00", "1am", "1:00am", "1:00:00am"},
		{1, 2, 3, 4, 5, "01", "01:02", "01:02:03", "1am", "1:02am", "1:02:03am"},
		{11, 0, 0, 0, 0, "11", "11:00", "11:00:00", "11am", "11:00am", "11:00:00am"},
		{12, 0, 0, 0, 0, "12", "12:00", "12:00:00", "12pm", "12:00pm", "12:00:00pm"},
		{13, 0, 0, 0, 0, "13", "13:00", "13:00:00", "1pm", "1:00pm", "1:00:00pm"},
		{24, 0, 0, 0, 0, "24", "24:00", "24:00:00", "12am", "12:00am", "12:00:00am"},
		{24, 0, 0, 1, 0, "00", "00:00", "00:00:00", "12am", "12:00am", "12:00:00am"},
		{-1, 0, 0, 0, 0, "23", "23:00", "23:00:00", "11pm", "11:00pm", "11:00:00pm"},
		{-1, -1, -1, -1, 0, "22", "22:58", "22:58:58", "10pm", "10:58pm", "10:58:58pm"},
	}
	for i, x := range cases {
		t.Run(fmt.Sprintf("%d %s", i, x.hhmmss), func(t *testing.T) {
			d := x.h*Hour + x.m*Minute + x.s*Second + x.ms*Millisecond + x.ns
			if d.Hh() != x.hh {
				t.Errorf("%d, %d, %d, %d, got %v, want %v (%d)", x.h, x.m, x.s, x.ns, d.Hh(), x.hh, d)
			}
			if d.HhMm() != x.hhmm {
				t.Errorf("%d, %d, %d, %d, got %v, want %v (%d)", x.h, x.m, x.s, x.ns, d.HhMm(), x.hhmm, d)
			}
			if d.HhMmSs() != x.hhmmss {
				t.Errorf("%d, %d, %d, %d, got %v, want %v (%d)", x.h, x.m, x.s, x.ns, d.HhMmSs(), x.hhmmss, d)
			}
			if d.Hh12() != x.h12 {
				t.Errorf("%d, %d, %d, %d, got %v, want %v (%d)", x.h, x.m, x.s, x.ns, d.Hh12(), x.h12, d)
			}
			if d.HhMm12() != x.hmm12 {
				t.Errorf("%d, %d, %d, %d, got %v, want %v (%d)", x.h, x.m, x.s, x.ns, d.HhMm12(), x.hmm12, d)
			}
			if d.HhMmSs12() != x.hmmss12 {
				t.Errorf("%d, %d, %d, %d, got %v, want %v (%d)", x.h, x.m, x.s, x.ns, d.HhMmSs12(), x.hmmss12, d)
			}
		})
	}
}

func TestClockString(t *testing.T) {
	cases := []struct {
		h, m, s, ms, ns Clock
		trunc, str      string
	}{
		{0, 0, 0, 0, 0, "00:00:00.000", "00:00:00.000"},
		{0, 0, 0, 1, 0, "00:00:00.001", "00:00:00.001"},
		{0, 0, 0, 0, 1, "00:00:00.000", "00:00:00.000000001"},
		{0, 0, 1, 0, 0, "00:00:01.000", "00:00:01.000"},
		{0, 1, 0, 0, 0, "00:01:00.000", "00:01:00.000"},
		{1, 0, 0, 0, 0, "01:00:00.000", "01:00:00.000"},
		{1, 2, 3, 4, 5, "01:02:03.004", "01:02:03.004000005"},
		{11, 0, 0, 0, 0, "11:00:00.000", "11:00:00.000"},
		{12, 0, 0, 0, 0, "12:00:00.000", "12:00:00.000"},
		{13, 0, 0, 0, 0, "13:00:00.000", "13:00:00.000"},
		{24, 0, 0, 0, 0, "24:00:00.000", "24:00:00.000"},
		{24, 0, 0, 1, 0, "00:00:00.001", "00:00:00.001"},
		{-1, 0, 0, 0, 0, "23:00:00.000", "23:00:00.000"},
		{-1, -1, -1, -1, 0, "22:58:58.999", "22:58:58.999"},
	}
	for i, x := range cases {
		t.Run(fmt.Sprintf("%d %s", i, x.str), func(t *testing.T) {
			d := x.h*Hour + x.m*Minute + x.s*Second + x.ms*Millisecond + x.ns
			tr := d.TruncateMillisecond()
			if tr.String() != x.trunc {
				t.Errorf("%d, %d, %d, %d, got %v, want %v (%d)", x.h, x.m, x.s, x.ns, tr.String(), x.str, d)
			}
			if d.String() != x.str {
				t.Errorf("%d, %d, %d, %d, got %v, want %v (%d)", x.h, x.m, x.s, x.ns, d.String(), x.str, d)
			}
			if ValueAsString(d) != x.str {
				t.Errorf("%d, %d, %d, %d, got %v, want %v (%d)", x.h, x.m, x.s, x.ns, d.String(), x.str, d)
			}
		})
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
	for i, x := range cases {
		t.Run(fmt.Sprintf("%d %s", i, x.str), func(t *testing.T) {
			str := MustParse(x.str)
			if str != x.want {
				t.Errorf("%s, got %v, want %v", x.str, str, x.want)
			}
		})
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
	for i, x := range cases {
		t.Run(fmt.Sprintf("%d %s", i, x), func(t *testing.T) {
			c, err := Parse(x.str)
			if err == nil {
				t.Errorf("%s, got %#v, want err", x.str, c)
				//		} else {
				//			println(err.Error())
			}
		})
	}
}

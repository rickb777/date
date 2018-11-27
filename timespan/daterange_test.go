// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package timespan

import (
	"fmt"
	"strings"
	"testing"
	"time"

	. "github.com/rickb777/date"
	"github.com/rickb777/date/period"
)

var d0320 = New(2015, time.March, 20)
var d0321 = New(2015, time.March, 21)
var d0325 = New(2015, time.March, 25)
var d0326 = New(2015, time.March, 26)
var d0327 = New(2015, time.March, 27)
var d0328 = New(2015, time.March, 28)
var d0329 = New(2015, time.March, 29) // n.b. clocks go forward (UK)
var d0330 = New(2015, time.March, 30)
var d0331 = New(2015, time.March, 31)
var d0401 = New(2015, time.April, 1)
var d0402 = New(2015, time.April, 2)
var d0403 = New(2015, time.April, 3)
var d0404 = New(2015, time.April, 4)
var d0407 = New(2015, time.April, 7)
var d0408 = New(2015, time.April, 8)
var d0409 = New(2015, time.April, 9)
var d0410 = New(2015, time.April, 10)
var d0501 = New(2015, time.May, 1)
var d1025 = New(2015, time.October, 25)

var london *time.Location = mustLoadLocation("Europe/London")

func mustLoadLocation(name string) *time.Location {
	loc, err := time.LoadLocation("Europe/London")
	if err != nil {
		panic(err)
	}
	return loc
}

func TestNewDateRangeOf(t *testing.T) {
	dr := NewDateRangeOf(t0327, 7*24*time.Hour)
	isEq(t, 0, dr.mark, d0327)
	isEq(t, 0, dr.Days(), PeriodOfDays(7))
	isEq(t, 0, dr.IsEmpty(), false)
	isEq(t, 0, dr.Start(), d0327)
	isEq(t, 0, dr.Last(), d0402)
	isEq(t, 0, dr.End(), d0403)

	dr2 := NewDateRangeOf(t0327, -7*24*time.Hour)
	isEq(t, 0, dr2.mark, d0327)
	isEq(t, 0, dr2.Days(), PeriodOfDays(7))
	isEq(t, 0, dr2.IsEmpty(), false)
	isEq(t, 0, dr2.Start(), d0321)
	isEq(t, 0, dr2.Last(), d0327)
	isEq(t, 0, dr2.End(), d0328)
}

func TestNewDateRangeWithNormalise(t *testing.T) {
	r1 := NewDateRange(d0327, d0402)
	isEq(t, 0, r1.Start(), d0327)
	isEq(t, 0, r1.Last(), d0401)
	isEq(t, 0, r1.End(), d0402)

	r2 := NewDateRange(d0402, d0327)
	isEq(t, 0, r2.Start(), d0327)
	isEq(t, 0, r2.Last(), d0401)
	isEq(t, 0, r2.End(), d0402)
}

func TestEmptyRange(t *testing.T) {
	drN0 := DateRange{d0327, -1}
	isEq(t, 0, drN0.Days(), PeriodOfDays(1))
	isEq(t, 0, drN0.IsZero(), false)
	isEq(t, 0, drN0.IsEmpty(), false)
	isEq(t, 0, drN0.Start(), d0327)
	isEq(t, 0, drN0.Last(), d0327)
	isEq(t, 0, drN0.String(), "1 day on 2015-03-26")

	dr0 := DateRange{}
	isEq(t, 0, dr0.Days(), PeriodOfDays(0))
	isEq(t, 0, dr0.IsZero(), true)
	isEq(t, 0, dr0.IsEmpty(), true)
	isEq(t, 0, dr0.String(), "0 days at 1970-01-01")

	dr1 := EmptyRange(Date{})
	isEq(t, 0, dr1.IsZero(), true)
	isEq(t, 0, dr1.IsEmpty(), true)
	isEq(t, 0, dr1.Days(), PeriodOfDays(0))

	dr2 := EmptyRange(d0327)
	isEq(t, 0, dr2.IsZero(), false)
	isEq(t, 0, dr2.IsEmpty(), true)
	isEq(t, 0, dr2.Start(), d0327)
	isEq(t, 0, dr2.Last().IsZero(), true)
	isEq(t, 0, dr2.End(), d0327)
	isEq(t, 0, dr2.Days(), PeriodOfDays(0))
	isEq(t, 0, dr2.String(), "0 days at 2015-03-27")
}

func TestOneDayRange(t *testing.T) {
	dr1 := OneDayRange(Date{})
	isEq(t, 0, dr1.IsZero(), false)
	isEq(t, 0, dr1.IsEmpty(), false)
	isEq(t, 0, dr1.Days(), PeriodOfDays(1))

	dr2 := OneDayRange(d0327)
	isEq(t, 0, dr2.Start(), d0327)
	isEq(t, 0, dr2.Last(), d0327)
	isEq(t, 0, dr2.End(), d0328)
	isEq(t, 0, dr2.Days(), PeriodOfDays(1))
	isEq(t, 0, dr2.String(), "1 day on 2015-03-27")
}

func TestDayRange(t *testing.T) {
	dr1 := DayRange(Date{}, 0)
	isEq(t, 0, dr1.IsZero(), true)
	isEq(t, 0, dr1.IsEmpty(), true)
	isEq(t, 0, dr1.Days(), PeriodOfDays(0))

	dr2 := DayRange(d0327, 2)
	isEq(t, 0, dr2.Start(), d0327)
	isEq(t, 0, dr2.Last(), d0328)
	isEq(t, 0, dr2.End(), d0329)
	isEq(t, 0, dr2.Days(), PeriodOfDays(2))
	isEq(t, 0, dr2.String(), "2 days from 2015-03-27 to 2015-03-28")

	dr3 := DayRange(d0327, -2)
	isEq(t, 0, dr3.Start(), d0325)
	isEq(t, 0, dr3.Last(), d0326)
	isEq(t, 0, dr3.End(), d0327)
	isEq(t, 0, dr3.Days(), PeriodOfDays(2))
	isEq(t, 0, dr3.String(), "2 days from 2015-03-25 to 2015-03-26")
}

func TestNewYearOf(t *testing.T) {
	dr := NewYearOf(2015)
	isEq(t, 0, dr.Days(), PeriodOfDays(365))
	isEq(t, 0, dr.Start(), New(2015, time.January, 1))
	isEq(t, 0, dr.Last(), New(2015, time.December, 31))
	isEq(t, 0, dr.End(), New(2016, time.January, 1))
}

func TestNewMonthOf(t *testing.T) {
	dr := NewMonthOf(2015, time.February)
	isEq(t, 0, dr.Days(), PeriodOfDays(28))
	isEq(t, 0, dr.Start(), New(2015, time.February, 1))
	isEq(t, 0, dr.Last(), New(2015, time.February, 28))
	isEq(t, 0, dr.End(), New(2015, time.March, 1))
}

func TestShiftAndExtend(t *testing.T) {
	cases := []struct {
		dr    DateRange
		n     PeriodOfDays
		start Date
		end   Date
		s     string
	}{
		{DayRange(d0327, 6).ShiftBy(0), 6, d0327, d0402, "6 days from 2015-03-27 to 2015-04-01"},
		{DayRange(d0327, 6).ShiftBy(7), 6, d0403, d0409, "6 days from 2015-04-03 to 2015-04-08"},
		{DayRange(d0327, 6).ShiftBy(-1), 6, d0326, d0401, "6 days from 2015-03-26 to 2015-03-31"},
		{DayRange(d0327, 6).ShiftBy(-7), 6, d0320, d0326, "6 days from 2015-03-20 to 2015-03-25"},
		{NewDateRange(d0327, d0402).ShiftBy(-7), 6, d0320, d0326, "6 days from 2015-03-20 to 2015-03-25"},

		{EmptyRange(d0327).ExtendBy(0), 0, d0327, d0327, "0 days at 2015-03-27"},
		{EmptyRange(d0327).ExtendBy(6), 6, d0327, d0402, "6 days from 2015-03-27 to 2015-04-01"},
		{DayRange(d0327, 6).ExtendBy(0), 6, d0327, d0402, "6 days from 2015-03-27 to 2015-04-01"},
		{DayRange(d0327, 6).ExtendBy(7), 13, d0327, d0409, "13 days from 2015-03-27 to 2015-04-08"},
		{DayRange(d0327, 6).ExtendBy(-6), 0, d0327, d0327, "0 days at 2015-03-27"},
		{DayRange(d0327, 6).ExtendBy(-8), 2, d0325, d0327, "2 days from 2015-03-25 to 2015-03-26"},

		{DayRange(d0327, 6).ShiftByPeriod(period.NewYMD(0, 0, 0)), 6, d0327, d0402, "6 days from 2015-03-27 to 2015-04-01"},
		{DayRange(d0327, 6).ShiftByPeriod(period.NewYMD(0, 0, 7)), 6, d0403, d0409, "6 days from 2015-04-03 to 2015-04-08"},
		{DayRange(d0327, 6).ShiftByPeriod(period.NewYMD(0, 0, -7)), 6, d0320, d0326, "6 days from 2015-03-20 to 2015-03-25"},

		{DayRange(d0327, 6).ExtendByPeriod(period.NewYMD(0, 0, 0)), 6, d0327, d0402, "6 days from 2015-03-27 to 2015-04-01"},
		{DayRange(d0327, 6).ExtendByPeriod(period.NewYMD(0, 0, 7)), 13, d0327, d0409, "13 days from 2015-03-27 to 2015-04-08"},
		{DayRange(d0327, 6).ExtendByPeriod(period.NewYMD(0, 0, -5)), 1, d0327, d0328, "1 day on 2015-03-27"},
		{DayRange(d0327, 6).ExtendByPeriod(period.NewYMD(0, 0, -6)), 0, d0327, d0327, "0 days at 2015-03-27"},
		{DayRange(d0327, 6).ExtendByPeriod(period.NewYMD(0, 0, -7)), 1, d0326, d0327, "1 day on 2015-03-26"},
	}

	for i, c := range cases {
		isEq(t, i, c.dr.Days(), c.n)
		isEq(t, i, c.dr.Start(), c.start)
		isEq(t, i, c.dr.End(), c.end)
		isEq(t, i, c.dr.String(), c.s)
	}
}

func TestContains0(t *testing.T) {
	old := time.Local
	time.Local = time.FixedZone("Test", 7200)
	dr := EmptyRange(d0326)
	isEq(t, 0, dr.Contains(d0320), false, dr, d0320)
	time.Local = old
}

func TestContains1(t *testing.T) {
	old := time.Local
	time.Local = time.FixedZone("Test", 7200)
	dr := DayRange(d0326, 2)
	isEq(t, 0, dr.Contains(d0320), false, dr, d0320)
	isEq(t, 0, dr.Contains(d0325), false, dr, d0325)
	isEq(t, 0, dr.Contains(d0326), true, dr, d0326)
	isEq(t, 0, dr.Contains(d0327), true, dr, d0327)
	isEq(t, 0, dr.Contains(d0328), false, dr, d0328)
	isEq(t, 0, dr.Contains(d0401), false, dr, d0401)
	isEq(t, 0, dr.Contains(d0410), false, dr, d0410)
	isEq(t, 0, dr.Contains(d0501), false, dr, d0501)
	time.Local = old
}

func TestContains2(t *testing.T) {
	old := time.Local
	time.Local = time.FixedZone("Test", 7200)
	dr := OneDayRange(d0326)
	isEq(t, 0, dr.Contains(d0325), false, dr, d0325)
	isEq(t, 0, dr.Contains(d0326), true, dr, d0326)
	isEq(t, 0, dr.Contains(d0327), false, dr, d0327)
	time.Local = old
}

func TestContainsTime0(t *testing.T) {
	old := time.Local
	time.Local = time.FixedZone("Test", 7200)
	t0328e := time.Date(2015, 3, 28, 23, 59, 59, 999999999, time.UTC)
	t0329 := time.Date(2015, 3, 29, 0, 0, 0, 0, time.UTC)

	dr := EmptyRange(d0327)
	isEq(t, 0, dr.StartUTC(), t0327, dr, t0327)
	isEq(t, 0, dr.EndUTC(), t0327, dr, t0327)
	isEq(t, 0, dr.ContainsTime(t0327), false, dr, t0327)
	isEq(t, 0, dr.ContainsTime(t0328), false, dr, t0328)
	isEq(t, 0, dr.ContainsTime(t0328e), false, dr, t0328e)
	isEq(t, 0, dr.ContainsTime(t0329), false, dr, t0329)
	time.Local = old
}

func TestContainsTimeUTC(t *testing.T) {
	old := time.Local
	time.Local = time.FixedZone("Test", 7200)
	t0328e := time.Date(2015, 3, 28, 23, 59, 59, 999999999, time.UTC)
	t0329 := time.Date(2015, 3, 29, 0, 0, 0, 0, time.UTC)

	dr := DayRange(d0327, 2)
	isEq(t, 0, dr.StartUTC(), t0327, dr, t0327)
	isEq(t, 0, dr.EndUTC(), t0329, dr, t0329)
	isEq(t, 0, dr.ContainsTime(t0327), true, dr, t0327)
	isEq(t, 0, dr.ContainsTime(t0328), true, dr, t0328)
	isEq(t, 0, dr.ContainsTime(t0328e), true, dr, t0328e)
	isEq(t, 0, dr.ContainsTime(t0329), false, dr, t0329)
	time.Local = old
}

func TestMerge1(t *testing.T) {
	dr1 := DayRange(d0327, 2)
	dr2 := DayRange(d0327, 8)
	m1 := dr1.Merge(dr2)
	m2 := dr2.Merge(dr1)
	isEq(t, 0, m1.Start(), d0327)
	isEq(t, 0, m1.End(), d0404)
	isEq(t, 0, m1, m2)
}

func TestMerge2(t *testing.T) {
	dr1 := DayRange(d0328, 2)
	dr2 := DayRange(d0327, 8)
	m1 := dr1.Merge(dr2)
	m2 := dr2.Merge(dr1)
	isEq(t, 0, m1.Start(), d0327)
	isEq(t, 0, m1.End(), d0404)
	isEq(t, 0, m1, m2)
}

func TestMergeOverlapping(t *testing.T) {
	dr1 := OneDayRange(d0320).ExtendBy(12)
	dr2 := OneDayRange(d0401).ExtendBy(6)
	m1 := dr1.Merge(dr2)
	m2 := dr2.Merge(dr1)
	isEq(t, 0, m1.Start(), d0320)
	isEq(t, 0, m1.End(), d0408)
	isEq(t, 0, m1, m2)
}

func TestMergeNonOverlapping(t *testing.T) {
	dr1 := OneDayRange(d0320).ExtendBy(2)
	dr2 := OneDayRange(d0401).ExtendBy(6)
	m1 := dr1.Merge(dr2)
	m2 := dr2.Merge(dr1)
	isEq(t, 0, m1.Start(), d0320)
	isEq(t, 0, m1.End(), d0408)
	isEq(t, 0, m1, m2)
}

func TestMergeEmpties(t *testing.T) {
	dr1 := EmptyRange(d0320)
	dr2 := EmptyRange(d0408) // curiously, this is *not* included because it has no size.
	m1 := dr1.Merge(dr2)
	m2 := dr2.Merge(dr1)
	isEq(t, 0, m1.Start(), d0320)
	isEq(t, 0, m1.End(), d0408)
	isEq(t, 0, m1, m2)
}

func TestMergeZeroes(t *testing.T) {
	dr0 := DateRange{}
	dr1 := OneDayRange(d0401).ExtendBy(6)
	m1 := dr1.Merge(dr0)
	m2 := dr0.Merge(dr1)
	m3 := dr0.Merge(dr0)
	isEq(t, 0, m1.Start(), d0401)
	isEq(t, 0, m1.End(), d0408)
	isEq(t, 0, m1, m2)
	isEq(t, 0, m3.IsZero(), true)
	isEq(t, 0, m3, dr0)
}

func TestDurationNormalUTC(t *testing.T) {
	dr := OneDayRange(d0329)
	isEq(t, 0, dr.Duration(), time.Hour*24)
}

func TestDurationInZoneWithDaylightSaving(t *testing.T) {
	isEq(t, 0, OneDayRange(d0328).DurationIn(london), time.Hour*24)
	isEq(t, 0, OneDayRange(d0329).DurationIn(london), time.Hour*23)
	isEq(t, 0, OneDayRange(d1025).DurationIn(london), time.Hour*25)
	isEq(t, 0, NewDateRange(d0328, d0331).DurationIn(london), time.Hour*71)
}

func isEq(t *testing.T, i int, a, b interface{}, msg ...interface{}) {
	t.Helper()
	if a != b {
		sa := make([]string, len(msg))
		for i, m := range msg {
			sa[i] = fmt.Sprintf(", %v", m)
		}
		t.Errorf("%d: %+v is not equal to %+v%s", i, a, b, strings.Join(sa, ""))
	}
}

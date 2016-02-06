// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package timespan

import (
	"fmt"
	. "github.com/rickb777/date"
	"github.com/rickb777/date/period"
	"runtime/debug"
	"strings"
	"testing"
	"time"
)

var d0320 = New(2015, time.March, 20)
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
	dr := NewDateRangeOf(t0327, time.Duration(7*24*60*60*1e9))
	isEq(t, dr.mark, d0327)
	isEq(t, dr.Days(), PeriodOfDays(7))
	isEq(t, dr.IsEmpty(), false)
	isEq(t, dr.Start(), d0327)
	isEq(t, dr.Last(), d0402)
	isEq(t, dr.End(), d0403)
}

func TestNewDateRangeWithNormalise(t *testing.T) {
	r1 := NewDateRange(d0327, d0402)
	isEq(t, r1.Start(), d0327)
	isEq(t, r1.Last(), d0401)
	isEq(t, r1.End(), d0402)

	r2 := NewDateRange(d0402, d0327)
	isEq(t, r2.Start(), d0327)
	isEq(t, r2.Last(), d0401)
	isEq(t, r2.End(), d0402)
}

func TestEmptyRange(t *testing.T) {
	drN0 := DateRange{d0327, -1}
	isEq(t, drN0.Days(), PeriodOfDays(1))
	isEq(t, drN0.IsZero(), false)
	isEq(t, drN0.IsEmpty(), false)
	isEq(t, drN0.Start(), d0327)
	isEq(t, drN0.Last(), d0327)
	isEq(t, drN0.String(), "1 day on 2015-03-26")

	dr0 := DateRange{}
	isEq(t, dr0.Days(), PeriodOfDays(0))
	isEq(t, dr0.IsZero(), true)
	isEq(t, dr0.IsEmpty(), true)
	isEq(t, dr0.String(), "0 days at 1970-01-01")

	dr1 := EmptyRange(Date{})
	isEq(t, dr1.IsZero(), true)
	isEq(t, dr1.IsEmpty(), true)
	isEq(t, dr1.Days(), PeriodOfDays(0))

	dr2 := EmptyRange(d0327)
	isEq(t, dr2.IsZero(), false)
	isEq(t, dr2.IsEmpty(), true)
	isEq(t, dr2.Start(), d0327)
	isEq(t, dr2.Last().IsZero(), true)
	isEq(t, dr2.End(), d0327)
	isEq(t, dr2.Days(), PeriodOfDays(0))
	isEq(t, dr2.String(), "0 days at 2015-03-27")
}

func TestOneDayRange(t *testing.T) {
	dr1 := OneDayRange(Date{})
	isEq(t, dr1.IsZero(), false)
	isEq(t, dr1.IsEmpty(), false)
	isEq(t, dr1.Days(), PeriodOfDays(1))

	dr2 := OneDayRange(d0327)
	isEq(t, dr2.Start(), d0327)
	isEq(t, dr2.Last(), d0327)
	isEq(t, dr2.End(), d0328)
	isEq(t, dr2.Days(), PeriodOfDays(1))
	isEq(t, dr2.String(), "1 day on 2015-03-27")
}

func TestNewYearOf(t *testing.T) {
	dr := NewYearOf(2015)
	isEq(t, dr.Days(), PeriodOfDays(365))
	isEq(t, dr.Start(), New(2015, time.January, 1))
	isEq(t, dr.Last(), New(2015, time.December, 31))
	isEq(t, dr.End(), New(2016, time.January, 1))
}

func TestNewMonthOf(t *testing.T) {
	dr := NewMonthOf(2015, time.February)
	isEq(t, dr.Days(), PeriodOfDays(28))
	isEq(t, dr.Start(), New(2015, time.February, 1))
	isEq(t, dr.Last(), New(2015, time.February, 28))
	isEq(t, dr.End(), New(2015, time.March, 1))
}

func TestShiftByPos(t *testing.T) {
	dr := NewDateRange(d0327, d0402).ShiftBy(7)
	isEq(t, dr.Days(), PeriodOfDays(6))
	isEq(t, dr.Start(), d0403)
	isEq(t, dr.Last(), d0408)
}

func TestShiftByNeg(t *testing.T) {
	dr := NewDateRange(d0403, d0408).ShiftBy(-7)
	isEq(t, dr.Days(), PeriodOfDays(5))
	isEq(t, dr.Start(), d0327)
	isEq(t, dr.Last(), d0331)
}

func TestExtendByPos(t *testing.T) {
	dr := OneDayRange(d0327).ExtendBy(6)
	isEq(t, dr.Days(), PeriodOfDays(7))
	isEq(t, dr.Start(), d0327)
	isEq(t, dr.Last(), d0402)
	isEq(t, dr.End(), d0403)
	isEq(t, dr.String(), "7 days from 2015-03-27 to 2015-04-02")
}

func TestExtendByNeg(t *testing.T) {
	dr := OneDayRange(d0327).ExtendBy(-8)
	isEq(t, dr.Days(), PeriodOfDays(7))
	isEq(t, dr.Start(), d0320)
	isEq(t, dr.Last(), d0326)
	isEq(t, dr.String(), "7 days from 2015-03-20 to 2015-03-26")
}

func TestShiftByPosPeriod(t *testing.T) {
	dr := NewDateRange(d0327, d0402).ShiftByPeriod(period.NewPeriod(0, 0, 7))
	isEq(t, dr.Days(), PeriodOfDays(6))
	isEq(t, dr.Start(), d0403)
	isEq(t, dr.Last(), d0408)
}

func TestShiftByNegPeriod(t *testing.T) {
	dr := NewDateRange(d0403, d0408).ShiftByPeriod(period.NewPeriod(0, 0, -7))
	isEq(t, dr.Days(), PeriodOfDays(5))
	isEq(t, dr.Start(), d0327)
	isEq(t, dr.Last(), d0331)
}

func TestExtendByPosPeriod(t *testing.T) {
	dr := OneDayRange(d0327).ExtendByPeriod(period.NewPeriod(0, 0, 6))
	isEq(t, dr.Days(), PeriodOfDays(7))
	isEq(t, dr.Start(), d0327)
	isEq(t, dr.Last(), d0402)
	isEq(t, dr.End(), d0403)
	isEq(t, dr.String(), "7 days from 2015-03-27 to 2015-04-02")
}

func TestExtendByNegPeriod(t *testing.T) {
	dr := OneDayRange(d0327).ExtendByPeriod(period.NewPeriod(0, 0, -8))
	//fmt.Printf("\ndr=%#v\n", dr)
	isEq(t, dr.Days(), PeriodOfDays(7))
	isEq(t, dr.Start(), d0320)
	isEq(t, dr.Last(), d0326)
	isEq(t, dr.String(), "7 days from 2015-03-20 to 2015-03-26")
}

func TestContains1(t *testing.T) {
	old := time.Local
	time.Local = time.FixedZone("Test", 7200)
	dr := OneDayRange(d0326).ExtendBy(1)
	isEq(t, dr.Contains(d0320), false, dr, d0320)
	isEq(t, dr.Contains(d0325), false, dr, d0325)
	isEq(t, dr.Contains(d0326), true, dr, d0326)
	isEq(t, dr.Contains(d0327), true, dr, d0327)
	isEq(t, dr.Contains(d0328), false, dr, d0328)
	isEq(t, dr.Contains(d0401), false, dr, d0401)
	isEq(t, dr.Contains(d0410), false, dr, d0410)
	isEq(t, dr.Contains(d0501), false, dr, d0501)
	time.Local = old
}

func TestContains2(t *testing.T) {
	old := time.Local
	time.Local = time.FixedZone("Test", 7200)
	dr := OneDayRange(d0326)
	isEq(t, dr.Contains(d0325), false, dr, d0325)
	isEq(t, dr.Contains(d0326), true, dr, d0326)
	isEq(t, dr.Contains(d0327), false, dr, d0327)
	time.Local = old
}

func TestContainsTimeUTC(t *testing.T) {
	old := time.Local
	time.Local = time.FixedZone("Test", 7200)
	t0328e := time.Date(2015, 3, 28, 23, 59, 59, 999999999, time.UTC)
	t0329 := time.Date(2015, 3, 29, 0, 0, 0, 0, time.UTC)

	dr := OneDayRange(d0327).ExtendBy(1)
	isEq(t, dr.StartUTC(), t0327, dr, t0327)
	isEq(t, dr.EndUTC(), t0329, dr, t0329)
	isEq(t, dr.ContainsTime(t0327), true, dr, t0327)
	isEq(t, dr.ContainsTime(t0328), true, dr, t0328)
	isEq(t, dr.ContainsTime(t0328e), true, dr, t0328e)
	isEq(t, dr.ContainsTime(t0329), false, dr, t0329)
	time.Local = old
}

func TestMerge1(t *testing.T) {
	dr1 := OneDayRange(d0327).ExtendBy(1)
	dr2 := OneDayRange(d0327).ExtendBy(7)
	m1 := dr1.Merge(dr2)
	m2 := dr2.Merge(dr1)
	isEq(t, m1.Start(), d0327)
	isEq(t, m1.End(), d0404)
	isEq(t, m1, m2)
}

func TestMerge2(t *testing.T) {
	dr1 := OneDayRange(d0327).ExtendBy(1).ShiftBy(1)
	dr2 := OneDayRange(d0327).ExtendBy(7)
	m1 := dr1.Merge(dr2)
	m2 := dr2.Merge(dr1)
	isEq(t, m1.Start(), d0327)
	isEq(t, m1.End(), d0404)
	isEq(t, m1, m2)
}

func TestMergeOverlapping(t *testing.T) {
	dr1 := OneDayRange(d0320).ExtendBy(12)
	dr2 := OneDayRange(d0401).ExtendBy(6)
	m1 := dr1.Merge(dr2)
	m2 := dr2.Merge(dr1)
	isEq(t, m1.Start(), d0320)
	isEq(t, m1.End(), d0408)
	isEq(t, m1, m2)
}

func TestMergeNonOverlapping(t *testing.T) {
	dr1 := OneDayRange(d0320).ExtendBy(2)
	dr2 := OneDayRange(d0401).ExtendBy(6)
	m1 := dr1.Merge(dr2)
	m2 := dr2.Merge(dr1)
	isEq(t, m1.Start(), d0320)
	isEq(t, m1.End(), d0408)
	isEq(t, m1, m2)
}

func TestMergeEmpties(t *testing.T) {
	dr1 := EmptyRange(d0320)
	dr2 := EmptyRange(d0408) // curiously, this is *not* included because it has no size.
	m1 := dr1.Merge(dr2)
	m2 := dr2.Merge(dr1)
	isEq(t, m1.Start(), d0320)
	isEq(t, m1.End(), d0408)
	isEq(t, m1, m2)
}

func TestMergeZeroes(t *testing.T) {
	dr0 := DateRange{}
	dr1 := OneDayRange(d0401).ExtendBy(6)
	m1 := dr1.Merge(dr0)
	m2 := dr0.Merge(dr1)
	m3 := dr0.Merge(dr0)
	isEq(t, m1.Start(), d0401)
	isEq(t, m1.End(), d0408)
	isEq(t, m1, m2)
	isEq(t, m3.IsZero(), true)
	isEq(t, m3, dr0)
}

func TestDurationNormalUTC(t *testing.T) {
	dr := OneDayRange(d0329)
	isEq(t, dr.Duration(), time.Hour*24)
}

func TestDurationInZoneWithDaylightSaving(t *testing.T) {
	isEq(t, OneDayRange(d0328).DurationIn(london), time.Hour*24)
	isEq(t, OneDayRange(d0329).DurationIn(london), time.Hour*23)
	isEq(t, OneDayRange(d1025).DurationIn(london), time.Hour*25)
	isEq(t, NewDateRange(d0328, d0331).DurationIn(london), time.Hour*71)
}

func isEq(t *testing.T, a, b interface{}, msg ...interface{}) {
	if a != b {
		sa := make([]string, len(msg))
		for i, m := range msg {
			sa[i] = fmt.Sprintf(", %v", m)
		}
		t.Errorf("%v (%#v) is not equal to %v (%#v)%s\n%s", a, a, b, b, strings.Join(sa, ""), debug.Stack())
	}
}

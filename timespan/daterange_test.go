// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package timespan

import (
	. "github.com/rickb777/date"
	"testing"
	"time"
	"fmt"
	"strings"
	"runtime/debug"
)

var d0320 = New(2015, time.March, 20)
var d0325 = New(2015, time.March, 25)
var d0326 = New(2015, time.March, 26)
var d0327 = New(2015, time.March, 27)
var d0328 = New(2015, time.March, 28)
var d0329 = New(2015, time.March, 29) // n.b. clocks go forward (UK)
var d0330 = New(2015, time.March, 30)
var d0401 = New(2015, time.April, 1)
var d0402 = New(2015, time.April, 2)
var d0403 = New(2015, time.April, 3)
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
	isEq(t, dr.Start(), d0327)
	isEq(t, dr.End(), d0402)
	isEq(t, dr.Next(), d0403)
}

func TestNewDateRangeWithNormalise(t *testing.T) {
	r1 := NewDateRange(d0327, d0401)
	isEq(t, r1.Start(), d0327)
	isEq(t, r1.End(), d0401)
	isEq(t, r1.Next(), d0402)

	r2 := NewDateRange(d0401, d0327)
	isEq(t, r2.Start(), d0327)
	isEq(t, r2.End(), d0401)
	isEq(t, r2.Next(), d0402)
}

func TestOneDayRange(t *testing.T) {
	drN0 := DateRange{d0327, -1}
	isEq(t, drN0.Days(), PeriodOfDays(-1))
	isEq(t, drN0.Start(), d0327)
	isEq(t, drN0.End(), d0327)
	isEq(t, drN0.String(), "1 day on 2015-03-27")

	dr0 := DateRange{}
	isEq(t, dr0.Days(), PeriodOfDays(0))
	isEq(t, dr0.String(), "0 days from 1970-01-01")

	dr1 := OneDayRange(Date{})
	isEq(t, dr1.Days(), PeriodOfDays(1))

	dr2 := OneDayRange(d0327)
	isEq(t, dr2.Start(), d0327)
	isEq(t, dr2.End(), d0327)
	isEq(t, dr2.Next(), d0328)
	isEq(t, dr2.Days(), PeriodOfDays(1))
	isEq(t, dr2.String(), "1 day on 2015-03-27")
}

func TestNewYearOf(t *testing.T) {
	dr := NewYearOf(2015)
	isEq(t, dr.Days(), PeriodOfDays(365))
	isEq(t, dr.Start(), New(2015, time.January, 1))
	isEq(t, dr.End(), New(2015, time.December, 31))
	isEq(t, dr.Next(), New(2016, time.January, 1))
}

func TestNewMonthOf(t *testing.T) {
	dr := NewMonthOf(2015, time.February)
	isEq(t, dr.Days(), PeriodOfDays(28))
	isEq(t, dr.Start(), New(2015, time.February, 1))
	isEq(t, dr.End(), New(2015, time.February, 28))
	isEq(t, dr.Next(), New(2015, time.March, 1))
}

func TestShiftByPos(t *testing.T) {
	dr := NewDateRange(d0327, d0401).ShiftBy(7)
	isEq(t, dr.Days(), PeriodOfDays(6))
	isEq(t, dr.Start(), d0403)
	isEq(t, dr.End(), d0408)
}

func TestShiftByNeg(t *testing.T) {
	dr := NewDateRange(d0403, d0408).ShiftBy(-7)
	isEq(t, dr.Days(), PeriodOfDays(6))
	isEq(t, dr.Start(), d0327)
	isEq(t, dr.End(), d0401)
}

func TestExtendByPos(t *testing.T) {
	dr := OneDayRange(d0327).ExtendBy(6)
	isEq(t, dr.Days(), PeriodOfDays(7))
	isEq(t, dr.Start(), d0327)
	isEq(t, dr.End(), d0402)
	isEq(t, dr.Next(), d0403)
	isEq(t, dr.String(), "7 days from 2015-03-27 to 2015-04-02")
}

func TestExtendByNeg(t *testing.T) {
	dr := OneDayRange(d0327).ExtendBy(-9)
	isEq(t, dr.Days(), PeriodOfDays(-8))
	isEq(t, dr.Start(), d0320)
	isEq(t, dr.End(), d0327)
	isEq(t, dr.String(), "8 days from 2015-03-20 to 2015-03-27")
}

func xTestContains1(t *testing.T) {
	old := time.Local
	time.Local = time.FixedZone("Test", 7200)
	dr := OneDayRange(d0326).ExtendBy(1)
	isEq(t, dr.Contains(d0320), false, dr, d0320)
	isEq(t, dr.Contains(d0325), false, dr, d0325)
	isEq(t, dr.Contains(d0326), true,  dr, d0326)
	isEq(t, dr.Contains(d0327), true,  dr, d0327)
	isEq(t, dr.Contains(d0328), false, dr, d0328)
	isEq(t, dr.Contains(d0401), false, dr, d0401)
	isEq(t, dr.Contains(d0410), false, dr, d0410)
	isEq(t, dr.Contains(d0501), false, dr, d0501)
	time.Local = old
}

func xTestContains2(t *testing.T) {
	old := time.Local
	time.Local = time.FixedZone("Test", 7200)
	dr := OneDayRange(d0326)
	isEq(t, dr.Contains(d0325), false, dr, d0325)
	isEq(t, dr.Contains(d0326), true,  dr, d0326)
	isEq(t, dr.Contains(d0327), false, dr, d0327)
	time.Local = old
}

func xTestContainsTimeUTC(t *testing.T) {
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

func xTestMerge1(t *testing.T) {
	dr1 := OneDayRange(d0327).ExtendBy(1)
	dr2 := OneDayRange(d0327).ExtendBy(7)
	m1 := dr1.Merge(dr2)
	m2 := dr2.Merge(dr1)
	isEq(t, m1.Start(), d0327)
	isEq(t, m1.End(), d0403)
	isEq(t, m1, m2)
}

func xTestMerge2(t *testing.T) {
	dr1 := OneDayRange(d0327).ExtendBy(1).ShiftBy(1)
	dr2 := OneDayRange(d0327).ExtendBy(7)
	m1 := dr1.Merge(dr2)
	m2 := dr2.Merge(dr1)
	isEq(t, m1.Start(), d0327)
	isEq(t, m1.End(), d0403)
	isEq(t, m1, m2)
}

func xTestMergeOverlapping(t *testing.T) {
	dr1 := OneDayRange(d0320).ExtendBy(12)
	dr2 := OneDayRange(d0401).ExtendBy(7)
	m1 := dr1.Merge(dr2)
	m2 := dr2.Merge(dr1)
	isEq(t, m1.Start(), d0320)
	isEq(t, m1.End(), d0408)
	isEq(t, m1, m2)
}

func xTestMergeNonOverlapping(t *testing.T) {
	dr1 := OneDayRange(d0320).ExtendBy(2)
	dr2 := OneDayRange(d0401).ExtendBy(7)
	m1 := dr1.Merge(dr2)
	m2 := dr2.Merge(dr1)
	isEq(t, m1.Start(), d0320)
	isEq(t, m1.End(), d0408)
	isEq(t, m1, m2)
}

func xTestDurationNormalUTC(t *testing.T) {
	dr := OneDayRange(d0329)
	isEq(t, dr.Duration(), time.Hour * 24)
}

func xTestDurationInZoneWithDaylightSaving(t *testing.T) {
	isEq(t, OneDayRange(d0328).DurationIn(london), time.Hour * 24)
	isEq(t, OneDayRange(d0329).DurationIn(london), time.Hour * 23)
	isEq(t, OneDayRange(d1025).DurationIn(london), time.Hour * 25)
	isEq(t, NewDateRange(d0328, d0330).DurationIn(london), time.Hour * 71)
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

// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package timsepan

import (
	. "github.com/rickb777/date"
	"testing"
	"time"
	"fmt"
	"strings"
	"runtime/debug"
)

var t0327 = time.Date(2015, 3, 27, 0, 0, 0, 0, time.UTC)
var t0328 = time.Date(2015, 3, 28, 0, 0, 0, 0, time.UTC)
var d0320 = New(2015, time.March, 20)
var d0325 = New(2015, time.March, 25)
var d0326 = New(2015, time.March, 26)
var d0327 = New(2015, time.March, 27)
var d0328 = New(2015, time.March, 28)
var d0329 = New(2015, time.March, 29)
var d0401 = New(2015, time.April, 1)
var d0403 = New(2015, time.April, 3)
var d0408 = New(2015, time.April, 8)
var d0410 = New(2015, time.April, 10)
var d0501 = New(2015, time.May, 1)
var d1025 = New(2015, time.October, 25)

func TestNewDateRangeOf(t *testing.T) {
	dr := NewDateRangeOf(t0327, time.Duration(7*24*60*60*1e9))
	isEq(t, dr.Start, d0327)
	isEq(t, dr.End, d0403)
}

func TestNewDateRangeWithNormalise(t *testing.T) {
	r1 := NewDateRange(d0327, d0401)
	isEq(t, r1.Start, d0327)
	isEq(t, r1.End, d0401)

	r2 := NewDateRange(d0401, d0327)
	isEq(t, r2.Start, d0327)
	isEq(t, r2.End, d0401)
}

func TestOneDayRange(t *testing.T) {
	dr := OneDayRange(d0327)
	isEq(t, dr.Start, d0327)
	isEq(t, dr.End, d0327)
}

func TestNewYearOf(t *testing.T) {
	dr := NewYearOf(2015)
	isEq(t, dr.Start, New(2015, time.January, 1))
	isEq(t, dr.End, New(2015, time.December, 31))
}

func TestNewMonthOf(t *testing.T) {
	dr := NewMonthOf(2015, time.February)
	isEq(t, dr.Start, New(2015, time.February, 1))
	isEq(t, dr.End, New(2015, time.February, 28))
}

func TestShiftByPos(t *testing.T) {
	dr := NewDateRange(d0327, d0401).ShiftBy(7)
	isEq(t, dr.Start, d0403)
	isEq(t, dr.End, d0408)
}

func TestShiftByNeg(t *testing.T) {
	dr := NewDateRange(d0403, d0408).ShiftBy(-7)
	isEq(t, dr.Start, d0327)
	isEq(t, dr.End, d0401)
}

func TestExtendByPos(t *testing.T) {
	dr := OneDayRange(d0327).ExtendBy(7)
	isEq(t, dr.Start, d0327)
	isEq(t, dr.End, d0403)
	isEq(t, dr.String(), "2015-03-27 to 2015-04-03")
}

func TestExtendByNeg(t *testing.T) {
	dr := OneDayRange(d0327).ExtendBy(-7)
	isEq(t, dr.Start, d0320)
	isEq(t, dr.End, d0327)
	isEq(t, dr.String(), "2015-03-20 to 2015-03-27")
}

func TestContains1(t *testing.T) {
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

func TestContains2(t *testing.T) {
	old := time.Local
	time.Local = time.FixedZone("Test", 7200)
	dr := OneDayRange(d0326)
	isEq(t, dr.Contains(d0325), false, dr, d0325)
	isEq(t, dr.Contains(d0326), true,  dr, d0326)
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
	isEq(t, m1.Start, d0327)
	isEq(t, m1.End, d0403)
	isEq(t, m1, m2)
}

func TestMerge2(t *testing.T) {
	dr1 := OneDayRange(d0327).ExtendBy(1).ShiftBy(1)
	dr2 := OneDayRange(d0327).ExtendBy(7)
	m1 := dr1.Merge(dr2)
	m2 := dr2.Merge(dr1)
	isEq(t, m1.Start, d0327)
	isEq(t, m1.End, d0403)
	isEq(t, m1, m2)
}

func TestMergeOverlapping(t *testing.T) {
	dr1 := OneDayRange(d0320).ExtendBy(12)
	dr2 := OneDayRange(d0401).ExtendBy(7)
	m1 := dr1.Merge(dr2)
	m2 := dr2.Merge(dr1)
	isEq(t, m1.Start, d0320)
	isEq(t, m1.End, d0408)
	isEq(t, m1, m2)
}

func TestMergeNonOverlapping(t *testing.T) {
	dr1 := OneDayRange(d0320).ExtendBy(2)
	dr2 := OneDayRange(d0401).ExtendBy(7)
	m1 := dr1.Merge(dr2)
	m2 := dr2.Merge(dr1)
	isEq(t, m1.Start, d0320)
	isEq(t, m1.End, d0408)
	isEq(t, m1, m2)
}

func TestDurationNormalUTC(t *testing.T) {
	dr := OneDayRange(d0329)
	isEq(t, dr.Duration(), time.Hour * 24)
}

func TestDurationInZoneWithDaylightSaving(t *testing.T) {
	london, err := time.LoadLocation("Europe/London")
	if err != nil {
		panic(err)
	}
	isEq(t, OneDayRange(d0328).DurationIn(london), time.Hour * 24)
	isEq(t, OneDayRange(d0329).DurationIn(london), time.Hour * 23)
	isEq(t, OneDayRange(d1025).DurationIn(london), time.Hour * 25)
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

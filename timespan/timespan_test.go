// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package timespan

import (
	"testing"
	"time"
)

const zero time.Duration = 0

var t0327 = time.Date(2015, 3, 27, 0, 0, 0, 0, time.UTC)
var t0328 = time.Date(2015, 3, 28, 0, 0, 0, 0, time.UTC)
var t0329 = time.Date(2015, 3, 29, 0, 0, 0, 0, time.UTC) // n.b. clocks go forward (UK)
var t0330 = time.Date(2015, 3, 30, 0, 0, 0, 0, time.UTC)

func TestZeroTimeSpan(t *testing.T) {
	ts := ZeroTimeSpan(t0327)
	isEq(t, ts.mark, t0327)
	isEq(t, ts.Duration(), zero)
	isEq(t, ts.End(), t0327)
}

func TestNewTimeSpan(t *testing.T) {
	ts1 := NewTimeSpan(t0327, t0327)
	isEq(t, ts1.mark, t0327)
	isEq(t, ts1.Duration(), zero)
	isEq(t, ts1.End(), t0327)

	ts2 := NewTimeSpan(t0327, t0328)
	isEq(t, ts2.mark, t0327)
	isEq(t, ts2.Duration(), time.Hour * 24)
	isEq(t, ts2.End(), t0328)

	ts3 := NewTimeSpan(t0329, t0327)
	isEq(t, ts3.mark, t0327)
	isEq(t, ts3.Duration(), time.Hour * 48)
	isEq(t, ts3.End(), t0329)
}

func TestTSEnd(t *testing.T) {
	ts1 := TimeSpan{t0328, time.Hour * 24}
	isEq(t, ts1.Start(), t0328)
	isEq(t, ts1.End(), t0329)

	// not normalised, deliberately
	ts2 := TimeSpan{t0328, -time.Hour * 24}
	isEq(t, ts2.Start(), t0327)
	isEq(t, ts2.End(), t0328)
}

func TestTSShiftBy(t *testing.T) {
	ts1 := NewTimeSpan(t0327, t0328).ShiftBy(time.Hour * 24)
	isEq(t, ts1.mark, t0328)
	isEq(t, ts1.Duration(), time.Hour * 24)
	isEq(t, ts1.End(), t0329)

	ts2 := NewTimeSpan(t0328, t0329).ShiftBy(-time.Hour * 24)
	isEq(t, ts2.mark, t0327)
	isEq(t, ts2.Duration(), time.Hour * 24)
	isEq(t, ts2.End(), t0328)
}

func TestTSExtendBy(t *testing.T) {
	ts1 := NewTimeSpan(t0327, t0328).ExtendBy(time.Hour * 24)
	isEq(t, ts1.mark, t0327)
	isEq(t, ts1.Duration(), time.Hour * 48)
	isEq(t, ts1.End(), t0329)

	ts2 := NewTimeSpan(t0328, t0329).ExtendBy(-time.Hour * 48)
	isEq(t, ts2.mark, t0327)
	isEq(t, ts2.Duration(), time.Hour * 24)
	isEq(t, ts2.End(), t0328)
}

func TestTSExtendWithoutWrapping(t *testing.T) {
	ts1 := NewTimeSpan(t0327, t0328).ExtendWithoutWrapping(time.Hour * 24)
	isEq(t, ts1.mark, t0327)
	isEq(t, ts1.Duration(), time.Hour * 48)
	isEq(t, ts1.End(), t0329)

	ts2 := NewTimeSpan(t0328, t0329).ExtendWithoutWrapping(-time.Hour * 48)
	isEq(t, ts2.mark, t0328)
	isEq(t, ts2.Duration(), zero)
	isEq(t, ts2.End(), t0328)
}

func TestTSString(t *testing.T) {
	s := NewTimeSpan(t0327, t0328).String()
	isEq(t, s, "24h0m0s from 2015-03-27 00:00:00 to 2015-03-28 00:00:00")
}

func TestTSContains(t *testing.T) {
	ts := NewTimeSpan(t0327, t0329)
	isEq(t, ts.Contains(t0327.Add(minusOneNano)), false)
	isEq(t, ts.Contains(t0327), true)
	isEq(t, ts.Contains(t0328), true)
	isEq(t, ts.Contains(t0329.Add(minusOneNano)), true)
	isEq(t, ts.Contains(t0329), false)
}

func TestTSIn(t *testing.T) {
	ts := ZeroTimeSpan(t0327).In(time.FixedZone("Test", 7200))
	isEq(t, ts.mark.Equal(t0327), true)
	isEq(t, ts.Duration(), zero)
	isEq(t, ts.End().Equal(t0327), true)
}

func TestTSMerge1(t *testing.T) {
	ts1 := NewTimeSpan(t0327, t0328)
	ts2 := NewTimeSpan(t0327, t0330)
	m1 := ts1.Merge(ts2)
	m2 := ts2.Merge(ts1)
	isEq(t, m1.mark, t0327)
	isEq(t, m1.End(), t0330)
	isEq(t, m1, m2)
}

func TestTSMerge2(t *testing.T) {
	ts1 := NewTimeSpan(t0328, t0329)
	ts2 := NewTimeSpan(t0327, t0330)
	m1 := ts1.Merge(ts2)
	m2 := ts2.Merge(ts1)
	isEq(t, m1.mark, t0327)
	isEq(t, m1.End(), t0330)
	isEq(t, m1, m2)
}

func TestTSMerge3(t *testing.T) {
	ts1 := NewTimeSpan(t0329, t0330)
	ts2 := NewTimeSpan(t0327, t0330)
	m1 := ts1.Merge(ts2)
	m2 := ts2.Merge(ts1)
	isEq(t, m1.mark, t0327)
	isEq(t, m1.End(), t0330)
	isEq(t, m1, m2)
}

func TestTSMergeOverlapping(t *testing.T) {
	ts1 := NewTimeSpan(t0327, t0329)
	ts2 := NewTimeSpan(t0328, t0330)
	m1 := ts1.Merge(ts2)
	m2 := ts2.Merge(ts1)
	isEq(t, m1.mark, t0327)
	isEq(t, m1.End(), t0330)
	isEq(t, m1, m2)
}

func xTestTSMergeNonOverlapping(t *testing.T) {
	ts1 := NewTimeSpan(t0327, t0328)
	ts2 := NewTimeSpan(t0329, t0330)
	m1 := ts1.Merge(ts2)
	m2 := ts2.Merge(ts1)
	isEq(t, m1.mark, t0327)
	isEq(t, m1.End(), t0330)
	isEq(t, m1, m2)
}

func xTestConversion1(t *testing.T) {
	ts1 := ZeroTimeSpan(t0327)
	dr := ts1.DateRangeIn(time.UTC)
	ts2 := dr.TimeSpanIn(time.UTC)
	isEq(t, dr.Start, d0327)
	isEq(t, dr.End, d0327)
	isEq(t, ts1, ts2)
	isEq(t, ts1.Duration(), zero)
}

func xTestConversion2(t *testing.T) {
	ts1 := NewTimeSpan(t0327, t0328)
	dr := ts1.DateRangeIn(time.UTC)
//	ts2 := dr.TimeSpanIn(time.UTC)
	isEq(t, dr.Start, d0327)
	isEq(t, dr.End, d0328)
//	isEq(t, ts1, ts2)
	isEq(t, ts1.Duration(), time.Hour * 24)
}

func xTestConversion3(t *testing.T) {
	dr1 := NewDateRange(d0327, d0330) // weekend of clocks changing
	ts1 := dr1.TimeSpanIn(london)
	dr2 := ts1.DateRangeIn(london)
//	ts2 := dr2.TimeSpanIn(london)
	isEq(t, dr1.Start, d0327)
	isEq(t, dr1.End, d0330)
	isEq(t, dr1, dr2)
//	isEq(t, ts1, ts2)
	isEq(t, ts1.Duration(), time.Hour * 71)
}


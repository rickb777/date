// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package daterange

import (
	. "github.com/rickb777/date"
	"testing"
	"time"
)

var t0327 = time.Date(2015, 3, 27, 0, 0, 0, 0, time.UTC)
var d0320 = New(2015, time.March, 20)
var d0327 = New(2015, time.March, 27)
var d0401 = New(2015, time.April, 1)
var d0403 = New(2015, time.April, 3)

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

//func TestAddToStart1(t *testing.T) {
//	timeSpan := ZeroTimeSpan(date0327).AddYMDToStart(0, 0, 7)
//	is1(t, timeSpan.Start, date0327)
//	is1(t, timeSpan.End, date0403)
//	is1(t, timeSpan.String(), "2015-03-27 to 2015-04-03")
//}
//
//func TestAddToStart2(t *testing.T) {
//	timeSpan := ZeroTimeSpan(date0327).AddYMDToStart(0, 0, -7)
//	is1(t, timeSpan.Start, date0320)
//	is1(t, timeSpan.End, date0327)
//	is1(t, timeSpan.String(), "2015-03-20 to 2015-03-27")
//}

func TestExtendBy1(t *testing.T) {
	timeSpan := OneDayRange(d0327).ExtendBy(7)
	isEq(t, timeSpan.Start, d0327)
	isEq(t, timeSpan.End, d0403)
	isEq(t, timeSpan.String(), "2015-03-27 to 2015-04-03")
}

func TestExtendBy2(t *testing.T) {
	timeSpan := OneDayRange(d0327).ExtendBy(-7)
	isEq(t, timeSpan.Start, d0320)
	isEq(t, timeSpan.End, d0327)
	isEq(t, timeSpan.String(), "2015-03-20 to 2015-03-27")
}

//func TestAddToEnd2(t *testing.T) {
//	timeSpan := ZeroTimeSpan(date0327).AddYMDToEnd(0, 0, -7)
//	is1(t, timeSpan.Start, date0320)
//	is1(t, timeSpan.End, date0327)
//	is1(t, timeSpan.String(), "2015-03-20 to 2015-03-27")
//}


func isEq(t *testing.T, a, b interface{}) {
	if a != b {
		t.Errorf("%s %#v is not equal to %s %#v", a, a, b, b)
	}
}

// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package date

import (
	"database/sql/driver"
	"testing"
)

func TestDateScan(t *testing.T) {
	cases := []struct {
		v        interface{}
		disallow bool
		expected PeriodOfDays
	}{
		{int64(0), false, 0},
		{int64(1000), false, 1000},
		{int64(10000), false, 10000},
		{int64(0), true, 0},
		{int64(1000), true, 1000},
		{int64(10000), true, 10000},
		{"0", false, 0},
		{"1000", false, 1000},
		{"10000", false, 10000},
		{[]byte("10000"), false, 10000},
		{PeriodOfDays(10000).Date().Local(), false, 10000},
	}

	prior := DisableTextStorage

	for i, c := range cases {
		DisableTextStorage = c.disallow
		r := new(Date)
		e := r.Scan(c.v)
		if e != nil {
			t.Errorf("%d: Got %v for %d", i, e, c.expected)
		}
		if r.DaysSinceEpoch() != c.expected {
			t.Errorf("%d: Got %v, want %d", i, *r, c.expected)
		}

		var d driver.Valuer = *r

		q, e := d.Value()
		if e != nil {
			t.Errorf("%d: Got %v for %d", i, e, c.expected)
		}
		if q.(int64) != int64(c.expected) {
			t.Errorf("%d: Got %v, want %d", i, q, c.expected)
		}
	}

	DisableTextStorage = prior
}

func TestDateScanWithJunk(t *testing.T) {
	cases := []struct {
		v        interface{}
		disallow bool
		expected string
	}{
		{true, false, "bool true is not a meaningful date"},
		{true, true, "bool true is not a meaningful date"},
	}

	prior := DisableTextStorage

	for i, c := range cases {
		DisableTextStorage = c.disallow
		r := new(Date)
		e := r.Scan(c.v)
		if e.Error() != c.expected {
			t.Errorf("%d: Got %q, want %q", i, e.Error(), c.expected)
		}
	}

	DisableTextStorage = prior
}

func TestDateScanWithNil(t *testing.T) {
	var r *Date
	e := r.Scan(nil)
	if e != nil {
		t.Errorf("Got %v", e)
	}
	if r != nil {
		t.Errorf("Got %v", r)
	}
}

// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package date

import (
	"database/sql/driver"
	"testing"
)

func TestDate_Scan(t *testing.T) {
	cases := []struct {
		v        interface{}
		expected PeriodOfDays
	}{
		{int64(0), 0},
		{int64(1000), 1000},
		{int64(10000), 10000},
		{int64(0), 0},
		{int64(1000), 1000},
		{int64(10000), 10000},
		{"0", 0},
		{"1000", 1000},
		{"10000", 10000},
		{"2018-12-31", 17896},
		{"31/12/2018", 17896},
		{[]byte("10000"), 10000},
		{PeriodOfDays(10000).Date().Local(), 10000},
	}

	for i, c := range cases {
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
}

func TestDateString_Scan(t *testing.T) {
	cases := []struct {
		v        interface{}
		expected string
	}{
		{int64(0), "1970-01-01"},
		{int64(15000), "2011-01-26"},
		{"0", "1970-01-01"},
		{"15000", "2011-01-26"},
		{"2018-12-31", "2018-12-31"},
		{"31/12/2018", "2018-12-31"},
		//{[]byte("10000"), ""},
		//{PeriodOfDays(10000).Date().Local(), ""},
	}

	for i, c := range cases {
		r := new(DateString)
		e := r.Scan(c.v)
		if e != nil {
			t.Errorf("%d: Got %v for %s", i, e, c.expected)
		}
		if r.Date().String() != c.expected {
			t.Errorf("%d: Got %v, want %s", i, r.Date(), c.expected)
		}

		var d driver.Valuer = *r

		q, e := d.Value()
		if e != nil {
			t.Errorf("%d: Got %v for %s", i, e, c.expected)
		}
		if q.(string) != c.expected {
			t.Errorf("%d: Got %v, want %s", i, q, c.expected)
		}
	}
}

func TestDate_Scan_with_junk(t *testing.T) {
	cases := []struct {
		v        interface{}
		expected string
	}{
		{true, "bool true is not a meaningful date"},
		{true, "bool true is not a meaningful date"},
	}

	for i, c := range cases {
		r := new(Date)
		e := r.Scan(c.v)
		if e.Error() != c.expected {
			t.Errorf("%d: Got %q, want %q", i, e.Error(), c.expected)
		}
	}
}

func TestDateString_Scan_with_junk(t *testing.T) {
	cases := []struct {
		v        interface{}
		expected string
	}{
		{true, "bool true is not a meaningful date"},
		{true, "bool true is not a meaningful date"},
	}

	for i, c := range cases {
		r := new(DateString)
		e := r.Scan(c.v)
		if e.Error() != c.expected {
			t.Errorf("%d: Got %q, want %q", i, e.Error(), c.expected)
		}
	}
}

func TestDate_Scan_with_nil(t *testing.T) {
	var r *Date
	e := r.Scan(nil)
	if e != nil {
		t.Errorf("Got %v", e)
	}
}

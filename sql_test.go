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
		expected Date
	}{
		{v: int64(0), expected: 0},
		{v: int64(1), expected: 1},
		{v: int64(1000), expected: 1000},
		{v: int64(10000), expected: 10000},
		//{v: "00000101", expected: 0},
		//{v: "00000102", expected: 1},
		{v: "0001-01-01", expected: 0},
		{v: "19700101", expected: zeroOffset},
		{v: "1970-01-01", expected: zeroOffset},
		{v: "1971-01-01", expected: 365 + zeroOffset},
		{v: "2018-12-31", expected: 17896 + zeroOffset},
		{v: "31/12/2018", expected: 17896 + zeroOffset},
		{v: []byte("19700101"), expected: zeroOffset},
		{v: Date(10000).Midnight(), expected: 10000},
	}

	for i, c := range cases {
		r := new(Date)
		e := r.Scan(c.v)
		if e != nil {
			t.Errorf("%d: Got %v for %d", i, e, c.expected)
		}
		if *r != c.expected {
			t.Errorf("%d: Got %v, want %d", i, *r, c.expected)
		}

		var d driver.Valuer = *r

		q, e := d.Value()
		if e != nil {
			t.Errorf("%d: Got %v for %d", i, e, c.expected)
		}
		if q.(string) != c.expected.String() {
			t.Errorf("%d: Got %v, want %d", i, q, c.expected)
		}

		q, e = ValueAsInt(*r)
		if e != nil {
			t.Errorf("%d: Got %v for %d", i, e, c.expected)
		}
		if q.(int64) != int64(c.expected) {
			t.Errorf("%d: Got %v, want %d", i, q, c.expected)
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

func TestDate_Scan_with_nil(t *testing.T) {
	var r *Date
	e := r.Scan(nil)
	if e != nil {
		t.Errorf("Got %v", e)
	}
}

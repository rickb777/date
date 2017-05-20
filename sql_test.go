// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package date

import (
	"database/sql/driver"
	"testing"
)

func TestDateScan(t *testing.T) {
	cases := []PeriodOfDays{
		0, 1, 28, 30, 31, 32, 364, 365, 366, 367, 500, 1000, 10000, 100000,
	}
	for _, c := range cases {
		var d driver.Valuer = NewOfDays(c)

		v, e := d.Value()
		if e != nil {
			t.Errorf("Got %v for %d", e, c)
		}
		if v.(int64) != int64(c) {
			t.Errorf("Got %v, want %d", v, c)
		}

		r := new(Date)
		e = r.Scan(v)
		if e != nil {
			t.Errorf("Got %v for %d", e, c)
		}
		if *r != d {
			t.Errorf("Got %v, want %d", *r, d)
		}
	}
}

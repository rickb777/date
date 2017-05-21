// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package date

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"time"
)

// These methods allow Date and PeriodOfDays to be fields stored in an
// SQL database by implementing the database/sql/driver interfaces.
// The underlying column type is simply an integer.

// Scan parses some value. It implements sql.Scanner,
// https://golang.org/pkg/database/sql/#Scanner
func (d *Date) Scan(value interface{}) (err error) {
	if DisableTextStorage {
		return d.scanInt(value)
	}
	var n int64
	err = nil
	switch value.(type) {
	case int64:
		*d = Date{PeriodOfDays(value.(int64))}
	case []byte:
		n, err = strconv.ParseInt(string(value.([]byte)), 10, 64)
		*d = Date{PeriodOfDays(n)}
	case string:
		n, err = strconv.ParseInt(value.(string), 10, 64)
		*d = Date{PeriodOfDays(n)}
	case time.Time:
		*d = NewAt(value.(time.Time))
	default:
		err = fmt.Errorf("%#v", value)
	}
	return
}

func (d *Date) scanInt(value interface{}) (err error) {
	err = nil
	switch value.(type) {
	case int64:
		*d = Date{PeriodOfDays(value.(int64))}
	default:
		err = fmt.Errorf("%#v", value)
	}
	return
}

// Value converts the value to an int64. It implements driver.Valuer,
// https://golang.org/pkg/database/sql/driver/#Valuer
func (d Date) Value() (driver.Value, error) {
	return int64(d.day), nil
}

// DisableTextStorage reduces the Scan method so that only integers are handled.
// Normally, database types int64, []byte, string and time.Time are supported.
// When set true, only int64 is supported; this mode allows optimisation of SQL
// result processing and would only be used during development.
var DisableTextStorage = false

// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package date

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// These methods allow Date to be stored in an SQL database by implementing the
// database/sql/driver interfaces.
// The underlying column type can be a string, an integer (period of days since
// year 0), or a DATE.

// Scan parses some value. If the value holds a string, the AutoParse function is used.
// Otherwise, if the value holds an integer, it is treated as the period of days
// since year 0 value that represents a Date.
//
// This implements sql.Scanner https://golang.org/pkg/database/sql/#Scanner
func (d *Date) Scan(value interface{}) (err error) {
	if value == nil {
		return nil
	}

	return d.scanAny(value)
}

func (d *Date) scanAny(value interface{}) (err error) {
	err = nil
	switch v := value.(type) {
	case int64:
		*d = Date(v)
	case []byte:
		return d.scanString(string(v))
	case string:
		return d.scanString(v)
	case time.Time:
		*d = NewAt(v)
	default:
		err = fmt.Errorf("%T %+v is not a meaningful date", value, value)
	}

	return err
}

func (d *Date) scanString(value string) error {
	var err1 error
	*d, err1 = AutoParse(value)
	return err1
}

// Value converts the value for DB storage. It uses Valuer, which returns strings
// by default.
//
// This implements driver.Valuer https://golang.org/pkg/database/sql/driver/#Valuer
func (d Date) Value() (driver.Value, error) {
	return Valuer(d)
}

// Valuer is the pluggable implementation function for converting dates to driver.Value.
// It is initialised with ValueAsString.
var Valuer = ValueAsString

// ValueAsInt converts a date for DB storage using an integer.
func ValueAsInt(d Date) (driver.Value, error) {
	return int64(d), nil
}

// ValueAsString converts a date for DB storage using an string.
func ValueAsString(d Date) (driver.Value, error) {
	return d.String(), nil
}

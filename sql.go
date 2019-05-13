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
// The underlying column type can be an integer (period of days since the epoch),
// a string, or a DATE.

// Scan parses some value. It implements sql.Scanner,
// https://golang.org/pkg/database/sql/#Scanner
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
		*d = Date{PeriodOfDays(v)}
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

func (d *Date) scanString(value string) (err error) {
	n, err := strconv.ParseInt(value, 10, 64)
	if err == nil {
		*d = Date{PeriodOfDays(n)}
		return nil
	}
	*d, err = AutoParse(value)
	return err
}

// Value converts the value to an int64. It implements driver.Valuer,
// https://golang.org/pkg/database/sql/driver/#Valuer
func (d Date) Value() (driver.Value, error) {
	return int64(d.day), nil
}

//-------------------------------------------------------------------------------------------------

// DateString alters Date to make database storage use a string column, or
// a similar derived column such as SQL DATE. (Otherwise, Date is stored as
// an integer).
type DateString Date

// Date provides a simple fluent type conversion to the underlying type.
func (ds DateString) Date() Date {
	return Date(ds)
}

// DateString provides a simple fluent type conversion from the underlying type.
func (d Date) DateString() DateString {
	return DateString(d)
}

// Scan parses some value. It implements sql.Scanner,
// https://golang.org/pkg/database/sql/#Scanner
func (ds *DateString) Scan(value interface{}) (err error) {
	if value == nil {
		return nil
	}
	return (*Date)(ds).Scan(value)
}

// Value converts the value to an int64. It implements driver.Valuer,
// https://golang.org/pkg/database/sql/driver/#Valuer
func (ds DateString) Value() (driver.Value, error) {
	return ds.Date().String(), nil
}

//-------------------------------------------------------------------------------------------------

// DisableTextStorage reduces the Scan method so that only integers are handled.
// Normally, database types int64, []byte, string and time.Time are supported.
// When set true, only int64 is supported; this mode allows optimisation of SQL
// result processing and would only be used during development.
//
// Deprecated: this is no longer used.
var DisableTextStorage = false

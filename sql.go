package date

import (
	"database/sql/driver"
	"fmt"
)

// These methods allow Date and PeriodOfDays to be fields stored in an
// SQL database by implementing the database/sql/driver interfaces.
// The underlying column type is simply an integer.

// Scan parses some value. It implements sql.Scanner,
// https://golang.org/pkg/database/sql/#Scanner
func (d *Date) Scan(value interface{}) (err error) {
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

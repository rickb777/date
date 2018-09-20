package clock

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// Scan parses some value. It implements sql.Scanner,
// https://golang.org/pkg/database/sql/#Scanner
func (c *Clock) Scan(value interface{}) (err error) {
	if value == nil {
		return nil
	}

	return c.scanAny(value)
}

func (c *Clock) scanAny(value interface{}) (err error) {
	err = nil
	switch value.(type) {
	case int64:
		*c = Clock(value.(int64))
	case []byte:
		*c, err = Parse(string(value.([]byte)))
	case string:
		*c, err = Parse(value.(string))
	case time.Time:
		*c = NewAt(value.(time.Time))
	default:
		err = fmt.Errorf("%T %+v is not a meaningful clock", value, value)
	}
	return
}

// Value converts the value to an int64. It implements driver.Valuer,
// https://golang.org/pkg/database/sql/driver/#Valuer
func (c Clock) Value() (driver.Value, error) {

	return int64(c), nil
}

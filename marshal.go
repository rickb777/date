// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package date

import (
	"bytes"
	"errors"
)

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (d Date) MarshalBinary() ([]byte, error) {
	enc := []byte{
		byte(d.day >> 24),
		byte(d.day >> 16),
		byte(d.day >> 8),
		byte(d.day),
	}
	return enc, nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (d *Date) UnmarshalBinary(data []byte) error {
	if len(data) == 0 {
		return errors.New("Date.UnmarshalBinary: no data")
	}
	if len(data) != 4 {
		return errors.New("Date.UnmarshalBinary: invalid length")
	}

	d.day = PeriodOfDays(data[3]) | PeriodOfDays(data[2])<<8 | PeriodOfDays(data[1])<<16 | PeriodOfDays(data[0])<<24

	return nil
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (ds DateString) MarshalBinary() ([]byte, error) {
	return Date(ds).MarshalBinary()
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (ds *DateString) UnmarshalBinary(data []byte) error {
	return (*Date)(ds).UnmarshalBinary(data)
}

// MarshalJSON implements the json.Marshaler interface.
// The date is given in ISO 8601 extended format (e.g. "2006-01-02").
// If the year of the date falls outside the [0,9999] range, this format
// produces an expanded year representation with possibly extra year digits
// beyond the prescribed four-digit minimum and with a + or - sign prefix
// (e.g. , "+12345-06-07", "-0987-06-05").
// Note that the zero value is marshalled as a blank string, which allows
// "omitempty" to work.
func (d Date) MarshalJSON() ([]byte, error) {
	buf := &bytes.Buffer{}
	buf.Grow(14)
	buf.WriteByte('"')
	d.WriteTo(buf)
	buf.WriteByte('"')
	return buf.Bytes(), nil
}

// MarshalText implements the encoding.TextMarshaler interface.
// The date is given in ISO 8601 extended format (e.g. "2006-01-02").
// If the year of the date falls outside the [0,9999] range, this format
// produces an expanded year representation with possibly extra year digits
// beyond the prescribed four-digit minimum and with a + or - sign prefix
// (e.g. , "+12345-06-07", "-0987-06-05").
func (d Date) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// The date is expected to be in ISO 8601 extended format
// (e.g. "2006-01-02", "+12345-06-07", "-0987-06-05");
// the year must use at least 4 digits and if outside the [0,9999] range
// must be prefixed with a + or - sign.
// Note that the a blank string is unmarshalled as the zero value.
func (d *Date) UnmarshalText(data []byte) (err error) {
	if len(data) == 0 {
		return nil
	}
	u, err := ParseISO(string(data))
	if err == nil {
		d.day = u.day
	}
	return err
}

// MarshalJSON implements the json.Marshaler interface.
// The date is given in ISO 8601 extended format (e.g. "2006-01-02").
// If the year of the date falls outside the [0,9999] range, this format
// produces an expanded year representation with possibly extra year digits
// beyond the prescribed four-digit minimum and with a + or - sign prefix
// (e.g. , "+12345-06-07", "-0987-06-05").
// Note that the zero value is marshalled as a blank string, which allows
// "omitempty" to work.
func (ds DateString) MarshalJSON() ([]byte, error) {
	return Date(ds).MarshalJSON()
}

// MarshalText implements the encoding.TextMarshaler interface.
// The date is given in ISO 8601 extended format (e.g. "2006-01-02").
// If the year of the date falls outside the [0,9999] range, this format
// produces an expanded year representation with possibly extra year digits
// beyond the prescribed four-digit minimum and with a + or - sign prefix
// (e.g. , "+12345-06-07", "-0987-06-05").
func (ds DateString) MarshalText() ([]byte, error) {
	return Date(ds).MarshalText()
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// The date is expected to be in ISO 8601 extended format
// (e.g. "2006-01-02", "+12345-06-07", "-0987-06-05");
// the year must use at least 4 digits and if outside the [0,9999] range
// must be prefixed with a + or - sign.
func (ds *DateString) UnmarshalText(data []byte) (err error) {
	return (*Date)(ds).UnmarshalText(data)
}

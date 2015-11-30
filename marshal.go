// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package date

import (
	"fmt"
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

	d.day = int32(data[3]) | int32(data[2]) << 8 | int32(data[1]) << 16 | int32(data[0]) << 24
	//	d.decoded = time.Time{}

	return nil
}

// GobEncode implements the gob.GobEncoder interface.
func (d Date) GobEncode() ([]byte, error) {
	return d.MarshalBinary()
}

// GobDecode implements the gob.GobDecoder interface.
func (d *Date) GobDecode(data []byte) error {
	return d.UnmarshalBinary(data)
}

// MarshalJSON implements the json.Marshaler interface.
// The date is a quoted string in ISO 8601 extended format (e.g. "2006-01-02").
// If the year of the date falls outside the [0,9999] range, this format
// produces an expanded year representation with possibly extra year digits
// beyond the prescribed four-digit minimum and with a + or - sign prefix
// (e.g. , "+12345-06-07", "-0987-06-05").
func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.String() + `"`), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// The date is expected to be a quoted string in ISO 8601 extended format
// (e.g. "2006-01-02", "+12345-06-07", "-0987-06-05");
// the year must use at least 4 digits and if outside the [0,9999] range
// must be prefixed with a + or - sign.
func (d *Date) UnmarshalJSON(data []byte) (err error) {
	value := string(data)
	n := len(value)
	if n < 2 || value[0] != '"' || value[n - 1] != '"' {
		return fmt.Errorf("Date.UnmarshalJSON: missing double quotes (%s)", value)
	}
	u, err := ParseISO(value[1 : n - 1])
	if err != nil {
		return err
	}
	d.day = u.day
	//	d.decoded = time.Time{}
	return nil
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
func (d *Date) UnmarshalText(data []byte) (err error) {
	u, err := ParseISO(string(data))
	if err != nil {
		return err
	}
	d.day = u.day
	//	d.decoded = time.Time{}
	return nil
}

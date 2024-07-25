// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package date

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
)

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (d Date) MarshalBinary() (b []byte, err error) {
	if math.MaxInt == math.MaxInt32 {
		b = make([]byte, 4)
		binary.LittleEndian.PutUint32(b, uint32(d))
	} else {
		b = make([]byte, 8)
		binary.LittleEndian.PutUint64(b, uint64(d))
	}
	return b, nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (d *Date) UnmarshalBinary(data []byte) error {
	switch len(data) {
	case 0:
		return errors.New("Date.UnmarshalBinary: no data")
	case 4:
		*d = Date(binary.LittleEndian.Uint32(data))
	case 8:
		*d = Date(binary.LittleEndian.Uint64(data))
	default:
		return fmt.Errorf("Date.UnmarshalBinary: invalid length %d bytes", len(data))
	}
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
	if err == nil {
		*d = u
	}
	return err
}

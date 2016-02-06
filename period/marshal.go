// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package period

import (
	"fmt"
)

//// MarshalBinary implements the encoding.BinaryMarshaler interface.
//func (d Date) MarshalBinary() ([]byte, error) {
//	enc := []byte{
//		byte(d.day >> 24),
//		byte(d.day >> 16),
//		byte(d.day >> 8),
//		byte(d.day),
//	}
//	return enc, nil
//}
//
//// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
//func (d *Date) UnmarshalBinary(data []byte) error {
//	if len(data) == 0 {
//		return errors.New("Date.UnmarshalBinary: no data")
//	}
//	if len(data) != 4 {
//		return errors.New("Date.UnmarshalBinary: invalid length")
//	}
//
//	d.day = PeriodOfDays(data[3]) | PeriodOfDays(data[2])<<8 | PeriodOfDays(data[1])<<16 | PeriodOfDays(data[0])<<24
//
//	return nil
//}
//
//// GobEncode implements the gob.GobEncoder interface.
//func (d Date) GobEncode() ([]byte, error) {
//	return d.MarshalBinary()
//}
//
//// GobDecode implements the gob.GobDecoder interface.
//func (d *Date) GobDecode(data []byte) error {
//	return d.UnmarshalBinary(data)
//}

// MarshalJSON implements the json.Marshaler interface for Period.
func (period Period) MarshalJSON() ([]byte, error) {
	return []byte(`"` + period.String() + `"`), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface for Period.
func (period *Period) UnmarshalJSON(data []byte) error {
	n := len(data)
	if n < 2 || data[0] != '"' || data[n-1] != '"' {
		return fmt.Errorf("Period.UnmarshalJSON: missing double quotes (%s)", string(data))
	}
	return period.UnmarshalText(data[1 : n-1])
}

// MarshalText implements the encoding.TextMarshaler interface for Periods.
func (period Period) MarshalText() ([]byte, error) {
	return []byte(period.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for Periods.
func (period *Period) UnmarshalText(data []byte) (err error) {
	u, err := ParsePeriod(string(data))
	if err == nil {
		*period = u
	}
	return err
}

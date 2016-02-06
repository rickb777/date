// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package period

import "errors"

// MarshalBinary implements the encoding.BinaryMarshaler interface.
// This also provides support for gob encoding.
func (p Period) MarshalBinary() ([]byte, error) {
	enc := []byte{
		byte(p.years >> 24),
		byte(p.years >> 16),
		byte(p.years >> 8),
		byte(p.years),
		byte(p.months >> 24),
		byte(p.months >> 16),
		byte(p.months >> 8),
		byte(p.months),
		byte(p.days >> 24),
		byte(p.days >> 16),
		byte(p.days >> 8),
		byte(p.days),
	}
	return enc, nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
// This also provides support for gob encoding.
func (p *Period) UnmarshalBinary(data []byte) error {
	if len(data) == 0 {
		return errors.New("Date.UnmarshalBinary: no data")
	}
	if len(data) != 12 {
		return errors.New("Date.UnmarshalBinary: invalid length")
	}

	p.years = int32(data[3]) | int32(data[2])<<8 | int32(data[1])<<16 | int32(data[0])<<24
	p.months = int32(data[7]) | int32(data[6])<<8 | int32(data[5])<<16 | int32(data[4])<<24
	p.days = int32(data[11]) | int32(data[10])<<8 | int32(data[9])<<16 | int32(data[8])<<24

	return nil
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

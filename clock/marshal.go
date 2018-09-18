package clock

import (
	"errors"
)

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (c Clock) MarshalBinary() ([]byte, error) {
	enc := []byte{
		byte(c >> 24),
		byte(c >> 16),
		byte(c >> 8),
		byte(c),
	}
	return enc, nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (c *Clock) UnmarshalBinary(data []byte) error {
	if len(data) == 0 {
		return errors.New("Clock.UnmarshalBinary: no data")
	}
	if len(data) != 4 {
		return errors.New("Clock.UnmarshalBinary: invalid length")
	}

	*c = Clock(data[3]) | Clock(data[2])<<8 | Clock(data[1])<<16 | Clock(data[0])<<24
	return nil
}

// MarshalText implements the encoding.TextMarshaler interface.
func (c Clock) MarshalText() ([]byte, error) {
	return []byte(c.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (c *Clock) UnmarshalText(data []byte) (err error) {
	clock, err := Parse(string(data))
	if err == nil {
		*c = clock
	}
	return err
}

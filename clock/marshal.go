package clock

import (
	"encoding/binary"
	"errors"
	"fmt"
)

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (c Clock) MarshalBinary() (b []byte, err error) {
	b = make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(c))
	return b, nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (c *Clock) UnmarshalBinary(data []byte) error {
	switch len(data) {
	case 0:
		return errors.New("Clock.UnmarshalBinary: no data")
	case 8:
		*c = Clock(binary.LittleEndian.Uint64(data))
	default:
		return fmt.Errorf("Clock.UnmarshalBinary: invalid length %d bytes", len(data))
	}

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

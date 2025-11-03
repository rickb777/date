package date

import (
	"encoding/binary"
	"errors"
	"fmt"
)

// MarshalBinary implements the [encoding.BinaryMarshaler] interface.
func (d Date) MarshalBinary() ([]byte, error) {
	var b [8]byte
	binary.LittleEndian.PutUint64(b[:], uint64(d))
	return b[:], nil
}

// UnmarshalBinary implements the [encoding.BinaryUnmarshaler] interface.
func (d *Date) UnmarshalBinary(data []byte) error {
	switch len(data) {
	case 0:
		return errors.New("Date.UnmarshalBinary: no data")
	case 8:
		*d = Date(binary.LittleEndian.Uint64(data))
	default:
		return fmt.Errorf("Date.UnmarshalBinary: invalid length %d bytes", len(data))
	}
	return nil
}

// MarshalText implements the [encoding.TextMarshaler] interface.
// The date is given in ISO 8601 extended format (e.g. "2006-01-02").
// If the year of the date falls outside the [0,9999] range, this format
// produces an expanded year representation with possibly extra year digits
// beyond the prescribed four-digit minimum and with a + or - sign prefix
// (e.g. , "+12345-06-07", "-0987-06-05").
func (d Date) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

// UnmarshalText implements the [encoding.TextUnmarshaler] interface.
// The date is typically expected to be in ISO 8601 extended format
// (e.g. "2006-01-02", "+12345-06-07", "-0987-06-05");
// the year must use at least 4 digits and if outside the [0,9999] range
// must be prefixed with a + or - sign.
// In practice, all inputs handled by [AutoParse] are accepted.
func (d *Date) UnmarshalText(data []byte) (err error) {
	u, err := AutoParse(string(data))
	if err == nil {
		*d = u
	}
	return err
}

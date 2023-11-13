package clock

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestGobEncoding(t *testing.T) {
	var b bytes.Buffer
	encoder := gob.NewEncoder(&b)
	decoder := gob.NewDecoder(&b)
	cases := []Clock{
		New(-1, -1, -1, -1),
		New(-1, -1, -1, -1).AddDuration(-1),
		New(0, 0, 0, 0),
		New(12, 40, 40, 80),
		New(13, 55, 0, 20),
		New(16, 20, 0, 0),
		New(20, 60, 59, 59),
		New(20, 60, 59, 111).AddDuration(222333),
		New(24, 0, 0, 0),
		New(24, 0, 0, 1),
		New(24, 0, 0, 0).AddDuration(1),
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("%d %s", i, c), func(t *testing.T) {
			var clock Clock
			err := encoder.Encode(&c)
			if err != nil {
				t.Errorf("Gob(%v) encode error %v", c, err)
			} else {
				err = decoder.Decode(&clock)
				if err != nil {
					t.Errorf("Gob(%v) decode error %v", c, err)
				} else if clock != c {
					t.Errorf("Gob(%v) decode got %v", c, clock)
				}
			}
		})
	}
}

func TestJSONMarshalling(t *testing.T) {
	cases := []struct {
		value Clock
		want  string
	}{
		{New(-1, -1, -1, -1), `"22:58:58.999"`},
		{New(0, 0, 0, 0), `"00:00:00.000"`},
		{New(12, 40, 40, 80), `"12:40:40.080"`},
		{New(13, 55, 0, 20), `"13:55:00.020"`},
		{New(16, 20, 0, 0), `"16:20:00.000"`},
		{New(20, 60, 59, 9), `"21:00:59.009"`},
		{New(20, 60, 59, 7).AddDuration(9), `"21:00:59.007000009"`},
		{New(24, 0, 0, 0), `"24:00:00.000"`},
		{New(24, 0, 0, 1), `"00:00:00.001"`},
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("%d %s", i, c.want), func(t *testing.T) {
			bb, err := json.Marshal(c.value)
			if err != nil {
				t.Errorf("JSON(%v) marshal error %v", c, err)
			} else if string(bb) != c.want {
				t.Errorf("JSON(%v) == %v, want %v", c.value, string(bb), c.want)
			}
		})
	}
}

func TestJSONUnmarshalling(t *testing.T) {
	cases := []struct {
		values []string
		want   Clock
	}{
		{[]string{`"22:58:58.999"`, `"22:58:58.999000000"`, `"10:58:58.999pm"`}, New(-1, -1, -1, -1)},
		{[]string{`"00:00:00.000"`, `"00:00:00.000000000"`, `"00:00:00.000AM"`}, New(0, 0, 0, 0)},
		{[]string{`"12:40:40.080"`, `"12:40:40.080000000"`, `"12:40:40.080PM"`}, New(12, 40, 40, 80)},
		{[]string{`"13:55:00.020"`, `"13:55:00.020000000"`, `"01:55:00.020PM"`}, New(13, 55, 0, 20)},
		{[]string{`"16:20:00.000"`, `"16:20:00.000000000"`, `"04:20:00.000pm"`}, New(16, 20, 0, 0)},
		{[]string{`"21:00:59.059"`, `"21:00:59.059000000"`, `"09:00:59.059PM"`}, New(20, 60, 59, 59)},
		{[]string{`"24:00:00.000"`, `"24:00:00.000000000"`, `"00:00:00.000am"`}, New(24, 0, 0, 0)},
		{[]string{`"00:00:00.001"`, `"00:00:00.001000000"`, `"00:00:00.001AM"`}, New(24, 0, 0, 1)},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("%d %s", i, c.want), func(t *testing.T) {
			for _, v := range c.values {
				var clock Clock
				err := json.Unmarshal([]byte(v), &clock)
				if err != nil {
					t.Errorf("JSON(%v) unmarshal error %v", v, err)
				} else if c.want.Mod24() != clock.Mod24() {
					t.Errorf("JSON(%v) == %v, want %v", v, clock, c.want)
				}
			}
		})
	}
}

func TestBinaryMarshalling(t *testing.T) {
	cases := []Clock{
		New(-1, -1, -1, -1),
		New(-1, -1, -1, -1).AddDuration(-1),
		New(0, 0, 0, 0),
		New(12, 40, 40, 80),
		New(13, 55, 0, 20),
		New(16, 20, 0, 0),
		New(20, 60, 59, 9),
		New(20, 60, 59, 7).AddDuration(9),
		New(24, 0, 0, 0),
		New(24, 0, 0, 1),
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("%d %s", i, c), func(t *testing.T) {
			bb, err := c.MarshalBinary()
			if err != nil {
				t.Errorf("Binary(%v) marshal error %v", c, err)
			} else {
				var clock Clock
				err = clock.UnmarshalBinary(bb)
				if err != nil {
					t.Errorf("Binary(% v) unmarshal error %v", c, err)
				} else if clock.Mod24() != c.Mod24() {
					t.Errorf("Binary(%v) unmarshal got %v", c, clock)
				}
			}
		})
	}
}

func TestBinaryUnmarshallingErrors(t *testing.T) {
	var c Clock
	err1 := c.UnmarshalBinary([]byte{})
	if err1 == nil {
		t.Errorf("unmarshal no empty data error")
	}

	err2 := c.UnmarshalBinary([]byte("12345"))
	if err2 == nil {
		t.Errorf("unmarshal no wrong length error")
	}
}

func TestInvalidClockText(t *testing.T) {
	cases := []struct {
		value string
		want  string
	}{
		{value: `not-a-clock`, want: `clock.Clock: cannot parse "not-a-clock"`},
		{value: `00:50:100.0`, want: `clock.Clock: cannot parse "00:50:100.0"`},
		{value: `24:00:00.0pM`, want: `clock.Clock: cannot parse "24:00:00.0pM"`},
	}
	for _, c := range cases {
		var clock Clock
		err := clock.UnmarshalText([]byte(c.value))
		if err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("InvalidText(%v) == %v, want %v", c.value, err, c.want)
		}
	}
}

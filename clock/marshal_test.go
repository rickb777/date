package clock

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"strings"
	"testing"
)

func TestGobEncoding(t *testing.T) {
	var b bytes.Buffer
	encoder := gob.NewEncoder(&b)
	decoder := gob.NewDecoder(&b)
	cases := []Clock{
		New(-1, -1, -1, -1),
		New(0, 0, 0, 0),
		New(12, 40, 40, 80),
		New(13, 55, 0, 20),
		New(16, 20, 0, 0),
		New(20, 60, 59, 59),
		New(24, 0, 0, 0),
		New(24, 0, 0, 1),
	}
	for _, c := range cases {
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
		{New(20, 60, 59, 59), `"21:00:59.059"`},
		{New(24, 0, 0, 0), `"24:00:00.000"`},
		{New(24, 0, 0, 1), `"00:00:00.001"`},
	}
	for _, c := range cases {
		bb, err := json.Marshal(c.value)
		if err != nil {
			t.Errorf("JSON(%v) marshal error %v", c, err)
		} else if string(bb) != c.want {
			t.Errorf("JSON(%v) == %v, want %v", c.value, string(bb), c.want)
		}
	}
}

func TestJSONUnmarshalling(t *testing.T) {
	cases := []struct {
		values []string
		want   Clock
	}{
		{[]string{`"22:58:58.999"`, `"10:58:58.999pm"`}, New(-1, -1, -1, -1)},
		{[]string{`"00:00:00.000"`, `"00:00:00.000AM"`}, New(0, 0, 0, 0)},
		{[]string{`"12:40:40.080"`, `"12:40:40.080PM"`}, New(12, 40, 40, 80)},
		{[]string{`"13:55:00.020"`, `"01:55:00.020PM"`}, New(13, 55, 0, 20)},
		{[]string{`"16:20:00.000"`, `"04:20:00.000pm"`}, New(16, 20, 0, 0)},
		{[]string{`"21:00:59.059"`, `"09:00:59.059PM"`}, New(20, 60, 59, 59)},
		{[]string{`"24:00:00.000"`, `"00:00:00.000am"`}, New(24, 0, 0, 0)},
		{[]string{`"00:00:00.001"`, `"00:00:00.001AM"`}, New(24, 0, 0, 1)},
	}

	for _, c := range cases {
		for _, v := range c.values {
			var clock Clock
			err := json.Unmarshal([]byte(v), &clock)
			if err != nil {
				t.Errorf("JSON(%v) unmarshal error %v", v, err)
			} else if c.want.Mod24() != clock.Mod24() {
				t.Errorf("JSON(%v) == %v, want %v", v, clock, c.want)
			}
		}
	}
}

func TestBinaryMarshalling(t *testing.T) {
	cases := []Clock{
		New(-1, -1, -1, -1),
		New(0, 0, 0, 0),
		New(12, 40, 40, 80),
		New(13, 55, 0, 20),
		New(16, 20, 0, 0),
		New(20, 60, 59, 59),
		New(24, 0, 0, 0),
		New(24, 0, 0, 1),
	}
	for _, c := range cases {
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
		{`not-a-clock`, `clock.Clock: cannot parse not-a-clock`},
		{`00:50:100.0`, `clock.Clock: cannot parse 00:50:100.0`},
		{`24:00:00.0pM`, `clock.Clock: cannot parse 24:00:00.0pM: strconv.Atoi: parsing "0pM": invalid syntax`},
	}
	for _, c := range cases {
		var clock Clock
		err := clock.UnmarshalText([]byte(c.value))
		if err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("InvalidText(%v) == %v, want %v", c.value, err, c.want)
		}
	}
}

// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package period

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"testing"
)

func TestGobEncoding(t *testing.T) {
	var b bytes.Buffer
	encoder := gob.NewEncoder(&b)
	decoder := gob.NewDecoder(&b)
	cases := []string{
		"P0D",
		"P1D",
		"P1W",
		"P1M",
		"P1Y",
		"PT1H",
		"PT1M",
		"PT1S",
		"P2Y3M4W5D",
		"-P2Y3M4W5D",
		"P2Y3M4W5DT1H7M9S",
		"-P2Y3M4W5DT1H7M9S",
	}
	for _, c := range cases {
		period := MustParse(c)
		var p Period
		err := encoder.Encode(&period)
		if err != nil {
			t.Errorf("Gob(%v) encode error %v", c, err)
		} else {
			err = decoder.Decode(&p)
			if err != nil {
				t.Errorf("Gob(%v) decode error %v", c, err)
			} else if p != period {
				t.Errorf("Gob(%v) decode got %v", c, p)
			}
		}
	}
}

func TestPeriodJSONMarshalling(t *testing.T) {
	cases := []struct {
		value Period
		want  string
	}{
		{New(-1111, -4, -3, -11, -59, -59), `"-P1111Y4M3DT11H59M59S"`},
		{New(-1, -10, -31, -5, -4, -20), `"-P1Y10M31DT5H4M20S"`},
		{New(0, 0, 0, 0, 0, 0), `"P0D"`},
		{New(0, 0, 0, 0, 0, 1), `"PT1S"`},
		{New(0, 0, 0, 0, 1, 0), `"PT1M"`},
		{New(0, 0, 0, 1, 0, 0), `"PT1H"`},
		{New(0, 0, 1, 0, 0, 0), `"P1D"`},
		{New(0, 1, 0, 0, 0, 0), `"P1M"`},
		{New(1, 0, 0, 0, 0, 0), `"P1Y"`},
	}
	for _, c := range cases {
		var p Period
		bytes, err := json.Marshal(c.value)
		if err != nil {
			t.Errorf("JSON(%v) marshal error %v", c, err)
		} else if string(bytes) != c.want {
			t.Errorf("JSON(%v) == %v, want %v", c.value, string(bytes), c.want)
		} else {
			err = json.Unmarshal(bytes, &p)
			if err != nil {
				t.Errorf("JSON(%v) unmarshal error %v", c.value, err)
			} else if p != c.value {
				t.Errorf("JSON(%v) unmarshal got %v", c.value, p)
			}
		}
	}
}

func TestPeriodTextMarshalling(t *testing.T) {
	cases := []struct {
		value Period
		want  string
	}{
		{New(-1111, -4, -3, -11, -59, -59), "-P1111Y4M3DT11H59M59S"},
		{New(-1, -9, -31, -5, -4, -20), "-P1Y9M31DT5H4M20S"},
		{New(0, 0, 0, 0, 0, 0), "P0D"},
		{New(0, 0, 0, 0, 0, 1), "PT1S"},
		{New(0, 0, 0, 0, 1, 0), "PT1M"},
		{New(0, 0, 0, 1, 0, 0), "PT1H"},
		{New(0, 0, 1, 0, 0, 0), "P1D"},
		{New(0, 1, 0, 0, 0, 0), "P1M"},
		{New(1, 0, 0, 0, 0, 0), "P1Y"},
	}
	for _, c := range cases {
		var p Period
		bytes, err := c.value.MarshalText()
		if err != nil {
			t.Errorf("Text(%v) marshal error %v", c, err)
		} else if string(bytes) != c.want {
			t.Errorf("Text(%v) == %v, want %v", c.value, string(bytes), c.want)
		} else {
			err = p.UnmarshalText(bytes)
			if err != nil {
				t.Errorf("Text(%v) unmarshal error %v", c.value, err)
			} else if p != c.value {
				t.Errorf("Text(%v) unmarshal got %v", c.value, p)
			}
		}
	}
}

func TestInvalidPeriodText(t *testing.T) {
	cases := []struct {
		value string
		want  string
	}{
		{``, `cannot parse a blank string as a period`},
		{`not-a-period`, `expected 'P' period mark at the start: not-a-period`},
		{`P000`, `expected 'Y', 'M', 'W', 'D', 'H', 'M', or 'S' marker: P000`},
	}
	for _, c := range cases {
		var p Period
		err := p.UnmarshalText([]byte(c.value))
		if err == nil || err.Error() != c.want {
			t.Errorf("InvalidText(%v) == %v, want %v", c.value, err, c.want)
		}
	}
}

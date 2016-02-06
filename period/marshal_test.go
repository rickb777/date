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
		"P2Y3M4W5D",
		"-P2Y3M4W5D",
	}
	for _, c := range cases {
		period := MustParsePeriod(c)
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
		{NewPeriod(-11111, -123, -3), `"-P11111Y123M3D"`},
		{NewPeriod(-1, -12, -31), `"-P1Y12M31D"`},
		{NewPeriod(0, 0, 0), `"P0D"`},
		{NewPeriod(0, 0, 1), `"P1D"`},
		{NewPeriod(0, 1, 0), `"P1M"`},
		{NewPeriod(1, 0, 0), `"P1Y"`},
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
		{NewPeriod(-11111, -123, -3), "-P11111Y123M3D"},
		{NewPeriod(-1, -12, -31), "-P1Y12M31D"},
		{NewPeriod(0, 0, 0), "P0D"},
		{NewPeriod(0, 0, 1), "P1D"},
		{NewPeriod(0, 1, 0), "P1M"},
		{NewPeriod(1, 0, 0), "P1Y"},
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
		{``, `Cannot parse a blank string as a period.`},
		{`not-a-period`, `Expected 'P' period mark at the start: not-a-period`},
		{`P000`, `Expected 'Y', 'M', 'W' or 'D' marker: P000`},
	}
	for _, c := range cases {
		var p Period
		err := p.UnmarshalText([]byte(c.value))
		if err == nil || err.Error() != c.want {
			t.Errorf("InvalidText(%v) == %v, want %v", c.value, err, c.want)
		}
	}
}

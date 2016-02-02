// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package date

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"testing"
	"time"
)

func TestGobEncoding(t *testing.T) {
	var b bytes.Buffer
	encoder := gob.NewEncoder(&b)
	decoder := gob.NewDecoder(&b)
	cases := []Date{
		New(-11111, time.February, 3),
		New(-1, time.December, 31),
		New(0, time.January, 1),
		New(1, time.January, 1),
		New(1970, time.January, 1),
		New(2012, time.June, 25),
		New(12345, time.June, 7),
	}
	for _, c := range cases {
		var d Date
		err := encoder.Encode(&c)
		if err != nil {
			t.Errorf("Gob(%v) encode error %v", c, err)
		} else {
			err = decoder.Decode(&d)
			if err != nil {
				t.Errorf("Gob(%v) decode error %v", c, err)
			} else if d != c {
				t.Errorf("Gob(%v) decode got %v", c, d)
			}
		}
	}
}

func TestInvalidGob(t *testing.T) {
	cases := []struct {
		bytes []byte
		want  string
	}{
		{[]byte{}, "Date.UnmarshalBinary: no data"},
		{[]byte{1, 2, 3}, "Date.UnmarshalBinary: invalid length"},
	}
	for _, c := range cases {
		var ignored Date
		err := ignored.GobDecode(c.bytes)
		if err == nil || err.Error() != c.want {
			t.Errorf("InvalidGobDecode(%v) == %v, want %v", c.bytes, err, c.want)
		}
		err = ignored.UnmarshalBinary(c.bytes)
		if err == nil || err.Error() != c.want {
			t.Errorf("InvalidUnmarshalBinary(%v) == %v, want %v", c.bytes, err, c.want)
		}
	}
}

func TestDateJSONMarshalling(t *testing.T) {
	cases := []struct {
		value Date
		want  string
	}{
		{New(-11111, time.February, 3), `"-11111-02-03"`},
		{New(-1, time.December, 31), `"-0001-12-31"`},
		{New(0, time.January, 1), `"0000-01-01"`},
		{New(1, time.January, 1), `"0001-01-01"`},
		{New(1970, time.January, 1), `"1970-01-01"`},
		{New(2012, time.June, 25), `"2012-06-25"`},
		{New(12345, time.June, 7), `"+12345-06-07"`},
	}
	for _, c := range cases {
		var d Date
		bytes, err := json.Marshal(c.value)
		if err != nil {
			t.Errorf("JSON(%v) marshal error %v", c, err)
		} else if string(bytes) != c.want {
			t.Errorf("JSON(%v) == %v, want %v", c.value, string(bytes), c.want)
		} else {
			err = json.Unmarshal(bytes, &d)
			if err != nil {
				t.Errorf("JSON(%v) unmarshal error %v", c.value, err)
			} else if d != c.value {
				t.Errorf("JSON(%v) unmarshal got %v", c.value, d)
			}
		}
	}
}

func TestInvalidDateJSON(t *testing.T) {
	cases := []struct {
		value string
		want  string
	}{
		{`"not-a-date"`, `Date.ParseISO: cannot parse not-a-date: incorrect syntax`},
		{`2015-08-15"`, `Date.UnmarshalJSON: missing double quotes (2015-08-15")`},
		{`"2015-08-15`, `Date.UnmarshalJSON: missing double quotes ("2015-08-15)`},
		{`"215-08-15"`, `Date.ParseISO: cannot parse 215-08-15: invalid year`},
	}
	for _, c := range cases {
		var d Date
		err := d.UnmarshalJSON([]byte(c.value))
		if err == nil || err.Error() != c.want {
			t.Errorf("InvalidJSON(%v) == %v, want %v", c.value, err, c.want)
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

func TestInvalidPeriodJSON(t *testing.T) {
	cases := []struct {
		value string
		want  string
	}{
		{`""`, `Cannot parse a blank string as a period.`},
		{`"not-a-period"`, `Expected 'P' period mark at the start: not-a-period`},
		{`"P000"`, `Expected 'Y', 'M', 'W' or 'D' marker: P000`},
	}
	for _, c := range cases {
		var p Period
		err := p.UnmarshalJSON([]byte(c.value))
		if err == nil || err.Error() != c.want {
			t.Errorf("InvalidJSON(%v) == %v, want %v", c.value, err, c.want)
		}
	}
}

func TestDateTextMarshalling(t *testing.T) {
	cases := []struct {
		value Date
		want  string
	}{
		{New(-11111, time.February, 3), "-11111-02-03"},
		{New(-1, time.December, 31), "-0001-12-31"},
		{New(0, time.January, 1), "0000-01-01"},
		{New(1, time.January, 1), "0001-01-01"},
		{New(1970, time.January, 1), "1970-01-01"},
		{New(2012, time.June, 25), "2012-06-25"},
		{New(12345, time.June, 7), "+12345-06-07"},
	}
	for _, c := range cases {
		var d Date
		bytes, err := c.value.MarshalText()
		if err != nil {
			t.Errorf("Text(%v) marshal error %v", c, err)
		} else if string(bytes) != c.want {
			t.Errorf("Text(%v) == %v, want %v", c.value, string(bytes), c.want)
		} else {
			err = d.UnmarshalText(bytes)
			if err != nil {
				t.Errorf("Text(%v) unmarshal error %v", c.value, err)
			} else if d != c.value {
				t.Errorf("Text(%v) unmarshal got %v", c.value, d)
			}
		}
	}
}

func TestInvalidDateText(t *testing.T) {
	cases := []struct {
		value string
		want  string
	}{
		{`not-a-date`, `Date.ParseISO: cannot parse not-a-date: incorrect syntax`},
		{`215-08-15`, `Date.ParseISO: cannot parse 215-08-15: invalid year`},
	}
	for _, c := range cases {
		var d Date
		err := d.UnmarshalText([]byte(c.value))
		if err == nil || err.Error() != c.want {
			t.Errorf("InvalidText(%v) == %v, want %v", c.value, err, c.want)
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

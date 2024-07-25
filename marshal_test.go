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

func TestDate_gob_Encode_round_tripe(t *testing.T) {
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
		var b bytes.Buffer
		encoder := gob.NewEncoder(&b)
		decoder := gob.NewDecoder(&b)

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

func TestDate_MarshalJSON_round_trip(t *testing.T) {
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
		bb1, err := json.Marshal(c.value)
		if err != nil {
			t.Errorf("JSON(%v) marshal error %v", c, err)
		} else if string(bb1) != c.want {
			t.Errorf("JSON(%v) == %v, want %v", c.value, string(bb1), c.want)
		} else {
			err = json.Unmarshal(bb1, &d)
			if err != nil {
				t.Errorf("JSON(%v) unmarshal error %v", c.value, err)
			} else if d != c.value {
				t.Errorf("JSON(%v) unmarshal got %v", c.value, d)
			}
		}
	}
}

func TestDate_MarshalText_round_trip(t *testing.T) {
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
		bb1, err := c.value.MarshalText()
		if err != nil {
			t.Errorf("Text(%v) marshal error %v", c, err)
		} else if string(bb1) != c.want {
			t.Errorf("Text(%v) == %q, want %q", c.value, string(bb1), c.want)
		} else {
			err = d.UnmarshalText(bb1)
			if err != nil {
				t.Errorf("Text(%v) unmarshal error %v", c.value, err)
			} else if d != c.value {
				t.Errorf("Text(%v) unmarshal got %v", c.value, d)
			}
		}
	}
}

func TestDate_MarshalBinary_round_trip(t *testing.T) {
	cases := []struct {
		value Date
	}{
		{New(-11111, time.February, 3)},
		{New(-1, time.December, 31)},
		{New(0, time.January, 1)},
		{New(1, time.January, 1)},
		{New(1970, time.January, 1)},
		{New(2012, time.June, 25)},
		{New(12345, time.June, 7)},
	}
	for _, c := range cases {
		bb1, err := c.value.MarshalBinary()
		if err != nil {
			t.Errorf("Binary(%v) marshal error %v", c, err)
		} else {
			var d Date
			err = d.UnmarshalBinary(bb1)
			if err != nil {
				t.Errorf("Binary(%v) unmarshal error %v", c.value, err)
			} else if d != c.value {
				t.Errorf("Binary(%v) unmarshal got %v", c.value, d)
			}
		}
	}
}

func TestDate_UnmarshalBinary_errors(t *testing.T) {
	var d Date
	err1 := d.UnmarshalBinary([]byte{})
	if err1 == nil {
		t.Errorf("unmarshal no empty data error")
	}

	err2 := d.UnmarshalBinary([]byte("12345"))
	if err2 == nil {
		t.Errorf("unmarshal no wrong length error")
	}
}

func TestDate_UnmarshalText_invalid_date_text(t *testing.T) {
	cases := []struct {
		value string
		want  string
	}{
		{`not-a-date`, "date.ParseISO: cannot parse \"not-a-date\": year has wrong length\nmonth has wrong length\nday has wrong length"},
		{`foot-of-og`, "date.ParseISO: cannot parse \"foot-of-og\": invalid year\ninvalid month\ninvalid day"},
		{`215-08-15`, `date.ParseISO: cannot parse "215-08-15": invalid year`},
		{``, `date.ParseISO: cannot parse "": invalid year`},
	}
	for _, c := range cases {
		var d Date
		err := d.UnmarshalText([]byte(c.value))
		if err == nil {
			t.Errorf("InvalidText(%v) want %v", c.value, c.want)
		}
	}
}

func TestDate_UnmarshalText_empty(t *testing.T) {
	cases := []struct {
		value    string
		expected Date
	}{
		{``, New(1, time.January, 1)},
	}
	for _, c := range cases {
		var d Date
		err := d.UnmarshalText([]byte(c.value))
		if err != nil {
			t.Errorf("Errored parsing date (%s): %s", c.value, err.Error())
		} else if c.expected != d {
			t.Errorf("Expected date (%s) doesn't match output: %v", c.expected, d)
		}
	}
}

func TestDate_UnmarshalTime_invalid_time_text(t *testing.T) {
	cases := []struct {
		value string
		want  string
	}{
		{`not-a-date`, "parsing time \"not-a-date\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"not-a-date\" as \"2006\""},
		{`foot-of-og`, "parsing time \"foot-of-og\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"foot-of-og\" as \"2006\""},
		{`215-08-15`, `parsing time "215-08-15" as "2006-01-02T15:04:05Z07:00": cannot parse "215-08-15" as "2006"`},
		{``, `parsing time "" as "2006-01-02T15:04:05Z07:00": cannot parse "" as "2006"`},
	}
	for _, c := range cases {
		var d time.Time
		err := d.UnmarshalText([]byte(c.value))
		if err == nil {
			t.Errorf("InvalidText(%v) want %v", c.value, c.want)
		} else if err.Error() != c.want {
			t.Errorf("Wrong error message(%s) want: %s", err.Error(), c.want)
		}

	}
}

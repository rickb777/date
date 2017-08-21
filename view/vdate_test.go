// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package view

import (
	"encoding/json"
	"github.com/rickb777/date"
	"testing"
	"time"
)

func TestBasicFormatting(t *testing.T) {
	d := NewVDate(date.New(2016, 2, 7))
	is(t, d.String(), "2016-02-07")
	is(t, d.Format(), "07/02/2016")
	is(t, d.WithFormat(MDYFormat).Format(), "02/07/2016")
	is(t, d.Mon(), "Sun")
	is(t, d.Monday(), "Sunday")
	is(t, d.Day2(), "7")
	is(t, d.Day02(), "07")
	is(t, d.Day2nd(), "7th")
	is(t, d.Month1(), "2")
	is(t, d.Month01(), "02")
	is(t, d.Jan(), "Feb")
	is(t, d.January(), "February")
	is(t, d.Year(), "2016")
}

func TestNext(t *testing.T) {
	d := NewVDate(date.New(2016, 2, 7))
	is(t, d.Next().Day().String(), "2016-02-08")
	is(t, d.Next().Week().String(), "2016-02-14")
	is(t, d.Next().Month().String(), "2016-03-07")
	is(t, d.Next().Year().String(), "2017-02-07")
}

func TestPrevious(t *testing.T) {
	d := NewVDate(date.New(2016, 2, 7))
	is(t, d.Previous().Day().String(), "2016-02-06")
	is(t, d.Previous().Week().String(), "2016-01-31")
	is(t, d.Previous().Month().String(), "2016-01-07")
	is(t, d.Previous().Year().String(), "2015-02-07")
}

func is(t *testing.T, s1, s2 string) {
	if s1 != s2 {
		t.Errorf("%s != %s", s1, s2)
	}
}

func TestJSONMarshalling(t *testing.T) {
	cases := []struct {
		value VDate
		want  string
	}{
		{NewVDate(date.New(-1, time.December, 31)), `"-0001-12-31"`},
		{NewVDate(date.New(2012, time.June, 25)), `"2012-06-25"`},
		{NewVDate(date.New(12345, time.June, 7)), `"+12345-06-07"`},
	}
	for _, c := range cases {
		var d VDate
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
				t.Errorf("JSON(%#v) unmarshal got %#v", c.value, d)
			}
		}
	}
}

func TestTextMarshalling(t *testing.T) {
	cases := []struct {
		value VDate
		want  string
	}{
		{NewVDate(date.New(-1, time.December, 31)), "-0001-12-31"},
		{NewVDate(date.New(2012, time.June, 25)), "2012-06-25"},
		{NewVDate(date.New(12345, time.June, 7)), "+12345-06-07"},
	}
	for _, c := range cases {
		var d VDate
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
				t.Errorf("Text(%#v) unmarshal got %#v", c.value, d)
			}
		}
	}
}

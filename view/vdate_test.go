package view

import (
	"github.com/rickb777/date"
	"testing"
)

func TestBasicFormatting(t *testing.T) {
	d := NewVDate(date.New(2016, 2, 7))
	is(t, d.String(), "2016-02-07")
	is(t, d.Web(), "07/02/2016")
	is(t, d.Mon(), "Sun")
	is(t, d.Monday(), "Sunday")
	is(t, d.Day2(), "7")
	is(t, d.Day02(), "07")
	is(t, d.Day02nd(), "7th")
	is(t, d.Month1(), "2")
	is(t, d.Month01(), "02")
	is(t, d.MonthJan(), "Feb")
	is(t, d.MonthJanuary(), "February")
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
		t.Error("%s != %s", s1, s2)
	}
}

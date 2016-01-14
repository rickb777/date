// This package provides a fluent API for formatting dates as strings. This is useful in view-models
// especially when using Go temapltes.
package view

import (
	"github.com/rickb777/date"
	"r3/roster3/util/r3date"
)

const WebDateFormat = "02/01/2006"
const SqlDateFormat = "2006-01-02"

type VDate struct {
	d date.Date
}

func NewVDate(d date.Date) VDate {
	return VDate{d}
}

func (d VDate) Next() VDateDelta {
	return VDateDelta{d.d, 1}
}

func (d VDate) Previous() VDateDelta {
	return VDateDelta{d.d, -1}
}

func (d VDate) String() string {
	return d.d.Format(r3date.SqlDateFormat)
}

func (d VDate) Web() string {
	return d.d.Format(r3date.WebDateFormat)
}

func (d VDate) Mon() string {
	return d.d.Format("Mon")
}

func (d VDate) Monday() string {
	return d.d.Format("Monday")
}

func (d VDate) Day2() string {
	return d.d.Format("2")
}

func (d VDate) Day02() string {
	return d.d.Format("02")
}

func (d VDate) Day02nd() string {
	return d.d.Format("2nd")
}

func (d VDate) Month1() string {
	return d.d.Format("1")
}

func (d VDate) Month01() string {
	return d.d.Format("01")
}

func (d VDate) MonthJan() string {
	return d.d.Format("Jan")
}

func (d VDate) MonthJanuary() string {
	return d.d.Format("January")
}

func (d VDate) Year() string {
	return d.d.Format("2006")
}

//-------------------------------------------------------------------------------------------------

type VDateDelta struct {
	d    date.Date
	sign date.PeriodOfDays
}

func (dd VDateDelta) Day() VDate {
	return VDate{dd.d.Add(dd.sign)}
}

func (dd VDateDelta) Week() VDate {
	return VDate{dd.d.Add(dd.sign * 7)}
}

func (dd VDateDelta) Month() VDate {
	return VDate{dd.d.AddDate(0, int(dd.sign), 0)}
}

func (dd VDateDelta) Year() VDate {
	return VDate{dd.d.AddDate(int(dd.sign), 0, 0)}
}

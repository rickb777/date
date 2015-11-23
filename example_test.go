// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package date

import (
	"fmt"
	"time"
)

func ExampleMax() {
	d := Max()
	fmt.Println(d)
	// Output: +5881580-07-11
}

func ExampleMin() {
	d := Min()
	fmt.Println(d)
	// Output: -5877641-06-23
}

func ExampleNew() {
	d := New(9999, time.December, 31)
	fmt.Printf("The world ends on %s\n", d)
	// Output: The world ends on 9999-12-31
}

func ExampleParse() {
	// longForm shows by example how the reference date would be
	// represented in the desired layout.
	const longForm = "Mon, January 2, 2006"
	d, _ := Parse(longForm, "Tue, February 3, 2013")
	fmt.Println(d)

	// shortForm is another way the reference date would be represented
	// in the desired layout.
	const shortForm = "2006-Jan-02"
	d, _ = Parse(shortForm, "2013-Feb-03")
	fmt.Println(d)

	// Output:
	// 2013-02-03
	// 2013-02-03
}

func ExampleParseISO() {
	d, _ := ParseISO("+12345-06-07")
	year, month, day := d.Date()
	fmt.Println(year)
	fmt.Println(month)
	fmt.Println(day)
	// Output:
	// 12345
	// June
	// 7
}

func ExampleDate_AddDate() {
	d := New(1000, time.January, 1)
	// Months and days do not need to be constrained to [1,12] and [1,365].
	u := d.AddDate(0, 14, -1)
	fmt.Println(u)
	// Output: 1001-02-28
}

func ExampleDate_Format() {
	// layout shows by example how the reference time should be represented.
	const layout = "Jan 2, 2006"
	d := New(2009, time.November, 10)
	fmt.Println(d.Format(layout))
	// Output: Nov 10, 2009
}

func ExampleDate_FormatISO() {
	// According to legend, Rome was founded on April 21, 753 BC.
	// Note that with astronomical year numbering, 753 BC becomes -752
	// because 1 BC is actually year 0.
	d := New(-752, time.April, 21)
	fmt.Println(d.FormatISO(5))
	// Output: -00752-04-21
}

// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This tool prints equivalences between the string representation and the internal numerical
// representation for dates and clocks.
package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/rickb777/date"
	"github.com/rickb777/date/clock"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func usage() {
	fmt.Printf("Usage: %s [-t] number | date | time\n\n", os.Args[0])
	fmt.Printf(" -t:    terse output\n")
	fmt.Printf(" date:  [+-]yyyy/mm/dd | yyyy.mm.dd | dd/mm/yyyy | dd.mm.yyyy\n")
	fmt.Printf(" time:  e.g. 11:15:20 | 2:45pm | 1:15:10.101\n")
	os.Exit(0)
}

var titled = false
var terse = false
var success = false
var printer = message.NewPrinter(language.English)

func sprintf(num interface{}) string {
	if terse {
		return fmt.Sprintf("%d", num)
	} else {
		return printer.Sprintf("%d", num)
	}
}

func title() {
	if !terse && !titled {
		titled = true
		fmt.Printf("%-15s %-15s %-15s %s\n", "input", "number", "clock", "date")
		fmt.Printf("%-15s %-15s %-15s %s\n", "-----", "------", "-----", "----")
	}
}

func printArg(arg string) {

	i, err := strconv.ParseInt(arg, 10, 64)
	if err == nil {
		title()
		d := date.NewOfDays(date.PeriodOfDays(i))
		c := clock.Clock(i)
		fmt.Printf("%-15s %-15s %-15s %-12s %s\n", arg, sprintf(i), c, d, d.Weekday())
		success = true
		return
	}

	d, e1 := date.AutoParse(arg)
	if e1 == nil {
		title()
		fmt.Printf("%-15s %-15s %15s %-12s %s\n", arg, sprintf(d.DaysSinceEpoch()), "", d, d.Weekday())
		success = true
	}

	c, err := clock.Parse(arg)
	if err == nil {
		title()
		fmt.Printf("%-15s %-15s %s\n", arg, sprintf(c), c)
		success = true
	}
}

func main() {
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) == 0 {
		usage()
	}

	if len(argsWithoutProg) > 0 && argsWithoutProg[0] == "-t" {
		terse = true
		argsWithoutProg = argsWithoutProg[1:]
	}

	for _, arg := range argsWithoutProg {
		printArg(arg)
	}

	if !success {
		usage()
	}

	if titled {
		fmt.Printf("\n# dates are counted using days since Thursday 1st Jan 1970\n")
		fmt.Printf("# clock operates via milliseconds since midnight\n")
	}
}

// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This tool prints equivalences between the string representation and the internal numerical
// representation for dates and clocks.
//
package main

import (
	"fmt"
	. "github.com/rickb777/date"
	"github.com/rickb777/date/clock"
	"os"
	"strconv"
	"strings"
)

func printPair(a string, b interface{}) {
	fmt.Printf("%-12s %12v\n", a, b)
}

func printOneDate(s string, d Date, err error) {
	if err != nil {
		printPair(s, err.Error())
	} else {
		printPair(s, d.Sub(Date{}))
	}
}

func printOneClock(s string, c clock.Clock, err error) {
	if err != nil {
		printPair(s, err.Error())
	} else {
		printPair(s, int32(c))
	}
}

func printArg(arg string) {
	d := Date{}

	d, e1 := AutoParse(arg)
	if e1 == nil {
		printPair(arg, d.Sub(Date{}))
	} else if strings.Index(arg, ":") == 2 {
		c, err := clock.Parse(arg)
		printOneClock(arg, c, err)
	} else {
		i, err := strconv.Atoi(arg)
		if err == nil {
			d = d.Add(PeriodOfDays(i))
			fmt.Printf("%-12s %12s  %s\n", arg, d, clock.Clock(i))
		} else {
			printPair(arg, err)
		}
	}
}

func main() {
	argsWithoutProg := os.Args[1:]
	for _, arg := range argsWithoutProg {
		printArg(arg)
	}
}

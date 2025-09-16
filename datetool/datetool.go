// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This tool prints equivalences between the string representation and the internal numerical
// representation for dates and clocks.
package main

import (
	"fmt"
	"github.com/govalues/decimal"
	"github.com/rickb777/date/v2"
	"github.com/rickb777/date/v2/clock"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"os"
	"strings"
	"time"
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

func printDate(d date.Date) string { return d.String() + " " + d.Weekday().String() }

func title() {
	if !terse && !titled {
		titled = true
		fmt.Printf("%-30s        possible value\n", "input")
		fmt.Printf("%-30s        --------------\n", "-----")
	}
}

const isoTimeNanos = "2006-01-02T15:04:05.999999999"

func printArg(arg string) {

	number, err := decimal.Parse(arg)
	if err == nil {
		title()
		i, _, _ := number.Int64(0)
		if number.IsInt() && i < 1000000 {
			d := date.Date(i)
			c := clock.Clock(i)
			fmt.Printf("%-30s clock: %-30s %sms since midnight\n", arg, c, sprintf(c))
			fmt.Printf("%-30s date:  %-30s %s days since 1AD\n", arg, printDate(d), sprintf(d))
			success = true
		}
		s, ns, _ := number.Int64(9)
		t := time.Unix(s, ns).UTC()
		fmt.Printf("%-30s time:  %s\n", arg, t.Format(isoTimeNanos))
		success = true
	}

	d, e1 := date.AutoParse(arg)
	if e1 == nil {
		title()
		fmt.Printf("%-30s date:  %-30s %s days since 1AD\n", arg, printDate(d), sprintf(d))
		success = true
	}

	c, err := clock.Parse(arg)
	if err == nil {
		title()
		fmt.Printf("%-30s clock: %-30s %sms since midnight\n", arg, c, sprintf(c))
		success = true
	}

	aug := arg
	if strings.IndexByte(arg, 'T') < 0 {
		aug += "T00:00:00"
	}
	t, err := time.Parse(isoTimeNanos, aug)
	if err == nil {
		if t.Year() > 1970+290 {
			number, _ := decimal.New(t.UTC().UnixMicro(), 6)
			fmt.Printf("%-30s time:  %ss\n", arg, number)
		} else {
			number, _ := decimal.New(t.UTC().UnixNano(), 9)
			fmt.Printf("%-30s time:  %ss\n", arg, number)
		}
		success = true
	}
	if success {
		fmt.Println()
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
}

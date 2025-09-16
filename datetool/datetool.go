// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This tool prints equivalences between the string representation and the internal numerical
// representation for dates and clocks.
package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/govalues/decimal"
	"github.com/rickb777/date/v2"
	"github.com/rickb777/date/v2/clock"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func usage() {
	fmt.Printf("Usage: %s [-t] number | date | time | now\n\n", os.Args[0])
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
		fmt.Printf("%-30s        interpretations\n", "input")
		fmt.Printf("%-30s        ---------------\n", "-----")
	}
}

const isoTimeNanos = "2006-01-02T15:04:05.999999999"

func print3Columns(a, b, c string) {
	fmt.Printf("%-30s %-6s %s\n", a, b, c)
}

func print4Columns(a, b, c, d string) {
	fmt.Printf("%-30s %-6s %-30s %s\n", a, b, c, d)
}

const (
	secondsSinceEpoch = "seconds since 1970 Unix epoch"
	daysSince1AD      = " days since 1AD"
	nsSinceMidnight   = "ns since midnight"
)

func printArg(arg string) {
	if arg == "now" {
		title()
		d := date.Today()
		t := time.Now()
		print4Columns(arg, "date:", printDate(d), sprintf(d)+daysSince1AD)
		number, _ := decimal.New(t.UTC().UnixMicro(), 6)
		print4Columns(arg, "time:", number.String(), secondsSinceEpoch)
		success = true
		return
	}

	number, err := decimal.Parse(arg)
	if err == nil {
		title()
		i, _, _ := number.Int64(0)
		if number.IsInt() && i < 1000000 {
			d := date.Date(i)
			c := clock.Clock(i)
			print4Columns(arg, "clock:", c.String(), sprintf(c)+nsSinceMidnight)
			print4Columns(arg, "date:", printDate(d), sprintf(d)+daysSince1AD)
			success = true
		}
		s, ns, _ := number.Int64(9)
		t := time.Unix(s, ns).UTC()
		print3Columns(arg, "time:", t.Format(isoTimeNanos))
		success = true
	}

	d, e1 := date.AutoParse(arg)
	if e1 == nil {
		title()
		print4Columns(arg, "date:", printDate(d), sprintf(d)+daysSince1AD)
		success = true
	}

	c, err := clock.Parse(arg)
	if err == nil {
		title()
		print4Columns(arg, "clock:", c.String(), sprintf(c)+nsSinceMidnight)
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
			print4Columns(arg, "time:", number.String(), secondsSinceEpoch)
		} else {
			number, _ := decimal.New(t.UTC().UnixNano(), 9)
			print4Columns(arg, "time:", number.String(), secondsSinceEpoch)
		}
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

	for i, arg := range argsWithoutProg {
		printArg(arg)

		if success && i < len(argsWithoutProg)-1 {
			fmt.Println()
		}
	}

	if !success {
		usage()
	}
}

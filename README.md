# date

[![GoDoc](https://img.shields.io/badge/api-Godoc-blue.svg)](https://pkg.go.dev/github.com/rickb777/date)
[![Build Status](https://api.travis-ci.org/rickb777/date.svg?branch=master)](https://travis-ci.org/rickb777/date/builds)
[![Coverage Status](https://coveralls.io/repos/rickb777/date/badge.svg?branch=master&service=github)](https://coveralls.io/github/rickb777/date?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/rickb777/date)](https://goreportcard.com/report/github.com/rickb777/date)
[![Issues](https://img.shields.io/github/issues/rickb777/date.svg)](https://github.com/rickb777/date/issues)

Package `date` provides functionality for working with dates.

This package introduces a light-weight `Date` type that is storage-efficient
and convenient for calendrical calculations and date parsing and formatting
(including years outside the [0,9999] interval).

It also provides

 * `clock.Clock` which expresses a wall-clock style hours-minutes-seconds with millisecond precision.
 * `period.Period` which expresses a period corresponding to the ISO-8601 form (e.g. "PT30S").
 * `timespan.DateRange` which expresses a period between two dates.
 * `timespan.TimeSpan` which expresses a duration of time between two instants.
 * `view.VDate` which wraps `Date` for use in templates etc.

See [package documentation](https://godoc.org/github.com/rickb777/date) for
full documentation and examples.

## Installation

    go get -u github.com/rickb777/date

or

    dep ensure -add github.com/rickb777/date

## Status

This library has been in reliable production use for some time. Versioning follows the well-known semantic version pattern.

## Credits

This package follows very closely the design of package
[`time`](http://golang.org/pkg/time/) in the standard library;
many of the `Date` methods are implemented using the corresponding methods
of the `time.Time` type and much of the documentation is copied directly
from that package.

The original [Good Work](https://github.com/fxtlabs/date) on which this was
based was done by Filippo Tampieri at Fxtlabs.

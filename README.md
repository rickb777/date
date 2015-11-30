# date

[![GoDoc](https://img.shields.io/badge/api-Godoc-blue.svg?style=flat-square)](https://godoc.org/github.com/rickb777/date)
[![Build Status](https://api.travis-ci.org/rickb777/date.svg?branch=master)](https://travis-ci.org/rickb777/date)
[![Coverage Status](https://coveralls.io/repos/rickb777/date/badge.svg?branch=master&service=github)](https://coveralls.io/github/rickb777/date?branch=master)

Package `date` provides functionality for working with dates.

This package introduces a light-weight `Date` type that is storage-efficient
and convenient for calendrical calculations and date parsing and formatting
(including years outside the [0,9999] interval).

It also provides

 * `TimeSpan` which expresses a duration of time between two instants, and
 * `DateRange` which expresses a period between two dats.

See [package documentation](https://godoc.org/github.com/rickb777/date) for
full documentation and examples.

## Installation

    go get -u github.com/rickb777/date

## Credits

This package follows very closely the design of package
[`time`](http://golang.org/pkg/time/) in the standard library;
many of the `Date` methods are implemented using the corresponding methods
of the `time.Time` type and much of the documentation is copied directly
from that package.

The original [Good Work](https://github.com/fxtlabs/date) on which this was
based was done by Fxtlabs.

# date

[![GoDoc](https://img.shields.io/badge/api-Godoc-blue.svg?style=flat-square)](https://godoc.org/github.com/fxtlabs/date)
[![Build Status](https://api.travis-ci.org/fxtlabs/date.svg?branch=master)](https://travis-ci.org/fxtlabs/date)
[![Coverage Status](https://coveralls.io/repos/fxtlabs/date/badge.svg?branch=master&service=github)](https://coveralls.io/github/fxtlabs/date?branch=master)

Package `date` provides functionality for working with dates.

This package introduces a light-weight `Date` type that is storage-efficient
and covenient for calendrical calculations and date parsing and formatting
(including years outside the [0,9999] interval).

See [package documentation](https://godoc.org/github.com/fxtlabs/date) for
full documentation and examples.

## Installation

    go get -u github.com/fxtlabs/date

## Credits

This package follows very closely the design of package
[`time`](http://golang.org/pkg/time/) in the standard library;
many of the `Date` methods are implemented using the corresponding methods
of the `time.Time` type and much of the documentation is copied directly
from that package.


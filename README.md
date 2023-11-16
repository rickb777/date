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
 * `timespan.DateRange` which expresses a period between two dates.
 * `timespan.TimeSpan` which expresses a duration of time between two instants (see RFC5545).
 * `view.VDate` which wraps `Date` for use in templates etc.

See [package documentation](https://godoc.org/github.com/rickb777/date) for
full documentation and examples.

See also [period.Period](https://pkg.go.dev/github.com/rickb777/period), which implements periods corresponding
to the ISO-8601 form (e.g. "PT30S").

## Installation

    go get github.com/rickb777/date/v2

## Status

This library has been in reliable production use for some time. Versioning follows the well-known semantic version pattern.

### Version 2

Changes since v1:

* The [period.Period](https://pkg.go.dev/github.com/rickb777/period) type has moved.
* `clock.Clock` now has nanosecond resolution (formerly millisecond resolution). 
* `date.Date` is now an integer that holds the number of days since year zero. Previously, it was a struct based on year 1970.
* `date.Date` time conversion methods have more explicit names - see table below.
* `date.Date` arithmetic and comparison operations now rely on Go operators; the corresponding methods have been deleted - see table below.
* `date.Date` zero value is now year 0 (Gregorian proleptic astronomical) so 1970 will no longer cause issues.
* `date.PeriodOfDays` has been moved to `timespan.PeriodOfDays`
* `date.DateString` has been deleted; the SQL `driver.Valuer` implementation is now pluggable and serves the same purpose more simply.

Renamed methods:

| Was        | Use instead         |
|------------|---------------------|
| Date`.Local` | Date`.Midnight`     |
| Date`.UTC`   | Date`.MidnightUTC`  |
| Date`.In`    | Date`.MidnightIn`   |

Deleted methods and functions:

| Was         | Use instead        |
|-------------|--------------------|
| Date`.Add`    | `+`                |
| Date`.Sub`    | `-`                |
| Date`.IsZero` | `== 0`             |
| Date`.Equal`  | `==`               |
| Date`.Before` | `<`                |
| Date`.After`  | `>`                |
| `date.IsLeap` | `gregorian.IsLeap` |
| `date.DaysIn` | `gregorian.DaysIn` |

Any v1 dates stored as integers will be incorrect; these can be corrected by adding 719528 to them, which is the number of days between year zero (v2) and 1970 (v1). Dates stored as strings will be unaffected.

## Credits

This package follows very closely the design of package
[`time`](http://golang.org/pkg/time/) in the standard library;
many of the `Date` methods are implemented using the corresponding methods
of the `time.Time` type and much of the documentation is copied directly
from that package.

The original [Good Work](https://github.com/fxtlabs/date) on which this was
based was done by Filippo Tampieri at Fxtlabs.

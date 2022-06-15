# date

This package is forked from https://github.com/rickb777/date. It is unmodified beyond using 32-bit integers to store Period components internally (instead of 16-bit components) to avoid integer overflows. The upstream package is itself being deprecated in favor of a replacement that addresses the overflow issue, but the replacement package is under active development and not yet ready for production use.

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

See [package documentation](https://godoc.org/github.com/voltusdev/date) for
full documentation and examples.

## Installation

    go get -u github.com/voltusdev/date

or

    dep ensure -add github.com/voltusdev/date

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

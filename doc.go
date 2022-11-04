// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package date provides functionality for working with dates.
// It implements a light-weight Date type that is storage-efficient
// and convenient for calendrical calculations and date parsing and formatting
// (including years outside the [0,9999] interval).
//
// Subpackages provide:
//
// * `clock.Clock` which expresses a wall-clock style hours-minutes-seconds with millisecond precision.
//
// * `period.Period` which expresses a period corresponding to the ISO-8601 form (e.g. "PT30S").
//
// * `timespan.DateRange` which expresses a period between two dates.
//
// * `timespan.TimeSpan` which expresses a duration of time between two instants.
//
// * `view.VDate` which wraps `Date` for use in templates etc.
//
// # Credits
//
// This package follows very closely the design of package time
// (http://golang.org/pkg/time/) in the standard library, many of the Date
// methods are implemented using the corresponding methods of the time.Time
// type, and much of the documentation is copied directly from that package.
//
// # References
//
// https://golang.org/src/time/time.go
//
// https://en.wikipedia.org/wiki/Gregorian_calendar
//
// https://en.wikipedia.org/wiki/Proleptic_Gregorian_calendar
//
// https://en.wikipedia.org/wiki/Astronomical_year_numbering
//
// https://en.wikipedia.org/wiki/ISO_8601
//
// https://tools.ietf.org/html/rfc822
//
// https://tools.ietf.org/html/rfc850
//
// https://tools.ietf.org/html/rfc1123
//
// https://tools.ietf.org/html/rfc3339
package date

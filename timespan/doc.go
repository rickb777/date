// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package timespan provides spans of time [timespan.TimeSpan], and ranges of dates [timespan.DateRange].
// Both are half-open intervals for which the start is included and the end is excluded.
// This allows for empty spans and also facilitates aggregating spans together.
//
// [timespan.TimeSpan] is compatible with RFC-5545 (iCalendar);
// see https://datatracker.ietf.org/doc/html/rfc5545#section-3.3.9
package timespan

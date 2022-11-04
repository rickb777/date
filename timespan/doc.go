// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package timespan provides spans of time (TimeSpan), and ranges of dates (DateRange).
// Both are half-open intervals for which the start is included and the end is excluded.
// This allows for empty spans and also facilitates aggregating spans together.
package timespan

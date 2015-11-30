// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package date

import "time"

const secondsPerDay = 60 * 60 * 24

// encode returns the number of days elapsed from date zero to the date
// corresponding to the given Time value.
func encode(t time.Time) int32 {
	// Compute the number of seconds elapsed since January 1, 1970 00:00:00
	// in the location specified by t and not necessarily UTC.
	// A Time value is represented internally as an offset from a UTC base
	// time; because we want to extract a date in the time zone specified
	// by t rather than in UTC, we need to compensate for the time zone
	// difference.
	_, offset := t.Zone()
	secs := t.Unix() + int64(offset)
	// Unfortunately operator / rounds towards 0, so negative values
	// must be handled differently
	if secs >= 0 {
		return int32(secs / secondsPerDay)
	}
	return -int32((secondsPerDay - 1 - secs) / secondsPerDay)
}

// decode returns the Time value corresponding to 00:00:00 UTC of the date
// represented by d, the number of days elapsed since date zero.
func decode(d int32) time.Time {
	secs := int64(d) * secondsPerDay
	return time.Unix(secs, 0).UTC()
}

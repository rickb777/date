// Copyright 2015 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package clock

import "fmt"

// Hh gets the clock-face number of hours as a two-digit string.
// It is calculated from the modulo time; see Mod24.
// Note the special case of midnight at the end of a day is "24".
func (c Clock) Hh() string {
	if c == Day {
		return "24"
	}
	cm := c.Mod24()
	return fmt.Sprintf("%02d", clockHour(cm))
}

// HhMm gets the clock-face number of hours and minutes as a five-character ISO-8601 time string.
// It is calculated from the modulo time; see Mod24.
// Note the special case of midnight at the end of a day is "24:00".
func (c Clock) HhMm() string {
	if c == Day {
		return "24:00"
	}
	cm := c.Mod24()
	return fmt.Sprintf("%02d:%02d", clockHour(cm), clockMinute(cm))
}

// HhMmSs gets the clock-face number of hours, minutes, seconds as an eight-character ISO-8601 time string.
// It is calculated from the modulo time; see Mod24.
// Note the special case of midnight at the end of a day is "24:00:00".
func (c Clock) HhMmSs() string {
	if c == Day {
		return "24:00:00"
	}
	cm := c.Mod24()
	return fmt.Sprintf("%02d:%02d:%02d", clockHour(cm), clockMinute(cm), clockSecond(cm))
}

// Hh12 gets the clock-face number of hours as a one- or two-digit string, followed by am or pm.
// Remember that midnight is 12am, noon is 12pm.
// It is calculated from the modulo time; see Mod24.
func (c Clock) Hh12() string {
	cm := c.Mod24()
	h, sfx := clockHour12(cm)
	return fmt.Sprintf("%d%s", h, sfx)
}

// HhMm12 gets the clock-face number of hours and minutes, followed by am or pm.
// Remember that midnight is 12am, noon is 12pm.
// It is calculated from the modulo time; see Mod24.
func (c Clock) HhMm12() string {
	cm := c.Mod24()
	h, sfx := clockHour12(cm)
	return fmt.Sprintf("%d:%02d%s", h, clockMinute(cm), sfx)
}

// HhMmSs12 gets the clock-face number of hours, minutes and seconds, followed by am or pm.
// Remember that midnight is 12am, noon is 12pm.
// It is calculated from the modulo time; see Mod24.
func (c Clock) HhMmSs12() string {
	cm := c.Mod24()
	h, sfx := clockHour12(cm)
	return fmt.Sprintf("%d:%02d:%02d%s", h, clockMinute(cm), clockSecond(cm), sfx)
}

// String gets the clock-face number of hours, minutes, seconds and fraction as an ISO-8601 time
// string.
//
// If the clock value has more than 24 hours, the excess is discarded (see Mod24).
//
// The number of decimal places depends on the clock value. If microsecond and nanosecond digits
// are non-zero, the result is given to nanosecond precision. Otherwise, a shorter form is used that
// only has millisecond precision.
//
// See TruncateMillisecond to obtain the shorter form always.
//
// The special case of midnight at the end of a day is "24:00:00.000".
func (c Clock) String() string {
	if c == Day {
		return "24:00:00.000"
	}
	cm := c.Mod24()
	if cm%Millisecond == 0 {
		return fmt.Sprintf("%02d:%02d:%02d.%03d", clockHour(cm), clockMinute(cm), clockSecond(cm), clockMillisecond(cm))
	}
	return fmt.Sprintf("%02d:%02d:%02d.%09d", clockHour(cm), clockMinute(cm), clockSecond(cm), clockNanosecond(cm))
}

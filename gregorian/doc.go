// Copyright 2016 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package gregorian provides utility functions for the Gregorian calendar calculations.
// The Gregorian calendar was officially introduced on 15th October 1582 so, strictly speaking,
// it only applies after that date. Some countries did not switch to the Gregorian calendar
// for many years after (such as Great Britain in 1782).
//
// Extending the Gregorian calendar backwards to dates preceding its official introduction
// produces a proleptic calendar that should be used with some caution for historic dates
// because it can lead to confusion.
//
// See https://en.wikipedia.org/wiki/Gregorian_calendar
// https://en.wikipedia.org/wiki/Proleptic_Gregorian_calendar
package gregorian

// Copyright 2016 Rick Beton. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package period provides functionality for periods of time using ISO-8601 conventions.
// This deals with years, months, weeks and days.
// Because of the vagaries of calendar systems, the meaning of year lengths, month lengths
// and even day lengths depends on context. So a period is not necessarily a fixed duration
// of time in terms of seconds.
//
// See https://en.wikipedia.org/wiki/ISO_8601#Periods
//
// Example representations:
//
// * "P4D" is four days;
//
// * "P3Y6M4W1D" is three years, 6 months, 4 weeks and one day.
//
// Note that ISO-8601 periods can also express periods of time in terms of hours, minutes and seconds,
// as well as years/months/weeks/days, but the hour/minute/second parts are not (yet) supported.
//
package period

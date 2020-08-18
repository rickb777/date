package period

import (
	"fmt"
	"math"
	"strings"
)

// used for stages in arithmetic
type period64 struct {
	centiMonths, centiDays, centiSeconds cent64
	neg                                  bool
	showAs                               string
}

func (period Period) toPeriod64() *period64 {
	return &period64{
		centiMonths:  cent64(period.centiMonths),
		centiDays:    cent64(period.centiDays),
		centiSeconds: cent64(period.centiSeconds),
		showAs:       period.showAs,
	}
}

func (p64 period64) toPeriod() Period {
	if p64.neg {
		return Period{
			centiMonths:  int32(-p64.centiMonths),
			centiDays:    int32(-p64.centiDays),
			centiSeconds: int32(-p64.centiSeconds),
			showAs:       p64.showAs,
		}
	}

	return Period{
		centiMonths:  int32(p64.centiMonths),
		centiDays:    int32(p64.centiDays),
		centiSeconds: int32(p64.centiSeconds),
		showAs:       p64.showAs,
	}
}

func (p64 period64) normalise64(precise bool) period64 {
	v := p64.abs().rippleUp(precise).moveFractionToRight()
	v.showAs = v.toIso8601String()
	return v
}

func (p64 period64) abs() period64 {
	if !p64.neg {
		if p64.centiMonths < 0 {
			p64.centiMonths = -p64.centiMonths
			p64.neg = true
		}

		if p64.centiDays < 0 {
			p64.centiDays = -p64.centiDays
			p64.neg = true
		}

		if p64.centiSeconds < 0 {
			p64.centiSeconds = -p64.centiSeconds
			p64.neg = true
		}
	}
	return p64
}

func (p64 period64) overflowedFields() []string {
	var m []string
	if p64.centiMonths > math.MaxInt32 {
		m = append(m, "months")
	}
	if p64.centiDays > math.MaxInt32 {
		m = append(m, "days")
	}
	if p64.centiSeconds > math.MaxInt32 {
		m = append(m, "seconds")
	}
	return m
}

func (p64 period64) rippleUp(precise bool) period64 {
	if !precise {
		if p64.centiSeconds > math.MaxInt32 {
			p64.centiDays += (p64.centiSeconds / 8640000) * 100
			p64.centiSeconds %= 8640000
		}

		if p64.centiDays > math.MaxInt32 {
			dE6 := p64.centiDays * oneE6
			p64.centiMonths += dE6 / daysPerMonthE6
			p64.centiDays = (dE6 % daysPerMonthE6) / oneE6
		}
	}

	return p64
}

// moveFractionToRight applies the rule that only the smallest field is permitted to have a decimal fraction.
func (p64 period64) moveFractionToRight() period64 {
	m100 := p64.centiMonths % 100
	if m100 != 0 && (p64.centiDays != 0 || p64.centiSeconds != 0) {
		p64.centiDays += (m100 * daysPerMonthE6) / oneE6
		p64.centiMonths = (p64.centiMonths / 100) * 100
	}

	d100 := p64.centiDays % 100
	if d100 != 0 && p64.centiSeconds != 0 {
		p64.centiSeconds += d100 * 86400
		p64.centiDays = (p64.centiDays / 100) * 100
	}

	return p64
}

func (p64 period64) toIso8601String() string {
	if p64.centiMonths == 0 && p64.centiDays == 0 && p64.centiSeconds == 0 {
		return "P0D"
	}

	buf := &strings.Builder{}
	if p64.neg {
		buf.WriteByte('-')
	}

	buf.WriteByte('P')

	if p64.centiMonths != 0 {
		years, centiMonths := p64.unpackYM()
		if years != 0 {
			fmt.Fprintf(buf, "%dY", years)
		}
		if centiMonths != 0 {
			fmt.Fprintf(buf, "%gM", absFloat100x(centiMonths))
		}
	}

	if p64.centiDays != 0 {
		if p64.centiDays%700 == 0 {
			fmt.Fprintf(buf, "%gW", absFloat100x(p64.centiDays/7))
		} else {
			fmt.Fprintf(buf, "%gD", absFloat100x(p64.centiDays))
		}
	}

	if p64.centiSeconds != 0 {
		hours, minutes, centiSeconds := p64.unpackHMS()
		buf.WriteByte('T')
		if hours != 0 {
			fmt.Fprintf(buf, "%dH", hours)
		}
		if minutes != 0 {
			fmt.Fprintf(buf, "%dM", minutes)
		}
		if centiSeconds != 0 {
			fmt.Fprintf(buf, "%gS", absFloat100x(centiSeconds))
		}
	}

	return buf.String()
}

func (p64 period64) unpackYM() (cent64, cent64) {
	years := p64.centiMonths / 1200
	months := p64.centiMonths - (years * 1200)
	return years, months
}

func (p64 period64) unpackHMS() (cent64, cent64, cent64) {
	hours := p64.centiSeconds / 360000
	seconds := p64.centiSeconds - hours*360000

	minutes := seconds / 6000
	seconds -= minutes * 6000
	return hours, minutes, seconds
}

func absFloat1x(v cent64) float32 {
	f := float32(v)
	if v < 0 {
		return -f
	}
	return f
}

func absFloat100x(v cent64) float32 {
	return absFloat1x(v) / 100
}

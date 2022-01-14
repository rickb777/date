package period

type period64MS struct {
	period64
	milliseconds int64
}

func (period PeriodMS) toPeriod64MS(input string) *period64MS {
	if period.IsNegative() {
		return &period64MS{
			period64: period64{
				years: int64(-period.years), months: int64(-period.months), days: int64(-period.days),
				hours: int64(-period.hours), minutes: int64(-period.minutes), seconds: int64(-period.seconds),
				neg:   true,
				input: input,
			},
			milliseconds: int64(-period.milliseconds),
		}
	}
	return &period64MS{
		period64: period64{
			years: int64(period.years), months: int64(period.months), days: int64(period.days),
			hours: int64(period.hours), minutes: int64(period.minutes), seconds: int64(period.seconds),
			input: input,
		},
		milliseconds: int64(period.milliseconds),
	}
}

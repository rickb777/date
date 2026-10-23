package main

import (
	"testing"
	"time"
)

func TestUnixSecondsString(t *testing.T) {
	cases := []struct {
		value time.Time
		want  string
	}{
		{time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC), "-62135596800.000000000"},
		{time.Date(1600, 12, 31, 23, 59, 59, 123456789, time.UTC), "-11644473600.876543211"},
		{time.Date(1970, 1, 1, 0, 0, 0, 1, time.UTC), "0.000000001"},
	}
	for _, c := range cases {
		if got := unixSecondsString(c.value); got != c.want {
			t.Errorf("unixSecondsString(%s) = %s, want %s", c.value, got, c.want)
		}
	}
}

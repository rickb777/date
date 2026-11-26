package date

import "testing"

// The expanded ISO 8601 year representation requires a sign followed by at
// least four digits. FormatOrdinal used a four-character field including the
// sign, so years from -1 to -999 were rendered with only three digits.

func TestFormatOrdinalNegativeShortYears(t *testing.T) {
	cases := []struct {
		input, want string
	}{
		{"-0001-001", "-0001-001"},
		{"-0009-032", "-0009-032"},
		{"-0099-100", "-0099-100"},
		{"-0999-365", "-0999-365"},
		{"-1000-001", "-1000-001"},
		{"-5000-001", "-5000-001"},
		{"0000-001", "0000-001"},
		{"0001-001", "0001-001"},
	}

	for _, c := range cases {
		d := MustParseISO(c.input)
		if got := d.FormatOrdinal(); got != c.want {
			t.Errorf("FormatOrdinal(MustParseISO(%q)) = %q, want %q", c.input, got, c.want)
		}
	}
}

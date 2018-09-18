package clock

import (
	"database/sql/driver"
	"testing"
	"time"
)

func TestClockScan(t *testing.T) {
	now := time.Now()

	cases := []struct {
		v        interface{}
		expected Clock
	}{
		{int64(New(-1, -1, -1, -1)), New(-1, -1, -1, -1)},
		{int64(New(10, 60, 10, 0)), New(10, 60, 10, 0)},
		{int64(New(24, 10, 0, 10)), New(0, 10, 0, 10)},
		{"12:00:00.400", New(12, 0, 0, 400)},
		{"01:40:50.000pm", New(13, 40, 50, 0)},
		{"4:20:00.000pm", New(16, 20, 0, 0)},
		{[]byte("23:60:60.000"), New(0, 1, 0, 0)},
		{now, NewAt(now)},
	}

	for i, c := range cases {
		var clock Clock
		e := clock.Scan(c.v)
		if e != nil {
			t.Errorf("%d: Got %v for %d", i, e, c.expected)
		} else if clock.Mod24() != c.expected.Mod24() {
			t.Errorf("%d: Got %v, want %d", i, clock, c.expected)
		}

		var d driver.Valuer = clock

		q, e := d.Value()
		if e != nil {
			t.Errorf("%d: Got %v for %d", i, e, c.expected)
		} else if Clock(q.(int64)).Mod24() != c.expected.Mod24() {
			t.Errorf("%d: Got %v, want %d", i, q, c.expected)
		}
	}
}

func TestClockScanWithJunk(t *testing.T) {
	cases := []struct {
		v        interface{}
		expected string
	}{
		{true, "bool true is not a meaningful clock"},
		{false, "bool false is not a meaningful clock"},
	}

	for i, c := range cases {
		var clock Clock
		e := clock.Scan(c.v)
		if e.Error() != c.expected {
			t.Errorf("%d: Got %q, want %q", i, e.Error(), c.expected)
		}
	}
}

func TestClockScanWithNil(t *testing.T) {
	var r *Clock
	e := r.Scan(nil)
	if e != nil {
		t.Errorf("Got %v", e)
	}
	if r != nil {
		t.Errorf("Got %v", r)
	}
}

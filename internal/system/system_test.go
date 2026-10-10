package system

import "testing"

func TestFormatUptime(t *testing.T) {
	cases := []struct {
		in   uint64
		want string
	}{
		{0, "0d 0h 0m"},
		{59, "0d 0h 0m"},
		{60, "0d 0h 1m"},
		{3661, "0d 1h 1m"},
		{90061, "1d 1h 1m"},
	}

	for _, c := range cases {
		if got := FormatUptime(c.in); got != c.want {
			t.Errorf("FormatUptime(%d) = %q, want %q", c.in, got, c.want)
		}
	}
}

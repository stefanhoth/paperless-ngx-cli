package cmd

import "testing"

func TestFormatStat(t *testing.T) {
	cases := []struct {
		in   interface{}
		want string
	}{
		{float64(43298748), "43298748"},
		{float64(4667), "4667"},
		{float64(1.5), "1.5"},
		{nil, "<nil>"},
		{"x", "x"},
	}
	for _, c := range cases {
		if got := formatStat(c.in); got != c.want {
			t.Errorf("formatStat(%v) = %q, want %q", c.in, got, c.want)
		}
	}
}

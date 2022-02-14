package main

import "testing"

func TestParseDuration(t *testing.T) {
	cases := []struct {
		in   string
		days int64
	}{
		{"5d", 5},
		{"5 days", 5},
		{"5 day", 5},
	}
	for _, tt := range cases {
		d, err := parseDuration(tt.in)
		if err != nil {
			t.Fatalf("parseDuration(%q): expected nil err, got %v", tt.in, err)
		}
		_ = d
	}
}

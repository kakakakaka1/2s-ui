package migration

import "testing"

func TestVersionBefore(t *testing.T) {
	cases := []struct {
		a, b string
		want bool
	}{
		{"1.5.6", "1.5.7", true},
		{"1.5.7", "1.5.7", false},
		{"1.5.8", "1.5.7", false},
		// The lexical-compare trap this function exists to fix.
		{"1.5.10", "1.5.7", false},
		{"1.5.9", "1.5.10", true},
		{"1.2", "1.5.4", true},
		{"1.4.2", "1.5.4", true},
		{"1.5.3", "1.5.4", true},
		{"2.0", "1.5.7", false},
		{"1.6", "1.5.7", false},
		// Padding: missing segments count as zero.
		{"1.5", "1.5.0", false},
		{"1.5.0", "1.5", false},
		// Empty sorts first, matching the old string-compare behaviour.
		{"", "1.5.7", true},
	}
	for _, c := range cases {
		if got := versionBefore(c.a, c.b); got != c.want {
			t.Errorf("versionBefore(%q, %q) = %v, want %v", c.a, c.b, got, c.want)
		}
	}
}

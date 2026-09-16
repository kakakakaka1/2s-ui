package service

import "testing"

// The self-updater used to ask "is the latest tag different from this binary",
// which answers yes in both directions.
//
// config/version is embedded and bumped by hand in the commit that prepares a
// release, so between that commit and the publish -- and on every build from
// main -- the running binary is newer than the newest release. A difference
// check calls that "not up to date" and installs the older tarball over it:
// the operator clicks Update and gets a downgrade, losing whatever the newer
// build fixed, while the database stays stamped at the newer version so no
// migration re-runs to meet it.
func TestShouldInstall(t *testing.T) {
	tests := []struct {
		name    string
		tag     string
		current string
		want    bool
	}{
		{"a newer release", "v1.8.3", "1.8.2", true},
		{"the same release", "v1.8.2", "1.8.2", false},
		// The v prefix is the tag's, not the version file's.
		{"the same release, no prefix", "1.8.2", "1.8.2", false},
		{"the same release, padded", "v1.8", "1.8.0", false},

		// The regression: main is ahead of the newest published release.
		{"this binary is newer", "v1.8.1", "1.8.2", false},
		{"this binary is a whole minor ahead", "v1.8.1", "1.9.0", false},

		// Two digits sort as numbers, which is the whole reason this shares
		// the migration gates' comparator rather than comparing text.
		{"ten is after nine", "v1.8.10", "1.8.9", true},
		{"nine is not after ten", "v1.8.9", "1.8.10", false},

		// A tag whose numbers do not order as newer is left alone. Refusing to
		// install something that cannot be ordered is the safe direction --
		// the operator can still install it by hand.
		{"a release candidate of the running version", "v1.8.2-rc1", "1.8.2", false},
		{"an unparseable tag", "vlatest", "1.8.2", false},
		{"whitespace around the tag", "  v1.8.3\n", "1.8.2", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldInstall(tt.tag, tt.current); got != tt.want {
				t.Errorf("shouldInstall(%q, %q) = %v, want %v", tt.tag, tt.current, got, tt.want)
			}
		})
	}
}

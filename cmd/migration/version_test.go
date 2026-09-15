package migration

import (
	"os"
	"regexp"
	"testing"

	"github.com/shenaba/2s-ui/config"
)

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

// Every migration gate has to name a version this binary claims to be, or it
// never runs on the installs it was written for.
//
// migrateDb returns at "Database is up to date" as soon as the stored version
// equals config.GetVersion(), so a gate on a version *newer* than the embedded
// one is unreachable for exactly the population that is already on the current
// release -- it is skipped on the upgrade that carries it, and the stored
// version never advances past it afterwards. That is how the 1.8.2 repair
// shipped dead: config/version was left at 1.8.1.
//
// The embedded version is edited by hand (config/version is //go:embed'ed, not
// injected with -X), and .github/workflows does not gate on it, so nothing else
// notices. This does.
func TestMigrationGatesAreReachable(t *testing.T) {
	source, err := os.ReadFile("main.go")
	if err != nil {
		t.Fatalf("read main.go: %v", err)
	}

	gates := regexp.MustCompile(`versionBefore\(dbVersion, "([^"]+)"\)`).FindAllStringSubmatch(string(source), -1)
	if len(gates) == 0 {
		t.Fatal("no versionBefore(dbVersion, ...) gates found; has migrateDb been rewritten?")
	}

	current := config.GetVersion()
	for _, gate := range gates {
		if versionBefore(current, gate[1]) {
			t.Errorf("migration gate %q is newer than config/version %q, so it can never run: bump config/version",
				gate[1], current)
		}
	}
}

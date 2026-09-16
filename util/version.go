package util

import (
	"strconv"
	"strings"
)

// VersionBefore reports whether release a is older than release b, comparing
// dot-separated numbers rather than text.
//
// The lexical compare it replaces got "1.5.10" < "1.5.7" wrong, which is the
// kind of mistake that silently skips a migration. Missing segments count as
// zero, so "1.5" and "1.5.0" are the same release, and a segment that is not a
// number counts as zero too -- this is handed whatever string the database or a
// release tag happens to hold, and an answer is more useful here than a panic.
//
// It lives in util because two callers need the same ordering and a second copy
// would drift: cmd/migration gates its one-off repairs on it, and
// service.shouldInstall decides whether a published release is actually newer
// than the running binary.
func VersionBefore(a, b string) bool {
	as := strings.Split(a, ".")
	bs := strings.Split(b, ".")
	n := len(as)
	if len(bs) > n {
		n = len(bs)
	}
	for i := 0; i < n; i++ {
		av, bv := 0, 0
		if i < len(as) {
			av, _ = strconv.Atoi(as[i])
		}
		if i < len(bs) {
			bv, _ = strconv.Atoi(bs[i])
		}
		if av != bv {
			return av < bv
		}
	}
	return false
}

//go:build !linux

package service

import "errors"

// Self-update never reaches the restart step off Linux (CanSelfUpdate gates
// it), so this stub only exists to keep non-Linux builds compiling.
func execRestart(installDir string) error {
	return errors.New("in-place restart is only supported on Linux")
}

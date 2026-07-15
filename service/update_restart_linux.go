//go:build linux

package service

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

// execRestart replaces the current process image with the container
// entrypoint, so the freshly installed binary boots exactly like a new
// container start (entrypoint.sh runs migrate, then execs sui as the same
// PID). The container itself never restarts, so no restart policy is needed.
// Falls back to a plain migrate+exec shell chain when entrypoint.sh is absent
// (e.g. custom images).
func execRestart(installDir string) error {
	env := os.Environ()
	entrypoint := filepath.Join(installDir, "entrypoint.sh")
	argv0, argv := entrypoint, []string{entrypoint}
	if _, err := os.Stat(entrypoint); err != nil {
		// Custom images without the bundled entrypoint: mirror what it does.
		bin := filepath.Join(installDir, "sui")
		script := fmt.Sprintf("%q migrate; exec %q", bin, bin)
		argv0, argv = "/bin/sh", []string{"sh", "-c", script}
	}
	// entrypoint.sh resolves ./sui relative to the working directory, which is
	// the install dir on a normal container start. Chdir last so a failure
	// before this point leaves the still-running process where it was.
	if err := os.Chdir(installDir); err != nil {
		return err
	}
	return syscall.Exec(argv0, argv, env)
}

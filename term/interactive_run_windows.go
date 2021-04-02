// +build windows
package term

import (
	"os"
	"os/exec"
	"os/signal"
)

// InteractiveRun runs the input command that starts a child process.
func (r Runner) InteractiveRun(name string, args []string) error {
	sig := make(chan os.Signal, 1)
	// See https://golang.org/pkg/os/signal/#hdr-Windows
	signal.Notify(sig, os.Interrupt)
	defer signal.Reset(os.Interrupt)
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

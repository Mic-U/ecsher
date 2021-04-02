// +build !windows

package term

import (
	"os"
	"os/exec"
	"os/signal"
)

func (r Runner) InteractiveRun(name string, args []string) error {
	// Ignore interrupt signal otherwise the program exits.
	signal.Ignore(os.Interrupt)
	defer signal.Reset(os.Interrupt)
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

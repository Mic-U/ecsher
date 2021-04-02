package term

import (
	"io"
	"os"
	"os/exec"
)

type Runner struct{}

func New() Runner {
	return Runner{}
}

type Option func(cmd *exec.Cmd)

// Stdin sets the internal *exec.Cmd's Stdin field.
func Stdin(r io.Reader) Option {
	return func(c *exec.Cmd) {
		c.Stdin = r
	}
}

// Stdout sets the internal *exec.Cmd's Stdout field.
func Stdout(writer io.Writer) Option {
	return func(c *exec.Cmd) {
		c.Stdout = writer
	}
}

// Stderr sets the internal *exec.Cmd's Stderr field.
func Stderr(writer io.Writer) Option {
	return func(c *exec.Cmd) {
		c.Stderr = writer
	}
}

func (r Runner) Run(name string, args []string, options ...Option) error {
	cmd := exec.Command(name, args...)

	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr

	for _, opt := range options {
		opt(cmd)
	}
	return cmd.Run()
}

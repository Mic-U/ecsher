package term

import (
	"os"
	"os/exec"
)

type Runner struct{}

func New() Runner {
	return Runner{}
}

type Option func(cmd *exec.Cmd)

func (r Runner) Run(name string, args []string, options ...Option) error {
	cmd := exec.Command(name, args...)

	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr

	for _, opt := range options {
		opt(cmd)
	}
	return cmd.Run()
}

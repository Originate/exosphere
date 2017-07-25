package processHelpers

import (
	"fmt"
	"os/exec"

	shellwords "github.com/mattn/go-shellwords"
	"github.com/pkg/errors"
)

// Run runs the given command, waits for the process to finish and
// returns the output string and error (if any)
func Run(dir string, commandWords ...string) (string, error) {
	if len(commandWords) == 1 {
		var err error
		commandWords, err = ParseCommand(commandWords[0])
		if err != nil {
			return "", err
		}
	}
	process := NewProcess(commandWords...)
	process.SetDir(dir)
	return process.Run()
}

// RunAndLog runs the given command, logs the process with the given
// function, waits for the process to finish and returns an error (if any)
func RunAndLog(dir string, env []string, log func(string), commandWords ...string) error {
	if len(commandWords) == 1 {
		var err error
		commandWords, err = ParseCommand(commandWords[0])
		if err != nil {
			return err
		}
	}
	cmd := exec.Command(commandWords[0], commandWords[1:]...) // nolint gas
	cmd.Dir = dir
	cmd.Env = env
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	go readPipe(stdoutPipe, log)
	if err := cmd.Start(); err != nil {
		return err
	}
	return cmd.Wait()
}

// RunSeries runs each command in commands and returns an error if any
func RunSeries(dir string, commands [][]string) error {
	for _, command := range commands {
		if output, err := Run(dir, command...); err != nil {
			return errors.Wrap(err, fmt.Sprintf("Command Failed:\nCommand: %s\nOutput:\n%s\n\n'", command, output))
		}
	}
	return nil
}

// ParseCommand parses the command string into a string array
func ParseCommand(command string) ([]string, error) {
	return shellwords.Parse(command)
}

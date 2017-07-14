package processHelpers

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"

	shellwords "github.com/mattn/go-shellwords"
	"github.com/pkg/errors"
)

// Run runs the given command, waits for the process to finish and
// returns the output string and error (if any)
func Run(dir string, commandWords ...string) (string, error) {
	if len(commandWords) == 1 {
		var err error
		commandWords, err = shellwords.Parse(commandWords[0])
		if err != nil {
			return "", err
		}
	}
	cmd := exec.Command(commandWords[0], commandWords[1:]...) // nolint gass
	cmd.Dir = dir
	outputArray, err := cmd.CombinedOutput()
	output := string(outputArray)
	return output, err
}

// Start runs the given command in the given dir directory, and returns
// the pointer to the command, stdout pipe, output buffer and error (if any)
func Start(dir string, commandWords ...string) (*exec.Cmd, io.WriteCloser, *bytes.Buffer, error) {
	if len(commandWords) == 1 {
		var err error
		commandWords, err = shellwords.Parse(commandWords[0])
		if err != nil {
			return nil, nil, nil, err
		}
	}
	cmd := exec.Command(commandWords[0], commandWords[1:]...) // nolint gas
	cmd.Dir = dir
	in, err := cmd.StdinPipe()
	var out bytes.Buffer
	cmd.Stdout = &out
	if err != nil {
		return nil, in, &out, err
	}
	if err = cmd.Start(); err != nil {
		return nil, in, &out, fmt.Errorf("Error running %s\nError:%s", commandWords, err)
	}
	return cmd, in, &out, nil
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

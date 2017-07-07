package processHelpers

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"

	shellwords "github.com/mattn/go-shellwords"
)

// Run runs the given command, waits for the process to finish and
// returns the output string and error (if any)
func Run(command string) (string, error) {
	commandWords, err := shellwords.Parse(command)
	if err != nil {
		return "", err
	}
	cmd := exec.Command(commandWords[0], commandWords[1:]...) // nolint gas
	outputArray, err := cmd.CombinedOutput()
	output := string(outputArray)
	return output, err
}

// Start runs the given command in the given dir directory, and returns
// the pointer to the command, stdout pipe, output buffer and error (if any)
func Start(command string, dir string, env []string) (*exec.Cmd, io.WriteCloser, *bytes.Buffer, error) {
	commandWords, err := shellwords.Parse(command)
	if err != nil {
		return nil, nil, &bytes.Buffer{}, err
	}
	cmd := exec.Command(commandWords[0], commandWords[1:]...) // nolint gas
	cmd.Dir = dir
	if len(env) > 0 {
		cmd.Env = env
	}
	in, err := cmd.StdinPipe()
	var out bytes.Buffer
	cmd.Stdout = &out
	if err != nil {
		return nil, in, &out, err
	}
	if err = cmd.Start(); err != nil {
		return nil, in, &out, fmt.Errorf("Error running %s\nError:%s", command, err)
	}
	return cmd, in, &out, nil
}

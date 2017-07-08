package processHelpers

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
)

// Run runs the given command, waits for the process to finish and
// returns the output string and error (if any)
func Run(dir string, command ...string) (string, error) {
	cmd := exec.Command(command[0], command[1:]...) // nolint gass
	cmd.Dir = dir
	outputArray, err := cmd.CombinedOutput()
	output := string(outputArray)
	return output, err
}

// Start runs the given command in the given dir directory, and returns
// the pointer to the command, stdout pipe, output buffer and error (if any)
func Start(dir string, command ...string) (*exec.Cmd, io.WriteCloser, *bytes.Buffer, error) {
	cmd := exec.Command(command[0], command[1:]...) // nolint gas
	cmd.Dir = dir
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

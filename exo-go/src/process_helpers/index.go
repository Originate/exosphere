package processHelpers

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"

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
func Start(command string, dir string) (*exec.Cmd, io.WriteCloser, *bytes.Buffer, error) {
	commandWords, err := shellwords.Parse(command)
	if err != nil {
		return nil, nil, &bytes.Buffer{}, err
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
		return nil, in, &out, fmt.Errorf("Error running %s\nError:%s", command, err)
	}
	return cmd, in, &out, nil
}

// Wait waits for the given text for the specified duration
func Wait(stdout fmt.Stringer, text string, done func()) {
	ticker := time.NewTicker(100 * time.Millisecond)
	var output string
	for !strings.Contains(output, text) {
		select {
		case <-ticker.C:
			output = stdout.String()
		}
	}
	done()
}

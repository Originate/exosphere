package processHelpers

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/Originate/exosphere/exo-go/src/util"
)

// Process represents a exec.Cmd process
type Process struct {
	Cmd        *exec.Cmd
	StdoutLog  func(string)
	StdoutPipe io.ReadCloser
	StdinPipe  io.WriteCloser
	Output     string
}

// NewProcess is Process's constructor
func NewProcess(commandWords ...string) *Process {
	process := &Process{Cmd: exec.Command(commandWords[0], commandWords[1:]...)} //nolint gas
	return process
}

// Log reads the stream from the given stdPipeReader, logs the
// output, and update process.Output
func (process *Process) Log(stdPipeReader io.Reader) {
	scanner := bufio.NewScanner(stdPipeReader)
	for scanner.Scan() {
		text := scanner.Text()
		if process.StdoutLog != nil {
			process.StdoutLog(text)
		}
		process.Output = process.Output + text
	}
}

// Run runs the process, waits for the process to finish and
// returns the output string and error (if any)
func (process *Process) Run() (string, error) {
	outputArray, err := process.Cmd.CombinedOutput()
	process.Output = string(outputArray)
	return process.Output, err
}

// SetDir sets the directory that the process should be run in
func (process *Process) SetDir(dir string) {
	process.Cmd.Dir = dir
}

// SetEnv sets the environment for the process
func (process *Process) SetEnv(env []string) {
	process.Cmd.Env = env
}

// SetStdoutLog sets the function that process should use to log
// the stdout output
func (process *Process) SetStdoutLog(log func(string)) {
	process.StdoutLog = log
}

// Start runs the process and returns an error if any
func (process *Process) Start() error {
	var err error
	process.StdinPipe, err = process.Cmd.StdinPipe()
	if err != nil {
		return err
	}
	stdoutPipe, err := process.Cmd.StdoutPipe()
	if err != nil {
		return err
	}
	logPipeReader, exposedPipeReader := duplicateReader(stdoutPipe)
	process.StdoutPipe = exposedPipeReader
	go process.Log(logPipeReader)
	return process.Cmd.Start()
}

// WaitForText reads from the StdoutPipe stream waiting for the given text
// for the given duration and updates Output with the output it reads while
// waiting
func (process *Process) WaitForText(text string, duration int) error {
	var err error
	var output string
	return util.WaitForf(func() bool {
		buffer := make([]byte, 1000)
		var count int
		count, err = process.StdoutPipe.Read(buffer)
		if count == 0 {
			if err == io.EOF {
				return false
			}
			if err != nil {
				return false
			}
		}
		output = output + string(buffer)
		return strings.Contains(output, text)
	}, func() error {
		if err != nil && err != io.EOF {
			return err
		}
		return fmt.Errorf("Expected '%s' to include '%s'", output, text)
	}, duration)
}

// Wait waits for the process to finish, can only be called after Start()
func (process *Process) Wait() error {
	return process.Cmd.Wait()
}

package processHelpers

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strings"
	"syscall"
	"time"
)

// Process represents a exec.Cmd process
type Process struct {
	Cmd        *exec.Cmd
	StdoutLog  func(string)
	stdoutPipe *bufio.Reader
	StdinPipe  io.WriteCloser
	Output     string
}

// NewProcess is Process's constructor
func NewProcess(commandWords ...string) *Process {
	process := &Process{Cmd: exec.Command(commandWords[0], commandWords[1:]...)} //nolint gas
	return process
}

func (process *Process) isRunning() bool {
	err := process.Cmd.Process.Signal(syscall.Signal(0))
	return fmt.Sprint(err) != "os: process already finished"
}

// Kill kills the process if it is running
func (process *Process) Kill() error {
	if process.isRunning() {
		return process.Cmd.Process.Kill()
	}
	return nil
}

// log reads the stream from the given stdPipeReader, logs the
// output, and update process.Output
func (process *Process) log(stdPipeReader io.Reader) {
	scanner := bufio.NewScanner(stdPipeReader)
	for scanner.Scan() {
		text := scanner.Text()
		if process.StdoutLog != nil {
			process.StdoutLog(text)
		}
		process.Output = process.Output + text
	}
}

// readStdoutPipe reads 1000 bytes of string from stdoutPipe and
// returns the string and an error if any
func (process *Process) readStdoutPipe() (string, error) {
	buffer := make([]byte, 1000)
	count, err := process.stdoutPipe.Read(buffer)
	if count == 0 {
		if err == io.EOF {
			return "", err
		}
		if err != nil {
			return "", err
		}
	}
	return string(buffer), nil
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
	process.stdoutPipe = bufio.NewReader(exposedPipeReader)
	go process.log(logPipeReader)
	return process.Cmd.Start()
}

// Wait waits for the process to finish, can only be called after Start()
func (process *Process) Wait() error {
	return process.Cmd.Wait()
}

func (process *Process) waitFor(condition func(string) bool) error {
	var output string
	for {
		text, err := process.readStdoutPipe()
		if err != nil {
			return err
		}
		output = output + text
		if condition(output) {
			return nil
		}
	}
}

// WaitForRegex waits for the given regex and returns an error if any
func (process *Process) WaitForRegex(regex *regexp.Regexp) error {
	if regex.MatchString(process.Output) {
		return nil
	}
	return process.waitFor(func(output string) bool {
		return regex.MatchString(output)
	})
}

func (process *Process) waitForText(text string, err chan<- error) {
	if strings.Contains(process.Output, text) {
		err <- nil
	}
	err <- process.waitFor(func(output string) bool {
		return strings.Contains(output, text)
	})
}

// WaitForTextWithTimeout waits for the given text and returns an error if any
func (process *Process) WaitForTextWithTimeout(text string, duration int) error {
	waitErr := make(chan error)
	go process.waitForText(text, waitErr)
	select {
	case <-waitErr:
		return nil
	case <-time.After(time.Duration(duration) * time.Millisecond):
		return fmt.Errorf("Timed out after %d", duration)
	}
}

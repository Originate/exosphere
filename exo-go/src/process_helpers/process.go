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

// Process represents a exec.Cmd p
type Process struct {
	Cmd        *exec.Cmd
	StdoutLog  func(string)
	stdoutPipe *bufio.Reader
	StdinPipe  io.WriteCloser
	Output     string
}

// NewProcess is Process's constructor
func NewProcess(commandWords ...string) *Process {
	p := &Process{Cmd: exec.Command(commandWords[0], commandWords[1:]...)} //nolint gas
	return p
}

func (p *Process) isRunning() bool {
	err := p.Cmd.Process.Signal(syscall.Signal(0))
	return fmt.Sprint(err) != "os: process already finished"
}

// Kill kills the p if it is running
func (p *Process) Kill() error {
	if p.isRunning() {
		return p.Cmd.Process.Kill()
	}
	return nil
}

// log reads the stream from the given stdPipeReader, logs the
// output, and update p.Output
func (p *Process) log(stdPipeReader io.Reader) {
	scanner := bufio.NewScanner(stdPipeReader)
	for scanner.Scan() {
		text := scanner.Text()
		if p.StdoutLog != nil {
			p.StdoutLog(text)
		}
		p.Output = p.Output + text
	}
}

// readStdoutPipe reads 1000 bytes of string from stdoutPipe and
// returns the string and an error if any
func (p *Process) readStdoutPipe() (string, error) {
	buffer := make([]byte, 1000)
	count, err := p.stdoutPipe.Read(buffer)
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

// Run runs the p, waits for the p to finish and
// returns the output string and error (if any)
func (p *Process) Run() (string, error) {
	outputArray, err := p.Cmd.CombinedOutput()
	p.Output = string(outputArray)
	return p.Output, err
}

// SetDir sets the directory that the p should be run in
func (p *Process) SetDir(dir string) {
	p.Cmd.Dir = dir
}

// SetEnv sets the environment for the p
func (p *Process) SetEnv(env []string) {
	p.Cmd.Env = env
}

// SetStdoutLog sets the function that p should use to log
// the stdout output
func (p *Process) SetStdoutLog(log func(string)) {
	p.StdoutLog = log
}

// Start runs the p and returns an error if any
func (p *Process) Start() error {
	var err error
	p.StdinPipe, err = p.Cmd.StdinPipe()
	if err != nil {
		return err
	}
	stdoutPipe, err := p.Cmd.StdoutPipe()
	if err != nil {
		return err
	}
	logPipeReader, exposedPipeReader := duplicateReader(stdoutPipe)
	p.stdoutPipe = bufio.NewReader(exposedPipeReader)
	go p.log(logPipeReader)
	stderrPipe, err := p.Cmd.StderrPipe()
	if err != nil {
		return err
	}
	go p.log(stderrPipe)
	return p.Cmd.Start()
}

// Wait waits for the p to finish, can only be called after Start()
func (p *Process) Wait() error {
	return p.Cmd.Wait()
}

func (p *Process) waitFor(condition func(string) bool) error {
	var output string
	for {
		text, err := p.readStdoutPipe()
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
func (p *Process) WaitForRegex(regex *regexp.Regexp) error {
	if regex.MatchString(p.Output) {
		return nil
	}
	return p.waitFor(func(output string) bool {
		return regex.MatchString(output)
	})
}

func (p *Process) waitForText(text string, err chan<- error) {
	if strings.Contains(p.Output, text) {
		err <- nil
	}
	err <- p.waitFor(func(output string) bool {
		return strings.Contains(output, text)
	})
}

// WaitForTextWithTimeout waits for the given text and returns an error if any
func (p *Process) WaitForTextWithTimeout(text string, duration int) error {
	waitErr := make(chan error)
	go p.waitForText(text, waitErr)
	select {
	case <-waitErr:
		return nil
	case <-time.After(time.Duration(duration) * time.Millisecond):
		return fmt.Errorf("Timed out after %d", duration)
	}
}

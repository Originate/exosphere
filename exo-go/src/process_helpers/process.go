package processHelpers

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"syscall"
	"time"

	uuid "github.com/satori/go.uuid"
)

// Process represents a exec.Cmd p
type Process struct {
	Cmd                *exec.Cmd
	onOutputFuncsMutex sync.Mutex // lock for reading / updating process.onOutputFuncs
	onOutputFuncs      map[string]func(string)
	StdinPipe          io.WriteCloser
	outputMutex        sync.Mutex // lock for reading / updating process.Output
	Output             string
}

// NewProcess is Process's constructor
func NewProcess(commandWords ...string) *Process {
	p := &Process{
		Cmd:                exec.Command(commandWords[0], commandWords[1:]...), //nolint gas
		onOutputFuncsMutex: sync.Mutex{},
		onOutputFuncs:      map[string]func(string){},
		outputMutex:        sync.Mutex{},
	}
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
	scanner.Split(scanLinesOrPrompt)
	for scanner.Scan() {
		text := scanner.Text()
		// Its important to lock the outputMutex before the onOutputFuncsMutex
		// becuase the waitFor method locks outputMutex and then may or may not lock
		// the onOutputFuncsMutex. We need to use the same order to avoid deadlock
		p.outputMutex.Lock()
		p.onOutputFuncsMutex.Lock()
		for _, fn := range p.onOutputFuncs {
			go fn(text)
		}
		p.Output = p.Output + text
		p.onOutputFuncsMutex.Unlock()
		p.outputMutex.Unlock()
	}
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

// AddOutputFunc adds a function that process should call anytime there is new output
func (p *Process) AddOutputFunc(key string, log func(string)) {
	p.onOutputFuncsMutex.Lock()
	p.onOutputFuncs[key] = log
	p.onOutputFuncsMutex.Unlock()
}

// RemoveOutputFunc removes a function that process should call anytime there is new output
func (p *Process) RemoveOutputFunc(key string) {
	p.onOutputFuncsMutex.Lock()
	delete(p.onOutputFuncs, key)
	p.onOutputFuncsMutex.Unlock()
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
	go p.log(stdoutPipe)
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

func (p *Process) waitFor(condition func(string) bool, err chan<- error) {
	p.outputMutex.Lock()
	if condition(p.Output) {
		err <- nil
	} else {
		id := uuid.NewV4().String()
		p.AddOutputFunc(id, func(output string) {
			if condition(output) {
				err <- nil
				p.RemoveOutputFunc(id)
			}
		})
	}
	p.outputMutex.Unlock()
}

// WaitForRegex waits for the given regex
func (p *Process) WaitForRegex(regex *regexp.Regexp) {
	waitErr := make(chan error)
	go p.waitFor(func(output string) bool {
		return regex.MatchString(output)
	}, waitErr)
	<-waitErr
}

// WaitForTextWithTimeout waits for the given text and returns an error if any
func (p *Process) WaitForTextWithTimeout(text string, duration int) error {
	waitErr := make(chan error)
	go p.waitFor(func(output string) bool {
		return strings.Contains(output, text)
	}, waitErr)
	select {
	case <-waitErr:
		return nil
	case <-time.After(time.Duration(duration) * time.Millisecond):
		return fmt.Errorf("Timed out after %d, expected: '%s', output: '%s'", duration, text, p.Output)
	}
}

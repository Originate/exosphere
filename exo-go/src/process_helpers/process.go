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

// Process represents a exec.Cmd process
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
	process := &Process{
		Cmd:                exec.Command(commandWords[0], commandWords[1:]...), //nolint gas
		onOutputFuncsMutex: sync.Mutex{},
		onOutputFuncs:      map[string]func(string){},
		outputMutex:        sync.Mutex{},
	}
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
	scanner.Split(scanLinesOrPrompt)
	for scanner.Scan() {
		text := scanner.Text()
		process.onOutputFuncsMutex.Lock()
		process.outputMutex.Lock()
		fns := []func(string){}
		for _, fn := range process.onOutputFuncs {
			fns = append(fns, fn)
		}
		process.Output = process.Output + text
		process.outputMutex.Unlock()
		process.onOutputFuncsMutex.Unlock()
		// calls fns after releasing the locks in case an fn calls RemoveOutputFunc
		// (otherwise would have deadlock with the onOutputFuncsMutex)
		for _, fn := range fns {
			fn(text)
		}
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

// AddOutputFunc adds a function that process should call anytime there is new output
func (process *Process) AddOutputFunc(key string, log func(string)) {
	process.onOutputFuncsMutex.Lock()
	process.onOutputFuncs[key] = log
	process.onOutputFuncsMutex.Unlock()
}

// RemoveOutputFunc removes a function that process should call anytime there is new output
func (process *Process) RemoveOutputFunc(key string) {
	process.onOutputFuncsMutex.Lock()
	delete(process.onOutputFuncs, key)
	process.onOutputFuncsMutex.Unlock()
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
	go process.log(stdoutPipe)
	stderrPipe, err := process.Cmd.StderrPipe()
	if err != nil {
		return err
	}
	go process.log(stderrPipe)
	return process.Cmd.Start()
}

// Wait waits for the process to finish, can only be called after Start()
func (process *Process) Wait() error {
	return process.Cmd.Wait()
}

func (process *Process) waitFor(condition func(string) bool, err chan<- error) {
	process.outputMutex.Lock()
	if condition(process.Output) {
		err <- nil
		return
	}
	id := uuid.NewV4().String()
	process.AddOutputFunc(id, func(output string) {
		if condition(output) {
			err <- nil
			process.RemoveOutputFunc(id)
		}
	})
	process.outputMutex.Unlock()
}

// WaitForRegex waits for the given regex and returns an error if any
func (process *Process) WaitForRegex(regex *regexp.Regexp) error {
	waitErr := make(chan error)
	go process.waitFor(func(output string) bool {
		return regex.MatchString(output)
	}, waitErr)
	return <-waitErr
}

// WaitForTextWithTimeout waits for the given text and returns an error if any
func (process *Process) WaitForTextWithTimeout(text string, duration int) error {
	waitErr := make(chan error)
	go process.waitFor(func(output string) bool {
		return strings.Contains(output, text)
	}, waitErr)
	select {
	case <-waitErr:
		return nil
	case <-time.After(time.Duration(duration) * time.Millisecond):
		return fmt.Errorf("Timed out after %d, expected: '%s', output: '%s'", duration, text, process.Output)
	}
}

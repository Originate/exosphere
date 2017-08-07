package execplus

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

// CmdPlus represents a more observable exec.Cmd command
type CmdPlus struct {
	Cmd       *exec.Cmd
	StdinPipe io.WriteCloser

	output         string
	outputChannels map[string]chan OutputChunk
	mutex          sync.Mutex // lock for updating Output and outputChannels
	stdoutClosed   chan bool
	stderrClosed   chan bool
}

// NewCmdPlus is CmdPlus's constructor
func NewCmdPlus(commandWords ...string) *CmdPlus {
	p := &CmdPlus{
		Cmd:            exec.Command(commandWords[0], commandWords[1:]...), //nolint gas
		mutex:          sync.Mutex{},
		outputChannels: map[string]chan OutputChunk{},
		stdoutClosed:   make(chan bool),
		stderrClosed:   make(chan bool),
	}
	return p
}

// Kill kills the CmdPlus if it is running
func (c *CmdPlus) Kill() error {
	if c.isRunning() {
		return c.Cmd.Process.Kill()
	}
	return nil
}

// GetOutput returns the output thus far
func (c *CmdPlus) GetOutput() string {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.output
}

// GetOutputChannel returns a channel that passes OutputChunk as they occur.
// It will immediately send an OutputChunk with an empty chunk and the full output
// thus far. It also returns a function that when called closes the channel
func (c *CmdPlus) GetOutputChannel() (chan OutputChunk, func()) {
	id := uuid.NewV4().String()
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.outputChannels[id] = make(chan OutputChunk)
	c.sendOutputChunk(c.outputChannels[id], OutputChunk{Full: c.output})
	return c.outputChannels[id], func() {
		c.mutex.Lock()
		defer c.mutex.Unlock()
		delete(c.outputChannels, id)
	}
}

// Run is shorthand for Start() followed by Wait()
func (c *CmdPlus) Run() error {
	if err := c.Start(); err != nil {
		return err
	}
	return c.Wait()
}

// SetDir sets the directory that the CmdPlus should be run in
func (c *CmdPlus) SetDir(dir string) {
	c.Cmd.Dir = dir
}

// SetEnv sets the environment for the CmdPlus
func (c *CmdPlus) SetEnv(env []string) {
	c.Cmd.Env = env
}

// Start runs the CmdPlus and returns an error if any
func (c *CmdPlus) Start() error {
	var err error
	c.StdinPipe, err = c.Cmd.StdinPipe()
	if err != nil {
		return err
	}
	stdoutPipe, err := c.Cmd.StdoutPipe()
	if err != nil {
		return err
	}
	go c.log(stdoutPipe, c.stdoutClosed)
	stderrPipe, err := c.Cmd.StderrPipe()
	if err != nil {
		return err
	}
	go c.log(stderrPipe, c.stderrClosed)
	return c.Cmd.Start()
}

// Wait waits for the CmdPlus to finish, can only be called after Start()
func (c *CmdPlus) Wait() error {
	<-c.stdoutClosed
	<-c.stderrClosed
	return c.Cmd.Wait()
}

// WaitForCondition calls the given function with the latest chunk of output
// and the full output until it returns true
// returns an error if it does not match after the given duration
func (c *CmdPlus) WaitForCondition(condition func(string, string) bool, duration time.Duration) error {
	success := make(chan bool)
	go c.waitForCondition(condition, success)
	select {
	case <-success:
		return nil
	case <-time.After(duration):
		return fmt.Errorf("Timed out after %v, full output:\n%s", duration, c.GetOutput())
	}
}

// WaitForRegexp waits for the full output to match the given regex
// returns an error if it does not match after the given duration
func (c *CmdPlus) WaitForRegexp(isValid *regexp.Regexp, duration time.Duration) error {
	return c.WaitForCondition(func(outputChunk, fullOutput string) bool {
		return isValid.MatchString(fullOutput)
	}, duration)
}

// WaitForText waits for the full output to contain the given text
// returns an error if it does not match after the given duration
func (c *CmdPlus) WaitForText(text string, duration time.Duration) error {
	return c.WaitForCondition(func(outputChunk, fullOutput string) bool {
		return strings.Contains(fullOutput, text)
	}, duration)
}

// Helpers

func (c *CmdPlus) isRunning() bool {
	err := c.Cmd.Process.Signal(syscall.Signal(0))
	return fmt.Sprint(err) != "os: process already finished"
}

func (c *CmdPlus) log(reader io.Reader, closed chan bool) {
	scanner := bufio.NewScanner(reader)
	scanner.Split(scanLinesOrPrompt)
	for scanner.Scan() {
		text := scanner.Text()
		c.mutex.Lock()
		if c.output != "" {
			c.output += "\n"
		}
		c.output += text
		outputChunk := OutputChunk{Chunk: text, Full: c.output}
		for _, outputChannel := range c.outputChannels {
			c.sendOutputChunk(outputChannel, outputChunk)
		}
		c.mutex.Unlock()
	}
	closed <- true
}

func (c *CmdPlus) sendOutputChunk(outputChannel chan OutputChunk, outputChunk OutputChunk) {
	go func() {
		outputChannel <- outputChunk
	}()
}

func (c *CmdPlus) waitForCondition(condition func(string, string) bool, success chan<- bool) {
	outputChannel, stopFunc := c.GetOutputChannel()
	for {
		outputChunk := <-outputChannel
		if condition(outputChunk.Chunk, outputChunk.Full) {
			success <- true
			stopFunc()
			return
		}
	}
}

package execplus

import (
	"bufio"
	"fmt"
	"io"
	"os"
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
	mutex          sync.RWMutex // lock for updating output and outputChannels
	stdoutClosed   chan bool
	stderrClosed   chan bool
}

// NewCmdPlus is CmdPlus's constructor
func NewCmdPlus(commandWords ...string) *CmdPlus {
	p := &CmdPlus{
		Cmd:            exec.Command(commandWords[0], commandWords[1:]...), //nolint gas
		outputChannels: map[string]chan OutputChunk{},
		stdoutClosed:   make(chan bool),
		stderrClosed:   make(chan bool),
	}
	return p
}

// Kill kills the CmdPlus if it is running
func (c *CmdPlus) Kill() error {
	if c.Cmd.Process != nil {
		err := c.Cmd.Process.Signal(syscall.Signal(0))
		if err == nil || err.Error() != "os: process already finished" {
			return c.Cmd.Process.Kill()
		}
	}
	return nil
}

// GetOutput returns the output thus far
func (c *CmdPlus) GetOutput() string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
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
	go c.sendInitialChunk(id, OutputChunk{Full: c.output})
	return c.outputChannels[id], c.getStopFunc(id)
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

// AppendEnv sets the environment for the CmdPlus to the current process environment
// appended with the given environment
func (c *CmdPlus) AppendEnv(env []string) {
	c.Cmd.Env = append(os.Environ(), env...)
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
	stderrPipe, err := c.Cmd.StderrPipe()
	if err != nil {
		return err
	}
	// Create the buffers before starting the command to ensure no output is lost
	stdoutScanner := bufio.NewScanner(stdoutPipe)
	stderrScanner := bufio.NewScanner(stderrPipe)
	defer func() {
		// Start scanning for output chunks after the command has started
		// in order to avoid a race condition around the stdout file descriptor
		// between scanning and c.Cmd.Start()
		go c.scanForOutputChunks(stdoutScanner, c.stdoutClosed)
		go c.scanForOutputChunks(stderrScanner, c.stderrClosed)
	}()
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

func (c *CmdPlus) getStopFunc(id string) func() {
	return func() {
		c.mutex.Lock()
		close(c.outputChannels[id])
		delete(c.outputChannels, id)
		c.mutex.Unlock()
	}
}

func (c *CmdPlus) sendInitialChunk(channelID string, chunk OutputChunk) {
	c.mutex.Lock()
	if c.outputChannels[channelID] != nil {
		c.outputChannels[channelID] <- chunk
	}
	c.mutex.Unlock()
}

func (c *CmdPlus) scanForOutputChunks(scanner *bufio.Scanner, closed chan bool) {
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
			outputChannel <- outputChunk
		}
		c.mutex.Unlock()
	}
	closed <- true
}

func (c *CmdPlus) waitForCondition(condition func(string, string) bool, success chan<- bool) {
	outputChannel, stopFunc := c.GetOutputChannel()
	stopping := false
	for {
		outputChunk, ok := <-outputChannel
		if !ok {
			return
		}
		if stopping {
			continue
		}
		if condition(outputChunk.Chunk, outputChunk.Full) {
			success <- true
			stopping = true
			go stopFunc()
		}
	}
}

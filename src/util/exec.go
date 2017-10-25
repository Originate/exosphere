package util

import (
	"io"
	"os"
	"os/exec"
	"strings"

	execplus "github.com/Originate/go-execplus"
	shellwords "github.com/mattn/go-shellwords"
	"github.com/pkg/errors"
)

// Run runs the given command, waits for the process to finish and
// returns the output string and error (if any)
func Run(dir string, commandWords ...string) (string, error) {
	if len(commandWords) == 1 {
		var err error
		commandWords, err = ParseCommand(commandWords[0])
		if err != nil {
			return "", err
		}
	}
	cmdPlus := execplus.NewCmdPlus(commandWords...)
	cmdPlus.SetDir(dir)
	if err := cmdPlus.Run(); err != nil {
		return cmdPlus.GetOutput(), errors.Wrapf(err, "Error running '%s'. Output:\n%s", strings.Join(commandWords, " "), cmdPlus.GetOutput())
	}
	return cmdPlus.GetOutput(), nil
}

// RunAndLog runs the given command, logs the process to the given
// channel, waits for the process to finish and returns an error (if any)
func RunAndLog(dir string, env []string, logger *Logger, commandWords ...string) error {
	if len(commandWords) == 1 {
		var err error
		commandWords, err = ParseCommand(commandWords[0])
		if err != nil {
			return err
		}
	}
	cmdPlus := execplus.NewCmdPlus(commandWords...)
	cmdPlus.SetDir(dir)
	cmdPlus.AppendEnv(env)
	ConnectLogChannel(cmdPlus, logger)
	if err := cmdPlus.Run(); err != nil {
		return errors.Wrapf(err, "Error running '%s'. Output:\n%s", strings.Join(commandWords, " "), cmdPlus.GetOutput())
	}
	return nil
}

// RunAndPipe runs the given command, and logs the process to the given writer
func RunAndPipe(dir string, env []string, writer io.Writer, commandWords ...string) error {
	if len(commandWords) == 1 {
		var err error
		commandWords, err = ParseCommand(commandWords[0])
		if err != nil {
			return err
		}
	}
	cmd := exec.Command(commandWords[0], commandWords[1:]...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), env...)
	cmd.Stdout = writer
	cmd.Stderr = writer
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "Error running '%s'", strings.Join(commandWords, " "))
	}
	return nil
}

// RunSeries runs each command in commands and returns an error if any
func RunSeries(dir string, commands [][]string) error {
	for _, command := range commands {
		if _, err := Run(dir, command...); err != nil {
			return err
		}
	}
	return nil
}

// ParseCommand parses the command string into a string array
func ParseCommand(command string) ([]string, error) {
	return shellwords.Parse(command)
}

// ConnectLogChannel connects a log channel that wants to receive only each new output
func ConnectLogChannel(cmdPlus *execplus.CmdPlus, logger *Logger) {
	outputChannel, _ := cmdPlus.GetOutputChannel()
	go func() {
		for {
			outputChunk := <-outputChannel
			if outputChunk.Chunk != "" {
				logger.Log(outputChunk.Chunk)
			}
		}
	}()
}

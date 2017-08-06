package util

import (
	"fmt"
	"os"

	execplus "github.com/Originate/go-execplus"
	shellwords "github.com/mattn/go-shellwords"
	"github.com/pkg/errors"
)

// Run runs the given command, waits for the process to finish and
// returns the output string and error (if any)
func Run(dir string, commandWords ...string) (string, error) {
	if len(commandWords) == 1 {
		var err error
		commandWords, err = shellwords.Parse(commandWords[0])
		if err != nil {
			return "", err
		}
	}
	cmdPlus := execplus.NewCmdPlus(commandWords...)
	cmdPlus.SetDir(dir)
	err := cmdPlus.Run()
	return cmdPlus.Output, err
}

// RunAndLog runs the given command, logs the process to the given
// channel, waits for the process to finish and returns an error (if any)
func RunAndLog(dir string, env []string, logChannel chan string, commandWords ...string) error {
	if len(commandWords) == 1 {
		var err error
		commandWords, err = shellwords.Parse(commandWords[0])
		if err != nil {
			return err
		}
	}
	cmdPlus := execplus.NewCmdPlus(commandWords...)
	cmdPlus.SetDir(dir)
	cmdPlus.SetEnv(append(env, os.Environ()...))
	ConnectLogChannel(cmdPlus, logChannel)
	return cmdPlus.Run()
}

// RunSeries runs each command in commands and returns an error if any
func RunSeries(dir string, commands [][]string) error {
	for _, command := range commands {
		if output, err := Run(dir, command...); err != nil {
			return errors.Wrap(err, fmt.Sprintf("Command Failed:\nCommand: %s\nOutput:\n%s\n\n'", command, output))
		}
	}
	return nil
}

// ConnectLogChannel connects a log channel that wants to receive only each new output
func ConnectLogChannel(cmdPlus *execplus.CmdPlus, logChannel chan string) {
	outputChannel, _ := cmdPlus.GetOutputChannel()
	go func() {
		for {
			outputChunk := <-outputChannel
			if outputChunk.Chunk != "" {
				logChannel <- outputChunk.Chunk
			}
		}
	}()
}

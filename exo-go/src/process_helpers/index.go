package processHelpers

import (
	"fmt"
	"log"

	shellwords "github.com/mattn/go-shellwords"
	"github.com/pkg/errors"
)

// Run runs the given command, waits for the process to finish and
// returns the output string and error (if any)
func Run(dir string, commandWords ...string) (string, error) {
	process := NewProcess(commandWords...)
	process.SetDir(dir)
	return process.Run()
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

// ParseCommand parses the command string into a string array
func ParseCommand(command string) []string {
	commandWords, err := shellwords.Parse(command)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to parse the command '%s': %s", command, err))
	}
	return commandWords
}

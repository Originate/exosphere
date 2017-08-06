package run

import (
	"fmt"

	"github.com/pkg/errors"
)

// Series runs each command in commands and returns an error if any
func Series(dir string, commands [][]string) error {
	for _, command := range commands {
		if output, err := Run(dir, command...); err != nil {
			return errors.Wrap(err, fmt.Sprintf("Command Failed:\nCommand: %s\nOutput:\n%s\n\n'", command, output))
		}
	}
	return nil
}

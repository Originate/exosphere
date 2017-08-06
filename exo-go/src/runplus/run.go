package runplus

import (
	execplus "github.com/Originate/go-execplus"
	shellwords "github.com/mattn/go-shellwords"
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

package run

import (
	"os"

	"github.com/Originate/exosphere/exo-go/src/util"
	execplus "github.com/Originate/go-execplus"
	shellwords "github.com/mattn/go-shellwords"
)

// AndLog runs the given command, logs the process to the given
// channel, waits for the process to finish and returns an error (if any)
func AndLog(dir string, env []string, logChannel chan string, commandWords ...string) error {
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
	util.ConnectLogChannel(cmdPlus, logChannel)
	return cmdPlus.Run()
}

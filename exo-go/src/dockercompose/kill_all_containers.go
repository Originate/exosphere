package dockercompose

import (
	"github.com/Originate/exosphere/exo-go/src/run"
	execplus "github.com/Originate/go-execplus"
)

// KillAllContainers kills all the containers
func KillAllContainers(dockerComposeDir string, logChannel chan string) (*execplus.CmdPlus, error) {
	cmdPlus := execplus.NewCmdPlus("docker-compose", "down")
	cmdPlus.SetDir(dockerComposeDir)
	run.ConnectLogChannel(cmdPlus, logChannel)
	return cmdPlus, cmdPlus.Start()
}

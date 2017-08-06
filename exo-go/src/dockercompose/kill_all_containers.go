package dockercompose

import (
	"github.com/Originate/exosphere/exo-go/src/util"
	execplus "github.com/Originate/go-execplus"
)

// KillAllContainers kills all the containers
func KillAllContainers(dockerComposeDir string, logChannel chan string) (*execplus.CmdPlus, error) {
	cmdPlus := execplus.NewCmdPlus("docker-compose", "down")
	cmdPlus.SetDir(dockerComposeDir)
	util.ConnectLogChannel(cmdPlus, logChannel)
	return cmdPlus, cmdPlus.Start()
}

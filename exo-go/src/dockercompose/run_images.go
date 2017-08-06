package dockercompose

import (
	"github.com/Originate/exosphere/exo-go/src/runplus"
	execplus "github.com/Originate/go-execplus"
)

// RunImages runs the given docker images
func RunImages(images []string, env []string, dockerComposeDir string, logChannel chan string) (*execplus.CmdPlus, error) {
	cmdPlus := execplus.NewCmdPlus(append([]string{"docker-compose", "up"}, images...)...)
	cmdPlus.SetDir(dockerComposeDir)
	cmdPlus.SetEnv(env)
	runplus.ConnectLogChannel(cmdPlus, logChannel)
	return cmdPlus, cmdPlus.Start()
}

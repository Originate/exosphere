package docker

import (
	"github.com/Originate/exosphere/exo-go/src/util"
	execplus "github.com/Originate/go-execplus"
)

// BuildAllImages builds all the docker images defined in docker-compose.yml
func BuildAllImages(dockerComposeDir string, logChannel chan string) error {
	return util.RunAndLog(dockerComposeDir, []string{}, logChannel, "docker-compose", "build")
}

// CreateNewContainer creates a new docker container for the given service
func CreateNewContainer(serviceName string, env []string, dockerComposeDir string, logChannel chan string) error {
	return util.RunAndLog(dockerComposeDir, env, logChannel, "docker-compose", "create", "--build", serviceName)
}

// KillAllContainers kills all the containers
func KillAllContainers(dockerComposeDir string, logChannel chan string) (*execplus.CmdPlus, error) {
	cmdPlus := execplus.NewCmdPlus("docker-compose", "down")
	cmdPlus.SetDir(dockerComposeDir)
	util.ConnectLogChannel(cmdPlus, logChannel)
	return cmdPlus, cmdPlus.Start()
}

// KillContainer kills the docker container of the given service
func KillContainer(serviceName, dockerComposeDir string, logChannel chan string) error {
	return util.RunAndLog(dockerComposeDir, []string{}, logChannel, "docker-compose", "kill", serviceName)
}

// PullAllImages pulls all the docker images defined in docker-compose.yml
func PullAllImages(dockerComposeDir string, logChannel chan string) error {
	return util.RunAndLog(dockerComposeDir, []string{}, logChannel, "docker-compose", "pull")
}

// RunImages runs the given docker images
func RunImages(images []string, env []string, dockerComposeDir string, logChannel chan string) (*execplus.CmdPlus, error) {
	cmdPlus := execplus.NewCmdPlus(append([]string{"docker-compose", "up"}, images...)...)
	cmdPlus.SetDir(dockerComposeDir)
	cmdPlus.SetEnv(env)
	util.ConnectLogChannel(cmdPlus, logChannel)
	return cmdPlus, cmdPlus.Start()
}

// RestartContainer starts the docker container of the given service
func RestartContainer(serviceName string, env []string, dockerComposeDir string, logChannel chan string) error {
	return util.RunAndLog(dockerComposeDir, env, logChannel, "docker-compose", "restart", serviceName)
}

package compose

import (
	"os"

	"github.com/Originate/exosphere/src/util"
	execplus "github.com/Originate/go-execplus"
)

// BuildAllImages builds all the docker images defined in docker-compose.yml
func BuildAllImages(opts BaseOptions) error {
	return util.RunAndLog(opts.DockerComposeDir, opts.Env, opts.LogChannel, "docker-compose", "build")
}

// CreateNewContainer creates a new docker container for the given service
func CreateNewContainer(opts ImageOptions) error {
	return util.RunAndLog(opts.DockerComposeDir, opts.Env, opts.LogChannel, "docker-compose", "create", "--build", opts.ImageName)
}

// KillAllContainers kills all the containers
func KillAllContainers(opts BaseOptions) (*execplus.CmdPlus, error) {
	cmdPlus := execplus.NewCmdPlus("docker-compose", "down")
	cmdPlus.SetDir(opts.DockerComposeDir)
	cmdPlus.SetEnv(append(opts.Env, os.Environ()...))
	util.ConnectLogChannel(cmdPlus, opts.LogChannel)
	return cmdPlus, cmdPlus.Start()
}

// KillContainer kills the docker container of the given service
func KillContainer(opts ImageOptions) error {
	return util.RunAndLog(opts.DockerComposeDir, opts.Env, opts.LogChannel, "docker-compose", "kill", opts.ImageName)
}

// PullAllImages pulls all the docker images defined in docker-compose.yml
func PullAllImages(opts BaseOptions) error {
	return util.RunAndLog(opts.DockerComposeDir, opts.Env, opts.LogChannel, "docker-compose", "pull")
}

// RunImages runs the given docker images
func RunImages(opts ImagesOptions) (*execplus.CmdPlus, error) {
	cmdPlus := execplus.NewCmdPlus(append([]string{"docker-compose", "up"}, opts.ImageNames...)...)
	cmdPlus.SetDir(opts.DockerComposeDir)
	cmdPlus.SetEnv(append(opts.Env, os.Environ()...))
	util.ConnectLogChannel(cmdPlus, opts.LogChannel)
	return cmdPlus, cmdPlus.Start()
}

// RestartContainer starts the docker container of the given service
func RestartContainer(opts ImageOptions) error {
	return util.RunAndLog(opts.DockerComposeDir, opts.Env, opts.LogChannel, "docker-compose", "restart", opts.ImageName)
}

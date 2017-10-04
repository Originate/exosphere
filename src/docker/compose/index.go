package compose

import (
	"github.com/Originate/exosphere/src/util"
	execplus "github.com/Originate/go-execplus"
)

// BuildAllImages builds all the docker images defined in docker-compose.yml
func BuildAllImages(opts BaseOptions) error {
	return util.RunAndLog(opts.DockerComposeDir, opts.Env, opts.Logger, "docker-compose", "build")
}

// CreateNewContainer creates a new docker container for the given service
func CreateNewContainer(opts ImageOptions) error {
	return util.RunAndLog(opts.DockerComposeDir, opts.Env, opts.Logger, "docker-compose", "create", "--build", opts.ImageName)
}

// KillAllContainers kills all the containers
func KillAllContainers(opts BaseOptions) error {
	return util.RunAndLog(opts.DockerComposeDir, opts.Env, opts.Logger, "docker-compose", "down")
}

// KillContainer kills the docker container of the given service
func KillContainer(opts ImageOptions) error {
	return util.RunAndLog(opts.DockerComposeDir, opts.Env, opts.Logger, "docker-compose", "kill", opts.ImageName)
}

// PullAllImages pulls all the docker images defined in docker-compose.yml
func PullAllImages(opts BaseOptions) error {
	return util.RunAndLog(opts.DockerComposeDir, opts.Env, opts.Logger, "docker-compose", "pull")
}

// RunImages runs the given docker images
func RunImages(opts ImagesOptions) (*execplus.CmdPlus, error) {
	cmdPlus := execplus.NewCmdPlus(append([]string{"docker-compose", "up"}, opts.ImageNames...)...)
	cmdPlus.SetDir(opts.DockerComposeDir)
	cmdPlus.AppendEnv(opts.Env)
	util.ConnectLogChannel(cmdPlus, opts.Logger)
	return cmdPlus, cmdPlus.Start()
}

// RestartContainer starts the docker container of the given service
func RestartContainer(opts ImageOptions) error {
	return util.RunAndLog(opts.DockerComposeDir, opts.Env, opts.Logger, "docker-compose", "restart", opts.ImageName)
}

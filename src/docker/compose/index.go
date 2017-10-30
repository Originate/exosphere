package compose

import (
	"github.com/Originate/exosphere/src/util"
)

// BuildImage builds the given docker image
func BuildImage(opts ImageOptions) error {
	return util.RunAndPipe(opts.DockerComposeDir, opts.Env, opts.Writer, "docker-compose", "build", opts.ImageName)
}

// BuildAllImages builds all the docker images defined in docker-compose.yml
func BuildAllImages(opts BaseOptions) error {
	return util.RunAndPipe(opts.DockerComposeDir, opts.Env, opts.Writer, "docker-compose", "build")
}

// KillAllContainers kills all the containers
func KillAllContainers(opts BaseOptions) error {
	return util.RunAndPipe(opts.DockerComposeDir, opts.Env, opts.Writer, "docker-compose", "down")
}

// PullImage builds the given docker image
func PullImage(opts ImageOptions) error {
	return util.RunAndPipe(opts.DockerComposeDir, opts.Env, opts.Writer, "docker-compose", "pull", opts.ImageName)
}

// PullAllImages pulls all the docker images defined in docker-compose.yml
func PullAllImages(opts BaseOptions) error {
	return util.RunAndPipe(opts.DockerComposeDir, opts.Env, opts.Writer, "docker-compose", "pull")
}

// RunImages runs the given docker images
func RunImages(opts ImagesOptions) error {
	cmd := []string{"docker-compose", "up"}
	if opts.AbortOnExit {
		cmd = append(cmd, "--abort-on-container-exit")
	}
	cmd = append(cmd, opts.ImageNames...)
	return util.RunAndPipe(opts.DockerComposeDir, opts.Env, opts.Writer, cmd...)
}

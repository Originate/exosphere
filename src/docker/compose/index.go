package compose

import (
	"github.com/Originate/exosphere/src/util"
)

// BuildImage builds the given docker image
func BuildImage(opts ImageOptions) error {
	cmd := []string{"docker-compose", "--file", opts.DockerComposeFileName, "build", opts.ImageName}
	return util.RunAndPipe(opts.DockerComposeDir, opts.Env, opts.Writer, cmd...)
}

// BuildAllImages builds all the docker images defined in docker-compose.yml
func BuildAllImages(opts BaseOptions) error {
	cmd := []string{"docker-compose", "--file", opts.DockerComposeFileName, "build"}
	return util.RunAndPipe(opts.DockerComposeDir, opts.Env, opts.Writer, cmd...)
}

// KillAllContainers kills all the containers
func KillAllContainers(opts BaseOptions) error {
	cmd := []string{"docker-compose", "--file", opts.DockerComposeFileName, "down"}
	return util.RunAndPipe(opts.DockerComposeDir, opts.Env, opts.Writer, cmd...)
}

// PullImage builds the given docker image
func PullImage(opts ImageOptions) error {
	cmd := []string{"docker-compose", "--file", opts.DockerComposeFileName, "pull", opts.ImageName}
	return util.RunAndPipe(opts.DockerComposeDir, opts.Env, opts.Writer, cmd...)
}

// PullAllImages pulls all the docker images defined in docker-compose.yml
func PullAllImages(opts BaseOptions) error {
	cmd := []string{"docker-compose", "--file", opts.DockerComposeFileName, "pull"}
	return util.RunAndPipe(opts.DockerComposeDir, opts.Env, opts.Writer, cmd...)
}

// RunImages runs the given docker images
func RunImages(opts ImagesOptions) error {
	cmd := []string{"docker-compose", "--file", opts.DockerComposeFileName, "up"}
	if opts.AbortOnExit {
		cmd = append(cmd, "--abort-on-container-exit")
	}
	cmd = append(cmd, opts.ImageNames...)
	return util.RunAndPipe(opts.DockerComposeDir, opts.Env, opts.Writer, cmd...)
}

// RunImage runs a given image
func RunImage(opts ImageOptions) error {
	cmd := []string{"docker-compose", "--file", opts.DockerComposeFileName, "up"}
	if opts.AbortOnExit {
		cmd = append(cmd, "--abort-on-container-exit")
	}
	if opts.Build {
		cmd = append(cmd, "--build")
	}
	cmd = append(cmd, opts.ImageName)
	return util.RunAndPipe(opts.DockerComposeDir, opts.Env, opts.Writer, cmd...)
}

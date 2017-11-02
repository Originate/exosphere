package compose

import (
	"github.com/Originate/exosphere/src/util"
)

var filename = "docker-compose.yml" //TODO inline opts.DockerComposeFileName once refactor is complete

// BuildImage builds the given docker image
func BuildImage(opts ImageOptions) error {
	if opts.DockerComposeFileName != "" {
		filename = opts.DockerComposeFileName
	}
	cmd := []string{"docker-compose", "--file", filename, "build", opts.ImageName}
	return util.RunAndPipe(opts.DockerComposeDir, opts.Env, opts.Writer, cmd...)
}

// BuildAllImages builds all the docker images defined in docker-compose.yml
func BuildAllImages(opts BaseOptions) error {
	if opts.DockerComposeFileName != "" {
		filename = opts.DockerComposeFileName
	}
	cmd := []string{"docker-compose", "--file", filename, "build"}
	return util.RunAndPipe(opts.DockerComposeDir, opts.Env, opts.Writer, cmd...)
}

// KillAllContainers kills all the containers
func KillAllContainers(opts BaseOptions) error {
	if opts.DockerComposeFileName != "" {
		filename = opts.DockerComposeFileName
	}
	cmd := []string{"docker-compose", "--file", filename, "down"}
	return util.RunAndPipe(opts.DockerComposeDir, opts.Env, opts.Writer, cmd...)
}

// PullImage builds the given docker image
func PullImage(opts ImageOptions) error {
	if opts.DockerComposeFileName != "" {
		filename = opts.DockerComposeFileName
	}
	cmd := []string{"docker-compose", "--file", filename, "pull", opts.ImageName}
	return util.RunAndPipe(opts.DockerComposeDir, opts.Env, opts.Writer, cmd...)
}

// PullAllImages pulls all the docker images defined in docker-compose.yml
func PullAllImages(opts BaseOptions) error {
	if opts.DockerComposeFileName != "" {
		filename = opts.DockerComposeFileName
	}
	cmd := []string{"docker-compose", "--file", filename, "pull"}
	return util.RunAndPipe(opts.DockerComposeDir, opts.Env, opts.Writer, cmd...)
}

// RunImages runs the given docker images
func RunImages(opts ImagesOptions) error {
	if opts.DockerComposeFileName != "" {
		filename = opts.DockerComposeFileName
	}
	cmd := []string{"docker-compose", "--file", filename, "up"}
	if opts.AbortOnExit {
		cmd = append(cmd, "--abort-on-container-exit")
	}
	cmd = append(cmd, opts.ImageNames...)
	return util.RunAndPipe(opts.DockerComposeDir, opts.Env, opts.Writer, cmd...)
}

// RunImage runs a given image
func RunImage(opts ImagesOptions, imageName string) error {
	if opts.DockerComposeFileName != "" {
		filename = opts.DockerComposeFileName
	}
	cmd := []string{"docker-compose", "--file", filename, "up"}
	if opts.AbortOnExit {
		cmd = append(cmd, "--abort-on-container-exit") //TODO
	}
	cmd = append(cmd, imageName)
	return util.RunAndPipe(opts.DockerComposeDir, opts.Env, opts.Writer, cmd...)
}

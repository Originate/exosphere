package compose

import (
	"github.com/Originate/exosphere/src/util"
)

// BuildImages builds images in opts.ImageNames, or builds all images if opts.ImageNames is empty
func BuildImages(opts CommandOptions) error {
	cmd := []string{"docker-compose", "--file", opts.DockerComposeFileName, "build"}
	cmd = append(cmd, opts.ImageNames...)
	return util.RunAndPipe(opts.DockerComposeDir, opts.Env, opts.Writer, cmd...)
}

// KillContainers kills images in opts.ImageNames, or kills all images if opts.ImageNames is empty
func KillContainers(opts CommandOptions) error {
	cmd := []string{"docker-compose", "--file", opts.DockerComposeFileName, "down"}
	cmd = append(cmd, opts.ImageNames...)
	return util.RunAndPipe(opts.DockerComposeDir, opts.Env, opts.Writer, cmd...)
}

// PullImages pulls images in opts.ImageNames, or pulls all images if opts.ImageNames is empty
func PullImages(opts CommandOptions) error {
	cmd := []string{"docker-compose", "--file", opts.DockerComposeFileName, "pull"}
	cmd = append(cmd, opts.ImageNames...)
	return util.RunAndPipe(opts.DockerComposeDir, opts.Env, opts.Writer, cmd...)
}

// RunImages runs the given docker images, or runs all images if opts.ImageNames is empty
func RunImages(opts CommandOptions) error {
	cmd := []string{"docker-compose", "--file", opts.DockerComposeFileName, "up"}
	if opts.AbortOnExit {
		cmd = append(cmd, "--abort-on-container-exit")
	}
	if opts.Build {
		cmd = append(cmd, "--build")
	}
	cmd = append(cmd, opts.ImageNames...)
	return util.RunAndPipe(opts.DockerComposeDir, opts.Env, opts.Writer, cmd...)
}

package composerunner

import (
	"fmt"

	"github.com/Originate/exosphere/src/docker/compose"
	"github.com/Originate/exosphere/src/docker/composebuilder"
	execplus "github.com/Originate/go-execplus"
)

// Run runs docker images based on the given options
func Run(options RunOptions) (*execplus.CmdPlus, error) {
	err := composebuilder.WriteYML(options.DockerComposeDir, options.DockerConfigs)
	if err != nil {
		return nil, err
	}
	err = buildAndPullImages(options)
	if err != nil {
		return nil, err
	}
	options.Logger.Log("setup complete")
	cmdPlus, err := compose.RunImages(compose.ImagesOptions{
		DockerComposeDir: options.DockerComposeDir,
		Logger:           options.Logger,
		Env:              []string{fmt.Sprintf("COMPOSE_PROJECT_NAME=%s", options.DockerComposeProjectName)},
	})
	return cmdPlus, err
}

func buildAndPullImages(options RunOptions) error {
	opts := compose.BaseOptions{
		DockerComposeDir: options.DockerComposeDir,
		Logger:           options.Logger,
		Env:              []string{fmt.Sprintf("COMPOSE_PROJECT_NAME=%s", options.DockerComposeProjectName)},
	}
	err := compose.PullAllImages(opts)
	if err != nil {
		return err
	}
	return compose.BuildAllImages(opts)
}

package runner

import (
	"fmt"

	"github.com/Originate/exosphere/src/docker/compose"
	"github.com/Originate/exosphere/src/docker/composebuilder"
)

// Run runs docker images based on the given options
func Run(options RunOptions) error {
	err := composebuilder.WriteYML(options.DockerComposeDir, options.DockerConfigs)
	if err != nil {
		return err
	}
	err = buildAndPullImages(options)
	if err != nil {
		return err
	}
	options.Logger.Log("setup complete")
	for _, imageGroup := range options.ImageGroups {
		err = runImageGroup(options, imageGroup)
		if err != nil {
			return err
		}
	}
	return nil
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

package application

import (
	"fmt"
	"path"

	"github.com/Originate/exosphere/src/docker/compose"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/util"
)

// CleanApplicationContainers cleans all Docker containers started by the
// docker-compose.yml file under appDir/tmp/
func CleanApplicationContainers(appDir, composeProjectName string, logger *util.Logger) error {
	dockerComposeFile := path.Join(appDir, "tmp", "docker-compose.yml")
	return killIfExists(dockerComposeFile, composeProjectName, logger)
}

// CleanServiceTestContainers cleans all Docker containers started by the
// docker-compose.yml files under appDir/serviceLocation/tests/tmp/
func CleanServiceTestContainers(appDir, composeProjectName string, logger *util.Logger) error {
	appConfig, err := types.NewAppConfig(appDir)
	if err != nil {
		return err
	}
	serviceData := appConfig.GetServiceData()
	for _, data := range serviceData {
		dockerComposeFile := path.Join(appDir, data.Location, "tests", "tmp", "docker-compose.yml")
		err = killIfExists(dockerComposeFile, composeProjectName, logger)
		if err != nil {
			return err
		}
	}
	return nil
}

func killIfExists(dockerComposeFile, composeProjectName string, logger *util.Logger) error {
	exists, err := util.DoesFileExist(dockerComposeFile)
	if err != nil {
		return err
	}
	if exists {
		_, err = compose.KillAllContainers(compose.BaseOptions{
			DockerComposeDir: path.Dir(dockerComposeFile),
			Logger:           logger,
			Env:              []string{fmt.Sprintf("COMPOSE_PROJECT_NAME=%s", composeProjectName)},
		})
	}
	return err
}

package application

import (
	"fmt"
	"path"

	"github.com/Originate/exosphere/src/config"
	"github.com/Originate/exosphere/src/docker/compose"
	"github.com/Originate/exosphere/src/docker/composebuilder"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/util"
)

// CleanApplicationContainers cleans all Docker containers started by the
// docker-compose.yml file under appDir/tmp/
func CleanApplicationContainers(appDir string, logChannel chan string) error {
	dockerComposeFile := path.Join(appDir, "tmp", "docker-compose.yml")
	exists, err := util.DoesFileExist(dockerComposeFile)
	if err != nil {
		return err
	}
	if exists {
		composeProjectName := composebuilder.GetDockerComposeProjectName(appDir)
		_, err = compose.KillAllContainers(compose.BaseOptions{
			DockerComposeDir: path.Join(appDir, "tmp"),
			LogChannel:       logChannel,
			Env:              []string{fmt.Sprintf("COMPOSE_PROJECT_NAME=%s", composeProjectName)},
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// CleanServiceTestContainers cleans all Docker containers started by the
// docker-compose.yml files under appDir/serviceLocation/tests/tmp/
func CleanServiceTestContainers(appDir string, logChannel chan string) error {
	appConfig, err := types.NewAppConfig(appDir)
	if err != nil {
		return err
	}
	serviceData := config.GetServiceData(appConfig.Services)
	for _, data := range serviceData {
		dockerComposeFile := path.Join(appDir, data.Location, "tests", "tmp", "docker-compose.yml")
		exists, err := util.DoesFileExist(dockerComposeFile)
		if err != nil {
			return err
		}
		if exists {
			composeProjectName := fmt.Sprintf("%stests", composebuilder.GetDockerComposeProjectName(appDir))
			_, err = compose.KillAllContainers(compose.BaseOptions{
				DockerComposeDir: path.Dir(dockerComposeFile),
				LogChannel:       logChannel,
				Env:              []string{fmt.Sprintf("COMPOSE_PROJECT_NAME=%s", composeProjectName)},
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

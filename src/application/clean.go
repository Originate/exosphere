package application

import (
	"path"

	"github.com/Originate/exosphere/src/config"
	"github.com/Originate/exosphere/src/docker/compose"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/util"
)

func CleanApplicationContainers(appDir string, logChannel chan string) error {
	dockerComposeFile := path.Join(appDir, "tmp", "docker-compose.yml")
	exists, err := util.DoesFileExist(dockerComposeFile)
	if err != nil {
		return err
	}
	if exists {
		compose.KillAllContainers(compose.BaseOptions{
			DockerComposeDir: path.Join(appDir, "tmp"),
			LogChannel:       logChannel,
		})
	}
	return nil
}

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
			compose.KillAllContainers(compose.BaseOptions{
				DockerComposeDir: path.Base(dockerComposeFile),
				LogChannel:       logChannel,
			})
		}
	}
	return nil
}

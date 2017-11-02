package application

import (
	"fmt"
	"io"
	"path"

	"github.com/Originate/exosphere/src/docker/compose"
	"github.com/Originate/exosphere/src/docker/composebuilder"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/util"
)

// CleanApplicationContainers cleans all Docker containers started by the
// docker-compose.yml file under appDir/tmp/
func CleanApplicationContainers(appDir, composeProjectName string, writer io.Writer) error {
	dockerComposeFiles := composebuilder.GetLocalRunComposeFileNames()
	for _, dockerComposeFile := range dockerComposeFiles {
		err = killIfExists(dockerComposeFile, composeProjectName, writer)
		if err != nil {
			return err
		}
	}
	return nil
}

// CleanServiceTestContainers cleans all Docker containers started by the
// docker-compose.yml files under appDir/serviceLocation/tests/tmp/
func CleanServiceTestContainers(appContext types.AppContext, composeProjectName string, writer io.Writer) error {
	serviceData := appContext.Config.GetServiceData()
	for _, data := range serviceData {
		dockerComposeFile := path.Join(appContext.Location, data.Location, "tests", "tmp", "docker-compose.yml")
		err := killIfExists(dockerComposeFile, composeProjectName, writer)
		if err != nil {
			return err
		}
	}
	return nil
}

func killIfExists(dockerComposeFile, composeProjectName string, writer io.Writer) error {
	exists, err := util.DoesFileExist(dockerComposeFile)
	if err != nil {
		return err
	}
	if exists {
		err = compose.KillAllContainers(compose.BaseOptions{
			DockerComposeDir: path.Dir(dockerComposeFile),
			Writer:           writer,
			Env:              []string{fmt.Sprintf("COMPOSE_PROJECT_NAME=%s", composeProjectName)},
		})
	}
	return err
}

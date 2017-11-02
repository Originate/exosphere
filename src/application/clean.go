package application

import (
	"fmt"
	"io"
	"path"

	"github.com/Originate/exosphere/src/docker/compose"
	"github.com/Originate/exosphere/src/docker/composebuilder"
	"github.com/Originate/exosphere/src/util"
)

// CleanServiceTestContainers cleans application Docker containers under appDir/docker-compose
func CleanApplicationContainers(appDir, composeProjectName string, writer io.Writer) error {
	dockerComposeFileNames := composebuilder.GetLocalRunComposeFileNames()
	for _, dockerComposeFileName := range dockerComposeFileNames {
		dockerComposeFilePath := path.Join(appDir, "docker-compose", dockerComposeFile)
		err = killIfExists(dockerComposeFilePath, composeProjectName, writer)
		if err != nil {
			return err
		}
	}
	return nil
}

// CleanServiceTestContainers cleans test Docker containers under appDir/docker-compose
func CleanServiceTestContainers(appDir, composeProjectName string, writer io.Writer) error {
	dockerComposeFilePath := path.Join(appDir, "docker-compose", composebuilder.LocalTestComposeFileName)
	return killIfExists(dockerComposeFilePath, composeProjectName, writer)
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

package application

import (
	"fmt"
	"io"
	"path"

	"github.com/Originate/exosphere/src/docker/compose"
	"github.com/Originate/exosphere/src/docker/composebuilder"
	"github.com/Originate/exosphere/src/util"
)

// CleanApplicationContainers cleans application Docker containers under appDir/docker-compose
func CleanApplicationContainers(appDir, composeProjectName string, writer io.Writer) error {
	dockerComposeFileNames := composebuilder.GetLocalRunComposeFileNames()
	for _, dockerComposeFileName := range dockerComposeFileNames {
		err := killIfExists(appDir, dockerComposeFileName, composeProjectName, writer)
		if err != nil {
			return err
		}
	}
	return nil
}

// CleanServiceTestContainers cleans test Docker containers under appDir/docker-compose
func CleanServiceTestContainers(appDir, composeProjectName string, writer io.Writer) error {
	return killIfExists(appDir, composebuilder.LocalTestComposeFileName, composeProjectName, writer)
}

func killIfExists(appDir, dockerComposeFileName, composeProjectName string, writer io.Writer) error {
	dockerComposeFilePath := path.Join(appDir, "docker-compose", dockerComposeFileName)
	exists, err := util.DoesFileExist(dockerComposeFilePath)
	if err != nil {
		return err
	}
	if exists {
		err = compose.KillContainers(compose.CommandOptions{
			DockerComposeDir:      path.Dir(dockerComposeFilePath),
			DockerComposeFileName: dockerComposeFileName,
			Writer:                writer,
			Env: []string{
				fmt.Sprintf("COMPOSE_PROJECT_NAME=%s", composeProjectName),
				fmt.Sprintf("APP_PATH=%s", appDir),
			},
		})
	}
	return err
}

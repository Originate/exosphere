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

// CleanContainers cleans all cantainers listed in the yaml files under appDir/docker-compose
func CleanContainers(appContext types.AppContext, writer io.Writer) error {
	err := GenerateComposeFiles(appContext)
	if err != nil {
		return err
	}
	composeProjectName := composebuilder.GetDockerComposeProjectName(appContext.Config.Name)
	for _, dockerComposeFileName := range composebuilder.GetComposeFileNames() {
		if dockerComposeFileName == composebuilder.LocalTestComposeFileName {
			composeProjectName = composebuilder.GetTestDockerComposeProjectName(appContext.Config.Name)
		}
		err = killIfExists(appContext.Location, dockerComposeFileName, composeProjectName, writer)
		if err != nil {
			return err
		}
	}
	return nil
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

package application

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/Originate/exosphere/src/docker/composebuilder"
	"github.com/Originate/exosphere/src/terraform"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/util"
	yaml "gopkg.in/yaml.v2"
)

var buildModes = []types.BuildMode{
	types.BuildMode{
		Type:        types.BuildModeTypeLocal,
		Environment: types.BuildModeEnvironmentTest,
	},
	types.BuildMode{
		Type:        types.BuildModeTypeLocal,
		Mount:       true,
		Environment: types.BuildModeEnvironmentDevelopment,
	},
	types.BuildMode{
		Type:        types.BuildModeTypeLocal,
		Environment: types.BuildModeEnvironmentProduction,
	},
}

// CheckGeneratedFiles checks if docker-compose and terraform files are up-to-date
func CheckGeneratedFiles(appContext types.AppContext, deployConfig types.DeployConfig) error {
	for _, buildMode := range buildModes {
		dockerCompose, err := composebuilder.GetApplicationDockerCompose(composebuilder.ApplicationOptions{
			AppConfig: appContext.Config,
			AppDir:    appContext.Location,
			BuildMode: buildMode,
		})
		if err != nil {
			return err
		}
		err = diffDockerCompose(dockerCompose, appContext.Location, buildMode.GetDockerComposeFileName())
		if err != nil {
			return err
		}
	}
	return terraform.GenerateCheck(deployConfig)
}

// GenerateComposeFiles generates all docker-compose files for exosphere commands
func GenerateComposeFiles(appContext types.AppContext) error {
	composeDir := path.Join(appContext.Location, "docker-compose")
	for _, buildMode := range buildModes {
		dockerCompose, err := composebuilder.GetApplicationDockerCompose(composebuilder.ApplicationOptions{
			AppConfig: appContext.Config,
			AppDir:    appContext.Location,
			BuildMode: buildMode,
		})
		if err != nil {
			return err
		}
		err = composebuilder.WriteYML(composeDir, buildMode.GetDockerComposeFileName(), dockerCompose)
		if err != nil {
			return err
		}
	}
	return nil
}

//GenerateTerraformFiles generates the terraform/main.tf file
func GenerateTerraformFiles(deployConfig types.DeployConfig) error {
	return terraform.GenerateFile(deployConfig)
}

func diffDockerCompose(newDockerCompose *types.DockerCompose, appDir, dockerComposeFileName string) error {
	dockerComposeRelativeFilePath := path.Join("docker-compose", dockerComposeFileName)
	dockerComposeAbsoluteFilePath := path.Join(appDir, dockerComposeRelativeFilePath)
	fileExists, err := util.DoesFileExist(dockerComposeAbsoluteFilePath)
	if err != nil {
		return err
	}
	if fileExists {
		currDockerComposeContent, err := ioutil.ReadFile(dockerComposeAbsoluteFilePath)
		if err != nil {
			return err
		}
		newDockerComposeContent, err := yaml.Marshal(newDockerCompose)
		if err != nil {
			return err
		}
		if string(currDockerComposeContent) != string(newDockerComposeContent) {
			return fmt.Errorf("'%s' is out of date. Please run 'exo generate'", dockerComposeRelativeFilePath)
		}
	} else {
		return fmt.Errorf("'%s' does not exist. Please run 'exo generate'", dockerComposeRelativeFilePath)
	}
	return nil
}

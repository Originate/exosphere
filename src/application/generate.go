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

var buildModes = []composebuilder.BuildMode{
	composebuilder.BuildMode{
		Type:        composebuilder.BuildModeTypeLocal,
		Environment: composebuilder.BuildModeEnvironmentTest,
	},
	composebuilder.BuildMode{
		Type:        composebuilder.BuildModeTypeLocal,
		Mount:       true,
		Environment: composebuilder.BuildModeEnvironmentDevelopment,
	},
	composebuilder.BuildMode{
		Type:        composebuilder.BuildModeTypeLocal,
		Environment: composebuilder.BuildModeEnvironmentProduction,
	},
}

// CheckGeneratedFiles checks if docker-compose and terraform files are up-to-date
func CheckGeneratedFiles(appContext types.AppContext, deployConfig types.DeployConfig) error {
	for _, buildMode := range buildModes {
		dockerConfigs, err := composebuilder.GetApplicationDockerConfigs(composebuilder.ApplicationOptions{
			AppConfig: appContext.Config,
			AppDir:    appContext.Location,
			BuildMode: buildMode,
		})
		if err != nil {
			return err
		}
		err = diffDockerCompose(dockerConfigs, appContext.Location, buildMode.GetDockerComposeFileName())
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
		dockerConfigs, err := composebuilder.GetApplicationDockerConfigs(composebuilder.ApplicationOptions{
			AppConfig: appContext.Config,
			AppDir:    appContext.Location,
			BuildMode: buildMode,
		})
		if err != nil {
			return err
		}
		err = composebuilder.WriteYML(composeDir, buildMode.GetDockerComposeFileName(), dockerConfigs)
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

func diffDockerCompose(newDockerConfig types.DockerConfigs, appDir, dockerComposeFileName string) error {
	dockerComposeRelativeFilePath := path.Join("docker-compose", dockerComposeFileName)
	dockerComposeAbsoluteFilePath := path.Join(appDir, dockerComposeRelativeFilePath)
	fileExists, err := util.DoesFileExist(dockerComposeAbsoluteFilePath)
	if err != nil {
		return err
	}
	if fileExists {
		currDockerCompose, err := ioutil.ReadFile(dockerComposeAbsoluteFilePath)
		if err != nil {
			return err
		}
		newDockerCompose, err := yaml.Marshal(types.DockerCompose{
			Version:  "3",
			Services: newDockerConfig,
		})
		if err != nil {
			return err
		}
		if string(currDockerCompose) != string(newDockerCompose) {
			return fmt.Errorf("'%s' is out of date. Please run 'exo generate'", dockerComposeRelativeFilePath)
		}
	} else {
		return fmt.Errorf("'%s' does not exist. Please run 'exo generate'", dockerComposeRelativeFilePath)
	}
	return nil
}

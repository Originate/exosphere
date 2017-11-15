package application

import (
	"errors"
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
		dockerComposeFilePath := path.Join(appContext.Location, "docker-compose", buildMode.GetDockerComposeFileName())
		err = diffDockerCompose(dockerConfigs, dockerComposeFilePath)
		if err != nil {
			return err
		}
	}

	currTerraformFileBytes, err := terraform.ReadTerraformFile(deployConfig)
	if err != nil {
		return err
	}
	newTerraformFileContents, err := terraform.Generate(deployConfig)
	if err != nil {
		return err
	}
	if newTerraformFileContents != string(currTerraformFileBytes) {
		return errors.New("terraform files out of date. Please run 'exo generate'")
	}
	return nil
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

func diffDockerCompose(newDockerConfig types.DockerConfigs, dockerComposeFilePath string) error {
	fileExists, err := util.DoesFileExist(dockerComposeFilePath)
	if err != nil {
		return err
	}
	if fileExists {
		currDockerCompose, err := ioutil.ReadFile(dockerComposeFilePath)
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
			return fmt.Errorf("'%s' is out of date. Please run 'exo generate'", dockerComposeFilePath)
		}
	} else {
		return fmt.Errorf("'%s' does not exist. Please run 'exo generate'", dockerComposeFilePath)
	}
	return nil
}

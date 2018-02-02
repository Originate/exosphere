package application

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/Originate/exosphere/src/application/deployer"
	"github.com/Originate/exosphere/src/docker/composebuilder"
	"github.com/Originate/exosphere/src/docker/tools"
	"github.com/Originate/exosphere/src/terraform"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/context"
	"github.com/Originate/exosphere/src/types/deploy"
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

// CheckGeneratedDockerComposeFiles checks if docker-compose files are up-to-date
func CheckGeneratedDockerComposeFiles(appContext *context.AppContext) error {
	for _, buildMode := range buildModes {
		dockerCompose, err := composebuilder.GetApplicationDockerCompose(composebuilder.ApplicationOptions{
			AppContext: appContext,
			BuildMode:  buildMode,
		})
		if err != nil {
			return err
		}
		err = diffDockerCompose(dockerCompose, appContext.Location, buildMode.GetDockerComposeFileName())
		if err != nil {
			return err
		}
	}
	return nil
}

// GenerateComposeFiles generates all docker-compose files for exosphere commands
func GenerateComposeFiles(appContext *context.AppContext) error {
	composeDir := path.Join(appContext.Location, "docker-compose")
	for _, buildMode := range buildModes {
		dockerCompose, err := composebuilder.GetApplicationDockerCompose(composebuilder.ApplicationOptions{
			AppContext: appContext,
			BuildMode:  buildMode,
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

// GenerateTerraformFiles validates application configuration and generates the terraform file
func GenerateTerraformFiles(deployConfig deploy.Config) error {
	err := deployer.ValidateConfigs(deployConfig)
	if err != nil {
		return err
	}
	return terraform.GenerateFiles(deployConfig)
}

// GenerateTerraformVarFile generates a tfvar file
func GenerateTerraformVarFiles(deployConfig deploy.Config, secrets map[string]string) error {
	err := GenerateComposeFiles(deployConfig.AppContext)
	if err != nil {
		return err
	}
	serviceDockerImagesMap, err := getImagesMap(deployConfig)
	if err != nil {
		return err
	}
	err = terraform.GenerateInfrastructureVarFile(deployConfig, secrets)
	if err != nil {
		return err
	}
	return terraform.GenerateServicesVarFile(deployConfig, secrets, serviceDockerImagesMap)
}

// returns a map of service role to tagged images (without actually tagging or pushing those images)
// used for generating .tfvars files before actual deployment
func getImagesMap(deployConfig deploy.Config) (map[string]string, error) {
	dockerCompose, err := tools.GetDockerCompose(path.Join(deployConfig.AppContext.GetDockerComposeDir(), types.LocalProductionComposeFileName))
	if err != nil {
		return nil, err
	}
	imagesMap := map[string]string{}
	for serviceRole, taggedImage := range deployer.GetServiceImageNames(deployConfig, dockerCompose) {
		serviceLocation := path.Join(deployConfig.AppContext.Location, deployConfig.AppContext.ServiceContexts[serviceRole].Source.Location)
		imageName, imageVersion, err := deployer.GetImageNameAndVersion(taggedImage, serviceLocation, deployConfig.AppContext.Location)
		if err != nil {
			return nil, err
		}
		imagesMap[serviceRole] = fmt.Sprintf("%s:%s", imageName, imageVersion)
	}
	return imagesMap, nil
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

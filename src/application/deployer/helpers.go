package deployer

import (
	"fmt"

	"github.com/Originate/exosphere/src/config"
	"github.com/Originate/exosphere/src/types"
)

// GetImageNames returns a mapping from service/dependency names to image name on the user's machine
func GetImageNames(deployConfig types.DeployConfig, dockerComposeDir string, dockerCompose types.DockerCompose) (map[string]string, error) {
	images := getServiceImageNames(deployConfig, dockerComposeDir, dockerCompose)
	dependencyImages, err := getDependencyImageNames(deployConfig, dockerComposeDir)
	if err != nil {
		return nil, err
	}
	for k, v := range dependencyImages {
		images[k] = v
	}
	return images, nil
}

func getServiceImageNames(deployConfig types.DeployConfig, dockerComposeDir string, dockerCompose types.DockerCompose) map[string]string {
	images := map[string]string{}
	for _, serviceRole := range deployConfig.AppConfig.GetSortedServiceRoles() {
		dockerConfig := dockerCompose.Services[serviceRole]
		images[serviceRole] = buildImageName(dockerConfig, deployConfig.DockerComposeProjectName, serviceRole)
	}
	return images
}

func getDependencyImageNames(deployConfig types.DeployConfig, dockerComposeDir string) (map[string]string, error) {
	images := map[string]string{}
	serviceConfigs, err := config.GetServiceConfigs(deployConfig.AppDir, deployConfig.AppConfig)
	if err != nil {
		return nil, err
	}
	dependencies := config.GetBuiltProductionDependencies(deployConfig.AppConfig, serviceConfigs, deployConfig.AppDir)
	for dependencyName, dependency := range dependencies {
		if dependency.HasDockerConfig() {
			dockerConfig, err := dependency.GetDockerConfig()
			if err != nil {
				return nil, err
			}
			images[dependencyName] = buildImageName(dockerConfig, deployConfig.DockerComposeProjectName, dependencyName)
		}
	}
	return images, nil
}

func buildImageName(dockerConfig types.DockerConfig, dockerComposeProjectName, serviceRole string) string {
	if dockerConfig.Image != "" {
		return dockerConfig.Image
	}
	return fmt.Sprintf("%s_%s", dockerComposeProjectName, serviceRole)
}

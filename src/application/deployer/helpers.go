package deployer

import (
	"fmt"

	"github.com/Originate/exosphere/src/config"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/deploy"
)

// GetImageNames returns a mapping from service/dependency names to image name on the user's machine
func GetImageNames(deployConfig deploy.Config, dockerCompose types.DockerCompose) (map[string]string, error) {
	images := getServiceImageNames(deployConfig, dockerCompose)
	dependencyImages, err := getDependencyImageNames(deployConfig)
	if err != nil {
		return nil, err
	}
	for k, v := range dependencyImages {
		images[k] = v
	}
	return images, nil
}

func getServiceImageNames(deployConfig deploy.Config, dockerCompose types.DockerCompose) map[string]string {
	images := map[string]string{}
	for _, serviceRole := range deployConfig.AppContext.Config.GetSortedServiceRoles() {
		dockerConfig := dockerCompose.Services[serviceRole]
		images[serviceRole] = buildImageName(dockerConfig, deployConfig.GetDockerComposeProjectName(), serviceRole)
	}
	return images
}

func getDependencyImageNames(deployConfig deploy.Config) (map[string]string, error) {
	dependencies := config.GetBuiltRemoteDependencies(deployConfig.AppContext)
	images := map[string]string{}
	for dependencyName, dependency := range dependencies {
		if dependency.HasDockerConfig() {
			dockerConfig, err := dependency.GetDockerConfig()
			if err != nil {
				return nil, err
			}
			images[dependencyName] = buildImageName(dockerConfig, deployConfig.GetDockerComposeProjectName(), dependencyName)
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

package deployer

import (
	"fmt"

	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/deploy"
)

// GetServiceImageNames returns a mapping from service names to image name on the user's machine
func GetServiceImageNames(deployConfig deploy.Config, dockerCompose types.DockerCompose) map[string]string {
	images := map[string]string{}
	for _, serviceRole := range deployConfig.AppContext.Config.GetSortedServiceRoles() {
		dockerConfig := dockerCompose.Services[serviceRole]
		images[serviceRole] = buildImageName(dockerConfig, deployConfig.GetDockerComposeProjectName(), serviceRole)
	}
	return images
}

func buildImageName(dockerConfig types.DockerConfig, dockerComposeProjectName, serviceRole string) string {
	if dockerConfig.Image != "" {
		return dockerConfig.Image
	}
	return fmt.Sprintf("%s_%s", dockerComposeProjectName, serviceRole)
}

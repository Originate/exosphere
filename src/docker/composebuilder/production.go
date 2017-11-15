package composebuilder

import (
	"fmt"
	"path"

	"github.com/Originate/exosphere/src/config"
	"github.com/Originate/exosphere/src/types"
)

// ProductionDockerComposeBuilder contains the docker-compose.yml config for a single service
type ProductionDockerComposeBuilder struct {
	ServiceData       types.ServiceData
	BuiltDependencies map[string]config.AppProductionDependency
	Role              string
	AppDir            string
}

// NewProductionDockerComposeBuilder is ProductionDockerComposeBuilder's constructor
func NewProductionDockerComposeBuilder(appConfig types.AppConfig, serviceConfig types.ServiceConfig, serviceData types.ServiceData, role, appDir string) *ProductionDockerComposeBuilder {
	return &ProductionDockerComposeBuilder{
		ServiceData:       serviceData,
		BuiltDependencies: config.GetBuiltServiceProductionDependencies(serviceConfig, appConfig, appDir),
		Role:              role,
		AppDir:            appDir,
	}
}

// getServiceDockerConfigs returns a DockerConfig object for a single service and its dependencies (if any(
func (d *ProductionDockerComposeBuilder) getServiceDockerCompose() (*types.DockerCompose, error) {
	if d.ServiceData.Location != "" {
		return d.getInternalServiceDockerCompose()
	}
	if d.ServiceData.DockerImage != "" {
		return d.getExternalServiceDockerCompose()
	}
	return nil, fmt.Errorf("No location or docker image listed for '%s'", d.Role)
}

func (d *ProductionDockerComposeBuilder) getInternalServiceDockerCompose() (*types.DockerCompose, error) {
	result := types.NewDockerCompose()
	result.Services[d.Role] = types.DockerConfig{
		Build: map[string]string{
			"context":    path.Join(d.AppDir, d.ServiceData.Location),
			"dockerfile": "Dockerfile.prod",
		},
	}
	dependencyDockerCompose, err := d.getServiceDependenciesDockerCompose()
	if err != nil {
		return result, err
	}
	return result.Merge(dependencyDockerCompose), nil
}

func (d *ProductionDockerComposeBuilder) getExternalServiceDockerCompose() (*types.DockerCompose, error) {
	result := types.NewDockerCompose()
	result.Services[d.Role] = types.DockerConfig{
		Image: d.ServiceData.DockerImage,
	}
	return result, nil
}

// returns the DockerConfigs object for a service's dependencies
func (d *ProductionDockerComposeBuilder) getServiceDependenciesDockerCompose() (*types.DockerCompose, error) {
	result := types.NewDockerCompose()
	for _, builtDependency := range d.BuiltDependencies {
		if builtDependency.HasDockerConfig() {
			dockerConfig, err := builtDependency.GetDockerConfig()
			if err != nil {
				return result, err
			}
			result.Services[builtDependency.GetServiceName()] = dockerConfig
		}
	}
	return result, nil
}

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
	Production        bool
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
func (d *ProductionDockerComposeBuilder) getServiceDockerConfigs() (types.DockerConfigs, error) {
	if d.ServiceData.Location != "" {
		return d.getInternalServiceDockerConfigs()
	}
	if d.ServiceData.DockerImage != "" {
		return d.getExternalServiceDockerConfigs()
	}
	return types.DockerConfigs{}, fmt.Errorf("No location or docker image listed for '%s'", d.Role)
}

func (d *ProductionDockerComposeBuilder) getInternalServiceDockerConfigs() (types.DockerConfigs, error) {
	result := types.DockerConfigs{}
	result[d.Role] = types.DockerConfig{
		Build: map[string]string{
			"context":    path.Join(d.AppDir, d.ServiceData.Location),
			"dockerfile": "Dockerfile.prod",
		},
	}
	dependencyDockerConfigs, err := d.getServiceDependenciesDockerConfigs()
	if err != nil {
		return result, err
	}
	return result.Merge(dependencyDockerConfigs), nil
}

func (d *ProductionDockerComposeBuilder) getExternalServiceDockerConfigs() (types.DockerConfigs, error) {
	result := types.DockerConfigs{}
	result[d.Role] = types.DockerConfig{
		Image: d.ServiceData.DockerImage,
	}
	return result, nil
}

// returns the DockerConfigs object for a service's dependencies
func (d *ProductionDockerComposeBuilder) getServiceDependenciesDockerConfigs() (types.DockerConfigs, error) {
	result := types.DockerConfigs{}
	for _, builtDependency := range d.BuiltDependencies {
		dockerConfig, err := builtDependency.GetDockerConfig()
		if err != nil {
			return result, err
		}
		result[builtDependency.GetServiceName()] = dockerConfig
	}
	return result, nil
}

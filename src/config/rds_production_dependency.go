package config

import (
	"fmt"

	"github.com/Originate/exosphere/src/types"
)

type rdsProductionDependency struct {
	config    types.ProductionDependencyConfig
	appConfig types.AppConfig
}

// GetDockerConfig returns docker configuration and an error if any
func (r *rdsProductionDependency) GetDockerConfig() (types.DockerConfig, error) {
	return types.DockerConfig{
		Image: fmt.Sprintf("%s:%s", r.config.Name, r.config.Version),
	}, nil
}

// GetServiceName returns the name used as the key of this dependency in docker-compose.yml
func (r *rdsProductionDependency) GetServiceName() string {
	return r.config.Name + r.config.Version
}

//GetDeploymentConfig returns configuration needed in deployment
func (r *rdsProductionDependency) GetDeploymentConfig() (map[string]string, error) {
	config := map[string]string{
		"engine":           r.config.Name,
		"engineVersion":    r.config.Version,
		"allocatedStorage": r.config.Config.AllocatedStorage,
		"instanceClass":    r.config.Config.InstanceClass,
		"name":             r.config.Config.DbName,
		"username":         r.config.Config.Username,
		"passwordEnvVar":   r.config.Config.PasswordEnvVar,
		"storageType":      r.config.Config.StorageType,
	}
	return config, nil
}

// GetDeploymentServiceEnvVariables returns configuration needed for each service in deployment
func (r *rdsProductionDependency) GetDeploymentServiceEnvVariables() map[string]string {
	return map[string]string{}
}

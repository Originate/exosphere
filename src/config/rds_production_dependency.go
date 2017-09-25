package config

import (
	"fmt"
	"strings"

	"github.com/Originate/exosphere/src/types"
)

type rdsProductionDependency struct {
	config    types.ProductionDependencyConfig
	appConfig types.AppConfig
}

// HasDockerConfig returns a boolean indicating if a docker-compose.yml entry should be generated for the dependency
func (r *rdsProductionDependency) HasDockerConfig() bool {
	return false
}

// GetDockerConfig returns docker configuration and an error if any
func (r *rdsProductionDependency) GetDockerConfig() (types.DockerConfig, error) {
	return types.DockerConfig{}, nil
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
		"allocatedStorage": r.config.Config.Rds.AllocatedStorage,
		"instanceClass":    r.config.Config.Rds.InstanceClass,
		"name":             r.config.Config.Rds.DbName,
		"username":         r.config.Config.Rds.Username,
		"passwordEnvVar":   r.config.Config.Rds.PasswordEnvVar,
		"storageType":      r.config.Config.Rds.StorageType,
	}
	return config, nil
}

// GetDeploymentServiceEnvVariables returns configuration needed for each service in deployment
func (r *rdsProductionDependency) GetDeploymentServiceEnvVariables() map[string]string {
	return map[string]string{
		strings.ToUpper(r.config.Name): fmt.Sprintf("%s.%s.local", r.config.Config.Rds.DbName, r.appConfig.Name),
	}
}

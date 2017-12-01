package config

import (
	"fmt"
	"strings"

	"github.com/Originate/exosphere/src/types"
)

type remoteRdsDependency struct {
	config     types.RemoteDependency
	appContext *types.AppContext
}

// HasDockerConfig returns a boolean indicating if a docker-compose.yml entry should be generated for the dependency
func (r *remoteRdsDependency) HasDockerConfig() bool {
	return false
}

// GetDockerConfig returns docker configuration and an error if any
func (r *remoteRdsDependency) GetDockerConfig() (types.DockerConfig, error) {
	return types.DockerConfig{}, nil
}

// GetServiceName returns the name used as the key of this dependency in docker-compose.yml
func (r *remoteRdsDependency) GetServiceName() string {
	return r.config.Name + r.config.Version
}

//GetDeploymentConfig returns configuration needed in deployment
func (r *remoteRdsDependency) GetDeploymentConfig() (map[string]string, error) {
	config := map[string]string{
		"engine":             r.config.Name,
		"engineVersion":      r.config.Version,
		"allocatedStorage":   r.config.Config.Rds.AllocatedStorage,
		"instanceClass":      r.config.Config.Rds.InstanceClass,
		"name":               r.config.Config.Rds.DbName,
		"username":           r.config.Config.Rds.Username,
		"passwordSecretName": r.config.Config.Rds.PasswordSecretName,
		"storageType":        r.config.Config.Rds.StorageType,
	}
	return config, nil
}

// GetDeploymentServiceEnvVariables returns env vars for a service
func (r *remoteRdsDependency) GetDeploymentServiceEnvVariables(secrets types.Secrets) map[string]string {
	return map[string]string{
		strings.ToUpper(r.config.Name):                  fmt.Sprintf("%s.%s.local", r.config.Config.Rds.DbName, r.appContext.Config.Name),
		r.config.Config.Rds.ServiceEnvVarNames.DbName:   r.config.Config.Rds.DbName,
		r.config.Config.Rds.ServiceEnvVarNames.Username: r.config.Config.Rds.Username,
		r.config.Config.Rds.ServiceEnvVarNames.Password: secrets[r.config.Config.Rds.PasswordSecretName],
	}
}

// GetDeploymentVariables returns a map from string to string of variables that a dependency Terraform module needs
func (r *remoteRdsDependency) GetDeploymentVariables() (map[string]string, error) {
	return map[string]string{}, nil
}

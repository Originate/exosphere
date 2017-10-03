package types

import (
	"sort"

	"github.com/Originate/exosphere/src/util"
)

// DeployConfig contains information needed for deployment
type DeployConfig struct {
	AppConfig                AppConfig
	ServiceConfigs           map[string]ServiceConfig
	Logger                   *util.Logger
	AppDir                   string
	HomeDir                  string
	DockerComposeProjectName string
	TerraformDir             string
	SecretsPath              string
	AwsConfig                AwsConfig
}

// GetSortedServiceNames returns the service names sorted alphabetically
func (d *DeployConfig) GetSortedServiceNames() []string {
	serviceNames := []string{}
	for serviceName := range d.ServiceConfigs {
		serviceNames = append(serviceNames, serviceName)
	}
	sort.Strings(serviceNames)
	return serviceNames
}

package types

import (
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
	DeployServicesOnly       bool
}

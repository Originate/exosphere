package types

import "io"

// DeployConfig contains information needed for deployment
type DeployConfig struct {
	AppConfig                AppConfig
	ServiceConfigs           map[string]ServiceConfig
	Writer                   io.Writer
	AppDir                   string
	HomeDir                  string
	DockerComposeProjectName string
	TerraformDir             string
	SecretsPath              string
	AwsConfig                AwsConfig
	DeployServicesOnly       bool
	TerraformModulesRef      string
}

package types

import "io"

// DeployConfig contains information needed for deployment
type DeployConfig struct {
	AppContext               AppContext
	ServiceConfigs           map[string]ServiceConfig
	Writer                   io.Writer
	HomeDir                  string
	DockerComposeProjectName string
	TerraformDir             string
	SecretsPath              string
	AwsConfig                AwsConfig
	DeployServicesOnly       bool
	TerraformModulesRef      string
}

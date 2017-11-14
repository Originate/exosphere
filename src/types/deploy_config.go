package types

import "io"

// DeployConfig contains information needed for deployment
type DeployConfig struct {
	AppContext               AppContext
	ServiceConfigs           map[string]ServiceConfig
	Writer                   io.Writer
	DockerComposeProjectName string
	DockerComposeDir         string
	TerraformDir             string
	SecretsPath              string
	AwsConfig                AwsConfig
	DeployServicesOnly       bool
	TerraformModulesRef      string
}

package types

import (
	"io"

	"github.com/Originate/exosphere/src/types/context"
)

// DeployConfig contains information needed for deployment
type DeployConfig struct {
	AppContext               context.AppContext
	Writer                   io.Writer
	DockerComposeProjectName string
	DockerComposeDir         string
	TerraformDir             string
	SecretsPath              string
	AwsConfig                AwsConfig
	DeployServicesOnly       bool
	TerraformModulesRef      string
}

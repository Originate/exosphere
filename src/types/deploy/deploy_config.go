package deploy

import (
	"io"

	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/context"
)

// Config contains information needed for deployment
type Config struct {
	AppContext               context.AppContext
	ServiceConfigs           map[string]types.ServiceConfig
	Writer                   io.Writer
	DockerComposeProjectName string
	DockerComposeDir         string
	TerraformDir             string
	SecretsPath              string
	AwsConfig                types.AwsConfig
	DeployServicesOnly       bool
	TerraformModulesRef      string
}

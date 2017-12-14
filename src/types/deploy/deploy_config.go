package deploy

import (
	"io"
	"path/filepath"

	"github.com/Originate/exosphere/src/docker/composebuilder"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/context"
)

// Config contains information needed for deployment
type Config struct {
	AppContext  *context.AppContext
	Writer      io.Writer
	AwsConfig   types.AwsConfig
	AutoApprove bool
}

// GetDockerComposeProjectName returns the docker compose project name
func (c Config) GetDockerComposeProjectName() string {
	return composebuilder.GetDockerComposeProjectName(c.AppContext.Config.Name)
}

// GetBuildMode returns the build mode
func (c Config) GetBuildMode() types.BuildMode {
	return types.BuildMode{
		Type:        types.BuildModeTypeDeploy,
		Environment: types.BuildModeEnvironmentProduction,
	}
}

// GetSecretsPath returns the path to the terraform secrets file
func (c Config) GetSecretsPath() string {
	return filepath.Join(c.AppContext.GetTerraformDir(), "secrets.tfvars")
}

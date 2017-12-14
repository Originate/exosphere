package deploy

import (
	"io"
	"path"
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
	RemoteID    string
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

// GetTerraformDir returns the path to the terraform directory
func (c Config) GetTerraformDir() string {
	return path.Join(c.AppContext.Location, c.GetRelativeTerraformDir())
}

// GetRelativeTerraformDir returns the relative file path to the terraform directory
func (c Config) GetRelativeTerraformDir() string {
	return path.Join("terraform", c.RemoteID)
}

// GetTerraformSecretsPath returns the path to the terraform secrets file
func (c Config) GetTerraformSecretsPath() string {
	return filepath.Join(c.GetTerraformDir(), "secrets.tfvars")
}

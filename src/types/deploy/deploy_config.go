package deploy

import (
	"io"
	"path"

	"github.com/Originate/exosphere/src/docker/composebuilder"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/context"
)

// Config contains information needed for deployment
type Config struct {
	AppContext          *context.AppContext
	Writer              io.Writer
	AwsConfig           types.AwsConfig
	AutoApprove         bool
	RemoteEnvironmentID string
}

// GetDockerComposeProjectName returns the docker compose project name
func (c Config) GetDockerComposeProjectName() string {
	return composebuilder.GetDockerComposeProjectName(c.AppContext.Config.Name)
}

// GetTerraformDir returns the path of the main terraform directory
func (c Config) GetTerraformDir() string {
	return path.Join(c.AppContext.Location, "terraform")
}

// GetInfrastructureTerraformDir returns the path of the infra structure terraform directory
func (c Config) GetInfrastructureTerraformDir() string {
	return path.Join(c.GetTerraformDir(), "infrastructure")
}

// GetServicesTerraformDir returns the path of the services terraform directory
func (c Config) GetServicesTerraformDir() string {
	return path.Join(c.GetTerraformDir(), "services")
}

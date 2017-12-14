package deploy

import (
	"io"
	"path"

	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/context"
)

// Config contains information needed for deployment
type Config struct {
	AppContext               *context.AppContext
	Writer                   io.Writer
	DockerComposeProjectName string
	AwsConfig                types.AwsConfig
	AutoApprove              bool
}

// GetTerraformDir returns the path of the terraform directory
func (c Config) GetTerraformDir() string {
	return path.Join(c.AppContext.Location, c.GetRelativeTerraformDir())
}

// GetRelativeTerraformDir returns the relative path of the terraform directory
func (c Config) GetRelativeTerraformDir() string {
	return "terraform"
}

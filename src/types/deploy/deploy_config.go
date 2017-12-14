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
	BuildMode                types.BuildMode
	Writer                   io.Writer
	DockerComposeProjectName string
	SecretsPath              string
	AwsConfig                types.AwsConfig
	AutoApprove              bool
}

// GetTerraformDir returns the file path to the directory containaing the terraform files
func (c Config) GetTerraformDir() string {
	return path.Join(c.AppContext.Location, c.GetRelativeTerraformDir())
}

// GetRelativeTerraformDir returns the relative file path to the directory containaing the terraform files
func (c Config) GetRelativeTerraformDir() string {
	return "terraform"
}

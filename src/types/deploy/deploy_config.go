package deploy

import (
	"fmt"
	"io"
	"path"

	"github.com/Originate/exosphere/src/aws"
	"github.com/Originate/exosphere/src/docker/composebuilder"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/context"
)

// Config contains information needed for deployment
type Config struct {
	AppContext          *context.AppContext
	AutoApprove         bool
	AwsProfile          string
	RemoteEnvironmentID string
	Writer              io.Writer
}

// GetDockerComposeProjectName returns the docker compose project name
func (c Config) GetDockerComposeProjectName() string {
	return composebuilder.GetDockerComposeProjectName(c.AppContext.Config.Name)
}

// GetTerraformDir returns the path of the terraform directory
func (c Config) GetTerraformDir() string {
	return path.Join(c.AppContext.Location, c.GetRelativeTerraformDir())
}

// GetRelativeTerraformDir returns the relative path of the terraform directory
func (c Config) GetRelativeTerraformDir() string {
	return "terraform"
}

// GetAwsBucketName returns the aws bucket name to use for storage
func (c Config) GetAwsBucketName() string {
	return fmt.Sprintf("%s-%s-%s-terraform", c.GetRemoteEnvironment().AccountID, c.AppContext.Config.Name, c.RemoteEnvironmentID)
}

// GetAwsTerraformLockTable returns the dynamodb table name for storing terraform locks
func (c Config) GetAwsTerraformLockTable() string {
	return "TerraformLocks"
}

// GetRemoteEnvironment returns the app remote environment
func (c Config) GetRemoteEnvironment() types.AppRemoteEnvironment {
	return c.AppContext.Config.Remote.Environments[c.RemoteEnvironmentID]
}

// GetAwsOptions returns options to aws functions
func (c Config) GetAwsOptions() aws.Options {
	return aws.Options{
		Profile:            c.AwsProfile,
		Region:             c.GetRemoteEnvironment().Region,
		TerraformLockTable: c.GetAwsTerraformLockTable(),
		BucketName:         c.GetAwsBucketName(),
	}
}

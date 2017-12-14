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
	AppContext   *context.AppContext
	BuildMode    types.BuildMode
	Writer       io.Writer
	TerraformDir string
	AwsConfig    types.AwsConfig
	AutoApprove  bool
}

func (c Config) GetDockerComposeProjectName() string {
	return composebuilder.GetDockerComposeProjectName(c.AppContext.Config.Name)
}

func (c Config) GetBuildMode() types.BuildMode {
	return types.BuildMode{
		Type:        types.BuildModeTypeDeploy,
		Environment: types.BuildModeEnvironmentProduction,
	}
}

func (c Config) GetSecretsPath() string {
	return filepath.Join(c.TerraformDir, "secrets.tfvars")
}

// AppContext:               appContext,
// DockerComposeProjectName: composebuilder.GetDockerComposeProjectName(appContext.Config.Name),
// DockerComposeDir:         path.Join(appContext.Location, "docker-compose"),
// TerraformDir:             terraformDir,
// SecretsPath:              filepath.Join(terraformDir, "secrets.tfvars"),
// AwsConfig:                awsConfig,
// BuildMode: types.BuildMode{
// 	Type:        types.BuildModeTypeDeploy,
// 	Environment: types.BuildModeEnvironmentProduction,
// },

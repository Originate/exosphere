package deploy

import (
	"io"

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

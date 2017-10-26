package composerunner

import (
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/util"
)

// RunOptions are the options passed into Run
type RunOptions struct {
	DockerComposeDir         string
	DockerComposeProjectName string
	DockerConfigs            types.DockerConfigs
	Logger                   *util.Logger
}

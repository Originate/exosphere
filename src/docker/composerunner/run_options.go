package composerunner

import (
	"io"

	"github.com/Originate/exosphere/src/types"
)

// RunOptions are the options passed into Run
type RunOptions struct {
	DockerComposeDir         string
	DockerComposeProjectName string
	DockerConfigs            types.DockerConfigs
	Writer                   io.Writer
	AbortOnExit              bool
}

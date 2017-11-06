package composerunner

import (
	"io"

	"github.com/Originate/exosphere/src/types"
)

// RunOptions are the options passed into Run
type RunOptions struct {
	DockerComposeDir         string
	DockerComposeFileName    string
	DockerComposeProjectName string
	DockerConfigs            types.DockerConfigs
	DockerServiceName        string
	Writer                   io.Writer
	AbortOnExit              bool
}

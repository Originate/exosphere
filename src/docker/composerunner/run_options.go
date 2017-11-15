package composerunner

import (
	"io"

	"github.com/Originate/exosphere/src/types"
)

// RunOptions are the options passed into Run
type RunOptions struct {
	AppDir                   string
	DockerComposeDir         string
	DockerComposeFileName    string
	DockerComposeProjectName string
	DockerCompose            *types.DockerCompose
	DockerServiceName        string
	Writer                   io.Writer
	AbortOnExit              bool
}

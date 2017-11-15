package composerunner

import (
	"io"
)

// RunOptions are the options passed into Run
type RunOptions struct {
	AppDir                   string
	DockerComposeDir         string
	DockerComposeFileName    string
	DockerComposeProjectName string
	DockerServiceName        string
	Writer                   io.Writer
	AbortOnExit              bool
}

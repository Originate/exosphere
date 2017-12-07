package composerunner

import (
	"io"
)

// RunOptions are the options passed into Run
type RunOptions struct {
	DockerComposeDir      string
	DockerComposeFileName string
	DockerServiceName     string
	Writer                io.Writer
	AbortOnExit           bool
	EnvironmentVariables  map[string]string
}

package runner

import (
	"io"

	"github.com/Originate/exosphere/src/types/context"
)

// RunOptions runs the overall application
type RunOptions struct {
	AppContext               *context.AppContext
	DockerComposeFileName    string
	DockerComposeProjectName string
	Writer                   io.Writer
}

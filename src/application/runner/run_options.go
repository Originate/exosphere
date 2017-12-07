package runner

import (
	"io"

	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/context"
)

// RunOptions runs the overall application
type RunOptions struct {
	AppContext               *context.AppContext
	DockerComposeProjectName string
	Writer                   io.Writer
	BuildMode                types.BuildMode
}

package runner

import (
	"io"

	"github.com/Originate/exosphere/src/docker/composebuilder"
	"github.com/Originate/exosphere/src/types"
)

// RunOptions runs the overall application
type RunOptions struct {
	AppContext               *types.AppContext
	DockerComposeProjectName string
	Writer                   io.Writer
	BuildMode                composebuilder.BuildMode
}

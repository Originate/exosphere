package composebuilder

import (
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/context"
)

// ApplicationOptions are the options to
type ApplicationOptions struct {
	AppContext *context.AppContext
	BuildMode  types.BuildMode
}

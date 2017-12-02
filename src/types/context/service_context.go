package context

import (
	"path"

	"github.com/Originate/exosphere/src/types"
)

// ServiceContext represents the exosphere service the user is running
type ServiceContext struct {
	Config     types.ServiceConfig
	AppContext *AppContext
	Source     *types.ServiceSource
	Role       string
}

// ID returns the identifier for the ServiceContext
func (s *ServiceContext) ID() string {
	return path.Base(s.Source.Location)
}

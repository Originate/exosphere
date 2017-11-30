package types

import "path"

// ServiceContext represents the exosphere service the user is running
type ServiceContext struct {
	Role       string
	Config     ServiceConfig
	AppContext *AppContext
	Source     *ServiceSource
}

// ID returns the identifier for the ServiceContext
func (s *ServiceContext) ID() string {
	return path.Base(s.Source.Location)
}

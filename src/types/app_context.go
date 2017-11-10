package types

import (
	"path"
)

// AppContext represents the exosphere application the user is running
type AppContext struct {
	Location string
	Config   AppConfig
}

// GetServiceContext returns a ServiceContext for the service found at the given directory base
func (a AppContext) GetServiceContext(serviceDirBase string) (ServiceContext, error) {
	serviceConfig, err := NewServiceConfig(a.Location, serviceDirBase)
	return ServiceContext{
		Dir:        serviceDirBase,
		Location:   path.Join(a.Location, serviceDirBase),
		Config:     serviceConfig,
		AppContext: a,
	}, err
}

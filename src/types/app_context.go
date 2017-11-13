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
func (a AppContext) GetServiceContext(serviceLocation string) (ServiceContext, error) {
	serviceConfig, err := NewServiceConfig(serviceLocation)
	return ServiceContext{
		Dir:        path.Base(serviceLocation),
		Location:   serviceLocation,
		Config:     serviceConfig,
		AppContext: a,
	}, err
}

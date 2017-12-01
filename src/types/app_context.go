package types

import "path"

// AppContext represents the exosphere application the user is running
type AppContext struct {
	Location string
	Config   AppConfig
}

// GetServiceContext returns a ServiceContext for the service found at the given directory base
func (a *AppContext) GetServiceContext(serviceRole string) (*ServiceContext, error) {
	source := a.Config.Services[serviceRole]
	serviceConfig, err := NewServiceConfig(path.Join(a.Location, source.Location))
	return &ServiceContext{
		Role:       serviceRole,
		Config:     serviceConfig,
		AppContext: a,
		Source:     &source,
	}, err
}

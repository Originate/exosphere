package types

import "path"

// AppContext represents the exosphere application the user is running
type AppContext struct {
	Location string
	Config   AppConfig
}

// GetServiceContext returns a ServiceContext for the service found at the given directory base
func (a *AppContext) GetServiceContext(serviceRole string, serviceData ServiceData) (*ServiceContext, error) {
	serviceConfig, err := NewServiceConfig(path.Join(a.Location, serviceData.Location))
	return &ServiceContext{
		Role:       serviceRole,
		Config:     serviceConfig,
		AppContext: a,
		AppData:    &serviceData,
	}, err
}

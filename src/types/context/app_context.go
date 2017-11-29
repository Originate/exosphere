package context

import (
	"fmt"
	"path"

	"github.com/Originate/exosphere/src/types"
)

// AppContext represents the exosphere application the user is running
type AppContext struct {
	Location        string
	Config          types.AppConfig
	ServiceContexts map[string]*ServiceContext
}

// NewAppContext returns an AppContext with all the service contexts loaded
func NewAppContext(location string, config types.AppConfig) (*AppContext, error) {
	appContext := &AppContext{
		Location: location,
		Config:   config,
	}
	return appContext, appContext.initializeServiceContexts()
}

// GetServiceContextByLocation returns the service context for the given location
func (a *AppContext) GetServiceContextByLocation(location string) *ServiceContext {
	for _, serviceContext := range a.ServiceContexts {
		if location == path.Join(a.Location, serviceContext.AppData.Location) {
			return serviceContext
		}
	}
	return nil
}

func (a *AppContext) getServiceContext(serviceRole string, serviceData types.ServiceData) (*ServiceContext, error) {
	var serviceConfig types.ServiceConfig
	var err error
	if len(serviceData.Location) > 0 {
		serviceConfig, err = types.NewServiceConfig(path.Join(a.Location, serviceData.Location))
		if err != nil {
			return nil, err
		}
	} else if len(serviceData.DockerImage) > 0 {
		serviceConfig, err = getExternalServiceConfig(serviceData.DockerImage)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("No location or docker image for service: %s", serviceRole)
	}
	return &ServiceContext{
		Role:       serviceRole,
		Config:     serviceConfig,
		AppContext: a,
		AppData:    &serviceData,
	}, nil
}

func (a *AppContext) initializeServiceContexts() error {
	a.ServiceContexts = map[string]*ServiceContext{}
	for serviceRole, serviceData := range a.Config.Services {
		serviceContext, err := a.getServiceContext(serviceRole, serviceData)
		if err != nil {
			return err
		}
		a.ServiceContexts[serviceRole] = serviceContext
	}
	return nil
}

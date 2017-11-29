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
func (a *AppContext) GetServiceContextByLocation(serviceLocation string) *ServiceContext {
	for _, serviceContext := range a.ServiceContexts {
		if serviceContext.Location == serviceLocation {
			return serviceContext
		}
	}
	return nil
}

func (a *AppContext) getServiceContext(serviceData types.ServiceData) (*ServiceContext, error) {
	if len(serviceData.Location) > 0 {
		return a.getInternalServiceContext(path.Join(a.Location, serviceData.Location))
	} else if len(serviceData.DockerImage) > 0 {
		return a.getExternalServiceContext(serviceData.DockerImage)
	} else {
		return nil, fmt.Errorf("No location or docker image")
	}
}

func (a *AppContext) getInternalServiceContext(serviceLocation string) (*ServiceContext, error) {
	serviceConfig, err := types.NewServiceConfig(serviceLocation)
	return &ServiceContext{
		Dir:        path.Base(serviceLocation),
		Location:   serviceLocation,
		Config:     serviceConfig,
		AppContext: a,
	}, err
}

func (a *AppContext) getExternalServiceContext(dockerImage string) (*ServiceContext, error) {
	serviceConfig, err := getExternalServiceConfig(dockerImage)
	return &ServiceContext{
		External:   true,
		Config:     serviceConfig,
		AppContext: a,
	}, err
}

func (a *AppContext) initializeServiceContexts() error {
	a.ServiceContexts = map[string]*ServiceContext{}
	for serviceRole, serviceData := range a.Config.Services {
		serviceContext, err := a.getServiceContext(serviceData)
		if err != nil {
			return err
		}
		a.ServiceContexts[serviceRole] = serviceContext
	}
	return nil
}

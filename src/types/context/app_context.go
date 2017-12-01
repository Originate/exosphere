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
		if location == path.Join(a.Location, serviceContext.Source.Location) {
			return serviceContext
		}
	}
	return nil
}

// GetServiceEndpoints returns the service endpoints for the given build mode
func (a *AppContext) GetServiceEndpoints(buildMode types.BuildMode) map[string]*ServiceEndpoints {
	portReservation := types.NewPortReservation()
	serviceEndpoints := map[string]*ServiceEndpoints{}
	for _, serviceRole := range a.Config.GetSortedServiceRoles() {
		serviceEndpoints[serviceRole] = NewServiceEndpoint(a, serviceRole, portReservation, buildMode)
	}
	return serviceEndpoints
}

func (a *AppContext) getServiceContext(serviceRole string, serviceSource types.ServiceSource) (*ServiceContext, error) {
	var serviceConfig types.ServiceConfig
	var err error
	if len(serviceSource.Location) > 0 {
		serviceConfig, err = types.NewServiceConfig(path.Join(a.Location, serviceSource.Location))
		if err != nil {
			return nil, err
		}
	} else if len(serviceSource.DockerImage) > 0 {
		serviceConfig, err = getExternalServiceConfig(serviceSource.DockerImage)
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
		Source:     &serviceSource,
	}, nil
}

func (a *AppContext) initializeServiceContexts() error {
	a.ServiceContexts = map[string]*ServiceContext{}
	for serviceRole, serviceSource := range a.Config.Services {
		serviceContext, err := a.getServiceContext(serviceRole, serviceSource)
		if err != nil {
			return err
		}
		a.ServiceContexts[serviceRole] = serviceContext
	}
	return nil
}

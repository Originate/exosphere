package context

import (
	"fmt"
	"path"
	"sort"

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

// GetServiceContext returns the service context for the given location
func (a *AppContext) GetServiceContext(location string) *ServiceContext {
	for _, serviceContext := range a.ServiceContexts {
		if location == path.Join(a.Location, serviceContext.Source.Location) {
			return serviceContext
		}
	}
	return nil
}

// GetRemoteDependencies returns all remote dependencies
func (a *AppContext) GetRemoteDependencies() map[string]types.RemoteDependency {
	result := map[string]types.RemoteDependency{}
	for dependencyName, dependency := range a.Config.Remote.Dependencies {
		result[dependencyName] = dependency
	}
	for _, serviceContext := range a.ServiceContexts {
		for dependencyName, dependency := range serviceContext.Config.Remote.Dependencies {
			result[dependencyName] = dependency
		}
	}
	return result
}

// GetSortedRemoteDependencyNames returns all remote dependency names in alphabetical order
func (a *AppContext) GetSortedRemoteDependencyNames() []string {
	result := []string{}
	for k := range a.Config.Remote.Dependencies {
		result = append(result, k)
	}
	for _, serviceContext := range a.ServiceContexts {
		for k := range serviceContext.Config.Remote.Dependencies {
			result = append(result, k)
		}
	}
	sort.Strings(result)
	return result
}

// GetDependencyServiceData returns a map from service role to service data for the given dependency
func (a *AppContext) GetDependencyServiceData(dependencyName string) map[string]map[string]interface{} {
	result := map[string]map[string]interface{}{}
	for _, serviceRole := range a.Config.GetSortedServiceRoles() {
		result[serviceRole] = a.ServiceContexts[serviceRole].GetDependencyData(dependencyName)
	}
	return result
}

// GetDockerComposeDir returns the file path to the directory containaing the docker compose files
func (a *AppContext) GetDockerComposeDir() string {
	return path.Join(a.Location, "docker-compose")
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

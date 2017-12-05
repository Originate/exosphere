package config

import (
	"encoding/json"
	"fmt"

	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/context"
)

type localExocomDependency struct {
	config     types.LocalDependency
	appContext *context.AppContext
}

func (e *localExocomDependency) compileServiceRoutes() []map[string]interface{} {
	routes := []map[string]interface{}{}
	serviceData := e.appContext.Config.Services
	for _, serviceRole := range e.appContext.Config.GetSortedServiceRoles() {
		serviceConfig := e.appContext.ServiceContexts[serviceRole].Config
		route := map[string]interface{}{
			"role":     serviceRole,
			"receives": serviceConfig.ServiceMessages.Receives,
			"sends":    serviceConfig.ServiceMessages.Sends,
		}
		messageTranslations := serviceData[serviceRole].MessageTranslations
		if messageTranslations != nil {
			route["messageTranslations"] = messageTranslations
		}
		routes = append(routes, route)
	}
	return routes
}

// GetServiceName returns the service name
func (e *localExocomDependency) GetServiceName() string {
	return e.config.Name + e.config.Version
}

// GetDockerConfig returns docker configuration and an error if any
func (e *localExocomDependency) GetDockerConfig() (types.DockerConfig, error) {
	serviceRoutes, err := e.getServiceRoutesString()
	if err != nil {
		return types.DockerConfig{}, err
	}
	return types.DockerConfig{
		Image: fmt.Sprintf("originate/exocom:%s", e.config.Version),
		Environment: map[string]string{
			"ROLE":           "exocom",
			"SERVICE_ROUTES": serviceRoutes,
		},
		Restart: "on-failure",
	}, nil
}

// GetServiceEnvVariables returns the environment variables that need to
// be passed to services that use it
func (e *localExocomDependency) GetServiceEnvVariables() map[string]string {
	return map[string]string{
		"EXOCOM_HOST": e.GetServiceName(),
	}
}

func (e *localExocomDependency) getServiceRoutesString() (string, error) {
	serviceRoutes := e.compileServiceRoutes()
	serviceRoutesBytes, err := json.Marshal(serviceRoutes)
	if err != nil {
		return "", err
	}
	return string(serviceRoutesBytes), nil
}

// GetVolumeNames returns the named volumes used by this dependency
func (e *localExocomDependency) GetVolumeNames() []string {
	return []string{}
}

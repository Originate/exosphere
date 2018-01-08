package endpoints

import (
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/context"
	"github.com/Originate/exosphere/src/util"
)

// ServiceEndpoints holds the information to build an endpoint at which a service can be reached
type ServiceEndpoints map[string]*ServiceEndpoint

// NewServiceEndpoints returns the constructed service endpoint objects for all services in an application
func NewServiceEndpoints(appContext *context.AppContext, buildMode types.BuildMode, remoteEnvironmentID string) ServiceEndpoints {
	portReservation := NewPortReservation()
	serviceEndpoints := map[string]*ServiceEndpoint{}
	for _, serviceRole := range appContext.Config.GetSortedServiceRoles() {
		serviceEndpoints[serviceRole] = newServiceEndpoint(appContext, serviceRole, portReservation, buildMode, remoteEnvironmentID)
	}
	return serviceEndpoints
}

// GetServiceEndpointEnvVars creates all the endpoint env vars for a service
func (s ServiceEndpoints) GetServiceEndpointEnvVars(serviceRole string) map[string]string {
	endpointEnvVars := map[string]string{}
	for role, serviceEndpoint := range s {
		if role == serviceRole {
			continue
		}
		util.Merge(endpointEnvVars, serviceEndpoint.GetEndpointMappings())
	}
	return endpointEnvVars
}

// GetServicePortMappings gets the port mapping for a particular service
func (s ServiceEndpoints) GetServicePortMappings(serviceRole string) []string {
	return s[serviceRole].GetPortMappings()
}

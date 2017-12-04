package types

import (
	"github.com/Originate/exosphere/src/util"
)

// ServiceEndpoints holds the information to build an endpoint at which a service can be reached
type ServiceEndpoints struct {
	ServiceEndpoints map[string]ServiceEndpoint
}

// NewServiceEndpoints returns the constructed service endpoint objects for all services in an application
func NewServiceEndpoints(appContext AppContext, buildMode BuildMode) *ServiceEndpoint {
	portReservation := NewPortReservation()
	serviceEndpoints := map[string]ServiceEndpoint{}
	for _, serviceRole := range appContext.Config.GetSortedServiceRoles() {
		serviceConfig := appContext.ServiceContexts[serviceRole].Config
		serviceEndpoints[serviceRole] = newServiceEndpoint(serviceRole, serviceConfig, portReservation, buildMode)
	}
	return &ServiceEndpoints{
		ServiceEndpoints: serviceEndpoints,
	}
}

// GetServiceEndpointEnvVars creates all the endpoint env vars for a service
func (s *ServiceEndpoints) GetServiceEndpointEnvVars(serviceRole string) map[string]string {
	endpointEnvVars := map[string]string{}
	for role, serviceEndpoint := range s.ServiceEndpoints {
		if role == serviceRole {
			continue
		}
		util.Merge(endpointEnvVars, serviceEndpoint.GetEndpointMappings())
	}
	return endpointEnvVars
}

func (s *ServiceEndpoints) GetServicePortMappings(serviceRole string) []string {
	s.ServiceEndpoints[serviceRole].GetPortMappings()
}

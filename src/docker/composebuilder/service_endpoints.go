package composebuilder

import (
	"fmt"
	"strings"

	"github.com/Originate/exosphere/src/types"
)

// ServiceEndpoints holds the information to build an endpoint at which a service can be reached
type ServiceEndpoints struct {
	ServiceRole   string
	ContainerPort string
	HostPort      string
}

// NewServiceEndpoint initializes a ServiceEndpoint struct
func NewServiceEndpoint(serviceRole string, serviceConfig types.ServiceConfig, portReservation *types.PortReservation, buildMode BuildMode) *ServiceEndpoints {
	containerPort := ""
	switch buildMode.Environment {
	case BuildModeEnvironmentDevelopment:
		containerPort = serviceConfig.Development.Port
	case BuildModeEnvironmentProduction:
		containerPort = serviceConfig.Production.Port
	}
	hostPort := ""
	if containerPort != "" {
		hostPort = portReservation.GetAvailablePort()
	}
	return &ServiceEndpoints{
		ServiceRole:   serviceRole,
		ContainerPort: containerPort,
		HostPort:      hostPort,
	}
}

// GetPortMappings returns a list of port mappings from a service's container ports to host ports
func (s *ServiceEndpoints) GetPortMappings() []string {
	if s.ContainerPort == "" {
		return []string{}
	}
	return []string{fmt.Sprintf("%s:%s", s.HostPort, s.ContainerPort)}
}

// GetEndpointMappings returns a map from env var name to env var value of a service endpoint
func (s *ServiceEndpoints) GetEndpointMappings() map[string]string {
	if s.ContainerPort == "" {
		return map[string]string{}
	}
	externalKey := fmt.Sprintf("SERVICE_%s_EXTERNAL_ORIGIN", strings.ToUpper(s.ServiceRole))
	externalValue := fmt.Sprintf("http://localhost:%s", s.HostPort)
	return map[string]string{externalKey: externalValue}
}

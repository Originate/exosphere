package composebuilder

import (
	"fmt"
	"strings"

	"github.com/Originate/exosphere/src/types"
)

// ServiceEndpoints holds the information to build an endpoint at which a service can be reached
type ServiceEndpoints struct {
	ServiceRole   string
	ServiceConfig types.ServiceConfig
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
		ServiceConfig: serviceConfig,
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
	switch s.ServiceConfig.Type {
	case "public":
		externalKey := fmt.Sprintf("%s_EXTERNAL_ORIGIN", toConstantCase(s.ServiceRole))
		externalValue := fmt.Sprintf("http://localhost:%s", s.HostPort)
		return map[string]string{externalKey: externalValue}
	default:
		return map[string]string{}
	}
}

// converts valid serviceRole strings to constant case
// see validateAppConfig() in types/app_config.go for valid serviceRole regex
func toConstantCase(serviceRole string) string {
	str := strings.Join(strings.Split(serviceRole, "-"), "_")
	return strings.ToUpper(str)
}

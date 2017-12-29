package endpoints

import (
	"fmt"
	"strings"

	"github.com/Originate/exosphere/src/types"
)

// ServiceEndpoint holds the information to build an endpoint at which a service can be reached
type ServiceEndpoint struct {
	ServiceRole         string
	ServiceConfig       types.ServiceConfig
	ContainerPort       string
	HostPort            string
	BuildMode           types.BuildMode
	RemoteEnvironmentID string
}

func newServiceEndpoint(serviceRole string, serviceConfig types.ServiceConfig, portReservation *PortReservation, buildMode types.BuildMode, remoteEnvironmentID string) *ServiceEndpoint {
	containerPort := ""
	hostPort := ""
	switch buildMode.Environment {
	case types.BuildModeEnvironmentDevelopment:
		containerPort = serviceConfig.Development.Port
	case types.BuildModeEnvironmentProduction:
		containerPort = serviceConfig.Production.Port
	}
	if buildMode.Type == types.BuildModeTypeLocal && containerPort != "" {
		hostPort = portReservation.GetAvailablePort()
	}
	return &ServiceEndpoint{
		ServiceRole:         serviceRole,
		ServiceConfig:       serviceConfig,
		ContainerPort:       containerPort,
		HostPort:            hostPort,
		BuildMode:           buildMode,
		RemoteEnvironmentID: remoteEnvironmentID,
	}
}

// GetPortMappings returns a list of port mappings from a service's container ports to host ports
func (s *ServiceEndpoint) GetPortMappings() []string {
	if s.ContainerPort == "" {
		return []string{}
	}
	return []string{fmt.Sprintf("%s:%s", s.HostPort, s.ContainerPort)}
}

// GetEndpointMappings returns a map from env var name to env var value of a service endpoint
func (s *ServiceEndpoint) GetEndpointMappings() map[string]string {
	switch s.ServiceConfig.Type {
	case "public":
		return s.getPublicEndpoints()
	case "worker":
		return s.getWorkerEndpoints()
	default:
		return map[string]string{}
	}
}

func (s *ServiceEndpoint) getPublicEndpoints() map[string]string {
	externalKey := fmt.Sprintf("%s_EXTERNAL_ORIGIN", toConstantCase(s.ServiceRole))
	internalKey := fmt.Sprintf("%s_INTERNAL_ORIGIN", toConstantCase(s.ServiceRole))
	endpoints := map[string]string{}
	if s.ContainerPort != "" {
		if s.BuildMode.Type == types.BuildModeTypeLocal {
			endpoints[externalKey] = fmt.Sprintf("http://localhost:%s", s.HostPort)
			endpoints[internalKey] = fmt.Sprintf("http://%s:%s", s.ServiceRole, s.ContainerPort)
		} else {
			endpoints[externalKey] = fmt.Sprintf("https://%s", s.ServiceConfig.Remote.Environments[s.RemoteEnvironmentID].URL)
			endpoints[internalKey] = fmt.Sprintf("http://%s.local", s.ServiceRole)
		}
	}
	return endpoints
}

func (s *ServiceEndpoint) getWorkerEndpoints() map[string]string {
	host := fmt.Sprintf("%s_HOST", toConstantCase(s.ServiceRole))
	endpoints := map[string]string{}
	if s.ContainerPort != "" {
		if s.BuildMode.Type == types.BuildModeTypeLocal {
			endpoints[host] = fmt.Sprintf("http://%s:%s", s.ServiceRole, s.ContainerPort)
		}
	}
	return endpoints
}

// converts valid serviceRole strings to constant case
// see validateAppConfig() in types/app_config.go for valid serviceRole regex
func toConstantCase(serviceRole string) string {
	str := strings.Join(strings.Split(serviceRole, "-"), "_")
	return strings.ToUpper(str)
}

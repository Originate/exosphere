package endpoints

import (
	"fmt"
	"strings"

	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/util"
)

// ServiceEndpoint holds the information to build an endpoint at which a service can be reached
type ServiceEndpoint struct {
	ServiceRole   string
	ServiceConfig types.ServiceConfig
	ContainerPort string
	HostPort      string
	BuildMode     types.BuildMode
}

func newServiceEndpoint(serviceRole string, serviceConfig types.ServiceConfig, portReservation *PortReservation, buildMode types.BuildMode) *ServiceEndpoint {
	containerPort := ""
	hostPort := ""
	if buildMode.Type == types.BuildModeTypeLocal {
		switch buildMode.Environment {
		case types.BuildModeEnvironmentDevelopment:
			containerPort = serviceConfig.Development.Port
		case types.BuildModeEnvironmentProduction:
			containerPort = serviceConfig.Production.Port
		}
		if containerPort != "" {
			hostPort = portReservation.GetAvailablePort()
		}
	}
	return &ServiceEndpoint{
		ServiceRole:   serviceRole,
		ServiceConfig: serviceConfig,
		ContainerPort: containerPort,
		HostPort:      hostPort,
		BuildMode:     buildMode,
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
		externalOrigin := s.getExternalOrigin()
		internalOrigin := s.getInternalOrigin()
		util.Merge(externalOrigin, internalOrigin)
		return externalOrigin
	default:
		return map[string]string{}
	}
}

func (s *ServiceEndpoint) getExternalOrigin() map[string]string {
	externalKey := fmt.Sprintf("%s_EXTERNAL_ORIGIN", toConstantCase(s.ServiceRole))
	if s.BuildMode.Type == types.BuildModeTypeLocal {
		if s.HostPort != "" {
			return map[string]string{externalKey: fmt.Sprintf("http://localhost:%s", s.HostPort)}
		}
	} else {
		return map[string]string{externalKey: fmt.Sprintf("https://%s", s.ServiceConfig.Remote.URL)}
	}
	return map[string]string{}
}

func (s *ServiceEndpoint) getInternalOrigin() map[string]string {
	externalKey := fmt.Sprintf("%s_ORIGIN", toConstantCase(s.ServiceRole))
	if s.BuildMode.Type == types.BuildModeTypeLocal {
		if s.ContainerPort != "" {
			return map[string]string{externalKey: fmt.Sprintf("http://localhost:%s", s.ContainerPort)}
		}
	} else {
		return map[string]string{externalKey: fmt.Sprintf("http://%s.local", s.ServiceRole)}
	}
	return map[string]string{}

}

// converts valid serviceRole strings to constant case
// see validateAppConfig() in types/app_config.go for valid serviceRole regex
func toConstantCase(serviceRole string) string {
	str := strings.Join(strings.Split(serviceRole, "-"), "_")
	return strings.ToUpper(str)
}

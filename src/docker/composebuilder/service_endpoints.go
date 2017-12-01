package composebuilder

//TODO types refactor: move to types package

import (
	"fmt"
	"strings"

	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/context"
)

// ServiceEndpoints holds the information to build an endpoint at which a service can be reached
type ServiceEndpoints struct {
	AppContext    *context.AppContext
	ServiceRole   string
	ServiceConfig types.ServiceConfig
	ContainerPort string
	HostPort      string
	BuildMode     BuildMode
}

// NewServiceEndpoint initializes a ServiceEndpoint struct
func NewServiceEndpoint(appContext *context.AppContext, serviceRole string, portReservation *types.PortReservation, buildMode BuildMode) *ServiceEndpoints {
	containerPort := ""
	hostPort := ""
	serviceConfig := appContext.ServiceContexts[serviceRole].Config
	if buildMode.Type == BuildModeTypeLocal {
		switch buildMode.Environment {
		case BuildModeEnvironmentDevelopment:
			containerPort = serviceConfig.Development.Port
		case BuildModeEnvironmentProduction:
			containerPort = serviceConfig.Production.Port
		}
		if containerPort != "" {
			hostPort = portReservation.GetAvailablePort()
		}
	}
	return &ServiceEndpoints{
		AppContext:    appContext,
		ServiceRole:   serviceRole,
		ServiceConfig: serviceConfig,
		ContainerPort: containerPort,
		HostPort:      hostPort,
		BuildMode:     buildMode,
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
	switch s.ServiceConfig.Type {
	case "public":
		externalKey := fmt.Sprintf("%s_EXTERNAL_ORIGIN", toConstantCase(s.ServiceRole))
		var externalValue string
		if s.BuildMode.Type == BuildModeTypeLocal {
			if s.HostPort != "" {
				externalValue = fmt.Sprintf("http://localhost:%s", s.HostPort)
			}
		} else {
			externalValue = fmt.Sprintf("https://%s", s.ServiceConfig.Production.URL)
		}
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

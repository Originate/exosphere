package context

import (
	"path"

	"github.com/Originate/exosphere/src/types"
)

// ServiceContext represents the exosphere service the user is running
type ServiceContext struct {
	Config     types.ServiceConfig
	AppContext *AppContext
	Source     *types.ServiceSource
	Role       string
}

// ID returns the identifier for the ServiceContext
func (s *ServiceContext) ID() string {
	return path.Base(s.Source.Location)
}

// GetDependencyData returns the data to pass to a dependency for theis service
func (s *ServiceContext) GetDependencyData(dependencyName string) map[string]interface{} {
	result := map[string]interface{}{}
	serviceDependencyDataList := []types.ServiceDependencyData{
		s.Config.DependencyData,
		s.Source.DependencyData,
	}
	for _, serviceDependencyData := range serviceDependencyDataList {
		for name, dependencyData := range serviceDependencyData {
			if name == dependencyName {
				for key, value := range dependencyData {
					result[key] = value
				}
			}
		}
	}
	return result
}

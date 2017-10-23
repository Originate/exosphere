package composebuilder

import (
	"github.com/Originate/exosphere/src/config"
	"github.com/Originate/exosphere/src/types"
)

// GetApplicationDockerConfigs returns the DockerConfigs for a application
func GetApplicationDockerConfigs(options ApplicationOptions) (types.DockerConfigs, error) {
	dependencyDockerConfigs, err := getDependenciesDockerConfigs(options)
	if err != nil {
		return nil, err
	}
	serviceDockerConfigs, err := getServicesDockerConfigs(options)
	if err != nil {
		return nil, err
	}
	return dependencyDockerConfigs.Merge(serviceDockerConfigs), nil
}

// Helpers

func getDependenciesDockerConfigs(options ApplicationOptions) (types.DockerConfigs, error) {
	result := types.DockerConfigs{}
	if options.BuildMode == BuildModeDeployProduction {
		appDependencies := config.GetBuiltAppProductionDependencies(options.AppConfig, options.AppDir)
		for _, builtDependency := range appDependencies {
			if builtDependency.HasDockerConfig() {
				dockerConfig, err := builtDependency.GetDockerConfig()
				if err != nil {
					return result, err
				}
				result[builtDependency.GetServiceName()] = dockerConfig
			}
		}
	} else {
		appDependencies := config.GetBuiltAppDevelopmentDependencies(options.AppConfig, options.AppDir, options.HomeDir)
		for _, builtDependency := range appDependencies {
			dockerConfig, err := builtDependency.GetDockerConfig()
			if err != nil {
				return result, err
			}
			result[builtDependency.GetContainerName()] = dockerConfig
		}
	}
	return result, nil
}

func getServicesDockerConfigs(options ApplicationOptions) (types.DockerConfigs, error) {
	result := types.DockerConfigs{}
	serviceConfigs, err := config.GetServiceConfigs(options.AppDir, options.AppConfig)
	if err != nil {
		return result, err
	}
	serviceData := options.AppConfig.GetServiceData()
	for serviceRole, serviceConfig := range serviceConfigs {
		dockerConfig, err := GetServiceDockerConfigs(options.AppConfig, serviceConfig, serviceData[serviceRole], serviceRole, options.AppDir, options.HomeDir, options.BuildMode)
		if err != nil {
			return result, err
		}
		result = result.Merge(result, dockerConfig)
	}
	return result, nil
}

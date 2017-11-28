package config

import (
	"github.com/Originate/exosphere/src/types"
)

// GetDependencyDataMap returns a map from dependencyName to service data
func getDependencyDataMap(appDir string, appConfig types.AppConfig) (map[string]map[string]interface{}, error) {
	serviceConfigs, err := GetServiceConfigs(appDir, appConfig)
	if err != nil {
		return nil, err
	}
	dataMap := map[string]map[string]interface{}{}
	for _, serviceRole := range appConfig.GetSortedServiceRoles() {
		for dependencyName, serviceData := range serviceConfigs[serviceRole].DependencyData {
			if _, ok := dataMap[dependencyName]; !ok {
				dataMap[dependencyName] = map[string]interface{}{}
			}
			for key, value := range serviceData {
				dataMap[dependencyName][key] = value
			}
		}
		for dependencyName, serviceData := range appConfig.Services[serviceRole].DependencyData {
			if _, ok := dataMap[dependencyName]; !ok {
				dataMap[dependencyName] = map[string]interface{}{}
			}
			for key, value := range serviceData {
				dataMap[dependencyName][key] = value
			}
		}
	}
	return dataMap, nil
}

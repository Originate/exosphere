package config

import "fmt"

// GetServiceConfigs reads the service.yml of all services and returns
// the serviceConfig objects and an error (if any)
func GetServiceConfigs(appDir string, appConfig AppConfig) (map[string]ServiceConfig, error) {
	result := map[string]ServiceConfig{}
	for service, serviceData := range appConfig.GetServiceData() {
		if len(serviceData.Location) > 0 {
			serviceConfig, err := getInternalServiceConfig(appDir, serviceData.Location)
			if err != nil {
				return result, err
			}
			result[service] = serviceConfig
		} else if len(serviceData.DockerImage) > 0 {
			serviceConfig, err := getExternalServiceConfig(service, serviceData)
			if err != nil {
				return result, err
			}
			result[service] = serviceConfig
		} else {
			return result, fmt.Errorf("No location or docker image listed for %s", service)
		}
	}
	return result, nil
}

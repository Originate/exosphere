package config

// GetInternalServiceConfigs reads the service.yml of all internal services and returns
// the serviceConfig objects and an error (if any)
func GetInternalServiceConfigs(appDir string, appConfig AppConfig) (map[string]ServiceConfig, error) {
	result := map[string]ServiceConfig{}
	for service, serviceData := range GetServiceData(appConfig.Services) {
		if len(serviceData.Location) > 0 {
			serviceConfig, err := getInternalServiceConfig(appDir, serviceData.Location)
			if err != nil {
				return result, err
			}
			result[service] = serviceConfig
		}
	}
	return result, nil
}

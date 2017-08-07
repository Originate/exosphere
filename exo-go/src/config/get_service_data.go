package config

// GetServiceData returns the configurations data for the given services
func GetServiceData(services Services) map[string]ServiceData {
	result := make(map[string]ServiceData)
	for serviceName, data := range services.Private {
		result[serviceName] = data
	}
	for serviceName, data := range services.Public {
		result[serviceName] = data
	}
	return result
}

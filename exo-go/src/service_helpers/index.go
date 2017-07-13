package serviceHelpers

import (
	"fmt"

	"github.com/Originate/exosphere/exo-go/src/types"
	"github.com/Originate/exosphere/exo-go/src/util"
)

// VerifyServiceDoesNotExist returns an error if the service serviceRole already
// exists in existingServices, and return nil otherwise.
func VerifyServiceDoesNotExist(serviceRole string, existingServices []string) error {
	if util.DoesStringArrayContain(existingServices, serviceRole) {
		return fmt.Errorf(`Service %v already exists in this application`, serviceRole)
	}
	return nil
}

// GetExistingServices returns a slice of all service names in the application
func GetExistingServices(services types.Services) []string {
	existingServices := []string{}
	for service := range services.Private {
		existingServices = append(existingServices, service)
	}
	for service := range services.Public {
		existingServices = append(existingServices, service)
	}
	return existingServices
}

// GetServiceData returns a map of service role to its info listed in application.yml
// i.e. [role] -> {serviceLocation, dockerImage, namespace}
func GetServiceData(services types.Services) map[string]types.ServiceData {
	serviceDataMap := make(map[string]types.ServiceData)
	for service, serviceData := range services.Private {
		serviceDataMap[service] = serviceData
	}
	for service, serviceData := range services.Public {
		serviceDataMap[service] = serviceData
	}
	return serviceDataMap
}

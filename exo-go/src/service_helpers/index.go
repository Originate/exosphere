package serviceHelpers

import (
	"fmt"
	"os"

	"github.com/Originate/exosphere/exo-go/src/types"
	"github.com/Originate/exosphere/exo-go/src/util"
)

// VerifyServiceDoesNotExist forces the program to exit with status code 1
// if the service serviceRole already exists in existingServices
func VerifyServiceDoesNotExist(serviceRole string, existingServices []string) {
	if util.DoesStringArrayContain(existingServices, serviceRole) {
		fmt.Printf(`Service %v already exists in this application`, serviceRole)
		os.Exit(1)
	}
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

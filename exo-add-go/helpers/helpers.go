package helpers

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/Originate/exosphere/exo-add-go/types"
)

func Contains(strings []string, targetString string) bool {
	for _, element := range strings {
		if element == targetString {
			return true
		}
	}
	return false
}

func CheckForService(serviceRole string, existingServices []string) {
	if Contains(existingServices, serviceRole) {
		fmt.Printf("Service %v already exists in this application\n", serviceRole)
		os.Exit(1)
	}
}

func GetExistingServices(services map[string]map[string]types.Service) []string {
	existingServices := []string{}
	for _, serviceConfigs := range services {
		for service := range serviceConfigs {
			existingServices = append(existingServices, service)
		}
	}
	return existingServices
}

func GetSubdirectories(directory string) []string {
	subDirectories := []string{}
	entries, _ := ioutil.ReadDir(directory)
  for _, entry := range entries {
  	if isDirectory(path.Join(directory, entry.Name())) {
  		subDirectories = append(subDirectories, entry.Name())
  	}
  }
  return subDirectories
}

func isDirectory(entry string) bool {
  f, err := os.Stat(entry)
  if err != nil {
    panic(err)
  }
  return f.IsDir()
}

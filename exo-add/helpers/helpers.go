package helpers

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/Originate/exosphere/exo-add/os_helpers"
	"github.com/Originate/exosphere/exo-add/types"
	"github.com/Originate/exosphere/exo-add/user_input"
	"github.com/tmrts/boilr/pkg/template"
	"gopkg.in/yaml.v2"
)

// VerifyServiceDoesNotExist forces the program to exit with status code 1
// if the service serviceRole already exists in existingServices
func VerifyServiceDoesNotExist(serviceRole string, existingServices []string) {
	if contains(existingServices, serviceRole) {
		fmt.Printf(`Service %v already exists in this application`, serviceRole)
		os.Exit(1)
	}
}

// GetExistingServices returns a slice of all service names in the application
func GetExistingServices(services types.Services) []string {
	existingServices := []string{}
	for service, _ := range services.Private {
		existingServices = append(existingServices, service)
	}
	for service, _ := range services.Public {
		existingServices = append(existingServices, service)
	}
	return existingServices
}

// GetTemplates returns a slice of all template names found in the ".exosphere"
// folder of the application
func GetTemplates() []string {
	templatesDir := ".exosphere"
	if !osHelpers.DirectoryExists(templatesDir) || osHelpers.IsEmpty(templatesDir) {
		fmt.Println("no templates found\n\nPlease add templates to the \".exosphere\" folder of your code base.")
		os.Exit(1)
	}
	templates := []string{}
	for _, directory := range osHelpers.GetSubdirectories(templatesDir) {
		if osHelpers.IsValidTemplateDir(path.Join(templatesDir, directory)) {
			templates = append(templates, directory)
		}
	}
	return templates
}

// CreateTmpServiceDir makes bolir scaffold the template chosenTemplate
// and store the scaffoled service folder in a tmp folder, and finally
// returns the path to the tmp folder
func CreateTmpServiceDir(chosenTemplate string) string {
	templateDir := path.Join(".exosphere", chosenTemplate)
	template, err := template.Get(templateDir)
	if err != nil {
		log.Fatalf("Failed to fetch %s template: %s", chosenTemplate, err)
	}
	serviceTmpDir, err := ioutil.TempDir("", "service-tmp")
	if err != nil {
		log.Fatalf(`Failed to create a tmp folder for the service "%s": %s`, chosenTemplate, err)
	}
	if err = template.Execute(serviceTmpDir); err != nil {
		log.Fatalf(`Failed to create the service "%s": %s`, chosenTemplate, err)
	}
	return serviceTmpDir
}

// GetAppConfig reads application.yml and returns the appConfig object
func GetAppConfig() types.AppConfig {
	yamlFile, err := ioutil.ReadFile("application.yml")
	var appConfig types.AppConfig
	err = yaml.Unmarshal(yamlFile, &appConfig)
	if err != nil {
		log.Fatalf("Failed to read application.yml: %s", err)
	}
	return appConfig
}

// UpdateAppConfig adds serviceRole to the appConfig object and updates
// application.yml
func UpdateAppConfig(serviceRole string, appConfig types.AppConfig) {
	reader := bufio.NewReader(os.Stdin)
	switch protectionLevel := userInput.Choose(reader, "Protection Level:", []string{"public", "private"}); protectionLevel {
	case "public":
		if appConfig.Services.Public == nil {
			appConfig.Services.Public = make(map[string]types.ServiceConfig)
		}
		appConfig.Services.Public[serviceRole] = types.ServiceConfig{Location: fmt.Sprintf("./%s", serviceRole)}
	case "private":
		if appConfig.Services.Private == nil {
			appConfig.Services.Private = make(map[string]types.ServiceConfig)
		}
		appConfig.Services.Private[serviceRole] = types.ServiceConfig{Location: fmt.Sprintf("./%s", serviceRole)}
	}
	bytes, _ := yaml.Marshal(appConfig)
	ioutil.WriteFile(path.Join("application.yml"), bytes, 0777)
}

func contains(strings []string, targetString string) bool {
	for _, element := range strings {
		if element == targetString {
			return true
		}
	}
	return false
}

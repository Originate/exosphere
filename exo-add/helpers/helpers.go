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

func CheckForService(serviceRole string, existingServices []string) {
	if contains(existingServices, serviceRole) {
		log.Fatalf(`Service "%v" already exists in this application\n`, serviceRole)
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

func GetTemplates() []string {
	templatesDir := ".exosphere"
	if !osHelpers.DirectoryExists(templatesDir) || osHelpers.IsEmpty(templatesDir) {
		log.Fatal("no templates found")
	}
	templates := []string{}
	for _, directory := range osHelpers.GetSubdirectories(templatesDir) {
		if osHelpers.IsValidTemplateDir(path.Join(templatesDir, directory)) {
			templates = append(templates, directory)
		}
	}
	return templates
}

func CreateTmpServiceDir(chosenTemplate string) string {
	templateDir := path.Join(".exosphere", chosenTemplate)
	template, err := template.Get(templateDir)
	if err != nil {
		log.Fatalf("Failed to fetch %s template: %s", chosenTemplate, err)
	}
	serviceTmpDir := path.Join("tmp", "service-tmp")
	if err = createServiceTmpDir(); err != nil {
		log.Fatalf(`Failed to create a tmp folder for the service "%s": %s`, chosenTemplate, err)
	}
	if err = template.Execute(serviceTmpDir); err != nil {
		log.Fatalf(`Failed to create the service "%s": %s`, chosenTemplate, err)
	}
	return serviceTmpDir
}

func GetAppConfig() types.AppConfig {
	yamlFile, err := ioutil.ReadFile("application.yml")
	var appConfig types.AppConfig
	err = yaml.Unmarshal(yamlFile, &appConfig)
	if err != nil {
		log.Fatalf("Failed to read application.yml: %s", err)
	}
	return appConfig
}

func UpdateAppConfig(serviceRole string, appConfig types.AppConfig) {
	reader := bufio.NewReader(os.Stdin)
	protectionLevel := userInput.Choose(reader, "Protection Level:", []string{"public", "private"})
	if len(appConfig.Services[protectionLevel]) == 0 {
		appConfig.Services[protectionLevel] = make(map[string]types.Service)
	}
	appConfig.Services[protectionLevel][serviceRole] = types.Service{fmt.Sprintf("./%s", serviceRole)}
	bytes, _ := yaml.Marshal(appConfig)
	ioutil.WriteFile(path.Join("application.yml"), bytes, 0777)
}

func createServiceTmpDir() error {
	serviceTmpDir := path.Join("tmp", "service-tmp")
	return os.MkdirAll(serviceTmpDir, 0777)
}

func contains(strings []string, targetString string) bool {
	for _, element := range strings {
		if element == targetString {
			return true
		}
	}
	return false
}

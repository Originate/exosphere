package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Originate/exosphere/exo-add-go/helpers"
	"github.com/Originate/exosphere/exo-add-go/types"
	"github.com/Originate/exosphere/exo-add-go/user_input"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var serviceRole, serviceType, description, author, templatePath, modelName, protectionLevel, appDir string
var serviceConfig types.ServiceConfig

var rootCmd = &cobra.Command{
	Use: "exo add",
	Short: "\nAdds a new service to the current application.\nThis command must be called in the root directory of the application.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 && args[0] == "help" {
			if err := cmd.Help(); err != nil {
				panic(err)
			}
			return
		}
		fmt.Print("We are about to add a new Exosphere service to the application!\n")
		reader := bufio.NewReader(os.Stdin)
		serviceRole := strings.TrimSpace(userInput.Ask(reader, "ServiceName: "))
		yamlFile, err := ioutil.ReadFile("application.yml")
		var appConfig types.AppConfig
		err = yaml.Unmarshal(yamlFile, &appConfig)
		helpers.CheckForService(serviceRole, helpers.GetExistingServices(appConfig.Services))
		chosenTemplate = userInput.Choose(reader, "Select a template: ", helpers.GetTemplateDirs())
		templatePath = path.Join(".exosphere", chosenTemplate)
		template, err := template.Get(templatePath)
		if err != nil {
			log.Fatalf("Failed to fetch %s template: %s", chosenTemplate, err)
		}
		if err = template.Execute("."); err != nil {
			log.Fatalf(`Failed to create the service "%s": %s`, chosenTemplate, err)
		}
		helpers.CreateServiceYML(serviceRole)
		// filin application.yml
		fmt.Println("\ndone")
	},
}


func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

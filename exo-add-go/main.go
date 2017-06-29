package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/Originate/exosphere/exo-add-go/helpers"
	"github.com/Originate/exosphere/exo-add-go/types"
	"github.com/Originate/exosphere/exo-add-go/user_input"
	"github.com/spf13/cobra"
	"github.com/tmrts/boilr/pkg/template"
	"gopkg.in/yaml.v2"
)

var rootCmd = &cobra.Command{
	Use:   "exo add",
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
		serviceRole := strings.TrimSpace(userInput.Ask(reader, "ServiceRole: "))

		yamlFile, err := ioutil.ReadFile("application.yml")
		var appConfig types.AppConfig
		err = yaml.Unmarshal(yamlFile, &appConfig)
		helpers.CheckForService(serviceRole, helpers.GetExistingServices(appConfig.Services))

		chosenTemplate := userInput.Choose(reader, "Please select a template:\n", helpers.GetTemplates())
		templateDir := path.Join(".exosphere", chosenTemplate)
		template, err := template.Get(templateDir)
		if err != nil {
			log.Fatalf("Failed to fetch %s template: %s", chosenTemplate, err)
		}
		if err = template.Execute(serviceRole); err != nil {
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

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/Originate/exosphere/exo-add-go/helpers"
	"github.com/Originate/exosphere/exo-add-go/os_helpers"
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
		chosenTemplate := userInput.Choose(reader, "Please select a template:\n", helpers.GetTemplates())
		templateDir := path.Join(".exosphere", chosenTemplate)
		template, err := template.Get(templateDir)
		if err != nil {
			log.Fatalf("Failed to fetch %s template: %s", chosenTemplate, err)
		}
		serviceTmpDir := path.Join("tmp", "service-tmp")
		if err = helpers.CreateServiceTmpDir(); err != nil {
			log.Fatalf(`Failed to create a temp folder for the service "%s": %s`, chosenTemplate, err)
		}
		if err = template.Execute(serviceTmpDir); err != nil {
			log.Fatalf(`Failed to create the service "%s": %s`, chosenTemplate, err)
		}

		serviceRole := osHelpers.GetSubdirectories(serviceTmpDir)[0]
		fmt.Println(serviceRole)
		yamlFile, err := ioutil.ReadFile("application.yml")
		var appConfig types.AppConfig
		err = yaml.Unmarshal(yamlFile, &appConfig)
		helpers.CheckForService(serviceRole, helpers.GetExistingServices(appConfig.Services))
		osHelpers.MoveDir(path.Join(serviceTmpDir, serviceRole), serviceRole)

		helpers.CreateServiceYML(serviceRole)
		// not removing serviceTmpDir yet

		// make service dir from the template in a folder named tmp/ or [random-string]/
		// tmp/ contains [service-dir]/, check that name against existing service
		// if service not already exists, make a service.yml in tmp/[service-dir] (because tmp/ is used for service.yml template, it's safer to choose random string as above)
		// mv everything in tmp/[service-dir] to its parent dir
		// update application.yml with the service name

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

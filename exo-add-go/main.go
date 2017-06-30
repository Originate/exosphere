package main

import (
	"bufio"
	"fmt"
	"os"
	"path"

	"github.com/Originate/exosphere/exo-add-go/helpers"
	"github.com/Originate/exosphere/exo-add-go/os_helpers"
	"github.com/Originate/exosphere/exo-add-go/service_yml_template"
	"github.com/Originate/exosphere/exo-add-go/user_input"
	"github.com/spf13/cobra"
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
		chosenTemplate := userInput.Choose(reader, "Please select a template:", helpers.GetTemplates())
		serviceTmpDir := helpers.CreateTmpServiceDir(chosenTemplate)

		serviceRole := osHelpers.GetSubdirectories(serviceTmpDir)[0]
		appConfig := helpers.GetAppConfig()
		helpers.CheckForService(serviceRole, helpers.GetExistingServices(appConfig.Services))

		osHelpers.MoveDir(path.Join(serviceTmpDir, serviceRole), serviceRole)
		serviceYmlTemplate.CreateServiceYML(serviceRole)
		os.RemoveAll("tmp")
		helpers.UpdateAppConfig(serviceRole, appConfig)
		fmt.Println("\ndone")
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

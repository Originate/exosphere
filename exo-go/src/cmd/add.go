package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path"

	"github.com/Originate/exosphere/exo-go/src/app_config_helpers"
	"github.com/Originate/exosphere/exo-go/src/os_helpers"
	"github.com/Originate/exosphere/exo-go/src/service_helpers"
	"github.com/Originate/exosphere/exo-go/src/template_helpers"
	"github.com/Originate/exosphere/exo-go/src/user_input_helpers"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a new service to the current application",
	Long:  "Adds a new service to the current application",
	Run: func(cmd *cobra.Command, args []string) {
		if printHelpIfNecessary(cmd, args) {
			return
		}
		fmt.Print("We are about to add a new Exosphere service to the application!\n")
		reader := bufio.NewReader(os.Stdin)
		if !templateHelpers.HasTemplateDirectory() {
			fmt.Println("no templates found\n\nPlease add templates to the \".exosphere\" folder of your code base.")
			os.Exit(1)
		}
		templatesChoices, err := templateHelpers.GetTemplates()
		if err != nil {
			panic(err)
		}
		chosenTemplate, err := userInputHelpers.Choose(reader, "Please choose a template:", templatesChoices)
		if err != nil {
			panic(err)
		}
		serviceTmpDir, err := templateHelpers.CreateTmpServiceDir(chosenTemplate)
		if err != nil {
			panic(err)
		}
		subdirectories, err := osHelpers.GetSubdirectories(serviceTmpDir)
		if err != nil {
			panic(err)
		}
		serviceRole := subdirectories[0]
		appConfig, err := appConfigHelpers.GetAppConfig()
		if err != nil {
			panic(err)
		}
		err = serviceHelpers.VerifyServiceDoesNotExist(serviceRole, serviceHelpers.GetExistingServices(appConfig.Services))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = osHelpers.MoveDir(path.Join(serviceTmpDir, serviceRole), serviceRole)
		if err != nil {
			panic(err)
		}
		if !osHelpers.FileExists(path.Join(serviceRole, "service.yml")) {
			err = templateHelpers.CreateServiceYML(serviceRole)
			if err != nil {
				panic(err)
			}
		}
		err = os.RemoveAll(serviceTmpDir)
		if err != nil {
			panic(err)
		}
		err = appConfigHelpers.UpdateAppConfig(reader, serviceRole, appConfig)
		if err != nil {
			panic(err)
		}
		fmt.Println("\ndone")
	},
}

func init() {
	RootCmd.AddCommand(addCmd)
}

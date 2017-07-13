package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/Originate/exosphere/exo-go/src/app_config_helpers"
	"github.com/Originate/exosphere/exo-go/src/os_helpers"
	"github.com/Originate/exosphere/exo-go/src/template_helpers"
	"github.com/Originate/exosphere/exo-go/src/user_input_helpers"
	"github.com/Originate/exosphere/exo-go/src/util"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a new service to the current application",
	Long:  "Adds a new service to the current application",
	Run: func(cmd *cobra.Command, args []string) {
		if util.PrintHelpIfNecessary(cmd, args) {
			return
		}
		fmt.Print("We are about to add a new Exosphere service to the application!\n")

		reader := bufio.NewReader(os.Stdin)
		chosenTemplate := userInputHelpers.Choose(reader, "Please choose a template:", templateHelpers.GetTemplates())
		serviceTmpDir := templateHelpers.CreateTmpServiceDir(chosenTemplate)

		serviceRole := osHelpers.GetSubdirectories(serviceTmpDir)[0]
		appConfig := appConfigHelpers.GetAppConfig()
		if err := appConfigHelpers.VerifyServiceDoesNotExist(serviceRole, appConfigHelpers.GetServiceNames(appConfig.Services)); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		osHelpers.MoveDir(path.Join(serviceTmpDir, serviceRole), serviceRole)
		if !osHelpers.FileExists(path.Join(serviceRole, "service.yml")) {
			templateHelpers.CreateServiceYML(serviceRole)
		}
		if err := os.RemoveAll(serviceTmpDir); err != nil {
			log.Fatal("Failed to remove service tmp folder")
		}
		appConfigHelpers.UpdateAppConfig(reader, serviceRole, appConfig)
		fmt.Println("\ndone")
	},
}

func init() {
	RootCmd.AddCommand(addCmd)
}

package cmd

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/Originate/exosphere/exo-go/src/app_config_helpers"
	"github.com/Originate/exosphere/exo-go/src/os_helpers"
	"github.com/Originate/exosphere/exo-go/src/service_helpers"
	"github.com/Originate/exosphere/exo-go/src/template_helpers"
	"github.com/Originate/exosphere/exo-go/src/util"
	prompt "github.com/segmentio/go-prompt"
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
		templates := templateHelpers.GetTemplates()
		chosenTemplate := templates[prompt.Choose("Please choose a template:", templates)]
		serviceTmpDir := templateHelpers.CreateTmpServiceDir(chosenTemplate)

		serviceRole := osHelpers.GetSubdirectories(serviceTmpDir)[0]
		appConfig := appConfigHelpers.GetAppConfig()
		if err := serviceHelpers.VerifyServiceDoesNotExist(serviceRole, serviceHelpers.GetExistingServices(appConfig.Services)); err != nil {
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
		appConfigHelpers.UpdateAppConfig(serviceRole, appConfig)
		fmt.Println("\ndone")
	},
}

func init() {
	RootCmd.AddCommand(addCmd)
}

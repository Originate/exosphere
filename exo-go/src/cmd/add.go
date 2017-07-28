package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/Originate/exosphere/exo-go/src/app_config_helpers"
	"github.com/Originate/exosphere/exo-go/src/exotemplate"
	"github.com/Originate/exosphere/exo-go/src/os_helpers"
	"github.com/Originate/exosphere/exo-go/src/types"
	prompt "github.com/segmentio/go-prompt"
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
		appDir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		if !exotemplate.HasTemplatesDir(appDir) {
			fmt.Println("no templates found\n\nPlease add templates to the \".exosphere\" folder of your code base.")
			os.Exit(1)
		}
		templatesChoices, err := exotemplate.GetTemplates(appDir)
		if err != nil {
			panic(err)
		}
		chosenTemplate := templatesChoices[prompt.Choose("Please choose a template:", templatesChoices)]
		if err != nil {
			panic(err)
		}
		serviceTmpDir, err := exotemplate.CreateTmpServiceDir(appDir, chosenTemplate)
		if err != nil {
			panic(err)
		}
		subdirectories, err := osHelpers.GetSubdirectories(serviceTmpDir)
		if err != nil {
			panic(err)
		}
		serviceRole := subdirectories[0]
		appConfig, err := types.NewAppConfig(appDir)
		if err != nil {
			panic(err)
		}
		err = appConfig.VerifyServiceDoesNotExist(serviceRole)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = osHelpers.MoveDir(path.Join(serviceTmpDir, serviceRole), path.Join(appDir, serviceRole))
		if err != nil {
			panic(err)
		}
		if !osHelpers.FileExists(path.Join(appDir, serviceRole, "service.yml")) {
			var templateDir string
			templateDir, err = exotemplate.CreateServiceTemplateDir(serviceRole)
			if err != nil {
				panic(err)
			}
			err = exotemplate.Run(templateDir, path.Join(appDir, serviceRole))
			if err != nil {
				panic(err)
			}
		}
		err = os.RemoveAll(serviceTmpDir)
		if err != nil {
			panic(err)
		}
		err = appConfigHelpers.UpdateAppConfig(appDir, serviceRole, appConfig)
		if err != nil {
			panic(err)
		}
		fmt.Println("\ndone")
	},
}

func init() {
	RootCmd.AddCommand(addCmd)
}

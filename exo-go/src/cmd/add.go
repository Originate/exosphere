package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/Originate/exosphere/exo-go/src/config"
	"github.com/Originate/exosphere/exo-go/src/template"
	"github.com/Originate/exosphere/exo-go/src/util"
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
		if !template.HasTemplatesDir(appDir) {
			fmt.Println("no templates found\n\nPlease add templates to the \".exosphere\" folder of your code base.")
			os.Exit(1)
		}
		templatesChoices, err := template.GetTemplates(appDir)
		if err != nil {
			panic(err)
		}
		chosenTemplate := templatesChoices[prompt.Choose("Please choose a template:", templatesChoices)]
		if err != nil {
			panic(err)
		}
		serviceTmpDir, err := template.CreateTmpServiceDir(appDir, chosenTemplate)
		if err != nil {
			panic(err)
		}
		subdirectories, err := util.GetSubdirectories(serviceTmpDir)
		if err != nil {
			panic(err)
		}
		serviceRole := subdirectories[0]
		appConfig, err := config.NewAppConfig(appDir)
		if err != nil {
			panic(err)
		}
		err = appConfig.VerifyServiceDoesNotExist(serviceRole)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = util.MoveDirectory(path.Join(serviceTmpDir, serviceRole), path.Join(appDir, serviceRole))
		if err != nil {
			panic(err)
		}
		if !util.DoesFileExist(path.Join(appDir, serviceRole, "service.yml")) {
			var templateDir string
			templateDir, err = template.CreateServiceTemplateDir(serviceRole)
			if err != nil {
				panic(err)
			}
			err = template.Run(templateDir, path.Join(appDir, serviceRole))
			if err != nil {
				panic(err)
			}
		}
		err = os.RemoveAll(serviceTmpDir)
		if err != nil {
			panic(err)
		}
		err = config.UpdateAppConfig(appDir, serviceRole, appConfig)
		if err != nil {
			panic(err)
		}
		fmt.Println("\ndone")
	},
}

func init() {
	RootCmd.AddCommand(addCmd)
}

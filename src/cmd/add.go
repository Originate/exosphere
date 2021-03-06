package cmd

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/Originate/exosphere/src/template"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/util"
	prompt "github.com/kofalt/go-prompt"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a new service to the current application",
	Long:  "Adds a new service to the current application",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("We are about to add a new Exosphere service to the application!\n")
		userContext, err := GetUserContext()
		if err != nil {
			panic(err)
		}
		appContext := userContext.AppContext
		appDir := appContext.Location
		hasTemplateDir, err := template.HasTemplatesDir(appDir)
		if err != nil {
			panic(err)
		}
		if !hasTemplateDir {
			exitBecauseNoTemplates()
		}
		templatesChoices, err := template.GetTemplates(appDir)
		if err != nil {
			panic(err)
		}
		if len(templatesChoices) == 0 {
			exitBecauseNoTemplates()
		}
		chosenTemplate := templatesChoices[prompt.Choose("Please choose a template", templatesChoices)]
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
		appConfig, err := types.NewAppConfig(appDir)
		if err != nil {
			panic(err)
		}
		err = appConfig.VerifyServiceRoleDoesNotExist(serviceRole)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = util.MoveDirectory(path.Join(serviceTmpDir, serviceRole), path.Join(appDir, serviceRole))
		if err != nil {
			panic(err)
		}
		isService, err := util.DoesFileExist(path.Join(appDir, serviceRole, "service.yml"))
		if err != nil {
			panic(err)
		}
		if !isService {
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
		err = appContext.AddService(serviceRole)
		if err != nil {
			panic(err)
		}
		fmt.Println("\ndone")
	},
}

func exitBecauseNoTemplates() {
	log.Fatalln("no templates found\n\nPlease add templates to the \".exosphere/service_templates\" folder of your code base.")
}

func init() {
	RootCmd.AddCommand(addCmd)
}

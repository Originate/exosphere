package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	yaml "gopkg.in/yaml.v2"

	"github.com/Originate/exosphere/exo-go/src/app_config_helpers"
	"github.com/Originate/exosphere/exo-go/src/os_helpers"
	"github.com/Originate/exosphere/exo-go/src/template_helpers"
	"github.com/Originate/exosphere/exo-go/src/util"
	"github.com/spf13/cobra"
)

var removeTemplateCmd = &cobra.Command{
	Use:   "remove-template",
	Short: "Removes an existing service template from .exosphere",
	Long: `Removes an existing service template from .exosphere
This command must be called in the root directory of the application`,
	Run: func(cmd *cobra.Command, args []string) {
		if util.PrintHelpIfNecessary(cmd, args) {
			return
		}
		if len(args) != 1 {
			fmt.Println("not enough arguments")
      os.Exit(1)
		}
		templateName := args[0]
		templateDir := path.Join(".exosphere", templateName)
		fmt.Printf("We are about to remove the template \"%s\"\n\n", templateName)
		if !osHelpers.DirectoryExists(templateDir) {
			fmt.Println("Error: template does not exist")
			os.Exit(1)
		} else {
			if err := templateHelpers.RemoveTemplate(templateName, templateDir); err != nil {
				log.Fatalf(`Failed to remove template "%s": %s`, templateName, err)
			}
			appConfig := appConfigHelpers.GetAppConfig()
			delete(appConfig.Templates, templateName)
			bytes, err := yaml.Marshal(appConfig)
			if err != nil {
				log.Fatalf("Failed to marshal application.yml: %s", err)
			}
			err = ioutil.WriteFile("application.yml", bytes, 0777)
			if err != nil {
				log.Fatalf("Failed to update application.yml: %s", err)
			}
			fmt.Println("\ndone")
		}
	},
}

func init() {
	RootCmd.AddCommand(removeTemplateCmd)
}

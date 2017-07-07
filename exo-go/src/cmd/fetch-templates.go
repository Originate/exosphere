package cmd

import (
	"fmt"
	"log"
	"path"

	"github.com/Originate/exosphere/exo-go/src/app_config_helpers"
	"github.com/Originate/exosphere/exo-go/src/os_helpers"
	"github.com/Originate/exosphere/exo-go/src/template_helpers"
	"github.com/Originate/exosphere/exo-go/src/util"
	"github.com/spf13/cobra"
)

var fetchTemplateCmd = &cobra.Command{
	Use:   "fetch-templates",
	Short: "Fetch service templates",
	Long: `Fetch service templates
This command must be called in the root directory of the application`,
	Run: func(cmd *cobra.Command, args []string) {
		if util.PrintHelpIfNecessary(cmd, args) {
			return
		}
		fmt.Print("We are about to fetch remote templates\n\n")
		appConfig := appConfigHelpers.GetAppConfig()
		if len(appConfig.Templates) > 0 {
			for templateName, gitURL := range appConfig.Templates {
				templateDir := path.Join(".exosphere", templateName)
				if osHelpers.DirectoryExists(templateDir) {
					if err := templateHelpers.UpdateTemplate(gitURL, templateName); err != nil {
						log.Fatalf(`Failed to update template "%s": %s`, templateName, err)
					}
				} else {
					if err := templateHelpers.FetchTemplate(gitURL, templateName, templateDir); err != nil {
						log.Fatalf(`Failed to fetch template "%s": %s`, templateName, err)
					}
				}
			}
		}
		fmt.Println("\ndone")
	},
}

func init() {
	RootCmd.AddCommand(fetchTemplateCmd)
}

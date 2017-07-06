package cmd

import (
	"fmt"
	"log"

	"github.com/Originate/exosphere/exo-go/src/template_helpers"
	"github.com/Originate/exosphere/exo-go/src/util"
	"github.com/spf13/cobra"
	"github.com/tmrts/boilr/pkg/template"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new Exosphere application",
	Run: func(cmd *cobra.Command, args []string) {
		if util.PrintHelpIfNecessary(cmd, args) {
			return
		}
		fmt.Print("We are about to create a new Exosphere application\n\n")
		templatePath, err := templateHelpers.CreateTemplateDir()
		if err != nil {
			log.Fatalf("Failed to create the template: %s", err)
		}
		applicationYmlTemplate, err := template.Get(templatePath)
		if err != nil {
			log.Fatalf("Failed to fetch the application template: %s", err)
		}
		if err = applicationYmlTemplate.Execute("."); err != nil {
			log.Fatalf("Failed to create the application: %s", err)
		}
		if err = templateHelpers.RemoveTemplateDir(); err != nil {
			log.Fatalf("Failed to remove the template: %s", err)
		}
		fmt.Println("\ndone")
	},
}

func init() {
	RootCmd.AddCommand(createCmd)
}

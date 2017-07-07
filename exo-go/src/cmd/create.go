package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/Originate/exosphere/exo-go/src/template_helpers"
	"github.com/Originate/exosphere/exo-go/src/util"
	"github.com/spf13/cobra"
	"github.com/tmrts/boilr/pkg/template"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new Exosphere application",
	Run: func(cmd *cobra.Command, args []string) {
		if printHelpIfNecessary(cmd, args) {
			return
		}
		fmt.Print("We are about to create a new Exosphere application\n\n")
		templateDir, err := templateHelpers.CreateApplicationTemplateDir()
		if err != nil {
			log.Fatalf("Failed to create the template: %s", err)
		}
		applicationTemplate, err := template.Get(templateDir)
		if err != nil {
			log.Fatalf("Failed to fetch the application template: %s", err)
		}
		if err = applicationTemplate.Execute("."); err != nil {
			log.Fatalf("Failed to create the application: %s", err)
		}
		if err = os.RemoveAll(templateDir); err != nil {
			log.Fatalf("Failed to remove the application template: %s", err)
		}
		fmt.Println("\ndone")
	},
}

func init() {
	RootCmd.AddCommand(createCmd)
}

package cmd

import (
	"fmt"
	"os"

	"github.com/Originate/exosphere/exo-go/src/template_helpers"
	"github.com/spf13/cobra"
	"github.com/tmrts/boilr/pkg/template"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a new Exosphere application",
	Run: func(cmd *cobra.Command, args []string) {
		if printHelpIfNecessary(cmd, args) {
			return
		}
		fmt.Print("We are about to create a new Exosphere application\n\n")
		templateDir, err := templateHelpers.CreateApplicationTemplateDir()
		if err != nil {
			panic(err)
		}
		applicationTemplate, err := template.Get(templateDir)
		if err != nil {
			panic(err)
		}
		if err = applicationTemplate.Execute("."); err != nil {
			panic(err)
		}
		if err = os.RemoveAll(templateDir); err != nil {
			panic(err)
		}
		fmt.Println("\ndone")
	},
}

func init() {
	RootCmd.AddCommand(createCmd)
}

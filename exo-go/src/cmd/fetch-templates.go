package cmd

import (
	"fmt"
	"log"

	"github.com/Originate/exosphere/exo-go/src/template_helpers"
	"github.com/Originate/exosphere/exo-go/src/util"
	"github.com/spf13/cobra"
)

var fetchTemplatesCmd = &cobra.Command{
	Use:   "fetch-templates",
	Short: "Fetches updates for all existing templates",
	Long:  "Fetches updates for all existing git submodules in the .exosphere folder",
	Run: func(cmd *cobra.Command, args []string) {
		if util.PrintHelpIfNecessary(cmd, args) {
			return
		}
		fmt.Print("We are about to fetch updates for the remote templates\n\n")
		if err := templateHelpers.FetchTemplates(); err != nil {
			log.Fatalf(`Failed to fetch templates: %s`, err)
		}
		fmt.Println("\ndone")
	},
}

func init() {
	RootCmd.AddCommand(fetchTemplatesCmd)
}

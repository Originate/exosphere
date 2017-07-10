package cmd

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/Originate/exosphere/exo-go/src/os_helpers"
	"github.com/Originate/exosphere/exo-go/src/template_helpers"
	"github.com/Originate/exosphere/exo-go/src/util"
	"github.com/spf13/cobra"
)

var addTemplateCmd = &cobra.Command{
	Use:   "add-template",
	Short: "Adds a remote service template to .exosphere",
	Long:  "Adds a remote service template to .exosphere",
	Run: func(cmd *cobra.Command, args []string) {
		if util.PrintHelpIfNecessary(cmd, args) {
			return
		}
		if len(args) != 2 {
			fmt.Println("not enough arguments")
			os.Exit(1)
		}
		fmt.Print("We are about to add a new service template\n")
		templateName, gitURL := args[0], args[1]
		templateDir := path.Join(".exosphere", templateName)
		if osHelpers.DirectoryExists(templateDir) {
			fmt.Printf(`The template "%s" already exists\n`, templateName)
			os.Exit(1)
		} else {
			if err := templateHelpers.AddTemplate(gitURL, templateName, templateDir); err != nil {
				log.Fatalf(`Failed to add template "%s": %s`, templateName, err)
			}
		}
		fmt.Println("\ndone")
	},
}

func init() {
	RootCmd.AddCommand(addTemplateCmd)
}

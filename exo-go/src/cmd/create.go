package cmd

import (
	"fmt"

	"github.com/Originate/exosphere/exo-go/src/exotemplate"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a new Exosphere application",
	Run: func(cmd *cobra.Command, args []string) {
		if printHelpIfNecessary(cmd, args) {
			return
		}
		fmt.Print("We are about to create a new Exosphere application\n\n")
		templateDir, err := exotemplate.CreateApplicationTemplateDir()
		if err != nil {
			panic(err)
		}
		if err := exotemplate.Run(templateDir, "."); err != nil {
			panic(err)
		}
		fmt.Println("\ndone")
	},
}

func init() {
	RootCmd.AddCommand(createCmd)
}

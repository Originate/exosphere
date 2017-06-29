package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Originate/exosphere/exo-create/helpers"
	"github.com/spf13/cobra"
	"github.com/tmrts/boilr/pkg/template"
)

var rootCmd = &cobra.Command{
	Use:   "exo create",
	Short: "Create a new Exosphere application",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 && args[0] == "help" {
			if err := cmd.Help(); err != nil {
				panic(err)
			}
			return
		}
		fmt.Print("We are about to create a new Exosphere application\n\n")
		templatePath, err := helpers.CreateTemplateDir()
		if err != nil {
			log.Fatalf("Failed to create the template: %s", err)
		}
		template, err := template.Get(templatePath)
		if err != nil {
			log.Fatalf("Failed to fetch the application template: %s", err)
		}
		if err = template.Execute("."); err != nil {
			log.Fatalf("Failed to create the application: %s", err)
		}
		if err = helpers.RemoveTemplateDir(); err != nil {
			log.Fatalf("Failed to remove the template: %s", err)
		}
		fmt.Println("\ndone")
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

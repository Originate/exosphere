package main

import (
	"context"
	"fmt"
	"os"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "exo create",
	Short: "Create a new Exosphere application\n- options: \"[<app-name>] [<app-version>] [<exocom-version>] [<app-description>]\"",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 && args[0] == "help" {
			err := cmd.Help()
			if err != nil {
				panic(err)
			}
			return
		}
		fmt.Print("We are about to create a new Exosphere application\n\n")
		c, err := client.NewEnvClient()
		if err != nil {
			panic(err)
		}
		err = removeImages(c)
		if err != nil {
			panic(err)
		}
		fmt.Println("removed all dangling images")
		err = removeVolumes(c)
		if err != nil {
			panic(err)
		}
		fmt.Println("removed all dangling volumes")
	},
}

// nolint unused
func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

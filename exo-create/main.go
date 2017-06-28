package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/Originate/exosphere/exo-create/helpers"
	"github.com/Originate/exosphere/exo-create/types"
	"github.com/Originate/exosphere/exo-create/user_input"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "exo create",
	Short: "Usage: \"exo create application\"\nCreate a new Exosphere application\n  - options: \"[<app-name>] [<app-version>] [<exocom-version>] [<app-description>]\"",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Error: missing argument for 'create' command")
			return
		} else if args[0] == "help" {
			if err := cmd.Help(); err != nil {
				panic(err)
			}
			return
		} else if args[0] != "application" {
			fmt.Println("Error: cannot create '" + args[0] + "'")
			return
		}
		fmt.Print("We are about to create a new Exosphere application\n\n")
		appConfig := getAppConfig(args)
		if err := helpers.CreateApplicationYAML(appConfig); err != nil {
			panic(err)
		}
		if err := os.Mkdir(path.Join(appConfig.AppName, ".exosphere"), os.FileMode(0777)); err != nil {
			panic(err)
		}
		fmt.Println("\ndone")
	},
}

func getAppConfig(args []string) types.AppConfig {
	var appName, appVersion, exocomVersion, appDescription string
	reader := bufio.NewReader(os.Stdin)
	if len(args) > 1 {
		appName = args[1]
	} else {
		appName = userInput.Ask(reader, "Name of the application to create: ", "", true)
	}
	if len(args) > 2 {
		appVersion = args[2]
	} else {
		appVersion = userInput.Ask(reader, "Initial version: ", "0.0.1", true)
	}
	if len(args) > 3 {
		exocomVersion = args[3]
	} else {
		exocomVersion = userInput.Ask(reader, "ExoCom version: ", "0.22.1", true)
	}
	if len(args) > 4 {
		appDescription = strings.Join(args[4:], " ")
	} else {
		appDescription = userInput.Ask(reader, "Description: ", "", false)
	}
	return types.AppConfig{AppName: appName, AppVersion: appVersion, ExocomVersion: exocomVersion, AppDescription: appDescription}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

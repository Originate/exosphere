package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/Originate/exosphere/exo-create/helpers"
	"github.com/Originate/exosphere/exo-create/types"
	"github.com/Originate/exosphere/exo-create/user_input"
	"github.com/spf13/cobra"
)

var appName, appVersion, exocomVersion, appDescription string

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
		appConfig := getAppConfig(args)
		if err := helpers.CreateApplicationYAML(appConfig); err != nil {
			log.Fatalf("Failed to create application.yml for the application")
		}
		if err := os.Mkdir(path.Join(appConfig.AppName, ".exosphere"), os.FileMode(0777)); err != nil {
			log.Fatalf("Failed to create .exosphere subdirectory for the application")
		}
		fmt.Println("\ndone")
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&appName, "app-name", "", "")
	rootCmd.PersistentFlags().StringVar(&appVersion, "app-version", "", "")
	rootCmd.PersistentFlags().StringVar(&exocomVersion, "exocom-version", "", "")
	rootCmd.PersistentFlags().StringVar(&appDescription, "app-description", "", "")
}

func getAppConfig(args []string) types.AppConfig {
	reader := bufio.NewReader(os.Stdin)
	if helpers.NotSet(appName) {
		appName = userInput.Ask(reader, "Application Name: ", "", true)
	}
	if helpers.NotSet(appVersion) {
		appVersion = userInput.Ask(reader, "Initial version: ", "0.0.1", true)
	}
	if helpers.NotSet(exocomVersion) {
		exocomVersion = userInput.Ask(reader, "ExoCom version: ", "0.22.1", true)
	}
	if helpers.NotSet(appDescription) {
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

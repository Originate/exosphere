package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/Originate/exosphere/exo-go/src/app_config_helpers"
	"github.com/Originate/exosphere/exo-go/src/types"
	"github.com/spf13/cobra"
)

func printHelpIfNecessary(cmd *cobra.Command, args []string) bool {
	if len(args) == 1 && args[0] == "help" {
		if err := cmd.Help(); err != nil {
			panic(err)
		}
		return true
	}
	return false
}

func getAppConfig() types.AppConfig {
	appDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	appConfig, err := appConfigHelpers.GetAppConfig(appDir)
	if err != nil {
		log.Fatalf("Cannot read application configuration: %s", err)
	}
	return appConfig
}

func getSecretsConfig() (string, string) {
	appConfig := getAppConfig()
	secretsBucket := fmt.Sprintf("%s-terraform-secrets", appConfig.Name)
	awsRegion := "us-west-2" //TODO: read from app.yml
	return secretsBucket, awsRegion
}

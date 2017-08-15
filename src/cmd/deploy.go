package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/Originate/exosphere/src/application"
	"github.com/spf13/cobra"
)

var deployInit = &cobra.Command{
	Use:   "init",
	Short: "Prepares an application to be deployed",
	Long:  "Prepares an application to be deployed. Sets up an AWS account, pushes Docker images to ECR and generates Terraform files.",
	Run: func(cmd *cobra.Command, args []string) {
		if printHelpIfNecessary(cmd, args) {
			return
		}
		fmt.Println("We are preparing an application for deployment!")

		deployConfig, err := getDeployConfig()
		if err != nil {
			log.Fatalf("Deploy initialization failed: %s", err)
		}
		err = application.InitDeploy(deployConfig)
		if err != nil {
			log.Fatalf("Deploy initialization failed: %s", err)
		}
	},
}

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploys Exosphere application to the cloud",
	Long:  "Deploys Exosphere application to the cloud. Should be run after 'exo deploy init'.",
	Run: func(cmd *cobra.Command, args []string) {
		if printHelpIfNecessary(cmd, args) {
			return
		}
		fmt.Println("We are about to deploy an application!")

		deployConfig, err := getDeployConfig()
		if err != nil {
			log.Fatalf("Deploy failed: %s", err)
		}

		err = application.StartDeploy(deployConfig)
		if removalErr := application.RemoveSecretsFile(deployConfig.SecretsPath); removalErr != nil {
			fmt.Fprintf(os.Stderr, removalErr.Error())
		}
		if err != nil {
			log.Fatalf("Deploy failed: %s", err)
		}
	},
}

func init() {
	deployCmd.AddCommand(deployInit)
	RootCmd.AddCommand(deployCmd)
}

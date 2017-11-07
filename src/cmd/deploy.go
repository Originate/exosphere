package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Originate/exosphere/src/application/deployer"
	"github.com/Originate/exosphere/src/config"
	"github.com/Originate/exosphere/src/docker/composebuilder"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/util"
	"github.com/spf13/cobra"
)

var deployProfileFlag string
var deployServicesFlag bool

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploys Exosphere application to the cloud",
	Long:  "Deploys Exosphere application to the cloud",
	Run: func(cmd *cobra.Command, args []string) {
		if printHelpIfNecessary(cmd, args) {
			return
		}
		fmt.Println("We are about to deploy an application!")
		context, err := GetContext()
		if err != nil {
			log.Fatal(err)
		}
		homeDir, err := util.GetHomeDirectory()
		if err != nil {
			panic(err)
		}
		serviceConfigs, err := config.GetServiceConfigs(context.AppContext.Location, context.AppContext.Config)
		if err != nil {
			log.Fatalf("Failed to read service configurations: %s", err)
		}
		awsConfig := getAwsConfig(context.AppContext.Config, deployProfileFlag)
		writer := os.Stdout
		terraformDir := filepath.Join(context.AppContext.Location, "terraform")
		deployConfig := types.DeployConfig{
			AppContext:               context.AppContext,
			ServiceConfigs:           serviceConfigs,
			HomeDir:                  homeDir,
			DockerComposeProjectName: composebuilder.GetDockerComposeProjectName(context.AppContext.Config.Name),
			Writer:             writer,
			TerraformDir:       terraformDir,
			SecretsPath:        filepath.Join(terraformDir, "secrets.tfvars"),
			AwsConfig:          awsConfig,
			DeployServicesOnly: deployServicesFlag,

			// git commit hash of the Terraform modules in Originate/exosphere we are using
			TerraformModulesRef: "f456f3f4",
		}

		err = deployer.StartDeploy(deployConfig)
		if err != nil {
			log.Fatalf("Deploy failed: %s", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(deployCmd)
	deployCmd.PersistentFlags().StringVarP(&deployProfileFlag, "profile", "p", "default", "AWS profile to use")
	deployCmd.PersistentFlags().BoolVarP(&deployServicesFlag, "update-services", "", false, "Deploy changes to service images and env vars only")
}

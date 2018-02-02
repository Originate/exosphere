package cmd

import (
	"fmt"
	"log"

	"github.com/Originate/exosphere/src/application"
	"github.com/Originate/exosphere/src/application/deployer"
	"github.com/spf13/cobra"
)

var autoApproveFlag bool

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploys Exosphere application to the cloud",
	Long:  "Deploys Exosphere application to the cloud",
}

var deployInfraCmd = &cobra.Command{
	Use:   "infrastructure [remote-environment-id]",
	Short: "Deploys Exosphere application infrastructure to the cloud",
	Long:  "Deploys Exosphere application infrastructure to the cloud",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("We are about to deploy application infrastructure!")
		deployConfig := getBaseDeployConfig(args[0], deployProfileFlag)
		err := application.GenerateComposeFiles(deployConfig.AppContext)
		if err != nil {
			log.Fatal(err)
		}
		err = deployer.DeployInfrastructure(deployConfig)
		if err != nil {
			log.Fatalf("Deploy failed: %s", err)
		}
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateArgCount(args, 1)
	},
}

var deployServicesCmd = &cobra.Command{
	Use:   "services [remote-environment-id]",
	Short: "Deploys Exosphere services to the cloud",
	Long:  "Deploys Exosphere services to the cloud",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("We are about to deploy an application!")
		deployConfig := getBaseDeployConfig(args[0], deployProfileFlag)
		deployConfig.AutoApprove = autoApproveFlag
		err := application.GenerateComposeFiles(deployConfig.AppContext)
		if err != nil {
			log.Fatal(err)
		}
		err = deployer.DeployServices(deployConfig)
		if err != nil {
			log.Fatalf("Deploy failed: %s", err)
		}
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateArgCount(args, 1)
	},
}

func init() {
	RootCmd.AddCommand(deployCmd)
	deployCmd.AddCommand(deployInfraCmd)
	deployCmd.AddCommand(deployServicesCmd)
	deployInfraCmd.PersistentFlags().StringVarP(&deployProfileFlag, "profile", "p", "default", "AWS profile to use")
	deployServicesCmd.PersistentFlags().StringVarP(&deployProfileFlag, "profile", "p", "default", "AWS profile to use")
	deployServicesCmd.PersistentFlags().BoolVarP(&autoApproveFlag, "auto-approve", "", false, "Deploy changes without prompting for approval")
}

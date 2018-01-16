package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/Originate/exosphere/src/application"
	"github.com/Originate/exosphere/src/application/deployer"
	"github.com/spf13/cobra"
)

var deployProfileFlag string
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
		userContext, err := GetUserContext()
		if err != nil {
			log.Fatal(err)
		}
		remoteEnvironmentID := args[0]
		err = validateRemoteEnvironmentID(userContext, remoteEnvironmentID)
		if err != nil {
			log.Fatal(err)
		}
		err = application.GenerateComposeFiles(userContext.AppContext)
		if err != nil {
			log.Fatal(err)
		}
		deployConfig := getBaseDeployConfig(userContext.AppContext, remoteEnvironmentID)
		writer := os.Stdout
		deployConfig.Writer = writer
		deployConfig.AutoApprove = autoApproveFlag
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
		fmt.Println("We are about to deploy services!")
		userContext, err := GetUserContext()
		if err != nil {
			log.Fatal(err)
		}
		remoteEnvironmentID := args[0]
		err = validateRemoteEnvironmentID(userContext, remoteEnvironmentID)
		if err != nil {
			log.Fatal(err)
		}
		err = application.GenerateComposeFiles(userContext.AppContext)
		if err != nil {
			log.Fatal(err)
		}
		deployConfig := getBaseDeployConfig(userContext.AppContext, remoteEnvironmentID)
		writer := os.Stdout
		deployConfig.Writer = writer
		deployConfig.AutoApprove = autoApproveFlag
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

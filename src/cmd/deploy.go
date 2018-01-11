package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/Originate/exosphere/src/application"
	"github.com/Originate/exosphere/src/application/deployer"
	"github.com/spf13/cobra"
)

var autoApproveFlag bool

var deployCmd = &cobra.Command{
	Use:   "deploy [remote-environment-id]",
	Short: "Deploys Exosphere application to the cloud",
	Long:  "Deploys Exosphere application to the cloud",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("We are about to deploy an application!")
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
		err = deployer.StartDeploy(deployConfig)
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
	deployCmd.PersistentFlags().StringVarP(&awsProfileFlag, "profile", "p", "default", "AWS profile to use")
	deployCmd.PersistentFlags().BoolVarP(&autoApproveFlag, "auto-approve", "", false, "Deploy changes without prompting for approval")
}

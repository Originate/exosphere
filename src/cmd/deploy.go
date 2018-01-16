package cmd

import (
	"fmt"
	"log"

	"github.com/Originate/exosphere/src/application"
	"github.com/Originate/exosphere/src/application/deployer"
	"github.com/spf13/cobra"
)

var deployProfileFlag string
var autoApproveFlag bool

var deployCmd = &cobra.Command{
	Use:   "deploy [remote-environment-id]",
	Short: "Deploys Exosphere application to the cloud",
	Long:  "Deploys Exosphere application to the cloud",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("We are about to deploy an application!")
		deployConfig := getBaseDeployConfig(args[0], deployProfileFlag)
		deployConfig.AutoApprove = autoApproveFlag
		err := application.GenerateComposeFiles(deployConfig.AppContext)
		if err != nil {
			log.Fatal(err)
		}
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
	deployCmd.PersistentFlags().StringVarP(&deployProfileFlag, "profile", "p", "default", "AWS profile to use")
	deployCmd.PersistentFlags().BoolVarP(&autoApproveFlag, "auto-approve", "", false, "Deploy changes without prompting for approval")
}

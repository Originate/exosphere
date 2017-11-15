package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/Originate/exosphere/src/application/deployer"
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
		deployConfig := getBaseDeployConfig(context.AppContext)
		writer := os.Stdout
		deployConfig.Writer = writer
		deployConfig.DeployServicesOnly = deployServicesFlag
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

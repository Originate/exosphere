package cmd

import (
	"fmt"

	"github.com/Originate/exosphere/exo-go/src/util"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploys Exosphere application to the cloud",
	Run: func(cmd *cobra.Command, args []string) {
		if util.PrintHelpIfNecessary(cmd, args) {
			return
		}
		fmt.Println("exo deploy dis ish!")
	},
}

func init() {
	RootCmd.AddCommand(deployCmd)
}

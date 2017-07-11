package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploys Exosphere application to the cloud",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("exo deploy dis ish!")
	},
}

func init() {
	RootCmd.AddCommand(deployCmd)
}

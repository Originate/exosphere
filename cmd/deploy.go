package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "A brief description of the deploy command",
	Long:  `A longer description of the deploy command.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("fill in functionality here")
	},
}

func init() {
	RootCmd.AddCommand(deployCmd)
}

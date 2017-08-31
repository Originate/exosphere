package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Displays the version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Exosphere v0.23.0.alpha.4")
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}

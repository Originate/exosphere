package cmd

import (
	"fmt"

	"github.com/Originate/exosphere/src"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Displays the version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Exosphere v%s\n", src.Version)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}

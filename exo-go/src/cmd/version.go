package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Exosphere go version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Exosphere-Go v0.1")
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}

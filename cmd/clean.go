package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Description of the clean command",
	Long:  `A longer description of the clean command.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("put some functionality here")
	},
}

func init() {
	RootCmd.AddCommand(cleanCmd)
}

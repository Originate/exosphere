package util

import "github.com/spf13/cobra"

// DoesStringArrayContain returns whether the given string slice
// contains the given string.
func DoesStringArrayContain(strings []string, targetString string) bool {
	for _, element := range strings {
		if element == targetString {
			return true
		}
	}
	return false
}

// PrintHelpIfNecessary prints the cmd help screen if the help is the only argument
func PrintHelpIfNecessary(cmd *cobra.Command, args []string) bool {
	if len(args) == 1 && args[0] == "help" {
		if err := cmd.Help(); err != nil {
			panic(err)
		}
		return true
	}
	return false
}

package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Originate/exocom/go/exocom/src/exocom"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "exocom",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		serviceRoutes := os.Getenv("SERVICE_ROUTES")
		exocom, err := exocom.New(serviceRoutes)
		if err != nil {
			panic(err)
		}
		port := getPort()
		err = exocom.Listen(port)
		if err != nil {
			panic(err)
		}
	},
}

func getPort() int {
	defaultPort := 3100
	userPort := os.Getenv("PORT")
	if userPort == "" {
		return defaultPort
	}
	port, err := strconv.Atoi(userPort)
	if err != nil {
		fmt.Printf("Given port is not numeric, using default port %d", defaultPort)
		return defaultPort
	}
	return port
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

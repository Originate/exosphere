package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Originate/exosphere/exo-run-go/app_runner"
	"github.com/Originate/exosphere/exo-run-go/helpers"
	"github.com/Originate/exosphere/exo-run-go/logger"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "exo run",
	Short: "Create a new Exosphere application",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 && args[0] == "help" {
			if err := cmd.Help(); err != nil {
				panic(err)
			}
			return
		}
		appConfig := helpers.GetAppConfig()
		fmt.Printf("Running %s %s\n", appConfig.Name, appConfig.Version)
		services, silencedServices := helpers.GetExistingServices(appConfig)
		silencedDependencies := helpers.GetSilencedDependencies(appConfig)
		logger := logger.Logger{Roles: services, SilencedRoles: append(silencedServices, silencedDependencies...)}
		fmt.Println(services, silencedServices, silencedDependencies)

		appRunner := appRunner.NewAppRunner(appConfig, logger)
		// TODO: define exocom struct
		appRunner.On("routing setup", func(exocom map[string]map[string]map[string]map[string]string) {
			logger.Log("exocom", "received routing setup", true)
			for command, routing := range exocom["clientRegistry"]["routes"] {
				text := fmt.Sprintf("%s -->", command)
				receivers := []string{}
				for _, receiver := range routing["receivers"] {
					receivers = append(receivers, fmt.Sprintf("%s (%s:%s)", receiver["name"], receiver["host"], receiver["port"]))
				}
				text = fmt.Sprintf("%s%s", strings.Join(receivers, " & "))
			}
		}).On("message", func(messages []map[string]string, receivers []string) {
			message := messages[0]
			if message["name"] != message["originalName"] {
				logger.Log("exocom", fmt.Sprintf("%s --[ %s ] - [ %s ]-> %s %s", message["sender"], message["originalName"], message["name"], strings.Join(receivers, " and ")), true)
			} else {
				logger.Log("exocom", fmt.Sprintf("%s --[ %s ]-> %s %s", message["sender"], message["name"], strings.Join(receivers, " and ")), true)
			}
			indent := strings.Repeat(" ", len(message["sender"])*2)
			if len(message["payload"]) > 0 {
				for _, line := range message["payload"] {
					logger.Log("exocom", fmt.Sprintf("%s%s", indent, line), false)
				}
			} else {
				logger.Log("exocom", fmt.Sprintf("%s%s", indent, "no payload"), false)
			}
		})
		appRunner.Start()
		fmt.Println("\ndone")
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

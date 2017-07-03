package main

import (
	"fmt"
	"os"

	"github.com/Originate/exosphere/exo-run-go/helpers"
	"github.com/Originate/exosphere/exo-run-go/logger"
	"github.com/spf13/cobra"
	"github.com/fatih/color"
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

		appRunner := appRunner.appRunner{}
		emitter := emission.NewEmitter()
		emitter.On("routing setup", func() {
			logger.Log("exocom", "received routing setup", true)
			for command routing := range appRunner.exocom.clientRegistry.routes
				text := fmt.Sprintf("%s -->", color.Bold(command))
				receivers := []string{}
				for _, receiver := range routing.receivers {
					receivers = append(receivers, fmt.Sprintf("%s (%s:%s)", receiver.name, receiver.host, receiver.port))
				}
				text = fmt.Sprintf("%s%s", strings.Join(receivers, " & "))
		}).On("message", func(messages, receivers) {
			message := messages[0]
			if message.name != message.originalName {
				logger.Log("exocom", fmt.Sprintf("%s --[ %s ] - [ %s ]-> %s %s", message.sender, message.originalName, message.name, strings.Join(receivers, " and ")), true)
			} else {
				logger.Log("exocom", fmt.Sprintf("%s --[ %s ]-> %s %s", message.sender, message.name, strings.Join(receivers, " and ")), true)
			}
			indent := strings.Repeat(" ", len(message.sender) * 2)
			if message.payload {
				for _, line := range TODO {
					logger.Log("exocom", indent(color.Dim(line), false))
				}
			} else {
				logger.Log("exocom", indent(color.Dim("no payload")), false
			}
		}).Start()

		fmt.Println("\ndone")
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

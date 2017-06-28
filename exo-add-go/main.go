package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Originate/exosphere/exo-add-go/helpers"
	"github.com/Originate/exosphere/exo-add-go/types"
	"github.com/Originate/exosphere/exo-add-go/user_input"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var serviceRole, serviceType, description, author, templatePath, modelName, protectionLevel string
var serviceConfig types.ServiceConfig

var rootCmd = &cobra.Command{
	Use:   "exo add",
	Short: "\nUsage: #{cyan 'exo add'} #{blue '[<entity-name>]'}\n\nAdds a new service to the current application.\nThis command must be called in the root directory of the application.\n\noptions: #{blue '--service-role=[<service-role>] --service-type=[<service-type>] --template-name=[<template-name>] --model-name=[<model-name>] --protection-level=[<protection-level>] --description=[<description>]'}",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 && args[0] == "help" {
			if err := cmd.Help(); err != nil {
				panic(err)
			}
			return
		}
		serviceConfig = getServiceConfig(args)
		fmt.Print("We are about to add a new Exosphere service to the application!\n")
		yamlFile, err := ioutil.ReadFile("application.yml")
		if err != nil {
			panic(err)
		}
		var appConfig types.AppConfig
		err = yaml.Unmarshal(yamlFile, &appConfig)
		helpers.CheckForService(serviceConfig.ServiceRole, helpers.GetExistingServices(appConfig.Services))

	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&serviceRole, "service-role", "", "")
	rootCmd.PersistentFlags().StringVar(&serviceType, "service-type", "", "")
	rootCmd.PersistentFlags().StringVar(&description, "description", "", "")
	rootCmd.PersistentFlags().StringVar(&author, "author", "", "")
	rootCmd.PersistentFlags().StringVar(&templatePath, "template-name", "", "")
	rootCmd.PersistentFlags().StringVar(&modelName, "model-name", "", "")
	rootCmd.PersistentFlags().StringVar(&protectionLevel, "protection-level", "", "")
}

func getServiceConfig(args []string) types.ServiceConfig {
	reader := bufio.NewReader(os.Stdin)
	if len(serviceRole) == 0 {
		serviceRole = userInput.Ask(reader, "Role of the service to create: ", true)
	}
	if len(serviceType) == 0 {
		serviceType = userInput.Ask(reader, "Type of the service to create: ", true)
	}
	if len(description) == 0 {
		description = userInput.Ask(reader, "Description: ", false)
	}
	if len(author) == 0 {
		author = userInput.Ask(reader, "Author: ", true)
	}
	if len(templatePath) == 0 {
		templates := []string{} // read .exosphere for choices
		templatePath = userInput.Choose(reader, "Template: ", templates)
	}
	if len(modelName) == 0 {
		modelName = userInput.Ask(reader, "Name of the data model: ", false)
	}
	if len(protectionLevel) == 0 {
		protectionLevel = userInput.Choose(reader, "Protection level: ", []string{"public", "private"})
	}
	return types.ServiceConfig{serviceRole, serviceType, description, author, templatePath, modelName, protectionLevel}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

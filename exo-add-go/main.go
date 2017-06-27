package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tcnksm/go-input"
	"gopkg.in/yaml.v2"
)

var rootCmd = &cobra.Command{
	Use:   "exo add",
	Short: "\nUsage: #{cyan 'exo add'} #{blue '[<entity-name>]'}\n\nAdds a new service to the current application.\nThis command must be called in the root directory of the application.\n\noptions: #{blue '--service-role=[<service-role>] --service-type=[<service-type>] --template-name=[<template-name>] --model-name=[<model-name>] --protection-level=[<protection-level>] --description=[<description>]'}",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if len(serviceRole) == 0 {
			serviceRole = ask("Role of the service to create:", true)
		}
		if len(serviceType) == 0 {
			serviceType = ask("Type of the service to create:", true)
		}
		if len(description) == 0 {
			description = ask("Description:", false)
		}
		if len(author) == 0 {
			author = ask("Author:", true)
		}
		if len(templatePath) == 0 {
			templates := []string{} // read .exosphere for choices
			templatePath = choose("Template:", templates)
		}
		if len(modelName) == 0 {
			modelName = ask("Name of the data model:", false)
		}
		if len(protectionLevel) == 0 {
			protectionLevel = choose("Protection level:", []string{"public", "private"})
		}
		serviceConfig := ServiceConfig{serviceRole, serviceType, description, author, templatePath, modelName, protectionLevel}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 && args[0] == "help" {
			err := cmd.Help()
			if err != nil {
				panic(err)
			}
			return
		}
		fmt.Print("We are about to add a new Exosphere service to the application!\n")
		yamlFile, err := ioutil.ReadFile("application.yml")
		if err != nil {
			panic(err)
		}
		var appConfig AppConfig
		err = yaml.Unmarshal(yamlFile, &appConfig)
		checkForService(serviceConfig.ServiceRole, getExistingServices(appConfig.Services))


	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&serviceRole, "service-role")
	rootCmd.PersistentFlags().StringVar(&serviceType, "service-type")
	rootCmd.PersistentFlags().StringVar(&description, "description")
	rootCmd.PersistentFlags().StringVar(&author, "author")
	rootCmd.PersistentFlags().StringVar(&templatePath, "template-name")
	rootCmd.PersistentFlags().StringVar(&modelName, "model-name")
	rootCmd.PersistentFlags().StringVar(&protectionLevel, "protection-level")
}

type Dependency struct {
	Name string
	Version string
	Config DependencyConfig
}

type DependencyConfig struct {
	Ports []string
	Volumes []string
	OnlineText string
	DependencyEnvironment map[string]string
	ServiceEnvironment map[string]string
}

type Service struct {
	Namespace string
	Location string
	DockerImage string
}

type AppConfig struct {
	Name string
	Description string
	Version string
	Dependencies []Dependency
	Services map[string]map[string]Service
}

type ServiceConfig struct {
  ServiceRole, ServiceType, Description, Author, TemplatePath, ModelName, ProtectionLevel string
}

func ask(query string, required bool) string {
  ui := &input.UI{
    Writer: os.Stdout,
    Reader: os.Stdin,
  }
  answer, err := ui.Ask(query, &input.Options{
    Required: required,
    Loop:     true,
  })
  if err != nil {
    panic(err)
  }
  return answer
}

func choose(query string, options []string) string {
	ui := &input.UI{
    Writer: os.Stdout,
    Reader: os.Stdin,
  }
  answer, err := &ui.Select(query, []string options, &input.Options{
    Required: true,
    Loop:     true,
  })
  if err != nil {
    panic(err)
  }
  return answer
}

func checkForService(serviceRole string, existingServices []string) {
	if contains(existingServices, serviceRole) {
		fmt.Println("Service %s already exists in this application", serviceRole)
		os.Exit(1)
	}
}

func getExistingServices(services []string) {
	existingServices := []string{}
	for protectionLevel, serviceConfigs := range services {
		for service := range serviceConfigs {
			append(existingServices, service)
		}
	}
}

func contains(strings []string, targetString string) bool {
	for _, element := range strings {
		if element == targetString {
			return true
		}
	}
	return false
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

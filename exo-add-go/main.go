package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var serviceRole, serviceType, description, author, templatePath, modelName, protectionLevel string
var serviceConfig ServiceConfig

var rootCmd = &cobra.Command{
	Use:   "exo add",
	Short: "\nUsage: #{cyan 'exo add'} #{blue '[<entity-name>]'}\n\nAdds a new service to the current application.\nThis command must be called in the root directory of the application.\n\noptions: #{blue '--service-role=[<service-role>] --service-type=[<service-type>] --template-name=[<template-name>] --model-name=[<model-name>] --protection-level=[<protection-level>] --description=[<description>]'}",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)
		if len(serviceRole) == 0 {
			serviceRole = ask(reader, "Role of the service to create:", true)
		}
		if len(serviceType) == 0 {
			serviceType = ask(reader, "Type of the service to create:", true)
		}
		if len(description) == 0 {
			description = ask(reader, "Description:", false)
		}
		if len(author) == 0 {
			author = ask(reader, "Author:", true)
		}
		if len(templatePath) == 0 {
			templates := []string{} // read .exosphere for choices
			templatePath = choose(reader, "Template:", templates)
		}
		if len(modelName) == 0 {
			modelName = ask(reader, "Name of the data model:", false)
		}
		if len(protectionLevel) == 0 {
			protectionLevel = choose(reader, "Protection level:", []string{"public", "private"})
		}
		serviceConfig = ServiceConfig{serviceRole, serviceType, description, author, templatePath, modelName, protectionLevel}
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
	rootCmd.PersistentFlags().StringVar(&serviceRole, "service-role", "", "")
	rootCmd.PersistentFlags().StringVar(&serviceType, "service-type", "", "")
	rootCmd.PersistentFlags().StringVar(&description, "description", "", "")
	rootCmd.PersistentFlags().StringVar(&author, "author", "", "")
	rootCmd.PersistentFlags().StringVar(&templatePath, "template-name", "", "")
	rootCmd.PersistentFlags().StringVar(&modelName, "model-name", "", "")
	rootCmd.PersistentFlags().StringVar(&protectionLevel, "protection-level", "", "")
}

// Dependency is an unexpected type
type Dependency struct {
	Name    string
	Version string
	Config  DependencyConfig
}

// DependencyConfig is an unexpected type
type DependencyConfig struct {
	Ports                 []string
	Volumes               []string
	OnlineText            string
	DependencyEnvironment map[string]string
	ServiceEnvironment    map[string]string
}

// Service is an unexpected type
type Service struct {
	Namespace   string
	Location    string
	DockerImage string
}

// AppConfig is an unexpected type
type AppConfig struct {
	Name         string
	Description  string
	Version      string
	Dependencies []Dependency
	Services     map[string]map[string]Service
}

// ServiceConfig is an unexpected type
type ServiceConfig struct {
	ServiceRole, ServiceType, Description, Author, TemplatePath, ModelName, ProtectionLevel string
}

func ask(reader *bufio.Reader, query string, required bool) string {
	fmt.Print(query)
	answer, err := reader.ReadString('\n')
	answer = strings.TrimSpace(answer)
	if err != nil {
		panic(err)
	}
	if len(answer) == 0 && required {
		fmt.Println("expect a non-empty string")
		return ask(reader, query, required)
	}
	return answer
}

func choose(reader *bufio.Reader, query string, options []string) string {
	fmt.Print(query)
	if len(options) == 0 {
		panic(fmt.Errorf("no options found\n"))
	}
	for i, option := range options {
		fmt.Printf("%v. %v\n", i+1, option)
	}
	answer, err := reader.ReadString('\n')
	answer = strings.TrimSpace(answer)
	if err != nil {
		panic(err)
	}
	chosenNumber, err := strconv.Atoi(answer)
	if err != nil || !(0 <= chosenNumber-1 && chosenNumber-1 < len(options)) {
		fmt.Printf("expect a number between 1 and %v\n", len(options))
		return choose(reader, query, options)
	}
	return options[chosenNumber-1]
}

func checkForService(serviceRole string, existingServices []string) {
	if contains(existingServices, serviceRole) {
		fmt.Printf("Service %v already exists in this application\n", serviceRole)
		os.Exit(1)
	}
}

func getExistingServices(services map[string]map[string]Service) []string {
	existingServices := []string{}
	for _, serviceConfigs := range services {
		for service := range serviceConfigs {
			existingServices = append(existingServices, service)
		}
	}
	return existingServices
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

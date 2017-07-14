package appRunner

import (
	"fmt"
	"math"
	"path"
	"sync"

	"github.com/Originate/exosphere/exo-go/src/app_config_helpers"
	"github.com/Originate/exosphere/exo-go/src/docker_compose"
	"github.com/Originate/exosphere/exo-go/src/logger"
	"github.com/Originate/exosphere/exo-go/src/service_config_helpers"
	"github.com/Originate/exosphere/exo-go/src/types"
	"github.com/fatih/color"
)

// AppRunner runs the overall application
type AppRunner struct {
	AppConfig            types.AppConfig
	Logger               *logger.Logger
	Env                  map[string]string
	DockerConfigLocation string
	OnlineTexts          map[string]string
}

// NewAppRunner is AppRunner's constructor
func NewAppRunner(appConfig types.AppConfig, logger *logger.Logger, cwd string) *AppRunner {
	appRunner := &AppRunner{AppConfig: appConfig, Logger: logger, Env: appConfigHelpers.GetEnvironmentVariables(appConfig), DockerConfigLocation: path.Join(cwd, "tmp")}
	return appRunner
}

func (appRunner *AppRunner) compileOnlineTexts() (map[string]string, error) {
	onlineTexts := make(map[string]string)
	for _, dependency := range appRunner.AppConfig.Dependencies {
		onlineTexts[dependency.Name] = dependency.GetOnlineText()
	}
	serviceConfigs, err := serviceConfigHelpers.GetServiceConfigs(appRunner.AppConfig)
	if err != nil {
		return map[string]string{}, err
	}
	for serviceName, serviceConfig := range serviceConfigs {
		onlineTexts[serviceName] = serviceConfig.Startup["online-text"]
	}
	return onlineTexts, nil
}

func (appRunner *AppRunner) getEnv() []string {
	formattedEnvVars := []string{}
	for variable, value := range appRunner.Env {
		formattedEnvVars = append(formattedEnvVars, fmt.Sprintf("%s=%s", variable, value))
	}
	return formattedEnvVars
}

// Shutdown shuts down the application and returns the process output and an error if any
func (appRunner *AppRunner) Shutdown(closeMessage, errorMessage string) (string, error) {
	if len(errorMessage) > 0 {
		color.Red(errorMessage)
	} else {
		fmt.Printf("\n\n%s", closeMessage)
	}
	process, err := dockerCompose.KillAllContainers([]string{}, appRunner.DockerConfigLocation, appRunner.Write)
	if err != nil {
		return process.Output, err
	}
	err = process.Wait()
	return process.Output, err
}

// Start runs the application and returns the process output and an error if any
func (appRunner *AppRunner) Start() (string, error) {
	process, err := dockerCompose.RunAllImages(appRunner.getEnv(), appRunner.DockerConfigLocation, appRunner.Write)
	if err != nil {
		return process.Output, err
	}
	onlineTexts, err := appRunner.compileOnlineTexts()
	if err != nil {
		return process.Output, err
	}
	wg := new(sync.WaitGroup)
	for role, onlineText := range onlineTexts {
		wg.Add(1)
		go func(role string, onlineText string) {
			if err = process.WaitForText(onlineText, int(math.Inf(1))); err == nil {
				appRunner.Logger.Log(role, fmt.Sprintf("'%s' is running", role), true)
			}
			wg.Done()
		}(role, onlineText)
	}
	wg.Wait()
	if err == nil {
		appRunner.Write("all services online")
	}
	return process.Output, err
}

// Write logs exo-run output
func (appRunner *AppRunner) Write(text string) {
	appRunner.Logger.Log("exo-run", text, true)
}

package appRunner

import (
	"fmt"
	"path"
	"sync"

	"github.com/Originate/exosphere/exo-go/src/app_config_helpers"
	"github.com/Originate/exosphere/exo-go/src/docker_compose"
	"github.com/Originate/exosphere/exo-go/src/logger"
	"github.com/Originate/exosphere/exo-go/src/process_helpers"
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

// Shutdown shuts down the application
func (appRunner *AppRunner) Shutdown(closeMessage, errorMessage string) error {
	if len(errorMessage) > 0 {
		color.Red(errorMessage)
	} else {
		fmt.Printf("\n\n%s", closeMessage)
	}
	cmd, _, err := dockerCompose.KillAllContainers([]string{}, appRunner.DockerConfigLocation, appRunner.Write)
	if err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}

// Start runs the application
func (appRunner *AppRunner) Start() error {
	_, stdoutBuffer, err := dockerCompose.RunAllImages(appRunner.getEnv(), appRunner.DockerConfigLocation, appRunner.Write)
	if err != nil {
		return err
	}
	onlineTexts, err := appRunner.compileOnlineTexts()
	if err != nil {
		return err
	}
	wg := new(sync.WaitGroup)
	for role, onlineText := range onlineTexts {
		wg.Add(1)
		go func(role string, stdoutBuffer fmt.Stringer, onlineText string) {
			processHelpers.Wait(stdoutBuffer, onlineText, func() {
				appRunner.Logger.Log(role, fmt.Sprintf("'%s' is running", role), true)
			})
			wg.Done()
		}(role, stdoutBuffer, onlineText)
	}
	wg.Wait()
	appRunner.Write("all services online")
	return nil
}

// Write logs exo-run output
func (appRunner *AppRunner) Write(text string) {
	appRunner.Logger.Log("exo-run", text, true)
}

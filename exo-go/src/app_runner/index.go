package appRunner

import (
	"fmt"
	"path"
	"regexp"

	"github.com/Originate/exosphere/exo-go/src/app_config_helpers"
	"github.com/Originate/exosphere/exo-go/src/app_dependency_helpers"
	"github.com/Originate/exosphere/exo-go/src/docker_compose"
	"github.com/Originate/exosphere/exo-go/src/logger"
	"github.com/Originate/exosphere/exo-go/src/process_helpers"
	"github.com/Originate/exosphere/exo-go/src/service_config_helpers"
	"github.com/Originate/exosphere/exo-go/src/types"
	"github.com/fatih/color"
	"github.com/pkg/errors"
)

// AppRunner runs the overall application
type AppRunner struct {
	AppConfig            types.AppConfig
	Logger               *logger.Logger
	AppDir               string
	homeDir              string
	Env                  map[string]string
	DockerConfigLocation string
	OnlineTexts          map[string]string
}

// NewAppRunner is AppRunner's constructor
func NewAppRunner(appConfig types.AppConfig, logger *logger.Logger, appDir, homeDir string) *AppRunner {
	return &AppRunner{
		AppConfig:            appConfig,
		Logger:               logger,
		AppDir:               appDir,
		homeDir:              homeDir,
		Env:                  appConfigHelpers.GetEnvironmentVariables(appConfig, appDir, homeDir),
		DockerConfigLocation: path.Join(appDir, "tmp"),
	}
}

func (a *AppRunner) compileServiceOnlineTexts(serviceConfigs map[string]types.ServiceConfig) map[string]string {
	onlineTexts := make(map[string]string)
	for serviceName, serviceConfig := range serviceConfigs {
		onlineTexts[serviceName] = serviceConfig.Startup["online-text"]
	}
	return onlineTexts
}

func (a *AppRunner) compileDependencyOnlineTexts(serviceConfigs map[string]types.ServiceConfig) map[string]string {
	onlineTexts := make(map[string]string)
	for _, dependency := range a.AppConfig.Dependencies {
		onlineTexts[dependency.Name] = appDependencyHelpers.Build(dependency, a.AppConfig, a.AppDir, a.homeDir).GetOnlineText()
	}
	for _, serviceConfig := range serviceConfigs {
		for _, dependency := range serviceConfig.Dependencies {
			if !dependency.Config.IsEmpty() {
				onlineTexts[dependency.Name] = appDependencyHelpers.Build(dependency, a.AppConfig, a.AppDir, a.homeDir).GetOnlineText()
			}
		}
	}
	return onlineTexts
}

func (a *AppRunner) getEnv() []string {
	formattedEnvVars := []string{}
	for variable, value := range a.Env {
		formattedEnvVars = append(formattedEnvVars, fmt.Sprintf("%s=%s", variable, value))
	}
	return formattedEnvVars
}

// Shutdown shuts down the application and returns the process output and an error if any
func (a *AppRunner) Shutdown(shutdownConfig types.ShutdownConfig) error {
	if len(shutdownConfig.ErrorMessage) > 0 {
		color.Red(shutdownConfig.ErrorMessage)
	} else {
		fmt.Printf("\n\n%s", shutdownConfig.CloseMessage)
	}
	process, err := dockerCompose.KillAllContainers(a.DockerConfigLocation, a.write)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to shutdown the app\nOutput: %s\nError: %s\n", process.Output, err))
	}
	err = process.Wait()
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to shutdown the app\nOutput: %s\nError: %s\n", process.Output, err))
	}
	return nil
}

// Start runs the application and returns the process and returns an error if any
func (a *AppRunner) Start() error {
	dependencyNames, err := appConfigHelpers.GetAllDependencyNames(a.AppDir, a.AppConfig)
	if err != nil {
		return err
	}
	serviceNames := appConfigHelpers.GetServiceNames(a.AppConfig.Services)
	serviceConfigs, err := serviceConfigHelpers.GetServiceConfigs(a.AppDir, a.AppConfig)
	if err != nil {
		return err
	}
	if err := a.startDependencies(dependencyNames, serviceConfigs); err != nil {
		return err
	}
	return a.startServices(serviceNames, serviceConfigs)
}

func (a *AppRunner) startDependencies(dependencyNames []string, serviceConfigs map[string]types.ServiceConfig) error {
	dependencyOnlineTexts := a.compileDependencyOnlineTexts(serviceConfigs)
	process, err := dockerCompose.RunImages(dependencyNames, a.getEnv(), a.DockerConfigLocation, a.write)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to run dependencies\nOutput: %s\nError: %s\n", process.Output, err))
	}
	for dependencyName, onlineText := range dependencyOnlineTexts {
		if err := a.waitForOnlineText(process, dependencyName, onlineText); err != nil {
			return err
		}
	}
	a.write("all dependencies online")
	return nil
}

func (a *AppRunner) startServices(serviceNames []string, serviceConfigs map[string]types.ServiceConfig) error {
	serviceOnlineTexts := a.compileServiceOnlineTexts(serviceConfigs)
	process, err := dockerCompose.RunImages(serviceNames, a.getEnv(), a.DockerConfigLocation, a.write)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to run services\nOutput: %s\nError: %s\n", process.Output, err))
	}
	for role, onlineText := range serviceOnlineTexts {
		if err := a.waitForOnlineText(process, role, onlineText); err != nil {
			return err
		}
	}
	a.write("all services online")
	return nil
}

func (a *AppRunner) waitForOnlineText(process *processHelpers.Process, role, onlineText string) error {
	onlineTextRegex, err := regexp.Compile(fmt.Sprintf("%s.*%s", role, onlineText))
	if err != nil {
		return err
	}
	if err = process.WaitForRegex(onlineTextRegex); err == nil {
		return a.Logger.Log(role, fmt.Sprintf("'%s' is running", role), true)
	}
	return nil
}

// write logs exo-run output
func (a *AppRunner) write(text string) {
	err := a.Logger.Log("exo-run", text, true)
	if err != nil {
		fmt.Printf("Error logging exo-run output: %v\n", err)
	}
}

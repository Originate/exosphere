package appRunner

import (
	"fmt"
	"path"
	"regexp"
	"sync"
	"time"

	"github.com/Originate/exosphere/exo-go/src/app_config_helpers"
	"github.com/Originate/exosphere/exo-go/src/app_dependency_helpers"
	"github.com/Originate/exosphere/exo-go/src/docker_compose"
	"github.com/Originate/exosphere/exo-go/src/logger"
	"github.com/Originate/exosphere/exo-go/src/service_config_helpers"
	"github.com/Originate/exosphere/exo-go/src/service_restarter"
	"github.com/Originate/exosphere/exo-go/src/types"
	execplus "github.com/Originate/go-execplus"
	"github.com/fatih/color"
	"github.com/pkg/errors"
)

// AppRunner runs the overall application
type AppRunner struct {
	AppConfig        types.AppConfig
	Logger           *logger.Logger
	AppDir           string
	homeDir          string
	Env              map[string]string
	DockerComposeDir string
	logChannel       chan string
}

// NewAppRunner is AppRunner's constructor
func NewAppRunner(appConfig types.AppConfig, logger *logger.Logger, appDir, homeDir string) *AppRunner {
	return &AppRunner{
		AppConfig:        appConfig,
		Logger:           logger,
		AppDir:           appDir,
		homeDir:          homeDir,
		Env:              appConfigHelpers.GetEnvironmentVariables(appConfig, appDir, homeDir),
		DockerComposeDir: path.Join(appDir, "tmp"),
		logChannel:       logger.GetLogChannel("exo-run"),
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

func (a *AppRunner) runImages(imageNames []string, imageOnlineTexts map[string]string, identifier string) error {
	cmdPlus, err := dockerCompose.RunImages(imageNames, a.getEnv(), a.DockerComposeDir, a.logChannel)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to run %s\nOutput: %s\nError: %s\n", identifier, cmdPlus.Output, err))
	}
	var wg sync.WaitGroup
	var onlineTextRegex *regexp.Regexp
	for role, onlineText := range imageOnlineTexts {
		wg.Add(1)
		onlineTextRegex, err = regexp.Compile(fmt.Sprintf("%s.*%s", role, onlineText))
		if err != nil {
			return err
		}
		go func(role string, onlineTextRegex *regexp.Regexp) {
			a.waitForOnlineText(cmdPlus, role, onlineTextRegex)
			wg.Done()
		}(role, onlineTextRegex)
	}
	wg.Wait()
	a.logChannel <- fmt.Sprintf("all %s online", identifier)
	return nil
}

// Shutdown shuts down the application and returns the process output and an error if any
func (a *AppRunner) Shutdown(shutdownConfig types.ShutdownConfig) error {
	if len(shutdownConfig.ErrorMessage) > 0 {
		color.Red(shutdownConfig.ErrorMessage)
	} else {
		fmt.Printf("\n\n%s", shutdownConfig.CloseMessage)
	}
	process, err := dockerCompose.KillAllContainers(a.DockerComposeDir, a.logChannel)
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
	if err := a.runImages(dependencyNames, a.compileDependencyOnlineTexts(serviceConfigs), "dependencies"); err != nil {
		return err
	}
	if err := a.runImages(serviceNames, a.compileServiceOnlineTexts(serviceConfigs), "services"); err != nil {
		return err
	}
	a.watchServices()
	return nil
}

func (a *AppRunner) waitForOnlineText(cmdPlus *execplus.CmdPlus, role string, onlineTextRegex *regexp.Regexp) {
	_ = cmdPlus.WaitForRegexp(onlineTextRegex, time.Hour) // No user will actually wait this long
	err := a.Logger.Log(role, fmt.Sprintf("'%s' is running", role), true)
	if err != nil {
		fmt.Printf("Error logging '%s' as online: %v\n", role, err)
	}
}

func (a *AppRunner) watchServices() {
	watcherErrChannel := make(chan error)
	go func() {
		err := <-watcherErrChannel
		if err != nil {
			closeMessage := fmt.Sprintf("Error watching services for changes: %v", err)
			if err := a.Shutdown(types.ShutdownConfig{CloseMessage: closeMessage}); err != nil {
				a.logChannel <- "Failed to shutdown"
			}
		}
	}()
	for serviceName, data := range serviceConfigHelpers.GetServiceData(a.AppConfig.Services) {
		if data.Location != "" {
			restarter := serviceRestarter.ServiceRestarter{
				ServiceName:      serviceName,
				ServiceDir:       data.Location,
				DockerComposeDir: a.DockerComposeDir,
				LogChannel:       a.logChannel,
				Env:              a.getEnv(),
			}
			restarter.Watch(watcherErrChannel)
		}
	}
}

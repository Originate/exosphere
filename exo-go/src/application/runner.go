package application

import (
	"fmt"
	"path"
	"regexp"
	"sync"
	"time"

	"github.com/Originate/exosphere/exo-go/src/config"
	"github.com/Originate/exosphere/exo-go/src/docker"
	"github.com/Originate/exosphere/exo-go/src/types"
	execplus "github.com/Originate/go-execplus"
	"github.com/fatih/color"
	"github.com/pkg/errors"
)

// Runner runs the overall application
type Runner struct {
	AppConfig         types.AppConfig
	ServiceConfigs    map[string]types.ServiceConfig
	BuiltDependencies map[string]config.AppDependency
	Env               map[string]string
	DockerComposeDir  string
	Logger            *Logger
	logChannel        chan string
}

// NewRunner is Runner's constructor
func NewRunner(appConfig types.AppConfig, logger *Logger, logRole, appDir, homeDir string) (*Runner, error) {
	serviceConfigs, err := config.GetServiceConfigs(appDir, appConfig)
	if err != nil {
		return &Runner{}, err
	}
	allBuiltDependencies := config.GetAllBuiltDependencies(appConfig, serviceConfigs, appDir, homeDir)
	appBuiltDependencies := config.GetAppBuiltDependencies(appConfig, appDir, homeDir)
	return &Runner{
		AppConfig:         appConfig,
		ServiceConfigs:    serviceConfigs,
		BuiltDependencies: allBuiltDependencies,
		Env:               config.GetEnvironmentVariables(appBuiltDependencies),
		DockerComposeDir:  path.Join(appDir, "tmp"),
		Logger:            logger,
		logChannel:        logger.GetLogChannel(logRole),
	}, nil
}

func (r *Runner) compileServiceOnlineTexts() map[string]string {
	onlineTexts := make(map[string]string)
	for serviceName, serviceConfig := range r.ServiceConfigs {
		onlineTexts[serviceName] = serviceConfig.Startup["online-text"]
	}
	return onlineTexts
}

func (r *Runner) compileDependencyOnlineTexts() map[string]string {
	onlineTexts := make(map[string]string)
	for dependencyName, builtDependency := range r.BuiltDependencies {
		onlineTexts[dependencyName] = builtDependency.GetOnlineText()
	}
	return onlineTexts
}

func (r *Runner) getDependencyContainerNames() []string {
	result := []string{}
	for _, builtDependency := range r.BuiltDependencies {
		result = append(result, builtDependency.GetContainerName())
	}
	return result
}

func (r *Runner) getEnv() []string {
	formattedEnvVars := []string{}
	for variable, value := range r.Env {
		formattedEnvVars = append(formattedEnvVars, fmt.Sprintf("%s=%s", variable, value))
	}
	return formattedEnvVars
}

func (r *Runner) runImages(imageNames []string, imageOnlineTexts map[string]string, identifier string) (string, error) {
	cmdPlus, err := docker.RunImages(imageNames, r.getEnv(), r.DockerComposeDir, r.logChannel)
	if err != nil {
		return cmdPlus.GetOutput(), errors.Wrap(err, fmt.Sprintf("Failed to run %s\nOutput: %s\nError: %s\n", identifier, cmdPlus.GetOutput(), err))
	}
	var wg sync.WaitGroup
	var onlineTextRegex *regexp.Regexp
	for role, onlineText := range imageOnlineTexts {
		wg.Add(1)
		onlineTextRegex, err = regexp.Compile(fmt.Sprintf("%s.*%s", role, onlineText))
		if err != nil {
			return cmdPlus.GetOutput(), err
		}
		go func(role string, onlineTextRegex *regexp.Regexp) {
			r.waitForOnlineText(cmdPlus, role, onlineTextRegex)
			wg.Done()
		}(role, onlineTextRegex)
	}
	wg.Wait()
	r.logChannel <- fmt.Sprintf("all %s online", identifier)
	return cmdPlus.GetOutput(), nil
}

// Shutdown shuts down the application and returns the process output and an error if any
func (r *Runner) Shutdown(shutdownConfig types.ShutdownConfig) error {
	if len(shutdownConfig.ErrorMessage) > 0 {
		color.Red(shutdownConfig.ErrorMessage)
	} else {
		fmt.Printf("\n\n%s", shutdownConfig.CloseMessage)
	}
	process, err := docker.KillAllContainers(r.DockerComposeDir, r.logChannel)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to shutdown the app\nOutput: %s\nError: %s\n", process.GetOutput(), err))
	}
	err = process.Wait()
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to shutdown the app\nOutput: %s\nError: %s\n", process.GetOutput(), err))
	}
	return nil
}

// Start runs the application and returns the process and returns an error if any
func (r *Runner) Start() error {
	dependencyNames := r.getDependencyContainerNames()
	serviceNames := r.AppConfig.GetServiceNames()
	if len(dependencyNames) > 0 {
		if _, err := r.runImages(dependencyNames, r.compileDependencyOnlineTexts(), "dependencies"); err != nil {
			return err
		}
	}
	if len(serviceNames) > 0 {
		if _, err := r.runImages(serviceNames, r.compileServiceOnlineTexts(), "services"); err != nil {
			return err
		}
	}
	r.watchServices()
	return nil
}

func (r *Runner) waitForOnlineText(cmdPlus *execplus.CmdPlus, role string, onlineTextRegex *regexp.Regexp) {
	err := cmdPlus.WaitForRegexp(onlineTextRegex, time.Hour) // No user will actually wait this long
	if err != nil {
		fmt.Printf("'%s' failed to come online after an hour", role)
	}
	if role == "" {
		return
	}
	err = r.Logger.Log(role, fmt.Sprintf("'%s' is running", role), true)
	if err != nil {
		fmt.Printf("Error logging '%s' as online: %v\n", role, err)
	}
}

func (r *Runner) watchServices() {
	watcherErrChannel := make(chan error)
	go func() {
		err := <-watcherErrChannel
		if err != nil {
			closeMessage := fmt.Sprintf("Error watching services for changes: %v", err)
			if err := r.Shutdown(types.ShutdownConfig{CloseMessage: closeMessage}); err != nil {
				r.logChannel <- "Failed to shutdown"
			}
		}
	}()
	for serviceName, data := range r.AppConfig.GetServiceData() {
		if data.Location != "" {
			restarter := serviceRestarter{
				ServiceName:      serviceName,
				ServiceDir:       data.Location,
				DockerComposeDir: r.DockerComposeDir,
				LogChannel:       r.logChannel,
				Env:              r.getEnv(),
			}
			restarter.Watch(watcherErrChannel)
		}
	}
}

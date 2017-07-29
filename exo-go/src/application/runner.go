package application

import (
	"fmt"
	"path"
	"regexp"
	"sync"
	"time"

	"github.com/Originate/exosphere/exo-go/src/config"
	"github.com/Originate/exosphere/exo-go/src/docker_compose"
	"github.com/Originate/exosphere/exo-go/src/logger"
	"github.com/Originate/exosphere/exo-go/src/service_restarter"
	"github.com/Originate/exosphere/exo-go/src/types"
	execplus "github.com/Originate/go-execplus"
	"github.com/fatih/color"
	"github.com/pkg/errors"
)

// Runner runs the overall application
type Runner struct {
	AppConfig        types.AppConfig
	Logger           *logger.Logger
	AppDir           string
	homeDir          string
	Env              map[string]string
	DockerComposeDir string
	logChannel       chan string
}

// NewRunner is Runner's constructor
func NewRunner(appConfig types.AppConfig, logger *logger.Logger, appDir, homeDir string) *Runner {
	return &Runner{
		AppConfig:        appConfig,
		Logger:           logger,
		AppDir:           appDir,
		homeDir:          homeDir,
		Env:              config.GetEnvironmentVariables(appConfig, appDir, homeDir),
		DockerComposeDir: path.Join(appDir, "tmp"),
		logChannel:       logger.GetLogChannel("exo-run"),
	}
}

func (r *Runner) compileServiceOnlineTexts(serviceConfigs map[string]types.ServiceConfig) map[string]string {
	onlineTexts := make(map[string]string)
	for serviceName, serviceConfig := range serviceConfigs {
		onlineTexts[serviceName] = serviceConfig.Startup["online-text"]
	}
	return onlineTexts
}

func (r *Runner) compileDependencyOnlineTexts(serviceConfigs map[string]types.ServiceConfig) map[string]string {
	onlineTexts := make(map[string]string)
	for _, dependency := range r.AppConfig.Dependencies {
		onlineTexts[dependency.Name] = config.NewAppDependency(dependency, r.AppConfig, r.AppDir, r.homeDir).GetOnlineText()
	}
	for _, serviceConfig := range serviceConfigs {
		for _, dependency := range serviceConfig.Dependencies {
			if !dependency.Config.IsEmpty() {
				onlineTexts[dependency.Name] = config.NewAppDependency(dependency, r.AppConfig, r.AppDir, r.homeDir).GetOnlineText()
			}
		}
	}
	return onlineTexts
}

func (r *Runner) getEnv() []string {
	formattedEnvVars := []string{}
	for variable, value := range r.Env {
		formattedEnvVars = append(formattedEnvVars, fmt.Sprintf("%s=%s", variable, value))
	}
	return formattedEnvVars
}

func (r *Runner) runImages(imageNames []string, imageOnlineTexts map[string]string, identifier string) error {
	cmdPlus, err := dockerCompose.RunImages(imageNames, r.getEnv(), r.DockerComposeDir, r.logChannel)
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
			r.waitForOnlineText(cmdPlus, role, onlineTextRegex)
			wg.Done()
		}(role, onlineTextRegex)
	}
	wg.Wait()
	r.logChannel <- fmt.Sprintf("all %s online", identifier)
	return nil
}

// Shutdown shuts down the application and returns the process output and an error if any
func (r *Runner) Shutdown(shutdownConfig types.ShutdownConfig) error {
	if len(shutdownConfig.ErrorMessage) > 0 {
		color.Red(shutdownConfig.ErrorMessage)
	} else {
		fmt.Printf("\n\n%s", shutdownConfig.CloseMessage)
	}
	process, err := dockerCompose.KillAllContainers(r.DockerComposeDir, r.logChannel)
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
func (r *Runner) Start() error {
	dependencyNames, err := config.GetAllDependencyNames(r.AppDir, r.AppConfig)
	if err != nil {
		return err
	}
	serviceNames := r.AppConfig.GetServiceNames()
	serviceConfigs, err := config.GetServiceConfigs(r.AppDir, r.AppConfig)
	if err != nil {
		return err
	}
	if err := r.runImages(dependencyNames, r.compileDependencyOnlineTexts(serviceConfigs), "dependencies"); err != nil {
		return err
	}
	if err := r.runImages(serviceNames, r.compileServiceOnlineTexts(serviceConfigs), "services"); err != nil {
		return err
	}
	r.watchServices()
	return nil
}

func (r *Runner) waitForOnlineText(cmdPlus *execplus.CmdPlus, role string, onlineTextRegex *regexp.Regexp) {
	err := cmdPlus.WaitForRegexp(onlineTextRegex, time.Hour) // No user will actually wait this long
	if err != nil {
		fmt.Printf("'%s' failed to come online after an hour", role)
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
			restarter := serviceRestarter.ServiceRestarter{
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

package application

import (
	"fmt"
	"path"
	"regexp"
	"sync"
	"time"

	"github.com/Originate/exosphere/src/config"
	"github.com/Originate/exosphere/src/docker/compose"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/util"
	execplus "github.com/Originate/go-execplus"
	"github.com/fatih/color"
	"github.com/pkg/errors"
)

// Runner runs the overall application
type Runner struct {
	AppConfig                types.AppConfig
	ServiceConfigs           map[string]types.ServiceConfig
	BuiltDependencies        map[string]config.AppDevelopmentDependency
	DockerComposeDir         string
	DockerComposeProjectName string
	logger                   *util.Logger
}

// NewRunner is Runner's constructor
func NewRunner(appConfig types.AppConfig, logger *util.Logger, appDir, homeDir, dockerComposeProjectName string) (*Runner, error) {
	serviceConfigs, err := config.GetServiceConfigs(appDir, appConfig)
	if err != nil {
		return &Runner{}, err
	}
	allBuiltDependencies := config.GetBuiltDevelopmentDependencies(appConfig, serviceConfigs, appDir, homeDir)
	return &Runner{
		AppConfig:                appConfig,
		ServiceConfigs:           serviceConfigs,
		BuiltDependencies:        allBuiltDependencies,
		DockerComposeDir:         path.Join(appDir, "tmp"),
		DockerComposeProjectName: dockerComposeProjectName,
		logger: logger,
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

func (r *Runner) runImages(imageNames []string, imageOnlineTexts map[string]string, identifier string) (string, error) {
	cmdPlus, err := compose.RunImages(compose.ImagesOptions{
		DockerComposeDir: r.DockerComposeDir,
		ImageNames:       imageNames,
		Logger:           r.logger,
		Env:              []string{fmt.Sprintf("COMPOSE_PROJECT_NAME=%s", r.DockerComposeProjectName)},
	})
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
	r.logger.Logf("all %s online", identifier)
	return cmdPlus.GetOutput(), nil
}

// Shutdown shuts down the application and returns the process output and an error if any
func (r *Runner) Shutdown(shutdownConfig types.ShutdownConfig) error {
	if len(shutdownConfig.ErrorMessage) > 0 {
		r.logger.Log(color.New(color.FgRed).Sprint(shutdownConfig.ErrorMessage))
	} else {
		r.logger.Log(shutdownConfig.CloseMessage)
	}
	err := compose.KillAllContainers(compose.BaseOptions{
		DockerComposeDir: r.DockerComposeDir,
		Logger:           r.logger,
		Env:              []string{fmt.Sprintf("COMPOSE_PROJECT_NAME=%s", r.DockerComposeProjectName)},
	})
	if err != nil {
		return errors.Wrap(err, "Failed to shutdown the app")
	}
	return nil
}

// Start runs the application and returns the process and returns an error if any
func (r *Runner) Start() error {
	dependencyNames := r.getDependencyContainerNames()
	serviceNames := r.AppConfig.GetSortedServiceNames()
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
	r.logger.Logf("'%s' is running", role)
}

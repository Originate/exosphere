package application

import (
	"fmt"
	"os"
	"os/signal"
	"path"
	"sync"

	"github.com/Originate/exosphere/src/config"
	"github.com/Originate/exosphere/src/docker/composebuilder"
	"github.com/Originate/exosphere/src/docker/composerunner"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/util"
)

// Runner runs the overall application
type Runner struct {
	AppConfig                types.AppConfig
	AppDir                   string
	HomeDir                  string
	ServiceConfigs           map[string]types.ServiceConfig
	BuiltDependencies        map[string]config.AppDevelopmentDependency
	DockerComposeDir         string
	DockerComposeProjectName string
	logger                   *util.Logger
	BuildMode                composebuilder.BuildMode
}

// NewRunner is Runner's constructor
func NewRunner(appConfig types.AppConfig, logger *util.Logger, appDir, homeDir, dockerComposeProjectName string, buildMode composebuilder.BuildMode) (*Runner, error) {
	serviceConfigs, err := config.GetServiceConfigs(appDir, appConfig)
	if err != nil {
		return &Runner{}, err
	}
	allBuiltDependencies := config.GetBuiltDevelopmentDependencies(appConfig, serviceConfigs, appDir, homeDir)
	return &Runner{
		AppDir:                   appDir,
		HomeDir:                  homeDir,
		AppConfig:                appConfig,
		ServiceConfigs:           serviceConfigs,
		BuiltDependencies:        allBuiltDependencies,
		DockerComposeDir:         path.Join(appDir, "tmp"),
		DockerComposeProjectName: dockerComposeProjectName,
		logger:    logger,
		BuildMode: buildMode,
	}, nil
}

func (r *Runner) compileServiceOnlineTexts() map[string]string {
	onlineTexts := make(map[string]string)
	for serviceRole, serviceConfig := range r.ServiceConfigs {
		onlineTexts[serviceRole] = serviceConfig.Startup["online-text"]
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

// Run runs the application with graceful shutdown
func (r *Runner) Run() error {
	dockerConfigs, err := composebuilder.GetApplicationDockerConfigs(composebuilder.ApplicationOptions{
		AppConfig: r.AppConfig,
		AppDir:    r.AppDir,
		BuildMode: r.BuildMode,
		HomeDir:   r.HomeDir,
	})
	if err != nil {
		return err
	}
	runOptions := composerunner.RunOptions{
		ImageGroups: []composerunner.ImageGroup{
			{
				ID:          "dependencies",
				Names:       r.getDependencyContainerNames(),
				OnlineTexts: r.compileDependencyOnlineTexts(),
			},
			{
				ID:          "services",
				Names:       r.AppConfig.GetSortedServiceRoles(),
				OnlineTexts: r.compileServiceOnlineTexts(),
			},
		},
		DockerConfigs:            dockerConfigs,
		DockerComposeDir:         r.DockerComposeDir,
		DockerComposeProjectName: r.DockerComposeProjectName,
		Logger: r.logger,
	}
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		shutdownChannel := make(chan os.Signal, 1)
		signal.Notify(shutdownChannel, os.Interrupt)
		<-shutdownChannel
		signal.Stop(shutdownChannel)
		if shutdownErr := composerunner.Shutdown(runOptions); shutdownErr != nil {
			fmt.Println("Error shutting down")
		}
		wg.Done()
	}()
	err = composerunner.Run(runOptions)
	if err != nil {
		_ = composerunner.Shutdown(runOptions)
		return err
	}
	wg.Wait()
	return nil
}

package appRunner

import (
	"fmt"
	"log"
	"os"
	"path"
	"sync"

	yaml "gopkg.in/yaml.v2"

	"github.com/Originate/exosphere/exo-go/src/docker_compose"
	"github.com/Originate/exosphere/exo-go/src/docker_helpers"
	"github.com/Originate/exosphere/exo-go/src/logger"
	"github.com/Originate/exosphere/exo-go/src/process_helpers"
	"github.com/Originate/exosphere/exo-go/src/service_config_helpers"
	"github.com/Originate/exosphere/exo-go/src/service_helpers"
	"github.com/Originate/exosphere/exo-go/src/types"
	"github.com/docker/docker/client"
	"github.com/fatih/color"
	"github.com/pkg/errors"
)

// AppRunner runs the overall application
type AppRunner struct {
	AppConfig            types.AppConfig
	Logger               *logger.Logger
	Env                  []string
	DockerConfigLocation string
	Cwd                  string
	OnlineTexts          map[string]string
}

// NewAppRunner is AppRunner's constructor
func NewAppRunner(appConfig types.AppConfig, logger *logger.Logger) *AppRunner {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get the current path: %s", err)
	}
	appRunner := &AppRunner{AppConfig: appConfig, Logger: logger, DockerConfigLocation: path.Join(cwd, "tmp"), Cwd: cwd}
	for _, dependency := range appConfig.Dependencies {
		for variable, value := range dependency.GetEnvVariables() {
			appRunner.Env = append(appRunner.Env, fmt.Sprintf("%s=%s", variable, value))
		}
	}
	return appRunner
}

// Start runs the application
func (appRunner *AppRunner) Start() {
	_, stdoutBuffer, err := dockerCompose.RunAllImages(appRunner.Env, appRunner.DockerConfigLocation, appRunner.Write)
	if err != nil {
		appRunner.Shutdown("", "Failed to run images")
	} else {
		appRunner.compileOnlineText(func(err error) {
			if err != nil {
				log.Fatal(err)
			} else {
				wg := new(sync.WaitGroup)
				for role, onlineText := range appRunner.OnlineTexts {
					wg.Add(1)
					go func(role string, stdoutBuffer fmt.Stringer, onlineText string) {
						defer wg.Done()
						processHelpers.Wait(stdoutBuffer, onlineText, func() {
							appRunner.Logger.Log(role, fmt.Sprintf("'%s' is running", role), true)
						})
					}(role, stdoutBuffer, onlineText)
				}
				wg.Wait()
				appRunner.Write("all services online")
			}
		})
	}
}

// Shutdown shuts down the application
func (appRunner *AppRunner) Shutdown(closeMessage, errorMessage string) {
	var exitCode int
	if len(errorMessage) > 0 {
		color.Red(errorMessage)
		exitCode = 1
	} else {
		fmt.Printf("\n\n%s", closeMessage)
		exitCode = 0
	}
	_, _, err := dockerCompose.KillAllContainers(appRunner.Env, appRunner.DockerConfigLocation, appRunner.Write)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(exitCode)
}

// Write logs exo-run output
func (appRunner *AppRunner) Write(text string) {
	appRunner.Logger.Log("exo-run", text, true)
}

func (appRunner *AppRunner) compileOnlineText(done func(error)) {
	appRunner.OnlineTexts = make(map[string]string)
	for _, dependency := range appRunner.AppConfig.Dependencies {
		appRunner.OnlineTexts[dependency.Name] = dependency.GetOnlineText()
	}
	wg := new(sync.WaitGroup)
	for service, serviceData := range serviceHelpers.GetServiceConfigs(appRunner.AppConfig.Services) {
		wg.Add(1)
		go func(service string, serviceData types.ServiceConfig) {
			defer wg.Done()
			appRunner.getOnlineText(service, serviceData, func(err error) {
				if err != nil {
					done(err)
				}
			})
		}(service, serviceData)
	}
	wg.Wait()
	done(nil)
}

func (appRunner *AppRunner) getOnlineText(role string, serviceData types.ServiceConfig, done func(error)) {
	if len(serviceData.Location) > 0 {
		serviceConfig := serviceConfigHelpers.GetServiceConfig(serviceData.Location)
		appRunner.OnlineTexts[role] = serviceConfig.Startup["online-text"]
		done(nil)
	} else if len(serviceData.DockerImage) > 0 {
		c, err := client.NewEnvClient()
		if err != nil {
			panic(err)
		}
		dockerHelpers.CatFile(c, serviceData.DockerImage, "service.yml", func(err error, yamlFile []byte) {
			var serviceConfig types.ServiceConfig
			err = yaml.Unmarshal(yamlFile, &serviceConfig)
			if err != nil {
				done(errors.Wrap(err, fmt.Sprintf("Failed to unmarshal service.yml for the service %s", role)))
			}
			appRunner.OnlineTexts[role] = serviceConfig.Startup["online-text"]
			done(nil)
		})
	} else {
		done(fmt.Errorf("No location or docker image listed for %s", role))
	}
}

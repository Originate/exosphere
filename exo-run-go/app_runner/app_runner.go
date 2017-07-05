package appRunner

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/Originate/exosphere/exo-run-go/docker_compose"
	"github.com/Originate/exosphere/exo-run-go/logger"
	"github.com/Originate/exosphere/exo-run-go/types"
	"github.com/chuckpreslar/emission"
)

type AppRunner struct {
	AppConfig            types.AppConfig
	Logger               logger.Logger
	Env                  map[string]interface{}
	DockerConfigLocation string
	emission.Emitter
	Cwd string
}

func NewAppRunner(appConfig types.AppConfig, logger logger.Logger) *AppRunner {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current path: %s", err)
	}
	appRunner := &AppRunner{AppConfig: appConfig, Logger: logger, Env: make(map[string]interface{}), DockerConfigLocation: path.Join(cwd, "tmp"), Cwd: cwd}

	for _, dependency := range appConfig.Dependencies {
		for variable, value := range dependency.GetEnvVariables() {
			appRunner.Env[variable] = value
		}
	}
	return appRunner
}

func (appRunner *AppRunner) Start() {
	dockerCompose.RunAllImages(appRunner.Env, appRunner.Cwd, appRunner.Write)
}

func (appRunner *AppRunner) Shutdown(closeMessage, errorMessage string) {
	var exitCode int
	if len(errorMessage) > 0 {
		fmt.Println(errorMessage)
		exitCode = 1
	} else {
		fmt.Println(closeMessage)
		exitCode = 0
	}
	dockerCompose.KillAllContainers(appRunner.Env, appRunner.DockerConfigLocation, appRunner.Write)
	os.Exit(exitCode)
}

func (appRunner *AppRunner) Write(text string) {
	appRunner.Logger.Log("exo-run", text, true)
}

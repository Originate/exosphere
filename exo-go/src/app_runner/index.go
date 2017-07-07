package appRunner

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/Originate/exosphere/exo-go/src/docker_compose"
	"github.com/Originate/exosphere/exo-go/src/logger"
	"github.com/Originate/exosphere/exo-go/src/types"
)

// AppRunner runs the overall application
type AppRunner struct {
	AppConfig            types.AppConfig
	Logger               *logger.Logger
	Env                  []string
	DockerConfigLocation string
	Cwd                  string
	OnlineTexts          []string
}

// NewAppRunner is AppRunner's constructor
func NewAppRunner(appConfig types.AppConfig, logger *logger.Logger) *AppRunner {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current path: %s", err)
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
	dockerCompose.RunAllImages(appRunner.Env, appRunner.DockerConfigLocation, appRunner.Write)
}

// Shutdown shuts down the application
func (appRunner *AppRunner) Shutdown(closeMessage, errorMessage string) {
	fmt.Println("shutting down")
	// var exitCode int
	// if len(errorMessage) > 0 {
	// 	fmt.Println(errorMessage)
	// 	exitCode = 1
	// } else {
	// 	fmt.Println(closeMessage)
	// 	exitCode = 0
	// }
	dockerCompose.KillAllContainers(appRunner.Env, appRunner.DockerConfigLocation, appRunner.Write)
	// os.Exit(exitCode)
}

// @_compile-online-text (err) ~>
//   | err => throw err
//   asynchronizer = new Asynchronizer Object.keys(@online-texts)
//   for role, online-text of @online-texts
//     let role, online-text
//       @process.wait (new RegExp(role + ".*" + online-text)), ~>
//         @logger.log {role, text: "'#{role}' is running"}
//         asynchronizer.check role
//   asynchronizer.then ~>
//     @write 'all services online'

// Write logs exo-run output
func (appRunner *AppRunner) Write(text string) {
	appRunner.Logger.Log("exo-run", text, true)
}

// _compile-online-text: (done) ~>
//   @online-texts = {}
//   for app-dependency in @app-config.dependencies
//     dependency = ApplicationDependency.build app-dependency
//     @online-texts[app-dependency.name] = dependency.get-online-text!
//   services = []
//   for protection-level of @app-config.services
//     for role, service-data of @app-config.services[protection-level]
//       services.push {role: role, service-data: service-data}
//   async.map-series services, @_get-online-text, (err) ~>
//     | err => done err
//     done!

func (appRunner *AppRunner) compileOnlineText(done interface{}) {
	for _, dependency := range appRunner.AppConfig.Dependencies {
		appRunner.OnlineTexts = append(appRunner.OnlineTexts, dependency.GetOnlineText())
	}
	// TODO: get service dependencies' online texts
}

// _get-online-text: ({role, service-data}, done) ~>
//   | service-data.location =>
//     service-config = yaml.safe-load fs.read-file-sync(path.join(process.cwd!, service-data.location, 'service.yml'))
//     @online-texts[role] = service-config.startup['online-text']
//     done!
//   | service-data['docker-image'] =>
//     DockerHelper.cat-file image: service-data['docker-image'], file-name: 'service.yml', (err, external-service-config) ~>
//       | err => done err
//       service-config = yaml.safe-load external-service-config
//       @online-texts[role] = service-config.startup['online-text']
//       done!
//   | otherwise => done new Error red "No location or docker image listed for '#{role}'"

func (appRunner *AppRunner) getOnlineText(role string, serviceConfig types.ServiceConfig, done interface{}) {
	if len(serviceConfig.Location) > 0 {
	}
}

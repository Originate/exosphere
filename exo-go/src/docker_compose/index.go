package dockerCompose

import "github.com/Originate/exosphere/exo-go/src/process_helpers"

// RunAllImages runs all the docker images
func RunAllImages(env []string, cwd string, log func(string)) error {
	process := processHelpers.NewProcess("docker-compose", "up")
	process.SetDir(cwd)
	process.SetEnv(env)
	process.SetStdoutLog(log)
	return process.Start()
}

// KillAllContainers kills all the containers
func KillAllContainers(env []string, cwd string, log func(string)) error {
	process := processHelpers.NewProcess("docker-compose", "down")
	process.SetDir(cwd)
	process.SetEnv(env)
	process.SetStdoutLog(log)
	return process.Start()
}

package dockerCompose

import "github.com/Originate/exosphere/exo-go/src/process_helpers"

// BuildAllImages builds all the docker images defined in docker-compose.yml
func BuildAllImages(dockerComposeLocation string, log func(string)) error {
	return processHelpers.RunAndLog(dockerComposeLocation, log, "docker-compose", "build")
}

// KillAllContainers kills all the containers
func KillAllContainers(cwd string, log func(string)) (*processHelpers.Process, error) {
	process := processHelpers.NewProcess("docker-compose", "down")
	process.SetDir(cwd)
	process.SetStdoutLog(log)
	return process, process.Start()
}

// PullAllImages pulls all the docker images defined in docker-compose.yml
func PullAllImages(dockerComposeLocation string, log func(string)) error {
	return processHelpers.RunAndLog(dockerComposeLocation, log, "docker-compose", "pull")
}

// RunAllImages runs all the docker images
func RunAllImages(env []string, cwd string, log func(string)) (*processHelpers.Process, error) {
	process := processHelpers.NewProcess("docker-compose", "up")
	process.SetDir(cwd)
	process.SetEnv(env)
	process.SetStdoutLog(log)
	return process, process.Start()
}

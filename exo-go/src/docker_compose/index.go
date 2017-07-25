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
	process.AddOutputFunc("log", log)
	return process, process.Start()
}

// PullAllImages pulls all the docker images defined in docker-compose.yml
func PullAllImages(dockerComposeLocation string, log func(string)) error {
	return processHelpers.RunAndLog(dockerComposeLocation, log, "docker-compose", "pull")
}

// RunImages runs the given docker images
func RunImages(images []string, env []string, cwd string, log func(string)) (*processHelpers.Process, error) {
	process := processHelpers.NewProcess(append([]string{"docker-compose", "up"}, images...)...)
	process.SetDir(cwd)
	process.SetEnv(env)
	process.AddOutputFunc("log", log)
	return process, process.Start()
}

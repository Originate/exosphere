package dockerCompose

import "github.com/Originate/exosphere/exo-go/src/process_helpers"

// BuildAllImages builds all the docker images defined in docker-compose.yml
func BuildAllImages(dockerComposeDir string, log func(string)) error {
	return processHelpers.RunAndLog(dockerComposeDir, []string{}, log, "docker-compose", "build")
}

// CreateNewContainer creates a new docker container for the given service
func CreateNewContainer(serviceName string, env []string, dockerComposeDir string, log func(string)) error {
	return processHelpers.RunAndLog(dockerComposeDir, env, log, "docker-compose", "create", "--build", serviceName)
}

// KillAllContainers kills all the containers
func KillAllContainers(dockerComposeDir string, log func(string)) (*processHelpers.Process, error) {
	process := processHelpers.NewProcess("docker-compose", "down")
	process.SetDir(dockerComposeDir)
	process.AddOutputFunc("log", log)
	return process, process.Start()
}

// KillContainer kills the docker container of the given service
func KillContainer(serviceName, dockerComposeDir string, log func(string)) error {
	return processHelpers.RunAndLog(dockerComposeDir, []string{}, log, "docker-compose", "kill", serviceName)
}

// PullAllImages pulls all the docker images defined in docker-compose.yml
func PullAllImages(dockerComposeDir string, log func(string)) error {
	return processHelpers.RunAndLog(dockerComposeDir, []string{}, log, "docker-compose", "pull")
}

// RunImages runs the given docker images
func RunImages(images []string, env []string, dockerComposeDir string, log func(string)) (*processHelpers.Process, error) {
	process := processHelpers.NewProcess(append([]string{"docker-compose", "up"}, images...)...)
	process.SetDir(dockerComposeDir)
	process.SetEnv(env)
	process.AddOutputFunc("log", log)
	return process, process.Start()
}

// StartContainer starts the docker container of the given service
func StartContainer(serviceName string, env []string, dockerComposeDir string, log func(string)) error {
	return processHelpers.RunAndLog(dockerComposeDir, env, log, "docker-compose", "restart", serviceName)
}

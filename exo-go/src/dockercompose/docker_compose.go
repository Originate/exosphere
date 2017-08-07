package dockercompose

// DockerCompose represents the docker compose object
type DockerCompose struct {
	Version  string
	Services DockerConfigs
}

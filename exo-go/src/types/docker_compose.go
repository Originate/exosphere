package types

// DockerCompose represents the docker compose object
type DockerCompose struct {
	Version  string
	Services map[string]Service
}

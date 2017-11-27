package types

// DockerCompose represents the docker compose object
type DockerCompose struct {
	Version  string
	Services DockerConfigs
	Volumes  map[string]interface{}
}

// NewDockerCompose returns a docker compose object
func NewDockerCompose() *DockerCompose {
	return &DockerCompose{
		Version:  "3",
		Services: DockerConfigs{},
		Volumes:  map[string]interface{}{},
	}
}

// Merge joins the given docker compose objects into one
func (d *DockerCompose) Merge(objs ...*DockerCompose) *DockerCompose {
	result := NewDockerCompose()
	for _, obj := range append(objs, d) {
		result.Services = result.Services.Merge(obj.Services)
		for key, val := range obj.Volumes {
			result.Volumes[key] = val
		}
	}
	return result
}

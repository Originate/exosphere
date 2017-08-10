package types

// DockerConfigs represents a map of DockerConfig structs
type DockerConfigs map[string]DockerConfig

// Merge joins the given docker config maps into one
func (d DockerConfigs) Merge(maps ...DockerConfigs) DockerConfigs {
	result := DockerConfigs{}
	for _, currentMap := range append(maps, d) {
		for key, val := range currentMap {
			result[key] = val
		}
	}
	return result
}

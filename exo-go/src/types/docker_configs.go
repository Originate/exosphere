package types

// DockerConfigMap represents the configuration
type DockerConfigMap map[string]DockerConfig

// Merge joins the given docker config maps into one
func (d DockerConfigMap) Merge(maps ...DockerConfigMap) (result DockerConfigMap) {
	for _, currentMap := range append(maps, d) {
		for key, val := range currentMap {
			result[key] = val
		}
	}
	return result
}

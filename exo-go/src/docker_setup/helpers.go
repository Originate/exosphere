package dockerSetup

import "github.com/Originate/exosphere/exo-go/src/types"

func joinDockerConfigMaps(maps ...map[string]types.DockerConfig) map[string]types.DockerConfig {
	result := map[string]types.DockerConfig{}
	for _, currentMap := range maps {
		for key, val := range currentMap {
			result[key] = val
		}
	}
	return result
}

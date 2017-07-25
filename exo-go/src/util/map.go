package util

import "github.com/Originate/exosphere/exo-go/src/types"

// JoinDockerConfigMaps joins the given docker config maps
func JoinDockerConfigMaps(maps ...map[string]types.DockerConfig) map[string]types.DockerConfig {
	result := map[string]types.DockerConfig{}
	for _, currentMap := range maps {
		for key, val := range currentMap {
			result[key] = val
		}
	}
	return result
}

// JoinStringMaps joins the given string maps
func JoinStringMaps(maps ...map[string]string) map[string]string {
	result := map[string]string{}
	for _, currentMap := range maps {
		for key, val := range currentMap {
			result[key] = val
		}
	}
	return result
}

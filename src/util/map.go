package util

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

// Merge merges two maps. It overwrites existing keys of map1 if conflicting keys exist
func Merge(map1 map[string]string, map2 map[string]string) {
	for k, v := range map2 {
		map1[k] = v
	}
}

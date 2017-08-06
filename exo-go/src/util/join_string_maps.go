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

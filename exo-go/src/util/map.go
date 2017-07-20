package util

// JoinMaps joins the given maps
func JoinMaps(maps ...map[string]string) map[string]string {
	result := map[string]string{}
	for _, currentMap := range maps {
		for key, val := range currentMap {
			result[key] = val
		}
	}
	return result
}

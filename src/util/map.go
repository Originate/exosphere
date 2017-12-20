package util

// Merge merges two maps. It overwrites existing keys of map1 if conflicting keys exist
func Merge(map1 map[string]string, map2 map[string]string) {
	for k, v := range map2 {
		map1[k] = v
	}
}

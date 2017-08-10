package util

// JoinStringSlices joins the given slices
func JoinStringSlices(slices ...[]string) []string {
	result := []string{}
	for _, slice := range slices {
		result = append(result, slice...)
	}
	return result
}

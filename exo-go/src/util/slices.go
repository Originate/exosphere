package util

// JoinSlices joins the given slices
func JoinSlices(slices ...[]string) []string {
	result := []string{}
	for _, slice := range slices {
		result = append(result, slice...)
	}
	return result
}

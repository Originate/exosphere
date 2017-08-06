package stringtools

// DoesStringArrayContain returns whether the given string slice
// contains the given string.
func DoesStringArrayContain(strings []string, targetString string) bool {
	for _, element := range strings {
		if element == targetString {
			return true
		}
	}
	return false
}

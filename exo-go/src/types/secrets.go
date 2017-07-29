package types

import (
	"fmt"
	"sort"
	"strings"
)

// Secrets map contains maps from secret keys to values
type Secrets map[string]string

// ToTfString converts map to .tfvars string:
// {a:b, c:d} ->
// a="b"
// c="d"
func (s Secrets) ToTfString() TFString {
	tfvars := []string{}

	for key, value := range s {
		tfvars = append(tfvars, fmt.Sprintf("%s=\"%s\"", key, value))
	}

	sort.Strings(tfvars)
	return TFString(strings.Join(tfvars, "\n"))
}

// HasConflictingKey checks that two secrets map doesn't have clonficting keys
func (s Secrets) HasConflictingKey(map2 Secrets) bool {
	for k := range s {
		if _, hasKey := map2[k]; hasKey {
			return true
		}
	}
	return false
}

// MergeSecrets merges two secret maps
func (s Secrets) MergeSecrets(map2 Secrets) Secrets {
	for k, v := range s {
		map2[k] = v
	}
	return map2
}

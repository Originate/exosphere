package types

import (
	"fmt"
	"sort"
	"strings"
)

type Secrets map[string]string

// converts map to .tfvars string:
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

func (map1 Secrets) HasConflictingKey(map2 Secrets) bool {
	for k := range map1 {
		if _, hasKey := map2[k]; hasKey {
			return true
		}
	}
	return false
}

func (map1 Secrets) MergeSecrets(map2 Secrets) Secrets {
	for k, v := range map1 {
		map2[k] = v
	}
	return map2
}
